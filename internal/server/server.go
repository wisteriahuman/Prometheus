package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wisteriahuman/prometheus/internal/config"
	"github.com/wisteriahuman/prometheus/internal/db"
	"github.com/wisteriahuman/prometheus/internal/handler"
	"github.com/wisteriahuman/prometheus/internal/service"
)

type Server struct {
	cfg    *config.Config
	db     *db.DB
	router http.Handler
}

func NewServer(cfg *config.Config) *Server {
	// Database
	database := db.New(cfg.DBPath)
	database.Migrate()

	// Services
	vault := service.NewVault(cfg.VaultPath)
	md := service.NewMarkdown()
	indexer := service.NewIndexer(database, vault, md)
	daily := service.NewDaily(vault)

	// Init vault with sample notes if empty
	if service.InitVault(vault) {
		log.Println("Initialized vault with sample notes")
	}

	// Index all notes on startup
	count := indexer.IndexAll()
	log.Printf("Indexed %d notes", count)

	// Handlers
	handlers := &Handlers{
		Notes:     handler.NewNotesHandler(vault, indexer),
		Search:    handler.NewSearchHandler(database),
		Graph:     handler.NewGraphHandler(database),
		Tasks:     handler.NewTasksHandler(database, vault, indexer),
		Tags:      handler.NewTagsHandler(database),
		Daily:     handler.NewDailyHandler(daily, indexer, md, vault),
		Preview:   handler.NewPreviewHandler(md),
		Tree:      handler.NewTreeHandler(vault),
		Folders:   handler.NewFoldersHandler(cfg.VaultPath),
		Backlinks: handler.NewBacklinksHandler(database),
		Config:    handler.NewConfigHandler(cfg),
		Export:    handler.NewExportHandler(vault, md),
		Assets:    handler.NewAssetsHandler(cfg.VaultPath),
	}

	router := NewRouter(handlers)

	return &Server{
		cfg:    cfg,
		db:     database,
		router: router,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	log.Printf("Prometheus server starting on http://localhost:%s", s.cfg.Port)
	log.Printf("Vault: %s", s.cfg.VaultPath)
	return http.ListenAndServe(addr, s.router)
}
