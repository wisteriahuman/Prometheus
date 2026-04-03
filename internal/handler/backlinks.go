package handler

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wisteriahuman/prometheus/internal/db"
)

type BacklinksHandler struct {
	db *db.DB
}

func NewBacklinksHandler(database *db.DB) *BacklinksHandler {
	return &BacklinksHandler{db: database}
}

func (h *BacklinksHandler) GetBacklinks(w http.ResponseWriter, r *http.Request) {
	notePath := chi.URLParam(r, "*")
	if notePath == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	// Extract slug from path: filename without .md, case-insensitive
	slug := strings.TrimSuffix(filepath.Base(notePath), ".md")

	rows, err := h.db.Query(`
		SELECT l.source_id, l.target_slug, l.context, n.path, n.title
		FROM links l
		JOIN notes n ON n.id = l.source_id
		WHERE LOWER(l.target_slug) LIKE LOWER(?)
	`, slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query backlinks")
		return
	}
	defer rows.Close()

	type backlink struct {
		SourceID   string `json:"sourceId"`
		TargetSlug string `json:"targetSlug"`
		Context    string `json:"context"`
		SourcePath string `json:"sourcePath"`
		SourceTitle string `json:"sourceTitle"`
	}

	links := make([]backlink, 0)
	for rows.Next() {
		var bl backlink
		if err := rows.Scan(&bl.SourceID, &bl.TargetSlug, &bl.Context, &bl.SourcePath, &bl.SourceTitle); err != nil {
			continue
		}
		links = append(links, bl)
	}

	writeJSON(w, http.StatusOK, links)
}
