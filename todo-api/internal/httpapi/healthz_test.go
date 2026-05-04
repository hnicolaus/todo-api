package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/internal/todo"
)

func TestHealthz(t *testing.T) {
	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })
	srv := NewServer(svc)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type=%q want=%q", got, "application/json")
	}
}
