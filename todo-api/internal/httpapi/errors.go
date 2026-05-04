package httpapi

import (
	"errors"
	"mime"
	"net/http"

	"todo-api/internal/todo"
)

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details []todo.FieldIssue `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, code, message string, details []todo.FieldIssue) {
	WriteJSON(w, status, errorResponse{
		Error: apiError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func WriteErrorFromErr(w http.ResponseWriter, err error) {
	var vErr *todo.ValidationError
	switch {
	case errors.Is(err, todo.ErrNotFound):
		WriteError(w, http.StatusNotFound, "not_found", "to-do not found", nil)
	case errors.As(err, &vErr):
		WriteError(w, http.StatusUnprocessableEntity, "validation_error", vErr.Error(), vErr.Details)
	default:
		WriteError(w, http.StatusInternalServerError, "internal", "internal server error", nil)
	}
}

func RequireJSON(w http.ResponseWriter, r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		WriteError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json", nil)
		return false
	}
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil || mediaType != "application/json" {
		WriteError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json", nil)
		return false
	}
	return true
}
