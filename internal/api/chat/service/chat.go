package service

import (
	"context"
	"errors"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/models/requests"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/store"
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/model"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
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

func (s *ChatService) CreateChat(ctx context.Context, request *v1.CreateChatRequest) (*v1.CreateChatResponse, error) {

	switch request.Payload.(type) {
	case *v1.CreateChatRequest_PrivateChat_:
		req, err := model.MakeRequest[requests.CreatePrivateChatRequest](request.GetPrivateChat())
		if err != nil {
			return nil, err
		}

		chatID, err := s.store.CreatePrivateChat(ctx, req)
		if err != nil {
			return nil, err
		}

		return &v1.CreateChatResponse{ChatId: chatID}, nil
	case *v1.CreateChatRequest_GroupChat_:
		req, err := model.MakeRequest[requests.CreateGroupChatRequest](request.GetGroupChat())
		if err != nil {
			return nil, err
		}

		chatID, err := s.store.CreateGroupChat(ctx, req)
		if err != nil {
			return nil, err
		}

		return &v1.CreateChatResponse{ChatId: chatID}, nil
	}

	return nil, apperrors.Internal(errors.New("invalid Create Chat payload"))
}

func (s *ChatService) GetPrivateChatIDBetweenUsers(ctx context.Context, requesterID, searchedID string) (string, error) {
	return s.store.GetPrivateChatIDBetweenUsers(ctx, requesterID, searchedID)
}

func (s *ChatService) GetGroupChatIDBetweenUsers(ctx context.Context, userIDs ...string) (string, error) {
	return s.store.GetGroupChatIDBetweenUsers(ctx, userIDs...)
}
