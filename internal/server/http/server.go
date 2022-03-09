package http

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/spf13/viper"

	"github.com/willbicks/epigram/internal/logger"
	"github.com/willbicks/epigram/internal/server/http/frontend"
	"github.com/willbicks/epigram/internal/server/http/paths"
	"github.com/willbicks/epigram/internal/service"
)

type Config struct {
	RootTD frontend.RootTD
	// BaseURL is the complete domain and path to access the root of this server, used for creating
	// callback URLs.
	BaseURL string
	// TrustProxy determines whether X-Forwarded-For header should be trusted to obtain the client IP,
	// or if the requestor IP shoud be used instead.
	TrustProxy bool
	// routes is a struct which stores the url paths to each page,
	// and should be used in place of magic strings to represent routes.
	paths paths.Paths
}
type QuoteServer struct {
	mux  *http.ServeMux
	tmpl frontend.TemplateEngine

	Logger logger.Logger

	QuoteService service.Quote
	UserService  service.User
	QuizService  service.EntryQuiz
	gOIDC        service.OIDC

	config Config
}

// Init initalizes the quote server, including Google OIDC proivder, http ServerMux, template engine,
// and server routes.
func (s *QuoteServer) Init(cfg Config) error {
	// TODO: Fix duplication and sprawl of route configuration, this is wholy uninituitive
	s.config = cfg
	s.config.paths = paths.Default()
	s.config.RootTD.Paths = s.config.paths

	// Initialize service for Google OpenID COnnect
	s.gOIDC = service.OIDC{
		Name:         "google",
		IssuerURL:    "https://accounts.google.com",
		ClientID:     viper.GetString("googleOIDC.clientID"),
		ClientSecret: viper.GetString("googleOIDC.clientSecret"),
	}
	if err := s.gOIDC.Init(cfg.BaseURL); err != nil {
		return err
	}

	// Create a http mux
	s.mux = http.NewServeMux()

	// Initialize serve mux
	pubFS, err := frontend.PublicFS()
	if err != nil {
		return err
	}
	s.routes(pubFS)

	// Create a new template engine
	tmpl, err := frontend.NewTemplateEngine(s.config.RootTD)
	if err != nil {
		return err
	}
	s.tmpl = tmpl
	return nil
}

// ServeHTTP serves as the entrypoint for HTTP requests to the quote server. It applies the appropriate globalmiddleware,
// and then serves request responses using the http ServeMux
func (s QuoteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gziphandler.GzipHandler(s.interpretSession(s.getIP(s.mux))).ServeHTTP(w, r)
}
