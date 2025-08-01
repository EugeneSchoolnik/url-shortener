package url_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/http/api"
	"url-shortener/internal/http/handler"
	"url-shortener/internal/http/route"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service/auth"
	"url-shortener/internal/service/url"
	"url-shortener/internal/service/user"
	"url-shortener/internal/testutils/testdb"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveHandler(t *testing.T) {
	db := testdb.New(t)
	testdb.TruncateTables(t, "users", "urls")

	log := slog.Default()

	// services
	userRepo := repo.NewUserRepo(db)
	urlRepo := repo.NewUrlRepo(db)
	userService := user.New(userRepo, log)
	jwtService := auth.NewJWTService("secret", time.Hour)
	authService := auth.New(userService, jwtService, log)
	urlService := url.New(urlRepo, log)

	// test user
	user, token, err := authService.Register(&dto.CreateUser{Email: "example@gmail.com", Password: "12345678"})
	require.NoError(t, err)
	// url for test
	url, err := urlService.Create(&dto.CreateUrl{Link: "https://test/123"}, user.ID)
	require.NoError(t, err)

	r := gin.New()
	route.Url(r, r, log, &handler.Dependencies{UrlService: urlService, JwtService: jwtService})

	tests := []struct {
		name       string
		id         string
		authHeader string
		wantCode   int
		wantError  string
	}{
		{
			name:       "success",
			id:         url.ID,
			authHeader: "Bearer " + token,
			wantCode:   http.StatusOK,
		},
		{
			name:      "without authorization header",
			id:        url.ID,
			wantCode:  http.StatusUnauthorized,
			wantError: "invalid authorization",
		},
		{
			name:       "delete what does not exist",
			id:         url.ID,
			authHeader: "Bearer " + token,
			wantCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/url/%s", tt.id), nil)
			req.Header.Set("Authorization", tt.authHeader)

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			assert.Equal(t, tt.wantCode, res.Code)

			if tt.wantError != "" {
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
