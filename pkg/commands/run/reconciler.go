package run

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/envoy"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/http"
	httplib "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/sources"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/utils/ocifs"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	sdksources "github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/errors"
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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	envoyCompiler = vpolcompiler.NewCompiler[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]()
	httpCompiler  = vpolcompiler.NewCompiler[dynamic.Interface, *httplib.CheckRequest, *httplib.CheckResponse]()
)

type entry struct {
	cancel func() error
}

type reconciler struct {
	client        client.Client
	servers       map[reconcile.Request]*entry
	certFile      string
	keyFile       string
	nestedRequest bool
	lock          *sync.Mutex
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var object v1alpha1.AuthorizationServer
	err := r.client.Get(ctx, req.NamespacedName, &object)
	if errors.IsNotFound(err) {
		// stop server and remove
		err := r.stopServer(req)
		return ctrl.Result{}, err
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, r.runServer(req, object)
}

func (r *reconciler) runServer(req ctrl.Request, object v1alpha1.AuthorizationServer) error {
	var envoyErr, httpErr error
	if object.Spec.Type.Envoy != nil {
		envoyErr = r.runEnvoyServer(req, object)
	}
	if object.Spec.Type.Envoy != nil {
		httpErr = r.runHttpServer(req, object)
	}
	return multierr.Combine(envoyErr, httpErr)
}

func (r *reconciler) runEnvoyServer(req ctrl.Request, object v1alpha1.AuthorizationServer) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	server := r.servers[req]
	if server == nil {
		// create server
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			nil,
		)
		config, err := kubeConfig.ClientConfig()
		if err != nil {
			return err
		}
		dynclient, err := dynamic.NewForConfig(config)
		if err != nil {
			return err
		}
		// create a controller manager
		scheme := runtime.NewScheme()
		if err := vpol.Install(scheme); err != nil {
			return err
		}
		mgr, err := ctrl.NewManager(config, ctrl.Options{
			Scheme: scheme,
			Metrics: metricsserver.Options{
				BindAddress: "0",
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
		// create a cancellable context
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		// create a wait group
		var group wait.Group
		// track errors
		var grpcErr, mgrErr error
		group.StartWithContext(ctx, func(ctx context.Context) {
			// cancel context at the end
			defer cancel()
			mgrErr = mgr.Start(ctx)
		})
		if !mgr.GetCache().WaitForCacheSync(ctx) {
			defer cancel()
			return fmt.Errorf("failed to wait for cache sync")
		}
		src, err := buildSources(mgr, envoyCompiler, object)
		if err != nil {
			return fmt.Errorf("failed to build engine source: %w", err)
		}
		grpc := envoy.NewServer(object.Spec.Type.Envoy.Network, object.Spec.Type.Envoy.Address, src, dynclient)
		group.StartWithContext(ctx, func(ctx context.Context) {
			// grpc auth server
			defer cancel()
			grpcErr = grpc.Run(ctx)
		})
		server = &entry{
			cancel: func() error {
				// cancel context
				cancel()
				// wait all tasks in the group are over
				group.Wait()
				return multierr.Combine(grpcErr, mgrErr)
			},
		}
		r.servers[req] = server
		return nil
	}
	// configure server
	// TODO
	return nil
}

func (r *reconciler) runHttpServer(req ctrl.Request, object v1alpha1.AuthorizationServer) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	server := r.servers[req]
	if server == nil {
		// create server
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			nil,
		)
		config, err := kubeConfig.ClientConfig()
		if err != nil {
			return err
		}
		dynclient, err := dynamic.NewForConfig(config)
		if err != nil {
			return err
		}
		// create a controller manager
		scheme := runtime.NewScheme()
		if err := vpol.Install(scheme); err != nil {
			return err
		}
		mgr, err := ctrl.NewManager(config, ctrl.Options{
			Scheme: scheme,
			Metrics: metricsserver.Options{
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
		// create a cancellable context
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		// create a wait group
		var group wait.Group
		// track errors
		var grpcErr, mgrErr error
		group.StartWithContext(ctx, func(ctx context.Context) {
			// cancel context at the end
			defer cancel()
			mgrErr = mgr.Start(ctx)
		})
		if !mgr.GetCache().WaitForCacheSync(ctx) {
			defer cancel()
			return fmt.Errorf("failed to wait for cache sync")
		}
		src, err := buildSources(mgr, httpCompiler, object)
		if err != nil {
			return fmt.Errorf("failed to build engine source: %w", err)
		}
		httpConfig := http.Config{
			Address:          object.Spec.Type.HTTP.Address,
			NestedRequest:    r.nestedRequest,
			InputExpression:  object.Spec.Type.HTTP.Modifiers.Request,
			OutputExpression: object.Spec.Type.HTTP.Modifiers.Response,
			CertFile:         r.certFile,
			KeyFile:          r.keyFile,
		}
		http := http.NewServer(httpConfig, src, dynclient)
		group.StartWithContext(ctx, func(ctx context.Context) {
			// grpc auth server
			defer cancel()
			grpcErr = http.Run(ctx)
		})
		server = &entry{
			cancel: func() error {
				// cancel context
				cancel()
				// wait all tasks in the group are over
				group.Wait()
				return multierr.Combine(grpcErr, mgrErr)
			},
		}
		r.servers[req] = server
		return nil
	}
	// configure server
	// TODO
	return nil
}

func (r *reconciler) stopServer(req ctrl.Request) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	server := r.servers[req]
	if server != nil {
		defer delete(r.servers, req)
		return server.cancel()
	}
	return nil
}

func buildSources[POLICY any](mgr ctrl.Manager, compiler engine.Compiler[POLICY], server v1alpha1.AuthorizationServer) (core.Source[POLICY], error) {
	var out []core.Source[POLICY]
	for _, src := range server.Spec.Sources {
		if src.Kubernetes != nil {
			// TODO: selector
			source, err := sources.NewKube(server.Name, mgr, compiler)
			if err != nil {
				return nil, fmt.Errorf("failed to create kube source: %w", err)
			}
			out = append(out, source)
		}
		if src.Git != nil {
			u, err := url.Parse("file://" + src.Fs.Path)
			if err != nil {
				return nil, err
			}
			fsys, err := gitfs.New(u)
			if err != nil {
				return nil, err
			}
			out = append(out, sdksources.NewOnce(sources.NewFs(fsys, compiler)))
		}
		if src.Fs != nil {
			u, err := url.Parse("file://" + src.Fs.Path)
			if err != nil {
				return nil, err
			}
			fsys, err := filefs.New(u)
			if err != nil {
				return nil, err
			}
			out = append(out, sdksources.NewOnce(sources.NewFs(fsys, compiler)))
		}
		if src.Oci != nil {
			u, err := url.Parse(src.Oci.URL)
			if err != nil {
				return nil, err
			}
			kubernetes, err := kubernetes.NewForConfig(mgr.GetConfig())
			if err != nil {
				return nil, err
			}
			rOpts, nOpts, err := ocifs.RegistryOpts(
				kubernetes.CoreV1().Secrets(server.Namespace),
				src.Oci.AllowInsecureRegistry,
				src.Oci.ImagePullSecrets...,
			)
			if err != nil {
				return nil, err
			}
			fsys, err := ocifs.New(u, nOpts, rOpts)
			if err != nil {
				return nil, err
			}
			out = append(out, sdksources.NewOnce(sources.NewFs(fsys, compiler)))
		}
	}
	return sdksources.NewComposite(out...), nil
}
