package http

import (
	"net/http"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/server/http/frontend"
)

// quotesHandler handles requests to the quotes page, either GET requests to render
// quotes, or POST requests to submit a new quote
func (s *QuoteServer) quotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		quotes, err := s.QuoteService.GetAllQuotes(r.Context())
		if err != nil {
			s.serverError(w, r, err)
			return
		}
		err = s.tmpl.RenderPage(w, "quotes.gohtml", frontend.QuotesTD{
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
			err = s.tmpl.RenderPage(w, "quotes.gohtml", frontend.QuotesTD{
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
		http.Redirect(w, r, s.paths.Quotes, http.StatusFound)
	default:
		s.methodNotAllowedError(w, r)
		return
	}
}
