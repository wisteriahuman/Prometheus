package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wisteriahuman/prometheus/internal/service"
)

type NotesHandler struct {
	vault   *service.Vault
	indexer *service.Indexer
}

func NewNotesHandler(vault *service.Vault, indexer *service.Indexer) *NotesHandler {
	return &NotesHandler{vault: vault, indexer: indexer}
}

func (h *NotesHandler) List(w http.ResponseWriter, r *http.Request) {
	paths, err := h.vault.ListNotes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list notes")
		return
	}

	type noteItem struct {
		Path     string   `json:"path"`
		Title    string   `json:"title"`
		Tags     []string `json:"tags"`
		Modified string   `json:"modified"`
		Theme    *string  `json:"theme"`
	}

	notes := make([]noteItem, 0, len(paths))
	for _, p := range paths {
		note, err := h.vault.ReadNote(p)
		if err != nil {
			continue
		}
		notes = append(notes, noteItem{
			Path:     note.Path,
			Title:    note.Frontmatter.Title,
			Tags:     note.Frontmatter.Tags,
			Modified: note.Frontmatter.Modified,
			Theme:    note.Frontmatter.Theme,
		})
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *NotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Path  string `json:"path"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}
	if !strings.HasSuffix(body.Path, ".md") {
		body.Path += ".md"
	}

	note, err := h.vault.CreateNote(body.Path, body.Title)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create note: "+err.Error())
		return
	}

	h.indexer.IndexNote(note.Path)

	writeJSON(w, http.StatusCreated, map[string]string{
		"path":  note.Path,
		"title": note.Title,
		"id":    note.Frontmatter.ID,
	})
}

func (h *NotesHandler) Read(w http.ResponseWriter, r *http.Request) {
	notePath := chi.URLParam(r, "*")
	if notePath == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	note, err := h.vault.ReadNote(notePath)
	if err != nil {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"path":        note.Path,
		"title":       note.Title,
		"content":     note.Content,
		"frontmatter": note.Frontmatter,
		"checksum":    note.Checksum,
	})
}

func (h *NotesHandler) Update(w http.ResponseWriter, r *http.Request) {
	notePath := chi.URLParam(r, "*")
	if notePath == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	var body struct {
		Content     string                  `json:"content"`
		Frontmatter map[string]interface{} `json:"frontmatter"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Read existing note to get current frontmatter
	existing, err := h.vault.ReadNote(notePath)
	if err != nil {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}

	fm := existing.Frontmatter

	// Merge frontmatter if provided
	if body.Frontmatter != nil {
		if v, ok := body.Frontmatter["title"]; ok {
			if s, ok := v.(string); ok {
				fm.Title = s
			}
		}
		if v, ok := body.Frontmatter["tags"]; ok {
			if arr, ok := v.([]interface{}); ok {
				tags := make([]string, 0, len(arr))
				for _, t := range arr {
					if s, ok := t.(string); ok {
						tags = append(tags, s)
					}
				}
				fm.Tags = tags
			}
		}
		// Handle theme: null deletes, string sets
		if v, exists := body.Frontmatter["theme"]; exists {
			if v == nil {
				fm.Theme = nil
			} else if s, ok := v.(string); ok {
				fm.Theme = &s
			}
		}
	}

	note, err := h.vault.WriteNote(notePath, body.Content, fm)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update note: "+err.Error())
		return
	}

	h.indexer.IndexNote(note.Path)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"path":     note.Path,
		"title":    note.Title,
		"checksum": note.Checksum,
	})
}

func (h *NotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	notePath := chi.URLParam(r, "*")
	if notePath == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	if err := h.vault.DeleteNote(notePath); err != nil {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}

	h.indexer.RemoveNote(notePath)

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
