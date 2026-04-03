import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { like, or, sql } from "drizzle-orm";

// GET /api/search?q=query — full-text search
export const GET: RequestHandler = async ({ url }) => {
  const query = url.searchParams.get("q")?.trim();

  if (!query) {
    return json([]);
  }

  // Use LIKE for simple search (FTS5 would need raw SQL)
  // For now, use pattern matching which works across both SQLite and PostgreSQL
  const pattern = `%${query}%`;

  const results = await db
    .select({
      id: schema.notes.id,
      path: schema.notes.path,
      title: schema.notes.title,
      content: schema.notes.content,
      modifiedAt: schema.notes.modifiedAt,
    })
    .from(schema.notes)
    .where(
      or(
        like(schema.notes.title, pattern),
        like(schema.notes.content, pattern),
      ),
    )
    .limit(50);

  // Generate snippets around matching text
  const resultsWithSnippets = results.map((r) => {
    const lowerContent = r.content.toLowerCase();
    const lowerQuery = query.toLowerCase();
    const idx = lowerContent.indexOf(lowerQuery);

    let snippet = "";
    if (idx >= 0) {
      const start = Math.max(0, idx - 60);
      const end = Math.min(r.content.length, idx + query.length + 60);
      snippet = (start > 0 ? "..." : "") +
        r.content.slice(start, end) +
        (end < r.content.length ? "..." : "");
    }

    return {
      id: r.id,
      path: r.path,
      title: r.title,
      snippet,
      modifiedAt: r.modifiedAt,
    };
  });

  return json(resultsWithSnippets);
};
