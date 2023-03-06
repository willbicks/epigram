package sqlite

import (
	"testing"

	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"
)

func TestQuoteRepository(t *testing.T) {
	validate.QuoteRepository(t, func() (repo service.QuoteRepository, close func()) {
		mc := &MigrationController{}
		db := makeSqliteTestDB(t)

		repo, err := NewQuoteRepository(db, mc)
		if err != nil {
			t.Fatalf("unable to create quote repository: %v", err)
		}

		return repo, func() {
			err = db.Close()
			if err != nil {
				t.Fatalf("unable to close database: %v", err)
			}
		}
	})
}
