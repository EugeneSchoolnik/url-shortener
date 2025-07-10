package url_test

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
	"url-shortener/internal/http/handler/url/create"
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

func TestCreateHandler(t *testing.T) {
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
	_, token, err := authService.Register(&dto.CreateUser{Email: "example@gmail.com", Password: "12345678"})
	require.NoError(t, err)
	// token with user id that doesn't exist
	invalidToken, err := jwtService.Generate("notfound")
	require.NoError(t, err)

	r := gin.New()
	route.Url(r, r, log, &handler.Dependencies{UrlService: urlService, JwtService: jwtService})

	type successType = create.SuccessResponse
	type errorType = create.ErrorResponse

	tests := []struct {
		name       string
		body       string
		authHeader string
		wantCode   int
		wantError  string
	}{
		{
			name:       "success",
			body:       `{"alias":"g","link":"https://google.com"}`,
			authHeader: "Bearer " + token,
			wantCode:   http.StatusCreated,
		},
		{
			name:       "without alias",
			body:       `{"link":"https://google.com"}`,
			authHeader: "Bearer " + token,
			wantCode:   http.StatusCreated,
		},
		{
			name:       "invalid input",
			body:       `123`,
			authHeader: "Bearer " + token,
			wantCode:   http.StatusBadRequest,
			wantError:  "invalid input",
		},
		{
			name:      "without authorization header",
			body:      `{"link":"https://google.com"}`,
			wantCode:  http.StatusUnauthorized,
			wantError: "invalid authorization",
		},
		{
			name:       "invalid authorization header",
			body:       `{"link":"https://google.com"}`,
			authHeader: "Bearer ",
			wantCode:   http.StatusUnauthorized,
			wantError:  "invalid token",
		},
		{
			name:       "duplicate",
			body:       `{"alias":"g","link":"https://google.com"}`,
			authHeader: "Bearer " + token,
			wantCode:   http.StatusConflict,
			wantError:  "this alias is already taken",
		},
		{
			name:       "invalid link",
			body:       `{"link":"google-website"}`,
			authHeader: "Bearer " + token,
			wantCode:   http.StatusBadRequest,
			wantError:  "field Link is not valid",
		},
		{
			name:       "hacked token with user id that doesn't exist",
			body:       `{"link":"https://google.com"}`,
			authHeader: "Bearer " + invalidToken,
			wantCode:   http.StatusUnprocessableEntity,
			wantError:  "invalid entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.authHeader)

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
