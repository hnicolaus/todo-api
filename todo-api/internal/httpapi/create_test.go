package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"todo-api/internal/todo"
)

func TestCreateTodo_201_LocationAndBody(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"  hello  ","description":"  world  "}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusCreated, rr.Body.String())
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type=%q want=%q", got, "application/json")
	}
	loc := rr.Header().Get("Location")
	if !strings.HasPrefix(loc, "/api/v1/todos/") || loc == "/api/v1/todos/" {
		t.Fatalf("location=%q want prefix=%q", loc, "/api/v1/todos/")
	}

	var resp struct {
		Data todoDTO `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v body=%s", err, rr.Body.String())
	}
	if resp.Data.ID == "" {
		t.Fatalf("id empty")
	}
	if loc != "/api/v1/todos/"+resp.Data.ID {
		t.Fatalf("location=%q want=%q", loc, "/api/v1/todos/"+resp.Data.ID)
	}
	if resp.Data.Title != "hello" {
		t.Fatalf("title=%q want=%q", resp.Data.Title, "hello")
	}
	if resp.Data.Description != "world" {
		t.Fatalf("description=%q want=%q", resp.Data.Description, "world")
	}
	if resp.Data.Completed != false {
		t.Fatalf("completed=%v want=%v", resp.Data.Completed, false)
	}
	if resp.Data.DueAt != nil {
		t.Fatalf("dueAt=%v want=nil", *resp.Data.DueAt)
	}
	if resp.Data.CreatedAt != "1970-01-01T00:00:00Z" {
		t.Fatalf("createdAt=%q want=%q", resp.Data.CreatedAt, "1970-01-01T00:00:00Z")
	}
	if resp.Data.UpdatedAt != "1970-01-01T00:00:00Z" {
		t.Fatalf("updatedAt=%q want=%q", resp.Data.UpdatedAt, "1970-01-01T00:00:00Z")
	}
}

func TestCreateTodo_415_MissingOrWrongContentType(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, time.Now)
	srv := NewServer(svc)

	t.Run("missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"x"}`))
		rr := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rr, req)
		if rr.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusUnsupportedMediaType, rr.Body.String())
		}
	})

	t.Run("wrong", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"x"}`))
		req.Header.Set("Content-Type", "text/plain")
		rr := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rr, req)
		if rr.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusUnsupportedMediaType, rr.Body.String())
		}
	})
}

func TestCreateTodo_400_InvalidJSON(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, time.Now)
	srv := NewServer(svc)

	t.Run("malformed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusBadRequest, rr.Body.String())
		}
		var resp errorResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Error.Code != "invalid_json" {
			t.Fatalf("code=%q want=%q", resp.Error.Code, "invalid_json")
		}
	})

	t.Run("unknown_field", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"x","nope":1}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusBadRequest, rr.Body.String())
		}
		var resp errorResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Error.Code != "invalid_json" {
			t.Fatalf("code=%q want=%q", resp.Error.Code, "invalid_json")
		}
	})
}

func TestCreateTodo_422_TitleValidation(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, time.Now)
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"   "}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusUnprocessableEntity, rr.Body.String())
	}
	var resp errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v body=%s", err, rr.Body.String())
	}
	if resp.Error.Code != "validation_error" {
		t.Fatalf("code=%q want=%q", resp.Error.Code, "validation_error")
	}
	if len(resp.Error.Details) != 1 || resp.Error.Details[0].Field != "title" || resp.Error.Details[0].Issue != "required" {
		t.Fatalf("details=%v want title/required", resp.Error.Details)
	}
}

func TestCreateTodo_DueAtRFC3339(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBufferString(`{"title":"x","dueAt":"2026-05-04T12:34:56+07:00"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status=%d want=%d body=%s", rr.Code, http.StatusCreated, rr.Body.String())
	}
	var resp struct {
		Data todoDTO `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v body=%s", err, rr.Body.String())
	}
	if resp.Data.DueAt == nil {
		t.Fatalf("dueAt=nil want non-nil")
	}
	if *resp.Data.DueAt != "2026-05-04T05:34:56Z" {
		t.Fatalf("dueAt=%q want=%q", *resp.Data.DueAt, "2026-05-04T05:34:56Z")
	}
}
