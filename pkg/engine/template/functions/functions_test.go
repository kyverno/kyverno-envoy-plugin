package functions

import (
	"fmt"
	"reflect"

	"testing"
)

func Test_jwt_decode(t *testing.T) {

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk"
	secret := "c2VjcmV0"
	type args struct {
		arguments []any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Positive case , function returns what we expected",
			args: args{[]any{token, secret}},
			want: map[string]interface{}{
				"header": map[string]interface{}{
					"alg": "HS256",
					"typ": "JWT",
				},
				"payload": map[string]interface{}{
					"exp":  2.241081539e+09,
					"nbf":  1.514851139e+09,
					"role": "guest",
					"sub":  "YWxpY2U=",
				},
				"sig": fmt.Sprintf("%x", []byte{0x6a, 0x61, 0x31, 0x62, 0x67, 0x76, 0x49, 0x74, 0x34, 0x37, 0x33, 0x39, 0x33, 0x62, 0x61, 0x5f, 0x57, 0x62, 0x53, 0x42, 0x6d, 0x33, 0x35, 0x4e, 0x72, 0x55, 0x68, 0x64, 0x78, 0x4d, 0x34, 0x6d, 0x4f, 0x56, 0x51, 0x4e, 0x38, 0x69, 0x58, 0x7a, 0x38, 0x6c, 0x6b}),
			},
			wantErr: false,
		},
		// Negative test case: passing incorrect arguments (invalid token)
		{
			name: "negative case - invalid token",
			args: args{[]any{"invalid_jwt_token", secret}},
			want: map[string]interface{}{
				"header":  nil,
				"payload": nil,
				"sig":     nil,
			},
			wantErr: true,
		},
		// Negative test case: passing incorrect arguments (invalid secret)
		{
			name: "negative case - invalid secret",
			args: args{[]any{token, "invalid_secret"}},
			want: map[string]interface{}{
				"header":  nil,
				"payload": nil,
				"sig":     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jwt_decode(tt.args.arguments)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwt_decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotValue := reflect.ValueOf(got)
				wantValue := reflect.ValueOf(tt.want)

				if !reflect.DeepEqual(gotValue.Interface(), wantValue.Interface()) {
					t.Errorf("jwt_decode() = %v, want %v", gotValue.Interface(), wantValue.Interface())
				}
			}
		})
	}
}
