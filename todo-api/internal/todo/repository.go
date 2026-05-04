package todo

import "time"

type SortField string

const (
	SortCreatedAt SortField = "createdAt"
	SortUpdatedAt SortField = "updatedAt"
	SortDueAt     SortField = "dueAt"
	SortTitle     SortField = "title"
)

type SortOrder string

const (
	OrderAsc  SortOrder = "asc"
	OrderDesc SortOrder = "desc"
)

type ListQuery struct {
	Completed *bool
	Limit     int
	Offset    int
	Sort      SortField
	Order     SortOrder
}

type Repository interface {
	Create(now time.Time, t Todo) (Todo, error)
	Get(id string) (Todo, error)
	List(q ListQuery) ([]Todo, int, error)
	Replace(now time.Time, id string, t Todo) (Todo, error)
	Patch(now time.Time, id string, patch TodoPatch) (Todo, error)
	Delete(id string) error
}

type TodoPatch struct {
	Title       *string
	Description *string
	Completed   *bool
	DueAt       **time.Time // nil means "not provided"; non-nil points to value (which can be nil to clear)
}
