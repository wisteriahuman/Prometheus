package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wisteriahuman/prometheus/internal/db"
	"github.com/wisteriahuman/prometheus/internal/service"
)

type TasksHandler struct {
	db      *db.DB
	vault   *service.Vault
	indexer *service.Indexer
}

func NewTasksHandler(database *db.DB, vault *service.Vault, indexer *service.Indexer) *TasksHandler {
	return &TasksHandler{db: database, vault: vault, indexer: indexer}
}

func (h *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}

	query := `
		SELECT t.id, t.note_id, t.content, t.completed, t.line_number, t.due_date,
		       n.path, n.title
		FROM tasks t
		JOIN notes n ON n.id = t.note_id
	`

	switch filter {
	case "pending":
		query += " WHERE t.completed = 0"
	case "completed":
		query += " WHERE t.completed = 1"
	}

	query += " ORDER BY t.completed ASC, t.id ASC"

	rows, err := h.db.Query(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query tasks")
		return
	}
	defer rows.Close()

	type taskResult struct {
		ID         int64   `json:"id"`
		NoteID     string  `json:"noteId"`
		Content    string  `json:"content"`
		Completed  bool    `json:"completed"`
		LineNumber int     `json:"lineNumber"`
		DueDate    *string `json:"dueDate"`
		NotePath   string  `json:"notePath"`
		NoteTitle  string  `json:"noteTitle"`
	}

	tasks := make([]taskResult, 0)
	for rows.Next() {
		var t taskResult
		var completed int
		var dueDate *string
		if err := rows.Scan(&t.ID, &t.NoteID, &t.Content, &completed, &t.LineNumber, &dueDate, &t.NotePath, &t.NoteTitle); err != nil {
			continue
		}
		t.Completed = completed == 1
		t.DueDate = dueDate
		tasks = append(tasks, t)
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TasksHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var body struct {
		Completed bool `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get task info
	var noteID string
	var lineNumber int
	err = h.db.QueryRow("SELECT note_id, line_number FROM tasks WHERE id = ?", taskID).Scan(&noteID, &lineNumber)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	// Get note path
	var notePath string
	err = h.db.QueryRow("SELECT path FROM notes WHERE id = ?", noteID).Scan(&notePath)
	if err != nil {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}

	// Read the note and toggle the checkbox
	note, err := h.vault.ReadNote(notePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read note")
		return
	}

	// Toggle checkbox in raw content
	lines := strings.Split(note.RawContent, "\n")
	if lineNumber < 1 || lineNumber > len(lines) {
		writeError(w, http.StatusBadRequest, "line number out of range")
		return
	}

	line := lines[lineNumber-1]
	if body.Completed {
		line = strings.Replace(line, "- [ ] ", "- [x] ", 1)
	} else {
		line = strings.Replace(line, "- [x] ", "- [ ] ", 1)
		line = strings.Replace(line, "- [X] ", "- [ ] ", 1)
	}
	lines[lineNumber-1] = line

	// Re-parse and write
	newRaw := strings.Join(lines, "\n")
	fm, content := service.ParseFrontmatterPublic([]byte(newRaw))

	if _, err := h.vault.WriteNote(notePath, content, fm); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to write note: %v", err))
		return
	}

	h.indexer.IndexNote(notePath)

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
