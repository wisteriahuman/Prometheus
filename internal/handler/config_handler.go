package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wisteriahuman/prometheus/internal/config"
)

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

type vaultConfig struct {
	Name              string            `json:"name,omitempty"`
	Theme             string            `json:"theme,omitempty"`
	DailyNoteTemplate string            `json:"dailyNoteTemplate,omitempty"`
	CustomThemes      []json.RawMessage `json:"customThemes,omitempty"`
	TutorialShown     bool              `json:"tutorialShown,omitempty"`
}

func (h *ConfigHandler) configPath() string {
	return filepath.Join(h.cfg.VaultPath, ".prometheus", "config.json")
}

func (h *ConfigHandler) readVaultConfig() vaultConfig {
	var vc vaultConfig
	data, err := os.ReadFile(h.configPath())
	if err != nil {
		return vc
	}
	json.Unmarshal(data, &vc)
	return vc
}

func (h *ConfigHandler) writeVaultConfig(vc vaultConfig) error {
	dir := filepath.Dir(h.configPath())
	os.MkdirAll(dir, 0o755)
	data, err := json.MarshalIndent(vc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.configPath(), data, 0o644)
}

// GET /api/config
func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	vc := h.readVaultConfig()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"vaultPath":     h.cfg.VaultPath,
		"theme":         vc.Theme,
		"customThemes":  vc.CustomThemes,
		"tutorialShown": vc.TutorialShown,
	})
}

// PUT /api/config
func (h *ConfigHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Theme         *string           `json:"theme"`
		CustomThemes  []json.RawMessage `json:"customThemes"`
		TutorialShown *bool             `json:"tutorialShown"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	vc := h.readVaultConfig()

	if req.Theme != nil {
		vc.Theme = *req.Theme
	}
	if req.CustomThemes != nil {
		vc.CustomThemes = req.CustomThemes
	}
	if req.TutorialShown != nil {
		vc.TutorialShown = *req.TutorialShown
	}

	if err := h.writeVaultConfig(vc); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save config")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"theme":        vc.Theme,
		"customThemes": vc.CustomThemes,
	})
}
