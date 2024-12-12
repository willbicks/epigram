// package legacy provides a command line utility to migrate data from a legacy quote server using a JSON export.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"

	"github.com/willbicks/epigram/internal/logutils"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage/sqlite"
)

func main() {
	// Initialize logger
	log := slog.Default()

	// Declare and parse flags
	inFile := flag.String("in", "in.json", "file containing legacy quotes in JSON format")
	dbFile := flag.String("db", "", "sqlite database to migrate quotes to")
	flag.Parse()

	if *dbFile == "" {
		log.Error("No database specified. Please specify one as '--db <file.db>'")
		os.Exit(1)
	}

	// Initialize repositories
	connStr := fmt.Sprint("file:./", *dbFile, "?cache=shared&mode=rw")
	log.Debug("opening database", "conn_str", connStr)
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		log.Error("unable to open database", logutils.Error(err))
		os.Exit(1)
	}
	defer db.Close()
	repo, u, err := initSqliteRepo(db)
	if err != nil {
		log.Error("unable to initialize repository", logutils.Error(err))
		os.Exit(1)
	}

	// process input file
	b, err := os.ReadFile(*inFile)
	if err != nil {
		log.Error("opening input file", logutils.Error(err))
		os.Exit(1)
	}
	quotes := findQuoteArray(b)
	if quotes == nil {
		log.Error("unable to find an array of legacy quotes in specified JSON file")
		os.Exit(1)
	}
	for k, v := range quotes {
		fmt.Printf("found quote %s: %+v", k, v)
	}

	// prompt the user for confirmation to continue
	fmt.Println("Would you like to migrate the above quotes to the specified repository? (y/N)")
	var resp string
	fmt.Scanln(&resp)
	if len(resp) == 0 || !strings.EqualFold(resp[0:1], "y") {
		log.Error("user cancelled legacy migration")
		os.Exit(1)
	}

	// migrate quotes
	var successful, total int
	for k, v := range quotes {
		total++
		if err := migrateQuote(repo, u.ID, v); err != nil {
			log.Warn("unable to migrate quote", "key", k, logutils.Error(err))
		} else {
			successful++
		}
	}
	fmt.Printf("Successfully migrated %v out of %v legacy quotes.", successful, total)
}

// initSqliteRepo accepts an sqlite db connection, creates a user and quote repository,
// creates a user to own migrated quotes, and returns the user and quote repo.
func initSqliteRepo(db *sql.DB) (*sqlite.QuoteRepository, *model.User, error) {
	mc := &sqlite.MigrationController{}

	userRepo, err := sqlite.NewUserRepository(db, mc)
	if err != nil {
		return nil, nil, err
	}
	// Create user to own imported quotes
	u := model.User{
		ID:   xid.New().String(),
		Name: "Legacy User",
	}
	if err := userRepo.Create(context.Background(), u); err != nil {
		return nil, nil, err
	}

	quoteRepo, err := sqlite.NewQuoteRepository(db, mc)
	if err != nil {
		return nil, nil, err
	}

	return quoteRepo, &u, nil
}
