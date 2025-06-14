package main

import (
	"log"
	"url-shortener/internal/config"
)

func main() {
	config := config.MustLoad()
	log.Printf("config loaded: %v\n", config)
}
