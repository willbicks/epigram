package frontend

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed public
var publicEmbedFS embed.FS

//go:embed templates
var templateEmbedFS embed.FS

// templateFS returns the filesystem containing template files, rooted inside the templates subdirectory
func templateFS() (fs.FS, error) {
	filesys, err := fs.Sub(templateEmbedFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("creating embedFS: %v", err)
	}
	return filesys, nil
}

// PublicFS returns the filesystem containing public files (css, js, etc), rooted inside the public subdirectory
func PublicFS() (fs.FS, error) {
	filesys, err := fs.Sub(publicEmbedFS, "public")
	if err != nil {
		return nil, fmt.Errorf("creating publicFS: %v", err)
	}
	return filesys, nil
}
