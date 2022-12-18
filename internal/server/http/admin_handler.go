package http

import (
	"net/http"

	"github.com/willbicks/epigram/internal/server/http/frontend"
)

// privacyHandler renders the privacy policy page in response to GET requests
func (s *QuoteServer) adminMainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, err := s.UserService.GetAllUsers(r.Context())
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		err = s.tmpl.RenderPage(w, frontend.AdminMainPage{
			Users: users,
		})
		if err != nil {
			s.serverError(w, r, err)
			return
		}
	default:
		s.methodNotAllowedError(w, r)
		return
	}
}
