package deletetask

import (
	"context"
	"main/internal/features/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type DeleteTask interface {
	DeleteTask(ctx context.Context, userID, id int) error
}

func New(delete DeleteTask, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			logger.Error("failed smth", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := delete.DeleteTask(r.Context(), userID, id); err != nil {
			logger.Error("failed smth", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
