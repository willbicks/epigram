package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/server/http/frontend"
	"github.com/willbicks/epigram/internal/storage"
)

func (s *QuoteServer) getQuotesPage(ctx context.Context) (frontend.QuotesPage, error) {
	quotes, err := s.QuoteService.GetAllQuotes(ctx)
	if err != nil {
		return frontend.QuotesPage{}, err
	}

	user := ctxval.UserFromContext(ctx)

	page := frontend.QuotesPage{
		User:   user,
		Quotes: quotes,
	}

	if user.Admin {
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
			return
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
		http.Redirect(w, r, s.paths.Quotes, http.StatusFound)
	default:
		s.methodNotAllowedError(w, r)
		return
	}
}

// quoteEditHandler handles requests to the edit a quote either GET requests to render
// an edit form, or POST requests to submit a new quote
func (s *QuoteServer) quoteEditHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := ctxval.UserFromContext(r.Context())
		q, err := s.QuoteService.GetQuote(r.Context(), r.URL.Query().Get("id"))
		if err == storage.ErrNotFound {
			s.clientError(w, r, errors.New("Quote not found"), http.StatusNotFound)
			return
		} else if err != nil {
			s.serverError(w, r, err)
			return
		}

		if !q.Editable(u) {
			s.clientError(w, r,
				errors.New("You cannot edit this quote. Quotes can only be edited by their submitters within an hour of submission."),
				http.StatusForbidden)
			return
		}

		err = s.tmpl.RenderPage(w, frontend.QuoteEditPage{
			Quote: q,
		})
		if err != nil {
			s.serverError(w, r, err)
			return
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			s.serverError(w, r, err)
			return
		}

		q, err := s.QuoteService.GetQuote(r.Context(), r.URL.Query().Get("id"))
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		q.Quote = r.FormValue("quote")
		q.Quotee = r.FormValue("quotee")
		q.Context = r.FormValue("context")

		updateErr := s.QuoteService.UpdateQuote(r.Context(), q)

		if updateErr != nil {
			err = s.tmpl.RenderPage(w, frontend.QuoteEditPage{
				Quote: q,
				Error: updateErr,
			})
			if err != nil {
				s.serverError(w, r, err)
				return
			}
			return
		}

		http.Redirect(w, r, s.paths.Quotes, http.StatusFound)

	case "DELETE":
		q, err := s.QuoteService.GetQuote(r.Context(), r.URL.Query().Get("id"))
		if err != nil {
			s.serverError(w, r, err)
			return
		}

		updateErr := s.QuoteService.DeleteQuote(r.Context(), q.ID)

		if updateErr != nil {
			err := s.tmpl.RenderPage(w, frontend.QuoteEditPage{
				Quote: q,
				Error: updateErr,
			})
			if err != nil {
				s.serverError(w, r, err)
				return
			}
			return
		}

		w.Header().Set("HX-Redirect", s.paths.Quotes)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})

	default:
		s.methodNotAllowedError(w, r)
		return
	}
}
