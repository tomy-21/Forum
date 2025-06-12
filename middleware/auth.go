package middleware

import (
	"net/http"
)

// La clé secrète est déclarée ICI une seule fois pour tout le package middleware.
var jwtKey = []byte("VOTRE_CLE_SECRETE_ULTRA_SECURISEE")

// AuthMiddleware vérifie simplement si un userID existe dans le contexte.
// Si ce n'est pas le cas, il redirige vers la page de connexion.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Le contexte est censé avoir déjà été peuplé par PopulateContextMiddleware.
		userID := r.Context().Value("userID")

		if userID == nil {
			// Pas d'ID, donc l'utilisateur n'est pas authentifié.
			http.Redirect(w, r, "/login?error=auth", http.StatusSeeOther)
			return
		}

		// L'utilisateur est authentifié, on continue vers la page protégée.
		next.ServeHTTP(w, r)
	})
}
