package createuser

import (
	"context"
	"encoding/json"
	"errors"
	appErrors "main/internal/core/errors"
	"main/internal/features/dto"
	"net/http"

	"go.uber.org/zap"
)

type CreateUser interface {
	RegisterUser(ctx context.Context, user dto.UserDTO) error
}

func New(create CreateUser, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user dto.UserDTO
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Error("failed to decode request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := create.RegisterUser(r.Context(), user); err != nil {
			if errors.Is(err, appErrors.ErrUserAlreadyExists) {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				logger.Error("failed to register user", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
