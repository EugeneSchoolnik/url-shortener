package main

import (
	"log"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/lib/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()
	log.Printf("config loaded: %v\n", config)

	log := setupLogger(config.Env)
	log.Info("starting url-shortener", slog.String("env", config.Env))
	log.Debug("debug messages are enabled")

	db, err := database.New(&config.Postgres)
	if err != nil {
		log.Error("failed to init database: %s", sl.Err(err))
	}
	log.Info("database initialized")

	_ = db
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
