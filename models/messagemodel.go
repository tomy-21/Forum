package models

import "time"

type Message struct {
	ID         int       `json:"id"`
	TopicID    int       `json:"topic_id"`
	UserID     int       `json:"user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	AuthorName string    `json:"author_name"`

	Likes    int `json:"likes"`
	Dislikes int `json:"dislikes"`
}
