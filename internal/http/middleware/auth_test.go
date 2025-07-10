package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/internal/http/middleware"
	"url-shortener/internal/service/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	tests := []struct {
		name   string
		header func(jwt *auth.JWTService) string
		want   string
		err    string
	}{
		{
			name: "success",
			header: func(jwt *auth.JWTService) string {
				token, err := jwt.Generate("1234")
				require.NoError(t, err)
				t.Log(token)
				return "Bearer " + token
			},
			want: "1234",
		},
		{
			name: "empty token",
			header: func(jwt *auth.JWTService) string {
				return "Bearer "
			},
			err: "invalid token",
		},
		{
			name: "empty header",
			header: func(jwt *auth.JWTService) string {
				return ""
			},
			err: "invalid authorization",
		},
		{
			name: "expired token",
			header: func(jwt *auth.JWTService) string {
				return "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0IiwiZXhwIjoxNzUxNDU1MzY5LCJpYXQiOjE3NTE0NTUzNjh9.bHxyiz4LoLm4MzblZkUAOttSp5ZUqQiTsN-9oMkvv6U"
			},
			err: "invalid token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtService := auth.NewJWTService("secret", time.Hour)
			router := gin.New()
			router.Use(middleware.Auth(jwtService))

			router.GET("/protected", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				c.JSON(http.StatusOK, gin.H{"user_id": userID})
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tt.header(jwtService))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if tt.err == "" {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			}
			assert.Contains(t, w.Body.String(), tt.want)
			assert.Contains(t, w.Body.String(), tt.err)
		})
	}
}
