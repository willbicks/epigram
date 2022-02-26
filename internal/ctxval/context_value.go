// Package ctxval is used by the server and service tiers to store request-scoped
// informtaion that is ancillary to the request being made (eg: authenticated user,
// request IP)
package ctxval

import (
	"context"

	"github.com/willbicks/epigram/internal/model"
)

type contextKey int

const userKey contextKey = 0
const ipKey contextKey = 1

func ContextWithUser(ctx context.Context, u model.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func UserFromContext(ctx context.Context) model.User {
	u, ok := ctx.Value(userKey).(model.User)
	if !ok {
		return model.User{}
	}
	return u
}

func ContextWithIP(ctx context.Context, IP string) context.Context {
	return context.WithValue(ctx, ipKey, IP)
}

func IPFromContext(ctx context.Context) string {
	ip, ok := ctx.Value(ipKey).(string)
	if !ok {
		return ""
	}
	return ip
}
