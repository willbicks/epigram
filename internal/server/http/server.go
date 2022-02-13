package http

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/willbicks/charisms/internal/model"
	"github.com/willbicks/charisms/internal/service"

	"github.com/spf13/viper"
)

type Config struct {
	RootTD  TemplateData
	BaseURL string
}
type CharismsServer struct {
	mux   *http.ServeMux
	views map[string]*template.Template

	QuoteService service.Quote
	UserService  service.User
	QuizService  service.EntryQuiz
	gOIDC        service.OIDC

	PubFS fs.FS

	Config Config
}

func (s *CharismsServer) Init(tmplFS fs.FS) error {
	// Initialize service for Google OpenID COnnect
	s.gOIDC = service.OIDC{
		Name:         "google",
		IssuerURL:    "https://accounts.google.com",
		ClientID:     viper.GetString("googleOIDC.clientID"),
		ClientSecret: viper.GetString("googleOIDC.clientSecret"),
	}
	if err := s.gOIDC.Init(viper.GetString("baseURL")); err != nil {
		return err
	}

	// Create a http mux
	s.mux = http.NewServeMux()

	// Create a new template cache for page views
	views, err := newTemplateCache(tmplFS)
	if err != nil {
		return err
	}
	s.views = views

	// Initialize server routes
	s.routes()

	return nil
}

func (s *CharismsServer) StuffFakeData() {
	s.QuoteService.CreateQuote(context.Background(), &model.Quote{
		Quote:   "Who can I fire over that?",
		Quotee:  "Rob Lewis",
		Context: "There's a fish on my door",
	})

	s.QuoteService.CreateQuote(context.Background(), &model.Quote{
		Quote:   "Austin you have to be gay, it's for your family",
		Quotee:  "Megin",
		Context: "The matrix is a trans allegory",
	})

	s.QuoteService.CreateQuote(context.Background(), &model.Quote{
		Quote:   "Evan Craska was born and people were like \"We need a genre for this\"",
		Quotee:  "Jamieson",
		Context: "Watching pop punk music videos",
	})

	qs, err := s.QuoteService.GetAllQuotes(context.Background())
	if err != nil {
		panic(fmt.Sprintf("unable to get stuffed quotes %v", err))
	}
	fmt.Printf("created %v dummy records \n", len(qs))
}

func (s CharismsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.interpretSession(s.mux).ServeHTTP(w, r)
}

func (s *CharismsServer) routes() {
	s.mux.Handle("/", requireQuizPassed(http.HandlerFunc(s.homeHandler)))
	s.mux.Handle("/static/", s.staticHandler())
	s.mux.HandleFunc("/login", s.googleLoginHandler)
	s.mux.HandleFunc("/quiz", s.quizHandler)
	s.mux.HandleFunc("/login/google/callback", s.googleCallbackHandler)
}

func (s *CharismsServer) staticHandler() http.Handler {
	// TODO: modify to disable directory listing
	return http.StripPrefix("/static/", http.FileServer(http.FS(s.PubFS)))
}
