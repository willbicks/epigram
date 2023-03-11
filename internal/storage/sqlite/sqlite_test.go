package sqlite

import (
	"database/sql"
	"testing"

	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"

	_ "github.com/mattn/go-sqlite3"
)

func makeSqliteTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("unable to open database: %v", err)
	}

	return db
}

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

func TestUserRepository(t *testing.T) {
	validate.UserRepository(t, func() (repo service.UserRepository, closer func()) {
		mc := &MigrationController{}
		db := makeSqliteTestDB(t)

		repo, err := NewUserRepository(db, mc)
		if err != nil {
			t.Fatalf("unable to create user repository: %v", err)
		}

		return repo, func() {
			err = db.Close()
			if err != nil {
				t.Fatalf("unable to close database: %v", err)
			}
		}
	})
}

func TestUserSessionRepository(t *testing.T) {
	validate.UserSessionRepository(t, func() (repo service.UserSessionRepository, closer func()) {
		mc := &MigrationController{}
		db := makeSqliteTestDB(t)

		repo, err := NewUserSessionRepository(db, mc)
		if err != nil {
			t.Fatalf("unable to create user session repository: %v", err)
		}

		return repo, func() {
			err = db.Close()
			if err != nil {
				t.Fatalf("unable to close database: %v", err)
			}
		}
	})
}
