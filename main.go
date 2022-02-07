package main

import (
	"embed"
	http2 "github.com/willbicks/charisms/internal/server/http"
	service2 "github.com/willbicks/charisms/internal/service"
	inmemory2 "github.com/willbicks/charisms/internal/storage/inmemory"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"

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

	var entryQuestions []service2.QuizQuestion
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
	cs := http2.CharismsServer{
		QuoteService: service2.NewQuoteService(inmemory2.NewQuoteRepository()),
		UserService:  service2.NewUserService(inmemory2.NewUserRepository(), inmemory2.NewUserSessionRepository()),
		QuizService:  service2.NewEntryQuizService(entryQuestions),
		TmplFS:       templateFS,
		PubFS:        publicFS,
		// TODO: Can viper.Unmarshall be used here?
		Cfg: http2.Config{
			BaseURL: viper.GetString("baseURL"),
			RootTD: http2.TemplateData{
				Title: viper.GetString("title"),
			},
		},
	}
	cs.Init()
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
