package http

import (
	"io/fs"
	"net/http"
)

// paths is a constsant, package scoped variable which stores the url paths to each page,
// and should be used in place of magic strings to represent routes.
var paths = struct {
	home  string
	quiz  string
	login string
}{
	home:  "/",
	quiz:  "/quiz",
	login: "/login",
}

// routes initializes the mux in the server struct with all of the desired routes.
func (s *CharismsServer) routes(pubFS fs.FS) {
	s.mux.Handle(paths.home, requireQuizPassed(http.HandlerFunc(s.homeHandler)))
	s.mux.Handle(paths.quiz, requireLoggedIn(http.HandlerFunc(s.quizHandler)))

	s.mux.HandleFunc(paths.login, s.googleLoginHandler)
	s.mux.HandleFunc(s.gOIDC.CallbackURL(), s.googleCallbackHandler)

	s.mux.Handle("/static/", s.staticHandler(pubFS))
}

// staticHandler accepts a file system containging files that should be publicly
// available, and returns a handler which serves them less the `/static/` prefix.
func (s *CharismsServer) staticHandler(fileSys fs.FS) http.Handler {
	// TODO: modify to disable directory listing
	return http.StripPrefix("/static/", http.FileServer(http.FS(fileSys)))
}
