package remove

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"

	"github.com/gin-gonic/gin"
)

type UrlDeleter interface {
	Delete(id, userID string) error
}

// @Summary Remove user's short url
// @Tags url
// @Produce  json
// @Param id path int true "short url id"
// @Success 200
// @Failure 401  {object}  api.ErrorResponse
// @Router /url/{id} [delete]
// @Security Bearer
func New(log *slog.Logger, urlDeleter UrlDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.remove"))

		id := c.Param("id")
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, api.ErrResponse("authorization error"))
			return
		}

		err := urlDeleter.Delete(id, userID.(string))
		if err != nil {
			// no need for logs
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		c.Status(http.StatusOK)
	}
}
