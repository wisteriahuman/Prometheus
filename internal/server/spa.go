package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:web
var webFS embed.FS

func spaHandler() http.Handler {
	fsys, err := fs.Sub(webFS, "web")
	if err != nil {
		// Fallback: serve a simple message if embed fails
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<h1>Prometheus</h1><p>Web UI not found. Run 'make web' first.</p>"))
		})
	}

	fileServer := http.FileServer(http.FS(fsys))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Skip API routes
		if strings.HasPrefix(path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the file directly
		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		if _, err := fs.Stat(fsys, cleanPath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for all non-file routes
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
