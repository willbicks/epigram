package service

import (
	"context"

	"github.com/willbicks/epigram/internal/ctxval"
)

// ErrNotAuhtenticated is returned when a request which requires an authentication is attempted without it.
var ErrNotAuhtenticated = Error{
	Issues:     []string{"Request requires authentication."},
	StatusCode: 401,
}

// ErrNotAuhthorized is returned when a request which requires an Auhthorized is attempted without it.
var ErrNotAuhthorized = Error{
	Issues:     []string{"Request requires authorization."},
	StatusCode: 401,
}

// verifySignedIn returns ErrNotAuhtenticated if the context lacks an authenticated user
func verifySignedIn(ctx context.Context) error {
	if u := ctxval.UserFromContext(ctx); u.ID == "" {
		return ErrNotAuhtenticated
	}
	return nil
}

// verifyUserPrivlege returns ErrNotAuthroized if the user on the Context is not authorized
func verifyUserPrivlege(ctx context.Context) error {
	if !ctxval.UserFromContext(ctx).IsAuthorized() {
		return ErrNotAuhthorized
	}
	return nil
}

// verifyAdminPrivlege returns ErrNotAuthroized if the user on the Context is not an admin
func verifyAdminPrivlege(ctx context.Context) error {
	if !ctxval.UserFromContext(ctx).IsAdmin() {
		return ErrNotAuhthorized
	}
	return nil
}
