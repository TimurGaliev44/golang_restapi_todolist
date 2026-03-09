package server

import (
	"encoding/json"
	"errors"
	appErrors "main/internal/core/errors"
	"main/internal/features/dto"
	"main/internal/features/middleware"
	"main/internal/features/todo"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HTTPHandlers struct {
	todo   *todo.TODO
	logger *zap.Logger
}

func NewHandlers(td *todo.TODO, logger *zap.Logger) *HTTPHandlers {
	return &HTTPHandlers{todo: td, logger: logger}
}

func (h *HTTPHandlers) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := user.ValidateToCreate(); err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	err := h.todo.RegisterUser(r.Context(), user)
	if err != nil {
		errDTO := dto.CreateNewError(err)
		if errors.Is(err, appErrors.ErrUserAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
			return
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPHandlers) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var userDTO dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := userDTO.ValidateToCreate(); err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	token, err := h.todo.LoginUser(r.Context(), userDTO)
	if err != nil {
		if errors.Is(err, appErrors.ErrUserNotFound) {
			errDTO := dto.CreateNewError(err)
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		}
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(token)); err != nil {
		h.logger.Error("Failed to write in response", zap.Error(err))
	}
}

func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var taskDTO dto.TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateToCreate(); err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	task, err := h.todo.CreateTask(r.Context(), userID, taskDTO)
	if err != nil {
		errDTO := dto.CreateNewError(err)
		if errors.Is(err, appErrors.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "   ")
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		h.logger.Error("Failed to convert in []byte", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write in response", zap.Error(err))
	}

}

func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	task, err := h.todo.SelectTask(r.Context(), userID, id)
	if err != nil {
		errDTO := dto.CreateNewError(err)
		if errors.Is(err, appErrors.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "   ")
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		h.logger.Error("Failed to convert in []byte", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write in response", zap.Error(err))
	}
}

func (h *HTTPHandlers) HandleGetAllTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tasks, err := h.todo.SelectAllTask(r.Context(), userID)
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(tasks, "", "   ")
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		h.logger.Error("Failed to convert in []byte", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write in response", zap.Error(err))
	}
}

func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	task, err := h.todo.CompleteTask(r.Context(), userID, id)
	if err != nil {
		errDTO := dto.CreateNewError(err)
		if errors.Is(err, appErrors.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(task, "", "   ")
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		h.logger.Error("Failed to convert in []byte", zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write in response", zap.Error(err))
	}
}

func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := h.todo.DeleteTask(r.Context(), userID, id); err != nil {
		errDTO := dto.CreateNewError(err)

		if errors.Is(err, appErrors.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandlers) HandleGetUncompletedTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	models, err := h.todo.SelectUncompletedTasks(r.Context(), userID)
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(models, "", "   ")
	if err != nil {
		errDTO := dto.CreateNewError(err)
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		h.logger.Error("Failed to convert in []byte", zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write in response", zap.Error(err))
	}
}
