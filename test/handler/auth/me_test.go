package auth_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/http/api"
	"url-shortener/internal/http/handler"
	"url-shortener/internal/http/handler/auth/me"
	"url-shortener/internal/http/route"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service/auth"
	"url-shortener/internal/service/user"
	"url-shortener/internal/testutils/testdb"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeHandler(t *testing.T) {
	db := testdb.New(t)
	testdb.TruncateTables(t, "users")

	log := slog.Default()

	// services
	userRepo := repo.NewUserRepo(db)
	userService := user.New(userRepo, log)
	jwtService := auth.NewJWTService("secret", time.Hour)
	authService := auth.New(userService, jwtService, log)

	r := gin.New()
	route.Auth(r, log, &handler.Dependencies{AuthService: authService, UserService: userService, JwtService: jwtService})

	// create test user
	user, token, err := authService.Register(&dto.CreateUser{Email: "example@gmail.com", Password: "12345678"})
	require.NoError(t, err)
	invalidToken, err := jwtService.Generate("notfound")
	require.NoError(t, err)

	type successType = me.SuccessResponse

	tests := []struct {
		name      string
		token     string
		wantCode  int
		wantError string
	}{
		{
			name:     "success",
			token:    token,
			wantCode: http.StatusOK,
		},
		{
			name:      "without token",
			wantCode:  http.StatusUnauthorized,
			wantError: "invalid token",
		},
		{
			name:      "invalid token with user id that doesn't exist",
			token:     invalidToken,
			wantCode:  http.StatusNotFound,
			wantError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, tt.wantCode, res.Code)

			// success
			if tt.wantError == "" {
				var body successType
				if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
					t.Error("response body is not success type")
					return
				}
				assert.Equal(t, dto.ToPublicUser(user), body.User)
			} else {
				// error
				var body api.ErrorResponse
				if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
					t.Error("response body is not error type")
					return
				}
				assert.Equal(t, tt.wantError, body.Error)
			}
		})
	}
}
