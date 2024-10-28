package utils

import (
	"reflect"

	"github.com/google/cel-go/common/types/ref"
)

func ConvertToNative[T any](value ref.Val) (T, error) {
	response, err := value.ConvertToNative(reflect.TypeFor[T]())
	if err != nil {
		var t T
		return t, err
	}
	return response.(T), nil
}
