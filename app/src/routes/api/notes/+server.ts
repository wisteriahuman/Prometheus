import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import {
  listNotes,
  createNote,
  readNote,
  noteExists,
} from "$lib/server/fs/vault.js";
import { indexNote } from "$lib/server/indexer.js";

// GET /api/notes — list all notes
export const GET: RequestHandler = async () => {
  const notes = await listNotes();
  const noteData = await Promise.all(
    notes.map(async (path) => {
      try {
        const note = await readNote(path);
        return {
          path: note.path,
          title: note.frontmatter.title,
          tags: note.frontmatter.tags,
          modified: note.frontmatter.modified,
          theme: note.frontmatter.theme ?? null,
        };
      } catch {
        return { path, title: path, tags: [], modified: "", theme: null };
      }
    }),
  );

  return json(noteData);
};

// POST /api/notes — create a new note
export const POST: RequestHandler = async ({ request }) => {
  const { path, title } = await request.json();

  if (!path || typeof path !== "string") {
    return json({ error: "path is required" }, { status: 400 });
  }

  const notePath = path.endsWith(".md") ? path : `${path}.md`;

  if (await noteExists(notePath)) {
    return json({ error: "Note already exists" }, { status: 409 });
  }

  const note = await createNote(notePath, title);
  await indexNote(notePath);

  return json(
    {
      path: note.path,
      title: note.frontmatter.title,
      id: note.frontmatter.id,
    },
    { status: 201 },
  );
};
