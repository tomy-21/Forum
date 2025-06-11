package services

import (
	"Forum/models"
	"database/sql"
	"fmt"
)

type MessageService struct {
	DB *sql.DB
}

func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{DB: db}
}

// CreateMessage insère un nouveau message dans la base de données.
func (s *MessageService) CreateMessage(message *models.Message) error {
	_, err := s.DB.Exec(
		"INSERT INTO messages (topic_id, user_id, content) VALUES (?, ?, ?)",
		message.TopicID,
		message.UserID,
		message.Content,
	)
	if err != nil {
		return fmt.Errorf("MessageService.CreateMessage: %w", err)
	}
	return nil
}

// GetMessagesByTopicID récupère tous les messages d'un sujet, avec leurs comptes de réactions, et les trie.
func (s *MessageService) GetMessagesByTopicID(topicID int, sortBy string) ([]models.Message, error) {
	// On construit la requête de base avec une jointure à gauche pour inclure les messages sans réaction
	baseQuery := `
        SELECT
            m.message_id, m.content, m.created_at, u.name as author_name,
            COALESCE(SUM(CASE WHEN r.type = 'like' THEN 1 ELSE 0 END), 0) as likes,
            COALESCE(SUM(CASE WHEN r.type = 'dislike' THEN 1 ELSE 0 END), 0) as dislikes
        FROM messages m
        JOIN Utilisateurs u ON m.user_id = u.user_id
        LEFT JOIN reaction r ON m.message_id = r.message_id
        WHERE m.topic_id = ?
        GROUP BY m.message_id, m.content, m.created_at, u.name
    `

	// On ajoute la clause de tri en fonction du paramètre `sortBy`
	var orderByClause string
	switch sortBy {
	case "top":
		// Trier par popularité (likes - dislikes), puis par date
		orderByClause = "ORDER BY (likes - dislikes) DESC, m.created_at DESC"
	case "old":
		// Trier par le plus ancien
		orderByClause = "ORDER BY m.created_at ASC"
	default:
		// Tri par défaut : le plus récent
		orderByClause = "ORDER BY m.created_at DESC"
	}

	query := fmt.Sprintf("%s %s", baseQuery, orderByClause)

	rows, err := s.DB.Query(query, topicID)
	if err != nil {
		return nil, fmt.Errorf("MessageService.GetMessagesByTopicID Query: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt, &msg.AuthorName, &msg.Likes, &msg.Dislikes); err != nil {
			return nil, fmt.Errorf("MessageService.GetMessagesByTopicID Scan: %w", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
