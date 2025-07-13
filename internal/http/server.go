package http_server

import (
	"log/slog"
	"url-shortener/internal/http/handler"
	"url-shortener/internal/http/route"

	"github.com/gin-gonic/gin"
)

func NewRouter(log *slog.Logger, deps *handler.Dependencies) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	// middleware

	// routes
	route.Auth(v1, log, deps)
	route.Url(r, v1, log, deps)

	return r
}
