package service_test

import (
	"context"
	"testing"
	"time"

	model2 "github.com/willbicks/charisms/internal/model"
	"github.com/willbicks/charisms/internal/service"
	inmemory "github.com/willbicks/charisms/internal/storage/inmemory"

	"github.com/matryer/is"
	"github.com/rs/xid"
)

func TestUserSessionService_CreateUserSession(t *testing.T) {
	is := is.New(t)

	userRepo := inmemory.NewUserRepository()
	sessionRepo := inmemory.NewUserSessionRepository()

	service := service.NewUserSessionService(sessionRepo)

	user := model2.User{
		ID:   xid.New().String(),
		Name: "Test user",
	}
	userRepo.Create(context.Background(), user)

	sess, err := service.CreateUserSession(context.Background(), user)
	is.NoErr(err)                                    // creating user session should not fail
	is.Equal(sess.UserID, user.ID)                   // user sesision UserID should match user's ID
	is.True(len(sess.ID) > 8)                        // session ID should be at least 8 characters
	is.True(time.Since(sess.Created) < time.Second)  // session should be created within the last second
	is.True(time.Until(sess.Expires) < 72*time.Hour) // session should expire within 72 hours
}

func TestUserSessionService_FindSessionByID_Valid(t *testing.T) {
	is := is.New(t)

	sessionRepo := inmemory.NewUserSessionRepository()

	service := service.NewUserSessionService(sessionRepo)

	user := model2.User{
		ID:   xid.New().String(),
		Name: "Test user",
	}

	validSess, _ := service.CreateUserSession(context.Background(), user)

	found, err := service.FindSessionByID(context.Background(), validSess.ID)
	is.NoErr(err)                   // lookup of valid session id should not fail
	is.Equal(found.UserID, user.ID) // found session's UserID should match users's ID
}

func TestUserSessionService_FindSessionByID_Invalid(t *testing.T) {
	is := is.New(t)

	sessionRepo := inmemory.NewUserSessionRepository()

	service := service.NewUserSessionService(sessionRepo)

	user := model2.User{
		ID:   xid.New().String(),
		Name: "Test user",
	}

	sessionRepo.Create(context.Background(), model2.UserSession{
		ID:      "ExPiReD000",
		UserID:  user.ID,
		Created: time.Now().Add(-25 * time.Hour),
		Expires: time.Now().Add(-1 * time.Hour),
	})

	_, err := service.FindSessionByID(context.Background(), "hXu02xmYGEkt5RTf6Z3gYCymW")
	is.True(err != nil) // lookup of invalid session id should return error

	_, err = service.FindSessionByID(context.Background(), "ExPiReD000")
	is.True(err != nil) // lookup of expired session id should return error
}