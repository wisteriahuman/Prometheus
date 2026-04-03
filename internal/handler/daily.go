package handler

import (
	"net/http"
	"time"

	"github.com/wisteriahuman/prometheus/internal/service"
)

type DailyHandler struct {
	daily   *service.DailyService
	indexer *service.Indexer
	md      *service.Markdown
	vault   *service.Vault
}

func NewDailyHandler(daily *service.DailyService, indexer *service.Indexer, md *service.Markdown, vault *service.Vault) *DailyHandler {
	return &DailyHandler{daily: daily, indexer: indexer, md: md, vault: vault}
}

func (h *DailyHandler) GetDaily(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	var date time.Time
	if dateStr == "" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
			return
		}
	}

	note, err := h.daily.EnsureDailyNote(date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to ensure daily note: "+err.Error())
		return
	}

	h.indexer.IndexNote(note.Path)

	html := h.md.ToHTML(note.RawContent)

	// Recent daily notes
	type recentNote struct {
		Date   string `json:"date"`
		Exists bool   `json:"exists"`
	}

	recentDates := service.RecentDailyDates(7)
	recentNotes := make([]recentNote, 0, len(recentDates))
	for _, d := range recentDates {
		p := service.DailyNotePath(d)
		recentNotes = append(recentNotes, recentNote{
			Date:   d.Format("2006-01-02"),
			Exists: h.vault.NoteExists(p),
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"path":        note.Path,
		"title":       note.Title,
		"content":     note.Content,
		"frontmatter": note.Frontmatter,
		"html":        html,
		"currentDate": date.Format("2006-01-02"),
		"recentNotes": recentNotes,
	})
}
