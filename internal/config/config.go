package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	Postgres   Postgres   `yaml:"postgres"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type Postgres struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SLLMode  string `yaml:"ssl_mode" env-default:"disable"`
}

type HTTPServer struct {
	Host         string        `yaml:"host" env-default:"localhost"`
	Port         string        `yaml:"port" env-default:"8080"`
	Timeout      time.Duration `yaml:"timeout" env-required:"true"`
	IddleTimeout time.Duration `yaml:"iddle_timeout" env-required:"true"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not defined")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	return &cfg
}
