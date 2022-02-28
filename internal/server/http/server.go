package http

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/spf13/viper"

	"github.com/willbicks/epigram/internal/logger"
	"github.com/willbicks/epigram/internal/service"
)

type Config struct {
	RootTD RootTD
	// BaseURL is the complete domain and path to access the root of this server, used for creating
	// callback URLs.
	BaseURL string
	// TrustProxy determines whether X-Forwarded-For header should be trusted to obtain the client IP,
	// or if the requestor IP shoud be used instead.
	TrustProxy bool
	// routes is a struct which stores the url paths to each page,
	// and should be used in place of magic strings to represent routes.
	routes routeStruct
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
	s.Config.routes = routeStruct{
		Home:   "/",
		Quotes: "/quotes",
		Quiz:   "/quiz",
		Login:  "/login",
	}
	s.Config.RootTD.Routes = s.Config.routes
	s.routes(pubFS)

	return nil
}

func (s QuoteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gziphandler.GzipHandler(s.interpretSession(s.getIP(s.mux))).ServeHTTP(w, r)
}
