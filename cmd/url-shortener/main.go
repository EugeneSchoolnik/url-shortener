package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/lib/logger/sl"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log.Printf("config loaded: %v\n", cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	db, err := database.New(&cfg.Postgres)
	if err != nil {
		log.Error("failed to init database: %s", sl.Err(err))
		return
	}
	log.Info("database initialized")

	_ = db

	// init http server
	router := NewRouter()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	server := NewServer(&cfg.HTTPServer, router)

	log.Info(fmt.Sprintf("Starting server at %s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server failed: %v", sl.Err(err))
	}
}

func NewServer(cfg *config.HTTPServer, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	return router
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
