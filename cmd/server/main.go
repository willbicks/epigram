package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/willbicks/epigram/internal/config"

	"github.com/willbicks/epigram/internal/logger"
	quoteserver "github.com/willbicks/epigram/internal/server/http"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/inmemory"
	"github.com/willbicks/epigram/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize logger
	log := logger.New(os.Stdout, os.Stderr, true)
	log.Level = logger.LevelDebug

	// Configuration parsing
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal(err.Error())
	}

	var userRepo service.UserRepository
	var userSessionRepo service.UserSessionRepository
	var quoteRepo service.QuoteRepository

	switch cfg.Repo {
	case config.Inmemory:
		userRepo = inmemory.NewUserRepository()
		userSessionRepo = inmemory.NewUserSessionRepository()
		quoteRepo = inmemory.NewQuoteRepository()
	case config.SQLite:
		mc := &sqlite.MigrationController{}
		db, err := sql.Open("sqlite3", fmt.Sprint("file:", cfg.DBLoc, "?cache=shared&mode=rwc"))
		if err != nil {
			log.Fatalf("unable to open database: %v", err)
		}

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatalf("unable to close database: %v", err)
			}
		}(db)

		userRepo, err = sqlite.NewUserRepository(db, mc)
		if err != nil {
			log.Fatalf("unable to create user repo: %v", err)
		}

		quoteRepo, err = sqlite.NewQuoteRepository(db, mc)
		if err != nil {
			log.Fatalf("unable to create quote repo: %v", err)
		}

		userSessionRepo, err = sqlite.NewUserSessionRepository(db, mc)
		if err != nil {
			log.Fatalf("unable to create user sess repo: %v", err)
		}
	}

	// Quote Server Initialization
	cs := quoteserver.QuoteServer{
		QuoteService: service.NewQuoteService(quoteRepo),
		UserService:  service.NewUserService(userRepo, userSessionRepo),
		QuizService:  service.NewEntryQuizService(cfg.EntryQuestions),
		Logger:       log,
		Config:       cfg,
	}

	if err := cs.Init(); err != nil {
		log.Fatalf("critical error while initializing server: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	log.Infof("Running server at %s ...", addr)
	s := http.Server{
		Addr:              addr,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      4 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           cs,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to listen and serve: %v", err)
	}
}
