package frontend

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
)

//go:embed public
var publicEmbedFS embed.FS

//go:embed templates
var templateEmbedFS embed.FS

// templateFS returns the filesystem containing template files, rooted inside the templates subdirectory
func (e TemplateEngine) templateFS() (fs.FS, error) {
	var fsys fs.FS
	var err error
	if e.DevMode {
		fsys = os.DirFS("internal/server/http/frontend/templates")
	} else {
		fsys = templateEmbedFS
		fsys, err = fs.Sub(fsys, "templates")
	}

	if err != nil {
		return nil, fmt.Errorf("creating templateFS: %v", err)
	}
	return fsys, nil
}

// PublicFS returns the filesystem containing public files (css, js, etc.), rooted inside the public subdirectory
func (e TemplateEngine) PublicFS() (fs.FS, error) {
	var fsys fs.FS
	var err error
	if e.DevMode {
		fsys = os.DirFS("internal/server/http/frontend/public")
	} else {
		fsys = publicEmbedFS
		fsys, err = fs.Sub(fsys, "public")
	}

	if err != nil {
		return nil, fmt.Errorf("creating publicFS: %v", err)
	}
	return fsys, nil
}
