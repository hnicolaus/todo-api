package httpapi

import (
	"net/http"

	"todo-api/internal/todo"
)

type Server struct {
	svc *todo.Service
	mux *http.ServeMux
}

func NewServer(svc *todo.Service) *Server {
	s := &Server{svc: svc, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealthz)
	s.mux.HandleFunc("/api/v1/todos", s.handleTodos)
	s.mux.HandleFunc("/api/v1/todos/", s.handleTodoByID)
}
