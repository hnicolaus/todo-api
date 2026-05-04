package todo

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrValidation     = errors.New("validation error")
	ErrNotImplemented = errors.New("not implemented")
)

type FieldIssue struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

type ValidationError struct {
	Message string
	Details []FieldIssue
}

func (e *ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "validation failed"
}

func (e *ValidationError) Unwrap() error { return ErrValidation }
