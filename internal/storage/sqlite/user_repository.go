package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage"
)

// UserRepository is an implementation of the service.UserRepository interface which stores Users in a SQLite database.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository returns a new UserRepository which stores Users in the provided SQLite database.
func NewUserRepository(db *sql.DB, c *MigrationController) (*UserRepository, error) {
	err := c.migrateRepository(db, "user", []migration{
		{
			version: 1,
			stmts: []string{
				`CREATE TABLE users (
					ID text PRIMARY KEY,
					Name text NOT NULL,
					Email text NOT NULL,
					PictureURL text NOT NULL,
					Created timestamp NOT NULL,
					QuizAttempts smallint NOT NULL,
					QuizPassed boolean NOT NULL,
					Banned boolean NOT NULL,
					Admin boolean NOT NULL
				);`,
			},
		},
	})

	return &UserRepository{db}, err
}

// Create adds a new User to the repository.
func (r *UserRepository) Create(ctx context.Context, u model.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (ID, Name, Email, PictureURL, Created, QuizAttempts, QuizPassed, Banned, Admin) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
		u.ID, u.Name, u.Email, u.PictureURL, u.Created, u.QuizAttempts, u.QuizPassed, u.Banned, u.Admin)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
		return storage.ErrAlreadyExists
	}
	return err
}

// Update updates an existing User in the repository.
func (r *UserRepository) Update(ctx context.Context, u model.User) error {
	result, err := r.db.ExecContext(ctx, "UPDATE users SET Name = ?, Email = ?, PictureURL = ?, Created = ?, QuizAttempts = ?, QuizPassed = ?, Banned = ?, Admin = ? WHERE ID = ?;",
		u.Name, u.Email, u.PictureURL, u.Created, u.QuizAttempts, u.QuizPassed, u.Banned, u.Admin, u.ID)

	if i, _ := result.RowsAffected(); i == 0 {
		return storage.ErrNotFound
	}
	return err
}

// FindByID returns the User with the provided ID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	var u model.User
	err := r.db.QueryRowContext(ctx, "SELECT ID, Name, Email, PictureURL, Created, QuizAttempts, QuizPassed, Banned, Admin FROM users WHERE ID = ?;", id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PictureURL, &u.Created, &u.QuizAttempts, &u.QuizPassed, &u.Banned, &u.Admin)

	if err == sql.ErrNoRows {
		return model.User{}, storage.ErrNotFound
	}

	return u, err
}

// FindAll returns all Users in the repository.
func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT ID, Name, Email, PictureURL, Created, QuizAttempts, QuizPassed, Banned, Admin FROM users;")
	if err != nil {
		return []model.User{}, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User

		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PictureURL, &u.Created, &u.QuizAttempts, &u.QuizPassed, &u.Banned, &u.Admin)
		if err != nil {
			return users, err
		}

		users = append(users, u)
	}

	return users, rows.Err()
}
