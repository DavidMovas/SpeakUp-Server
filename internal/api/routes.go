package api

import (
	"context"
	"net/http"

	"github.com/DavidMovas/SpeakUp-Server/internal/config"
	"github.com/DavidMovas/SpeakUp-Server/internal/log"
	"github.com/labstack/echo/v4"
)

func RegisterAPI(_ context.Context, e *echo.Echo, _ *log.Logger, _ *config.Config) error {
	api := e.Group("/api")

	api.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return nil
}
