import { json, error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq } from "drizzle-orm";

// GET /api/tags/:name — get notes with this tag
export const GET: RequestHandler = async ({ params }) => {
  const tagName = params.name;
  if (!tagName) throw error(400, "tag name is required");

  const notes = await db
    .select({
      id: schema.notes.id,
      path: schema.notes.path,
      title: schema.notes.title,
      modifiedAt: schema.notes.modifiedAt,
    })
    .from(schema.notes)
    .innerJoin(schema.noteTags, eq(schema.notes.id, schema.noteTags.noteId))
    .innerJoin(schema.tags, eq(schema.noteTags.tagId, schema.tags.id))
    .where(eq(schema.tags.name, tagName));

  return json(notes);
};
