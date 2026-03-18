package logoutuser

import (
	"context"
	"main/internal/core/jwt"
	"main/internal/features/middleware"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type LogoutUser interface {
	LogoutUser(ctx context.Context, token string, ttl time.Duration) error
}

func New(logout LogoutUser, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.ClaimsKey).(*jwt.Claims)
		ttl := time.Until(claims.ExpiresAt.Time)
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		err := logout.LogoutUser(r.Context(), token, ttl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("failed to log out", zap.Error(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
