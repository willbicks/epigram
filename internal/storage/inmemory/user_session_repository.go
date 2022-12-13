package inmemory

import (
	"context"
	"sync"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage"
)

// UserSessionRepository is an in-memory implementation of the service.UserSessionRepository interface.
type UserSessionRepository struct {
	mu sync.RWMutex
	m  map[string]model.UserSession
}

// NewUserSessionRepository returns a new UserSessionRepository which stores UserSessions in memory.
func NewUserSessionRepository() service.UserSessionRepository {
	return &UserSessionRepository{
		m: make(map[string]model.UserSession, 0),
	}
}

// Create adds a new UserSession to the repository.
func (r *UserSessionRepository) Create(ctx context.Context, us model.UserSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.m[us.ID] = us
	return nil
}

// FindByID returns the UserSession with the provided ID.
func (r *UserSessionRepository) FindByID(ctx context.Context, id string) (model.UserSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.m[id]
	if !ok {
		return model.UserSession{}, storage.ErrNotFound
	}

	return session, nil
}
