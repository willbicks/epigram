package sqlite

import (
	"testing"

	"github.com/willbicks/epigram/internal/storage/validate"
)

func TestUserSessionRepository(t *testing.T) {
	mc := &MigrationController{}
	db := makeSqliteTestDB(t)
	defer db.Close()

	repo, err := NewUserSessionRepository(db, mc)
	if err != nil {
		t.Fatalf("unable to create user session repository: %v", err)
	}

	validate.UserSessionRepository(t, repo)
}
