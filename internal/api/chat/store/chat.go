package store

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/dbx"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"github.com/Masterminds/squirrel"
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

func (s *ChatsStore) GetPrivateChatIDBetweenUsers(ctx context.Context, requestID, searchedID string) (string, error) {
	builder := dbx.StatementBuilder.
		Select("c.id").
		From("chats c").
		Join("chat_members cm1 ON c.id = cm1.chat_id").
		Join("chat_members cm2 ON c.id = cm2.chat_id").
		Where(squirrel.Eq{"cm1.user_id": requestID}).
		Where(squirrel.Eq{"cm2.user_id": searchedID}).
		Where(squirrel.Eq{"c.type": "private"})

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	var chatID string
	err = s.db.QueryRow(ctx, query, args...).Scan(&chatID)

	switch {
	case dbx.IsNoRows(err):
		return "", apperrors.NotFound("chat", "user_id", searchedID)
	case err != nil:
		return "", apperrors.Internal(err)
	}

	return chatID, nil
}

func (s *ChatsStore) GetGroupChatIDBetweenUsers(ctx context.Context, userIDs ...string) (string, error) {
	builder := dbx.StatementBuilder.
		Select("c.id").
		From("chats c").
		Join("chats_members cm ON c.id = cm.chat_id").
		Where(squirrel.Eq{"cm.user_id": userIDs}).
		Where(squirrel.Eq{"c.type": "group"}).
		GroupBy("c.id").
		Having("COUNT(DISTINCT cm.user_id) = ?", len(userIDs))

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	var chatID string
	err = s.db.QueryRow(ctx, query, args...).Scan(&chatID)

	switch {
	case dbx.IsNoRows(err):
		return "", apperrors.NotFound("chat", "user_id", userIDs)
	case err != nil:
		return "", apperrors.Internal(err)
	}

	return chatID, nil
}
