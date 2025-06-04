package jwt

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	jwklib "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/jwk"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/stretchr/testify/assert"
)

func Test_decode_string_string(t *testing.T) {
	tests := []struct {
		name       string
		token      ref.Val
		key        ref.Val
		wantHeader map[string]any
		wantClaims map[string]any
		wantValid  bool
	}{{
		name:  "HS256 - valid",
		token: types.String("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"),
		key:   types.String("secret"),
		wantHeader: map[string]any{
			"alg": "HS256",
			"typ": "JWT",
		},
		wantClaims: map[string]any{
			"exp":  float64(2241081539),
			"nbf":  float64(1514851139),
			"role": "guest",
			"sub":  "YWxpY2U=",
		},
		wantValid: true,
	}, {
		name:  "HS256 - expired",
		token: types.String("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ4NTExNTAsIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.3737qU6BTn8QHhGXxcgZ6EA2hJFY-a00V28F4mD8_98"),
		key:   types.String("secret"),
		wantHeader: map[string]any{
			"alg": "HS256",
			"typ": "JWT",
		},
		wantClaims: map[string]any{
			"exp":  float64(1514851150),
			"nbf":  float64(1514851139),
			"role": "guest",
			"sub":  "YWxpY2U=",
		},
		wantValid: false,
	}, {
		name:  "HS256 - not yet valid",
		token: types.String("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1NTAsIm5iZiI6MjI0MTA4MTUzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.nzWl9VIwNT3RYgF8IlTG_IczpIZrFFsUWqnUeexWC64"),
		key:   types.String("secret"),
		wantHeader: map[string]any{
			"alg": "HS256",
			"typ": "JWT",
		},
		wantClaims: map[string]any{
			"exp":  float64(2241081550),
			"nbf":  float64(2241081539),
			"role": "guest",
			"sub":  "YWxpY2U=",
		},
		wantValid: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				Lib(),
			)
			assert.NoError(t, err)
			impl := impl{env.CELTypeAdapter()}
			out := impl.decode_string_string(tt.token, tt.key)
			got, err := utils.ConvertToNative[Token](out)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantClaims, got.Claims.AsMap())
			assert.Equal(t, tt.wantValid, got.Valid)
		})
	}
}

func Test_decode_string_set(t *testing.T) {
	jwks := `
	{
		"keys": [
			{
				"alg": "ES256",
				"crv": "P-256",
				"kid": "my-key-id",
				"kty": "EC",
				"use": "sig",
				"x": "iTV4PECbWuDaNBMTLmwH0jwBTD3xUXR0S-VWsCYv8Gc",
				"y": "-Cnw8d0XyQztrPZpynrFn8t10lyEb6oWqWcLJWPUB5A"
			}
		]
	}`
	set, err := jwk.Parse([]byte(jwks))
	assert.NoError(t, err)
	tests := []struct {
		name       string
		token      ref.Val
		set        jwk.Set
		wantHeader map[string]any
		wantClaims map[string]any
		wantValid  bool
	}{{
		name:  "ES256 - valid",
		token: types.String("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJwa2kuZXhhbXBsZS5jb20ifQ.ViJTHHv5FuJM9LsRrTpzts6tZkN8deKiu5x49-M8-nq6Rs6ta-Wn8fN_YVLlpZvwhFu_yfxpfUGhBRc33QSSsw"),
		set:   set,
		wantHeader: map[string]any{
			"alg": "ES256",
			"typ": "JWT",
		},
		wantClaims: map[string]any{
			"iss": "pki.example.com",
		},
		wantValid: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				Lib(),
			)
			assert.NoError(t, err)
			impl := impl{env.CELTypeAdapter()}
			set := env.CELTypeAdapter().NativeToValue(jwklib.Set{Set: tt.set})
			out := impl.decode_string_set(tt.token, set)
			got, err := utils.ConvertToNative[Token](out)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantClaims, got.Claims.AsMap())
			assert.Equal(t, tt.wantValid, got.Valid)
		})
	}
}
