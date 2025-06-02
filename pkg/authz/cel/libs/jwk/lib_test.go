package jwk

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	env, err := cel.NewEnv(
		Lib(),
	)
	assert.NoError(t, err)
	ast, issues := env.Compile("jwks.Fetch('https://www.googleapis.com/oauth2/v3/certs')")
	assert.NoError(t, issues.Err())
	prog, err := env.Program(ast)
	assert.NoError(t, err)
	jwks, _, err := prog.Eval(map[string]any{})
	assert.NoError(t, err)
	assert.NotNil(t, jwks.Value())
}
