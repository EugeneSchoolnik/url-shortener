package route

import (
	"log/slog"
	"url-shortener/internal/http/handler"
	by_user "url-shortener/internal/http/handler/url/by-user"
	"url-shortener/internal/http/handler/url/create"
	"url-shortener/internal/http/handler/url/redirect"
	"url-shortener/internal/http/handler/url/remove"
	"url-shortener/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func Url(root gin.IRouter, router gin.IRouter, log *slog.Logger, deps *handler.Dependencies) {
	r := router.Group("/url", middleware.Auth(deps.JwtService))

	root.GET("/:alias", redirect.New(log, deps.UrlService))
	r.POST("", create.New(log, deps.UrlService))
	r.GET("", by_user.New(log, deps.UrlService))
	r.DELETE(":id", remove.New(log, deps.UrlService))
}
