package jwt

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"github.com/stretchr/testify/assert"
)

func Test_decode(t *testing.T) {
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
			out := impl.decode(tt.token, tt.key)
			got, err := utils.ConvertToNative[Token](out)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantClaims, got.Claims.AsMap())
			assert.Equal(t, tt.wantValid, got.Valid)
		})
	}
}
