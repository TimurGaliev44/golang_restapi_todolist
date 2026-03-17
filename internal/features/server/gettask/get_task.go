package gettask

import (
	"context"
	"encoding/json"
	"errors"
	appErrors "main/internal/core/errors"
	"main/internal/features/middleware"
	"main/internal/features/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type GetTask interface {
	SelectTask(ctx context.Context, userID, id int) (*models.TaskModel, error)
}

func New(get GetTask, logger *zap.Logger) http.HandlerFunc {
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

		task, err := get.SelectTask(r.Context(), userID, id)
		if err != nil {
			if errors.Is(err, appErrors.ErrTaskNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				logger.Error("failed to get task", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
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
