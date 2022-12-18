package http

import (
	"net/http"

	"github.com/willbicks/epigram/internal/server/http/frontend"
)

func (s *QuoteServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path != "/" {
			s.notFoundError(w, r)
			return
		}
		err := s.tmpl.RenderPage(w, frontend.HomePage{})
		if err != nil {
			s.serverError(w, r, err)
			return
		}
		return

	default:
		s.methodNotAllowedError(w, r)
		return
	}
}
