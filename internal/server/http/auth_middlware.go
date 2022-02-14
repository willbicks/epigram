package http

import (
	"net/http"
)

// interpretSession wraps the request's context with the authenticated user, if they are known.
// Otherwise, execution passes to the next handler.
func (s *CharismsServer) interpretSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieName)
		if err != nil {
			// unable to read cookie (probably not set / no session)
			next.ServeHTTP(w, r)
			return
		}

		u, err := s.UserService.GetUserFromSessionID(r.Context(), c.Value)
		if err != nil {
			// session token is invalid
			next.ServeHTTP(w, r)
			return
		}

		ctx := ContextWithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireLoggedIn requires that the request has a valid session which has been translated to a user.
// If the user is not logged in, they will be redirected to the login page.
func requireLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r.Context())
		if u.ID != "" {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, paths.login, http.StatusFound)
		}
	})
}

// requireQuizPassed requires that the user has passed the entry quiz before proceeding, and if not,
// redirects them to the quiz page.
func requireQuizPassed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r.Context())
		if u.QuizPassed {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, paths.quiz, http.StatusFound)
		}
	})
}
