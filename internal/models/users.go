package models

import "time"

type User struct {
	ID          string
	Email       string
	Username    string
	AvatarURL   string
	FullName    string
	Bio         string
	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserWithPassword struct {
	*User
	PassHash string
}
