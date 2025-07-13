package auth_test

import (
	"testing"
	"time"
	"url-shortener/internal/service/auth"

	"github.com/stretchr/testify/require"
)

func TestJWTService_Generate(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "success",
			userID:  "1234",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.NewJWTService("secret", time.Hour)
			got, err := s.Generate(tt.userID)
			t.Log(got)

			if (err != nil) != tt.wantErr {
				t.Errorf("JWTService.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("JWTService.Generate() got is empty")
			}
		})
	}
}

func TestJWTService_Parse(t *testing.T) {
	s := auth.NewJWTService("secret", time.Hour)

	tokenStr, err := s.Generate("1234")
	require.NoError(t, err)

	tests := []struct {
		name     string
		tokenStr string
		want     string
		wantErr  bool
	}{
		{
			name:     "success",
			tokenStr: tokenStr,
			want:     "1234",
			wantErr:  false,
		},
		{
			name:     "invalid token",
			tokenStr: "invalidtoken",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Parse(tt.tokenStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTService.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JWTService.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
