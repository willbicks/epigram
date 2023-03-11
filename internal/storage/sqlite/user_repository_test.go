package sqlite

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"
)

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
