package loginuser

import (
	"context"
	"encoding/json"
	"errors"
	appErrors "main/internal/core/errors"
	"main/internal/features/dto"
	"net/http"

	"go.uber.org/zap"
)

type LoginUser interface {
	LoginUser(ctx context.Context, user dto.UserDTO) (string, error)
}

func New(login LoginUser, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userDTO dto.UserDTO
		if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
			logger.Error("failed to decode request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := login.LoginUser(r.Context(), userDTO)
		if err != nil {
			if errors.Is(err, appErrors.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(token)); err != nil {
			logger.Error("Failed to write in response", zap.Error(err))
		}
	}
}
