package handler

import (
	"net/http"
	"strings"

	"github.com/wisteriahuman/prometheus/internal/db"
)

type GraphHandler struct {
	db *db.DB
}

func NewGraphHandler(database *db.DB) *GraphHandler {
	return &GraphHandler{db: database}
}

func (h *GraphHandler) Graph(w http.ResponseWriter, r *http.Request) {
	type node struct {
		ID    string   `json:"id"`
		Path  string   `json:"path"`
		Title string   `json:"title"`
		Tags  []string `json:"tags"`
	}

	type link struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Slug   string `json:"slug"`
	}

	// Fetch all notes
	rows, err := h.db.Query("SELECT id, path, title FROM notes ORDER BY path")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query notes")
		return
	}
	defer rows.Close()

	nodes := make([]node, 0)
	pathToID := make(map[string]string)    // path -> id
	idSet := make(map[string]bool)

	for rows.Next() {
		var n node
		if err := rows.Scan(&n.ID, &n.Path, &n.Title); err != nil {
			continue
		}
		n.Tags = []string{}
		nodes = append(nodes, n)
		pathToID[n.Path] = n.ID
		idSet[n.ID] = true
	}

	// Fetch tags for each node
	for i, n := range nodes {
		tagRows, err := h.db.Query(`
			SELECT t.name FROM tags t
			JOIN note_tags nt ON nt.tag_id = t.id
			WHERE nt.note_id = ?
		`, n.ID)
		if err != nil {
			continue
		}
		for tagRows.Next() {
			var tag string
			tagRows.Scan(&tag)
			nodes[i].Tags = append(nodes[i].Tags, tag)
		}
		tagRows.Close()
	}

	// Fetch all links and resolve targets
	linkRows, err := h.db.Query("SELECT source_id, target_slug FROM links")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query links")
		return
	}
	defer linkRows.Close()

	type linkKey struct{ source, target string }
	seen := make(map[linkKey]bool)
	links := make([]link, 0)

	for linkRows.Next() {
		var sourceID, slug string
		if err := linkRows.Scan(&sourceID, &slug); err != nil {
			continue
		}

		targetID := h.resolveSlug(slug, pathToID)
		if targetID == "" || targetID == sourceID {
			continue
		}

		key := linkKey{sourceID, targetID}
		if seen[key] {
			continue
		}
		seen[key] = true

		links = append(links, link{
			Source: sourceID,
			Target: targetID,
			Slug:   slug,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"nodes": nodes,
		"links": links,
	})
}

func (h *GraphHandler) resolveSlug(slug string, pathToID map[string]string) string {
	// Try exact match: slug.md
	candidates := []string{
		slug + ".md",
		strings.ReplaceAll(slug, " ", "-") + ".md",
		strings.ReplaceAll(slug, "-", " ") + ".md",
	}

	for _, c := range candidates {
		// Case-insensitive path match
		for path, id := range pathToID {
			if strings.EqualFold(path, c) || strings.EqualFold(path, strings.ToLower(c)) {
				return id
			}
		}
	}

	return ""
}
