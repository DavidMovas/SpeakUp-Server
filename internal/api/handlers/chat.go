package handlers

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/services"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var _ v1.ChatServiceServer = (*ChatHandler)(nil)

type ChatHandler struct {
	service *services.ChatService
	logger  *zap.Logger

	v1.UnimplementedChatServiceServer
}

func NewChatHandler(service *services.ChatService, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ChatHandler) CreateRoom(ctx context.Context, request *v1.CreateRoomRequest) (*v1.CreateRoomResponse, error) {
	return nil, nil
}

func (h *ChatHandler) JoinRoom(stream grpc.BidiStreamingServer[v1.JoinRoomRequest, v1.JoinRoomResponse]) error {
	return nil
}
