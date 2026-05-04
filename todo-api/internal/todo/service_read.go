package todo

func (s *Service) Get(id string) (Todo, error) {
	return s.repo.Get(id)
}

func (s *Service) List(q ListQuery) ([]Todo, int, error) {
	nq, vErr := normalizeListQuery(q)
	if vErr != nil {
		return nil, 0, vErr
	}
	return s.repo.List(nq)
}

func normalizeListQuery(q ListQuery) (ListQuery, *ValidationError) {
	const (
		defaultLimit = 50
		maxLimit     = 100
	)

	if q.Limit == 0 {
		q.Limit = defaultLimit
	}
	if q.Limit < 1 {
		return ListQuery{}, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "limit", Issue: "invalid"}},
		}
	}
	if q.Limit > maxLimit {
		return ListQuery{}, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "limit", Issue: "too_large"}},
		}
	}
	if q.Offset < 0 {
		return ListQuery{}, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "offset", Issue: "invalid"}},
		}
	}
	if q.Sort == "" {
		q.Sort = SortCreatedAt
	}
	switch q.Sort {
	case SortCreatedAt, SortUpdatedAt, SortDueAt, SortTitle:
	default:
		return ListQuery{}, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "sort", Issue: "invalid"}},
		}
	}
	if q.Order == "" {
		q.Order = OrderDesc
	}
	switch q.Order {
	case OrderAsc, OrderDesc:
	default:
		return ListQuery{}, &ValidationError{
			Message: "validation failed",
			Details: []FieldIssue{{Field: "order", Issue: "invalid"}},
		}
	}

	return q, nil
}
