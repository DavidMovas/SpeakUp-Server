package services

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/stores"
	chat "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
)

type ChatService struct {
	store  *stores.ChatsStore
	logger *zap.Logger
}

func NewChatService(repo *stores.ChatsStore, logger *zap.Logger) *ChatService {
	return &ChatService{
		store:  repo,
		logger: logger,
	}
}

func (h *ChatService) CreateRoom(ctx context.Context, request *chat.CreateRoomRequest) (*chat.CreateRoomResponse, error) {
	return nil, nil
}
