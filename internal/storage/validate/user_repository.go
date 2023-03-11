package validate

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage"
)

var (
	u1 = model.User{
		ID:           "user_id",
		Name:         "Ficky Neldo",
		Email:        "fn@example.com",
		PictureURL:   "https://example.com/fn.jpg",
		Created:      time.Now(),
		QuizPassed:   false,
		QuizAttempts: 1,
		Banned:       true,
		Admin:        false,
	}

	u2 = model.User{
		ID:           "user_id2",
		Name:         "Bill Wicks",
		Email:        "bw@example.com",
		PictureURL:   "https://example.com/bw.jpg",
		Created:      time.Now(),
		QuizPassed:   true,
		QuizAttempts: 2,
		Banned:       false,
		Admin:        true,
	}

	u3 = model.User{
		ID:           "user_id3",
		Name:         "Rawb",
		Email:        "rlhllc@rtcexample.edu",
		PictureURL:   "https://www.rtcexample.edu/wp-content/uploads/2018/08/rlhllc-300x300.jpg",
		Created:      time.Now(),
		QuizPassed:   true,
		QuizAttempts: 1,
		Banned:       false,
		Admin:        false,
	}
)

// UserRepository validates a type implementing the UserRepository interface
func UserRepository(t *testing.T, repoFactory func() (repo service.UserRepository, close func())) {
	t.Run("Create_FindByID", func(t *testing.T) {
		repo, close := repoFactory()
		defer close()
		t.Parallel()
		userRepository_Create_FindByID(t, repo)
	})
	t.Run("Update", func(t *testing.T) {
		repo, close := repoFactory()
		defer close()
		t.Parallel()
		userRepository_Update(t, repo)
	})
	t.Run("FindAll", func(t *testing.T) {
		repo, close := repoFactory()
		defer close()
		t.Parallel()
		userRepository_FindAll(t, repo)
	})
}

func userRepository_Create_FindByID(t *testing.T, repo service.UserRepository) {
	if err := repo.Create(context.Background(), u1); err != nil {
		t.Errorf("create user u1: %v", err)
	}
	gotu1, err := repo.FindByID(context.Background(), u1.ID)
	if err != nil {
		t.Errorf("find u1: %v", err)
	}
	if !cmp.Equal(gotu1, u1) {
		t.Errorf("got user %v, want %v", gotu1, u1)
	}

	gotu2, err := repo.FindByID(context.Background(), u2.ID)
	if err != storage.ErrNotFound {
		t.Errorf("find u2 before created: got error %v, want %v", err, storage.ErrNotFound)
	}
	if gotu2 != (model.User{}) {
		t.Errorf("find u2 before created: got user %v, want %v", gotu2, model.User{})
	}

	if err := repo.Create(context.Background(), u2); err != nil {
		t.Errorf("create user u2: %v", err)
	}
	gotu2, err = repo.FindByID(context.Background(), u2.ID)
	if err != nil {
		t.Errorf("find u2: %v", err)
	}
	if !cmp.Equal(gotu2, u2) {
		t.Errorf("got user %v, want %v", gotu2, u2)
	}

	if err := repo.Create(context.Background(), u1); err != storage.ErrAlreadyExists {
		t.Errorf("create user u1 again: got error %v, want %v", err, storage.ErrAlreadyExists)
	}
}

func userRepository_Update(t *testing.T, repo service.UserRepository) {
	if err := repo.Create(context.Background(), u3); err != nil {
		t.Errorf("create user u3: %v", err)
	}

	uEdit := u1

	if err := repo.Create(context.Background(), uEdit); err != nil {
		t.Errorf("create user uEdit: %v", err)
	}
	gotuEdit, err := repo.FindByID(context.Background(), uEdit.ID)
	if err != nil {
		t.Errorf("find uEdit: %v", err)
	}
	if !cmp.Equal(gotuEdit, uEdit) {
		t.Errorf("got user %v, want %v", gotuEdit, uEdit)
	}

	uEdit.Name = "Ficky Neldo II"
	uEdit.Email = "fn2@example.com"
	uEdit.PictureURL = "https://example.com/fn2.jpg"
	uEdit.QuizAttempts = 2
	uEdit.QuizPassed = true
	uEdit.Banned = false
	uEdit.Admin = true

	if err := repo.Update(context.Background(), uEdit); err != nil {
		t.Errorf("update user uEdit: %v", err)
	}
	gotuEdit, err = repo.FindByID(context.Background(), uEdit.ID)
	if err != nil {
		t.Errorf("find uEdit: %v", err)
	}
	if !cmp.Equal(gotuEdit, uEdit) {
		t.Errorf("got user %v, want %v", gotuEdit, uEdit)
	}

	if err := repo.Update(context.Background(), u2); err != storage.ErrNotFound {
		t.Errorf("update user u2: got error %v, want %v", err, storage.ErrNotFound)
	}

	gotu3, err := repo.FindByID(context.Background(), u3.ID)
	if err != nil {
		t.Errorf("find u3: %v", err)
	}
	if !cmp.Equal(gotu3, u3) {
		t.Errorf("u3 should be unchanged, got user %v, want %v", gotu3, u3)
	}
}

func userRepository_FindAll(t *testing.T, repo service.UserRepository) {
	gotUsers, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("finding no users: %v", err)
	}
	wantUsers := []model.User{}
	if !cmp.Equal(gotUsers, wantUsers) {
		t.Errorf("finding no users, got %v, want %v", gotUsers, wantUsers)
	}

	if err := repo.Create(context.Background(), u1); err != nil {
		t.Errorf("create user u1: %v", err)
	}

	gotUsers, err = repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("finding no users: %v", err)
	}
	wantUsers = []model.User{u1}
	if !cmp.Equal(gotUsers, wantUsers) {
		t.Errorf("finding no users, got %v, want %v", gotUsers, wantUsers)
	}

	if err := repo.Create(context.Background(), u2); err != nil {
		t.Errorf("create user u2: %v", err)
	}
	if err := repo.Create(context.Background(), u3); err != nil {
		t.Errorf("create user u3: %v", err)
	}

	gotUsers, err = repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("find all users: %v", err)
	}
	wantUsers = []model.User{u1, u2, u3}
	if !cmp.Equal(gotUsers, wantUsers, cmpopts.SortSlices(func(u1, u2 model.User) bool {
		return u1.ID < u2.ID
	})) {
		t.Errorf("got users %v, want %v", gotUsers, wantUsers)
	}
}
