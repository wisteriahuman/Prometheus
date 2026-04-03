package handler

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/wisteriahuman/prometheus/internal/db"
)

type SearchHandler struct {
	db *db.DB
}

func NewSearchHandler(database *db.DB) *SearchHandler {
	return &SearchHandler{db: database}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}

	type result struct {
		ID         string `json:"id"`
		Path       string `json:"path"`
		Title      string `json:"title"`
		Snippet    string `json:"snippet"`
		ModifiedAt int64  `json:"modifiedAt"`
	}

	results := make([]result, 0)

	// Try FTS5 first
	rows, err := h.searchFTS5(q)
	if err != nil {
		// Fallback to LIKE
		rows, err = h.searchLIKE(q)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "search failed")
			return
		}
	}
	defer rows.Close()

	for rows.Next() {
		var id, path, title, content string
		var modifiedAt int64
		if err := rows.Scan(&id, &path, &title, &content, &modifiedAt); err != nil {
			continue
		}

		snippet := generateSnippet(content, q)
		if snippet == "" {
			snippet = generateSnippet(title, q)
		}

		results = append(results, result{
			ID:         id,
			Path:       path,
			Title:      title,
			Snippet:    snippet,
			ModifiedAt: modifiedAt,
		})
	}

	writeJSON(w, http.StatusOK, results)
}

func (h *SearchHandler) searchFTS5(q string) (*sql.Rows, error) {
	// FTS5 MATCH query — escape special characters
	ftsQuery := strings.ReplaceAll(q, `"`, `""`)

	return h.db.Query(`
		SELECT n.id, n.path, n.title, n.content, n.modified_at
		FROM notes_fts f
		JOIN notes n ON n.rowid = f.rowid
		WHERE notes_fts MATCH ?
		LIMIT 50
	`, `"`+ftsQuery+`"`)
}

func (h *SearchHandler) searchLIKE(q string) (*sql.Rows, error) {
	like := "%" + q + "%"
	return h.db.Query(`
		SELECT id, path, title, content, modified_at
		FROM notes
		WHERE title LIKE ? OR content LIKE ?
		LIMIT 50
	`, like, like)
}

func generateSnippet(text, term string) string {
	lower := strings.ToLower(text)
	lowerTerm := strings.ToLower(term)
	idx := strings.Index(lower, lowerTerm)
	if idx < 0 {
		return ""
	}

	start := idx - 60
	if start < 0 {
		start = 0
	}
	end := idx + len(term) + 60
	if end > len(text) {
		end = len(text)
	}

	snippet := text[start:end]
	snippet = strings.ReplaceAll(snippet, "\n", " ")

	prefix := ""
	suffix := ""
	if start > 0 {
		prefix = "..."
	}
	if end < len(text) {
		suffix = "..."
	}

	return prefix + strings.TrimSpace(snippet) + suffix
}
