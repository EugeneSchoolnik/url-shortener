package login

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

type Request struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type SuccessResponse struct {
	User  *dto.UserPublic `json:"user"`
	Token string          `json:"token"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

type UserAuthenticator interface {
	Login(email, password string) (*model.User, string, error)
}

func New(log *slog.Logger, userAuthenticator UserAuthenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handler.auth.login"
		log = log.With(slog.String("op", op))

		var req Request
		if err := c.ShouldBind(&req); err != nil {
			log.Info("invalid input", sl.Err(err))
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input"})
			return
		}

		user, token, err := userAuthenticator.Login(req.Email, req.Password)
		if err != nil {
			// no need for logs
			var code int
			switch {
			case errors.Is(err, service.ErrValidation) || errors.Is(err, service.ErrInvalidCredentials):
				code = http.StatusBadRequest
			case errors.Is(err, service.ErrNotFound):
				code = http.StatusNotFound
			default:
				code = http.StatusInternalServerError
			}
			c.JSON(code, ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{User: dto.ToUserPublic(user), Token: token})
	}
}
