package services

import (
	"Forum/models"
	"database/sql"
	"fmt"
)

type TopicService struct {
	DB *sql.DB
}

func NewTopicService(db *sql.DB) *TopicService {
	return &TopicService{DB: db}
}

// Create insère un nouveau sujet en base de données.
func (s *TopicService) Create(topic *models.Topic) (int, error) {
	// Par défaut, un sujet est créé dans le forum 1 et avec le statut "ouvert" (true)
	result, err := s.DB.Exec(
		"INSERT INTO sujet (forum_id, user_id, title, status, description) VALUES (?, ?, ?, ?, ?)",
		1, // Forum ID par défaut, à adapter plus tard
		topic.UserID,
		topic.Title,
		true, // Statut "ouvert"
		topic.Description,
	)
	if err != nil {
		return 0, fmt.Errorf("TopicService.Create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("TopicService.Create LastInsertId: %w", err)
	}
	return int(id), nil
}

// GetAllTopics récupère tous les sujets pour affichage.
func (s *TopicService) GetAllTopics() ([]models.Topic, error) {
	// On utilise une jointure pour récupérer le nom de l'auteur directement
	// CORRECTION : On retire s.description de la requête SELECT
	query := `
        SELECT s.topic_id, s.title, s.created_at, u.name
        FROM sujet s
        JOIN Utilisateurs u ON s.user_id = u.user_id
        ORDER BY s.created_at DESC`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("TopicService.GetAllTopics Query: %w", err)
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		// CORRECTION : On retire &t.Description de la méthode Scan
		if err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.AuthorName); err != nil {
			return nil, fmt.Errorf("TopicService.GetAllTopics Scan: %w", err)
		}
		topics = append(topics, t)
	}

	return topics, nil
}
