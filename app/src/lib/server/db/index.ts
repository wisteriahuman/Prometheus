import Database from "better-sqlite3";
import { drizzle } from "drizzle-orm/better-sqlite3";
import * as schema from "./schema.js";
import { dirname } from "node:path";
import { mkdirSync } from "node:fs";
import { DB_PATH } from "../env.js";

try {
  mkdirSync(dirname(DB_PATH), { recursive: true });
} catch {
  // directory already exists
}

const sqlite = new Database(DB_PATH);

sqlite.pragma("journal_mode = WAL");
sqlite.pragma("foreign_keys = ON");

export const db = drizzle(sqlite, { schema });
export { schema };
