import { json, error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { VAULT_PATH } from "$lib/server/env.js";
import { mkdir, rm, readdir } from "node:fs/promises";
import { resolve, join } from "node:path";

// POST /api/folders/:path — create a folder
export const POST: RequestHandler = async ({ params }) => {
  const folderPath = params.path;
  if (!folderPath) throw error(400, "path is required");

  const fullPath = resolve(VAULT_PATH, folderPath);

  // Security: prevent path traversal
  if (!fullPath.startsWith(resolve(VAULT_PATH))) {
    throw error(400, "Invalid path");
  }

  await mkdir(fullPath, { recursive: true });
  return json({ success: true }, { status: 201 });
};

// DELETE /api/folders/:path — delete a folder (must be empty or force)
export const DELETE: RequestHandler = async ({ params, url }) => {
  const folderPath = params.path;
  if (!folderPath) throw error(400, "path is required");

  const fullPath = resolve(VAULT_PATH, folderPath);

  if (!fullPath.startsWith(resolve(VAULT_PATH))) {
    throw error(400, "Invalid path");
  }

  const force = url.searchParams.get("force") === "true";

  try {
    const entries = await readdir(fullPath);
    const nonHidden = entries.filter((e) => !e.startsWith("."));

    if (nonHidden.length > 0 && !force) {
      throw error(400, "フォルダが空ではありません。中のノートも含めて削除しますか？");
    }

    await rm(fullPath, { recursive: true });
    return json({ success: true });
  } catch (e: any) {
    if (e.status) throw e;
    throw error(500, "Failed to delete folder");
  }
};
