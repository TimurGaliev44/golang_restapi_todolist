package middleware

import (
	"context"
	"errors"
	appErrors "main/internal/core/errors"
	"main/internal/core/jwt"
	"main/internal/features/dto"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			errDTO := dto.CreateNewError(err)
			if errors.Is(err, appErrors.ErrInvalidJWT) {
				http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
				return
			}
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	},
	)
}
