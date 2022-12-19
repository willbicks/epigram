package frontend

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
)

// TemplateEngine is responsible for storing cached html templates, and rendering them on-demand with
// template data.
type TemplateEngine struct {
	views  map[string]*template.Template
	rootTD RootTD
	// If devMode is true, the TemplateEngine will use local files instead of the embedded filesystems,
	// and will reload the templates on every render.
	DevMode bool
}

// NewTemplateEngine creates a new template engine, and initializes it's views so that it is ready
// to render responses.
func NewTemplateEngine(rootTD RootTD) (TemplateEngine, error) {
	e := TemplateEngine{
		rootTD: rootTD,
	}

	err := e.initViews()
	if err != nil {
		return TemplateEngine{}, err
	}

	return e, nil
}

// initViews initializes a cache of view templates from the provided filesystem in the form of a
// map relating page view names to complete templates, with components and layouts included.
func (e *TemplateEngine) initViews() error {
	tmplFS, err := e.templateFS()
	if err != nil {
		return err
	}

	e.views = make(map[string]*template.Template)

	views, err := fs.Glob(tmplFS, "views/*.gohtml")
	if err != nil {
		return err
	}

	for _, view := range views {
		name := filepath.Base(view)
		// Parse the view with functions
		t, err := template.New(name).Funcs(templateFuncs).ParseFS(tmplFS, view)
		if err != nil {
			return err
		}

		// Parse associated components and layouts
		t, err = t.ParseFS(tmplFS, "components/*.gohtml", "base.gohtml")
		if err != nil {
			return err
		}

		// Store the template
		e.views[name] = t
	}

	return nil
}

// RenderPage renders the specified view with the provided data joined to the RootTD.
func (e TemplateEngine) RenderPage(w io.Writer, page Page) error {

	// If devMode is true, reload the templates on every render
	if e.DevMode {
		err := e.initViews()
		if err != nil {
			return err
		}
	}

	t, ok := e.views[page.viewName()]
	if !ok {
		return errors.New("template not found in views")
	}

	return t.ExecuteTemplate(w, page.viewName(), e.rootTD.joinPage(page))
}
