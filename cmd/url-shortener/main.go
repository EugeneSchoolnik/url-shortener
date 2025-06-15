package main

import (
	"log"
	"url-shortener/internal/config"
	"url-shortener/internal/database"
)

func main() {
	config := config.MustLoad()
	log.Printf("config loaded: %v\n", config)

	db, err := database.New(&config.Postgres)
	if err != nil {
		log.Fatalf("failed to init database: %s", err.Error())
	}
	log.Println("database initialized")

	_ = db
}
