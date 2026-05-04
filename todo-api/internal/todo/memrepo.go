package todo

import (
	"sort"
	"sync"
	"time"
)

type InMemoryRepo struct {
	mu    sync.RWMutex
	items map[string]Todo
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		items: make(map[string]Todo),
	}
}

func (r *InMemoryRepo) Create(now time.Time, t Todo) (Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, err := newID()
	if err != nil {
		return Todo{}, err
	}
	t.ID = id
	t.CreatedAt = now.UTC()
	t.UpdatedAt = t.CreatedAt
	r.items[t.ID] = t
	return t, nil
}

func (r *InMemoryRepo) Get(id string) (Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

func (r *InMemoryRepo) List(q ListQuery) ([]Todo, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all := make([]Todo, 0, len(r.items))
	for _, t := range r.items {
		if q.Completed != nil && t.Completed != *q.Completed {
			continue
		}
		all = append(all, t)
	}

	sort.Slice(all, func(i, j int) bool {
		// Compute ascending "less" and "equal" for the selected sort key.
		var less, equal bool
		switch q.Sort {
		case SortUpdatedAt:
			less = all[i].UpdatedAt.Before(all[j].UpdatedAt)
			equal = all[i].UpdatedAt.Equal(all[j].UpdatedAt)
		case SortDueAt:
			ai, aj := all[i].DueAt, all[j].DueAt
			switch {
			case ai == nil && aj == nil:
				less, equal = false, true
			case ai == nil && aj != nil:
				less = true
			case ai != nil && aj == nil:
				less = false
			default:
				less = ai.Before(*aj)
				equal = ai.Equal(*aj)
			}
		case SortTitle:
			less = all[i].Title < all[j].Title
			equal = all[i].Title == all[j].Title
		case SortCreatedAt:
			fallthrough
		default:
			less = all[i].CreatedAt.Before(all[j].CreatedAt)
			equal = all[i].CreatedAt.Equal(all[j].CreatedAt)
		}

		if equal {
			return all[i].ID < all[j].ID
		}
		if q.Order == OrderAsc {
			return less
		}
		return !less
	})

	total := len(all)
	start := q.Offset
	if start > total {
		start = total
	}
	end := start + q.Limit
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (r *InMemoryRepo) Replace(now time.Time, id string, t Todo) (Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	t.ID = id
	t.CreatedAt = existing.CreatedAt
	t.UpdatedAt = now.UTC()
	r.items[id] = t
	return t, nil
}

func (r *InMemoryRepo) Patch(now time.Time, id string, patch TodoPatch) (Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	if patch.Title != nil {
		t.Title = *patch.Title
	}
	if patch.Description != nil {
		t.Description = *patch.Description
	}
	if patch.Completed != nil {
		t.Completed = *patch.Completed
	}
	if patch.DueAt != nil {
		t.DueAt = *patch.DueAt
	}
	t.UpdatedAt = now.UTC()
	r.items[id] = t
	return t, nil
}

func (r *InMemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
