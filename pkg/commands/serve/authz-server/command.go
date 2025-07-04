package authzserver

import (
	"context"
	"fmt"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/blobfs"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/policy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
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
	var externalPolicySources []string
	command := &cobra.Command{
		Use:   "authz-server",
		Short: "Start the Kyverno Authz Server",
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
					// create a wait group
					var group wait.Group
					// wait all tasks in the group are over
					defer group.Wait()
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
					})
					if err != nil {
						return fmt.Errorf("failed to construct manager: %w", err)
					}
					// create compiler
					compiler := policy.NewCompiler()
					// create kube provider
					kubeProvider, err := policy.NewKubeProvider(mgr, compiler)
					if err != nil {
						return err
					}
					// create external providers
					externalProvider, err := getExternalProviders(compiler, externalPolicySources...)
					if err != nil {
						return err
					}
					// create final provider
					provider := policy.NewComposite(
						kubeProvider,
						policy.NewComposite(externalProvider...),
					)
					// create a cancellable context
					ctx, cancel := context.WithCancel(ctx)
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
					grpc := authz.NewServer(grpcNetwork, grpcAddress, provider)
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
					return nil
				}(ctx)
				return multierr.Combine(err, httpErr, grpcErr, mgrErr)
			})
		},
	}
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&grpcAddress, "grpc-address", ":9081", "Address to listen on")
	command.Flags().StringVar(&grpcNetwork, "grpc-network", "tcp", "Network to listen on")
	command.Flags().StringVar(&metricsAddress, "metrics-address", ":9082", "Address to listen on for metrics")
	command.Flags().StringArrayVar(&externalPolicySources, "external-policy-source", nil, "External policy sources")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return command
}

func getExternalProviders(compiler policy.Compiler, urls ...string) ([]policy.Provider, error) {
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	mux.Add(httpfs.FS)
	mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)
	var providers []policy.Provider
	for _, url := range urls {
		fsys, err := mux.Lookup(url)
		if err != nil {
			return nil, err
		}
		providers = append(providers, policy.NewFsProvider(compiler, fsys))
	}
	return providers, nil
}
