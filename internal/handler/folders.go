package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type FoldersHandler struct {
	vaultPath string
}

func NewFoldersHandler(vaultPath string) *FoldersHandler {
	return &FoldersHandler{vaultPath: vaultPath}
}

func (h *FoldersHandler) Create(w http.ResponseWriter, r *http.Request) {
	folderPath := chi.URLParam(r, "*")
	if folderPath == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	// Path traversal protection
	if strings.Contains(folderPath, "..") {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	fullPath := filepath.Join(h.vaultPath, folderPath)
	// Ensure it's still within vault
	if !strings.HasPrefix(fullPath, h.vaultPath) {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	if err := os.MkdirAll(fullPath, 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create folder: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]bool{"success": true})
}

func (h *FoldersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	folderPath := chi.URLParam(r, "*")
	if folderPath == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	if strings.Contains(folderPath, "..") {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	fullPath := filepath.Join(h.vaultPath, folderPath)
	if !strings.HasPrefix(fullPath, h.vaultPath) {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	force := r.URL.Query().Get("force") == "true"

	var err error
	if force {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete folder: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}
