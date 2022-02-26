package http

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/willbicks/epigram/internal/logger"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"

	"github.com/spf13/viper"
)

type Config struct {
	RootTD     TemplateData
	BaseURL    string
	TrustProxy bool
}
type QuoteServer struct {
	mux   *http.ServeMux
	views map[string]*template.Template

	Logger logger.Logger

	QuoteService service.Quote
	UserService  service.User
	QuizService  service.EntryQuiz
	gOIDC        service.OIDC

	Config Config
}

func (s *QuoteServer) Init(tmplFS fs.FS, pubFS fs.FS) error {
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

	if err := s.initViewCache(tmplFS); err != nil {
		return err
	}

	// Initialize server routes
	s.routes(pubFS)

	return nil
}

func (s *QuoteServer) StuffFakeData() {
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
		s.Logger.Fatalf("unable to get stuffed quotes: %v", err)
	}
	s.Logger.Infof("created %v dummy records \n", len(qs))
}

func (s QuoteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.interpretSession(s.getIP(s.mux)).ServeHTTP(w, r)
}
