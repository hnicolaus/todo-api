package httpapi

import (
	"net/http"

	"todo-api/internal/todo"
)

func (s *Server) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	if !RequireJSON(w, r) {
		return
	}

	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Completed   *bool   `json:"completed"`
		DueAt       *string `json:"dueAt"`
	}
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON", nil)
		return
	}

	completed := false
	if req.Completed != nil {
		completed = *req.Completed
	}

	created, err := s.svc.Create(todo.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		Completed:   completed,
		DueAt:       req.DueAt,
	})
	if err != nil {
		WriteErrorFromErr(w, err)
		return
	}

	w.Header().Set("Location", "/api/v1/todos/"+created.ID)
	WriteJSON(w, http.StatusCreated, dataResponse[todoDTO]{Data: toTodoDTO(created)})
}
