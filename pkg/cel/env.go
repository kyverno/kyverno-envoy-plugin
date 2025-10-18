package cel

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	jsonimpl "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/impl"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/envoy"
	httpauth "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	jsoncel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/json"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/jwt"

	"github.com/kyverno/kyverno/pkg/cel/libs/http"
	"github.com/kyverno/kyverno/pkg/cel/libs/resource"

	"k8s.io/apiserver/pkg/cel/library"
)

func NewEnv() (*cel.Env, error) {
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
		// register our libs
		envoy.Lib(),
		jwt.Lib(),
		jsoncel.Lib(&jsonimpl.JsonImpl{}),
		// register kyverno libs
		http.Lib(),
		httpauth.Lib(),
		resource.Lib(),
	)
}
