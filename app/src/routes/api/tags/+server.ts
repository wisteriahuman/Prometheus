import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq, sql } from "drizzle-orm";

// GET /api/tags — list all tags with note counts
export const GET: RequestHandler = async () => {
  const tags = await db
    .select({
      id: schema.tags.id,
      name: schema.tags.name,
      count: sql<number>`count(${schema.noteTags.noteId})`,
    })
    .from(schema.tags)
    .leftJoin(schema.noteTags, eq(schema.tags.id, schema.noteTags.tagId))
    .groupBy(schema.tags.id, schema.tags.name)
    .orderBy(schema.tags.name);

  return json(tags);
};
