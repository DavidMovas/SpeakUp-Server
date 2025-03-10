package store

import (
	"context"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/models/requests"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/dbx"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
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

func (s *ChatsStore) CreatePrivateChat(ctx context.Context, request *requests.CreatePrivateChatRequest) (string, error) {
	chatBuilder := dbx.StatementBuilder.
		Insert("chats").
		Columns("id", "slug", "name", "type", "created_at").
		Values(request.ID, request.Slug, request.Name, request.Type, request.CreatedAt)

	query, args, err := chatBuilder.ToSql()

	err = dbx.InTransaction(ctx, s.db, func(ctx context.Context, tx pgx.Tx) error {
		_, err = s.db.Exec(ctx, query, args...)

		if err != nil {
			return apperrors.Internal(err)
		}

		now := time.Now()

		relations := []*chatServiceRelation{
			{chatID: request.ID, userID: request.InitiatorID, role: "admin", joinedAt: now},
			{chatID: request.ID, userID: request.MemberID, role: "admin", joinedAt: now},
		}

		membersBuilder := squirrel.StatementBuilder.
			Insert("chats_members").
			Columns("chat_id", "user_id", "role", "joined_at").
			Values(s.buildChatMembersRelation(relations))

		query, args, err = membersBuilder.ToSql()

		_, err = tx.Exec(ctx, query, args...)

		if err != nil {
			return apperrors.Internal(err)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return request.ID, nil
}

func (s *ChatsStore) CreateGroupChat(ctx context.Context, request *requests.CreateGroupChatRequest) (string, error) {
	chatBuilder := dbx.StatementBuilder.
		Insert("chats").
		Columns("id", "slug", "name", "type", "created_at").
		Values(request.ID, request.Slug, request.Name, request.Type, request.CreatedAt)

	query, args, err := chatBuilder.ToSql()

	err = dbx.InTransaction(ctx, s.db, func(ctx context.Context, tx pgx.Tx) error {
		_, err = s.db.Exec(ctx, query, args...)

		if err != nil {
			return apperrors.Internal(err)
		}

		now := time.Now()

		var relations []*chatServiceRelation
		for _, memberID := range request.MemberIDs {
			relations = append(relations, &chatServiceRelation{chatID: request.ID, userID: memberID, role: "member", joinedAt: now})
		}

		relations = append(relations, &chatServiceRelation{chatID: request.ID, userID: request.InitiatorID, role: "admin", joinedAt: now})

		membersBuilder := squirrel.StatementBuilder.
			Insert("chats_members").
			Columns("chat_id", "user_id", "role", "joined_at").
			Values(s.buildChatMembersRelation(relations))

		query, args, err = membersBuilder.ToSql()

		_, err = tx.Exec(ctx, query, args...)

		if err != nil {
			return apperrors.Internal(err)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return request.ID, nil
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

type chatServiceRelation struct {
	chatID   string
	userID   string
	role     string
	joinedAt time.Time
}

func (s *ChatsStore) buildChatMembersRelation(relations []*chatServiceRelation) [][]any {
	var values [][]any

	for _, r := range relations {
		values = append(values, []any{r.chatID, r.userID, r.role, r.joinedAt})
	}

	return values
}
