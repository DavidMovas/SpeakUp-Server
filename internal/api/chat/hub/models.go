package hub

import "time"

type Message struct {
	ChatId    string
	SenderId  string
	Message   string
	CreatedAt time.Time
}
