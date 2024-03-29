package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage"
)

// QuoteRepository is an implementation of the storage.QuoteRepository interface which stores Quotes in a SQLite database.
type QuoteRepository struct {
	db *sql.DB
}

// NewQuoteRepository returns a new QuoteRepository which stores Quotes in the specified SQLite database.
func NewQuoteRepository(db *sql.DB, c *MigrationController) (*QuoteRepository, error) {
	err := c.migrateRepository(db, "quote", []migration{
		{
			version: 1,
			stmts: []string{
				`CREATE TABLE quotes (
					ID text PRIMARY KEY,
					SubmitterID text NOT NULL,
					Quotee text NOT NULL,
					Context text NOT NULL,
					Quote text NOT NULL,
					Created timestamp NOT NULL
				);`,
			},
		},
	})

	return &QuoteRepository{db}, err
}

// Create adds a new Quote to the repository.
func (r *QuoteRepository) Create(ctx context.Context, q model.Quote) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO quotes (ID, SubmitterID, Quotee, Context, Quote, Created) VALUES (?, ?, ?, ?, ?, ?);",
		q.ID, q.SubmitterID, q.Quotee, q.Context, q.Quote, q.Created)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
		return storage.ErrAlreadyExists
	}
	return err
}

// Update updates an existing Quote in the repository.
func (r *QuoteRepository) Update(ctx context.Context, q model.Quote) error {
	result, err := r.db.ExecContext(ctx, "UPDATE quotes SET SubmitterID = ?, Quotee = ?, Context = ?, Quote = ?, Created = ? WHERE ID = ?;",
		q.SubmitterID, q.Quotee, q.Context, q.Quote, q.Created, q.ID)

	if i, _ := result.RowsAffected(); i == 0 {
		return storage.ErrNotFound
	}
	return err
}

// FindByID returns a Quote with the provided ID.
func (r *QuoteRepository) FindByID(ctx context.Context, id string) (model.Quote, error) {
	var q model.Quote
	err := r.db.QueryRowContext(ctx, "SELECT ID, SubmitterID, Quotee, Context, Quote, Created FROM quotes WHERE ID = ?;", id).Scan(
		&q.ID, &q.SubmitterID, &q.Quotee, &q.Context, &q.Quote, &q.Created)

	if err == sql.ErrNoRows {
		return model.Quote{}, storage.ErrNotFound
	}

	return q, err
}

// FindAll returns all Quotes in the repository.
func (r *QuoteRepository) FindAll(ctx context.Context) ([]model.Quote, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT ID, SubmitterID, Quotee, Context, Quote, Created FROM quotes;")
	if err != nil {
		return []model.Quote{}, err
	}
	defer rows.Close()

	quotes := []model.Quote{}
	for rows.Next() {
		var q model.Quote

		err := rows.Scan(&q.ID, &q.SubmitterID, &q.Quotee, &q.Context, &q.Quote, &q.Created)
		if err != nil {
			return quotes, err
		}

		quotes = append(quotes, q)
	}

	return quotes, err
}
