// Package controllerruntime provides integration utilities between the
// controller-runtime framework and the SDK's core abstractions.
//
// It defines a generic apiSource that acts as a live, controller-managed
// Source of Kubernetes API objects. This allows the SDK to automatically
// track, reconcile, and serve Kubernetes resources as typed data sources
// within an evaluation engine.
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

// apiSource is a generic implementation of a Kubernetes-backed data source.
// It watches and maintains an in-memory collection of API objects of type API,
// synchronized through a controller-runtime controller.
//
// Each reconcile event updates the internal resource map accordingly,
// allowing Load to return a consistent, snapshot-like view of known objects.
//
// Type parameter:
//
//	API — a concrete type implementing client.Object, such as *v1.ConfigMap or *v1.Pod.
type apiSource[API client.Object] struct {
	client    client.Client  // Kubernetes client used to get objects
	lock      *sync.Mutex    // protects concurrent access to resources
	resources map[string]API // in-memory map of reconciled objects, keyed by namespace/name
}

// NewApiSource creates and registers a new controller-runtime managed source
// for a specific Kubernetes API type.
//
// It sets up a controller that watches resources of type API and maintains
// an in-memory cache of their current state. The resulting apiSource can be
// used as a core.Source for higher-level engine components.
//
// Generic type parameter:
//
//	API — a Kubernetes API object type (must satisfy client.Object)
//
// Parameters:
//
//	mgr      — the controller-runtime manager used to build and register the controller
//	options  — controller options such as concurrency, rate limiting, etc.
//
// Returns:
//
//	*apiSource[API] — a live-updating data source of API objects
//	error           — if controller setup fails
//
// Example:
//
//	type MyResource v1.ConfigMap
//	src, err := controllerruntime.NewApiSource[*v1.ConfigMap](mgr, controller.Options{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Later, call src.Load(ctx) to get the current set of ConfigMaps.
func NewApiSource[API client.Object](mgr ctrl.Manager, options controller.Options) (*apiSource[API], error) {
	provider := newApiSource[API](mgr.GetClient())

	// A zero value of API is needed for controller registration.
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

// Reconcile implements the controller-runtime Reconciler interface.
//
// It ensures that the internal state of apiSource matches the current state
// of the Kubernetes API for the given object key. If the object is deleted,
// it is removed from the internal cache; if created or updated, it is stored.
//
// This keeps apiSource.resources synchronized with cluster state.
//
// Example (simplified flow):
//   - Receive reconcile request for a ConfigMap.
//   - Attempt to Get() the object from the API.
//   - If found: store/update it in resources map.
//   - If deleted: remove it from resources map.
func (r *apiSource[API]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Reconcile triggered", "object", req.String())

	var policy API
	err := r.client.Get(ctx, req.NamespacedName, policy)
	if errors.IsNotFound(err) {
		// Object deleted — remove from cache.
		r.lock.Lock()
		defer r.lock.Unlock()
		delete(r.resources, req.String())
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	// Object created or updated — store in cache.
	r.lock.Lock()
	defer r.lock.Unlock()
	r.resources[req.String()] = policy
	return ctrl.Result{}, nil
}

// Load returns a snapshot of all currently known API objects managed by this source.
//
// It implements the core.Source interface and provides a consistent, sorted view
// of the internal resource map. Sorting ensures deterministic ordering of results.
//
// Returns:
//
//	[]API — list of managed API objects
//	error — always nil unless internal state access fails (unlikely)
//
// Example:
//
//	objs, _ := src.Load(context.Background())
//	for _, obj := range objs {
//	    fmt.Println(obj.GetName())
//	}
func (r *apiSource[API]) Load(ctx context.Context) ([]API, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.resources == nil {
		return nil, nil
	}

	out := make([]API, 0, len(r.resources))

	// Collect and sort keys for deterministic iteration.
	keys := maps.Keys(r.resources)
	slices.Sort(keys)

	for _, key := range keys {
		out = append(out, r.resources[key])
	}
	return out, nil
}

// newApiSource creates a bare apiSource instance with initialized fields.
//
// This helper is primarily used internally by NewApiSource. It initializes
// synchronization primitives and the in-memory resource map.
func newApiSource[API client.Object](client client.Client) *apiSource[API] {
	return &apiSource[API]{
		client:    client,
		lock:      &sync.Mutex{},
		resources: map[string]API{},
	}
}
