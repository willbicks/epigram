package frontend

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
)

// TemplateEngine is responsible for storing cached html templates, and renderning them ondemand with
// template data.
type TemplateEngine struct {
	views  map[string]*template.Template
	rootTD RootTD
}

// NewTemplateEngine creates a new template engine, and initializes it's views so that it is ready
// to render responses.
func NewTemplateEngine(rootTD RootTD) (TemplateEngine, error) {
	e := TemplateEngine{
		rootTD: rootTD,
	}

	tmplFS, err := templateFS()
	if err != nil {
		return TemplateEngine{}, err
	}

	e.initViews(tmplFS)

	return e, nil
}

// initViews initializes a cache of view templates from the provided filesystem in the form of a
// map relating page view names to complete templates, with components and layouts included.
func (e *TemplateEngine) initViews(fileSys fs.FS) error {
	e.views = make(map[string]*template.Template)

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

		// Store the template
		e.views[name] = t
	}

	return nil
}

// RenderPage renders the specified view with the provided data joined to the RootTD.
func (e *TemplateEngine) RenderPage(w io.Writer, name string, data interface{}) error {
	t, ok := e.views[name]
	if !ok {
		return errors.New("template not found in views")
	}

	return t.ExecuteTemplate(w, name, e.rootTD.joinPage(data))
}
