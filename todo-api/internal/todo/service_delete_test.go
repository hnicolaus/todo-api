package todo

import (
	"testing"
	"time"
)

func TestServiceDelete_Success(t *testing.T) {
	repo := NewInMemoryRepo()
	now := func() time.Time { return time.Unix(0, 0).UTC() }
	svc := NewService(repo, now)

	created, err := repo.Create(now(), Todo{Title: "t1"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.Get(created.ID); err != ErrNotFound {
		t.Fatalf("get after delete err=%v want=%v", err, ErrNotFound)
	}
}

func TestServiceDelete_NotFound(t *testing.T) {
	repo := NewInMemoryRepo()
	svc := NewService(repo, func() time.Time { return time.Unix(0, 0).UTC() })

	if err := svc.Delete("missing"); err != ErrNotFound {
		t.Fatalf("err=%v want=%v", err, ErrNotFound)
	}
}
