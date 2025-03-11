package hub

import (
	"context"
	"io"
	"sync"
	"time"

	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type clientStream = grpc.BidiStreamingServer[v1.ConnectRequest, v1.ConnectResponse]
type clientSet = map[string]*client

// Room представляет комнату чата
type Room struct {
	sync.RWMutex
	id      string
	clients map[string]*client
}

func newRoom(id string) *Room {
	return &Room{
		id:      id,
		clients: make(map[string]*client),
	}
}

func (r *Room) addClient(c *client) {
	r.Lock()
	defer r.Unlock()
	r.clients[c.userID] = c
}

func (r *Room) removeClient(userID string) {
	r.Lock()
	defer r.Unlock()
	delete(r.clients, userID)
}

func (r *Room) isEmpty() bool {
	r.RLock()
	defer r.RUnlock()
	return len(r.clients) == 0
}

func (r *Room) broadcast(message *v1.ConnectResponse, logger *zap.Logger) {
	r.RLock()
	defer r.RUnlock()

	for _, c := range r.clients {
		if err := c.Send(message); err != nil {
			logger.Warn("Failed to send message to c",
				zap.String("user_id", c.userID),
				zap.Error(err))
		}
	}
}

type Hub struct {
	sync.RWMutex
	ctx          context.Context
	clients      clientSet
	rooms        map[string]*Room
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
		rooms:        make(map[string]*Room),
		logger:       logger,
	}

	go hub.handleStorageMessages()

	return hub
}

func (h *Hub) handleStorageMessages() {
	for {
		select {
		case <-h.ctx.Done():
			return
		case msg := <-h.readStoreCh:
			if msg == nil {
				continue
			}

			resp := &v1.ConnectResponse{
				Payload: &v1.ConnectResponse_Message{
					Message: msg,
				},
			}

			h.RLock()
			if room, ok := h.rooms[msg.ChatId]; ok {
				room.broadcast(resp, h.logger)
			}
			h.RUnlock()
		}
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

func (h *Hub) joinRoom(userID, roomID string) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = newRoom(roomID)
	}

	if c, ok := h.clients[userID]; ok {
		h.rooms[roomID].addClient(c)
		h.notifyRoomJoin(roomID, userID)
	}
}

func (h *Hub) leaveRoom(userID, roomID string) {
	h.Lock()
	defer h.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		room.removeClient(userID)
		if room.isEmpty() {
			delete(h.rooms, roomID)
		}
		h.notifyRoomJoin(roomID, userID)
	}
}

func (h *Hub) notifyRoomJoin(roomID, userID string) {
	msg := &v1.ConnectResponse{
		Payload: &v1.ConnectResponse_Message{
			Message: &v1.Message{
				ChatId:    roomID,
				SenderId:  userID,
				Message:   "joined the chat",
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	}

	if room, ok := h.rooms[roomID]; ok {
		room.broadcast(msg, h.logger)
	}
}

func (h *Hub) notifyRoomLeave(roomID, userID string) {
	msg := &v1.ConnectResponse{
		Payload: &v1.ConnectResponse_Message{
			Message: &v1.Message{
				ChatId:    roomID,
				SenderId:  userID,
				Message:   "left the chat",
				CreatedAt: timestamppb.New(time.Now()),
			},
		},
	}

	if room, ok := h.rooms[roomID]; ok {
		room.broadcast(msg, h.logger)
	}
}

func (h *Hub) manage(c *client) {
	defer func() {
		h.Lock()
		for roomID, room := range h.rooms {
			room.removeClient(c.userID)
			if room.isEmpty() {
				delete(h.rooms, roomID)
			}
		}
		delete(h.clients, c.userID)
		h.Unlock()
	}()

	for {
		select {
		case <-h.ctx.Done():
			return
		default:
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
}

func (h *Hub) handleJoinChat(userID string, joinChat *v1.ConnectRequest_JoinChat) {
	chatID := joinChat.GetChatId()
	h.joinRoom(userID, chatID)

	// Запрашиваем непрочитанные сообщения из хранилища
	if h.writeStoreCh != nil {
		lastRead := joinChat.GetLastReadAt()
		h.writeStoreCh <- &v1.Message{
			ChatId:    chatID,
			SenderId:  userID,
			Message:   "HISTORY_REQUEST",
			CreatedAt: lastRead,
		}
	}
}

func (h *Hub) handleMessage(msg *v1.Message) {
	msg.CreatedAt = timestamppb.New(time.Now())

	// Сохраняем сообщение в хранилище
	if h.writeStoreCh != nil {
		h.writeStoreCh <- msg
	}

	resp := &v1.ConnectResponse{
		Payload: &v1.ConnectResponse_Message{
			Message: msg,
		},
	}

	// Отправляем сообщение всем участникам комнаты
	if room, ok := h.rooms[msg.ChatId]; ok {
		room.broadcast(resp, h.logger)
	}
}
