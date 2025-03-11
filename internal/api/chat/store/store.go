package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "github.com/DavidMovas/SpeakUp-Server/internal/shared/grpc/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Префикс для ключей чата в Redis
	chatKeyPrefix = "chat:%s:messages"
	// Лимит сообщений в кеше
	cacheLimit = 100
)

type ChatStore struct {
	pg     *pgxpool.Pool
	redis  *redis.Client
	logger *zap.Logger
}

func NewChatStore(pg *pgxpool.Pool, redis *redis.Client, logger *zap.Logger) *ChatStore {
	return &ChatStore{
		pg:     pg,
		redis:  redis,
		logger: logger,
	}
}

// SaveMessage сохраняет сообщение в PostgreSQL и Redis
func (s *ChatStore) SaveMessage(ctx context.Context, msg *v1.Message) error {
	// Сохраняем в PostgreSQL
	query := `
		INSERT INTO messages (chat_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var msgID int64
	err := s.pg.QueryRow(ctx, query,
		msg.ChatId,
		msg.SenderId,
		msg.Message,
		msg.CreatedAt.AsTime(),
	).Scan(&msgID)

	if err != nil {
		return fmt.Errorf("failed to save message to postgres: %w", err)
	}

	// Сохраняем в Redis
	return s.cacheMessage(ctx, msg)
}

// GetUnreadMessages получает непрочитанные сообщения для пользователя
func (s *ChatStore) GetUnreadMessages(ctx context.Context, chatID string, lastReadAt time.Time) ([]*v1.Message, error) {
	// Сначала пробуем получить из Redis
	messages, err := s.getFromCache(ctx, chatID, lastReadAt)
	if err == nil && len(messages) > 0 {
		s.logger.Debug("Got messages from cache",
			zap.String("chat_id", chatID),
			zap.Int("count", len(messages)))
		return messages, nil
	}

	// Если в кеше нет или произошла ошибка, берем из PostgreSQL
	return s.getFromPostgres(ctx, chatID, lastReadAt)
}

func (s *ChatStore) cacheMessage(ctx context.Context, msg *v1.Message) error {
	// Сериализуем сообщение
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	key := fmt.Sprintf(chatKeyPrefix, msg.ChatId)
	pipe := s.redis.Pipeline()

	// Добавляем сообщение в отсортированный список
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(msg.CreatedAt.AsTime().Unix()),
		Member: data,
	})

	// Оставляем только последние cacheLimit сообщений
	pipe.ZRemRangeByRank(ctx, key, 0, -(cacheLimit + 1))

	// Устанавливаем TTL на 24 часа
	pipe.Expire(ctx, key, 24*time.Hour)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to cache message: %w", err)
	}

	return nil
}

func (s *ChatStore) getFromCache(ctx context.Context, chatID string, lastReadAt time.Time) ([]*v1.Message, error) {
	key := fmt.Sprintf(chatKeyPrefix, chatID)

	// Получаем сообщения после lastReadAt
	data, err := s.redis.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", lastReadAt.Unix()),
		Max:    "+inf",
		Offset: 0,
		Count:  cacheLimit,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to get messages from cache: %w", err)
	}

	messages := make([]*v1.Message, 0, len(data))
	for _, item := range data {
		msg := &v1.Message{}
		if err := json.Unmarshal([]byte(item), msg); err != nil {
			s.logger.Error("Failed to unmarshal message from cache",
				zap.Error(err),
				zap.String("data", item))
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *ChatStore) getFromPostgres(ctx context.Context, chatID string, lastReadAt time.Time) ([]*v1.Message, error) {
	query := `
		SELECT user_id, content, created_at
		FROM messages
		WHERE chat_id = $1 AND created_at > $2
		ORDER BY created_at ASC
		LIMIT $3`

	rows, err := s.pg.Query(ctx, query, chatID, lastReadAt, cacheLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	messages := make([]*v1.Message, 0, cacheLimit)
	for rows.Next() {
		msg := &v1.Message{
			ChatId: chatID,
		}
		var createdAt time.Time

		err = rows.Scan(
			&msg.SenderId,
			&msg.Message,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		msg.CreatedAt = timestamppb.New(createdAt)
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	// Кешируем полученные сообщения
	for _, msg := range messages {
		if err = s.cacheMessage(ctx, msg); err != nil {
			s.logger.Error("Failed to cache message from postgres",
				zap.Error(err),
				zap.String("chat_id", chatID),
				zap.String("message_id", msg.SenderId))
		}
	}

	return messages, nil
}
