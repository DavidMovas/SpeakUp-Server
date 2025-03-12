package hub

import (
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"sync"
)

type room struct {
	sync.RWMutex
	id          string
	clients     map[string]*client
	messageChan chan *v1.Message
}

func newRoom(id string) *room {
	var r room
	r.id = id
	r.clients = make(map[string]*client)
	r.messageChan = make(chan *v1.Message, 50)

	go r.broadcast()

	return &r
}

func (r *room) addClient(c *client) {
	r.Lock()
	defer r.Unlock()

	r.clients[c.userID] = c
}

func (r *room) addMessage(msg *v1.Message) {
	r.Lock()
	defer r.Unlock()

	r.messageChan <- msg
}

func (r *room) broadcast() {
	for msg := range r.messageChan {
		r.RLock()

		pMsg := &v1.ConnectResponse{
			Payload: &v1.ConnectResponse_Message{
				Message: msg,
			},
		}

		for _, c := range r.clients {
			_ = c.Send(pMsg)
		}

		r.RUnlock()
	}
}

func (r *room) removeClient(userID string) {
	r.Lock()
	defer r.Unlock()

	delete(r.clients, userID)
}
