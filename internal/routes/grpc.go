package routes

import (
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/handler"
	handler2 "github.com/DavidMovas/SpeakUp-Server/internal/api/users/handler"
	"github.com/DavidMovas/SpeakUp-Server/internal/config"
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RegisterGRPCAPI(grpcServer *grpc.Server, usersHandler *handler2.UsersHandler, chatHandler *handler.ChatHandler, _ *trace.TracerProvider, _ *metric.MeterProvider, _ *zap.Logger, _ *config.Config) error {

	v1.RegisterUsersServiceServer(grpcServer, usersHandler)
	v1.RegisterChatServiceServer(grpcServer, chatHandler)

	return nil
}
