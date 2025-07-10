package auth_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/http/api"
	"url-shortener/internal/http/handler"
	"url-shortener/internal/http/handler/auth/login"
	"url-shortener/internal/http/route"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service/auth"
	"url-shortener/internal/service/user"
	"url-shortener/internal/testutils/testdb"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	db := testdb.New(t)
	testdb.TruncateTables(t, "users")

	log := slog.Default()

	// services
	userRepo := repo.NewUserRepo(db)
	userService := user.New(userRepo, log)
	jwtService := auth.NewJWTService("secret", time.Hour)
	authService := auth.New(userService, jwtService, log)

	r := gin.New()
	route.Auth(r, log, &handler.Dependencies{AuthService: authService})

	// create test user
	_, _, err := authService.Register(&dto.CreateUser{Email: "example@gmail.com", Password: "12345678"})
	require.NoError(t, err)

	type successType = login.SuccessResponse

	tests := []struct {
		name      string
		body      string
		wantCode  int
		wantError string
	}{
		{
			name:     "success",
			body:     `{"email":"example@gmail.com","password":"12345678"}`,
			wantCode: http.StatusOK,
		},
		{
			name:      "invalid input",
			body:      `{"pochta":"example@gmail.com","pass":"12345678"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "invalid input",
		},
		{
			name:      "invalid email",
			body:      `{"email":"invalid.gmail.com","password":"12345678"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "field Email is not a valid email",
		},
		{
			name:      "invalid password",
			body:      `{"email":"example@gmail.com","password":"12345"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "field Password is not valid",
		},
		{
			name:      "user not found",
			body:      `{"email":"asdfasdf@gmail.com","password":"12345678"}`,
			wantCode:  http.StatusNotFound,
			wantError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

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
