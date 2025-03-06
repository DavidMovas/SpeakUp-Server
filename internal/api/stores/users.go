package stores

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UsersStore struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewUsersStore(db *pgxpool.Pool, logger *zap.Logger) *UsersStore {
	return &UsersStore{
		db:     db,
		logger: logger,
	}
}

func (s *UsersStore) CreateUser(ctx context.Context) error {
	return nil
}
