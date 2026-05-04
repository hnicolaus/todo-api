package todo

import (
	"strings"
	"time"
)

const (
	TitleMaxLen       = 200
	DescriptionMaxLen = 2000
)

func ValidateTitle(title string) (string, *ValidationError) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return "", &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "title", Issue: "required"}},
		}
	}
	if len(trimmed) > TitleMaxLen {
		return "", &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "title", Issue: "too_long"}},
		}
	}
	return trimmed, nil
}

func ValidateDescription(desc string) (string, *ValidationError) {
	trimmed := strings.TrimSpace(desc)
	if len(trimmed) > DescriptionMaxLen {
		return "", &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "description", Issue: "too_long"}},
		}
	}
	return trimmed, nil
}

func ParseDueAtRFC3339(field string, v *string) (*time.Time, *ValidationError) {
	if v == nil {
		return nil, nil
	}
	if *v == "" {
		return nil, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: field, Issue: "invalid_format"}},
		}
	}
	t, err := time.Parse(time.RFC3339, *v)
	if err != nil {
		return nil, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: field, Issue: "invalid_format"}},
		}
	}
	utc := t.UTC()
	return &utc, nil
}
