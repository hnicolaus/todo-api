package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/internal/todo"
)

type testErrorResponse struct {
	Error struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Details []todo.FieldIssue `json:"details,omitempty"`
	} `json:"error"`
}

func TestReplaceTodo_Success(t *testing.T) {
	repo := todo.NewInMemoryRepo()

	createdAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	existing, err := repo.Create(createdAt, todo.Todo{
		Title:       "old",
		Description: "old desc",
		Completed:   false,
		DueAt:       nil,
	})
	if err != nil {
		t.Fatalf("seed create: %v", err)
	}

	updatedAt := time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC)
	svc := todo.NewService(repo, func() time.Time { return updatedAt })
	srv := NewServer(svc)

	body := `{"title":"  new  ","description":"  desc  ","completed":true,"dueAt":null}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/todos/"+existing.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type=%q want=%q", got, "application/json")
	}

	var resp dataResponse[todoDTO]
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Data.ID != existing.ID {
		t.Fatalf("id=%q want=%q", resp.Data.ID, existing.ID)
	}
	if resp.Data.Title != "new" {
		t.Fatalf("title=%q want=%q", resp.Data.Title, "new")
	}
	if resp.Data.Description != "desc" {
		t.Fatalf("description=%q want=%q", resp.Data.Description, "desc")
	}
	if resp.Data.Completed != true {
		t.Fatalf("completed=%v want=%v", resp.Data.Completed, true)
	}
	if resp.Data.DueAt != nil {
		t.Fatalf("dueAt=%v want=nil", *resp.Data.DueAt)
	}
	if resp.Data.CreatedAt != createdAt.Format(time.RFC3339) {
		t.Fatalf("createdAt=%q want=%q", resp.Data.CreatedAt, createdAt.Format(time.RFC3339))
	}
	if resp.Data.UpdatedAt != updatedAt.Format(time.RFC3339) {
		t.Fatalf("updatedAt=%q want=%q", resp.Data.UpdatedAt, updatedAt.Format(time.RFC3339))
	}
	if resp.Data.UpdatedAt == resp.Data.CreatedAt {
		t.Fatalf("updatedAt should change from createdAt")
	}
}

func TestPatchTodo_Success(t *testing.T) {
	repo := todo.NewInMemoryRepo()

	createdAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	existing, err := repo.Create(createdAt, todo.Todo{
		Title:       "title",
		Description: "desc",
		Completed:   false,
		DueAt:       nil,
	})
	if err != nil {
		t.Fatalf("seed create: %v", err)
	}

	updatedAt := time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC)
	svc := todo.NewService(repo, func() time.Time { return updatedAt })
	srv := NewServer(svc)

	dueAt := "2026-06-01T00:00:00Z"
	body := `{"completed":true,"dueAt":"` + dueAt + `"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/"+existing.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var resp dataResponse[todoDTO]
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Data.Title != "title" {
		t.Fatalf("title=%q want=%q", resp.Data.Title, "title")
	}
	if resp.Data.Completed != true {
		t.Fatalf("completed=%v want=%v", resp.Data.Completed, true)
	}
	if resp.Data.DueAt == nil || *resp.Data.DueAt != dueAt {
		if resp.Data.DueAt == nil {
			t.Fatalf("dueAt=nil want=%q", dueAt)
		}
		t.Fatalf("dueAt=%q want=%q", *resp.Data.DueAt, dueAt)
	}
	if resp.Data.UpdatedAt != updatedAt.Format(time.RFC3339) {
		t.Fatalf("updatedAt=%q want=%q", resp.Data.UpdatedAt, updatedAt.Format(time.RFC3339))
	}
}

func TestReplaceTodo_ValidationMissingFields(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	createdAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	existing, err := repo.Create(createdAt, todo.Todo{Title: "t"})
	if err != nil {
		t.Fatalf("seed create: %v", err)
	}

	svc := todo.NewService(repo, func() time.Time { return time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC) })
	srv := NewServer(svc)

	body := `{"title":"x"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/todos/"+existing.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusUnprocessableEntity, rr.Body.String())
	}

	var resp testErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "validation_error" {
		t.Fatalf("code=%q want=%q", resp.Error.Code, "validation_error")
	}

	wantMissing := map[string]bool{"description": false, "completed": false, "dueAt": false}
	for _, d := range resp.Error.Details {
		if d.Issue != "required" {
			continue
		}
		if _, ok := wantMissing[d.Field]; ok {
			wantMissing[d.Field] = true
		}
	}
	for field, ok := range wantMissing {
		if !ok {
			t.Fatalf("missing required detail for field %q; details=%v", field, resp.Error.Details)
		}
	}
}

func TestPatchTodo_ValidationEmptyBody(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	createdAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	existing, err := repo.Create(createdAt, todo.Todo{Title: "t"})
	if err != nil {
		t.Fatalf("seed create: %v", err)
	}

	svc := todo.NewService(repo, func() time.Time { return time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC) })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/"+existing.ID, bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusUnprocessableEntity, rr.Body.String())
	}

	var resp testErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "validation_error" {
		t.Fatalf("code=%q want=%q", resp.Error.Code, "validation_error")
	}
	if len(resp.Error.Details) != 1 || resp.Error.Details[0].Field != "body" || resp.Error.Details[0].Issue != "required" {
		t.Fatalf("details=%v want body required", resp.Error.Details)
	}
}

func TestPatchTodo_NotFound(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC) })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/does-not-exist", bytes.NewBufferString(`{"completed":true}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}

	var resp testErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "not_found" {
		t.Fatalf("code=%q want=%q", resp.Error.Code, "not_found")
	}
}

func TestReplaceTodo_UnknownFieldRejected(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	createdAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	existing, err := repo.Create(createdAt, todo.Todo{Title: "t"})
	if err != nil {
		t.Fatalf("seed create: %v", err)
	}

	svc := todo.NewService(repo, func() time.Time { return time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC) })
	srv := NewServer(svc)

	body := `{"title":"x","description":"y","completed":false,"dueAt":null,"nope":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/todos/"+existing.ID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}

	var resp testErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error.Code != "invalid_json" {
		t.Fatalf("code=%q want=%q", resp.Error.Code, "invalid_json")
	}
}

func TestPatchTodo_UnsupportedMediaType(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	createdAt := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)
	existing, err := repo.Create(createdAt, todo.Todo{Title: "t"})
	if err != nil {
		t.Fatalf("seed create: %v", err)
	}

	svc := todo.NewService(repo, func() time.Time { return time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC) })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/"+existing.ID, bytes.NewBufferString(`{"completed":true}`))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusUnsupportedMediaType, rr.Body.String())
	}
}
