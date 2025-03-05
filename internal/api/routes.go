package api

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/handlers/http"
	"github.com/DavidMovas/SpeakUp-Server/internal/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/DavidMovas/SpeakUp-Server/internal/config"
	"github.com/labstack/echo/v4"
)

func RegisterAPI(_ context.Context, e *echo.Echo, handler *http.Handler, _ *trace.TracerProvider, _ *metric.MeterProvider, logger *zap.Logger, _ *config.Config) error {
	api := e.Group("/api")

	api.Use(middlewares.NewLoggingMiddleware(logger))

	api.GET("/health", handler.Health)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return nil
}
