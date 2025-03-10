package handler

import (
	"context"
	"errors"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/hub"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/service"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/pipe"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var _ v1.ChatServiceServer = (*ChatHandler)(nil)

type ChatHandler struct {
	hub     *hub.Hub
	service *service.ChatService
	pipe    *pipe.Pipe
	logger  *zap.Logger

	v1.UnimplementedChatServiceServer
}

func NewChatHandler(hub *hub.Hub, service *service.ChatService, pipe *pipe.Pipe, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		hub:     hub,
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

func (h *ChatHandler) Connect(stream grpc.BidiStreamingServer[v1.ConnectRequest, v1.ConnectResponse]) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	if _, ok := msg.Payload.(*v1.ConnectRequest_JointChat); ok {
		h.hub.RegisterStream(msg.GetJointChat().GetUserId(), stream)
	}

	return apperrors.Internal(errors.New("invalid steam payload"))
}

func (h *ChatHandler) GetChatHistory(_ context.Context, _ *v1.GetChatHistoryRequest) (*v1.GetChatHistoryResponse, error) {
	panic("implement me")
}
