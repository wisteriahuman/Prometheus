import type { PageServerLoad } from "./$types.js";
import { readNote, noteExists, createNote } from "$lib/server/fs/vault.js";
import { markdownToHtml } from "@prometheus/core/markdown";
import { indexNote } from "$lib/server/indexer.js";
import { error } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ params }) => {
  const notePath = params.path;
  if (!notePath) throw error(400, "path is required");

  // Auto-create if it doesn't exist
  if (!(await noteExists(notePath))) {
    if (notePath.endsWith(".md")) {
      await createNote(notePath);
      await indexNote(notePath);
    } else {
      throw error(404, "Note not found");
    }
  }

  const note = await readNote(notePath);
  const html = await markdownToHtml(note.content);

  return {
    path: note.path,
    title: note.frontmatter.title,
    content: note.content,
    frontmatter: note.frontmatter,
    html,
  };
};
