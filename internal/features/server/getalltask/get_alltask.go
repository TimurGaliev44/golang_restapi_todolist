package getalltask

import (
	"context"
	"encoding/json"
	"main/internal/features/middleware"
	"main/internal/features/models"
	"net/http"

	"go.uber.org/zap"
)

type GetAllTask interface {
	SelectAllTask(ctx context.Context, userID int) ([]models.TaskModel, error)
}

func New(get GetAllTask, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tasks, err := get.SelectAllTask(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.MarshalIndent(tasks, "", "   ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Failed to convert in []byte", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			logger.Error("Failed to write in response", zap.Error(err))
		}
	}
}
