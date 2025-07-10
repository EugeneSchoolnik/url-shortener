package create

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type Request = dto.CreateUrl
type SuccessResponse = *dto.PublicUrl
type ErrorResponse struct {
	Error string `json:"error"`
}

type UrlCreator interface {
	Create(urlDto *dto.CreateUrl, userID string) (*model.Url, error)
}

func New(log *slog.Logger, urlCreator UrlCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handler.url.create"))

		var req Request
		if err := c.ShouldBind(&req); err != nil {
			log.Info("invalid input", sl.Err(err))
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input"})
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authorization error"})
			return
		}

		url, err := urlCreator.Create(&req, userID.(string))
		if err != nil {
			// no need for logs
			var code int
			switch {
			case errors.Is(err, service.ErrValidation):
				code = http.StatusBadRequest
			case errors.Is(err, service.ErrAliasTaken):
				code = http.StatusConflict
			case errors.Is(err, service.ErrRelatedResourceNotFound):
				code = http.StatusUnprocessableEntity
			default:
				code = http.StatusInternalServerError
			}
			c.JSON(code, ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, dto.ToPublicUrl(url))
	}
}
