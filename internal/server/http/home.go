package http

import (
	"net/http"

	"github.com/willbicks/charisms/internal/model"
)

// homeTD represents the template data (TD) needed to render the home page
type homeTD struct {
	Error  error
	Quote  model.Quote
	Quotes []model.Quote
}

// homeHandler handles requests to the homepage (/)
func (s *CharismsServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		qs, err := s.QuoteService.GetAllQuotes(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.Logger.Warn(err.Error())
			return
		}

		err = s.renderPage(w, "home.gohtml", homeTD{
			Quotes: qs,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.Logger.Warn(err.Error())
		}
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

		createErr := s.QuoteService.CreateQuote(r.Context(), &q)

		if createErr != nil {
			qs, err := s.QuoteService.GetAllQuotes(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				s.Logger.Warn(err.Error())
				return
			}

			err = s.renderPage(w, "home.gohtml", homeTD{
				Error:  createErr,
				Quote:  q,
				Quotes: qs,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				s.Logger.Warn(err.Error())
			}
			return
		}
		http.Redirect(w, r, paths.home, http.StatusFound)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}
}
