package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &c, nil
}

type Config struct {
	Local           bool          `env:"LOCAL"`
	LogLevel        string        `env:"LOG_LEVEL"`
	HTTPPort        string        `env:"HTTP_PORT"`
	StartTimeout    time.Duration `env:"START_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}
