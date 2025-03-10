package service

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/store"
	chat "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
)

type ChatService struct {
	store  *store.ChatsStore
	logger *zap.Logger
}

func NewChatService(store *store.ChatsStore, logger *zap.Logger) *ChatService {
	return &ChatService{
		store:  store,
		logger: logger,
	}
}

func (s *ChatService) CreateRoom(_ context.Context, _ *chat.CreateRoomRequest) (*chat.CreateRoomResponse, error) {
	return nil, nil
}

func (s *ChatService) GetPrivateChatIDBetweenUsers(ctx context.Context, requesterID, searchedID string) (string, error) {
	return s.store.GetPrivateChatIDBetweenUsers(ctx, requesterID, searchedID)
}

func (s *ChatService) GetGroupChatIDBetweenUsers(ctx context.Context, userIDs ...string) (string, error) {
	return s.store.GetGroupChatIDBetweenUsers(ctx, userIDs...)
}
