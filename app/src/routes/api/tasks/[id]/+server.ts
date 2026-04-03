import { json, error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq } from "drizzle-orm";
import { readNote, writeNote } from "$lib/server/fs/vault.js";
import { indexNote } from "$lib/server/indexer.js";

// PATCH /api/tasks/:id — toggle task completion
export const PATCH: RequestHandler = async ({ params, request }) => {
  const taskId = parseInt(params.id);
  if (isNaN(taskId)) throw error(400, "Invalid task ID");

  const { completed } = await request.json();
  if (typeof completed !== "boolean") {
    throw error(400, "completed must be a boolean");
  }

  // Get the task and its note
  const [task] = await db
    .select({
      id: schema.tasks.id,
      noteId: schema.tasks.noteId,
      lineNumber: schema.tasks.lineNumber,
      content: schema.tasks.content,
    })
    .from(schema.tasks)
    .where(eq(schema.tasks.id, taskId))
    .limit(1);

  if (!task) throw error(404, "Task not found");

  // Get the note path
  const [note] = await db
    .select({ path: schema.notes.path })
    .from(schema.notes)
    .where(eq(schema.notes.id, task.noteId))
    .limit(1);

  if (!note) throw error(404, "Note not found");

  // Read the note and toggle the checkbox in the actual markdown
  const noteData = await readNote(note.path);
  const lines = noteData.rawContent.split("\n");

  // Find the line with the task (1-indexed in AST, but rawContent includes frontmatter)
  // We need to search for the task by content
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const uncheckedMatch = line.match(/^(\s*[-*+]\s)\[ \](.*)/);
    const checkedMatch = line.match(/^(\s*[-*+]\s)\[x\](.*)/i);

    if (completed && uncheckedMatch && line.includes(task.content.slice(0, 30))) {
      lines[i] = `${uncheckedMatch[1]}[x]${uncheckedMatch[2]}`;
      break;
    } else if (!completed && checkedMatch && line.includes(task.content.slice(0, 30))) {
      lines[i] = `${checkedMatch[1]}[ ]${checkedMatch[2]}`;
      break;
    }
  }

  // Write back
  const newRawContent = lines.join("\n");
  const matter = await import("gray-matter");
  const parsed = matter.default(newRawContent);
  await writeNote(note.path, parsed.content, noteData.frontmatter);
  await indexNote(note.path);

  return json({ success: true });
};
