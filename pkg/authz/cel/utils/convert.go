package utils

import (
	"reflect"

	"github.com/google/cel-go/common/types/ref"
)

func ConvertToNative[T any](value ref.Val) (T, error) {
	// try to convert value to native type
	response, err := value.ConvertToNative(reflect.TypeFor[T]())
	// if it failed return default value for T and error
	if err != nil {
		var t T
		return t, err
	}
	// return the converted value
	return response.(T), nil
}
