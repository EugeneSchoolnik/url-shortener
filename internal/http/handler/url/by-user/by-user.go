package by_user

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type SuccessResponse = []*dto.PublicUrl
type ErrorResponse struct {
	Error string `json:"error"`
}

type UrlsGetter interface {
	ByUserID(id string, limit int, offset int) ([]model.Url, error)
}

func New(log *slog.Logger, urlGetter UrlsGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.create"))

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "16"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "query parameter `limit` is invalid"})
			return
		}
		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "query parameter `offset` is invalid"})
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authorization error"})
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
			c.JSON(code, ErrorResponse{Error: err.Error()})
			return
		}

		publicUrls := make([]*dto.PublicUrl, len(urls))
		for i, v := range urls {
			publicUrls[i] = dto.ToPublicUrl(&v)
		}

		c.JSON(http.StatusOK, publicUrls)
	}
}
