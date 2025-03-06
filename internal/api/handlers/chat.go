package handlers

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/services"
	"github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

var _ v1.ChatServiceServer = (*ChatHandler)(nil)

type ChatHandler struct {
	service *services.ChatService
	logger  *zap.Logger

	v1.UnimplementedChatServiceServer
}

func (h *ChatHandler) CreateRoom(ctx context.Context, request *v1.CreateRoomRequest) (*v1.CreateRoomResponse, error) {
	return nil, nil
}

func (h *ChatHandler) JoinRoom(stream grpc.BidiStreamingServer[v1.JoinRoomRequest, v1.JoinRoomResponse]) error {
	return nil
}

func NewChatHandler(service *services.ChatService, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ChatHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
