package internal

import (
	"context"
	"fmt"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/handlers/http"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/repository"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/services"
	"github.com/DavidMovas/SpeakUp-Server/utils/clients"
	"github.com/DavidMovas/SpeakUp-Server/utils/metrics"
	"github.com/DavidMovas/SpeakUp-Server/utils/telemetry"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"

	"github.com/DavidMovas/SpeakUp-Server/internal/api"
	"github.com/DavidMovas/SpeakUp-Server/internal/config"
	"github.com/DavidMovas/SpeakUp-Server/internal/log"
	"github.com/DavidMovas/SpeakUp-Server/utils/echox"
	"github.com/DavidMovas/SpeakUp-Server/utils/helpers"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	e         *echo.Echo
	logger    *log.Logger
	telemetry *trace.TracerProvider
	metrics   *metric.MeterProvider
	cfg       *config.Config
	closers   []func() error
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	var closers []func() error
	logger, err := log.NewLogger(cfg.Local, cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("logger create faild: %w", err)
	}

	closers = append(closers, logger.Close)

	telem, err := telemetry.NewTelemetry(cfg.TracerURL, cfg.ServerName, cfg.Version)
	if err != nil {
		logger.Error("Failed to create telem provider", zap.Error(err))
		return nil, err
	}

	promet, err := metrics.NewMetrics(cfg.ServerName, cfg.Version)
	if err != nil {
		logger.Error("Failed to create metrics provider", zap.Error(err))
		return nil, err
	}

	postgres, err := clients.NewPostgresClient(ctx, cfg.PostgresURL, nil)
	if err != nil {
		logger.Error("Failed to create postgres client", zap.Error(err))
		return nil, err
	}

	closers = append(closers, func() error {
		postgres.Close()
		return nil
	})

	redis, err := clients.NewRedisClient(cfg.RedisURL, nil)
	if err != nil {
		logger.Error("Failed to create redis client", zap.Error(err))
		return nil, err
	}

	closers = append(closers, redis.Close)

	repo := repository.NewRepository(postgres, redis, logger.Logger)
	service := services.NewService(repo, logger.Logger)
	handler := http.NewHandler(service, logger.Logger)

	e := echo.New()
	e.HTTPErrorHandler = echox.NewErrorHandler(logger.Logger)
	e.HideBanner = true
	e.HidePort = true

	err = api.RegisterAPI(ctx, e, handler, telem, promet, logger.Logger, cfg)
	if err != nil {
		logger.Warn("register api failed", zap.Error(err))
		return nil, fmt.Errorf("register api: %w", err)
	}

	return &Server{
		e:         e,
		logger:    logger,
		telemetry: telem,
		metrics:   promet,
		cfg:       cfg,
		closers:   closers,
	}, nil
}

func (s *Server) Start() error {
	port := s.cfg.HTTPPort

	s.logger.Info("Starting server...", zap.String("port", port), zap.Time("start_time", time.Now()))

	return s.e.Start(fmt.Sprintf(":%s", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...", zap.Time("stop_time", time.Now()))

	err := s.telemetry.Shutdown(ctx)
	if err != nil {
		s.logger.Warn("Failed to shutdown telemetry server", zap.Error(err))
	}

	err = s.metrics.Shutdown(ctx)
	if err != nil {
		s.logger.Warn("Failed to shutdown metrics server", zap.Error(err))
	}

	err = s.Close()
	if err != nil {
		s.logger.Warn("Failed to close closers", zap.Error(err))
	}

	return s.e.Shutdown(ctx)
}

func (s *Server) Close() error {
	return helpers.WithClosers(s.closers, nil)
}
