package engine

import (
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
)

type EnvoySource = core.Source[EnvoyPolicy]
type HTTPSource = core.Source[HTTPPolicy]
