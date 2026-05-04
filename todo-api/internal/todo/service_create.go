package todo

type CreateInput struct {
	Title       string
	Description string
	Completed   bool
	DueAt       *string // RFC3339
}

func (s *Service) Create(input CreateInput) (Todo, error) {
	title, vErr := ValidateTitle(input.Title)
	if vErr != nil {
		return Todo{}, vErr
	}
	desc, vErr := ValidateDescription(input.Description)
	if vErr != nil {
		return Todo{}, vErr
	}
	dueAt, vErr := ParseDueAtRFC3339("dueAt", input.DueAt)
	if vErr != nil {
		return Todo{}, vErr
	}

	t := Todo{
		Title:       title,
		Description: desc,
		Completed:   input.Completed,
		DueAt:       dueAt,
	}
	return s.repo.Create(s.now(), t)
}
