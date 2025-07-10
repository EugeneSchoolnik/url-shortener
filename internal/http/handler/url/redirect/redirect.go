package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/model"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UrlGetter interface {
	ByID(id string) (*model.Url, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.create"))

		alias := c.Param("alias")
		if alias == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid alias"})
			return
		}

		url, err := urlGetter.ByID(alias)
		if err != nil {
			// no need for logs
			var code int
			switch {
			case errors.Is(err, service.ErrValidation):
				code = http.StatusBadRequest
			case errors.Is(err, service.ErrUrlNotFound):
				code = http.StatusNotFound
			default:
				code = http.StatusInternalServerError
			}
			c.JSON(code, ErrorResponse{Error: err.Error()})
			return
		}

		c.Redirect(http.StatusFound, url.Link)
	}
}
