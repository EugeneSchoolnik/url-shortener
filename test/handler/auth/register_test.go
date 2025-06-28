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
	"url-shortener/internal/http/handler"
	"url-shortener/internal/http/handler/auth/register"
	"url-shortener/internal/http/route"
	"url-shortener/internal/service/auth"
	"url-shortener/internal/service/user"
	"url-shortener/internal/testutils/testdb"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
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

	type successType = register.SuccessResponse
	type errorType = register.ErrorResponse

	tests := []struct {
		name      string
		body      string
		wantCode  int
		wantError string
	}{
		{
			name:     "success",
			body:     `{"user":{"email":"example@gmail.com","password":"12345678"}}`,
			wantCode: http.StatusCreated,
		},
		{
			name:      "invalid input",
			body:      `{"Human":{"pochta":"example@gmail.com","pass":"12345678"}}`,
			wantCode:  http.StatusBadRequest,
			wantError: "invalid input",
		},
		{
			name:      "invalid email",
			body:      `{"user":{"email":"invalid.gmail.com","password":"12345678"}}`,
			wantCode:  http.StatusBadRequest,
			wantError: "validation error: field Email is not a valid email",
		},
		{
			name:      "invalid password",
			body:      `{"user":{"email":"example@gmail.com","password":"12345"}}`,
			wantCode:  http.StatusBadRequest,
			wantError: "validation error: field Password is not valid",
		},
		{
			name:      "email taken",
			body:      `{"user":{"email":"example@gmail.com","password":"12345678"}}`,
			wantCode:  http.StatusConflict,
			wantError: "email's already taken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tt.body))
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
				var body errorType
				if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
					t.Error("response body is not error type")
					return
				}
				assert.Equal(t, tt.wantError, body.Error)
			}
		})
	}
}
