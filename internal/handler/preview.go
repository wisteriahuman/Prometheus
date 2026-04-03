package handler

import (
	"encoding/json"
	"net/http"

	"github.com/wisteriahuman/prometheus/internal/service"
)

type PreviewHandler struct {
	md *service.Markdown
}

func NewPreviewHandler(md *service.Markdown) *PreviewHandler {
	return &PreviewHandler{md: md}
}

func (h *PreviewHandler) Preview(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	html := h.md.ToHTML(body.Content)

	writeJSON(w, http.StatusOK, map[string]string{"html": html})
}
