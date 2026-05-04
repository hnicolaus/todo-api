package todo

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
