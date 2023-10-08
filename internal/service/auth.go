package service

import (
	"context"

	"github.com/willbicks/epigram/internal/ctxval"
)

// ErrNotAuthenticated is returned when a request which requires an authentication is attempted without it.
var ErrNotAuthenticated = Error{
	Issues:     []string{"Request requires authentication."},
	StatusCode: 401,
}

// ErrNotAuthorized is returned when a request which requires an Authorized is attempted without it.
var ErrNotAuthorized = Error{
	Issues:     []string{"Request requires authorization."},
	StatusCode: 401,
}

// verifySignedIn returns ErrNotAuthenticated if the context lacks an authenticated user
func verifySignedIn(ctx context.Context) error {
	if u := ctxval.UserFromContext(ctx); u.ID == "" {
		return ErrNotAuthenticated
	}
	return nil
}

// verifyUserPrivilege returns ErrNotAuthorized if the user on the Context is not authorized
func verifyUserPrivilege(ctx context.Context) error {
	if !ctxval.UserFromContext(ctx).IsAuthorized() {
		return ErrNotAuthorized
	}
	return nil
}

// verifyAdminPrivilege returns ErrNotAuthorized if the user on the Context is not an admin
func verifyAdminPrivilege(ctx context.Context) error {
	if !ctxval.UserFromContext(ctx).IsAdmin() {
		return ErrNotAuthorized
	}
	return nil
}
