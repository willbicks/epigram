package http

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"log"
	"path/filepath"
)

// newTemplateCache creates a cache of page templates from the provided filesystem, and returns
// a map relating page view names to complete templates, with components and layouts included.
func newTemplateCache(fileSys fs.FS) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	views, err := fs.Glob(fileSys, "views/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, view := range views {
		// Parse the view and all associated components and layouts
		t, err := template.ParseFS(fileSys, view)
		if err != nil {
			return nil, err
		}

		// Parse the view and all associated components and layouts
		t, err = t.ParseFS(fileSys, "components/*.gohtml", "base.gohtml")
		if err != nil {
			return nil, err
		}

		// Cache the template
		cache[filepath.Base(view)] = t

		// Print templates captured for debug
		log.Printf("cached view %s%v", view, t.DefinedTemplates())
	}

	return cache, nil
}

// renderPage renders the specified view with the provided data joined to the RootTD.
func (s *CharismsServer) renderPage(w io.Writer, name string, data interface{}) error {
	t, ok := s.views[name]
	if !ok {
		return errors.New("template not found in cache")
	}

	return t.ExecuteTemplate(w, name, s.Config.RootTD.joinPage(data))
}
