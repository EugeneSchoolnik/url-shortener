package by_user

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"url-shortener/internal/http/api"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type SuccessResponse = []*dto.PublicUrl

type UrlsGetter interface {
	ByUserID(id string, limit int, offset int) ([]model.Url, error)
}

func New(log *slog.Logger, urlGetter UrlsGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.create"))

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
			var code int
			switch {
			case errors.Is(err, service.ErrValidation):
				code = http.StatusBadRequest
			default:
				code = http.StatusInternalServerError
			}
			c.JSON(code, api.ErrResponse(err.Error()))
			return
		}

		publicUrls := make([]*dto.PublicUrl, len(urls))
		for i, v := range urls {
			publicUrls[i] = dto.ToPublicUrl(&v)
		}

		c.JSON(http.StatusOK, publicUrls)
	}
}
