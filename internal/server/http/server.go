package http

import (
	"log/slog"
	"net/http"

	"github.com/klauspost/compress/gzhttp"

	"github.com/willbicks/epigram/internal/config"
	"github.com/willbicks/epigram/internal/server/http/frontend"
	"github.com/willbicks/epigram/internal/server/http/paths"
	"github.com/willbicks/epigram/internal/service"
)

// QuoteServer provides a web UI over HTTP for the application
type QuoteServer struct {
	mux  *http.ServeMux
	tmpl frontend.TemplateEngine

	Logger *slog.Logger

	QuoteService service.Quote
	UserService  service.User
	QuizService  service.EntryQuiz
	OIDCService  service.OIDC

	// paths is a struct which stores the url paths to each page,
	// and should be used in place of magic strings to represent rout
	paths paths.Paths

	Config config.Application
}

// Init initializes the quote server, including Google OIDC provider, http ServerMux, template engine,
// and server routes.
func (s *QuoteServer) Init() error {
	// Initialize paths
	s.paths = paths.Default()

	// Initialize service for OpenID Connect
	s.OIDCService = service.OIDC{
		Name:         s.Config.OIDCProvider.Name,
		IssuerURL:    s.Config.OIDCProvider.IssuerURL,
		ClientID:     s.Config.OIDCProvider.ClientID,
		ClientSecret: s.Config.OIDCProvider.ClientSecret,
	}
	if err := s.OIDCService.Init(s.Config.BaseURL); err != nil {
		return err
	}

	// Initialize template engine
	tmpl, err := frontend.NewTemplateEngine(frontend.RootTD{
		Title:       s.Config.Title,
		Description: s.Config.Description,
		Paths:       s.paths,
	})
	if err != nil {
		return err
	}
	if s.Config.DevMode {
		s.Logger.Warn("Running in development mode. Template engine performance reduced.")
		tmpl.DevMode = true
	}
	s.tmpl = tmpl

	// Initialize http mux and routes
	s.mux = http.NewServeMux()
	pubFS, err := tmpl.PublicFS()
	if err != nil {
		return err
	}
	s.routes(pubFS)

	return nil
}

// ServeHTTP serves as the entrypoint for HTTP requests to the quote server. It applies the appropriate global middleware,
// and then serves request responses using the http ServeMux
func (s QuoteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gzhttp.GzipHandler(s.interpretSession(s.getIP(s.mux))).ServeHTTP(w, r)
}
