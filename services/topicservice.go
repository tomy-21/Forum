package services

import (
	"Forum/models"
	"database/sql"
	"fmt"
	"strings" // Import du package strings pour le query builder
)

type TopicService struct {
	DB *sql.DB
}

// NewTopicService est le constructeur pour le TopicService.
func NewTopicService(db *sql.DB) *TopicService {
	return &TopicService{DB: db}
}

// Create insère un nouveau sujet en base de données.
func (s *TopicService) Create(topic *models.Topic) (int, error) {
	result, err := s.DB.Exec(
		"INSERT INTO sujet (forum_id, user_id, title, status) VALUES (?, ?, ?, ?)",
		1, // Forum ID par défaut
		topic.UserID,
		topic.Title,
		true, // Statut "ouvert" par défaut
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

// GetTopicByID récupère un seul sujet par son ID, avec le nom de l'auteur.
func (s *TopicService) GetTopicByID(id int) (models.Topic, error) {
	var t models.Topic
	query := `
        SELECT s.topic_id, s.user_id, s.title, s.created_at, u.name as author_name
        FROM sujet s
        JOIN Utilisateurs u ON s.user_id = u.user_id
        WHERE s.topic_id = ?`

	err := s.DB.QueryRow(query, id).Scan(&t.ID, &t.UserID, &t.Title, &t.CreatedAt, &t.AuthorName)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("aucun sujet trouvé avec l'ID %d", id)
		}
		return t, fmt.Errorf("TopicService.GetTopicByID: %w", err)
	}
	return t, nil
}

// GetTotalTopicCount retourne le nombre total de sujets, optionnellement filtré par catégorie.
// MODIFIÉ : La fonction accepte maintenant un categoryID.
func (s *TopicService) GetTotalTopicCount(categoryID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM sujet"
	var args []interface{}

	// Si un ID de catégorie est fourni, on ajoute une clause WHERE.
	if categoryID > 0 {
		// Votre colonne s'appelle 'forum_id' mais représente une catégorie, nous l'utilisons ici.
		query += " WHERE forum_id = ?"
		args = append(args, categoryID)
	}

	err := s.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetTotalTopicCount: %w", err)
	}
	return count, nil
}

// GetTopics récupère une liste paginée de sujets, optionnellement filtrée par catégorie.
// MODIFIÉ : La fonction accepte maintenant un categoryID.
func (s *TopicService) GetTopics(limit, offset, categoryID int) ([]models.Topic, error) {
	var query strings.Builder
	var args []interface{}

	query.WriteString(`
        SELECT s.topic_id, s.title, s.created_at, u.name
        FROM sujet s
        JOIN Utilisateurs u ON s.user_id = u.user_id`)

	// Si un ID de catégorie est fourni, on ajoute une clause WHERE.
	if categoryID > 0 {
		query.WriteString(" WHERE s.forum_id = ?")
		args = append(args, categoryID)
	}

	query.WriteString(" ORDER BY s.created_at DESC LIMIT ? OFFSET ?")
	args = append(args, limit, offset)

	rows, err := s.DB.Query(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("GetTopics Query: %w", err)
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.AuthorName); err != nil {
			return nil, fmt.Errorf("GetTopics Scan: %w", err)
		}
		topics = append(topics, t)
	}
	return topics, nil
}

// DeleteTopic supprime un sujet de la base de données.
func (s *TopicService) DeleteTopic(topicID int) error {
	_, err := s.DB.Exec("DELETE FROM sujet WHERE topic_id = ?", topicID)
	if err != nil {
		return fmt.Errorf("TopicService.DeleteTopic: %w", err)
	}
	return nil
}
