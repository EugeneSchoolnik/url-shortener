package url_test

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
	"url-shortener/internal/http/handler/url/stats"
	"url-shortener/internal/http/route"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service/auth"
	clickstat "url-shortener/internal/service/click-stat"
	"url-shortener/internal/service/url"
	"url-shortener/internal/service/user"
	"url-shortener/internal/testutils/testdb"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatsHandler(t *testing.T) {
	db := testdb.New(t)
	testdb.TruncateTables(t, "users")

	log := slog.Default()

	// services
	userRepo := repo.NewUserRepo(db)
	urlRepo := repo.NewUrlRepo(db)
	clickStatRepo := repo.NewClickStatRepo(db)
	userService := user.New(userRepo, log)
	jwtService := auth.NewJWTService("secret", time.Hour)
	authService := auth.New(userService, jwtService, log)
	urlService := url.New(urlRepo, log)
	clickStatService := clickstat.New(clickStatRepo, log)

	// test user
	user, token, err := authService.Register(&dto.CreateUser{Email: "example@gmail.com", Password: "12345678"})
	require.NoError(t, err)
	// urls for test
	testUrl, err := urlService.Create(&dto.CreateUrl{Link: "https://google.com"}, user.ID)
	require.NoError(t, err)

	r := gin.New()
	route.Url(r, r, log, &handler.Dependencies{UrlService: urlService, JwtService: jwtService, ClickStatService: clickStatService})

	type successType = stats.SuccessResponse

	tests := []struct {
		name           string
		alias          string
		simulateClicks int64
		wantCode       int
		wantError      string
	}{
		{
			name:           "success123 redirects",
			alias:          testUrl.ID,
			simulateClicks: 123,
			wantCode:       http.StatusOK,
		},
		{
			name:      "alias doesn't exist",
			alias:     "notfound",
			wantCode:  http.StatusNotFound,
			wantError: "url statistics not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range tt.simulateClicks {
				req := httptest.NewRequest(http.MethodGet, "/"+tt.alias, nil)
				req.Header.Set("Authorization", "Bearer "+token)
				res := httptest.NewRecorder()
				r.ServeHTTP(res, req)
				// require.Equal(t, http.StatusFound, res.Code)
			}

			req := httptest.NewRequest(http.MethodGet, "/url/"+tt.alias, nil)
			req.Header.Set("Authorization", "Bearer "+token)
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
				t.Log("=================", body)
				assert.Equal(t, tt.simulateClicks, body[len(body)-1].Count)
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
