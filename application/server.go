package application

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/willbicks/charisms/model"
	"github.com/willbicks/charisms/service"
)

type Config struct {
	RootTD  TemplateData
	BaseURL string
}
type CharismsServer struct {
	mux          http.ServeMux
	tmpl         *template.Template
	QuoteService service.Quote
	UserService  service.User
	QuizService  service.EntryQuiz
	Cfg          Config
	gOIDC        service.OIDC
	TmplFS       fs.FS
	PubFS        fs.FS
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
	t := template.New("t")
	t, err := s.tmpl.ParseFS(s.TmplFS, "components/*.gohtml", "base.gohtml")
	if err != nil {
		log.Panicf("parsing TmplFS: %v", err)
	}
	s.tmpl = t

	// FOR DEBUG
	fmt.Println(s.tmpl.DefinedTemplates())
}

// renderPage renders the speciifed template, incorporating the base template and component templates, and
// joins the page data to the global site data.
func (s *CharismsServer) renderPage(w io.Writer, name string, data interface{}) error {
	t := template.Must(s.tmpl.Clone())
	t = template.Must(t.ParseFS(s.TmplFS, "views/"+name))
	return t.ExecuteTemplate(w, name, s.Cfg.RootTD.joinPage(data))
}

func (s *CharismsServer) routes() {
	s.mux.HandleFunc("/", s.homeHandler)
	s.mux.Handle("/static/", s.staticHandler())
	s.mux.HandleFunc("/login", s.googleLoginHandler)
	s.mux.HandleFunc("/quiz", s.quizHandler)
	s.mux.HandleFunc("/login/google/callback", s.googleCallbackHandler)
}

func (s *CharismsServer) staticHandler() http.Handler {
	// TODO: modify to disable directory listing
	return http.StripPrefix("/static/", http.FileServer(http.FS(s.PubFS)))
}
