package http

import (
	"context"
	"github.com/willbicks/charisms/internal/model"
)

// TODO: Consider moving this to another package

type contextKey int

const contextUserKey contextKey = 0

func ContextWithUser(ctx context.Context, u model.User) context.Context {
	return context.WithValue(ctx, contextUserKey, u)
}

func UserFromContext(ctx context.Context) model.User {
	u, ok := ctx.Value(contextUserKey).(model.User)
	if !ok {
		return model.User{}
	}
	return u
}
