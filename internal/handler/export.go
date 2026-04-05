package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wisteriahuman/prometheus/internal/service"
)

type ExportHandler struct {
	vault *service.Vault
	md    *service.Markdown
}

func NewExportHandler(vault *service.Vault, md *service.Markdown) *ExportHandler {
	return &ExportHandler{vault: vault, md: md}
}

// GET /api/export/{path...}?format=html|md|pdf
func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) {
	notePath := chi.URLParam(r, "*")
	if notePath == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "html"
	}

	note, err := h.vault.ReadNote(notePath)
	if err != nil {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}

	baseName := strings.TrimSuffix(filepath.Base(notePath), ".md")
	inline := r.URL.Query().Get("inline") == "true"

	switch format {
	case "html":
		h.exportHTML(w, note, baseName, inline)
	case "md":
		h.exportMarkdown(w, note, baseName)
	default:
		writeError(w, http.StatusBadRequest, "unsupported format: "+format+". Use html or md")
	}
}

func (h *ExportHandler) exportHTML(w http.ResponseWriter, note *service.VaultNote, baseName string, inline bool) {
	bodyHTML := h.md.ToHTML(note.RawContent)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Noto Sans JP', sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 40px 20px;
    line-height: 1.7;
    color: #1a1a1a;
    background: #fff;
  }
  h1 { font-size: 2em; border-bottom: 1px solid #ddd; padding-bottom: 0.3em; }
  h2 { font-size: 1.5em; margin-top: 1.5em; }
  h3 { font-size: 1.2em; margin-top: 1.2em; }
  a { color: #2563eb; }
  code { background: #f3f4f6; padding: 0.15em 0.4em; border-radius: 3px; font-size: 0.9em; }
  pre { background: #f3f4f6; padding: 1em; border-radius: 6px; overflow-x: auto; }
  pre code { background: none; padding: 0; }
  blockquote { border-left: 3px solid #6366f1; padding-left: 1em; color: #666; }
  table { border-collapse: collapse; width: 100%%; margin: 1em 0; }
  th, td { border: 1px solid #ddd; padding: 0.5em 0.75em; text-align: left; }
  th { background: #f9fafb; }
  ul, ol { padding-left: 1.5em; }
  li input[type="checkbox"] { margin-right: 0.5em; }
  img { max-width: 100%%; }
  .wikilink { color: #6366f1; text-decoration: underline dotted; }
  @media print {
    body { padding: 0; }
    a { color: inherit; text-decoration: none; }
  }
</style>
</head>
<body>
%s
<footer style="margin-top: 3em; padding-top: 1em; border-top: 1px solid #eee; font-size: 0.8em; color: #999;">
  Exported from Prometheus
</footer>
</body>
</html>`, note.Frontmatter.Title, bodyHTML)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !inline {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.html"`, baseName))
	}
	w.Write([]byte(html))
}

func (h *ExportHandler) exportMarkdown(w http.ResponseWriter, note *service.VaultNote, baseName string) {
	// Export clean markdown without frontmatter
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.md"`, baseName))
	w.Write([]byte(note.Content))
}
