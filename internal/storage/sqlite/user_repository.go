package sqlite

import (
	"context"
	"database/sql"

	"github.com/willbicks/epigram/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB, c *MigrationController) (*UserRepository, error) {
	err := c.migrateRepository(db, "user", []migration{
		{
			version: 1,
			stmts: []string{
				`CREATE TABLE users (
					ID varchar(20) PRIMARY KEY,
					Name text,
					Email text,
					PictureURL text,
					Created timestamp,
					QuizAttempts smallint,
					Banned boolean,
					Admin boolean
				);`,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &UserRepository{db}, nil
}

func (r *UserRepository) Create(ctx context.Context, q model.User) error {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) Update(ctx context.Context, q model.User) error {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	panic("not implemented") // TODO: Implement
}
