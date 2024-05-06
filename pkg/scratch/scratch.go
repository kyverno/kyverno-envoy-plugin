package scratch

import (
	"context"
	"fmt"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	function "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/template/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

var Caller = func() interpreter.FunctionCaller {
	var funcs []jpfunctions.FunctionEntry
	funcs = append(funcs, template.GetFunctions(context.Background())...)
	funcs = append(funcs, function.GetFunctions()...)
	return interpreter.NewFunctionCaller(funcs...)
}()

func GetUser(authorisation string) (string, error) {
	vm := interpreter.NewInterpreter(nil, nil)
	parser := parsing.NewParser()
	statement := `base64_decode(@)`
	compiled, err := parser.Parse(statement)
	if err != nil {
		return "", err
	}
	out, err := vm.Execute(compiled, authorisation, interpreter.WithFunctionCaller(Caller))
	if err != nil {
		return "", err
	}
	return out.(string), nil
}

func GetFormJWTToken(arguments []any) (map[string]interface{}, error) {
	vm := interpreter.NewInterpreter(nil, nil)
	parser := parsing.NewParser()

	// Construct JMESPath expression with arguments
	arg1 := fmt.Sprintf("'%s'", arguments[0])
	arg2 := fmt.Sprintf("'%s'", arguments[1])
	statement := fmt.Sprintf("jwt_decode(%s, %s)", arg1, arg2)

	compiled, err := parser.Parse(statement)
	if err != nil {
		return nil, fmt.Errorf("error on compiling , %w", err)
	}
	out, err := vm.Execute(compiled, arguments, interpreter.WithFunctionCaller(Caller))
	if err != nil {
		return nil, fmt.Errorf("error on execute , %w", err)
	}
	return out.(map[string]interface{}), nil
}

func GetFormJWTTokenPayload(arguments []any) (map[string]interface{}, error) {
	vm := interpreter.NewInterpreter(nil, nil)
	parser := parsing.NewParser()

	// Construct JMESPath expression with arguments
	arg1 := fmt.Sprintf("'%s'", arguments[0])
	arg2 := fmt.Sprintf("'%s'", arguments[1])
	statement := fmt.Sprintf("jwt_decode(%s, %s).payload", arg1, arg2)

	compiled, err := parser.Parse(statement)
	if err != nil {
		return nil, fmt.Errorf("error on compiling , %w", err)
	}
	out, err := vm.Execute(compiled, arguments, interpreter.WithFunctionCaller(Caller))
	if err != nil {
		return nil, fmt.Errorf("error on execute , %w", err)
	}
	return out.(map[string]interface{}), nil
}
