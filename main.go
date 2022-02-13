package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

	quote_server "github.com/willbicks/charisms/internal/server/http"
	"github.com/willbicks/charisms/internal/service"
	"github.com/willbicks/charisms/internal/storage/inmemory"

	"github.com/spf13/viper"
)

//go:embed frontend/public
var publicEmbedFS embed.FS

//go:embed frontend/templates
var templateEmbedFS embed.FS

func main() {
	// Viper Configuration Management
	viper.SetDefault("Port", 8080)
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // TODO: establish other configuration paths

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panic("required configuation file not found: config")
		} else {
			log.Panicf("unable to read configuration file: %v", err)
		}
	}

	var entryQuestions []service.QuizQuestion
	viper.UnmarshalKey("entryQuestions", &entryQuestions)

	// embedded fs initialization
	templateFS, err := fs.Sub(templateEmbedFS, "frontend/templates")
	if err != nil {
		log.Panicf("creating templateFS: %v", err)
	}

	publicFS, err := fs.Sub(publicEmbedFS, "frontend/public")
	if err != nil {
		log.Panicf("creating publicFS: %v", err)
	}

	// Charisms Server Initialization
	cs := quote_server.CharismsServer{
		QuoteService: service.NewQuoteService(inmemory.NewQuoteRepository()),
		UserService:  service.NewUserService(inmemory.NewUserRepository(), inmemory.NewUserSessionRepository()),
		QuizService:  service.NewEntryQuizService(entryQuestions),
		PubFS:        publicFS,
		// TODO: Can viper.Unmarshall be used here?
		Config: quote_server.Config{
			BaseURL: viper.GetString("baseURL"),
			RootTD: quote_server.TemplateData{
				Title: viper.GetString("title"),
			},
		},
	}

	if err := cs.Init(templateFS); err != nil {
		log.Panicf("ciritical error while initializing server: %v", err)
	}

	cs.StuffFakeData()

	port := viper.GetInt("Port")
	log.Printf("Running server at http://localhost:%v ...", port)
	s := http.Server{
		Addr:              "localhost:" + strconv.Itoa(port),
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           cs,
	}
	s.ListenAndServe()
}
