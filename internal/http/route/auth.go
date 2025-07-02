package route

import (
	"log/slog"
	"url-shortener/internal/http/handler"
	"url-shortener/internal/http/handler/auth/login"
	"url-shortener/internal/http/handler/auth/register"

	"github.com/gin-gonic/gin"
)

func Auth(router gin.IRouter, log *slog.Logger, deps *handler.Dependencies) {
	r := router.Group("/auth")

	r.POST("/register", register.New(log, deps.AuthService))
	r.POST("/login", login.New(log, deps.AuthService))
}
