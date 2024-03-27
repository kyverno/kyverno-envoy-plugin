package functions

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/golang-jwt/jwt"
	"github.com/jmespath-community/go-jmespath/pkg/functions"
)

func GetFunctions() []functions.FunctionEntry {
	return []functions.FunctionEntry{{
		Name: "jwt_decode",
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jwt_decode,
	}}
}

func jwt_decode(arguments []any) (any, error) {

	// Validate argument
	tokenString, err := validateArg(" ", arguments, 0, reflect.String)
	if err != nil {
		return nil, fmt.Errorf("invalidArgumentTypeError: %w", err)
	}
	tokenStringVal := tokenString.String()

	secretkey, err := validateArg(" ", arguments, 1, reflect.String)
	if err != nil {
		return nil, fmt.Errorf("invalidArgumentTypeError: %w", err)
	}

	// Attempt to decode the base64 encoded secret key
	decodedKey, err := base64.StdEncoding.DecodeString(secretkey.String())
	if err != nil {
		// If decoding fails, assume the secret key is not base64 encoded
		decodedKey = []byte(secretkey.String())
	}

	token, err := jwt.Parse(tokenStringVal, func(token *jwt.Token) (interface{}, error) {
		return decodedKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid JWT token: %w", err)
	}

	result := map[string]any{
		"header":  jwt.MapClaims(token.Header),
		"payload": jwt.MapClaims(token.Claims.(jwt.MapClaims)),
		"sig":     fmt.Sprintf("%x", token.Signature),
	}
	return result, nil
}

func validateArg(f string, arguments []any, index int, expectedType reflect.Kind) (reflect.Value, error) {
	if index >= len(arguments) {
		return reflect.Value{}, formatError(argOutOfBoundsError, f, index+1, len(arguments))
	}
	if arguments[index] == nil {
		return reflect.Value{}, formatError(invalidArgumentTypeError, f, index+1, expectedType.String())
	}
	arg := reflect.ValueOf(arguments[index])
	if arg.Type().Kind() != expectedType {
		return reflect.Value{}, formatError(invalidArgumentTypeError, f, index+1, expectedType.String())
	}
	return arg, nil
}
