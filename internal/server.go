package internal

import (
	"context"
	"fmt"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/hub"
	pipe "github.com/DavidMovas/SpeakUp-Server/internal/shared/pipe"
	"google.golang.org/grpc/reflection"
	"net"
	"time"

	chatHnd "github.com/DavidMovas/SpeakUp-Server/internal/api/chat/handler"
	chatSrv "github.com/DavidMovas/SpeakUp-Server/internal/api/chat/service"
	chatStr "github.com/DavidMovas/SpeakUp-Server/internal/api/chat/store"
	usersHnd "github.com/DavidMovas/SpeakUp-Server/internal/api/users/handler"
	usersSrv "github.com/DavidMovas/SpeakUp-Server/internal/api/users/service"
	usersStr "github.com/DavidMovas/SpeakUp-Server/internal/api/users/store"
	"github.com/DavidMovas/SpeakUp-Server/internal/routes"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/clients"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/echox"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/helpers"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/interceptors"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/jwt"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/metrics"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/telemetry"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/DavidMovas/SpeakUp-Server/internal/config"
	"github.com/DavidMovas/SpeakUp-Server/internal/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	e          *echo.Echo
	grpcServer *grpc.Server
	hubCancel  context.CancelFunc
	logger     *log.Logger
	telemetry  *trace.TracerProvider
	metrics    *metric.MeterProvider
	cfg        *config.Config
	closers    []func() error
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

	jwtService := jwt.NewService(cfg.JWTConfig.Secret, cfg.JWTConfig.AccessExpiration, cfg.JWTConfig.RefreshExpiration)

	usersStore := usersStr.NewUsersStore(postgres, logger.Logger)
	chatStore := chatStr.NewChatsStore(postgres, redis, logger.Logger)

	usersService := usersSrv.NewUsersService(usersStore, jwtService, logger.Logger)
	chatService := chatSrv.NewChatService(chatStore, logger.Logger)

	p := pipe.NewPipe(chatService, usersService)

	hubCtx, hubCancel := context.WithCancel(context.Background())
	h := hub.NewHub(hubCtx, chatStore, logger.Logger)

	usersHandler := usersHnd.NewUsersHandler(usersService, p, logger.Logger)
	chatHandler := chatHnd.NewChatHandler(h, chatService, p, logger.Logger)

	e := echo.New()
	e.HTTPErrorHandler = echox.NewErrorHandler(logger.Logger)
	e.HideBanner = true
	e.HidePort = true

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.NewChainUnaryErrorInterceptor(logger.Logger)))

	reflection.Register(grpcServer)

	err = routes.RegisterHTTPAPI(e, telem, promet, logger.Logger, cfg)
	if err != nil {
		logger.Warn("register api failed", zap.Error(err))
		hubCancel()
		return nil, fmt.Errorf("register api: %w", err)
	}

	err = routes.RegisterGRPCAPI(grpcServer, usersHandler, chatHandler, telem, promet, logger.Logger, cfg)
	if err != nil {
		logger.Warn("register api failed", zap.Error(err))
		hubCancel()
		return nil, fmt.Errorf("register api: %w", err)
	}

	return &Server{
		e:          e,
		grpcServer: grpcServer,
		hubCancel:  hubCancel,
		logger:     logger,
		telemetry:  telem,
		metrics:    promet,
		cfg:        cfg,
		closers:    closers,
	}, nil
}

func (s *Server) Start() error {
	httpPort := s.cfg.HTTPPort
	tcpPort := s.cfg.TCPPort

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.TCPPort))
	if err != nil {
		s.logger.Warn("Failed to create listener", zap.Error(err))
		return fmt.Errorf("create listener: %w", err)
	}

	s.closers = append(s.closers, listener.Close)

	startGroup := errgroup.Group{}

	startGroup.Go(func() error {
		s.logger.Info("Starting HTTP server...", zap.String("port", httpPort))
		err = s.e.Start(fmt.Sprintf(":%s", httpPort))
		if err != nil {
			s.logger.Warn("Failed to start HTTP server", zap.Error(err))
			return fmt.Errorf("start HTTP server: %w", err)
		}

		return nil
	})

	startGroup.Go(func() error {
		s.logger.Info("Starting TCP server...", zap.String("port", tcpPort))
		err = s.grpcServer.Serve(listener)
		if err != nil {
			s.logger.Warn("Failed to start TCP server", zap.Error(err))
			return fmt.Errorf("start TCP server: %w", err)
		}

		return nil
	})

	return startGroup.Wait()
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

	s.hubCancel()

	err = s.Close()
	if err != nil {
		s.logger.Warn("Failed to close closers", zap.Error(err))
	}

	err = s.e.Shutdown(ctx)
	if err != nil {
		s.logger.Warn("Failed to shutdown HTTP server", zap.Error(err))
	}

	s.grpcServer.GracefulStop()
	s.logger.Info("Server shutdown completed")

	return nil
}

func (s *Server) Close() error {
	return helpers.WithClosers(s.closers, nil)
}
