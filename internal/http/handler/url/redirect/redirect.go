package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"
	"url-shortener/internal/model"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type UrlGetter interface {
	ByID(id string) (*model.Url, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.create"))

		alias := c.Param("alias")
		if alias == "" {
			c.JSON(http.StatusBadRequest, api.ErrResponse("invalid alias"))
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
			c.JSON(code, api.ErrResponse(err.Error()))
			return
		}

		c.Redirect(http.StatusFound, url.Link)
	}
}
