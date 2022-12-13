package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/willbicks/epigram/internal/model"
)

const (
	// _idRandBytes represents the number of cryptographically secure random bytes that should be
	// generated for each UserSession's ID, and serves as it's unique token for user authenticaiton.
	// Multiples of 3 are preferred to make optimal usage of the base64 encoding scheme, but not
	// required.
	_idRandBytes = 18

	// DefaultExpirty represents the default ammount of time after which a UserSession will expire.
	_defaultExpiry = time.Hour * 24 * 14
)

// UserSessionRepository provides methods for storing and retriving UserSessions.
type UserSessionRepository interface {
	Create(ctx context.Context, us model.UserSession) error
	FindByID(ctx context.Context, id string) (model.UserSession, error)
}

// UserSession is a service for managing UserSessions.
type UserSession struct {
	repo UserSessionRepository
}

// NewUserSessionService returns a new UserSession service with the provided UserSessionRepository.
func NewUserSessionService(repo UserSessionRepository) UserSession {
	return UserSession{
		repo,
	}
}

// CreateUserSession creates a new UserSession for the provided User and returns it.
func (s UserSession) CreateUserSession(ctx context.Context, u model.User, IP string) (model.UserSession, error) {
	session := model.UserSession{}

	if u.ID == "" {
		return model.UserSession{}, errors.New("userSession: specified user has no id")
	}
	session.UserID = u.ID

	randBytes := make([]byte, _idRandBytes)
	if _, err := rand.Read(randBytes); err != nil {
		return model.UserSession{}, fmt.Errorf("generate randBytes for UserSession: %w", err)
	}
	session.ID = base64.URLEncoding.EncodeToString(randBytes)

	session.Created = time.Now()
	session.Expires = session.Created.Add(_defaultExpiry)
	session.IP = IP

	return session, s.repo.Create(ctx, session)
}

// FindSessionByID returns the UserSession with the specified ID, if it exists and is not expired.
func (s UserSession) FindSessionByID(ctx context.Context, id string) (model.UserSession, error) {
	session, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.UserSession{}, fmt.Errorf("UserSession: %w", err)
	}

	if session.IsExpired(time.Now()) {
		return model.UserSession{}, errors.New("UserSession is expired")
	}

	return session, nil
}
