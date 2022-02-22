package sqlite

import (
	"context"
	"database/sql"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage"
)

type UserSessionRepository struct {
	db *sql.DB
}

func NewUserSessionRepository(db *sql.DB, c *MigrationController) (*UserSessionRepository, error) {
	err := c.migrateRepository(db, "usersession", []migration{
		{
			version: 1,
			stmts: []string{
				`CREATE TABLE usersessions (
					ID text PRIMARY KEY,
					UserID text NOT NULL,
					Created timestamp NOT NULL,
					Expires timestamp NOT NULL,
					IP text NOT NULL
				);`,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &UserSessionRepository{db}, nil
}

func (r *UserSessionRepository) Create(ctx context.Context, us model.UserSession) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO usersessions (ID, UserID, Created, Expires, IP) VALUES (?, ?, ?, ?, ?);",
		us.ID, us.UserID, us.Created, us.Expires, us.IP)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserSessionRepository) Update(ctx context.Context, us model.UserSession) error {
	result, err := r.db.ExecContext(ctx, "UPDATE usersessions SET UserID = ?, Created =?, Expires = ?, IP = ? WHERE ID = ?;",
		us.UserID, us.Created, us.Expires, us.IP, us.ID)

	// TODO: Return ErrNotFound if quote does not exist
	if i, _ := result.RowsAffected(); i == 0 {
		return storage.ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func (r *UserSessionRepository) FindByID(ctx context.Context, id string) (model.UserSession, error) {
	var us model.UserSession
	err := r.db.QueryRowContext(ctx, "SELECT ID, UserID, Created, Expires, IP FROM usersessions WHERE ID = ?;", id).Scan(
		&us.ID, &us.UserID, &us.Created, &us.Expires, &us.IP)

	if err == sql.ErrNoRows {
		return model.UserSession{}, storage.ErrNotFound
	} else if err != nil {
		return model.UserSession{}, err
	}

	return us, nil
}

func (r *UserSessionRepository) FindAll(ctx context.Context) ([]model.UserSession, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT ID, UserID, Created, Expires, IP FROM usersessions;")
	if err != nil {
		return []model.UserSession{}, err
	}
	defer rows.Close()

	sessions := []model.UserSession{}
	for rows.Next() {
		var us model.UserSession

		err := rows.Scan(&us.ID, &us.UserID, &us.Created, &us.Expires, &us.IP)
		if err != nil {
			return sessions, err
		}

		sessions = append(sessions, us)
	}

	if err := rows.Err(); err != nil {
		return sessions, err
	}

	return sessions, nil
}
