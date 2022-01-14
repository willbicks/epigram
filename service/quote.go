package service

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/willbicks/charisms/model"
)

type QuoteRepository interface {
	Create(ctx context.Context, q model.Quote) error
	Update(ctx context.Context, q model.Quote) error
	FindByID(ctx context.Context, id string) (model.Quote, error)
	FindAll(ctx context.Context) ([]model.Quote, error)
}

type QuoteService struct {
	repo QuoteRepository
}

func NewQuoteService(repo QuoteRepository) QuoteService {
	return QuoteService{
		repo,
	}
}

func (s QuoteService) CreateQuote(ctx context.Context, q *model.Quote) error {
	q.ID = xid.New().String()
	q.Created = time.Now()
	return s.repo.Create(ctx, *q)
}

func (s QuoteService) GetAllQuotes(ctx context.Context) ([]model.Quote, error) {
	quotes, err := s.repo.FindAll(ctx)
	return quotes, err
}