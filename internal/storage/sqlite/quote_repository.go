package sqlite

import (
	"context"
	"database/sql"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage"
)

type QuoteRepository struct {
	db *sql.DB
}

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

	if err != nil {
		return nil, err
	}

	return &QuoteRepository{db}, nil
}

func (r *QuoteRepository) Create(ctx context.Context, q model.Quote) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO quotes (ID, SubmitterID, Quotee, Context, Quote, Created) VALUES (?, ?, ?, ?, ?, ?);",
		q.ID, q.SubmitterID, q.Quotee, q.Context, q.Quote, q.Created)

	if err != nil {
		return err
	}
	return nil
}

func (r *QuoteRepository) Update(ctx context.Context, q model.Quote) error {
	result, err := r.db.ExecContext(ctx, "UPDATE quotes SET SubmitterID = ?, Quotee = ?, Context = ?, Quote = ?, Created = ? WHERE ID = ?;",
		q.SubmitterID, q.Quotee, q.Context, q.Quote, q.Created, q.ID)

	// TODO: Return ErrNotFound if quote does not exist
	if i, _ := result.RowsAffected(); i == 0 {
		return storage.ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func (r *QuoteRepository) FindByID(ctx context.Context, id string) (model.Quote, error) {
	var q model.Quote
	err := r.db.QueryRowContext(ctx, "SELECT ID, SubmitterID, Quotee, Context, Quote, Created FROM quotes WHERE ID = ?;", id).Scan(
		&q.ID, &q.SubmitterID, &q.Quotee, &q.Context, &q.Quote, &q.Created)

	if err == sql.ErrNoRows {
		return model.Quote{}, storage.ErrNotFound
	} else if err != nil {
		return model.Quote{}, err
	}

	return q, nil
}

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

	if err := rows.Err(); err != nil {
		return quotes, err
	}

	return quotes, nil
}
