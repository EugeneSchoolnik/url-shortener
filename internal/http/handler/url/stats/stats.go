package stats

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/http/api"

	"github.com/gin-gonic/gin"
)

type SuccessResponse = []repo.DailyCount

type StatsGetter interface {
	Stats(urlID string, userID string) ([]repo.DailyCount, error)
}

// @Summary Get user's url stats
// @Tags url
// @Produce  json
// @Param id path int true "short url id"
// @Success 200  {object}  SuccessResponse
// @Failure 404  {object}  api.ErrorResponse
// @Router /url/{id} [get]
// @Security Bearer
func New(log *slog.Logger, statsGetter StatsGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.stats"))

		urlID := c.Param("id")
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, api.ErrResponse("authorization error"))
			return
		}

		stats, err := statsGetter.Stats(urlID, userID.(string))
		if err != nil {
			// no need for logs
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}
