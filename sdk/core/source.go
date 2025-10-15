package core

import "context"

type Source[DATA any] interface {
	Load(context.Context) ([]DATA, error)
}

type SourceFunc[DATA any] func(context.Context) ([]DATA, error)

func (f SourceFunc[DATA]) Load(ctx context.Context) ([]DATA, error) {
	return f(ctx)
}

func MakeSourceFunc[DATA any](f func(context.Context) ([]DATA, error)) SourceFunc[DATA] {
	return f
}

func MakeSource[DATA any](data ...DATA) SourceFunc[DATA] {
	return MakeSourceFunc(func(ctx context.Context) ([]DATA, error) {
		return data, nil
	})
}
