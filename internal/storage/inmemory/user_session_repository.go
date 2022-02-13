package inmemory

import (
	"context"
	"sync"

	"github.com/willbicks/charisms/internal/model"
	"github.com/willbicks/charisms/internal/service"
	storage "github.com/willbicks/charisms/internal/storage/common"
)

type UserSessionRepository struct {
	mu sync.RWMutex
	m  map[string]model.UserSession
}

func NewUserSessionRepository() service.UserSessionRepository {
	return &UserSessionRepository{
		m: make(map[string]model.UserSession, 0),
	}
}

func (r *UserSessionRepository) Create(ctx context.Context, us model.UserSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.m[us.ID] = us
	return nil
}

func (r *UserSessionRepository) FindByID(ctx context.Context, id string) (model.UserSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, ok := r.m[id]
	if !ok {
		return model.UserSession{}, storage.ErrNotFound
	}

	return session, nil
}
