package validate

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage"
)

// QuoteRepository validates a type implementing the QuoteRepository interface
func QuoteRepository(t *testing.T, repoFactory func() (repo service.QuoteRepository, close func())) {
	t.Run("Create_FindByID", func(t *testing.T) {
		repo, close := repoFactory()
		defer close()
		t.Parallel()
		quoteRepository_Create_FindByID(t, repo)
	})

	t.Run("Update", func(t *testing.T) {
		repo, close := repoFactory()
		defer close()
		t.Parallel()
		quoteRepository_Update(t, repo)
	})

	t.Run("FindAll", func(t *testing.T) {
		repo, close := repoFactory()
		defer close()
		t.Parallel()
		quoteRepository_FindAll(t, repo)
	})
}

func quoteRepository_Create_FindByID(t *testing.T, repo service.QuoteRepository) {
	q1 := model.Quote{
		ID:          "quote_id",
		SubmitterID: "user_id",
		Quotee:      "AJBR",
		Quote:       "I'm a quote",
		Context:     "mail trucks",
	}
	if err := repo.Create(context.Background(), q1); err != nil {
		t.Errorf("create quote q1: %v", err)
	}

	gotQ1, err := repo.FindByID(context.Background(), q1.ID)
	if err != nil {
		t.Errorf("find q1: %v", err)
	}
	if !cmp.Equal(gotQ1, q1) {
		t.Errorf("got quote %v, want %v", gotQ1, q1)
	}

	q2 := model.Quote{
		ID:          "quote_id2",
		SubmitterID: "user_id2",
		Quotee:      "Charlene",
		Quote:       "I'm also a quote",
	}

	gotQ2, err := repo.FindByID(context.Background(), q2.ID)
	if err != storage.ErrNotFound {
		t.Errorf("non-existent quote should return ErrNotFound, got %v", err)
	}
	if gotQ2 != (model.Quote{}) {
		t.Errorf("non-existent quote should return empty Quote, got %v", gotQ2)
	}

	err = repo.Create(context.Background(), q2)
	if err != nil {
		t.Errorf("create quote q2: %v", err)
	}

	gotQ2, err = repo.FindByID(context.Background(), q2.ID)
	if err != nil {
		t.Errorf("find q2: %v", err)
	}

	if !cmp.Equal(gotQ2, q2) {
		t.Errorf("got quote %v, want %v", gotQ2, q2)
	}

	err = repo.Create(context.Background(), q2)
	if err != storage.ErrAlreadyExists {
		t.Errorf("create duplicate quote should return ErrAlreadyExists, got %v", err)
	}
}

func quoteRepository_Update(t *testing.T, repo service.QuoteRepository) {
	unchanged := model.Quote{
		ID:          "a934",
		SubmitterID: "user_id",
		Quotee:      "DJ JD",
		Quote:       "This shouldn't change",
		Context:     "",
	}
	if err := repo.Create(context.Background(), unchanged); err != nil {
		t.Errorf("create quote unchanged: %v", err)
	}

	q1 := model.Quote{
		ID:          "x179",
		SubmitterID: "user_id",
		Quotee:      "DJ JD",
		Quote:       "Isn't every truck a hand truck cuz of the steering wheel?",
		Context:     "",
	}
	if err := repo.Create(context.Background(), q1); err != nil {
		t.Errorf("create quote q1: %v", err)
	}

	q1.Quote = "Isn't every truck a hand truck because of the steering wheel?"
	if err := repo.Update(context.Background(), q1); err != nil {
		t.Errorf("update quote q1: %v", err)
	}

	gotQ1, err := repo.FindByID(context.Background(), q1.ID)
	if err != nil {
		t.Errorf("find q1: %v", err)
	}
	if !cmp.Equal(gotQ1, q1) {
		t.Errorf("got quote after first update %v, want %v", gotQ1, q1)
	}

	q1.Quotee = "J-Man"
	q1.Context = "honk"

	if err := repo.Update(context.Background(), q1); err != nil {
		t.Errorf("update quote q1: %v", err)
	}

	gotQ1, err = repo.FindByID(context.Background(), q1.ID)
	if err != nil {
		t.Errorf("find q1: %v", err)
	}
	if !cmp.Equal(gotQ1, q1) {
		t.Errorf("got quote after second update %v, want %v", gotQ1, q1)
	}

	err = repo.Update(context.Background(), model.Quote{ID: "non-existent"})
	if err != storage.ErrNotFound {
		t.Errorf("update non-existent quote should return ErrNotFound, got %v", err)
	}

	gotUnchanged, err := repo.FindByID(context.Background(), unchanged.ID)
	if err != nil {
		t.Errorf("find unchanged: %v", err)
	}
	if !cmp.Equal(gotUnchanged, unchanged) {
		t.Errorf("quote was unexpectedly changed, got %v, want %v", gotUnchanged, unchanged)
	}
}

func quoteRepository_FindAll(t *testing.T, repo service.QuoteRepository) {
	got, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("finding all from empty repo: %v", err)
	}
	want := []model.Quote{}
	if !cmp.Equal(got, want) {
		t.Errorf("finding all from empty repo, got %v, want %v", got, want)
	}

	q1 := model.Quote{
		ID:          "quote_id",
		SubmitterID: "user_id",
		Quotee:      "AJBR",
		Quote:       "I'm a quote",
		Context:     "mail trucks",
	}
	if err := repo.Create(context.Background(), q1); err != nil {
		t.Errorf("create quote q1: %v", err)
	}

	got, err = repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("finding all from repo with one quote: %v", err)
	}
	want = []model.Quote{q1}
	if !cmp.Equal(got, want) {
		t.Errorf("finding all from repo with one quote, got %v, want %v", got, want)
	}

	q2 := model.Quote{
		ID:          "quote_id2",
		SubmitterID: "user_id2",
		Quotee:      "Charlene",
		Quote:       "I'm also a quote",
	}

	if err := repo.Create(context.Background(), q2); err != nil {
		t.Errorf("create quote q2: %v", err)
	}

	got, err = repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("finding all from repo with two quotes: %v", err)
	}
	want = []model.Quote{q1, q2}
	if !cmp.Equal(got, want, cmpopts.SortSlices(func(x, y model.Quote) bool {
		return x.ID < y.ID
	})) {
		t.Errorf("finding all from repo with two quotes, got %v, want %v", got, want)
	}
}
