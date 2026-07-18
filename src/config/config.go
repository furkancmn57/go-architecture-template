package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds environment-driven application settings.
type Config struct {
	AppEnv  string `env:"APP_ENV" envDefault:"local"`
	AppPort string `env:"APP_PORT" envDefault:"8080"`

	Postgres Postgres
	Redis    Redis
	RabbitMQ RabbitMQ
}

// Load reads .env (if present) and parses process environment into Config.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("config: no .env file found, relying on process environment")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse environment: %w", err)
	}
	return cfg, nil
}
