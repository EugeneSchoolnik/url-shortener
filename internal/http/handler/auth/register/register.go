package register

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"

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

// @Summary Registers the user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body Request true "Register credentials"
// @Success 201  {object}  SuccessResponse
// @Failure 400  {object}  api.ErrorResponse
// @Failure 409  {object}  api.ErrorResponse
// @Router /auth/register [post]
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
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse{User: dto.ToPublicUser(user), Token: token})
	}
}
