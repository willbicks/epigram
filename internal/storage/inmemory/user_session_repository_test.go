package inmemory

import (
	"testing"

	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"
)

func TestUserSessionRepository(t *testing.T) {
	validate.UserSessionRepository(t, func() (repo service.UserSessionRepository, closer func()) {
		return NewUserSessionRepository(), func() {}
	})
}
