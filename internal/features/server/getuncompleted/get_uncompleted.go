package getuncompleted

import (
	"context"
	"encoding/json"
	"main/internal/features/middleware"
	"main/internal/features/models"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type GetUncompleted interface {
	SelectUncompletedTasks(ctx context.Context, userID int) ([]models.TaskModel, error)
}

func New(get GetUncompleted, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		models, err := get.SelectUncompletedTasks(r.Context(), userID)
		if err != nil {
			logger.Error("failed to select uncompleted tasks", zap.String("id-task", strconv.Itoa(userID)))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.MarshalIndent(models, "", "   ")
		if err != nil {
			logger.Error("Failed to convert in []byte", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			logger.Error("Failed to write in response", zap.Error(err))
		}
	}
}
