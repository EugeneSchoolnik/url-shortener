package route

import (
	"log/slog"
	"url-shortener/internal/http/handler"
	"url-shortener/internal/http/handler/auth/login"
	"url-shortener/internal/http/handler/auth/me"
	"url-shortener/internal/http/handler/auth/register"
	"url-shortener/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func Auth(router gin.IRouter, log *slog.Logger, deps *handler.Dependencies) {
	r := router.Group("/auth")

	r.GET("/me", middleware.Auth(deps.JwtService), me.New(log, deps.UserService))
	r.POST("/register", register.New(log, deps.AuthService))
	r.POST("/login", login.New(log, deps.AuthService))
}
