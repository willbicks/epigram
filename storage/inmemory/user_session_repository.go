package inmemory

import (
	"context"
	"sync"

	"github.com/willbicks/charisms/model"
	"github.com/willbicks/charisms/service"
	storagecommon "github.com/willbicks/charisms/storage/storage_common"
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
		return model.UserSession{}, storagecommon.ErrNotFound
	}

	return session, nil
}
