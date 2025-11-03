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

// api is a helper generic interface for pointers to API objects.
//
// It ensures that the generic API pointer type implements client.Object,
// which is required for controller-runtime operations.
//
// API_PTR must be a pointer to API (e.g., *v1.ConfigMap).
type api[API any] interface {
	*API
	client.Object
}

// apiSource is a generic implementation of a Kubernetes-backed data source.
//
// It maintains an in-memory map of API objects (resources) that are kept
// up-to-date by a controller-runtime controller watching a specific Kubernetes
// resource type.
//
// Each reconcile event updates this internal map, so Load can return a consistent
// snapshot of all known objects.
//
// Type parameters:
//
//	API     — the concrete Kubernetes API object type (e.g., v1.ConfigMap)
//	API_PTR — a pointer to API (e.g., *v1.ConfigMap)
type apiSource[API any, API_PTR api[API]] struct {
	client    client.Client      // Kubernetes client used to Get objects
	lock      *sync.Mutex        // protects concurrent access to resources map
	resources map[string]API_PTR // in-memory map of reconciled objects keyed by namespace/name
}

// NewApiSource creates and registers a new controller-runtime managed source
// for a specific Kubernetes API type.
//
// The returned apiSource watches and maintains an in-memory cache of the
// specified resource type. It can be used as a core.Source for engine components.
//
// Type parameters:
//
//	API — a Kubernetes API object type (must satisfy client.Object)
//
// Parameters:
//
//	mgr     — controller-runtime manager used to build/register the controller
//	options — controller options like concurrency, rate limiting, etc.
//
// Returns:
//
//	*apiSource[API] — a live-updating source of API objects
//	error           — if controller setup fails
//
// Example:
//
//	type MyResource v1.ConfigMap
//	src, err := controllerruntime.NewApiSource[*v1.ConfigMap](mgr, controller.Options{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	objs, _ := src.Load(ctx) // snapshot of current ConfigMaps
func NewApiSource[API any, API_PTR api[API]](name string, mgr ctrl.Manager, options controller.Options) (*apiSource[API, API_PTR], error) {
	provider := newApiSource[API, API_PTR](mgr.GetClient())

	// Zero-value API for controller registration
	var api API
	var ptr API_PTR = &api

	builder := ctrl.
		NewControllerManagedBy(mgr).
		For(ptr).
		Named(name).
		WithOptions(options)

	if err := builder.Complete(provider); err != nil {
		return nil, err
	}
	return provider, nil
}

// Reconcile implements the controller-runtime Reconciler interface.
//
// It updates the internal resource map to match the current Kubernetes state
// for a given object key. Deleted objects are removed, created/updated objects
// are stored.
//
// Example reconcile flow:
//   - Receive request for a ConfigMap
//   - Attempt to Get() it from the API
//   - If found: store/update in resources map
//   - If deleted: remove from resources map
func (r *apiSource[API, API_PTR]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("Reconcile triggered", "object", req.String())

	var api API
	var ptr API_PTR = &api

	// Attempt to retrieve the object from the API
	err := r.client.Get(ctx, req.NamespacedName, ptr)
	if errors.IsNotFound(err) {
		// Object deleted — remove from cache
		r.lock.Lock()
		defer r.lock.Unlock()
		delete(r.resources, req.String())
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	// Object created or updated — store in cache
	r.lock.Lock()
	defer r.lock.Unlock()
	r.resources[req.String()] = ptr
	return ctrl.Result{}, nil
}

// Load returns a snapshot of all currently known API objects.
//
// Implements the core.Source interface. Returns a sorted slice for deterministic
// ordering.
//
// Returns:
//
//	[]API_PTR — list of managed API objects
//	error     — always nil unless internal state access fails (rare)
//
// Example:
//
//	objs, _ := src.Load(context.Background())
//	for _, obj := range objs {
//	    fmt.Println(obj.GetName())
//	}
func (r *apiSource[API, API_PTR]) Load(ctx context.Context) ([]API_PTR, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.resources == nil {
		return nil, nil
	}

	out := make([]API_PTR, 0, len(r.resources))

	// Collect and sort keys for deterministic iteration
	keys := maps.Keys(r.resources)
	slices.Sort(keys)

	for _, key := range keys {
		out = append(out, r.resources[key])
	}
	return out, nil
}

// newApiSource creates a bare apiSource instance with initialized fields.
//
// Primarily used internally by NewApiSource. Initializes the mutex and
// in-memory resource map.
func newApiSource[API any, API_PTR api[API]](client client.Client) *apiSource[API, API_PTR] {
	return &apiSource[API, API_PTR]{
		client:    client,
		lock:      &sync.Mutex{},
		resources: map[string]API_PTR{},
	}
}
