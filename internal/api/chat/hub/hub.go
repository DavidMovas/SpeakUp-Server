package hub

import (
	"context"
	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"sync"
	"time"
)

type clientStream = grpc.BidiStreamingServer[v1.ConnectRequest, v1.ConnectResponse]
type clientSet = map[string]*client

type Hub struct {
	sync.RWMutex
	ctx          context.Context
	clients      clientSet
	readStoreCh  <-chan *v1.Message
	writeStoreCh chan<- *v1.Message
	logger       *zap.Logger
}

func NewHub(ctx context.Context, readStoreCh <-chan *v1.Message, writeStoreCh chan<- *v1.Message, logger *zap.Logger) *Hub {
	return &Hub{
		ctx:          ctx,
		RWMutex:      sync.RWMutex{},
		readStoreCh:  readStoreCh,
		writeStoreCh: writeStoreCh,
		clients:      make(clientSet),
		logger:       logger,
	}
}

func (h *Hub) RegisterStream(userID string, stream clientStream) {
	h.Lock()
	defer h.Unlock()

	h.logger.Info("Register stream", zap.String("user_id", userID))

	if c, ok := h.clients[userID]; ok {
		c.BidiStreamingServer = stream
		return
	}

	c := newClient(userID, stream)
	h.clients[userID] = c

	go h.manage(c)
}

func (h *Hub) manage(c *client) {
	defer func() {
		h.Lock()
		delete(h.clients, c.userID)
		h.Unlock()
	}()

	for {
		if sug := h.ctx.Done(); sug != nil {
			return
		}

		req, err := c.Recv()
		if err == io.EOF {
			h.logger.Debug("Client disconnected", zap.String("user_id", c.userID))
			return
		}
		if err != nil {
			h.logger.Warn("Error reading from client", zap.String("user_id", c.userID), zap.Error(err))
			return
		}

		switch pl := req.GetPayload().(type) {
		case *v1.ConnectRequest_JointChat:
			h.handleJoinChat(c.userID, pl.JointChat)
		case *v1.ConnectRequest_Message:
			h.handleMessage(pl.Message)
		}
	}
}

func (h *Hub) handleJoinChat(userID string, joinChat *v1.ConnectRequest_JoinChat) {
	// Получаем историю сообщений из БД
	//messages := GetChatHistoryFromDB(joinChat.ChatId, joinChat.LastReadAt)

	// Отправляем историю пользователю
	c, ok := h.clients[userID]
	if !ok {
		return
	}

	reps := &v1.ConnectResponse{
		Payload: &v1.ConnectResponse_Message{
			Message: &v1.Message{
				ChatId:    joinChat.GetChatId(),
				SenderId:  "SYSTEM",
				Message:   "Joined Chat Successfully",
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	}

	if err := c.Send(reps); err != nil {
		h.logger.Warn("Error sending join chat", zap.String("user_id", c.userID), zap.Error(err))
	}
}

func (h *Hub) handleMessage(msg *v1.Message) {
	// Сохраняем сообщение в БД
	//SaveMessageToDB(msg)

	// Рассылаем сообщение всем участникам чата
	h.Lock()
	defer h.Unlock()

	now := time.Now()

	reps := &v1.ConnectResponse{
		Payload: &v1.ConnectResponse_Message{
			Message: &v1.Message{
				ChatId:    msg.ChatId,
				SenderId:  msg.SenderId,
				Message:   msg.Message,
				CreatedAt: timestamppb.New(now),
			},
		},
	}

	for _, c := range h.clients {
		if err := c.Send(reps); err != nil {
			h.logger.Warn("Error sending message to client", zap.String("user_id", c.userID))
		}
	}
}
