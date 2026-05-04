package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/internal/todo"
)

func TestDeleteTodo_Success(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	now := func() time.Time { return time.Unix(0, 0).UTC() }
	svc := todo.NewService(repo, now)
	srv := NewServer(svc)

	created, err := repo.Create(now(), todo.Todo{Title: "t1"})
	if err != nil {
		t.Fatalf("create todo: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/"+created.ID, nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusNoContent)
	}
	if rr.Body.Len() != 0 {
		t.Fatalf("body=%q want empty", rr.Body.String())
	}
	if _, err := repo.Get(created.ID); err == nil {
		t.Fatalf("expected deleted todo to be missing")
	} else if err != todo.ErrNotFound {
		t.Fatalf("repo get error=%v want=%v", err, todo.ErrNotFound)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/does-not-exist", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusNotFound)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type=%q want=%q", got, "application/json")
	}
	var got errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Error.Code != "not_found" {
		t.Fatalf("error.code=%q want=%q", got.Error.Code, "not_found")
	}
	if got.Error.Message == "" {
		t.Fatalf("error.message want non-empty")
	}
}
