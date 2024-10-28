package cel

import (
	"reflect"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/cel-go/interpreter"
	"github.com/stretchr/testify/assert"
)

func TestNewEnv(t *testing.T) {
	source := `
envoy
	.Denied(401)
	.WithBody("Authentication Failed")
	.WithHeader(envoy.Header("foo", "bar").KeepEmptyValue())
	.Response()
	.WithMetadata({"my-new-metadata": "my-new-value"})
`
	env, err := NewEnv()
	assert.NoError(t, err)
	ast, issues := env.Compile(source)
	assert.Nil(t, issues)
	prog, err := env.Program(ast)
	assert.NoError(t, err)
	assert.NotNil(t, prog)
	out, _, err := prog.Eval(interpreter.EmptyActivation())
	assert.NoError(t, err)
	assert.NotNil(t, out)
	a, err := out.ConvertToNative(reflect.TypeFor[*authv3.CheckResponse]())
	assert.NoError(t, err)
	assert.NotNil(t, a)
}
