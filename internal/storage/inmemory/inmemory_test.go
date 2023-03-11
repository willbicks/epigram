package inmemory

import (
	"testing"

	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/validate"
)

func TestQuoteRepository(t *testing.T) {
	validate.QuoteRepository(t, func() (repo service.QuoteRepository, closer func()) {
		return NewQuoteRepository(), func() {}
	})
}

func TestUserRepository(t *testing.T) {
	validate.UserRepository(t, func() (repo service.UserRepository, closer func()) {
		return NewUserRepository(), func() {}
	})
}

func TestUserSessionRepository(t *testing.T) {
	validate.UserSessionRepository(t, func() (repo service.UserSessionRepository, closer func()) {
		return NewUserSessionRepository(), func() {}
	})
}
