package completetask

import (
	"context"
	"encoding/json"
	"main/internal/features/middleware"
	"main/internal/features/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type CompleteTask interface {
	CompleteTask(ctx context.Context, userID, id int) (*models.TaskModel, error)
}

func New(complete CompleteTask, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task, err := complete.CompleteTask(r.Context(), userID, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		b, err := json.MarshalIndent(task, "", "   ")
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
