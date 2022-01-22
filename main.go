package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/willbicks/charisms/application"
	"github.com/willbicks/charisms/service"
	"github.com/willbicks/charisms/storage/inmemory"
)

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

	// Charisms Server Initialization
	cs := application.CharismsServer{
		QuoteService: service.NewQuoteService(inmemory.NewQuoteRepository()),
		UserService:  service.NewUserService(inmemory.NewUserRepository()),
		QuizService:  service.NewEntryQuizService(entryQuestions),
		// TODO: Can viper.Unmarshall be used here?
		Cfg: application.Config{
			BaseURL:    viper.GetString("baseURL"),
			ViewsPath:  viper.GetString("viewsPath"),
			PublicPath: viper.GetString("publicPath"),
			RootTD: application.TemplateData{
				Title: viper.GetString("title"),
			},
		},
	}
	cs.Init()
	cs.StuffFakeData()

	port := viper.GetInt("Port")
	log.Printf("Running server at http://localhost:%v ...", port)
	s := http.Server{
		Addr:         "localhost:" + strconv.Itoa(port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      cs,
	}
	s.ListenAndServe()
}
