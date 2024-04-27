package template

import (
	"context"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	plugin "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/template/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template/functions"
	kyvernofunctions "github.com/kyverno/kyverno-json/pkg/engine/template/kyverno"
)

func GetFunctions(ctx context.Context) []jpfunctions.FunctionEntry {
	var funcs []jpfunctions.FunctionEntry
	funcs = append(funcs, jpfunctions.GetDefaultFunctions()...)
	funcs = append(funcs, functions.GetFunctions()...)
	funcs = append(funcs, kyvernofunctions.GetBareFunctions()...)
	funcs = append(funcs, plugin.GetFunctions()...)
	return funcs
}
