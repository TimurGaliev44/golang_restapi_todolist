package server

import (
	"errors"
	"main/internal/features/middleware"
	"main/internal/features/server/completetask"
	"main/internal/features/server/createtask"
	"main/internal/features/server/createuser"
	"main/internal/features/server/deletetask"
	"main/internal/features/server/getalltask"
	"main/internal/features/server/gettask"
	"main/internal/features/server/getuncompleted"
	"main/internal/features/server/loginuser"
	"main/internal/features/todo"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HTTPServer struct {
	td     *todo.TODO
	logger *zap.Logger
}

func NewServer(td *todo.TODO, logger *zap.Logger) *HTTPServer {
	return &HTTPServer{td: td, logger: logger}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	public := router.PathPrefix("").Subrouter()

	public.Path("/register").Methods("POST").HandlerFunc(createuser.New(s.td, s.logger))
	public.Path("/login").Methods("POST").HandlerFunc(loginuser.New(s.td, s.logger))

	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth)
	protected.Path("/tasks").Methods("POST").HandlerFunc(createtask.New(s.td, s.logger))
	protected.Path("/tasks/{id}").Methods("GET").HandlerFunc(gettask.New(s.td, s.logger))
	protected.Path("/tasks").Methods("GET").HandlerFunc(getalltask.New(s.td, s.logger))
	protected.Path("/tasks").Methods("GET").Queries("completed", "true").HandlerFunc(getuncompleted.New(s.td, s.logger))
	protected.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(completetask.New(s.td, s.logger))
	protected.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(deletetask.New(s.td, s.logger))

	if err := http.ListenAndServe(":5050", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
