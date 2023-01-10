package inmemory

import (
	"context"
	"sync"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage"
)

// QuoteRepository is an in-memory implementation of the service.QuoteRepository interface.
type QuoteRepository struct {
	mu sync.RWMutex
	m  map[string]model.Quote
}

// NewQuoteRepository returns a new QuoteRepository which stores Quotes in memory.
func NewQuoteRepository() service.QuoteRepository {
	return &QuoteRepository{
		m: make(map[string]model.Quote, 0),
	}
}

// Create adds a new Quote to the repository.
func (r *QuoteRepository) Create(ctx context.Context, q model.Quote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[q.ID] = q
	return nil
}

// Update updates an existing Quote in the repository.
func (r *QuoteRepository) Update(ctx context.Context, q model.Quote) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.m[q.ID]

	if !ok {
		return storage.ErrNotFound
	}

	r.m[q.ID] = q
	return nil
}

// FindByID returns a Quote with the provided ID.
func (r *QuoteRepository) FindByID(ctx context.Context, id string) (model.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	q, ok := r.m[id]
	if !ok {
		return model.Quote{}, storage.ErrNotFound
	}

	return q, nil
}

// FindAll returns all Quotes in the repository.
func (r *QuoteRepository) FindAll(ctx context.Context) ([]model.Quote, error) {
	v := make([]model.Quote, 0, len(r.m))

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, value := range r.m {
		v = append(v, value)
	}

	return v, nil
}

// Delete removes a Quote with the provided ID from the repository.
func (r *QuoteRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.m[id]
	if !ok {
		return storage.ErrNotFound
	}

	delete(r.m, id)
	return nil
}
