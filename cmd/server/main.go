package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/willbicks/epigram/internal/config"
	"github.com/willbicks/epigram/internal/logutils"

	quoteserver "github.com/willbicks/epigram/internal/server/http"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/inmemory"
	"github.com/willbicks/epigram/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize logger
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	slogOptions := slog.HandlerOptions{
		Level: lvl,
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slogOptions))

	// Configuration parsing
	cfg, err := config.Parse()
	if err != nil {
		log.Error("Cannot parse config to start server.", logutils.Error(err))
	}

	// Switch to pretty logging if not JSON specified
	noColor, _ := os.LookupEnv("NO_COLOR")
	if !cfg.LogJSON {
		log = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level:   lvl.Level(),
			NoColor: noColor != "",
		}))
	}

	log.Debug("Parsed config", "config", cfg)

	var userRepo service.UserRepository
	var userSessionRepo service.UserSessionRepository
	var quoteRepo service.QuoteRepository

	switch cfg.Repo {
	case config.InMemory:
		userRepo = inmemory.NewUserRepository()
		userSessionRepo = inmemory.NewUserSessionRepository()
		quoteRepo = inmemory.NewQuoteRepository()
	case config.SQLite:
		mc := &sqlite.MigrationController{}
		db, err := sql.Open("sqlite3", fmt.Sprint("file:", cfg.DBLoc, "?cache=shared&mode=rwc"))
		if err != nil {
			log.Error("unable to open database", logutils.Error(err))
			os.Exit(1)
		}

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Error("unable to close database", logutils.Error(err))
			}
		}(db)

		userRepo, err = sqlite.NewUserRepository(db, mc)
		if err != nil {
			log.Error("unable to create user repo", logutils.Error(err))
			os.Exit(1)
		}

		quoteRepo, err = sqlite.NewQuoteRepository(db, mc)
		if err != nil {
			log.Error("unable to create quote repo", logutils.Error(err))
			os.Exit(1)
		}

		userSessionRepo, err = sqlite.NewUserSessionRepository(db, mc)
		if err != nil {
			log.Error("unable to create user sess repo", logutils.Error(err))
			os.Exit(1)
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
		log.Error("critical error while initializing server", logutils.Error(err))
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	log.Info("Server starting", "addr", addr)
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
		log.Error("unable to listen and serve", logutils.Error(err))
		os.Exit(1)
	}
}
