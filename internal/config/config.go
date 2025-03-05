package config

type Config struct {
	Local    bool   `env:"LOCAL"`
	LogLevel string `env:"LOG_LEVEL"`
	HTTPPort string `env:"HTTP_PORT"`
}
