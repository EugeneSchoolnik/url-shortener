package register

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type Request struct {
	User *dto.CreateUser `json:"user" binding:"required"`
}
type SuccessResponse struct {
	User  *dto.PublicUser `json:"user"`
	Token string          `json:"token"`
}

type UserRegistrar interface {
	Register(userDto *dto.CreateUser) (*model.User, string, error)
}

func New(log *slog.Logger, userRegisterer UserRegistrar) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handler.auth.register"
		log = log.With(slog.String("op", op))

		var req Request
		if err := c.ShouldBind(&req); err != nil {
			log.Info("invalid input", sl.Err(err))
			c.JSON(http.StatusBadRequest, api.ErrResponse("invalid input"))
			return
		}

		user, token, err := userRegisterer.Register(req.User)
		if err != nil {
			// no need for logs
			var code int
			switch {
			case errors.Is(err, service.ErrValidation):
				code = http.StatusBadRequest
			case errors.Is(err, service.ErrEmailTaken):
				code = http.StatusConflict
			default:
				code = http.StatusInternalServerError
			}
			c.JSON(code, api.ErrResponse(err.Error()))
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse{User: dto.ToPublicUser(user), Token: token})
	}
}
