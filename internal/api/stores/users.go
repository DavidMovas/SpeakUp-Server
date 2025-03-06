package stores

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/models"
	"github.com/DavidMovas/SpeakUp-Server/internal/models/requests"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/dbx"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
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

func (s *UsersStore) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error) {
	createdAt := time.Now()

	builder := dbx.StatementBuilder.
		Insert("users").
		Columns("id", "email", "username", "full_name", "pass_hash", "created_at").
		Values(request.ID, request.Email, request.Username, request.FullName, request.Password, createdAt)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := &models.User{
		ID:        request.ID,
		Email:     request.Email,
		Username:  request.Username,
		FullName:  request.FullName,
		CreatedAt: createdAt,
	}

	return user, nil
}
