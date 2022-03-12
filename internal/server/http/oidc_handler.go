package http

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/service"
)

const sessionCookieName = "sess"

// oidcLoginHandler generates state and nonce keys, adds them to the client, and redirects to the
// oidc provider for authentication
func (s *QuoteServer) oidcLoginHandler(oidc service.OIDC) http.Handler {

	// randString is a helper function used by OIDC to generate random strings for state and nonce.
	randString := func(nByte int) (string, error) {
		b := make([]byte, nByte)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return "", err
		}
		return base64.RawURLEncoding.EncodeToString(b), nil
	}

	// setCallbackCookie is a helper function used to create secure cookes containing state and nonce information.
	setCallbackCookie := func(w http.ResponseWriter, r *http.Request, name, value string) {
		c := &http.Cookie{
			Name:     name,
			Value:    value,
			MaxAge:   int(time.Hour.Seconds()),
			Secure:   r.TLS != nil,
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}

	// http handler func to set coookies and redirect to provider
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// first, check if user is already signed in. If so, redirect to the quotes page.
		if ctxval.UserFromContext(r.Context()).ID != "" {
			http.Redirect(w, r, s.paths.Quotes, http.StatusFound)
			return
		}

		state, err := randString(16)
		if err != nil {
			s.serverError(w, r, fmt.Errorf("unable to generate state key: %v", err))
			return
		}
		nonce, err := randString(16)
		if err != nil {
			s.serverError(w, r, fmt.Errorf("unable to generate nonce key: %v", err))
			return
		}
		setCallbackCookie(w, r, "state", state)
		setCallbackCookie(w, r, "nonce", nonce)

		http.Redirect(w, r, oidc.RedirectURL(state, nonce), http.StatusFound)
	})
}

// oidcCallbackHandler handles callbacks from the OIDC provider
func (s *QuoteServer) oidcCallbackHandler(oidc service.OIDC) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := oidc.ValidateCallback(*r)
		if err != nil {
			s.serverError(w, r, fmt.Errorf("validating callback: %v", err))
			return
		}

		user, err := s.UserService.GetUserFromIDToken(r.Context(), token)
		if err != nil {
			s.serverError(w, r, fmt.Errorf("getting user from OIDC token: %v", err))
			return
		}

		ip := ctxval.IPFromContext(r.Context())

		sess, err := s.UserService.CreateUserSession(r.Context(), user, ip)
		if err != nil {
			s.serverError(w, r, fmt.Errorf("creating user session: %v", err))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   sessionCookieName,
			Value:  sess.ID,
			Path:   "/",
			Secure: r.TLS != nil,
			//SameSite: http.SameSiteStrictMode, // breaks redirect from after oidc callback?
			HttpOnly: true,
			// Session expires on client one hour before server to account for sync differences.
			Expires: sess.Expires.Add(-time.Hour),
		})
		http.Redirect(w, r, s.paths.Quotes, http.StatusFound)
	})
}
