package httpapi

import "net/http"

func (s *Server) handleDeleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	if err := s.svc.Delete(id); err != nil {
		WriteErrorFromErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
