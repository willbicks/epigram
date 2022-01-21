package application

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/willbicks/charisms/model"
	"github.com/willbicks/charisms/service"
)

type Config struct {
	RootTD     TemplateData
	ViewsPath  string
	PublicPath string
	BaseURL    string
}
type CharismsServer struct {
	mux          http.ServeMux
	tmpl         *template.Template
	QuoteService service.Quote
	UserService  service.User
	Cfg          Config
	gOIDC        service.OIDC
}

func (s *CharismsServer) Init() {
	s.gOIDC = service.OIDC{
		Name:         "google",
		IssuerURL:    "https://accounts.google.com",
		ClientID:     viper.GetString("googleOIDC.clientID"),
		ClientSecret: viper.GetString("googleOIDC.clientSecret"),
	}

	if err := s.gOIDC.Init(viper.GetString("baseURL")); err != nil {
		log.Panic(err)
	}
	s.templates()
	s.routes()
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
	s.mux.ServeHTTP(w, r)
}

func (s *CharismsServer) templates() {
	args := func(vs ...interface{}) []interface{} { return vs }
	s.tmpl = template.New("t").Funcs(template.FuncMap{"args": args})
	// TODO: use embed and template.ParseFS to embed html in go binary
	// TODO: use config value for template directory
	s.tmpl = template.Must(s.tmpl.ParseGlob("frontend/views/components/*.gohtml"))
	s.tmpl = template.Must(s.tmpl.ParseGlob("frontend/views/*.gohtml"))
	fmt.Println(s.tmpl.DefinedTemplates())
}

func (s *CharismsServer) routes() {
	s.mux.HandleFunc("/", s.homeHandler)
	s.mux.Handle("/static/", s.staticHandler())
	s.mux.HandleFunc("/login", s.googleLoginHandler)
	s.mux.HandleFunc("/login/google/callback", s.googleCallbackHandler)
}

func (s *CharismsServer) staticHandler() http.Handler {
	// also requires refactor for embed
	// TODO: modify to disable directory listing
	return http.StripPrefix("/static/", http.FileServer(http.Dir(s.Cfg.PublicPath)))
}
