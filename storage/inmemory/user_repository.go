package inmemory

import (
	"context"
	"sync"

	"github.com/willbicks/charisms/model"
	"github.com/willbicks/charisms/service"
	storagecommon "github.com/willbicks/charisms/storage/storage_common"
)

type UserRepository struct {
	sync.Mutex
	m map[string]model.User
}

func NewUserRepository() service.UserRepository {
	return &UserRepository{
		m: make(map[string]model.User, 0),
	}
}

func (r *UserRepository) Create(ctx context.Context, q model.User) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[q.ID]; ok {
		return storagecommon.ErrAllreadyExists
	}

	r.m[q.ID] = q
	return nil
}

func (r *UserRepository) Update(ctx context.Context, q model.User) error {
	r.Lock()
	defer r.Unlock()

	_, ok := r.m[q.ID]

	if !ok {
		return storagecommon.ErrNotFound
	}

	r.m[q.ID] = q
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	r.Lock()
	defer r.Unlock()

	q, ok := r.m[id]
	if !ok {
		return model.User{}, storagecommon.ErrNotFound
	}

	return q, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	v := make([]model.User, 0, len(r.m))

	r.Lock()
	defer r.Unlock()

	for _, value := range r.m {
		v = append(v, value)
	}

	return v, nil
}
