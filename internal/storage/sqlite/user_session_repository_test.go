package sqlite

import (
	"testing"

	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"
)

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
