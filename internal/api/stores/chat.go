package stores

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ChatsStore struct {
	db     *pgxpool.Pool
	redis  *redis.Client
	logger *zap.Logger
}

func NewChatsStore(db *pgxpool.Pool, redis *redis.Client, logger *zap.Logger) *ChatsStore {
	return &ChatsStore{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}
