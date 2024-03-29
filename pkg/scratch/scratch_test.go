package scratch

import (
	"reflect"
	"testing"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name          string
		authorisation string
		want          string
		wantErr       bool
	}{
		{
			authorisation: "YWxpY2U6cGFzc3dvcmQ=",
			want:          "alice:password",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.authorisation)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFormJWTToken(t *testing.T) {

	type args struct {
		arguments []any
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			args: args{[]any{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk", "c2VjcmV0"}},
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
				"sig": "6a61316267764974343733393362615f576253426d33354e72556864784d346d4f56514e3869587a386c6b",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFormJWTToken(tt.args.arguments)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFormJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFormJWTToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
