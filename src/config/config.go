package config

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

//go:embed env/env.*
var envFiles embed.FS

const defaultAppEnv = "local"

// Config holds environment-driven application settings.
type Config struct {
	AppEnv  string `env:"APP_ENV" envDefault:"local"`
	AppPort string `env:"APP_PORT" envDefault:"7090"`

	Postgres Postgres
	Redis    Redis
	GraphQL  GraphQL
}

// Load reads src/config/env/env.{APP_ENV} (default: local) and parses into Config.
// Process environment variables already set (e.g. Docker) take precedence.
func Load() (*Config, error) {
	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == "" {
		appEnv = defaultAppEnv
	}

	name := "env/env." + appEnv
	raw, err := envFiles.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("config: unknown APP_ENV %q (expected env file %s): %w", appEnv, name, err)
	}

	vals, err := godotenv.Unmarshal(string(raw))
	if err != nil {
		return nil, fmt.Errorf("config: failed to parse %s: %w", name, err)
	}
	for key, value := range vals {
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
	slog.Info("config: loaded env file", "env", appEnv, "file", name)

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse environment: %w", err)
	}
	return cfg, nil
}
