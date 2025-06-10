package models

import "time"

type Topic struct {
	ID          int       `json:"id"`
	ForumID     int       `json:"forum_id"` // Pour plus tard
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorName  string    `json:"author_name"` // Champ supplémentaire pour l'affichage
}
