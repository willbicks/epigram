package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/model"

	"github.com/rs/xid"
)

// QuoteRepository provides methods for storing, manipulating, and retrieving Quotes
type QuoteRepository interface {
	Create(ctx context.Context, q model.Quote) error
	Update(ctx context.Context, q model.Quote) error
	FindByID(ctx context.Context, id string) (model.Quote, error)
	FindAll(ctx context.Context) ([]model.Quote, error)
}

// Quote provides a service for interacting with Quotes
type Quote struct {
	repo QuoteRepository
}

// NewQuoteService returns a new QuoteService with the provided QuoteRepository
func NewQuoteService(repo QuoteRepository) Quote {
	return Quote{
		repo,
	}
}

func verifyQuote(q model.Quote) error {
	var err Error
	if strings.TrimSpace(q.Quote) == "" {
		err.addIssue("Quote must not be blank.")
	}
	if strings.TrimSpace(q.Quotee) == "" {
		err.addIssue("This quote must be attributed to someone.")
	}

	if err.HasIssues() {
		err.StatusCode = http.StatusBadRequest
		return err
	}
	return nil
}

// CreateQuote creates a new Quote, setting its ID, Created, and SubmitterID fields
func (s *Quote) CreateQuote(ctx context.Context, q *model.Quote) error {

	if err := verifyUserPrivilege(ctx); err != nil {
		return err
	}

	if err := verifyQuote(*q); err != nil {
		return err
	}

	q.ID = xid.New().String()
	q.Created = time.Now()
	q.SubmitterID = ctxval.UserFromContext(ctx).ID
	return s.repo.Create(ctx, *q)
}

// GetQuote returns the Quote with the specified ID
func (s *Quote) GetQuote(ctx context.Context, id string) (model.Quote, error) {
	if err := verifyUserPrivilege(ctx); err != nil {
		return model.Quote{}, err
	}

	return s.repo.FindByID(ctx, id)
}

// UpdateQuote updates the specified Quote
func (s *Quote) UpdateQuote(ctx context.Context, q model.Quote) error {
	if err := verifyUserPrivilege(ctx); err != nil {
		return err
	}

	if !q.Editable(ctxval.UserFromContext(ctx)) {
		return Error{
			StatusCode: http.StatusUnauthorized,
			Issues:     []string{"You do not have permission to edit this quote. Quotes can only be edited by their submitter within an hour of submission."},
		}
	}

	if err := verifyQuote(q); err != nil {
		return err
	}

	return s.repo.Update(ctx, q)
}

// GetAllQuotes returns all Quotes
func (s *Quote) GetAllQuotes(ctx context.Context) ([]model.Quote, error) {
	if err := verifyUserPrivilege(ctx); err != nil {
		return nil, err
	}

	return s.repo.FindAll(ctx)
}
