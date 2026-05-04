package httpapi

import (
	"net/http"

	"todo-api/internal/todo"
)

func (s *Server) handleReplaceTodo(w http.ResponseWriter, r *http.Request, id string) {
	if !RequireJSON(w, r) {
		return
	}

	type replaceRequest struct {
		Title       *string        `json:"title"`
		Description *string        `json:"description"`
		Completed   *bool          `json:"completed"`
		DueAt       OptionalString `json:"dueAt"`
	}

	var req replaceRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON", nil)
		return
	}

	var dueAt **string
	if req.DueAt.Set {
		dueAt = &req.DueAt.Value
	}

	updated, err := s.svc.Replace(id, todo.ReplaceInput{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
		DueAt:       dueAt,
	})
	if err != nil {
		WriteErrorFromErr(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dataResponse[todoDTO]{Data: toTodoDTO(updated)})
}

func (s *Server) handlePatchTodo(w http.ResponseWriter, r *http.Request, id string) {
	if !RequireJSON(w, r) {
		return
	}

	type patchRequest struct {
		Title       *string        `json:"title"`
		Description *string        `json:"description"`
		Completed   *bool          `json:"completed"`
		DueAt       OptionalString `json:"dueAt"`
	}

	var req patchRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid_json", "invalid JSON", nil)
		return
	}

	var dueAt **string
	if req.DueAt.Set {
		dueAt = &req.DueAt.Value
	}

	updated, err := s.svc.Patch(id, todo.PatchInput{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
		DueAt:       dueAt,
	})
	if err != nil {
		WriteErrorFromErr(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, dataResponse[todoDTO]{Data: toTodoDTO(updated)})
}
