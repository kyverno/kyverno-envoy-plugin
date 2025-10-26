package cel

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	jsonimpl "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/impl"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/envoy"
	httpauth "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	jsoncel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/json"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/jwt"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/pkg/cel/libs/http"
	"github.com/kyverno/kyverno/pkg/cel/libs/image"
	"github.com/kyverno/kyverno/pkg/cel/libs/imagedata"
	"github.com/kyverno/kyverno/pkg/cel/libs/resource"
	"k8s.io/apiserver/pkg/cel/library"
)

func NewBaseEnv() (*cel.Env, error) {
	// create new cel env
	return cel.NewEnv(
		// configure env
		cel.HomogeneousAggregateLiterals(),
		cel.EagerlyValidateDeclarations(true),
		cel.DefaultUTCTimeZone(true),
		cel.CrossTypeNumericComparisons(true),
		// register common libs
		cel.OptionalTypes(),
		ext.Bindings(),
		ext.Encoders(),
		ext.Lists(),
		ext.Math(),
		ext.Protos(),
		ext.Sets(),
		ext.Strings(),
		// register kubernetes libs
		library.CIDR(),
		library.Format(),
		library.IP(),
		library.Lists(),
		library.Regex(),
		library.URLs(),
		library.Quantity(),
		library.SemverLib(),
		// register our libs
		jwt.Lib(),
		jsoncel.Lib(&jsonimpl.JsonImpl{}),
		http.Lib(),
		resource.Lib(),
		image.Lib(),
		imagedata.Lib(),
	)
}

func NewEnvoyEnv() (*cel.Env, error) {
	base, err := NewBaseEnv()
	if err != nil {
		return nil, err
	}
	return base.Extend(envoy.Lib())
}

func NewHttpEnv() (*cel.Env, error) {
	base, err := NewBaseEnv()
	if err != nil {
		return nil, err
	}
	return base.Extend(httpauth.Lib())
}

func NewEnv(evalMode vpol.EvaluationMode) (*cel.Env, error) {
	switch evalMode {
	case v1alpha1.EvaluationModeEnvoy:
		return NewEnvoyEnv()
	case v1alpha1.EvaluationModeHTTP:
		return NewHttpEnv()
	default:
		return nil, fmt.Errorf("invalid evaluation mode passed for env builder")
	}
}
