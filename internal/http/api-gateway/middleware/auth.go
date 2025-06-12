package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SlashLight/todo-list/internal/domain/models"
	"github.com/SlashLight/todo-list/internal/lib/jwt"
)

func AuthMiddleware(next http.Handler, secretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := headerParts[1]

		sess, err := jwt.ParseToken(tokenString, secretKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		ctx := models.ContextWithSession(context.Background(), sess)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
