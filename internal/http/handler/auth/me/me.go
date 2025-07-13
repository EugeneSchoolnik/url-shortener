package me

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/http/api"
	"url-shortener/internal/model"
	"url-shortener/internal/model/dto"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	User *dto.PublicUser `json:"user"`
}

type UserGetter interface {
	ById(id string, withContext ...bool) (*model.User, error)
}

// @Summary Get user data by token
// @Tags auth
// @Produce  json
// @Success 200  {object}  SuccessResponse
// @Failure 400  {object}  api.ErrorResponse
// @Failure 401  {object}  api.ErrorResponse
// @Failure 404  {object}  api.ErrorResponse
// @Router /auth/me [get]
// @Security Bearer
func New(log *slog.Logger, userGetter UserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handler.auth.me"
		log = log.With(slog.String("op", op))

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, api.ErrResponse("authorization error"))
			return
		}
		user, err := userGetter.ById(userID.(string), true)
		if err != nil {
			c.JSON(api.ErrReponseFromServiceError(err))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{User: dto.ToPublicUser(user)})
	}
}
