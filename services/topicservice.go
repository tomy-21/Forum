// Fichier : Forum/services/productservice.go

package services

import (
	"Forum/models"
	"database/sql" // MODIFIÉ : On n'importe plus "Forum/database"
	"fmt"
)

// MODIFIÉ : La struct contient maintenant un champ pour la connexion à la BDD.
type ProductService struct {
	DB *sql.DB
}

// MODIFIÉ : Le constructeur accepte maintenant un pointeur *sql.DB.
func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{DB: db}
}

// Create insère un nouveau produit en base et retourne l'ID généré.
// MODIFIÉ : On utilise `s.DB` au lieu de `database.DB`.
func (s *ProductService) Create(p models.Topic) (int, error) {
	result, err := s.DB.Exec(
		`INSERT INTO products (name, description, price, categorie_id) VALUES (?, ?, ?, ?)`,
		p.Name, p.Description, p.Price, p.CategorieId,
	)
	if err != nil {
		return 0, fmt.Errorf("ProductService.Create – Exec error: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ProductService.Create – LastInsertId error: %w", err)
	}

	return int(id), nil
}

// ReadAll récupère tous les produits de la table.
// MODIFIÉ : On utilise `s.DB` au lieu de `database.DB`.
func (s *ProductService) ReadAll() ([]models.Topic, error) {
	rows, err := s.DB.Query(`SELECT id, name, description, price, categorie_id FROM products`)
	if err != nil {
		return nil, fmt.Errorf("ProductService.ReadAll – Query error: %w", err)
	}
	defer rows.Close()

	products := []models.Topic{}
	for rows.Next() {
		var p models.Topic
		if scanErr := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategorieId); scanErr != nil {
			return nil, fmt.Errorf("ProductService.ReadAll – Scan error: %w", scanErr)
		}
		products = append(products, p)
	}

	return products, nil
}

// ... le reste des méthodes (ReadById, etc.) doit aussi utiliser `s.DB` ...
func (s *ProductService) ReadById(id int) (models.Topic, error) {
	var p models.Topic
	// MODIFIÉ : On utilise `s.DB`
	err := s.DB.QueryRow(
		`SELECT id, name, description, price, categorie_id FROM products WHERE id = ?`,
		id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategorieId)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("ProductService.ReadById – pas de produit pour l'id %d", id)
		}
		return p, fmt.Errorf("ProductService.ReadById – Scan error: %w", err)
	}
	return p, nil
}
