package controllerruntime

import (
	"context"
	"slices"
	"sync"

	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

func NewApiSource[API client.Object](mgr ctrl.Manager, options controller.Options) (*apiSource[API], error) {
	provider := newApiSource[API](mgr.GetClient())
	var api API
	builder := ctrl.
		NewControllerManagedBy(mgr).
		For(api).
		WithOptions(options)
	if err := builder.Complete(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

type apiSource[API client.Object] struct {
	client    client.Client
	lock      *sync.Mutex
	resources map[string]API
}

func (r *apiSource[API]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Reconcile...")
	var policy API
	err := r.client.Get(ctx, req.NamespacedName, policy)
	if errors.IsNotFound(err) {
		r.lock.Lock()
		defer r.lock.Unlock()
		delete(r.resources, req.String())
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.resources[req.String()] = policy
	return ctrl.Result{}, nil
}

func (r *apiSource[API]) Load(ctx context.Context) ([]API, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.resources == nil {
		return nil, nil
	}
	out := make([]API, 0, len(r.resources))
	keys := maps.Keys(r.resources)
	slices.Sort(keys)
	for _, key := range keys {
		out = append(out, r.resources[key])
	}
	return out, nil
}

func newApiSource[API client.Object](client client.Client) *apiSource[API] {
	return &apiSource[API]{
		client:    client,
		lock:      &sync.Mutex{},
		resources: map[string]API{},
	}
}
