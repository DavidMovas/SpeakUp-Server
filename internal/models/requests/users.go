package requests

import (
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/helpers"
	random "github.com/DavidMovas/SpeakUp-Server/internal/utils/models/helpers"
)

var _ requestable[CreateUserRequest, *v1.RegisterRequest] = (*CreateUserRequest)(nil)

type CreateUserRequest struct {
	ID       string
	Email    string
	Password string
	FullName string
	Username string
}

func (r CreateUserRequest) make(req *v1.RegisterRequest) (*CreateUserRequest, error) {
	id := helpers.GenerateID()
	r.ID = id
	r.Username = random.GenerateRandomUsername(req.GetFullName())
	r.Email = req.GetEmail()
	r.FullName = req.GetFullName()
	r.Password = req.GetPassword()

	return &r, nil
}

var _ requestable[GetUserByEmailRequest, *v1.LoginRequest] = (*GetUserByEmailRequest)(nil)

type GetUserByEmailRequest struct {
	Email    string
	Password string
}

func (r GetUserByEmailRequest) make(req *v1.LoginRequest) (*GetUserByEmailRequest, error) {
	r.Email = req.GetEmail()
	r.Password = req.GetPassword()

	return &r, nil
}

var _ requestable[GetUserByUsernameRequest, *v1.GetUserRequest] = (*GetUserByUsernameRequest)(nil)

type GetUserByUsernameRequest struct {
	Username string
}

func (r GetUserByUsernameRequest) make(req *v1.GetUserRequest) (*GetUserByUsernameRequest, error) {
	r.Username = req.GetUsername()

	return &r, nil
}
