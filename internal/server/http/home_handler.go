package http

import (
	"net/http"

	"github.com/willbicks/epigram/internal/model"
)

// homeTD represents the template data (TD) needed to render the home page
type homeTD struct {
	Error  error
	Quote  model.Quote
	Quotes []model.Quote
}

// homeHandler handles requests to the homepage (/)
func (s *QuoteServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		quotes, err := s.QuoteService.GetAllQuotes(r.Context())
		if err != nil {
			s.serverError(w, r, err)
			return
		}
		err = s.renderPage(w, "home.gohtml", homeTD{
			Quotes: quotes,
		})
		if err != nil {
			s.serverError(w, r, err)
			s.Logger.Warn(err.Error())
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			s.serverError(w, r, err)
			return
		}
		q := model.Quote{
			Quotee:  r.FormValue("quotee"),
			Quote:   r.FormValue("quote"),
			Context: r.FormValue("context"),
		}

		createErr := s.QuoteService.CreateQuote(r.Context(), &q)

		if createErr != nil {
			quotes, err := s.QuoteService.GetAllQuotes(r.Context())
			if err != nil {
				s.serverError(w, r, err)
				return
			}
			err = s.renderPage(w, "home.gohtml", homeTD{
				Error:  createErr,
				Quote:  q,
				Quotes: quotes,
			})
			if err != nil {
				s.serverError(w, r, err)
				return
			}
			return
		}
		http.Redirect(w, r, paths.home, http.StatusFound)
	default:
		s.methodNotAllowedError(w, r)
		return
	}
}
