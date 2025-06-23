package auth_test

import (
	"testing"
	"time"
	"url-shortener/internal/service/auth"
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
	tests := []struct {
		name     string
		tokenStr string
		want     string
		wantErr  bool
	}{
		{
			name:     "success",
			tokenStr: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwiZXhwIjoxNzUwNTgzMzk3LCJpYXQiOjE3NTA1Nzk3OTd9.9vUWCCID7qaawZD2Y10_Tmg1d4iVwK7aXdOi1l559Tc",
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
			s := auth.NewJWTService("secret", time.Hour)

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
