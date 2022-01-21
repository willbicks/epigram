package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/willbicks/charisms/model"
)

type UserRepository interface {
	Create(ctx context.Context, q model.User) error
	Update(ctx context.Context, q model.User) error
	FindByID(ctx context.Context, id string) (model.User, error)
	FindAll(ctx context.Context) ([]model.User, error)
}

type User struct {
	repo UserRepository
}

func NewUserService(r UserRepository) User {
	return User{
		repo: r,
	}
}

// FromIDToken returns a user from the specified OIDC token (assumed to be allready verified).
// If a user allready exists with the specified ID (derrived from the issuer URL and sub claim),
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
	if u, err := s.repo.FindByID(ctx, id); err == nil {
		return u, nil
	}

	// User does not exist, create them
	u := model.User{
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

func (s User) CreateUser(ctx context.Context, u *model.User) error {
	err := ServiceError{
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

	return s.repo.Create(ctx, *u)
}

func (s User) FindUserById(ctx context.Context, id string) (model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s User) UpdateUser(ctx context.Context, u model.User) error {
	return s.repo.Update(ctx, u)
}
