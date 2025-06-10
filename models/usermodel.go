package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // Le '-' empêche le mot de passe d'être encodé en JSON
	RoleID   int    `json:"role_id"`
}
