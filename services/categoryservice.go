package services

import (
	"Forum/models"
	"database/sql"
	"fmt"
)

type CategoryService struct {
	DB *sql.DB
}

func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{DB: db}
}

// GetAllCategories récupère toutes les catégories de la BDD.
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	rows, err := s.DB.Query("SELECT category_id, name, description, created_at FROM categories ORDER BY created_at ASC")
	if err != nil {
		return nil, fmt.Errorf("CategoryService.GetAllCategories: Query error: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("CategoryService.GetAllCategories: Scan error: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}
