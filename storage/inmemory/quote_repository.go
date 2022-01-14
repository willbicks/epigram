package inmemory

import (
	"context"
	"sync"

	"github.com/willbicks/charisms/model"
	"github.com/willbicks/charisms/service"
	storagecommon "github.com/willbicks/charisms/storage/storage_common"
)

type QuoteRepository struct {
	sync.Mutex
	m map[string]model.Quote
}

func NewQuoteRepository() service.QuoteRepository {
	return &QuoteRepository{
		m: make(map[string]model.Quote, 0),
	}
}

func (r *QuoteRepository) Create(ctx context.Context, q model.Quote) error {
	r.Lock()
	defer r.Unlock()
	r.m[q.ID] = q
	return nil
}

func (r *QuoteRepository) Update(ctx context.Context, q model.Quote) error {
	_, ok := r.m[q.ID]

	if !ok {
		return storagecommon.ErrNotFound
	}

	r.Lock()
	defer r.Unlock()
	r.m[q.ID] = q
	return nil
}

func (r *QuoteRepository) FindByID(ctx context.Context, id string) (model.Quote, error) {
	r.Lock()
	defer r.Unlock()
	q, ok := r.m[id]
	if !ok {
		return model.Quote{}, storagecommon.ErrNotFound
	}

	return q, nil
}

func (r *QuoteRepository) FindAll(ctx context.Context) ([]model.Quote, error) {
	v := make([]model.Quote, len(r.m))

	r.Lock()
	defer r.Unlock()
	for _, value := range r.m {
		v = append(v, value)
	}

	return v, nil
}
