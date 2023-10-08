package validate

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage"
)

// UserSessionRepository tests a type implementing the UserSessionRepository interface
func UserSessionRepository(t *testing.T, repoFactory func() (repo service.UserSessionRepository, closer func())) {

	repo, close := repoFactory()
	defer close()

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

	gotUS1, err := repo.FindByID(context.Background(), us1.ID)
	if err != nil {
		t.Errorf("find us1: %v", err)
	}
	if !cmp.Equal(gotUS1, us1) {
		t.Errorf("got user session %v, want %v", gotUS1, us1)
	}

	us2 := model.UserSession{
		ID:      "sess_id2",
		UserID:  "user_id2",
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour),
		IP:      "24.197.123.1",
	}
	gotUS2, err := repo.FindByID(context.Background(), us2.ID)
	if err != storage.ErrNotFound {
		t.Errorf("non-existent user session should return ErrNotFound, got %v", err)
	}
	if gotUS2 != (model.UserSession{}) {
		t.Errorf("non-existent user session should return empty UserSession, got %v", gotUS2)
	}

	err = repo.Create(context.Background(), us2)
	if err != nil {
		t.Errorf("create user session us2: %v", err)
	}

	gotUS2, err = repo.FindByID(context.Background(), us2.ID)
	if err != nil {
		t.Errorf("find us2: %v", err)
	}
	if !cmp.Equal(gotUS2, us2) {
		t.Errorf("got user session %v, want %v", gotUS2, us2)
	}

	err = repo.Create(context.Background(), us2)
	if err != storage.ErrAlreadyExists {
		t.Errorf("creating duplicate user session should return ErrAlreadyExists, got %v", err)
	}
}
