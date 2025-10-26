package controlplane

import (
	"context"
	"fmt"
	"time"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/sources"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/stream/sender"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func Command() *cobra.Command {
	var probesAddress string
	var metricsAddress string
	var grpcAddress string
	var grpcNetwork string
	var kubeConfigOverrides clientcmd.ConfigOverrides
	var initialSendPolicyWait time.Duration
	var maxSendPolicyInterval time.Duration
	var maxClientInactiveDuration time.Duration
	var clientFlushInterval time.Duration
	var leaderElection bool
	var leaderElectionID string
	command := &cobra.Command{
		Use:   "control-plane",
		Short: "Start the Kyverno authorizer control plane",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var httpErr, grpcErr, mgrErr error
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
					// create a cancellable context
					ctx, cancel := context.WithCancel(ctx)
					// cancel context at the end
					defer cancel()
					// create a wait group
					var group wait.Group
					// wait all tasks in the group are over
					defer group.Wait()

					s := sender.NewPolicySender(
						ctx,
						initialSendPolicyWait,
						maxSendPolicyInterval,
						clientFlushInterval,
						maxClientInactiveDuration,
					)

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
					// create policy reconciler
					r := sources.NewPolicyReconciler(mgr.GetClient(), s, nil)
					if err := ctrl.NewControllerManagedBy(mgr).For(&v1alpha1.ValidatingPolicy{}).Complete(r); err != nil {
						return fmt.Errorf("failed to register controller to manager: %w", err)
					}
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

					// create http and grpc servers
					http := probes.NewServer(probesAddress)
					// pass the validating policy stream server as the required argument
					grpc := envoy.NewServer(grpcNetwork, grpcAddress, nil, nil, s)

					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						httpErr = http.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						grpcErr = grpc.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// start dead client flush
						s.StartHealthCheckMonitor(ctx)
					})
					return nil
				}(ctx)
				return multierr.Combine(err, httpErr, grpcErr, mgrErr)
			})
		},
	}
	command.Flags().DurationVar(&initialSendPolicyWait, "initial-send-wait", 5*time.Second, "Duration to wait before retrying a send to a client")
	command.Flags().DurationVar(&maxSendPolicyInterval, "max-send-interval", 10*time.Second, "Duration to wait before stopping attempts of sending a policy to a client")
	command.Flags().DurationVar(&clientFlushInterval, "client-flush-interval", 180*time.Second, "Interval for how often to remove dead client connections")
	command.Flags().DurationVar(&maxClientInactiveDuration, "max-client-inactive-duration", 240*time.Second, "Duration to wait before declaring a client as inactive")
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&grpcAddress, "grpc-address", ":9081", "Address to listen on")
	command.Flags().StringVar(&grpcNetwork, "grpc-network", "tcp", "Network to listen on")
	command.Flags().StringVar(&metricsAddress, "metrics-address", ":9082", "Address to listen on for metrics")
	command.Flags().BoolVar(&leaderElection, "leader-election", false, "Enable leader election")
	command.Flags().StringVar(&leaderElectionID, "leader-election-id", "", "Leader election ID")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return command
}
