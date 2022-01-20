package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// googleOIDC is a struct containing the Oauth2 Config and OIDC IDTokenVerifier required to validate OIDC callbacks.
type googleOIDC struct {
	config   oauth2.Config
	provider *oidc.Provider
}

// randString is a helper function used by googleOIDC to generate random strings for state and nonce.
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// setCallbackCookie is a helper function used by googleOIDC to set secure cookes containing state and nonce information.
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

// newGoogleOIDC returns a new googleOIDC object.
func newGoogleOIDC() (googleOIDC, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return googleOIDC{}, fmt.Errorf("count not create oidc provider: %w", err)
	}

	config := oauth2.Config{
		ClientID:     viper.GetString("googleOIDC.clientID"),
		ClientSecret: viper.GetString("googleOIDC.clientSecret"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  viper.GetString("baseURL") + "/login/google/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return googleOIDC{
		config:   config,
		provider: provider,
	}, nil
}

// verifier returns and OIDC IDTokenVerifier to verify OIDC tokens.
func (gOIDC *googleOIDC) verifier() *oidc.IDTokenVerifier {
	return gOIDC.provider.Verifier(&oidc.Config{ClientID: gOIDC.config.ClientID})
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

	http.Redirect(w, r, s.gOIDC.config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

func (s CharismsServer) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := s.gOIDC.config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}
	idToken, err := s.gOIDC.verifier().Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		http.Error(w, "nonce not found", http.StatusBadRequest)
		return
	}
	if idToken.Nonce != nonce.Value {
		http.Error(w, "nonce did not match", http.StatusBadRequest)
		return
	}

	oauth2Token.AccessToken = "*REDACTED*"

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
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
