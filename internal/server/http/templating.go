package http

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
)

var templateFuncs template.FuncMap = template.FuncMap{
	// getIssues takes an error and returns a slice of issues. Derrived from issues slice if ServiceError,
	// otherwise a slice with a single entry of the error.Error()
	"getIssues": func(err error) []string {
		if err == nil {
			return []string{}
		}

		serr, ok := err.(service.ServiceError)
		if ok {
			return serr.Issues
		} else {
			return []string{err.Error()}
		}
	},
	// quotesByYear takes a slice of quotes, and returns them as a map where the key is the year.
	"quotesByYear": func(quotes []model.Quote) map[int][]model.Quote {
		// sort quotes from newest to oldest
		sort.Slice(quotes, func(i, j int) bool {
			return quotes[i].Created.After(quotes[j].Created)
		})

		byYear := make(map[int][]model.Quote)

		for _, q := range quotes {
			byYear[q.Created.Year()] = append(byYear[q.Created.Year()], q)
		}

		return byYear
	},
	// orderedYearKeys takes a map of quotes, and returns a slice of years (keys of map) in
	// reverse cronological order
	"orderedYearKeys": func(quotes map[int][]model.Quote) []int {
		years := make([]int, len(quotes))
		i := 0
		for y := range quotes {
			years[i] = y
			i++
		}
		sort.Slice(years, func(i, j int) bool {
			return years[i] > years[j]
		})
		return years
	},
}

// initViewCache initializes a cache of view templates from the provided filesystem in the form of a
// map relating page view names to complete templates, with components and layouts included.
func (s *QuoteServer) initViewCache(fileSys fs.FS) error {
	s.views = make(map[string]*template.Template)

	views, err := fs.Glob(fileSys, "views/*.gohtml")
	if err != nil {
		return err
	}

	for _, view := range views {
		name := filepath.Base(view)
		// Parse the view with functions
		t, err := template.New(name).Funcs(templateFuncs).ParseFS(fileSys, view)
		if err != nil {
			return err
		}

		// Parse associated components and layouts
		t, err = t.ParseFS(fileSys, "components/*.gohtml", "base.gohtml")
		if err != nil {
			return err
		}

		// Cache the template
		s.views[name] = t

		// Print templates captured for debug
		s.Logger.Debugf("cached view %s%v", view, t.DefinedTemplates())
	}

	return nil
}

// renderPage renders the specified view with the provided data joined to the RootTD.
func (s *QuoteServer) renderPage(w io.Writer, name string, data interface{}) error {
	t, ok := s.views[name]
	if !ok {
		return errors.New("template not found in cache")
	}

	return t.ExecuteTemplate(w, name, s.Config.RootTD.joinPage(data))
}
