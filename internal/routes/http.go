package routes

import (
	"github.com/DavidMovas/SpeakUp-Server/internal/config"
	"github.com/DavidMovas/SpeakUp-Server/internal/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func RegisterHTTPAPI(e *echo.Echo, _ *trace.TracerProvider, _ *metric.MeterProvider, logger *zap.Logger, _ *config.Config) error {
	api := e.Group("/api")

	api.Use(middlewares.NewLoggingMiddleware(logger))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return nil
}
