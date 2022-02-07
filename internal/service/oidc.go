package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OIDC is a struct containing the Oauth2 Config and OIDC IDTokenVerifier required to validate OIDC callbacks.
type OIDC struct {
	// Name is a unique identifier used by this OIDC service to build a callback url.
	Name string

	// IssuerURL is the URL for OIDC endpoint discovery.
	IssuerURL string

	ClientID     string
	ClientSecret string

	config   oauth2.Config
	provider *oidc.Provider
}

// RedirectURL returns the complete redirect URL (including provided state and nonce) to launch the oauth flow with the selected provider.
func (o OIDC) RedirectURL(state string, nonce string) (url string) {
	return o.config.AuthCodeURL(state, oidc.Nonce(nonce))
}

// CallbackURL returns the partial URL to be used for this OIDC service.
func (o OIDC) CallbackURL() string {
	return "/login/" + o.Name + "/callback"
}

// NewOIDC returns a new OIDC object with the specified issuer. Requires the baseURL of this server in order to build an oauth Redirect URL.
func (o *OIDC) Init(baseURL string) error {
	if baseURL == "" {
		return errors.New("baseURL is required to generate oauth callback url")
	}

	if o.IssuerURL == "" || o.ClientID == "" || o.ClientSecret == "" {
		return errors.New("at least one required field (IssuerURL, ClientID, ClientSecret) is missing from OIDC object")
	}

	provider, err := oidc.NewProvider(context.Background(), o.IssuerURL)
	if err != nil {
		return fmt.Errorf("could not create oidc provider: %w", err)
	}
	o.provider = provider

	o.config = oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  baseURL + o.CallbackURL(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return nil
}

// ValidateCallback accepts a callback http request, exchanges it for an oauth token, verifies that it contains a valid OIDC token, and then returns it.
func (o OIDC) ValidateCallback(r http.Request) (oidc.IDToken, error) {
	state, err := r.Cookie("state")
	if err != nil {
		return oidc.IDToken{}, ServiceError{
			StatusCode: http.StatusBadRequest,
			Issues:     []string{"State cookie not found."},
		}
	}
	if r.URL.Query().Get("state") != state.Value {
		return oidc.IDToken{}, ServiceError{
			StatusCode: http.StatusBadRequest,
			Issues:     []string{"State values do not match."},
		}
	}

	oauth2Token, err := o.config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		return oidc.IDToken{}, ServiceError{
			StatusCode: http.StatusInternalServerError,
			Issues:     []string{"Failed to exchange token: " + err.Error()},
		}
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return oidc.IDToken{}, ServiceError{
			StatusCode: http.StatusBadRequest,
			Issues:     []string{"No id_token field in oauth2 token."},
		}
	}

	v := o.provider.Verifier(&oidc.Config{ClientID: o.ClientID})
	idToken, err := v.Verify(r.Context(), rawIDToken)
	if err != nil {
		return oidc.IDToken{}, ServiceError{
			StatusCode: http.StatusInternalServerError,
			Issues:     []string{"Failed to verify ID Token: " + err.Error()},
		}
	}

	nonce, err := r.Cookie("nonce")
	if err != nil {
		return oidc.IDToken{}, ServiceError{
			StatusCode: http.StatusBadRequest,
			Issues:     []string{"Nonce cookie not found."},
		}
	}
	if idToken.Nonce != nonce.Value {
		return oidc.IDToken{}, ServiceError{
			StatusCode: http.StatusBadRequest,
			Issues:     []string{"Nonce values do not patch."},
		}
	}

	return *idToken, nil
}
