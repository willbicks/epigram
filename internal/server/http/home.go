package http

import (
	"fmt"
	"github.com/willbicks/charisms/internal/model"
	"github.com/willbicks/charisms/internal/service"
	"net/http"
)

// homeTD represents the template data (TD) needed to render the home page
type homeTD struct {
	Issues []string
	Quote  model.Quote
	Quotes []model.Quote
}

// homeHandler handles requests to the homepage (/)
func (s *CharismsServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// continue to render page after switch
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
		q := model.Quote{
			Quotee:  r.FormValue("quotee"),
			Quote:   r.FormValue("quote"),
			Context: r.FormValue("context"),
		}

		err := s.QuoteService.CreateQuote(r.Context(), &q)

		if err != nil {
			var issues []string

			serr, ok := err.(*service.ServiceError)
			if ok {
				issues = serr.Issues
			} else {
				issues = []string{err.Error()}
			}

			qs, err := s.QuoteService.GetAllQuotes(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			err = s.renderPage(w, "home.gohtml", homeTD{
				Issues: issues,
				Quote:  q,
				Quotes: qs,
			},
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Println(err)
			}
			return
		}

		if err := s.QuoteService.CreateQuote(r.Context(), &q); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	qs, err := s.QuoteService.GetAllQuotes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	err = s.renderPage(w, "home.gohtml", homeTD{
		Quotes: qs,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}
