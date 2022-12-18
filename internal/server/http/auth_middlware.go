package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/willbicks/epigram/internal/ctxval"
)

// interpretSession wraps the request's context with the authenticated user, if they are known.
// Otherwise, execution passes to the next handler.
func (s *QuoteServer) interpretSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieName)
		if err != nil {
			// unable to read cookie (probably not set / no session)
			s.Logger.Debug("session cookie not found")
			next.ServeHTTP(w, r)
			return
		}

		u, err := s.UserService.GetUserFromSessionID(r.Context(), c.Value)
		if err != nil {
			// session token is invalid
			s.Logger.Warnf("unable to get user from session ID: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := ctxval.ContextWithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireLoggedIn requires that the request has a valid session which has been translated to a user.
// If the user is not logged in, they will be redirected to the login page.
func (s *QuoteServer) requireLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := ctxval.UserFromContext(r.Context())
		if u.ID != "" {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, s.paths.Login, http.StatusFound)
		}
	})
}

// requireQuizPassed requires that the user has passed the entry quiz before proceeding, and if not,
// redirects them to the quiz page.
func (s *QuoteServer) requireQuizPassed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := ctxval.UserFromContext(r.Context())
		if u.QuizPassed {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, s.paths.Quiz, http.StatusFound)
		}
	})
}

// requireAdmin requires that the user be an admin, otherwise a 401 error is returned.
func (s *QuoteServer) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := ctxval.UserFromContext(r.Context())
		if u.Admin {
			next.ServeHTTP(w, r)
		} else {
			s.clientError(w, r, errors.New("you are not authorized to access this page"), http.StatusUnauthorized)
		}
	})
}

// getIp gets the IP of client making the request (using Config.TrustProxy to determine whether
// to use the X-Forwarded-For header), and stores it in the context of the request passed to the
// next handler.
func (s *QuoteServer) getIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		if s.Config.TrustProxy {
			ip = strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
		} else {
			ip = r.RemoteAddr
		}

		ctx := ctxval.ContextWithIP(r.Context(), ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
