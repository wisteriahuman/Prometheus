import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq, and } from "drizzle-orm";

// GET /api/tasks?filter=all|completed|pending
export const GET: RequestHandler = async ({ url }) => {
  const filter = url.searchParams.get("filter") ?? "all";

  let query = db
    .select({
      id: schema.tasks.id,
      noteId: schema.tasks.noteId,
      content: schema.tasks.content,
      completed: schema.tasks.completed,
      lineNumber: schema.tasks.lineNumber,
      dueDate: schema.tasks.dueDate,
      notePath: schema.notes.path,
      noteTitle: schema.notes.title,
    })
    .from(schema.tasks)
    .innerJoin(schema.notes, eq(schema.tasks.noteId, schema.notes.id));

  let results;
  if (filter === "completed") {
    results = await query.where(eq(schema.tasks.completed, true));
  } else if (filter === "pending") {
    results = await query.where(eq(schema.tasks.completed, false));
  } else {
    results = await query;
  }

  return json(results);
};
