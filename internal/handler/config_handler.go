package handler

import (
	"net/http"

	"github.com/wisteriahuman/prometheus/internal/config"
)

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"vaultPath": h.cfg.VaultPath,
	})
}
