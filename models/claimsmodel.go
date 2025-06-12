package models

import "github.com/golang-jwt/jwt/v5"

// Claims définit les informations stockées dans le token JWT.
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	RoleID   int    `json:"role_id"` // Ce champ est essentiel
	jwt.RegisteredClaims
}
