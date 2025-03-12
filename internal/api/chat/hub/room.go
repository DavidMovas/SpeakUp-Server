package hub

import (
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

type room struct {
	sync.RWMutex
	id            string
	clients       map[string]*client
	onlineClients uint
	messageChan   chan *Message

	removeFromHubFunc func(id string)
}

func newRoom(id string, remFn func(id string)) *room {
	var r room
	r.id = id
	r.clients = make(map[string]*client)
	r.messageChan = make(chan *Message, 50)
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

func (r *room) addMessage(msg *Message) {
	r.Lock()
	defer r.Unlock()

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
	if r.onlineClients < 1 {
		r.removeFromHubFunc(r.id)
	}
}

func (r *room) formMessage(msg *Message) *v1.ConnectResponse {
	return &v1.ConnectResponse{
		Payload: &v1.ConnectResponse_Message{
			Message: &v1.Message{
				ChatId:    msg.ChatId,
				SenderId:  msg.SenderId,
				Message:   msg.Message,
				CreatedAt: timestamppb.New(msg.CreatedAt),
			},
		},
	}
}
