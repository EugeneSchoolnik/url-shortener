package redirect

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"

	"github.com/gin-gonic/gin"
)

type LinkGetter interface {
	RedirectLinkByID(id string) (string, error)
}
type ClickRecorder interface {
	Record(urlID string) error
}

// @Summary Redirect
// @Produce  json
// @Param alias path string true "alias for long url"
// @Success 302
// @Failure 404  {object}  api.ErrorResponse
// @Router /{alias} [get]
func New(log *slog.Logger, linkGetter LinkGetter, clickRecorder ClickRecorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.create"))

		alias := c.Param("alias")
		if alias == "" {
			c.JSON(http.StatusBadRequest, api.ErrResponse("invalid alias"))
			return
		}

		link, err := linkGetter.RedirectLinkByID(alias)
		if err != nil {
			// no need for logs
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}
		err = clickRecorder.Record(alias)
		if err != nil {
			// no need for logs
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		c.Redirect(http.StatusFound, link)
	}
}
