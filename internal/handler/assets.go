package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type AssetsHandler struct {
	vaultPath string
}

func NewAssetsHandler(vaultPath string) *AssetsHandler {
	return &AssetsHandler{vaultPath: vaultPath}
}

func (h *AssetsHandler) assetsDir() string {
	return filepath.Join(h.vaultPath, "assets")
}

// POST /api/assets — upload a file
func (h *AssetsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Max 50MB
	r.ParseMultipartForm(50 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	// Generate unique filename to avoid conflicts
	ext := filepath.Ext(header.Filename)
	baseName := strings.TrimSuffix(header.Filename, ext)
	// Sanitize filename
	baseName = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '-'
		}
		return r
	}, baseName)

	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s%s", baseName, timestamp, ext)

	// Create assets directory
	dir := h.assetsDir()
	os.MkdirAll(dir, 0o755)

	// Write file
	destPath := filepath.Join(dir, filename)
	dest, err := os.Create(destPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create file")
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to write file")
		return
	}

	// Return markdown reference
	contentType := header.Header.Get("Content-Type")
	isImage := strings.HasPrefix(contentType, "image/")

	var markdown string
	if isImage {
		markdown = fmt.Sprintf("![%s](/api/assets/%s)", baseName, filename)
	} else {
		markdown = fmt.Sprintf("[%s](/api/assets/%s)", header.Filename, filename)
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"filename": filename,
		"path":     "/api/assets/" + filename,
		"markdown": markdown,
	})
}

// GET /api/assets/{filename} — serve a file
func (h *AssetsHandler) Serve(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "*")
	if filename == "" {
		writeError(w, http.StatusBadRequest, "filename is required")
		return
	}

	// Security: prevent path traversal
	cleanName := filepath.Clean(filename)
	if strings.Contains(cleanName, "..") {
		writeError(w, http.StatusBadRequest, "invalid filename")
		return
	}

	filePath := filepath.Join(h.assetsDir(), cleanName)

	// Check file exists
	info, err := os.Stat(filePath)
	if err != nil {
		writeError(w, http.StatusNotFound, "file not found")
		return
	}

	// Detect content type
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := "application/octet-stream"
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".svg":
		contentType = "image/svg+xml"
	case ".webp":
		contentType = "image/webp"
	case ".pdf":
		contentType = "application/pdf"
	case ".mp4":
		contentType = "video/mp4"
	case ".webm":
		contentType = "video/webm"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeFile(w, r, filePath)
	_ = info // used for potential future size checks
}

// GET /api/assets — list all assets
func (h *AssetsHandler) List(w http.ResponseWriter, r *http.Request) {
	dir := h.assetsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}

	type assetInfo struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Size int64  `json:"size"`
	}

	var assets []assetInfo
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		assets = append(assets, assetInfo{
			Name: entry.Name(),
			Path: "/api/assets/" + entry.Name(),
			Size: info.Size(),
		})
	}

	if assets == nil {
		assets = []assetInfo{}
	}
	writeJSON(w, http.StatusOK, assets)
}

// DELETE /api/assets/{filename} — delete a file
func (h *AssetsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "*")
	if filename == "" || strings.Contains(filepath.Clean(filename), "..") {
		writeError(w, http.StatusBadRequest, "invalid filename")
		return
	}

	filePath := filepath.Join(h.assetsDir(), filepath.Clean(filename))
	if err := os.Remove(filePath); err != nil {
		writeError(w, http.StatusNotFound, "file not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
