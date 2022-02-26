package service

import (
	"context"

	"github.com/willbicks/epigram/internal/ctxval"
)

var ErrNotAuhtenticated = ServiceError{
	Issues:     []string{"Request requires authentication."},
	StatusCode: 401,
}

var ErrNotAuhthorized = ServiceError{
	Issues:     []string{"Request requires authorixation."},
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

// noUserPrivlege returns ErrNotAuthroized if the user on the Context is not an admin
func verifyAdminPrivlege(ctx context.Context) error {
	if !ctxval.UserFromContext(ctx).IsAdmin() {
		return ErrNotAuhthorized
	}
	return nil
}
