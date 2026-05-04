package httpapi

import (
	"net/http"
	"strconv"

	"todo-api/internal/todo"
)

func (s *Server) handleListTodos(w http.ResponseWriter, r *http.Request) {
	q, vErr := parseListQuery(r)
	if vErr != nil {
		WriteErrorFromErr(w, vErr)
		return
	}

	items, total, err := s.svc.List(q)
	if err != nil {
		WriteErrorFromErr(w, err)
		return
	}

	dtos := make([]todoDTO, 0, len(items))
	for _, t := range items {
		dtos = append(dtos, toTodoDTO(t))
	}
	WriteJSON(w, http.StatusOK, listResponse[todoDTO]{
		Data: dtos,
		Meta: meta{
			Limit:  q.Limit,
			Offset: q.Offset,
			Count:  total,
		},
	})
}

func (s *Server) handleGetTodo(w http.ResponseWriter, r *http.Request, id string) {
	t, err := s.svc.Get(id)
	if err != nil {
		WriteErrorFromErr(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, dataResponse[todoDTO]{Data: toTodoDTO(t)})
}

func parseListQuery(r *http.Request) (todo.ListQuery, *todo.ValidationError) {
	const (
		defaultLimit = 50
		maxLimit     = 100
	)

	raw := r.URL.Query()
	var q todo.ListQuery

	if v := raw.Get("completed"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return todo.ListQuery{}, &todo.ValidationError{
				Message: "validation failed",
				Details: []todo.FieldIssue{{Field: "completed", Issue: "invalid"}},
			}
		}
		q.Completed = &b
	}

	q.Limit = defaultLimit
	if v := raw.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return todo.ListQuery{}, &todo.ValidationError{
				Message: "validation failed",
				Details: []todo.FieldIssue{{Field: "limit", Issue: "invalid"}},
			}
		}
		q.Limit = n
	}
	if q.Limit < 1 || q.Limit > maxLimit {
		return todo.ListQuery{}, &todo.ValidationError{
			Message: "validation failed",
			Details: []todo.FieldIssue{{Field: "limit", Issue: "invalid"}},
		}
	}

	if v := raw.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return todo.ListQuery{}, &todo.ValidationError{
				Message: "validation failed",
				Details: []todo.FieldIssue{{Field: "offset", Issue: "invalid"}},
			}
		}
		q.Offset = n
	}

	q.Sort = todo.SortCreatedAt
	if v := raw.Get("sort"); v != "" {
		q.Sort = todo.SortField(v)
	}
	switch q.Sort {
	case todo.SortCreatedAt, todo.SortUpdatedAt, todo.SortDueAt, todo.SortTitle:
	default:
		return todo.ListQuery{}, &todo.ValidationError{
			Message: "validation failed",
			Details: []todo.FieldIssue{{Field: "sort", Issue: "invalid"}},
		}
	}

	q.Order = todo.OrderDesc
	if v := raw.Get("order"); v != "" {
		q.Order = todo.SortOrder(v)
	}
	switch q.Order {
	case todo.OrderAsc, todo.OrderDesc:
	default:
		return todo.ListQuery{}, &todo.ValidationError{
			Message: "validation failed",
			Details: []todo.FieldIssue{{Field: "order", Issue: "invalid"}},
		}
	}

	return q, nil
}
