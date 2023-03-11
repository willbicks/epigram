package inmemory

import (
	"testing"

	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"
)

func TestUserRepository(t *testing.T) {
	validate.UserRepository(t, func() (repo service.UserRepository, closer func()) {
		return NewUserRepository(), func() {}
	})
}
