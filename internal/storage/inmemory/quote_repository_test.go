package inmemory

import (
	"testing"

	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"
)

func TestQuoteRepository(t *testing.T) {
	validate.QuoteRepository(t, func() (repo service.QuoteRepository, closer func()) {
		return NewQuoteRepository(), func() {}
	})
}
