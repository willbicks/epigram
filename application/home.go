package application

import (
	"fmt"
	"net/http"

	"github.com/willbicks/charisms/model"
)

// homeTD represents the template data (TD) needed to render the home page
type homeTD struct {
	Errors []string
	Quote  model.Quote
	Quotes []model.Quote
}

// homeHandler handles requests to the homepage (/)
func (s *CharismsServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

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

		var issues []string
		if q.Quote == "" {
			issues = append(issues, "Please enter a quote.")
		}
		if q.Quotee == "" {
			issues = append(issues, "Please attribute this quote to someone.")
		}

		if len(issues) > 0 {
			qs, err := s.QuoteService.GetAllQuotes(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			err = s.tmpl.ExecuteTemplate(w, "home.gohtml", s.Config.RootTD.joinPage(
				homeTD{
					Errors: issues,
					Quote:  q,
					Quotes: qs,
				},
			))
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

	err = s.tmpl.ExecuteTemplate(w, "home.gohtml", s.Config.RootTD.joinPage(homeTD{
		Quotes: qs,
	}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}
