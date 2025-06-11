package services

import (
	"database/sql"
	"fmt"
)

type ReactionService struct {
	DB *sql.DB
}

// NewReactionService est le constructeur pour le ReactionService.
func NewReactionService(db *sql.DB) *ReactionService {
	return &ReactionService{DB: db}
}

// HandleReaction gère la logique d'ajout/modification/suppression d'une réaction.
func (s *ReactionService) HandleReaction(userID, messageID int, reactionType string) error {
	var existingType string
	// Vérifier s'il y a une réaction existante
	err := s.DB.QueryRow("SELECT type FROM reaction WHERE user_id = ? AND message_id = ?", userID, messageID).Scan(&existingType)

	tx, errTx := s.DB.Begin()
	if errTx != nil {
		return fmt.Errorf("erreur début de transaction: %w", errTx)
	}

	if err == sql.ErrNoRows {
		// Aucune réaction existante, on en insère une nouvelle
		_, err = tx.Exec("INSERT INTO reaction (user_id, message_id, type) VALUES (?, ?, ?)", userID, messageID, reactionType)
	} else if err != nil {
		tx.Rollback()
		return fmt.Errorf("erreur lors de la vérification de la réaction: %w", err)
	} else {
		// Une réaction existe déjà
		if existingType == reactionType {
			// L'utilisateur clique sur le même bouton (ex: like sur un message déjà liké) -> on supprime la réaction
			_, err = tx.Exec("DELETE FROM reaction WHERE user_id = ? AND message_id = ?", userID, messageID)
		} else {
			// L'utilisateur change d'avis (ex: dislike sur un message liké) -> on met à jour
			_, err = tx.Exec("UPDATE reaction SET type = ? WHERE user_id = ? AND message_id = ?", reactionType, userID, messageID)
		}
	}

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erreur lors de l'opération sur la réaction: %w", err)
	}

	return tx.Commit()
}

// GetReactionCountsForTopic récupère les comptes de likes/dislikes pour tous les messages d'un sujet.
func (s *ReactionService) GetReactionCountsForTopic(topicID int) (map[int]map[string]int, error) {
	query := `
        SELECT message_id,
               SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END) as likes,
               SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END) as dislikes
        FROM reaction
        WHERE message_id IN (SELECT message_id FROM messages WHERE topic_id = ?)
        GROUP BY message_id`

	rows, err := s.DB.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[int]map[string]int)
	for rows.Next() {
		var messageID, likes, dislikes int
		if err := rows.Scan(&messageID, &likes, &dislikes); err != nil {
			return nil, err
		}
		counts[messageID] = map[string]int{"likes": likes, "dislikes": dislikes}
	}

	return counts, nil
}
