package http

import "net/http"

func (s *CharismsServer) routes() {
	s.mux.Handle("/", requireQuizPassed(http.HandlerFunc(s.homeHandler)))
	s.mux.Handle("/quiz", requireLoggedIn(http.HandlerFunc(s.quizHandler)))

	s.mux.HandleFunc("/login", s.googleLoginHandler)
	s.mux.HandleFunc("/login/google/callback", s.googleCallbackHandler)

	s.mux.Handle("/static/", s.staticHandler())
}

func (s *CharismsServer) staticHandler() http.Handler {
	// TODO: modify to disable directory listing
	return http.StripPrefix("/static/", http.FileServer(http.FS(s.PubFS)))
}
