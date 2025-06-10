package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	TopicID   int       `json:"topic_id"`
	UserID    int       `json:"user_id"` // L'auteur du message
	CreatedAt time.Time `json:"created_at"`
}
