package http

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

// staticHandler accepts a file system containing files that should be publicly
// available, and returns a handler which serves them less the `/static/` prefix.
//
// If dev mode is not enabled, the contents are assumed to be unchanging,
// and as such the handler will pre-calculate hashes for all files and add
// headers for cache-control and ETags. If the file system is not an embed.FS
// (for example if DevMode is enabled and files are being served from the local
// file system), the handler will not pre-calculate hashes and will not enable
// caching.
func (s *QuoteServer) staticHandler(fileSys fs.FS) http.Handler {
	if s.Config.DevMode {
		return http.StripPrefix("/static/", http.FileServer(http.FS(fileSys)))
	}

	// If the file system is an embed.FS, pre-calculate hashes for all files
	// and enable caching with ETags.
	type file struct {
		hash string
		data []byte
	}
	files := make(map[string]file)

	// Walk the file system, load files into memory, and pre-calculate hashes for all files.
	err := fs.WalkDir(fileSys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := fileSys.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		files[path] = file{
			hash: fmt.Sprintf("%x", sha256.Sum256(data)),
			data: data,
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/static/"):]

		f, ok := files[path]
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", fmt.Sprintf(`"%s"`, f.hash))

		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, f.hash) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		http.ServeContent(w, r, path, time.Time{}, bytes.NewReader(f.data))
	})
}
