package models

import "time"

type RelationRoomToUser struct {
	ChatID   string
	UserID   string
	Role     string
	JoinedAt time.Time
}
