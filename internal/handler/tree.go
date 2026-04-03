package handler

import (
	"net/http"

	"github.com/wisteriahuman/prometheus/internal/service"
)

type TreeHandler struct {
	vault *service.Vault
}

func NewTreeHandler(vault *service.Vault) *TreeHandler {
	return &TreeHandler{vault: vault}
}

func (h *TreeHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	tree, err := h.vault.GetFileTree("")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get file tree")
		return
	}

	if tree == nil {
		tree = []*service.FileEntry{}
	}

	writeJSON(w, http.StatusOK, tree)
}
