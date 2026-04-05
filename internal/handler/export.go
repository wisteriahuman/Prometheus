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

// GET /api/export/{path...}?format=html|md&theme=ocean&inline=true
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

	// Resolve theme: query param > frontmatter > light (default for export)
	themeSlug := r.URL.Query().Get("theme")
	if themeSlug == "" && note.Frontmatter.Theme != nil && *note.Frontmatter.Theme != "" {
		themeSlug = *note.Frontmatter.Theme
	}
	if themeSlug == "" {
		themeSlug = "light"
	}

	switch format {
	case "html":
		h.exportHTML(w, note, baseName, themeSlug, inline)
	case "md":
		h.exportMarkdown(w, note, baseName)
	default:
		writeError(w, http.StatusBadRequest, "unsupported format: "+format+". Use html or md")
	}
}

func (h *ExportHandler) exportHTML(w http.ResponseWriter, note *service.VaultNote, baseName, themeSlug string, inline bool) {
	bodyHTML := h.md.ToHTML(note.RawContent)
	tc := service.GetThemeColors(themeSlug)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
  :root {
    --bg: %s;
    --bg-card: %s;
    --text: %s;
    --text-muted: %s;
    --text-dim: %s;
    --primary: %s;
    --primary-light: %s;
    --accent: %s;
    --border: %s;
    --success: %s;
    --warning: %s;
    --error: %s;
  }

  @page {
    margin: 0;
  }

  * {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
  }

  html, body {
    margin: 0;
    background: var(--bg);
  }

  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Noto Sans JP', sans-serif;
    max-width: 800px;
    margin: 0 auto;
    padding: 15mm 20mm;
    line-height: 1.7;
    color: var(--text);
  }

  h1 {
    font-size: 2em;
    color: var(--primary-light);
    border-bottom: 1px solid var(--border);
    padding-bottom: 0.3em;
    margin-top: 0;
  }
  h2 {
    font-size: 1.5em;
    color: var(--primary-light);
    margin-top: 1.5em;
  }
  h3 { font-size: 1.2em; margin-top: 1.2em; }

  a { color: var(--accent); }
  a.wikilink { color: var(--primary-light); text-decoration: underline dotted; }

  code {
    background: var(--bg-card);
    padding: 0.15em 0.4em;
    border-radius: 3px;
    font-size: 0.9em;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    color: var(--accent);
  }
  pre {
    background: var(--bg-card);
    border: 1px solid var(--border);
    padding: 1em;
    border-radius: 6px;
    overflow-x: auto;
  }
  pre code {
    background: none;
    padding: 0;
    color: var(--text);
  }

  blockquote {
    border-left: 3px solid var(--primary);
    padding-left: 1em;
    color: var(--text-muted);
    font-style: italic;
    margin: 1em 0;
  }

  table { border-collapse: collapse; width: 100%%; margin: 1em 0; }
  th, td { border: 1px solid var(--border); padding: 0.5em 0.75em; text-align: left; }
  th { background: var(--bg-card); font-weight: 600; }

  ul, ol { padding-left: 1.5em; }
  li { margin: 0.25em 0; }
  li input[type="checkbox"] { margin-right: 0.5em; accent-color: var(--primary); }

  hr { border: none; border-top: 1px solid var(--border); margin: 2em 0; }
  img { max-width: 100%%; border-radius: 6px; }
  strong { color: var(--text); }
  em { color: var(--text-muted); }

  @media print {
    body { padding: 0; max-width: none; }
    a { text-decoration: none; }
    pre { white-space: pre-wrap; word-wrap: break-word; }
  }
</style>
</head>
<body>
%s
<footer style="margin-top: 3em; padding-top: 1em; border-top: 1px solid var(--border); font-size: 0.8em; color: var(--text-dim);">
  Exported from Prometheus
</footer>
</body>
</html>`,
		note.Frontmatter.Title,
		tc.BgDark, tc.BgCard,
		tc.TextMain, tc.TextMuted, tc.TextDim,
		tc.Primary, tc.PrimaryLight, tc.Accent,
		tc.Border, tc.Success, tc.Warning, tc.Error,
		bodyHTML,
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !inline {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.html"`, baseName))
	}
	w.Write([]byte(html))
}

func (h *ExportHandler) exportMarkdown(w http.ResponseWriter, note *service.VaultNote, baseName string) {
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.md"`, baseName))
	w.Write([]byte(note.Content))
}
