package sources

import (
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	controllerruntime "github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/controller-runtime"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

func NewKube(mgr ctrl.Manager, compiler compiler.Compiler) (engine.Source, error) {
	options := controller.Options{
		NeedLeaderElection: ptr.To(false),
	}
	apis, err := controllerruntime.NewApiSource[v1alpha1.ValidatingPolicy](mgr, options)
	if err != nil {
		return nil, err
	}
	transform := sources.NewTransformErr(apis, func(in *v1alpha1.ValidatingPolicy) (engine.Policy, error) {
		policy, err := compiler.Compile(in)
		return policy, err.ToAggregate()
	})
	// TODO: cache
	return transform, nil
}
