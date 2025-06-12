package middleware

import (
	"Forum/models"
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// PopulateContextMiddleware vérifie s'il y a un token valide et, si oui,
// enrichit le contexte de la requête avec les informations de l'utilisateur.
func PopulateContextMiddleware(next http.Handler) http.Handler {
	// La ligne "var jwtKey" est supprimée d'ici.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := c.Value
		claims := &models.Claims{}

		// La variable `jwtKey` est accessible car elle est dans le même package (déclarée dans auth.go)
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err == nil && token.Valid {
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "roleID", claims.RoleID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
