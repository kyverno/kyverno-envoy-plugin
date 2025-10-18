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
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/sources"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/httpauth"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/processor"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/utils/ocifs"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	sdksources "github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	vpolv1alpha1 "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/sirupsen/logrus"
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
	var leaderElection bool
	var leaderElectionID string
	var imagePullSecrets []string
	var allowInsecureRegistry bool
	var nestedRequest bool
	var httpAuthAddress string
	command := &cobra.Command{
		Use:   "authz-server",
		Short: "Start the Kyverno Authz Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var probesErr, httpAuthErr, grpcErr, mgrErr error
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
					logger := logrus.New()
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

					envoyCompiler := vpolcompiler.NewCompiler[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]()
					httpCompiler := vpolcompiler.NewCompiler[dynamic.Interface, *http.Request, *http.Response]()

					extForEnvoy, err := getExternalProviders(envoyCompiler, nOpts, rOpts, externalPolicySources...)
					if err != nil {
						return err
					}
					extForHTTP, err := getExternalProviders(httpCompiler, nOpts, rOpts, externalPolicySources...)
					if err != nil {
						return err
					}

					envoyProvider := sdksources.NewComposite(extForEnvoy...)
					httpProvider := sdksources.NewComposite(extForHTTP...)
					// if kube policy source is enabled
					if kubePolicySource {
						// create a controller manager
						scheme := runtime.NewScheme()
						if err := vpolv1alpha1.Install(scheme); err != nil {
							return err
						}
						mgr, err := ctrl.NewManager(config, ctrl.Options{
							Scheme: scheme,
							Metrics: metricsserver.Options{
								BindAddress: metricsAddress,
							},
							Cache: cache.Options{
								ByObject: map[client.Object]cache.ByObject{
									&vpolv1alpha1.ValidatingPolicy{}: {
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
						envoyProcessor := processor.NewPolicyAccessor(envoyCompiler, logger)
						httpProcessor := processor.NewPolicyAccessor(httpCompiler, logger)

						r := sources.NewPolicyReconciler(mgr.GetClient(), nil, []processor.Processor{httpProcessor, envoyProcessor})
						if err := ctrl.NewControllerManagedBy(mgr).For(&vpolv1alpha1.ValidatingPolicy{}).Complete(r); err != nil {
							return fmt.Errorf("failed to register controller to manager: %w", err)
						}
						envoyProvider = sdksources.NewComposite(envoyProcessor, envoyProvider)
						httpProvider = sdksources.NewComposite(httpProcessor, httpProvider)
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
					authorizer := httpauth.NewAuthorizer(dynclient, httpProvider, nestedRequest, logger)
					httpAuthServer := httpauth.NewServer(httpAuthAddress, authorizer)
					grpc := authz.NewServer(grpcNetwork, grpcAddress, envoyProvider, dynclient, nil)
					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// probes
						defer cancel()
						probesErr = http.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// grpc auth server
						defer cancel()
						grpcErr = grpc.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						defer cancel()
						httpAuthErr = httpAuthServer.Run(ctx)
					})
					return nil
				}(ctx)
				return multierr.Combine(err, probesErr, httpAuthErr, grpcErr, mgrErr)
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
	command.Flags().BoolVar(&leaderElection, "leader-election", false, "Enable leader election")
	command.Flags().StringVar(&leaderElectionID, "leader-election-id", "", "Leader election ID")
	command.Flags().StringVar(&httpAuthAddress, "http-auth-server-address", ":9083", "Address to serve the http authorization server on")
	command.Flags().BoolVar(&nestedRequest, "nested-request", false, "Expect the requests to validate to be in the body of the original request")
	return command
}

func getExternalProviders[DATA, IN, OUT any](vpolCompiler engine.Compiler[DATA, IN, OUT], nOpts []name.Option, rOpts []remote.Option, urls ...string) ([]core.Source[policy.Policy[DATA, IN, OUT]], error) {
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	// mux.Add(httpfs.FS)
	// mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)

	// Create a configured ocifs.FS with registry options
	configuredOCIFS := ocifs.ConfigureOCIFS(nOpts, rOpts)
	mux.Add(configuredOCIFS)

	var providers []core.Source[policy.Policy[DATA, IN, OUT]]
	for _, url := range urls {
		fsys, err := mux.Lookup(url)
		if err != nil {
			return nil, err
		}
		providers = append(
			providers,
			sdksources.NewOnce(sources.NewFsProvider(vpolCompiler, fsys)),
		)
	}
	return providers, nil
}
