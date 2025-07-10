package url_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/http/handler"
	by_user "url-shortener/internal/http/handler/url/by-user"
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

func TestByUserHandler(t *testing.T) {
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
	// urls for test
	for i := range 10 {
		_, err := urlService.Create(&dto.CreateUrl{Link: "https://test/" + strconv.Itoa(i)}, user.ID)
		require.NoError(t, err)
	}

	r := gin.New()
	route.Url(r, r, log, &handler.Dependencies{UrlService: urlService, JwtService: jwtService})

	type successType = by_user.SuccessResponse
	type errorType = by_user.ErrorResponse
	type query struct {
		limit  any
		offset any
	}

	tests := []struct {
		name       string
		query      query
		authHeader string
		wantCode   int
		wantError  string
	}{
		{
			name:       "success",
			query:      query{limit: 5, offset: 0},
			authHeader: "Bearer " + token,
			wantCode:   http.StatusOK,
		},
		{
			name:      "without authorization header",
			query:     query{limit: 5, offset: 0},
			wantCode:  http.StatusUnauthorized,
			wantError: "invalid authorization",
		},
		{
			name:       "invalid query limit",
			query:      query{limit: "a", offset: 0},
			authHeader: "Bearer " + token,
			wantCode:   http.StatusBadRequest,
			wantError:  "query parameter `limit` is invalid",
		},
		{
			name:       "invalid query offset",
			query:      query{limit: 5, offset: "a"},
			authHeader: "Bearer " + token,
			wantCode:   http.StatusBadRequest,
			wantError:  "query parameter `offset` is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/url"+fmt.Sprintf("?limit=%v&offset=%v", tt.query.limit, tt.query.offset), nil)
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
				assert.Equal(t, tt.query.limit, len(body))
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
