package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage"

	"github.com/coreos/go-oidc/v3/oidc"
)

// UserRepository provides methods for storing, manipulating, and retrieving Users.
type UserRepository interface {
	Create(ctx context.Context, u model.User) error
	Update(ctx context.Context, u model.User) error
	FindByID(ctx context.Context, id string) (model.User, error)
	FindAll(ctx context.Context) ([]model.User, error)
}

// User provides a service for interacting with Users.
type User struct {
	ur   UserRepository
	sess UserSession
}

// NewUserService returns a new UserService with the provided UserRepository.
func NewUserService(ur UserRepository, sr UserSessionRepository) User {
	return User{
		ur:   ur,
		sess: NewUserSessionService(sr),
	}
}

// GetUserFromIDToken returns a user from the specified OIDC token (assumed to be already verified).
// If a user already exists with the specified ID (derived from the issuer URL and subclass),
// that user is returned. If no such user exists, a new user is created based on the token
// details and returned.
func (s User) GetUserFromIDToken(ctx context.Context, token oidc.IDToken) (model.User, error) {
	var claims struct {
		Issuer     string `json:"iss"`
		Subject    string `json:"sub"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		PictureURL string `json:"picture"`
	}
	if err := token.Claims(&claims); err != nil {
		return model.User{}, fmt.Errorf("unmarshalling token claims: %w", err)
	}

	domain := claims.Issuer
	if strings.Contains(domain, "://") {
		domain = strings.Split(domain, "://")[1]
	}
	if strings.Contains(domain, "/") {
		domain = strings.Split(domain, "/")[0]
	}
	id := domain + "/" + claims.Subject

	// Check if the user exists, and if so, return them
	u, err := s.ur.FindByID(ctx, id)
	if err == nil {
		return u, nil
	} else if err != storage.ErrNotFound {
		return model.User{}, fmt.Errorf("unable to find from user repo: %w", err)
	}

	// User does not exist, create them
	u = model.User{
		ID:         id,
		Name:       claims.Name,
		Email:      claims.Email,
		PictureURL: claims.PictureURL,
	}

	if err := s.CreateUser(ctx, &u); err != nil {
		return model.User{}, fmt.Errorf("creating user from id token: %w", err)
	}

	return u, nil
}

// CreateUser stores the provided User in the database, updating their Created time.
func (s User) CreateUser(ctx context.Context, u *model.User) error {
	err := Error{
		StatusCode: 400,
	}

	if u.ID == "" {
		err.addIssue("User ID required.")
	}
	if u.Email == "" {
		err.addIssue("User Email required.")
	}
	if u.Name == "" {
		err.addIssue("User Name required.")
	}

	if err.HasIssues() {
		return err
	}

	u.Created = time.Now()

	return s.ur.Create(ctx, *u)
}

// FindUserByID returns the user with the specified ID.
func (s User) FindUserByID(ctx context.Context, id string) (model.User, error) {
	return s.ur.FindByID(ctx, id)
}

// UpdateUser updates the user in the database.
func (s User) UpdateUser(ctx context.Context, u model.User) error {
	return s.ur.Update(ctx, u)
}

// CreateUserSession creates a new UserSession for the specified user, and returns it.
func (s User) CreateUserSession(ctx context.Context, u model.User, IP string) (model.UserSession, error) {
	return s.sess.CreateUserSession(ctx, u, IP)
}

// GetUserFromSessionID returns the user associated with the specified session ID.
func (s User) GetUserFromSessionID(ctx context.Context, sessID string) (model.User, error) {
	sess, err := s.sess.FindSessionByID(ctx, sessID)
	if err != nil {
		return model.User{}, err
	}
	return s.ur.FindByID(ctx, sess.UserID)
}

// RecordQuizAttempt records that the user attempted to complete a quiz, updates their information, and returns
// either an empty string (pass), or the reason they failed (either got a question wrong or too many attempts),
func (s *User) RecordQuizAttempt(ctx context.Context, u *model.User, passed bool) (failReason string, err error) {
	u.QuizAttempts++
	u.QuizPassed = passed

	if err := s.UpdateUser(ctx, *u); err != nil {
		return "Unable to update user", err
	}

	if u.QuizAttempts > model.MaxQuizAttempts {
		return "Too many failed quiz attempts, please contact an administrator.", nil
	}

	if !passed {
		return "Sorry, at least one answer was incorrect.", nil
	}

	return "", nil
}

// GetAllUsers returns a slice of all users, and can ony be accessed by admins.
func (s *User) GetAllUsers(ctx context.Context) ([]model.User, error) {
	if err := verifyAdminPrivilege(ctx); err != nil {
		return nil, err
	}

	users, err := s.ur.FindAll(ctx)

	return users, err
}
