package db

func (d *DB) Migrate() {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			path TEXT NOT NULL UNIQUE,
			title TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			created_at INTEGER NOT NULL,
			modified_at INTEGER NOT NULL,
			theme TEXT,
			checksum TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
			target_id TEXT REFERENCES notes(id) ON DELETE SET NULL,
			target_slug TEXT NOT NULL,
			context TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS note_tags (
			note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
			tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (note_id, tag_id)
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			completed INTEGER NOT NULL DEFAULT 0,
			line_number INTEGER NOT NULL,
			due_date TEXT
		)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(title, content, tokenize='unicode61')`,
		`CREATE INDEX IF NOT EXISTS idx_links_source ON links(source_id)`,
		`CREATE INDEX IF NOT EXISTS idx_links_target ON links(target_id)`,
		`CREATE INDEX IF NOT EXISTS idx_links_slug ON links(target_slug)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_note ON tasks(note_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed)`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_note ON note_tags(note_id)`,
		`CREATE INDEX IF NOT EXISTS idx_note_tags_tag ON note_tags(tag_id)`,
	}

	for _, stmt := range stmts {
		_, err := d.Exec(stmt)
		if err != nil {
			// FTS5 might already exist with different schema, ignore
			continue
		}
	}
}
