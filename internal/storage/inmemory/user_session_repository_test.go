package inmemory

import (
	"testing"

	"github.com/willbicks/epigram/internal/storage/validate"
)

func TestUserSessionRepository(t *testing.T) {
	validate.UserSessionRepository(t, NewUserSessionRepository())
}
