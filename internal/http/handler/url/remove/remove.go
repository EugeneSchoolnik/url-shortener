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
