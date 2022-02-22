package http

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/willbicks/epigram/internal/service"
)

const sessionCookieName = "sess"

// randString is a helper function used by OIDC to generate random strings for state and nonce.
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// setCallbackCookie is a helper function used to create secure cookes containing state and nonce information.
func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func (s QuoteServer) oidcLoginHandler(oidc service.OIDC) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func (s QuoteServer) oidcCallbackHandler(oidc service.OIDC) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := oidc.ValidateCallback(*r)
		if err != nil {
			// TODO: use error handler helper to use ServiceError StatusCode when available
			s.serverError(w, r, fmt.Errorf("validating callback: %v", err))
			return
		}

		user, err := s.UserService.GetUserFromIDToken(r.Context(), token)
		if err != nil {
			s.serverError(w, r, fmt.Errorf("getting user from OIDC token: %v", err))
			return
		}

		sess, err := s.UserService.CreateUserSession(r.Context(), user)
		if err != nil {
			s.serverError(w, r, fmt.Errorf("creating user session: %v", err))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sess.ID,
			Path:     "/",
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			// Session expires on client one hour before server to account for sync differences.
			Expires: sess.Expires.Add(-time.Hour),
		})
		http.Redirect(w, r, paths.home, http.StatusFound)
	})
}
