package jwk

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	// {
	// 	jwks := `
	// 	{
	// 		"keys": [
	// 			{
	// 				"alg": "ES256",
	// 				"crv": "P-256",
	// 				"kid": "my-key-id",
	// 				"kty": "EC",
	// 				"use": "sig",
	// 				"x": "iTV4PECbWuDaNBMTLmwH0jwBTD3xUXR0S-VWsCYv8Gc",
	// 				"y": "-Cnw8d0XyQztrPZpynrFn8t10lyEb6oWqWcLJWPUB5A"
	// 			}
	// 		]
	// 	}`
	// 	set, err := jwk.Parse([]byte(jwks))
	// 	assert.NoError(t, err)
	// 	token := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJwa2kuZXhhbXBsZS5jb20ifQ.ViJTHHv5FuJM9LsRrTpzts6tZkN8deKiu5x49-M8-nq6Rs6ta-Wn8fN_YVLlpZvwhFu_yfxpfUGhBRc33QSSsw"
	// 	tok, err := jwt.Parse([]byte(token), jwt.WithKeySet(set, jws.WithUseDefault(true) /*, jws.WithRequireKid(false)*/))
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, tok)
	// }
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
