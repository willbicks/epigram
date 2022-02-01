package application

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"
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

// googleLoginHandler adds cookeies for state and nonce, and then redirects the user to Google sign in.
func (s CharismsServer) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := randString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	nonce, err := randString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	setCallbackCookie(w, r, "state", state)
	setCallbackCookie(w, r, "nonce", nonce)

	http.Redirect(w, r, s.gOIDC.RedirectURL(state, nonce), http.StatusFound)
}

func (s CharismsServer) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := s.gOIDC.ValidateCallback(*r)
	if err != nil {
		// TODO: use error handler helper to use ServiceError StatusCode when available
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err := s.UserService.GetUserFromIDToken(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sess, err := s.UserService.CreateUserSession(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sess.ID,
		Path:     "/",
		Secure:   r.TLS != nil,
		HttpOnly: true,
		// Session expires on client one hour before server to account for sync differences.
		Expires: sess.Expires.Add(-time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
