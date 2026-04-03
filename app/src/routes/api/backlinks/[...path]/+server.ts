import { json, error } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq, like, or } from "drizzle-orm";

// GET /api/backlinks/:path — get backlinks for a note
export const GET: RequestHandler = async ({ params }) => {
  const notePath = params.path;
  if (!notePath) throw error(400, "path is required");

  // Get the note's slug (filename without extension, case-insensitive)
  const slug = notePath
    .replace(/\.md$/, "")
    .split("/")
    .pop() ?? "";

  // Find all links pointing to this slug
  const backlinks = await db
    .select({
      sourceId: schema.links.sourceId,
      targetSlug: schema.links.targetSlug,
      context: schema.links.context,
      sourcePath: schema.notes.path,
      sourceTitle: schema.notes.title,
    })
    .from(schema.links)
    .innerJoin(schema.notes, eq(schema.links.sourceId, schema.notes.id))
    .where(
      or(
        like(schema.links.targetSlug, slug),
        like(schema.links.targetSlug, `%/${slug}`),
      ),
    );

  return json(backlinks);
};
