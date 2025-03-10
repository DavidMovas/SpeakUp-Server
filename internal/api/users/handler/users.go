package handler

import (
	"context"

	"github.com/DavidMovas/SpeakUp-Server/internal/api/users/service"
	"github.com/DavidMovas/SpeakUp-Server/internal/models/requests"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ v1.UsersServiceServer = (*UsersHandler)(nil)

type UsersHandler struct {
	service *service.UsersService
	logger  *zap.Logger

	v1.UnimplementedUsersServiceServer
}

func NewUsersHandler(service *service.UsersService, logger *zap.Logger) *UsersHandler {
	return &UsersHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UsersHandler) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	req, err := requests.MakeRequest[requests.CreateUserRequest](request)
	if err != nil {
		return nil, err
	}

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	token, err := h.service.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	res := &v1.RegisterResponse{
		AccessToken: token,
		User: &v1.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FullName:  user.FullName,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	return res, nil
}

func (h *UsersHandler) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	req, err := requests.MakeRequest[requests.GetUserByEmailRequest](request)
	if err != nil {
		return nil, err
	}

	user, err := h.service.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	token, err := h.service.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	res := &v1.LoginResponse{
		AccessToken: token,
		User: &v1.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FullName:  user.FullName,
			AvatarUrl: user.AvatarURL,
			Bio:       user.Bio,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	if user.LastLoginAt != nil {
		res.User.LastLoginAt = timestamppb.New(*user.LastLoginAt)
	}

	if user.UpdatedAt != nil {
		res.User.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return res, nil
}

func (h *UsersHandler) Logout(_ context.Context, _ *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	panic("implement me")
}

func (h *UsersHandler) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	req, err := requests.MakeRequest[requests.GetUserByUsernameRequest](request)
	if err != nil {
		return nil, err
	}

	user, err := h.service.GetUserByUsername(ctx, req)
	if err != nil {
		return nil, err
	}

	res := &v1.GetUserResponse{
		User: &v1.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FullName:  user.FullName,
			AvatarUrl: user.AvatarURL,
			Bio:       user.Bio,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	if user.LastLoginAt != nil {
		res.User.LastLoginAt = timestamppb.New(*user.LastLoginAt)
	}

	return res, nil
}
