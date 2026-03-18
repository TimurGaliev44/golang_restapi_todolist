package createtask

import (
	"context"
	"encoding/json"
	"errors"
	appErrors "main/internal/core/errors"
	"main/internal/features/dto"
	"main/internal/features/middleware"
	"main/internal/features/models"
	"net/http"

	"go.uber.org/zap"
)

type CreateTask interface {
	CreateTask(ctx context.Context, userID int, task dto.TaskDTO) (*models.TaskModel, error)
}

func New(create CreateTask, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var taskDTO dto.TaskDTO
		if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
			logger.Error("failed to decode request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := taskDTO.ValidateToCreateTask(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		task, err := create.CreateTask(r.Context(), userID, taskDTO)
		if err != nil {
			if errors.Is(err, appErrors.ErrTaskAlreadyExists) {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				logger.Error("failed to create task", zap.Error(err))
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

		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(b); err != nil {
			logger.Error("Failed to write in response", zap.Error(err))
		}
	}
}
