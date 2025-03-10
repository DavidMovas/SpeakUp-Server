package pipe

import (
	chat "github.com/DavidMovas/SpeakUp-Server/internal/api/chat/service"
	users "github.com/DavidMovas/SpeakUp-Server/internal/api/users/service"
)

type Pipe struct {
	chat  *chat.ChatService
	users *users.UsersService
}

func NewPipe(chat *chat.ChatService, users *users.UsersService) *Pipe {
	return &Pipe{
		chat:  chat,
		users: users,
	}
}

func (p *Pipe) Chat() *chat.ChatService {
	return p.chat
}

func (p *Pipe) Users() *users.UsersService {
	return p.users
}
