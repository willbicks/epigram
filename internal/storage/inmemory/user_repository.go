package inmemory

import (
	"context"
	"sync"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
	storage "github.com/willbicks/epigram/internal/storage/common"
)

type UserRepository struct {
	mu sync.RWMutex
	m  map[string]model.User
}

func NewUserRepository() service.UserRepository {
	return &UserRepository{
		m: make(map[string]model.User, 0),
	}
}

func (r *UserRepository) Create(ctx context.Context, u model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.m[u.ID]; ok {
		return storage.ErrAlreadyExists
	}

	r.m[u.ID] = u
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.m[u.ID]

	if !ok {
		return storage.ErrNotFound
	}

	r.m[u.ID] = u
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.m[id]
	if !ok {
		return model.User{}, storage.ErrNotFound
	}

	return u, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	v := make([]model.User, 0, len(r.m))

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, value := range r.m {
		v = append(v, value)
	}

	return v, nil
}
