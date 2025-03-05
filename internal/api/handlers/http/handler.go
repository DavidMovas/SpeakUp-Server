package http

import (
	"github.com/DavidMovas/SpeakUp-Server/internal/api/services"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	service *services.Service
	logger  *zap.Logger
}

func NewHandler(service *services.Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
