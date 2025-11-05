package run

import (
	"context"
	"fmt"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Command() *cobra.Command {
	var probesAddress string
	var metricsAddress string
	var kubeConfigOverrides clientcmd.ConfigOverrides
	var leaderElection bool
	var leaderElectionID string
	var certFile string
	var keyFile string
	var nestedRequest bool
	command := &cobra.Command{
		Use:   "run",
		Short: "Run authz-server controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// create a rest config
				config, err := loadConfig(kubeConfigOverrides)
				if err != nil {
					return err
				}
				// create a controller manager
				mgr, err := createManager(config, metricsAddress, leaderElection, leaderElectionID)
				if err != nil {
					return fmt.Errorf("failed to construct manager: %w", err)
				}
				// register controller
				if err := setup(mgr); err != nil {
					return fmt.Errorf("failed to setup controller: %w", err)
				}
				// run
				return run(ctx, mgr, probesAddress)
			})
		},
	}
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&metricsAddress, "metrics-address", ":9082", "Address to listen on for metrics")
	command.Flags().BoolVar(&leaderElection, "leader-election", false, "Enable leader election")
	command.Flags().StringVar(&leaderElectionID, "leader-election-id", "", "Leader election ID")
	command.Flags().StringVar(&certFile, "cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&keyFile, "key-file", "", "File containing tls private key")
	command.Flags().BoolVar(&nestedRequest, "nested-request", false, "Expect the requests to validate to be in the body of the original request")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return command
}

func loadConfig(kubeConfigOverrides clientcmd.ConfigOverrides) (*rest.Config, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&kubeConfigOverrides,
	)
	return kubeConfig.ClientConfig()
}

func createManager(config *rest.Config, metricsAddress string, leaderElection bool, leaderElectionID string) (manager.Manager, error) {
	scheme := runtime.NewScheme()
	if err := v1alpha1.Install(scheme); err != nil {
		return nil, err
	}
	return ctrl.NewManager(config, ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddress,
		},
		LeaderElection:   leaderElection,
		LeaderElectionID: leaderElectionID,
	})
}

func setup(mgr manager.Manager) error {
	reconciler := &reconciler{
		client:  mgr.GetClient(),
		servers: map[reconcile.Request]*entry{},
		lock:    &sync.Mutex{},
	}
	return ctrl.
		NewControllerManagedBy(mgr).
		For(&v1alpha1.AuthorizationServer{}).
		Complete(reconciler)
}

func run(ctx context.Context, mgr manager.Manager, probesAddress string) error {
	// track errors
	var probesErr, mgrErr error
	// create a cancellable context
	ctx, cancel := context.WithCancel(ctx)
	// cancel context at the end
	defer cancel()
	// create a wait group
	var group wait.Group
	// wait all tasks in the group are over
	defer group.Wait()
	// start manager
	group.StartWithContext(ctx, func(ctx context.Context) {
		// cancel context at the end
		defer cancel()
		mgrErr = mgr.Start(ctx)
	})
	// warm up
	if !mgr.GetCache().WaitForCacheSync(ctx) {
		defer cancel()
		return fmt.Errorf("failed to wait for cache sync")
	}
	// create http and grpc servers
	probesServer := probes.NewServer(probesAddress)
	// run servers
	group.StartWithContext(ctx, func(ctx context.Context) {
		// probes
		defer cancel()
		probesErr = probesServer.Run(ctx)
	})
	return multierr.Combine(probesErr, mgrErr)
}
