package inmemory

import (
	"context"
	"github.com/willbicks/charisms/internal/model"
	"github.com/willbicks/charisms/internal/service"
	storage "github.com/willbicks/charisms/internal/storage/common"
	"sync"
)

type UserSessionRepository struct {
	sync.Mutex
	m map[string]model.UserSession
}

func NewUserSessionRepository() service.UserSessionRepository {
	return &UserSessionRepository{
		m: make(map[string]model.UserSession, 0),
	}
}

func (r *UserSessionRepository) Create(ctx context.Context, us model.UserSession) error {
	r.Lock()
	defer r.Unlock()

	r.m[us.ID] = us
	return nil
}

func (r *UserSessionRepository) FindByID(ctx context.Context, id string) (model.UserSession, error) {
	r.Lock()
	defer r.Unlock()

	session, ok := r.m[id]
	if !ok {
		return model.UserSession{}, storage.ErrNotFound
	}

	return session, nil
}
