package todo

import "time"

type ReplaceInput struct {
	Title       *string
	Description *string
	Completed   *bool
	DueAt       **string // nil means missing; non-nil points to value (which can be nil to clear)
}

type PatchInput struct {
	Title       *string
	Description *string
	Completed   *bool
	DueAt       **string // nil means absent; non-nil points to value (which can be nil to clear)
}

func (s *Service) Replace(id string, input ReplaceInput) (Todo, error) {
	var missing []FieldIssue
	if input.Title == nil {
		missing = append(missing, FieldIssue{Field: "title", Issue: "required"})
	}
	if input.Description == nil {
		missing = append(missing, FieldIssue{Field: "description", Issue: "required"})
	}
	if input.Completed == nil {
		missing = append(missing, FieldIssue{Field: "completed", Issue: "required"})
	}
	if input.DueAt == nil {
		missing = append(missing, FieldIssue{Field: "dueAt", Issue: "required"})
	}
	if len(missing) > 0 {
		return Todo{}, &ValidationError{
			Message: "validation failed",
			Details: missing,
		}
	}

	title, vErr := ValidateTitle(*input.Title)
	if vErr != nil {
		return Todo{}, vErr
	}
	desc, vErr := ValidateDescription(*input.Description)
	if vErr != nil {
		return Todo{}, vErr
	}
	dueAt, vErr := ParseDueAtRFC3339("dueAt", *input.DueAt)
	if vErr != nil {
		return Todo{}, vErr
	}

	now := s.now()
	return s.repo.Replace(now, id, Todo{
		Title:       title,
		Description: desc,
		Completed:   *input.Completed,
		DueAt:       dueAt,
	})
}

func (s *Service) Patch(id string, input PatchInput) (Todo, error) {
	if input.Title == nil && input.Description == nil && input.Completed == nil && input.DueAt == nil {
		return Todo{}, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "body", Issue: "required"}},
		}
	}

	var patch TodoPatch
	if input.Title != nil {
		title, vErr := ValidateTitle(*input.Title)
		if vErr != nil {
			return Todo{}, vErr
		}
		patch.Title = &title
	}
	if input.Description != nil {
		desc, vErr := ValidateDescription(*input.Description)
		if vErr != nil {
			return Todo{}, vErr
		}
		patch.Description = &desc
	}
	if input.Completed != nil {
		patch.Completed = input.Completed
	}
	if input.DueAt != nil {
		if *input.DueAt == nil {
			var clear *time.Time
			patch.DueAt = &clear
		} else {
			dueAt, vErr := ParseDueAtRFC3339("dueAt", *input.DueAt)
			if vErr != nil {
				return Todo{}, vErr
			}
			patch.DueAt = &dueAt
		}
	}

	now := s.now()
	return s.repo.Patch(now, id, patch)
}
