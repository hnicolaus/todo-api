package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/internal/todo"
)

func TestListTodosEmpty(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type=%q want=%q", got, "application/json")
	}

	var resp listResponse[todoDTO]
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data == nil {
		t.Fatalf("data is nil; want empty array")
	}
	if got := len(resp.Data); got != 0 {
		t.Fatalf("len(data)=%d want=0", got)
	}
	if resp.Meta.Limit != 50 || resp.Meta.Offset != 0 || resp.Meta.Count != 0 {
		t.Fatalf("meta=%+v want limit=50 offset=0 count=0", resp.Meta)
	}
}

func TestListTodosFilteredCompleted(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	now := time.Unix(1000, 0).UTC()
	if _, err := repo.Create(now, todo.Todo{Title: "a", Completed: true}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := repo.Create(now.Add(time.Second), todo.Todo{Title: "b", Completed: false}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := repo.Create(now.Add(2*time.Second), todo.Todo{Title: "c", Completed: true}); err != nil {
		t.Fatalf("create: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos?completed=true", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusOK)
	}

	var resp listResponse[todoDTO]
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got := len(resp.Data); got != 2 {
		t.Fatalf("len(data)=%d want=2", got)
	}
	for _, it := range resp.Data {
		if !it.Completed {
			t.Fatalf("expected all completed=true, got id=%s completed=%v", it.ID, it.Completed)
		}
	}
	if resp.Meta.Count != 2 {
		t.Fatalf("meta.count=%d want=2", resp.Meta.Count)
	}
}

func TestListTodosPagingDefaultSortCreatedAtDesc(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	t1, err := repo.Create(time.Unix(1, 0).UTC(), todo.Todo{Title: "oldest"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	t2, err := repo.Create(time.Unix(2, 0).UTC(), todo.Todo{Title: "middle"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	t3, err := repo.Create(time.Unix(3, 0).UTC(), todo.Todo{Title: "newest"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Default sort is createdAt desc, so order is [t3, t2, t1]. With limit=1, offset=1, we should get t2.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos?limit=1&offset=1", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusOK)
	}

	var resp listResponse[todoDTO]
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Meta.Limit != 1 || resp.Meta.Offset != 1 || resp.Meta.Count != 3 {
		t.Fatalf("meta=%+v want limit=1 offset=1 count=3", resp.Meta)
	}
	if got := len(resp.Data); got != 1 {
		t.Fatalf("len(data)=%d want=1", got)
	}
	if resp.Data[0].ID != t2.ID {
		t.Fatalf("got id=%s want=%s (createdAt order t3=%s t2=%s t1=%s)", resp.Data[0].ID, t2.ID, t3.ID, t2.ID, t1.ID)
	}
}

func TestGetTodoSuccess(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	created, err := repo.Create(time.Unix(10, 0).UTC(), todo.Todo{
		Title:       "hello",
		Description: "world",
		Completed:   true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos/"+created.ID, nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type=%q want=%q", got, "application/json")
	}

	var resp dataResponse[todoDTO]
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.ID != created.ID {
		t.Fatalf("id=%s want=%s", resp.Data.ID, created.ID)
	}
	if resp.Data.Title != "hello" || resp.Data.Description != "world" || resp.Data.Completed != true {
		t.Fatalf("data=%+v want title/description/completed match", resp.Data)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos/does-not-exist", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusNotFound)
	}
	var resp errorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Error.Code != "not_found" {
		t.Fatalf("error.code=%q want=%q", resp.Error.Code, "not_found")
	}
}
