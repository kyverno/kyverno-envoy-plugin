package authzserver

import (
	"context"
	"fmt"
	"log"
	"os"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/sources"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/utils/ocifs"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	sdksources "github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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
	var imagePullSecrets []string
	var allowInsecureRegistry bool
	command := &cobra.Command{
		Use:   "authz-server",
		Short: "Start the Kyverno Authz Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var probesErr, grpcErr, envoyMgrErr error
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
					secrets := make([]string, 0)
					if len(imagePullSecrets) > 0 {
						secrets = append(secrets, imagePullSecrets...)
					}
					dynclient, err := dynamic.NewForConfig(config)
					if err != nil {
						return err
					}
					// Create kubernetes client
					kubeclient, err := kubernetes.NewForConfig(config)
					if err != nil {
						return err
					}
					namespace, _, err := kubeConfig.Namespace()
					if err != nil {
						return fmt.Errorf("failed to get namespace from kubeconfig: %w", err)
					}
					if namespace == "" || namespace == "default" {
						// Log a warning or require explicit namespace setting
						log.Printf("Using namespace '%s' - consider setting explicit namespace", namespace)
					}
					rOpts, nOpts, err := ocifs.RegistryOpts(kubeclient.CoreV1().Secrets(namespace), allowInsecureRegistry, secrets...)
					if err != nil {
						log.Fatalf("failed to initialize registry opts: %v", err)
						os.Exit(1)
					}
					// initialize compiler
					envoyCompiler := vpolcompiler.NewCompiler[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]()
					extForEnvoy, err := getExternalProviders(envoyCompiler, nOpts, rOpts, externalPolicySources...)
					if err != nil {
						return err
					}
					envoyProvider := sdksources.NewComposite(extForEnvoy...)
					// if kube policy source is enabled
					if kubePolicySource {
						// create a controller manager
						scheme := runtime.NewScheme()
						if err := vpol.Install(scheme); err != nil {
							return err
						}
						envoyMgr, err := ctrl.NewManager(config, ctrl.Options{
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
						})
						if err != nil {
							return fmt.Errorf("failed to construct manager: %w", err)
						}
						envoySource, err := sources.NewKube("envoy", envoyMgr, envoyCompiler)
						if err != nil {
							return fmt.Errorf("failed to create envoy source: %w", err)
						}
						envoyProvider = sdksources.NewComposite(envoySource, envoyProvider)
						if err != nil {
							return fmt.Errorf("failed to construct manager: %w", err)
						}
						// start managers
						group.StartWithContext(ctx, func(ctx context.Context) {
							// cancel context at the end
							defer cancel()
							envoyMgrErr = envoyMgr.Start(ctx)
						})
						if !envoyMgr.GetCache().WaitForCacheSync(ctx) {
							defer cancel()
							return fmt.Errorf("failed to wait for envoy cache sync")
						}
					}
					// create http and grpc servers
					probesServer := probes.NewServer(probesAddress)
					grpc := envoy.NewServer(grpcNetwork, grpcAddress, envoyProvider, dynclient)
					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// probes
						defer cancel()
						probesErr = probesServer.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// grpc auth server
						defer cancel()
						grpcErr = grpc.Run(ctx)
					})
					return nil
				}(ctx)
				return multierr.Combine(err, probesErr, grpcErr, envoyMgrErr)
			})
		},
	}
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&grpcAddress, "grpc-address", ":9081", "Address to listen on")
	command.Flags().StringVar(&grpcNetwork, "grpc-network", "tcp", "Network to listen on")
	command.Flags().StringVar(&metricsAddress, "metrics-address", ":9082", "Address to listen on for metrics")
	command.Flags().StringArrayVar(&externalPolicySources, "external-policy-source", nil, "External policy sources")
	command.Flags().StringArrayVar(&imagePullSecrets, "image-pull-secret", nil, "Image pull secrets")
	command.Flags().BoolVar(&allowInsecureRegistry, "allow-insecure-registry", false, "Allow insecure registry")
	command.Flags().BoolVar(&kubePolicySource, "kube-policy-source", true, "Enable in-cluster kubernetes policy source")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))

	return command
}

func getExternalProviders[POLICY any](vpolCompiler engine.Compiler[POLICY], nOpts []name.Option, rOpts []remote.Option, urls ...string) ([]core.Source[POLICY], error) {
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	// mux.Add(httpfs.FS)
	// mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)

	// Create a configured ocifs.FS with registry options
	configuredOCIFS := ocifs.ConfigureOCIFS(nOpts, rOpts)
	mux.Add(configuredOCIFS)

	var providers []core.Source[POLICY]
	for _, url := range urls {
		fsys, err := mux.Lookup(url)
		if err != nil {
			return nil, err
		}
		providers = append(
			providers,
			sdksources.NewOnce(sources.NewFs(fsys, vpolCompiler)),
		)
	}
	return providers, nil
}
