package handler

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/pipe"
	"google.golang.org/grpc"

	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/service"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
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

func (h *ChatHandler) CreateChat(ctx context.Context, request *v1.CreateChatRequest) (*v1.CreateChatResponse, error) {
	result, err := h.service.CreateChat(ctx, request)
	if err != nil {
		return nil, err
	}

	h.logger.Debug("Chat created", zap.String("chat_id", result.ChatId))

	return result, nil
}

func (h *ChatHandler) Connect(stream grpc.BidiStreamingServer[v1.ChatMessage, v1.ChatMessage]) error {
	panic("implement me")
}

func (h *ChatHandler) GetChatHistory(ctx context.Context, request *v1.GetChatHistoryRequest) (*v1.GetChatHistoryResponse, error) {
	panic("implement me")
}
