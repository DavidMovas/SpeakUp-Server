package services

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/stores"
	"github.com/DavidMovas/SpeakUp-Server/internal/models"
	"github.com/DavidMovas/SpeakUp-Server/internal/models/requests"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	store      *stores.UsersStore
	jwtService *jwt.Service
	logger     *zap.Logger
}

func NewUsersService(store *stores.UsersStore, jwtService *jwt.Service, logger *zap.Logger) *UsersService {
	return &UsersService{
		store:      store,
		jwtService: jwtService,
		logger:     logger,
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
		return nil, err
	}

	return user, nil
}

func (s *UsersService) GenerateAccessToken(userID string) (string, error) {
	token, err := s.jwtService.GenerateToken(userID)
	if err != nil {
		return "", apperrors.Internal(err)
	}

	return token, nil
}

func (s *UsersService) GetUserByEmail(ctx context.Context, request *requests.GetUserByEmailRequest) (*models.User, error) {
	user, err := s.store.GetUserByEmail(ctx, request)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(request.Password))
	if err != nil {
		return nil, apperrors.BadRequestHidden(err, "invalid password")
	}

	return user.User, nil
}

func (s *UsersService) GetUserByUsername(ctx context.Context, request *requests.GetUserByUsernameRequest) (*models.User, error) {
	user, err := s.store.GetUserByUsername(ctx, request)
	if err != nil {
		return nil, err
	}

	return user, nil
}
