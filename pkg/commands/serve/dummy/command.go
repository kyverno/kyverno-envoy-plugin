package dummy

import (
	"context"
	"fmt"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func Command() *cobra.Command {
	var probesAddress string
	var metricsAddress string
	var kubeConfigOverrides clientcmd.ConfigOverrides
	var leaderElection bool
	var leaderElectionID string
	command := &cobra.Command{
		Use:   "dummy",
		Short: "Start the Kyverno Authz Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var probesErr, mgrErr error
				err := func(ctx context.Context) error {
					// create a rest config
					kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
						clientcmd.NewDefaultClientConfigLoadingRules(),
						&kubeConfigOverrides,
					)
					config, err := kubeConfig.ClientConfig()
					if err != nil {
						return err
					}
					// create a controller manager
					scheme := runtime.NewScheme()
					if err := v1alpha1.Install(scheme); err != nil {
						return err
					}
					mgr, err := ctrl.NewManager(config, ctrl.Options{
						Scheme: scheme,
						Metrics: metricsserver.Options{
							BindAddress: metricsAddress,
						},
						LeaderElection:   leaderElection,
						LeaderElectionID: leaderElectionID,
					})
					if err != nil {
						return fmt.Errorf("failed to construct manager: %w", err)
					}
					builder := ctrl.
						NewControllerManagedBy(mgr).
						For(&v1alpha1.AuthorizationServer{})
					// .
					// WithOptions(options)
					// TODO: lock map
					servers := map[reconcile.Request]*server{}
					reconciler := reconcile.Func(func(_ context.Context, req reconcile.Request) (reconcile.Result, error) {
						var object v1alpha1.AuthorizationServer
						err := mgr.GetClient().Get(ctx, req.NamespacedName, &object)
						if errors.IsNotFound(err) {
							// stop server and remove
							srv := servers[req]
							if srv != nil {
								srv.Stop()
								delete(servers, req)
							}
							return ctrl.Result{}, nil
						}
						if err != nil {
							return ctrl.Result{}, err
						}
						srv := servers[req]
						if srv == nil {
							// create server
							fmt.Printf("CREATED Name: %s, Namespace: %s", req.Name, req.Namespace)
							ctx, cancel := context.WithCancel(ctx)
							srv = &server{
								grpcNetwork: "tcp",
								grpcAddress: fmt.Sprintf(":%d", object.Spec.Type.Envoy.Port),
								cancel:      cancel,
							}
							srv.Start(ctx)
							servers[req] = srv
						}
						// configure server
						fmt.Printf("UPDATED Name: %s, Namespace: %s", req.Name, req.Namespace)
						return ctrl.Result{}, nil // TODO
					})
					if err := builder.Complete(reconciler); err != nil {
						return fmt.Errorf("failed to construct controller: %w", err)
					}
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
					if !mgr.GetCache().WaitForCacheSync(ctx) {
						defer cancel()
						return fmt.Errorf("failed to wait for cache sync")
					}
					// create probes servers
					probesServer := probes.NewServer(probesAddress)
					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// probes
						defer cancel()
						probesErr = probesServer.Run(ctx)
					})
					return nil
				}(ctx)
				return multierr.Combine(err, probesErr, mgrErr)
			})
		},
	}
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&metricsAddress, "metrics-address", ":9082", "Address to listen on for metrics")
	command.Flags().BoolVar(&leaderElection, "leader-election", false, "Enable leader election")
	command.Flags().StringVar(&leaderElectionID, "leader-election-id", "", "Leader election ID")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return command
}
