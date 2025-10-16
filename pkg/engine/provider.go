package engine

import (
	"context"
)

type Source interface {
	Load(context.Context) ([]CompiledPolicy, error)
}
