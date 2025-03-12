package hub

import (
	"context"
	"fmt"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/models"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/store"
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	met "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

var (
	roomsAmountHist met.Int64Histogram
)

type clientStream = grpc.BidiStreamingServer[v1.ConnectRequest, v1.ConnectResponse]
type clientSet = map[string]*client

type Hub struct {
	sync.RWMutex
	ctx     context.Context
	clients clientSet
	rooms   map[string]*room
	store   *store.ChatsStore

	logger *zap.Logger
	promet *metric.MeterProvider
}

func NewHub(ctx context.Context, store *store.ChatsStore, logger *zap.Logger, promet *metric.MeterProvider) *Hub {
	hub := &Hub{
		ctx:     ctx,
		store:   store,
		clients: make(clientSet),
		rooms:   make(map[string]*room),
		logger:  logger,
		promet:  promet,
	}

	roomsAmountHist, _ = hub.promet.Meter("chat_hub").Int64Histogram("rooms_amount")

	return hub
}

func (h *Hub) HandleStream(stream clientStream) error {
	return h.handleStream(stream)
}

func (h *Hub) handleStream(stream clientStream) error {
	for {
		if h.ctx.Err() != nil {
			return nil
		}

		in, err := stream.Recv()
		switch {
		case err == io.EOF:
			return nil
		case stream.Context().Err() != nil:
			return nil
		case err != nil:
			h.logger.Error("Error receiving from stream", zap.Error(err))
			return fmt.Errorf("error receiving from stream: %w", err)
		}

		switch p := in.Payload.(type) {
		case *v1.ConnectRequest_Message:
			msg := &models.Message{
				ChatID:    p.Message.ChatId,
				SenderID:  p.Message.SenderId,
				Message:   p.Message.Message,
				CreatedAt: time.Now(),
			}

			h.handleMessage(p.Message.ChatId, msg)

		case *v1.ConnectRequest_JoinChat_:
			h.handleJoin(p.JoinChat.ChatId, p.JoinChat.UserId, stream)
		}
	}
}

func (h *Hub) handleJoin(chatID, userID string, stream clientStream) {
	h.getRoom(chatID).addClient(newClient(userID, stream))
}

func (h *Hub) handleMessage(chatID string, msg *models.Message) {
	h.rooms[chatID].addMessage(msg)
}

func (h *Hub) getRoom(chatID string) *room {
	h.RLock()
	defer h.RUnlock()

	r, ok := h.rooms[chatID]
	if ok {
		return r
	}

	r = newRoom(chatID, h.store, h.removeRoom)
	h.rooms[chatID] = r

	roomsAmountHist.Record(h.ctx, int64(len(h.rooms)))

	return r
}

func (h *Hub) removeRoom(chatID string) {
	h.Lock()
	defer h.Unlock()

	delete(h.rooms, chatID)

	roomsAmountHist.Record(h.ctx, int64(len(h.rooms)))
}
