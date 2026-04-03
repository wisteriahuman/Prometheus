import { createNote, noteExists, readNote, getVaultPath } from "./vault.js";
import { readFile } from "node:fs/promises";
import { resolve } from "node:path";

function formatDate(date: Date): string {
  return date.toISOString().slice(0, 10);
}

function formatDateJa(date: Date): string {
  return date.toLocaleDateString("ja-JP", {
    year: "numeric",
    month: "long",
    day: "numeric",
    weekday: "long",
  });
}

async function getTemplate(): Promise<string> {
  const vaultPath = await getVaultPath();
  const configPath = resolve(vaultPath, ".prometheus/config.json");

  try {
    const config = JSON.parse(await readFile(configPath, "utf-8"));
    if (config.dailyNoteTemplate) {
      return config.dailyNoteTemplate;
    }
  } catch {
    // Use default template
  }

  return `---
title: "{{date}}"
tags: [daily]
---

# {{dateJa}}

## Tasks

- [ ]

## Notes

`;
}

function applyTemplate(template: string, date: Date): string {
  return template
    .replace(/\{\{date\}\}/g, formatDate(date))
    .replace(/\{\{dateJa\}\}/g, formatDateJa(date));
}

export async function getDailyNotePath(date: Date): Promise<string> {
  return `daily/${formatDate(date)}.md`;
}

export async function ensureDailyNote(date: Date) {
  const path = await getDailyNotePath(date);

  if (await noteExists(path)) {
    return readNote(path);
  }

  const template = await getTemplate();
  const content = applyTemplate(template, date);

  // Parse out frontmatter and body from template
  const fmMatch = content.match(/^---\n([\s\S]*?)\n---\n([\s\S]*)$/);

  if (fmMatch) {
    // Template has frontmatter - create note with the full raw content
    const { writeFile, mkdir } = await import("node:fs/promises");
    const { resolve, dirname } = await import("node:path");
    const { ulid } = await import("ulid");
    const vaultPath = await getVaultPath();
    const fullPath = resolve(vaultPath, path);

    // Inject id and dates into frontmatter
    const now = new Date().toISOString();
    let fm = fmMatch[1];
    if (!fm.includes("id:")) fm = `id: "${ulid()}"\n${fm}`;
    if (!fm.includes("created:")) fm = `${fm}\ncreated: "${now}"`;
    if (!fm.includes("modified:")) fm = `${fm}\nmodified: "${now}"`;

    const rawContent = `---\n${fm}\n---\n${fmMatch[2]}`;

    await mkdir(dirname(fullPath), { recursive: true });
    await writeFile(fullPath, rawContent, "utf-8");

    return readNote(path);
  }

  // Fallback: simple creation
  return createNote(path, formatDate(date));
}

export function getRecentDailyDates(count: number = 7): Date[] {
  const dates: Date[] = [];
  const today = new Date();

  for (let i = 0; i < count; i++) {
    const date = new Date(today);
    date.setDate(today.getDate() - i);
    dates.push(date);
  }

  return dates;
}
