package http

import (
	"io/fs"
	"net/http"
)

// routes initializes the mux in the server struct with all of application routes
func (s *QuoteServer) routes(pubFS fs.FS) {
	paths := s.config.paths

	s.mux.Handle("/favicon.ico", http.FileServer(http.FS(pubFS)))

	s.mux.Handle(paths.Home, http.HandlerFunc(s.homeHandler))
	s.mux.Handle(paths.Quotes, s.requireQuizPassed(http.HandlerFunc(s.quotesHandler)))
	s.mux.Handle(paths.Quiz, s.requireLoggedIn(http.HandlerFunc(s.quizHandler)))

	// TODO: factor out into registerOIDCService(service.OIDC) method to prepare
	// for multiple OIDC providers
	s.mux.Handle(paths.Login, s.oidcLoginHandler(s.gOIDC))
	s.mux.Handle(s.gOIDC.CallbackURL(), s.oidcCallbackHandler(s.gOIDC))

	s.mux.Handle("/static/", s.staticHandler(pubFS))
}

// staticHandler accepts a file system containging files that should be publicly
// available, and returns a handler which serves them less the `/static/` prefix.
func (s *QuoteServer) staticHandler(fileSys fs.FS) http.Handler {
	// TODO: modify to disable directory listing
	return http.StripPrefix("/static/", http.FileServer(http.FS(fileSys)))
}
