package service

import (
	"fmt"
	"log"

	"github.com/wisteriahuman/prometheus/internal/db"
)

type Indexer struct {
	db    *db.DB
	vault *Vault
	md    *Markdown
}

func NewIndexer(database *db.DB, vault *Vault, md *Markdown) *Indexer {
	return &Indexer{db: database, vault: vault, md: md}
}

func (idx *Indexer) IndexNote(notePath string) error {
	note, err := idx.vault.ReadNote(notePath)
	if err != nil {
		return fmt.Errorf("read note %s: %w", notePath, err)
	}

	// Check if already indexed with same checksum
	var existingChecksum string
	row := idx.db.QueryRow("SELECT checksum FROM notes WHERE id = ?", note.Frontmatter.ID)
	row.Scan(&existingChecksum)
	if existingChecksum == note.Checksum {
		return nil
	}

	wikilinks := ExtractWikilinks(note.Content)
	tasks := ExtractTasks(note.Content)

	tx, err := idx.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing entry for this path (handles ID changes)
	tx.Exec("DELETE FROM notes WHERE path = ?", note.Path)

	// Upsert note
	_, err = tx.Exec(`
		INSERT INTO notes (id, path, title, content, created_at, modified_at, theme, checksum)
		VALUES (?, ?, ?, ?, strftime('%s', ?), strftime('%s', ?), ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			path = excluded.path,
			title = excluded.title,
			content = excluded.content,
			modified_at = excluded.modified_at,
			theme = excluded.theme,
			checksum = excluded.checksum
	`, note.Frontmatter.ID, note.Path, note.Frontmatter.Title, note.Content,
		note.Frontmatter.Created, note.Frontmatter.Modified,
		themePtr(note.Frontmatter.Theme), note.Checksum)
	if err != nil {
		return fmt.Errorf("upsert note: %w", err)
	}

	// Update FTS5 index
	tx.Exec("DELETE FROM notes_fts WHERE rowid IN (SELECT rowid FROM notes_fts WHERE title = ? OR content = ?)",
		note.Frontmatter.Title, note.Content)
	tx.Exec("INSERT INTO notes_fts (rowid, title, content) VALUES ((SELECT rowid FROM notes WHERE id = ?), ?, ?)",
		note.Frontmatter.ID, note.Frontmatter.Title, note.Content)

	// Re-index links
	tx.Exec("DELETE FROM links WHERE source_id = ?", note.Frontmatter.ID)
	for _, wl := range wikilinks {
		ctx := wl.DisplayText
		if ctx == "" {
			ctx = wl.Slug
		}
		tx.Exec("INSERT INTO links (source_id, target_id, target_slug, context) VALUES (?, NULL, ?, ?)",
			note.Frontmatter.ID, wl.Slug, ctx)
	}

	// Re-index tasks
	tx.Exec("DELETE FROM tasks WHERE note_id = ?", note.Frontmatter.ID)
	for _, t := range tasks {
		completed := 0
		if t.Completed {
			completed = 1
		}
		tx.Exec("INSERT INTO tasks (note_id, content, completed, line_number) VALUES (?, ?, ?, ?)",
			note.Frontmatter.ID, t.Content, completed, t.Line)
	}

	// Re-index tags
	tx.Exec("DELETE FROM note_tags WHERE note_id = ?", note.Frontmatter.ID)
	for _, tagName := range note.Frontmatter.Tags {
		tx.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", tagName)

		var tagID int64
		row := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName)
		if err := row.Scan(&tagID); err != nil {
			continue
		}
		tx.Exec("INSERT OR IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)",
			note.Frontmatter.ID, tagID)
	}

	return tx.Commit()
}

func (idx *Indexer) IndexAll() int {
	notes, err := idx.vault.ListNotes()
	if err != nil {
		log.Printf("Failed to list notes: %v", err)
		return 0
	}

	vaultPaths := make(map[string]bool)
	for _, n := range notes {
		vaultPaths[n] = true
	}

	count := 0
	for _, notePath := range notes {
		if err := idx.IndexNote(notePath); err != nil {
			log.Printf("Failed to index %s: %v", notePath, err)
			continue
		}
		count++
	}

	// Remove stale entries
	rows, err := idx.db.Query("SELECT id, path FROM notes")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, path string
			rows.Scan(&id, &path)
			if !vaultPaths[path] {
				idx.db.Exec("DELETE FROM notes WHERE id = ?", id)
			}
		}
	}

	return count
}

func (idx *Indexer) RemoveNote(notePath string) {
	idx.db.Exec("DELETE FROM notes WHERE path = ?", notePath)
}

func themePtr(t *string) interface{} {
	if t == nil {
		return nil
	}
	return *t
}
