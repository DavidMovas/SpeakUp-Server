package handlers

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/services"
	"github.com/DavidMovas/SpeakUp-Server/internal/models/requests"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ v1.UsersServiceServer = (*UsersHandler)(nil)

type UsersHandler struct {
	service *services.UsersService
	logger  *zap.Logger

	v1.UnimplementedUsersServiceServer
}

func NewUsersHandler(service *services.UsersService, logger *zap.Logger) *UsersHandler {
	return &UsersHandler{
		service: service,
		logger:  logger,
	}
}

func (u *UsersHandler) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	req, err := requests.MakeRequest[requests.CreateUserRequest](request)
	if err != nil {
		return nil, err
	}

	user, err := u.service.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	res := &v1.RegisterResponse{
		AccessToken: "token",
		User: &v1.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FullName:  user.FullName,
			AvatarUrl: &user.AvatarURL,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	return res, nil
}

func (u *UsersHandler) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UsersHandler) Logout(ctx context.Context, request *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UsersHandler) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}
