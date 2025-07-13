package create

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"

	"github.com/gin-gonic/gin"
)

type Request = dto.CreateUrl
type SuccessResponse = *dto.PublicUrl

type UrlCreator interface {
	Create(urlDto *dto.CreateUrl, userID string) (*model.Url, error)
}

// @Summary Create a short url
// @Tags url
// @Accept  json
// @Produce  json
// @Param request body Request true "alias is optional"
// @Success 201  {object}  SuccessResponse
// @Failure 400  {object}  api.ErrorResponse
// @Failure 401  {object}  api.ErrorResponse
// @Failure 409  {object}  api.ErrorResponse
// @Failure 422  {object}  api.ErrorResponse
// @Router /url [post]
// @Security Bearer
func New(log *slog.Logger, urlCreator UrlCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.create"))

		var req Request
		if err := c.ShouldBind(&req); err != nil {
			log.Info("invalid input", sl.Err(err))
			c.JSON(http.StatusBadRequest, api.ErrResponse("invalid input"))
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, api.ErrResponse("authorization error"))
			return
		}

		url, err := urlCreator.Create(&req, userID.(string))
		if err != nil {
			// no need for logs
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		c.JSON(http.StatusCreated, dto.ToPublicUrl(url))
	}
}
