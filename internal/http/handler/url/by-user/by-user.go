package by_user

import (
	"log/slog"
	"net/http"
	"strconv"
	"url-shortener/internal/http/api"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"

	"github.com/gin-gonic/gin"
)

type SuccessResponse = []*dto.PublicUrl

type UrlsGetter interface {
	ByUserID(id string, limit int, offset int) ([]model.Url, error)
}

// @Summary Get user's short urls
// @Tags url
// @Accept  json
// @Produce  json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200  {object}  SuccessResponse
// @Failure 400  {object}  api.ErrorResponse
// @Failure 401  {object}  api.ErrorResponse
// @Router /url [get]
// @Security Bearer
func New(log *slog.Logger, urlGetter UrlsGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.byUser"))

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "16"))
		if err != nil {
			c.JSON(http.StatusBadRequest, api.ErrResponse("query parameter `limit` is invalid"))
			return
		}
		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, api.ErrResponse("query parameter `offset` is invalid"))
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, api.ErrResponse("authorization error"))
			return
		}

		urls, err := urlGetter.ByUserID(userID.(string), limit, offset)
		if err != nil {
			// no need for logs
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		publicUrls := make([]*dto.PublicUrl, len(urls))
		for i, v := range urls {
			publicUrls[i] = dto.ToPublicUrl(&v)
		}

		c.JSON(http.StatusOK, publicUrls)
	}
}
