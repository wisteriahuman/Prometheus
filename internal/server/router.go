package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wisteriahuman/prometheus/internal/handler"
)

type Handlers struct {
	Notes     *handler.NotesHandler
	Search    *handler.SearchHandler
	Graph     *handler.GraphHandler
	Tasks     *handler.TasksHandler
	Tags      *handler.TagsHandler
	Daily     *handler.DailyHandler
	Preview   *handler.PreviewHandler
	Tree      *handler.TreeHandler
	Folders   *handler.FoldersHandler
	Backlinks *handler.BacklinksHandler
	Config    *handler.ConfigHandler
	Export    *handler.ExportHandler
}

func NewRouter(h *Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(corsMiddleware)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Notes
		r.Get("/notes", h.Notes.List)
		r.Post("/notes", h.Notes.Create)
		r.Get("/notes/*", h.Notes.Read)
		r.Put("/notes/*", h.Notes.Update)
		r.Delete("/notes/*", h.Notes.Delete)

		// Search
		r.Get("/search", h.Search.Search)

		// Graph
		r.Get("/graph", h.Graph.Graph)

		// Tasks
		r.Get("/tasks", h.Tasks.List)
		r.Patch("/tasks/{id}", h.Tasks.Toggle)

		// Tags
		r.Get("/tags", h.Tags.List)
		r.Get("/tags/{name}", h.Tags.NotesByTag)

		// Daily
		r.Get("/daily", h.Daily.GetDaily)

		// Preview
		r.Post("/preview", h.Preview.Preview)

		// Tree
		r.Get("/tree", h.Tree.GetTree)

		// Folders
		r.Post("/folders/*", h.Folders.Create)
		r.Delete("/folders/*", h.Folders.Delete)

		// Backlinks
		r.Get("/backlinks/*", h.Backlinks.GetBacklinks)

		// Config
		r.Get("/config", h.Config.GetConfig)
		r.Put("/config", h.Config.UpdateConfig)

		// Export
		r.Get("/export/*", h.Export.Export)
	})

	// SPA static file serving (catch-all)
	r.Handle("/*", spaHandler())

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
