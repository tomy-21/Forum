package middleware

import (
	"Forum/models"
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// IMPORTANT : Cette clé doit être EXACTEMENT la même que celle dans votre `userservice.go`
// Pensez à la déplacer dans un fichier de configuration ou une variable d'environnement.
var jwtKey = []byte("VOTRE_CLE_SECRETE_ULTRA_SECURISEE")

// AuthMiddleware vérifie le token JWT pour protéger une route.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le cookie contenant le token
		c, err := r.Cookie("token")
		if err != nil {
			// Si le cookie n'existe pas, on redirige vers la page de connexion
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			// Pour toute autre erreur de cookie, on retourne une erreur
			http.Error(w, "Requête invalide", http.StatusBadRequest)
			return
		}

		// Extraire la chaîne de caractères du token
		tokenStr := c.Value
		claims := &models.Claims{}

		// Parser et valider le token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Si le token est invalide ou expiré, on redirige vers la page de connexion
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Le token est valide ! On enrichit la requête avec les informations de l'utilisateur.
		// On stocke l'ID et le nom de l'utilisateur dans le "contexte" de la requête.
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		// On passe la requête modifiée au prochain handler (le contrôleur final)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
