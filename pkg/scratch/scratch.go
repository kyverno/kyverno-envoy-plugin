package scratch

import (
	"context"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

var Caller = func() interpreter.FunctionCaller {
	var funcs []jpfunctions.FunctionEntry
	funcs = append(funcs, template.GetFunctions(context.Background())...)
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
