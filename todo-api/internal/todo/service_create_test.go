package todo

import (
	"errors"
	"testing"
	"time"
)

func TestServiceCreate_TrimsAndSetsDefaults(t *testing.T) {
	repo := NewInMemoryRepo()
	now := time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC)
	svc := NewService(repo, func() time.Time { return now })

	due := "2026-05-04T12:34:56+07:00"
	created, err := svc.Create(CreateInput{
		Title:       "  hello  ",
		Description: "  world  ",
		Completed:   true,
		DueAt:       &due,
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("id empty")
	}
	if created.Title != "hello" {
		t.Fatalf("title=%q want=%q", created.Title, "hello")
	}
	if created.Description != "world" {
		t.Fatalf("description=%q want=%q", created.Description, "world")
	}
	if created.Completed != true {
		t.Fatalf("completed=%v want=%v", created.Completed, true)
	}
	if created.DueAt == nil {
		t.Fatalf("dueAt nil")
	}
	if got := created.DueAt.UTC().Format(time.RFC3339); got != "2026-05-04T05:34:56Z" {
		t.Fatalf("dueAt=%q want=%q", got, "2026-05-04T05:34:56Z")
	}
	if got := created.CreatedAt.UTC().Format(time.RFC3339); got != now.Format(time.RFC3339) {
		t.Fatalf("createdAt=%q want=%q", got, now.Format(time.RFC3339))
	}
	if got := created.UpdatedAt.UTC().Format(time.RFC3339); got != now.Format(time.RFC3339) {
		t.Fatalf("updatedAt=%q want=%q", got, now.Format(time.RFC3339))
	}
}

func TestServiceCreate_InvalidDueAt(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, time.Now)

	due := "nope"
	_, err := svc.Create(CreateInput{Title: "x", DueAt: &due})
	var vErr *ValidationError
	if err == nil || !errors.As(err, &vErr) {
		t.Fatalf("err=%v want ValidationError", err)
	}
	if len(vErr.Details) != 1 || vErr.Details[0].Field != "dueAt" || vErr.Details[0].Issue != "invalid_format" {
		t.Fatalf("details=%v want dueAt/invalid_format", vErr.Details)
	}
}
