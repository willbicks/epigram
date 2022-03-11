// package legacy provides a command line utility to migrate data from a legacy quote server using a JSON export.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"

	"github.com/willbicks/epigram/internal/logger"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage/sqlite"
)

func main() {
	// Initialize logger
	log := logger.New(os.Stdout, false)
	log.Level = logger.LevelDebug

	// Declare and parse flags
	inFile := flag.String("in", "in.json", "file containing legacy quotes in JSON format")
	dbFile := flag.String("db", "", "sqlite database to migrate quotes to")
	flag.Parse()

	if *dbFile == "" {
		log.Fatal("No database specified. Please specify one as '--db <file.db>'")
	}

	// Initialize repositories
	connstr := fmt.Sprint("file:./", *dbFile, "?cache=shared&mode=rw")
	log.Debugf("Opening datatabes with connection string: %s", connstr)
	db, err := sql.Open("sqlite3", connstr)
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	defer db.Close()
	repo, u, err := initSqliteRepo(db)
	if err != nil {
		log.Fatalf("unable to initialize repository: %v", err)
	}

	// process input file
	b, err := ioutil.ReadFile(*inFile)
	if err != nil {
		log.Fatalf("opening input file: %v", err)
	}
	quotes := findQuoteArray(log, b)
	if quotes == nil {
		log.Fatal("Unable to find an array of legacy quotes within specified JSON file.")
	}
	for k, v := range quotes {
		log.Infof("found quote %s: %+v", k, v)
	}

	// prompt the user for confirmation to continue
	fmt.Println("Would you like to migrate the logged quotes to the specified repository? (y/N)")
	var resp string
	fmt.Scanln(&resp)
	if len(resp) == 0 || !strings.EqualFold(resp[0:1], "y") {
		log.Fatal("User cancelled legacy migration.")
	}

	// migrate quotes
	var successful, total int
	for k, v := range quotes {
		total++
		if err := migrateQuote(repo, u.ID, v); err != nil {
			log.Warnf("unable to migrate quote %v: %v", k, err)
		} else {
			successful++
		}
	}
	log.Infof("Successfully migrated %v out of %v legacy quotes.", successful, total)
}

// initSqliteRepo accepts an sqlite db conncetion, creates a user and quote repository,
// creates a user to own migrated quotes, and returns the user and quote repo.
func initSqliteRepo(db *sql.DB) (*sqlite.QuoteRepository, *model.User, error) {
	mc := &sqlite.MigrationController{}

	userRepo, err := sqlite.NewUserRepository(db, mc)
	if err != nil {
		return nil, nil, err
	}
	// Creater legacy charisms user
	u := model.User{
		ID:   xid.New().String(),
		Name: "Legacy Charisms User",
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
