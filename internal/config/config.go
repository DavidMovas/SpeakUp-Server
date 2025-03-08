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
	ServerName      string        `env:"SERVER_NAME" envDefault:"api-server"`
	Local           bool          `env:"LOCAL" envDefault:"true"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"info"`
	HTTPPort        string        `env:"HTTP_PORT" envDefault:"8080"`
	TCPPort         string        `env:"TCP_PORT" envDefault:"8081"`
	StartTimeout    time.Duration `env:"START_TIMEOUT" envDefault:"15s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"15s"`
	PostgresURL     string        `env:"POSTGRES_URL"`
	RedisURL        string        `env:"REDIS_URL"`
	JWTConfig       JWTConfig     `envPrefix:"JWT_"`
	TracerURL       string        `env:"TRACER_URL" envDefault:"http://localhost:14268/api/traces"`
	Version         string        `env:"VERSION" envDefault:"0.0.1"`
}

type JWTConfig struct {
	Secret            string        `env:"SECRET" envDefault:"secret"`
	AccessExpiration  time.Duration `env:"ACCESS_EXPIRATION" envDefault:"24h"`
	RefreshExpiration time.Duration `env:"REFRESH_EXPIRATION" envDefault:"168h"`
}
