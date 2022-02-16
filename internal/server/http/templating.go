package http

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/willbicks/charisms/internal/service"
)

// initViewCache initializes a cache of view templates from the provided filesystem in the form of a
// map relating page view names to complete templates, with components and layouts included.
func (s *CharismsServer) initViewCache(fileSys fs.FS) error {
	s.views = make(map[string]*template.Template)

	views, err := fs.Glob(fileSys, "views/*.gohtml")
	if err != nil {
		return err
	}

	// define to transform an error into a list of issues.
	getIssues := func(err error) []string {
		if err == nil {
			return []string{}
		}

		serr, ok := err.(service.ServiceError)
		if ok {
			return serr.Issues
		} else {
			return []string{err.Error()}
		}
	}

	for _, view := range views {
		// Parse the view and all associated components and layouts
		t, err := template.ParseFS(fileSys, view)
		if err != nil {
			return err
		}

		t = t.Funcs(template.FuncMap{
			"getIssues": getIssues,
		})

		// Parse the view and all associated components and layouts
		t, err = t.ParseFS(fileSys, "components/*.gohtml", "base.gohtml")
		if err != nil {
			return err
		}

		// Cache the template
		s.views[filepath.Base(view)] = t

		// Print templates captured for debug
		s.Logger.Debugf("cached view %s%v", view, t.DefinedTemplates())
	}

	return nil
}

// renderPage renders the specified view with the provided data joined to the RootTD.
func (s *CharismsServer) renderPage(w io.Writer, name string, data interface{}) error {
	t, ok := s.views[name]
	if !ok {
		return errors.New("template not found in cache")
	}

	return t.ExecuteTemplate(w, name, s.Config.RootTD.joinPage(data))
}
