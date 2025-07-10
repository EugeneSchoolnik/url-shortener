package remove

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UrlDeleter interface {
	Delete(id, userID string) error
}

func New(log *slog.Logger, urlDeleter UrlDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.remove"))

		id := c.Param("id")
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authorization error"})
			return
		}

		err := urlDeleter.Delete(id, userID.(string))
		if err != nil {
			// no need for logs
			var code int
			switch {
			case errors.Is(err, service.ErrValidation):
				code = http.StatusBadRequest
			default:
				code = http.StatusInternalServerError
			}
			c.JSON(code, ErrorResponse{Error: err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
