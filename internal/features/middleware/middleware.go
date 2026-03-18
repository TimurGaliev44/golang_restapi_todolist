package middleware

import (
	"context"
	"main/internal/core/jwt"
	"main/internal/core/storage"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
)

type contextKey string

const ClaimsKey contextKey = "claims"
const UserIDKey contextKey = "user_id"

func Auth(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			b, err := storage.IsTokenBlocked(r.Context(), rdb, tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if b {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			claims, err := jwt.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		},
		)
	}
}
