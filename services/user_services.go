package services

import (
	"Forum/models" // REMPLACEZ forum_app par votre nom de module réel
	"database/sql"

	"golang.org/x/crypto/bcrypt" // Pour le hachage de mot de passe
)

type UserService struct {
	DB *sql.DB
}

// NewUserService crée un nouveau UserService
func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// CreateUser crée un nouvel utilisateur
func (s *UserService) CreateUser(payload *models.RegistrationPayload) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO users (username, email, password_hash, role) VALUES (?, ?, ?, 'user')"
	res, err := s.DB.Exec(query, payload.Username, payload.Email, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId() // Bonne pratique: vérifier l'erreur de LastInsertId
	if err != nil {
		return nil, err
	}

	// Retourne l'utilisateur créé
	createdUser := &models.User{
		ID:       id,
		Username: payload.Username,
		Email:    payload.Email,
		Role:     "user",
		// CreatedAt sera défini par défaut par la DB
	}
	return createdUser, nil
}

// Ajoutez d'autres méthodes : GetUserByEmail, AuthenticateUser, GetUserByID etc.
