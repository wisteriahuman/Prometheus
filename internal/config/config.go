package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	VaultPath string
	DBPath    string
	Port      string
}

func expandHome(p string) string {
	if len(p) >= 2 && p[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, p[2:])
	}
	if p == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	return p
}

func New(vaultPath, dbPath, port string) *Config {
	if vaultPath == "" {
		vaultPath = getEnv("PROMETHEUS_VAULT_PATH", "./vault")
	}
	if dbPath == "" {
		dbPath = getEnv("PROMETHEUS_DB_PATH", "./data/prometheus.db")
	}
	if port == "" {
		port = getEnv("PORT", "3000")
	}

	vaultPath = expandHome(vaultPath)
	dbPath = expandHome(dbPath)

	abs, err := filepath.Abs(vaultPath)
	if err == nil {
		vaultPath = abs
	}
	abs, err = filepath.Abs(dbPath)
	if err == nil {
		dbPath = abs
	}

	return &Config{
		VaultPath: vaultPath,
		DBPath:    dbPath,
		Port:      port,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
