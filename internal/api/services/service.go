package services

import (
	"github.com/DavidMovas/SpeakUp-Server/internal/api/repository"
	"go.uber.org/zap"
)

type Service struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewService(repo *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}
