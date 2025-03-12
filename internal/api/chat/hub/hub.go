package hub

import (
	"context"
	"fmt"
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"sync"
)

type clientStream = grpc.BidiStreamingServer[v1.ConnectRequest, v1.ConnectResponse]
type clientSet = map[string]*client

type Hub struct {
	sync.RWMutex
	ctx          context.Context
	clients      clientSet
	rooms        map[string]*room
	readStoreCh  <-chan *v1.Message
	writeStoreCh chan<- *v1.Message
	logger       *zap.Logger
}

func NewHub(ctx context.Context, readStoreCh <-chan *v1.Message, writeStoreCh chan<- *v1.Message, logger *zap.Logger) *Hub {
	hub := &Hub{
		ctx:          ctx,
		RWMutex:      sync.RWMutex{},
		readStoreCh:  readStoreCh,
		writeStoreCh: writeStoreCh,
		clients:      make(clientSet),
		rooms:        make(map[string]*room),
		logger:       logger,
	}

	return hub
}

func (h *Hub) Connect(stream clientStream) error {
	h.logger.Debug("HUB CONNECT")

	in, err := stream.Recv()
	if err != nil {
		h.logger.Error("Error receiving from stream", zap.Error(err))
		return fmt.Errorf("error receiving from stream: %w", err)
	}

	connectRoomReq, ok := in.Payload.(*v1.ConnectRequest_JoinChat_)
	if !ok {
		h.logger.Error("Error receiving from stream", zap.Any("in", in))
		return fmt.Errorf("error receiving from stream: %w", err)
	}

	h.logger.Debug("HUB JoinChat Request", zap.String("chat_id", connectRoomReq.JoinChat.ChatId), zap.String("user_id", connectRoomReq.JoinChat.UserId))

	userID := connectRoomReq.JoinChat.UserId

	r := h.getRoom(connectRoomReq.JoinChat.ChatId)
	r.addClient(newClient(userID, stream))

	for {
		in, err = stream.Recv()
		switch {
		case err == io.EOF:
			r.removeClient(userID)
			return nil
		case err != nil:
			r.removeClient(userID)
			h.logger.Error("Error receiving from stream", zap.Error(err))
			return fmt.Errorf("error receiving from stream: %w", err)
		case stream.Context().Err() != nil:
			r.removeClient(userID)
			h.logger.Error("Error receiving from stream", zap.Error(stream.Context().Err()))
			return fmt.Errorf("error receiving from stream: %w", stream.Context().Err())
		}

		switch p := in.Payload.(type) {
		case *v1.ConnectRequest_Message:
			r.addMessage(p.Message)
		}
	}
}

func (h *Hub) getRoom(chatID string) *room {
	h.RLock()
	defer h.RUnlock()

	if r, ok := h.rooms[chatID]; ok {
		return r
	} else {
		r = newRoom(chatID)
		h.rooms[chatID] = r
		return r
	}
}
