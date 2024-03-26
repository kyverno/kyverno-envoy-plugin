package scratch

import "testing"

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
