package sqlite

import (
	"database/sql"
	"fmt"
	"sort"
	"sync"
)

// MigrationController is responsible for running migrations on the SQLite database.
type MigrationController struct {
	// createTable ensures that an attempt to create a migrations table
	// is only made once per application start.
	createTable sync.Once
}

type migration struct {
	version int
	stmts   []string
}

func (c *MigrationController) migrateRepository(db *sql.DB, repoName string, migrations []migration) error {
	// one time only, create a migration table if it doesn't exist
	c.createTable.Do(func() {
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY AUTOINCREMENT, repo TEXT,  version INT);")
		if err != nil {
			panic(fmt.Errorf("creating migrations table: %v", err))
		}
	})

	// determine the maximum version that this repository has been migrated to
	var maxVer sql.NullInt32
	err := db.QueryRow("SELECT MAX(version) FROM migrations WHERE repo = ?;", repoName).Scan(&maxVer)
	if err != nil {
		return fmt.Errorf("selecting max version: %w", err)
	}

	// sort the provided list of migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	for _, m := range migrations {
		if m.version > int(maxVer.Int32) {
			// start a transaction for this migration
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("starting migration tx: %w", err)
			}

			// execute all the statements in this migration
			for i, stmt := range m.stmts {
				_, err := tx.Exec(stmt)
				if err != nil {
					if err := tx.Rollback(); err != nil {
						return fmt.Errorf("err executing statement %v of migration %v on %v AND unable to rollback transaction: %w", i, m.version, repoName, err)
					}
					return fmt.Errorf("err executing statement %v of migration %v on %v: %w", i, m.version, repoName, err)
				}
			}

			// update the migrations table
			_, err = tx.Exec("INSERT INTO migrations (repo, version) VALUES (?, ?);", repoName, m.version)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return fmt.Errorf("unable to record migration %v of %v AND unable to rollback transaction: %w", m.version, repoName, err)
				}
				return fmt.Errorf("unable to record migration %v of %v: %w", m.version, repoName, err)
			}

			// commit this migration
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("err commiting migration %v on %v: %w", m.version, repoName, err)
			}
		}
	}

	return nil
}
