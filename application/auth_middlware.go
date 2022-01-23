package application

import "net/http"

// interpretSession wraps the request's context with the authenticated user, if they are known.
// Otherwise, execution passes to the next handler.
func (s *CharismsServer) interpretSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieName)
		if err != nil {
			// unable to read cookie (probably not set / no session)
			next.ServeHTTP(w, r)
		}

		u, err := s.UserService.GetUserFromSessionID(r.Context(), c.Value)
		if err != nil {
			// session token is invalid
			next.ServeHTTP(w, r)
		}

		ctx := ContextWithUser(r.Context(), u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireQuizPassed requires that the user has passed the entry quiz before proceeding, and if not,
// returns HTTP 403: Forbidden.
func requireQuizPassed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := UserFromContext(r.Context())
		if u.QuizPassed {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "You are not authorized to access this resource.", http.StatusForbidden)
		}
	})
}
