package engine

import (
	"context"
)

type Provider interface {
	CompiledPolicies(context.Context) ([]CompiledPolicy, error)
}
