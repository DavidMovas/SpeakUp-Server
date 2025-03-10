package handler

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/pipe"

	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/service"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var _ v1.ChatServiceServer = (*ChatHandler)(nil)

type ChatHandler struct {
	service *service.ChatService
	pipe    *pipe.Pipe
	logger  *zap.Logger

	v1.UnimplementedChatServiceServer
}

func NewChatHandler(service *service.ChatService, pipe *pipe.Pipe, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		service: service,
		pipe:    pipe,
		logger:  logger,
	}
}

func (h *ChatHandler) CreateRoom(ctx context.Context, request *v1.CreateRoomRequest) (*v1.CreateRoomResponse, error) {
	return nil, nil
}

func (h *ChatHandler) JoinRoom(stream grpc.BidiStreamingServer[v1.JoinRoomRequest, v1.JoinRoomResponse]) error {
	return nil
}
