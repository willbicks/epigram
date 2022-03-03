package service

import (
	"context"
	"time"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/model"

	"github.com/rs/xid"
)

type QuoteRepository interface {
	Create(ctx context.Context, q model.Quote) error
	Update(ctx context.Context, q model.Quote) error
	FindByID(ctx context.Context, id string) (model.Quote, error)
	FindAll(ctx context.Context) ([]model.Quote, error)
}

type Quote struct {
	repo QuoteRepository
}

func NewQuoteService(repo QuoteRepository) Quote {
	return Quote{
		repo,
	}
}

func (s *Quote) CreateQuote(ctx context.Context, q *model.Quote) error {

	if err := verifyUserPrivlege(ctx); err != nil {
		return err
	}

	var err ServiceError
	if q.Quote == "" {
		err.addIssue("Quote must not be blank.")
	}
	if q.Quotee == "" {
		err.addIssue("This quote must be attributed to someone.")
	}

	if err.HasIssues() {
		return err
	}

	q.ID = xid.New().String()
	q.Created = time.Now()
	q.SubmitterID = ctxval.UserFromContext(ctx).ID
	return s.repo.Create(ctx, *q)
}

func (s *Quote) GetAllQuotes(ctx context.Context) ([]model.Quote, error) {
	if err := verifyUserPrivlege(ctx); err != nil {
		return nil, err
	}

	quotes, err := s.repo.FindAll(ctx)
	return quotes, err
}
