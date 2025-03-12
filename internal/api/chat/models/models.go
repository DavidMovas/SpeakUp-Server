package models

import "time"

type Message struct {
	ChatID    string    `json:"chatId"`
	SenderID  string    `json:"senderId"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
