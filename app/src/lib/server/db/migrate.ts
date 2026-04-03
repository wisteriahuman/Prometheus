import Database from "better-sqlite3";
import { resolve, dirname } from "node:path";
import { mkdirSync } from "node:fs";

const DB_PATH = resolve(
  process.env.PROMETHEUS_DB_PATH ?? "./data/prometheus.db",
);

try {
  mkdirSync(dirname(DB_PATH), { recursive: true });
} catch {
  // directory already exists
}

const sqlite = new Database(DB_PATH);

sqlite.pragma("journal_mode = WAL");
sqlite.pragma("foreign_keys = ON");

// Create tables
sqlite.exec(`
  CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    modified_at INTEGER NOT NULL,
    theme TEXT,
    checksum TEXT NOT NULL
  )
`);

sqlite.exec(`
  CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    target_id TEXT REFERENCES notes(id) ON DELETE SET NULL,
    target_slug TEXT NOT NULL,
    context TEXT
  )
`);

sqlite.exec(`
  CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
  )
`);

sqlite.exec(`
  CREATE TABLE IF NOT EXISTS note_tags (
    note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (note_id, tag_id)
  )
`);

sqlite.exec(`
  CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    completed INTEGER NOT NULL DEFAULT 0,
    line_number INTEGER NOT NULL,
    due_date TEXT
  )
`);

// Auth tables
sqlite.exec(`
  CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at INTEGER NOT NULL
  )
`);

sqlite.exec(`
  CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at INTEGER NOT NULL
  )
`);

// Workspace tables
sqlite.exec(`
  CREATE TABLE IF NOT EXISTS workspaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    vault_path TEXT NOT NULL,
    engine TEXT NOT NULL DEFAULT 'sqlite',
    db_url TEXT,
    theme TEXT,
    created_at INTEGER NOT NULL,
    owner_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
  )
`);

sqlite.exec(`
  CREATE TABLE IF NOT EXISTS workspace_members (
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'editor',
    PRIMARY KEY (workspace_id, user_id)
  )
`);

// FTS5 virtual table for full-text search
sqlite.exec(`
  CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
    title,
    content,
    content_rowid='rowid'
  )
`);

// Indexes
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_links_source ON links(source_id)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_links_target ON links(target_id)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_links_slug ON links(target_slug)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_tasks_note ON tasks(note_id)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_note_tags_note ON note_tags(note_id)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_note_tags_tag ON note_tags(tag_id)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id)`);
sqlite.exec(`CREATE INDEX IF NOT EXISTS idx_workspaces_owner ON workspaces(owner_id)`);

console.log("Database migrated successfully:", DB_PATH);
sqlite.close();
