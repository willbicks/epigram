package http

import "net/http"

func (s *QuoteServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path != "/" {
			s.notFoundError(w, r)
			return
		}
		err := s.renderPage(w, "home.gohtml", nil)
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
