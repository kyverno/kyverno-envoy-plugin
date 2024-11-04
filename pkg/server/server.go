package server

import (
	"context"
)

type Server interface {
	Run(context.Context) error
}

type ServerFunc func(context.Context) error

func (f ServerFunc) Run(ctx context.Context) error {
	return f(ctx)
}
