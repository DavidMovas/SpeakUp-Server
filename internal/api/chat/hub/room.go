package hub

import (
	"context"
	"sync"

	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/models"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/store"
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type room struct {
	sync.RWMutex
	id            string
	clients       map[string]*client
	onlineClients uint
	messageChan   chan *models.Message
	store         *store.ChatsStore

	removeFromHubFunc func(id string)
}

func newRoom(id string, store *store.ChatsStore, remFn func(id string)) *room {
	var r room
	r.id = id
	r.store = store
	r.clients = make(map[string]*client)
	r.messageChan = make(chan *models.Message, 50)
	r.removeFromHubFunc = remFn

	go r.broadcast()

	return &r
}

func (r *room) addClient(c *client) {
	r.Lock()
	defer r.Unlock()

	r.clients[c.userID] = c
	r.onlineClients++
}

func (r *room) addMessage(msg *models.Message) {
	r.Lock()
	defer r.Unlock()

	_ = r.store.SaveMessage(context.Background(), msg)

	r.messageChan <- msg
}

func (r *room) broadcast() {
	for msg := range r.messageChan {
		r.RLock()

		var err error
		for _, c := range r.clients {
			if err = c.Context().Err(); err != nil {
				r.removeClient(c.userID)
			}

			err = c.Send(r.formMessage(msg))
			if err != nil {
				r.removeClient(c.userID)
			}
		}

		r.RUnlock()
	}
}

func (r *room) removeClient(userID string) {
	r.Lock()
	defer r.Unlock()

	delete(r.clients, userID)

	r.onlineClients--
	if r.onlineClients == 0 {
		r.removeFromHubFunc(r.id)
	}
}

func (r *room) formMessage(msg *models.Message) *v1.ConnectResponse {
	return &v1.ConnectResponse{
		Payload: &v1.ConnectResponse_Message{
			Message: &v1.Message{
				ChatId:    msg.ChatID,
				SenderId:  msg.SenderID,
				Message:   msg.Message,
				CreatedAt: timestamppb.New(msg.CreatedAt),
			},
		},
	}
}
