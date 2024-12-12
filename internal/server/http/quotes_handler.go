package http

import (
	"context"
	"net/http"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/server/http/frontend"
)

func (s *QuoteServer) getQuotesPage(ctx context.Context) (frontend.QuotesPage, error) {
	quotes, err := s.QuoteService.GetAllQuotes(ctx)
	if err != nil {
		return frontend.QuotesPage{}, err
	}

	page := frontend.QuotesPage{
		Quotes: quotes,
	}

	if ctxval.UserFromContext(ctx).Admin {
		page.RenderAdmin = true

		users, err := s.UserService.GetAllUsers(ctx)
		if err != nil {
			return frontend.QuotesPage{}, err
		}

		page.Users = make(map[string]model.User)
		for _, u := range users {
			page.Users[u.ID] = u
		}
	}

	return page, nil
}

// quotesHandler handles requests to the quotes page, either GET requests to render
// quotes, or POST requests to submit a new quote
func (s *QuoteServer) quotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		page, err := s.getQuotesPage(r.Context())
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		err = s.tmpl.RenderPage(w, page)
		if err != nil {
			s.serverError(w, r, err)
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
			page, err := s.getQuotesPage(r.Context())
			if err != nil {
				s.serverError(w, r, err)
				return
			}
			page.Quote = q
			page.Error = createErr

			err = s.tmpl.RenderPage(w, page)
			if err != nil {
				s.serverError(w, r, err)
				return
			}
			return
		}
		http.Redirect(w, r, s.paths.Quotes, http.StatusSeeOther)
	default:
		s.methodNotAllowedError(w, r)
		return
	}
}
