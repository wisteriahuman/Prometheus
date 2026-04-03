import { json, error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { readNote, writeNote, deleteNote, noteExists } from "$lib/server/fs/vault.js";
import { indexNote, removeNoteFromIndex } from "$lib/server/indexer.js";

// GET /api/notes/:path — read a note
export const GET: RequestHandler = async ({ params }) => {
  const notePath = params.path;
  if (!notePath) throw error(400, "path is required");

  if (!(await noteExists(notePath))) {
    throw error(404, "Note not found");
  }

  const note = await readNote(notePath);

  return json({
    path: note.path,
    title: note.frontmatter.title,
    content: note.content,
    frontmatter: note.frontmatter,
    checksum: note.checksum,
  });
};

// PUT /api/notes/:path — update a note
export const PUT: RequestHandler = async ({ params, request }) => {
  const notePath = params.path;
  if (!notePath) throw error(400, "path is required");

  const { content, frontmatter } = await request.json();

  if (typeof content !== "string") {
    throw error(400, "content is required");
  }

  let note;
  if (await noteExists(notePath)) {
    const existing = await readNote(notePath);
    const updatedFrontmatter = {
      ...existing.frontmatter,
      ...frontmatter,
    };

    // Handle explicit null/undefined removal (e.g. theme: null means "remove theme")
    for (const [key, value] of Object.entries(frontmatter ?? {})) {
      if (value === null || value === undefined) {
        delete (updatedFrontmatter as Record<string, unknown>)[key];
      }
    }

    note = await writeNote(notePath, content, updatedFrontmatter);
  } else {
    throw error(404, "Note not found");
  }

  await indexNote(notePath);

  return json({
    path: note.path,
    title: note.frontmatter.title,
    checksum: note.checksum,
  });
};

// DELETE /api/notes/:path — delete a note
export const DELETE: RequestHandler = async ({ params }) => {
  const notePath = params.path;
  if (!notePath) throw error(400, "path is required");

  if (!(await noteExists(notePath))) {
    throw error(404, "Note not found");
  }

  await removeNoteFromIndex(notePath);
  await deleteNote(notePath);

  return json({ success: true });
};
