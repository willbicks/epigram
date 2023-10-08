package service

import (
	"context"
	"testing"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/model"
)

func Test_notSignedIn(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
	}{
		{
			"Background context",
			context.Background(),
			ErrNotAuthenticated,
		},
		{
			"Context with empty user",
			ctxval.ContextWithUser(context.Background(), model.User{}),
			ErrNotAuthenticated,
		},
		{
			"Context with valid user",
			ctxval.ContextWithUser(context.Background(), model.User{
				ID: "f000",
			}),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifySignedIn(tt.ctx)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("verifySignedIn() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil && tt.wantErr.Error() != err.Error() {
				t.Errorf("verifySignedIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_noUserPrivilege(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
	}{
		{
			"Background context",
			context.Background(),
			ErrNotAuthorized,
		},
		{
			"Context with empty user",
			ctxval.ContextWithUser(context.Background(), model.User{}),
			ErrNotAuthorized,
		},
		{
			"Context with valid non-privileged user",
			ctxval.ContextWithUser(context.Background(), model.User{
				ID: "f000",
			}),
			ErrNotAuthorized,
		},
		{
			"Context with quiz passed user",
			ctxval.ContextWithUser(context.Background(), model.User{
				ID:         "f000",
				QuizPassed: true,
			}),
			nil,
		},
		{
			"Context with admin user",
			ctxval.ContextWithUser(context.Background(), model.User{
				ID:    "f000",
				Admin: true,
			}),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyUserPrivilege(tt.ctx)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("verifyUserPrivilege() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil && tt.wantErr.Error() != err.Error() {
				t.Errorf("verifyUserPrivilege() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_noAdminPrivilege(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr error
	}{
		{
			"Background context",
			context.Background(),
			ErrNotAuthorized,
		},
		{
			"Context with empty user",
			ctxval.ContextWithUser(context.Background(), model.User{}),
			ErrNotAuthorized,
		},
		{
			"Context with valid non-privileged user",
			ctxval.ContextWithUser(context.Background(), model.User{
				ID: "f000",
			}),
			ErrNotAuthorized,
		},
		{
			"Context with quiz passed user",
			ctxval.ContextWithUser(context.Background(), model.User{
				ID:         "f000",
				QuizPassed: true,
			}),
			ErrNotAuthorized,
		},
		{
			"Context with admin user",
			ctxval.ContextWithUser(context.Background(), model.User{
				ID:    "f000",
				Admin: true,
			}),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyAdminPrivilege(tt.ctx)
			if (tt.wantErr == nil) != (err == nil) {
				t.Errorf("verifyAdminPrivilege() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != nil && tt.wantErr.Error() != err.Error() {
				t.Errorf("verifyAdminPrivilege() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
