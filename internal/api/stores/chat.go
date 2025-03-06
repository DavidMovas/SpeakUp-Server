package stores

import (
	"github.com/DavidMovas/SpeakUp-Server/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ChatsStore struct {
	db     *pgxpool.Pool
	redis  *redis.Client
	logger *zap.Logger

	roomHubs map[string]*models.RoomHub
}

func NewChatsStore(db *pgxpool.Pool, redis *redis.Client, logger *zap.Logger) *ChatsStore {
	return &ChatsStore{
		db:       db,
		redis:    redis,
		logger:   logger,
		roomHubs: make(map[string]*models.RoomHub),
	}
}
