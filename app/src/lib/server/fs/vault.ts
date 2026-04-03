import { readdir, readFile, writeFile, unlink, mkdir, stat } from "node:fs/promises";
import { resolve, relative, join, dirname, extname, basename } from "node:path";
import { createHash } from "node:crypto";
import matter from "gray-matter";
import { ulid } from "ulid";
import type { NoteFrontmatter } from "@prometheus/core";
import { VAULT_PATH } from "../env.js";

export interface VaultNote {
  path: string;
  title: string;
  content: string;
  rawContent: string;
  frontmatter: NoteFrontmatter;
  checksum: string;
}

export interface VaultFileEntry {
  name: string;
  path: string;
  isDirectory: boolean;
  children?: VaultFileEntry[];
}

function computeChecksum(content: string): string {
  return createHash("sha256").update(content).digest("hex");
}

function ensureFrontmatter(
  rawContent: string,
  filePath: string,
): { content: string; frontmatter: NoteFrontmatter; body: string } {
  const parsed = matter(rawContent);
  const now = new Date().toISOString();

  const frontmatter: NoteFrontmatter = {
    id: parsed.data.id ?? ulid(),
    title:
      parsed.data.title ??
      basename(filePath, extname(filePath)),
    created: parsed.data.created ?? now,
    modified: parsed.data.modified ?? now,
    tags: parsed.data.tags ?? [],
    theme: parsed.data.theme ?? undefined,
  };

  return {
    content: parsed.content,
    frontmatter,
    body: parsed.content,
  };
}

function serializeNote(frontmatter: NoteFrontmatter, body: string): string {
  // Remove undefined/null values — gray-matter's YAML dumper cannot handle them
  const clean: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(frontmatter)) {
    if (value !== undefined && value !== null) {
      clean[key] = value;
    }
  }
  return matter.stringify(body, clean);
}

export async function getVaultPath(): Promise<string> {
  await mkdir(VAULT_PATH, { recursive: true });
  return VAULT_PATH;
}

export async function readNote(notePath: string): Promise<VaultNote> {
  const vaultPath = await getVaultPath();
  const fullPath = resolve(vaultPath, notePath);
  const rawContent = await readFile(fullPath, "utf-8");
  const checksum = computeChecksum(rawContent);
  const { content, frontmatter } = ensureFrontmatter(rawContent, notePath);

  return {
    path: notePath,
    title: frontmatter.title,
    content,
    rawContent,
    frontmatter,
    checksum,
  };
}

export async function writeNote(
  notePath: string,
  body: string,
  frontmatter: NoteFrontmatter,
): Promise<VaultNote> {
  const vaultPath = await getVaultPath();
  const fullPath = resolve(vaultPath, notePath);

  frontmatter.modified = new Date().toISOString();

  const rawContent = serializeNote(frontmatter, body);
  await mkdir(dirname(fullPath), { recursive: true });
  await writeFile(fullPath, rawContent, "utf-8");

  const checksum = computeChecksum(rawContent);

  return {
    path: notePath,
    title: frontmatter.title,
    content: body,
    rawContent,
    frontmatter,
    checksum,
  };
}

export async function createNote(
  notePath: string,
  title?: string,
): Promise<VaultNote> {
  const noteTitle = title ?? basename(notePath, extname(notePath));
  const now = new Date().toISOString();

  const frontmatter: NoteFrontmatter = {
    id: ulid(),
    title: noteTitle,
    created: now,
    modified: now,
    tags: [],
  };

  const body = `\n# ${noteTitle}\n\n`;

  return writeNote(notePath, body, frontmatter);
}

export async function deleteNote(notePath: string): Promise<void> {
  const vaultPath = await getVaultPath();
  const fullPath = resolve(vaultPath, notePath);
  await unlink(fullPath);
}

export async function listNotes(dirPath: string = ""): Promise<string[]> {
  const vaultPath = await getVaultPath();
  const targetPath = resolve(vaultPath, dirPath);
  const notes: string[] = [];

  async function walk(dir: string) {
    let entries;
    try {
      entries = await readdir(dir, { withFileTypes: true });
    } catch {
      return;
    }

    for (const entry of entries) {
      if (entry.name.startsWith(".")) continue;

      const fullPath = join(dir, entry.name);
      if (entry.isDirectory()) {
        await walk(fullPath);
      } else if (extname(entry.name) === ".md") {
        notes.push(relative(vaultPath, fullPath));
      }
    }
  }

  await walk(targetPath);
  return notes.sort();
}

export async function getFileTree(
  dirPath: string = "",
): Promise<VaultFileEntry[]> {
  const vaultPath = await getVaultPath();
  const targetPath = resolve(vaultPath, dirPath);

  let entries;
  try {
    entries = await readdir(targetPath, { withFileTypes: true });
  } catch {
    return [];
  }

  const result: VaultFileEntry[] = [];

  for (const entry of entries) {
    if (entry.name.startsWith(".")) continue;

    const entryRelPath = dirPath ? join(dirPath, entry.name) : entry.name;

    if (entry.isDirectory()) {
      const children = await getFileTree(entryRelPath);
      result.push({
        name: entry.name,
        path: entryRelPath,
        isDirectory: true,
        children,
      });
    } else if (extname(entry.name) === ".md") {
      result.push({
        name: entry.name,
        path: entryRelPath,
        isDirectory: false,
      });
    }
  }

  return result.sort((a, b) => {
    if (a.isDirectory && !b.isDirectory) return -1;
    if (!a.isDirectory && b.isDirectory) return 1;
    return a.name.localeCompare(b.name);
  });
}

export async function noteExists(notePath: string): Promise<boolean> {
  const vaultPath = await getVaultPath();
  const fullPath = resolve(vaultPath, notePath);
  try {
    await stat(fullPath);
    return true;
  } catch {
    return false;
  }
}
