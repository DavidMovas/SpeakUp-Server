package stores

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/models"
	"github.com/DavidMovas/SpeakUp-Server/internal/models/requests"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/dbx"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"github.com/Masterminds/squirrel"
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

func (s *UsersStore) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error) {
	builder := dbx.StatementBuilder.
		Insert("users").
		Columns("id", "email", "username", "full_name", "pass_hash").
		Values(request.ID, request.Email, request.Username, request.FullName, request.Password).
		Suffix("RETURNING created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := models.User{
		ID:       request.ID,
		Email:    request.Email,
		Username: request.Username,
		FullName: request.FullName,
	}

	err = s.db.QueryRow(ctx, query, args...).Scan(&user.CreatedAt)

	switch {
	case dbx.IsUniqueViolation(err, "email"):
		return nil, apperrors.AlreadyExists("user", "email", request.Email)
	case dbx.IsUniqueViolation(err, "username"):
		return nil, apperrors.AlreadyExists("user", "username", request.Username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &user, nil
}

func (s *UsersStore) GetUserByEmail(ctx context.Context, request *requests.GetUserByEmailRequest) (*models.UserWithPassword, error) {
	builder := dbx.StatementBuilder.
		Select("id", "email", "username", "avatar_url", "full_name", "bio", "pass_hash", "last_login_at", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": request.Email})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var user models.UserWithPassword
	user.User = &models.User{}

	err = s.db.QueryRow(ctx, query, args...).
		Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.AvatarURL,
			&user.FullName,
			&user.Bio,
			&user.PassHash,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", request.Email)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &user, nil
}
