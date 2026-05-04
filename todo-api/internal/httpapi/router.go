package httpapi

import (
	"net/http"
	"strings"
)

func (s *Server) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleCreateTodo(w, r)
	case http.MethodGet:
		s.handleListTodos(w, r)
	default:
		methodNotAllowed(w, []string{http.MethodPost, http.MethodGet})
	}
}

func (s *Server) handleTodoByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "invalid_path", "missing id", nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleGetTodo(w, r, id)
	case http.MethodPut:
		s.handleReplaceTodo(w, r, id)
	case http.MethodPatch:
		s.handlePatchTodo(w, r, id)
	case http.MethodDelete:
		s.handleDeleteTodo(w, r, id)
	default:
		methodNotAllowed(w, []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete})
	}
}

func methodNotAllowed(w http.ResponseWriter, allow []string) {
	w.Header().Set("Allow", strings.Join(allow, ", "))
	WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", nil)
}
