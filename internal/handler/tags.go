package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wisteriahuman/prometheus/internal/db"
)

type TagsHandler struct {
	db *db.DB
}

func NewTagsHandler(database *db.DB) *TagsHandler {
	return &TagsHandler{db: database}
}

func (h *TagsHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT t.id, t.name, COUNT(nt.note_id) as count
		FROM tags t
		LEFT JOIN note_tags nt ON nt.tag_id = t.id
		GROUP BY t.id, t.name
		ORDER BY t.name
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query tags")
		return
	}
	defer rows.Close()

	type tagResult struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	tags := make([]tagResult, 0)
	for rows.Next() {
		var t tagResult
		if err := rows.Scan(&t.ID, &t.Name, &t.Count); err != nil {
			continue
		}
		tags = append(tags, t)
	}

	writeJSON(w, http.StatusOK, tags)
}

func (h *TagsHandler) NotesByTag(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "tag name is required")
		return
	}

	rows, err := h.db.Query(`
		SELECT n.id, n.path, n.title, n.modified_at
		FROM notes n
		JOIN note_tags nt ON nt.note_id = n.id
		JOIN tags t ON t.id = nt.tag_id
		WHERE t.name = ?
		ORDER BY n.modified_at DESC
	`, name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query notes by tag")
		return
	}
	defer rows.Close()

	type noteResult struct {
		ID         string `json:"id"`
		Path       string `json:"path"`
		Title      string `json:"title"`
		ModifiedAt int64  `json:"modifiedAt"`
	}

	notes := make([]noteResult, 0)
	for rows.Next() {
		var n noteResult
		if err := rows.Scan(&n.ID, &n.Path, &n.Title, &n.ModifiedAt); err != nil {
			continue
		}
		notes = append(notes, n)
	}

	writeJSON(w, http.StatusOK, notes)
}
