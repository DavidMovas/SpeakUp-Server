package services

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/stores"
	"github.com/DavidMovas/SpeakUp-Server/internal/models"
	"github.com/DavidMovas/SpeakUp-Server/internal/models/requests"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	store  *stores.UsersStore
	logger *zap.Logger
}

func NewUsersService(store *stores.UsersStore, logger *zap.Logger) *UsersService {
	return &UsersService{
		store:  store,
		logger: logger,
	}
}

func (s *UsersService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	request.Password = string(passHash)

	user, err := s.store.CreateUser(ctx, request)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return user, nil
}
