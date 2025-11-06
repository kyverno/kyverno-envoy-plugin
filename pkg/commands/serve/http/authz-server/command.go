package authzserver

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/http"
	httplib "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/listener"
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
	var serverAddress string
	var kubeConfigOverrides clientcmd.ConfigOverrides
	var externalPolicySources []string
	var kubePolicySource bool
	var imagePullSecrets []string
	var allowInsecureRegistry bool
	var controlPlaneAddr string
	var controlPlaneReconnectWait time.Duration
	var controlPlaneMaxDialInterval time.Duration
	var healthCheckInterval time.Duration
	var nestedRequest bool
	var certFile string
	var keyFile string
	command := &cobra.Command{
		Use:   "authz-server",
		Short: "Start the Kyverno Authz Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var probesErr, connErr, httpErr, httpMgrErr error
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
					httpCompiler := vpolcompiler.NewCompiler[dynamic.Interface, *httplib.CheckRequest, *httplib.CheckResponse]()
					extForHTTP, err := getExternalProviders(httpCompiler, nOpts, rOpts, externalPolicySources...)
					if err != nil {
						return err
					}
					httpProvider := sdksources.NewComposite(extForHTTP...)
					// if we have a control plane source
					if controlPlaneAddr != "" {
						httpListener := sources.NewListener()
						clientAddr := os.Getenv("POD_IP")
						if clientAddr == "" {
							panic("can't start auth server, no POD_IP has been passed")
						}
						policyListener := listener.NewPolicyListener(
							controlPlaneAddr,
							clientAddr,
							map[vpol.EvaluationMode]listener.Processor{
								v1alpha1.EvaluationModeHTTP: httpListener,
							},
							controlPlaneReconnectWait,
							controlPlaneMaxDialInterval,
							healthCheckInterval,
						)
						group.StartWithContext(ctx, func(ctx context.Context) {
							for {
								select {
								case <-ctx.Done():
									return
								default:
									if connErr = policyListener.Start(ctx); connErr != nil {
										ctrl.LoggerFrom(ctx).Error(connErr, "error connecting to the control plane, sleeping 10 seconds then retrying")
										time.Sleep(time.Second * 10)
									}
									continue
								}
							}
						})
					}
					// if kube policy source is enabled
					if kubePolicySource {
						// create a controller manager
						scheme := runtime.NewScheme()
						if err := vpol.Install(scheme); err != nil {
							return err
						}
						httpMgr, err := ctrl.NewManager(config, ctrl.Options{
							Scheme: scheme,
							Metrics: metricsserver.Options{
								// TODO
								// BindAddress: metricsAddress,
								BindAddress: "0",
							},
							Cache: cache.Options{
								ByObject: map[client.Object]cache.ByObject{
									&vpol.ValidatingPolicy{}: {
										Field: fields.OneTermEqualSelector("spec.evaluation.mode", string(v1alpha1.EvaluationModeHTTP)),
									},
								},
							},
						})
						if err != nil {
							return fmt.Errorf("failed to construct manager: %w", err)
						}
						httpSource, err := sources.NewKube("http", httpMgr, httpCompiler)
						if err != nil {
							return fmt.Errorf("failed to create http source: %w", err)
						}
						httpProvider = sdksources.NewComposite(httpSource, httpProvider)
						// start managers
						group.StartWithContext(ctx, func(ctx context.Context) {
							// cancel context at the end
							defer cancel()
							httpMgrErr = httpMgr.Start(ctx)
						})
						if !httpMgr.GetCache().WaitForCacheSync(ctx) {
							defer cancel()
							return fmt.Errorf("failed to wait for http cache sync")
						}
					}
					// create http and grpc servers
					probesServer := probes.NewServer(probesAddress)
					httpConfig := http.Config{
						Address:       serverAddress,
						NestedRequest: nestedRequest,
						CertFile:      certFile,
						KeyFile:       keyFile,
					}
					httpAuthServer := http.NewServer(httpConfig, httpProvider, dynclient) // run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// probes
						defer cancel()
						probesErr = probesServer.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						defer cancel()
						httpErr = httpAuthServer.Run(ctx)
					})
					return nil
				}(ctx)
				return multierr.Combine(err, probesErr, httpErr, httpMgrErr)
			})
		},
	}
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&metricsAddress, "metrics-address", ":9082", "Address to listen on for metrics")
	command.Flags().StringArrayVar(&externalPolicySources, "external-policy-source", nil, "External policy sources")
	command.Flags().StringArrayVar(&imagePullSecrets, "image-pull-secret", nil, "Image pull secrets")
	command.Flags().BoolVar(&allowInsecureRegistry, "allow-insecure-registry", false, "Allow insecure registry")
	command.Flags().BoolVar(&kubePolicySource, "kube-policy-source", true, "Enable in-cluster kubernetes policy source")
	command.Flags().StringVar(&serverAddress, "server-address", ":9083", "Address to serve the http authorization server on")
	command.Flags().BoolVar(&nestedRequest, "nested-request", false, "Expect the requests to validate to be in the body of the original request")
	command.Flags().DurationVar(&controlPlaneReconnectWait, "control-plane-reconnect-wait", 3*time.Second, "Duration to wait before retrying connecting to the control plane")
	command.Flags().DurationVar(&controlPlaneMaxDialInterval, "control-plane-max-dial-interval", 8*time.Second, "Duration to wait before stopping attempts of sending a policy to a client")
	command.Flags().DurationVar(&healthCheckInterval, "health-check-interval", 30*time.Second, "Interval for sending health checks")
	command.Flags().StringVar(&controlPlaneAddr, "control-plane-address", "", "Control plane address")
	command.Flags().StringVar(&certFile, "cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&keyFile, "key-file", "", "File containing tls private key")
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
