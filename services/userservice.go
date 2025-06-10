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

var jwtKey = []byte("VOTRE_CLE_SECRETE_ULTRA_SECURISEE") // Changez ceci et mettez-le dans vos variables d'environnement !

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// Register gère la logique d'inscription (FT-1)
func (s *UserService) Register(user *models.User) error {
	//  Vérifier l'unicité du nom d'utilisateur et de l'email
	var exists int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM Utilisateurs WHERE name = ? OR email = ?", user.Name, user.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de l'unicité: %w", err)
	}
	if exists > 0 {
		return errors.New("le nom d'utilisateur ou l'email est déjà utilisé")
	}

	//  Valider le mot de passe
	if err := validatePassword(user.Password); err != nil {
		return err
	}

	//  Hacher le mot de passe avec bcrypt (plus sûr que SHA512)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erreur lors du hachage du mot de passe: %w", err)
	}
	user.Password = string(hashedPassword)

	// Insérer l'utilisateur en base de données
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

// Login gère la logique de connexion (FT-2)
func (s *UserService) Login(identifier, password string) (string, error) {
	var user models.User
	//  Récupérer l'utilisateur par nom ou email
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

	//  Comparer le mot de passe haché avec celui fourni
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("identifiants invalides")
	}

	//  Générer le token JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Name,
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

// validatePassword vérifie la complexité du mot de passe
func validatePassword(password string) error {
	//  Doit faire au moins 12 caractères
	if len(password) < 12 {
		return errors.New("le mot de passe doit contenir au moins 12 caractères")
	}
	//  Doit contenir au moins une majuscule
	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return errors.New("le mot de passe doit contenir au moins une majuscule")
	}
	//  Doit contenir au moins un caractère spécial
	if matched, _ := regexp.MatchString(`[\W_]`, password); !matched {
		return errors.New("le mot de passe doit contenir au moins un caractère spécial")
	}
	return nil
}
