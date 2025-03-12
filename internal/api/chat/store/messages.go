package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DavidMovas/SpeakUp-Server/internal/api/chat/models"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/dbx"
	"github.com/DavidMovas/SpeakUp-Server/internal/utils/helpers"
	"golang.org/x/sync/errgroup"
	"time"
)

const (
	cacheMessageCacheKey = "chat:%s:messages"
	cacheMessagesAmount  = 100
	cacheMessagesTLL     = time.Hour
)

func (s *ChatsStore) SaveMessage(ctx context.Context, msg *models.Message) error {
	group := errgroup.Group{}

	group.Go(func() error {
		return s.saveMessagesToCache(ctx, msg)
	})

	group.Go(func() error {
		return s.saveMessages(ctx, msg)
	})

	return group.Wait()
}

func (s *ChatsStore) saveMessagesToCache(ctx context.Context, msg *models.Message) error {
	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error while marshaling message: %w", err)
	}

	key := fmt.Sprintf(cacheMessageCacheKey, msg.ChatID)

	cmd := s.redis.LPush(ctx, key, msgData)
	if cmd.Err() != nil {
		return fmt.Errorf("error while saving message to cache: %w", cmd.Err())
	}

	s.redis.ExpireLT(ctx, key, cacheMessagesTLL)
	s.redis.LTrim(ctx, key, 0, cacheMessagesAmount-1)

	return nil
}

func (s *ChatsStore) saveMessages(ctx context.Context, msg *models.Message) error {
	id := helpers.GenerateID()

	builder := dbx.StatementBuilder.
		Insert("messages").
		Columns("id", "chat_id", "user_id", "content", "created_at").
		Values(id, msg.ChatID, msg.SenderID, msg.Message, msg.CreatedAt)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("error while building query: %w", err)
	}

	cmd, err := s.db.Exec(ctx, query, args...)
	switch {
	case err != nil:
		return fmt.Errorf("error while saving message: %w", err)
	case cmd.RowsAffected() == 0:
		return fmt.Errorf("message did not saved to postgres: %w", err)
	}

	return nil
}
