package application

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

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

	resp := struct {
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{new(json.RawMessage)}

	if err := token.Claims(&resp.IDTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
