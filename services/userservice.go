package services

import (
	"Forum/models"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("VOTRE_CLE_SECRETE_ULTRA_SECURISEE")

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// Register gère la logique d'inscription.
func (s *UserService) Register(user *models.User) error {
	var exists int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM Utilisateurs WHERE name = ? OR email = ?", user.Name, user.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de l'unicité: %w", err)
	}
	if exists > 0 {
		return errors.New("le nom d'utilisateur ou l'email est déjà utilisé")
	}

	if err := validatePassword(user.Password); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erreur lors du hachage du mot de passe: %w", err)
	}
	user.Password = string(hashedPassword)

	_, err = s.DB.Exec(
		"INSERT INTO Utilisateurs (role_id, name, email, password) VALUES (?, ?, ?, ?)",
		3, // Par défaut, rôle 'Utilisateur'
		user.Name,
		user.Email,
		user.Password,
	)
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion de l'utilisateur: %w", err)
	}
	return nil
}

// Login gère la connexion et la création du token JWT.
func (s *UserService) Login(identifier, password string) (string, error) {
	var user models.User
	err := s.DB.QueryRow(
		"SELECT user_id, name, email, password, role_id FROM Utilisateurs WHERE name = ? OR email = ?",
		identifier, identifier,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("identifiants invalides")
		}
		return "", fmt.Errorf("erreur lors de la récupération de l'utilisateur: %w", err)
	}
	if user.RoleID == 4 {
		return "", errors.New("Votre compte a été banni. Veuillez contacter un administrateur.")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("identifiants invalides")
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Name,
		RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération du token: %w", err)
	}
	return tokenString, nil
}

// validatePassword vérifie la complexité du mot de passe.
func validatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("le mot de passe doit contenir au moins 12 caractères")
	}
	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return errors.New("le mot de passe doit contenir au moins une majuscule")
	}
	if matched, _ := regexp.MatchString(`[\W_]`, password); !matched {
		return errors.New("le mot de passe doit contenir au moins un caractère spécial")
	}
	return nil
}

// GetAllUsers récupère tous les utilisateurs de la base de données.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	query := `SELECT user_id, name, email, role_id FROM Utilisateurs`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserByID récupère les informations d'un utilisateur par son ID.
func (s *UserService) GetUserByID(userID int) (models.User, error) {
	var u models.User
	query := "SELECT user_id, name, email, role_id FROM Utilisateurs WHERE user_id = ?"
	err := s.DB.QueryRow(query, userID).Scan(&u.ID, &u.Name, &u.Email, &u.RoleID)
	return u, err
}

// GetUserTopicCount retourne le nombre de sujets créés par un utilisateur.
func (s *UserService) GetUserTopicCount(userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM sujet WHERE user_id = ?"
	err := s.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

// GetUserMessageCount retourne le nombre de messages postés par un utilisateur.
func (s *UserService) GetUserMessageCount(userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM messages WHERE user_id = ?"
	err := s.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

// BanUser met à jour le rôle d'un utilisateur pour le bannir.
func (s *UserService) BanUser(userID int) error {
	// On change le role_id à 4 (Banni)
	_, err := s.DB.Exec("UPDATE Utilisateurs SET role_id = 4 WHERE user_id = ?", userID)
	return err
}

// --- MÉTHODE MANQUANTE À AJOUTER CI-DESSOUS ---

// UnbanUser met à jour le rôle d'un utilisateur pour le débannir.
func (s *UserService) UnbanUser(userID int) error {
	// On remet le role_id à 3 (Utilisateur)
	_, err := s.DB.Exec("UPDATE Utilisateurs SET role_id = 3 WHERE user_id = ?", userID)
	return err
}
