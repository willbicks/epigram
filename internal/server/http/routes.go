package http

import "net/http"

var paths = struct {
	home           string
	quiz           string
	login          string
	googleCallback string
}{
	home:  "/",
	quiz:  "/quiz",
	login: "/login",
}

func (s *CharismsServer) routes() {
	s.mux.Handle(paths.home, requireQuizPassed(http.HandlerFunc(s.homeHandler)))
	s.mux.Handle(paths.quiz, requireLoggedIn(http.HandlerFunc(s.quizHandler)))

	s.mux.HandleFunc(paths.login, s.googleLoginHandler)
	s.mux.HandleFunc("/login/google/callback", s.googleCallbackHandler)

	s.mux.Handle("/static/", s.staticHandler())
}

func (s *CharismsServer) staticHandler() http.Handler {
	// TODO: modify to disable directory listing
	return http.StripPrefix("/static/", http.FileServer(http.FS(s.PubFS)))
}
