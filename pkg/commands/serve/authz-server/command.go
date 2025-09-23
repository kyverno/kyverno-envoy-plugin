package authzserver

import (
	"context"
	"fmt"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	apolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/apol/compiler"
	apolprovider "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/apol/provider"
	genericproviders "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/providers"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/vpol/compiler"
	vpolprovider "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/vpol/provider"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func Command() *cobra.Command {
	var probesAddress string
	var metricsAddress string
	var grpcAddress string
	var grpcNetwork string
	var kubeConfigOverrides clientcmd.ConfigOverrides
	var externalPolicySources []string
	var kubePolicySource bool
	var leaderElection bool
	var leaderElectionID string
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
					// create a cancellable context
					ctx, cancel := context.WithCancel(ctx)
					// cancel context at the end
					defer cancel()
					// create a wait group
					var group wait.Group
					// wait all tasks in the group are over
					defer group.Wait()
					// create compilers
					apolCompiler := apolcompiler.NewCompiler()
					vpolCompiler := vpolcompiler.NewCompiler()
					// create external providers
					externalProviders, err := getExternalProviders(apolCompiler, vpolCompiler, externalPolicySources...)
					if err != nil {
						return err
					}
					provider := genericproviders.NewComposite(externalProviders...)
					// if kube policy source is enabled
					if kubePolicySource {
						// create a controller manager
						scheme := runtime.NewScheme()
						if err := v1alpha1.Install(scheme); err != nil {
							return err
						}
						if err := vpol.Install(scheme); err != nil {
							return err
						}
						mgr, err := ctrl.NewManager(config, ctrl.Options{
							Scheme: scheme,
							Metrics: metricsserver.Options{
								BindAddress: metricsAddress,
							},
							Cache: cache.Options{
								ByObject: map[client.Object]cache.ByObject{
									&vpol.ValidatingPolicy{}: {
										Field: fields.OneTermEqualSelector("spec.evaluation.mode", string(v1alpha1.EvaluationModeEnvoy)),
									},
								},
							},
							LeaderElection:   leaderElection,
							LeaderElectionID: leaderElectionID,
						})
						if err != nil {
							return fmt.Errorf("failed to construct manager: %w", err)
						}
						// create kube providers
						apolProvider, err := apolprovider.NewKubeProvider(mgr, apolCompiler)
						if err != nil {
							return err
						}
						vpolProvider, err := vpolprovider.NewKubeProvider(mgr, vpolCompiler)
						if err != nil {
							return err
						}
						// create final provider
						provider = genericproviders.NewComposite(apolProvider, vpolProvider, provider)
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
	command.Flags().BoolVar(&kubePolicySource, "kube-policy-source", true, "Enable in-cluster kubernetes policy source")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	command.Flags().BoolVar(&leaderElection, "leader-election", true, "Enable leader election")
	command.Flags().StringVar(&leaderElectionID, "leader-election-id", "", "Leader election ID")
	return command
}

func getExternalProviders(apolCompiler apolcompiler.Compiler, vpolCompiler vpolcompiler.Compiler, urls ...string) ([]engine.Provider, error) {
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	// mux.Add(httpfs.FS)
	// mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)
	var providers []engine.Provider
	for _, url := range urls {
		fsys, err := mux.Lookup(url)
		if err != nil {
			return nil, err
		}
		providers = append(providers, genericproviders.NewFsProvider(apolCompiler, vpolCompiler, fsys))
	}
	return providers, nil
}
