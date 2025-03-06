package requests

import (
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/helpers"
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
	r.Username = id
	r.Email = req.Email
	r.FullName = req.FullName
	r.Password = req.Password

	return &r, nil
}
