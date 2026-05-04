package httpapi

import (
	"time"

	"todo-api/internal/todo"
)

type todoDTO struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Completed   bool    `json:"completed"`
	DueAt       *string `json:"dueAt"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type dataResponse[T any] struct {
	Data T `json:"data"`
}

type meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type listResponse[T any] struct {
	Data []T  `json:"data"`
	Meta meta `json:"meta"`
}

func toTodoDTO(t todo.Todo) todoDTO {
	var dueAt *string
	if t.DueAt != nil {
		s := t.DueAt.UTC().Format(time.RFC3339)
		dueAt = &s
	}
	return todoDTO{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		DueAt:       dueAt,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
