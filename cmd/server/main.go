package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/willbicks/epigram/internal/logger"
	quote_server "github.com/willbicks/epigram/internal/server/http"
	"github.com/willbicks/epigram/internal/server/http/frontend"
	"github.com/willbicks/epigram/internal/service"
	"github.com/willbicks/epigram/internal/storage/inmemory"
	"github.com/willbicks/epigram/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {
	// Initialize logger
	log := logger.New(os.Stdout, true)
	log.Level = logger.LevelDebug

	// Viper Configuration Management
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Bind", "0.0.0.0")
	viper.SetDefault("database", "inmemory")
	viper.SetDefault("sqlite-db", "/var/epigram/database.db")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/epigram")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("required configuation file not found: config")
		} else {
			log.Fatalf("unable to read configuration file: %v", err)
		}
	}

	var entryQuestions []service.QuizQuestion
	viper.UnmarshalKey("entryQuestions", &entryQuestions)

	var userRepo service.UserRepository
	var userSessionRepo service.UserSessionRepository
	var quoteRepo service.QuoteRepository

	switch viper.GetString("database") {
	case "inmemory":
		userRepo = inmemory.NewUserRepository()
		userSessionRepo = inmemory.NewUserSessionRepository()
		quoteRepo = inmemory.NewQuoteRepository()
	case "sqlite":
		mc := &sqlite.MigrationController{}
		db, err := sql.Open("sqlite3", fmt.Sprint("file:", viper.GetString("sqlite-db"), "?cache=shared&mode=rwc"))
		if err != nil {
			log.Fatalf("unable to open database: %v", err)
		}
		defer db.Close()

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
	cs := quote_server.QuoteServer{
		QuoteService: service.NewQuoteService(quoteRepo),
		UserService:  service.NewUserService(userRepo, userSessionRepo),
		QuizService:  service.NewEntryQuizService(entryQuestions),
		Logger:       log,
	}

	cfg := quote_server.Config{
		BaseURL: viper.GetString("baseURL"),
		RootTD: frontend.RootTD{
			Title:       viper.GetString("title"),
			Description: viper.GetString("description"),
		},
	}

	if err := cs.Init(cfg); err != nil {
		log.Fatalf("ciritical error while initializing server: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", viper.GetString("bind"), viper.GetInt("port"))
	log.Infof("Running server at %s ...", addr)
	s := http.Server{
		Addr:              addr,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           cs,
	}
	s.ListenAndServe()
}
