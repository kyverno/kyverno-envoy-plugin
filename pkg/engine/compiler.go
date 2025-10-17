package engine

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Compiler[T any] interface {
	Compile(T) (Policy, field.ErrorList)
}
