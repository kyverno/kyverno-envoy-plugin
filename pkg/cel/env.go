package cel

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	jsonimpl "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/impl"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/envoy"
	httpauth "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	httpserver "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http-server"
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
	)
}

func NewEnv(evalMode vpol.EvaluationMode) (*cel.Env, error) {
	base, err := NewBaseEnv()
	if err != nil {
		return nil, err
	}
	// register our libs
	switch evalMode {
	case v1alpha1.EvaluationModeEnvoy:
		base, err = base.Extend(
			envoy.Lib(),
		)
	case v1alpha1.EvaluationModeHTTP:
		base, err = base.Extend(
			httpauth.Lib(),
			httpserver.Lib(),
		)
	default:
		err = fmt.Errorf("invalid evaluation mode passed for env builder")
	}
	if err != nil {
		return nil, err
	}
	// create new cel env
	return base.Extend(
		http.Lib(),
		jwt.Lib(),
		jsoncel.Lib(&jsonimpl.JsonImpl{}),
		resource.Lib(),
		image.Lib(),
		imagedata.Lib(),
	)
}
