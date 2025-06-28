package register

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"
	"url-shortener/internal/service"
	userService "url-shortener/internal/service/user"

	"github.com/gin-gonic/gin"
)

type Request struct {
	User *dto.CreateUser `json:"user" binding:"required"`
}
type SuccessResponse struct {
	User  *dto.UserPublic `json:"user"`
	Token string          `json:"token"`
}
type ErrorResponse struct {
	Error string `json:"error"`
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
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input"})
			return
		}

		user, token, err := userRegisterer.Register(req.User)
		if err != nil {
			// no need for logs
			var code int
			switch {
			case errors.Is(err, service.ErrValidation):
				code = http.StatusBadRequest
			case errors.Is(err, userService.ErrEmailTaken):
				code = http.StatusConflict
			default:
				code = http.StatusInternalServerError
			}
			c.JSON(code, ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse{User: dto.ToUserPublic(user), Token: token})
	}
}
