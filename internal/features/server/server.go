package server

import (
	"errors"
	"main/internal/features/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewServer(h *HTTPHandlers) *HTTPServer {
	return &HTTPServer{httpHandlers: h}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	public := router.PathPrefix("").Subrouter()

	public.Path("/register").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateUser)
	public.Path("/login").Methods("POST").HandlerFunc(s.httpHandlers.HandleLoginUser)

	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth)

	protected.Path("/tasks").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateTask)
	protected.Path("/tasks/{id}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetTask)
	protected.Path("/tasks").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetAllTask)
	protected.Path("/tasks").Methods("GET").Queries("completed", "true").HandlerFunc(s.httpHandlers.HandleGetUncompletedTask)
	protected.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(s.httpHandlers.HandleCompleteTask)
	protected.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteTask)

	if err := http.ListenAndServe(":5050", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
