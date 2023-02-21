package http

import (
	"io/fs"
	"net/http"
)

// routes initializes the mux in the server struct with all application routes
func (s *QuoteServer) routes(pubFS fs.FS) {
	s.mux.Handle("/favicon.ico", http.FileServer(http.FS(pubFS)))

	s.mux.Handle(s.paths.Home, http.HandlerFunc(s.homeHandler))
	s.mux.Handle(s.paths.Quotes, s.requireQuizPassed(http.HandlerFunc(s.quotesHandler)))
	s.mux.Handle(s.paths.Quiz, s.requireLoggedIn(http.HandlerFunc(s.quizHandler)))

	s.mux.Handle(s.paths.Admin, s.requireLoggedIn(s.requireAdmin(http.HandlerFunc(s.adminMainHandler))))

	// TODO: factor out into registerOIDCService(service.OIDC) method to prepare
	// for multiple OIDC providers
	s.mux.Handle(s.paths.Login, s.oidcLoginHandler(s.OIDCService))
	s.mux.Handle(s.OIDCService.CallbackURL(), s.oidcCallbackHandler(s.OIDCService))

	s.mux.Handle(s.paths.Privacy, http.HandlerFunc(s.privacyHandler))
	s.mux.Handle("/static/", s.staticHandler(pubFS))
}
