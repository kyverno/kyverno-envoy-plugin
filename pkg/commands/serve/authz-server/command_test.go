package authzserver

import (
	"testing"

	apolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/apol/compiler"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/vpol/compiler"
)

func Test_getExternalProviders(t *testing.T) {
	apolcompiler := apolcompiler.NewCompiler()
	vpolcompiler := vpolcompiler.NewCompiler()
	getExternalProviders(apolcompiler, vpolcompiler, "https://raw.githubusercontent.com/kyverno/kyverno-envoy-plugin/refs/heads/main/README.md/")
	getExternalProviders(apolcompiler, vpolcompiler, "git+https://github.com/kyverno/kyverno-envoy-plugin.git")
}
