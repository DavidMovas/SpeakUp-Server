package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Repository struct {
	db     *pgxpool.Pool
	redis  *redis.Client
	logger *zap.Logger
}

func NewRepository(db *pgxpool.Pool, redis *redis.Client, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}
