package models

type Room struct {
	ID        string `json:"id"`
	OwnerID   string `json:"owner_id"`
	Name      string `json:"name"`
	Broadcast chan []byte
}
