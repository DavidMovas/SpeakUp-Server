package api

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/middlewares"
	"go.uber.org/zap"
	"net/http"

	"github.com/DavidMovas/SpeakUp-Server/internal/config"
	"github.com/labstack/echo/v4"
)

func RegisterAPI(_ context.Context, e *echo.Echo, logger *zap.Logger, _ *config.Config) error {
	api := e.Group("/api")

	api.Use(middlewares.NewLoggingMiddleware(logger))

	api.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return nil
}
