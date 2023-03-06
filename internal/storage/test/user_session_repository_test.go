package storagetest

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage"
	"github.com/willbicks/epigram/internal/storage/inmemory"
	"github.com/willbicks/epigram/internal/storage/sqlite"
)

func TestUserSessionRepositories(t *testing.T) {
	t.Run("UserSessionRepository_inmemory", func(t *testing.T) {
		testUserSessionRepository(t, inmemory.NewUserSessionRepository())
	})

	t.Run("UserSessionRepository_sqlite", func(t *testing.T) {
		mc := &sqlite.MigrationController{}
		db := makeSqliteTestDB(t)
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				t.Fatalf("unable to close database: %v", err)
			}
		}(db)

		repo, err := sqlite.NewUserSessionRepository(db, mc)
		if err != nil {
			t.Fatalf("unable to create user session repo: %v", err)
		}

		testUserSessionRepository(t, repo)
	})

}

// TestUserSessionRepository tests a type implementing the UserSessionRepository interface
func testUserSessionRepository(t *testing.T, repo service.UserSessionRepository) {
	us1 := model.UserSession{
		ID:      "sess_id",
		UserID:  "user_id",
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour),
		IP:      "192.168.0.1",
	}
	if err := repo.Create(context.Background(), us1); err != nil {
		t.Errorf("create user session us2: %v", err)
	}

	gotus1, err := repo.FindByID(context.Background(), us1.ID)
	if err != nil {
		t.Errorf("find us1: %v", err)
	}
	if !cmp.Equal(gotus1, us1) {
		t.Errorf("got user session %v, want %v", gotus1, us1)
	}

	us2 := model.UserSession{
		ID:      "sess_id2",
		UserID:  "user_id2",
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour),
		IP:      "24.197.123.1",
	}
	gotus2, err := repo.FindByID(context.Background(), us2.ID)
	if err != storage.ErrNotFound {
		t.Errorf("non-existent user session should return ErrNotFound, got %v", err)
	}
	if gotus2 != (model.UserSession{}) {
		t.Errorf("non-existent user session should return empty UserSession, got %v", gotus2)
	}

	err = repo.Create(context.Background(), us2)
	if err != nil {
		t.Errorf("create user session us2: %v", err)
	}

	gotus2, err = repo.FindByID(context.Background(), us2.ID)
	if err != nil {
		t.Errorf("find us2: %v", err)
	}
	if !cmp.Equal(gotus2, us2) {
		t.Errorf("got user session %v, want %v", gotus2, us2)
	}

	err = repo.Create(context.Background(), us2)
	if err != storage.ErrAlreadyExists {
		t.Errorf("creating duplicate user session should return ErrAlreadyExists, got %v", err)
	}
}
