package models

// User représente un utilisateur de la base de données.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"` // Le champ s'appelle bien "Name"
	Email    string `json:"email"`
	Password string `json:"-"`
	RoleID   int    `json:"role_id"`
	Username string // Ce champ n'est pas utilisé, mais pour éviter les erreurs
}
