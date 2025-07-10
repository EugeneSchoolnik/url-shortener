package redirect

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"
	"url-shortener/internal/model"

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
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		c.Redirect(http.StatusFound, url.Link)
	}
}
