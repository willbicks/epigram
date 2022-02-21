package sqlite

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/matryer/is"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage"
)

func createRepo(t *testing.T) *UserRepository {
	db, err := sql.Open("sqlite3", "file:./test.db?cache=shared&mode=rwc")
	if err != nil {
		t.Errorf("unable to open database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		os.Remove("./test.db")
	})

	ur, err := NewUserRepository(db, &MigrationController{})
	if err != nil {
		t.Errorf("unable to create repository: %v", err)
	}

	return ur
}

func usersEqual(u1 model.User, u2 model.User) bool {
	return u1.ID == u2.ID &&
		u1.Name == u2.Name &&
		u1.PictureURL == u2.PictureURL &&
		u1.Email == u2.Email &&
		u1.Created.Equal(u2.Created) &&
		u1.QuizPassed == u2.QuizPassed &&
		u1.QuizAttempts == u2.QuizAttempts &&
		u1.Banned == u2.Banned &&
		u1.Admin == u2.Admin
}

func TestUserRepository_Create(t *testing.T) {
	is := is.New(t)

	repo := createRepo(t)

	u1 := model.User{
		ID:      xid.New().String(),
		Name:    "George",
		Email:   "george@example.com",
		Created: time.Now(),
	}
	u2 := model.User{
		ID:         xid.New().String(),
		Name:       "Leo",
		Email:      "leo33@example.com",
		PictureURL: "https://google.com/test",
		Created:    time.Now(),
	}

	is.NoErr(repo.Create(context.Background(), u1)) // should create user 1
	is.NoErr(repo.Create(context.Background(), u2)) // should create user 2

	err := repo.Create(context.Background(), u1)
	is.True(err != nil) // should fail attempt to re-create user 1
}

func TestUserRepository_FindByID(t *testing.T) {
	is := is.New(t)

	repo := createRepo(t)

	u1 := model.User{
		ID:      xid.New().String(),
		Name:    "George",
		Email:   "george@example.com",
		Created: time.Now(),
	}
	u2 := model.User{
		ID:         xid.New().String(),
		Name:       "Leo",
		Email:      "leo33@example.com",
		PictureURL: "https://google.com/test",
		Created:    time.Now(),
	}

	is.NoErr(repo.Create(context.Background(), u1)) // should create user 1
	is.NoErr(repo.Create(context.Background(), u2)) // should create user 2

	uFound, err := repo.FindByID(context.Background(), u1.ID)
	is.NoErr(err) // should find u1
	//t.Logf("Got:    %+v\nWanted: %+v\n", uFound, u1)
	is.True(usersEqual(u1, uFound)) // u1 should equal uFound

	uFound, err = repo.FindByID(context.Background(), u2.ID)
	is.NoErr(err) // should find u2
	//t.Logf("Got:    %+v\nWanted: %+v\n", uFound, u2)
	is.True(usersEqual(u2, uFound)) // u2 should equal uFound

	uFound, err = repo.FindByID(context.Background(), xid.New().String())
	is.Equal(err, storage.ErrNotFound) // random should not be found
}

func TestUserRepository_Update(t *testing.T) {
	is := is.New(t)
	repo := createRepo(t)

	u1 := model.User{
		ID:      xid.New().String(),
		Name:    "George",
		Email:   "george@example.com",
		Created: time.Now(),
	}
	u2 := model.User{
		ID:         xid.New().String(),
		Name:       "Leo",
		Email:      "leo33@example.com",
		PictureURL: "https://google.com/test",
		Created:    time.Now(),
	}
	u3 := model.User{
		ID:   xid.New().String(),
		Name: "Un Created",
	}

	is.NoErr(repo.Create(context.Background(), u1)) // should create user 1
	is.NoErr(repo.Create(context.Background(), u2)) // should create user 2

	u1.QuizAttempts = 1
	u1.QuizPassed = true
	u1.Email = "redlight@example.com"
	u1.Name = "George Alcam"

	is.NoErr(repo.Update(context.Background(), u1)) // should update u1
	uFound, err := repo.FindByID(context.Background(), u1.ID)
	is.NoErr(err)                   // finding updated u1 should suceed
	is.True(usersEqual(u1, uFound)) // uFound should match updated u1

	u2.Name = "Jimothy Jacobs"
	u2.Admin = true
	u2.PictureURL = "https://u3content.xyz/3323418980.jpg"
	u2.Banned = true

	is.NoErr(repo.Update(context.Background(), u2)) // should update u2
	uFound, err = repo.FindByID(context.Background(), u2.ID)
	is.NoErr(err)                   // finding updated u2 should suceed
	is.True(usersEqual(u2, uFound)) // uFound should match updated u2

	err = repo.Update(context.Background(), u3)
	is.True(err != nil) // should return error if updated user does not exist
}

func TestUserRepository_FindAll(t *testing.T) {
	is := is.New(t)
	repo := createRepo(t)

	u1 := model.User{
		ID:      xid.New().String(),
		Name:    "George",
		Email:   "george@example.com",
		Created: time.Now(),
	}
	u2 := model.User{
		ID:         xid.New().String(),
		Name:       "Leo",
		Email:      "leo33@example.com",
		PictureURL: "https://google.com/test",
		Created:    time.Now(),
	}
	u3 := model.User{
		ID:   xid.New().String(),
		Name: "Un Created",
	}

	_, err := repo.FindAll(context.Background())
	is.NoErr(err) // FindAll should not fail with empty table

	is.NoErr(repo.Create(context.Background(), u1)) // should create user 1
	is.NoErr(repo.Create(context.Background(), u2)) // should create user 2

	res, err := repo.FindAll(context.Background())
	is.NoErr(err)         // FindAll should not fail with two entries
	is.Equal(len(res), 2) // should be two entries in result

	userSliceContains := func(slice []model.User, want model.User) bool {
		for _, u := range slice {
			if usersEqual(u, want) {
				return true
			}
		}
		return false
	}

	is.True(userSliceContains(res, u1))  // returned slice should contain u1
	is.True(userSliceContains(res, u2))  // returned slice should contain u2
	is.True(!userSliceContains(res, u3)) // returned slice should NOT contain u3

	is.NoErr(repo.Create(context.Background(), u3)) // should create user 3

	res, err = repo.FindAll(context.Background())
	is.NoErr(err)         // FindAll should not fail with three entries
	is.Equal(len(res), 3) // should be three entries in result

	is.True(userSliceContains(res, u1)) // returned slice should contain u1
	is.True(userSliceContains(res, u2)) // returned slice should contain u2
	is.True(userSliceContains(res, u3)) // returned slice should contain u3
}
