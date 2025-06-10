package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
