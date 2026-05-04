package todo

import "time"

type Todo struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	DueAt       *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
