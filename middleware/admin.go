package middleware

import (
	"Forum/models"
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// AdminMiddleware vérifie si l'utilisateur est bien un administrateur.
func AdminMiddleware(next http.Handler) http.Handler {
	// La ligne "var jwtKey" est supprimée d'ici.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenStr := c.Value
		claims := &models.Claims{}

		// La variable `jwtKey` est accessible car elle est dans le même package (déclarée dans auth.go)
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if claims.RoleID != 1 {
			http.Error(w, "Accès interdit : vous n'êtes pas administrateur.", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "roleID", claims.RoleID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
