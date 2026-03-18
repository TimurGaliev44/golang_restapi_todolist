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
	"main/internal/features/server/logoutuser"
	"main/internal/features/todo"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HTTPServer struct {
	td *todo.TODO
}

func NewServer(td *todo.TODO) *HTTPServer {
	return &HTTPServer{td: td}
}

func (s *HTTPServer) StartServer(logger *zap.Logger) error {
	router := mux.NewRouter()

	public := router.PathPrefix("").Subrouter()

	public.Path("/register").Methods("POST").HandlerFunc(createuser.New(s.td, logger))
	public.Path("/login").Methods("POST").HandlerFunc(loginuser.New(s.td, logger))

	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth(s.td.GetRedis()))
	protected.Path("/logout").Methods("POST").HandlerFunc(logoutuser.New(s.td, logger))
	protected.Path("/tasks").Methods("POST").HandlerFunc(createtask.New(s.td, logger))
	protected.Path("/tasks/{id}").Methods("GET").HandlerFunc(gettask.New(s.td, logger))
	protected.Path("/tasks").Methods("GET").HandlerFunc(getalltask.New(s.td, logger))
	protected.Path("/tasks").Methods("GET").Queries("completed", "true").HandlerFunc(getuncompleted.New(s.td, logger))
	protected.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(completetask.New(s.td, logger))
	protected.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(deletetask.New(s.td, logger))

	if err := http.ListenAndServe(":5050", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
