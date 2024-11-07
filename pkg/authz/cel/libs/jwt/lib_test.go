package jwt

// import (
// 	"testing"

// 	"github.com/golang-jwt/jwt"
// 	"github.com/google/cel-go/common/types"
// 	"github.com/google/cel-go/common/types/ref"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_decode(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		token ref.Val
// 		key   ref.Val
// 		want  jwt.Token
// 	}{{
// 		token: types.String("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"),
// 		key:   types.String("secret"),
// 		want: jwt.Token{
// 			Header: map[string]any{
// 				"alg": "HS256",
// 				"typ": "JWT",
// 			},
// 			Claims: jwt.MapClaims{
// 				"exp":  float64(2241081539),
// 				"nbf":  float64(1514851139),
// 				"role": "guest",
// 				"sub":  "YWxpY2U=",
// 			},
// 			Valid: true,
// 		},
// 	}}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			out := decode(tt.token, tt.key)
// 			got, ok := out.(Token)
// 			assert.True(t, ok)
// 			assert.Equal(t, tt.want.Header, got.Header)
// 			assert.Equal(t, tt.want.Claims, got.Claims)
// 			assert.Equal(t, tt.want.Valid, got.Valid)
// 		})
// 	}
// }
