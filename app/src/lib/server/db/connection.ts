import Database from "better-sqlite3";
import { drizzle } from "drizzle-orm/better-sqlite3";
import * as schema from "./schema.js";
import { resolve, dirname } from "node:path";
import { mkdirSync } from "node:fs";
import type { DbEngine } from "@prometheus/core";

type DrizzleDb = ReturnType<typeof drizzle>;

const connections = new Map<string, DrizzleDb>();

/**
 * Get or create a database connection for a workspace.
 * For SQLite workspaces, creates a local .db file.
 * For PostgreSQL workspaces, connects to the provided URL.
 * (PostgreSQL support will be added in a future phase)
 */
export function getWorkspaceDb(
  workspaceId: string,
  engine: DbEngine,
  dbUrl?: string | null,
): DrizzleDb {
  const cacheKey = `${engine}:${workspaceId}`;

  if (connections.has(cacheKey)) {
    return connections.get(cacheKey)!;
  }

  let db: DrizzleDb;

  if (engine === "sqlite") {
    const dbPath = resolve(`./data/workspaces/${workspaceId}/prometheus.db`);
    try {
      mkdirSync(dirname(dbPath), { recursive: true });
    } catch {
      // exists
    }

    const sqlite = new Database(dbPath);
    sqlite.pragma("journal_mode = WAL");
    sqlite.pragma("foreign_keys = ON");
    db = drizzle(sqlite, { schema });
  } else {
    // PostgreSQL — placeholder for Phase 5 expansion
    // When implemented:
    // import { drizzle } from "drizzle-orm/postgres-js";
    // import postgres from "postgres";
    // const client = postgres(dbUrl!);
    // db = drizzle(client, { schema: pgSchema });
    throw new Error(
      "PostgreSQL support is not yet implemented. Use engine: sqlite for now.",
    );
  }

  connections.set(cacheKey, db);
  return db;
}

/**
 * Close and remove a workspace connection from the cache.
 */
export function closeWorkspaceDb(workspaceId: string): void {
  for (const [key, _] of connections) {
    if (key.includes(workspaceId)) {
      connections.delete(key);
    }
  }
}
