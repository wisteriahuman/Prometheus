import { config } from "dotenv";
import { resolve, join } from "node:path";
import { homedir } from "node:os";

// Load .env file
config();

function expandHome(p: string): string {
  if (p.startsWith("~/") || p === "~") {
    return join(homedir(), p.slice(1));
  }
  return p;
}

export const VAULT_PATH = resolve(
  expandHome(process.env.PROMETHEUS_VAULT_PATH ?? "./vault"),
);

export const DB_PATH = resolve(
  expandHome(process.env.PROMETHEUS_DB_PATH ?? "./data/prometheus.db"),
);

console.log(`Vault: ${VAULT_PATH}`);
console.log(`DB: ${DB_PATH}`);
