import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types.js";
import { db, schema } from "$lib/server/db/index.js";
import { eq } from "drizzle-orm";

export interface GraphNode {
  id: string;
  path: string;
  title: string;
  tags: string[];
}

export interface GraphLink {
  source: string;
  target: string;
  slug: string;
}

export interface GraphData {
  nodes: GraphNode[];
  links: GraphLink[];
}

// GET /api/graph — get graph data for all notes
export const GET: RequestHandler = async () => {
  // Get all notes
  const notes = await db
    .select({
      id: schema.notes.id,
      path: schema.notes.path,
      title: schema.notes.title,
    })
    .from(schema.notes);

  // Build path-to-id and slug-to-id maps
  const pathToId = new Map<string, string>();
  const slugToId = new Map<string, string>();

  for (const note of notes) {
    pathToId.set(note.path, note.id);
    // Map multiple slug variations to the same note
    const fileName = note.path.replace(/\.md$/, "").split("/").pop() ?? "";
    const fullSlug = note.path.replace(/\.md$/, "");
    slugToId.set(fileName.toLowerCase(), note.id);
    slugToId.set(fullSlug.toLowerCase(), note.id);
    // Also map with spaces converted to hyphens (for [[Graph View]] -> graph-view.md)
    slugToId.set(fileName.toLowerCase().replace(/-/g, " "), note.id);
  }

  // Get all links
  const allLinks = await db
    .select({
      sourceId: schema.links.sourceId,
      targetSlug: schema.links.targetSlug,
    })
    .from(schema.links);

  // Get tags for each note
  const noteTags = await db
    .select({
      noteId: schema.noteTags.noteId,
      tagName: schema.tags.name,
    })
    .from(schema.noteTags)
    .innerJoin(schema.tags, eq(schema.noteTags.tagId, schema.tags.id));

  const tagsByNote = new Map<string, string[]>();
  for (const nt of noteTags) {
    if (!tagsByNote.has(nt.noteId)) tagsByNote.set(nt.noteId, []);
    tagsByNote.get(nt.noteId)!.push(nt.tagName);
  }

  // Build graph nodes
  const nodes: GraphNode[] = notes.map((n) => ({
    id: n.id,
    path: n.path,
    title: n.title,
    tags: tagsByNote.get(n.id) ?? [],
  }));

  // Build graph links (only resolved ones)
  const links: GraphLink[] = [];
  const nodeIds = new Set(notes.map((n) => n.id));

  for (const link of allLinks) {
    if (!nodeIds.has(link.sourceId)) continue;

    const slug = link.targetSlug.toLowerCase();
    // Try multiple slug variations
    const targetId =
      slugToId.get(slug) ??
      slugToId.get(slug.replace(/\s+/g, "-")) ??
      slugToId.get(slug.replace(/-/g, " "));

    if (targetId && nodeIds.has(targetId) && targetId !== link.sourceId) {
      // Avoid duplicate edges
      const edgeKey = [link.sourceId, targetId].sort().join(":");
      if (!links.some((l) => {
        const src = typeof l.source === "string" ? l.source : l.source;
        const tgt = typeof l.target === "string" ? l.target : l.target;
        return [src, tgt].sort().join(":") === edgeKey;
      })) {
        links.push({
          source: link.sourceId,
          target: targetId,
          slug: link.targetSlug,
        });
      }
    }
  }

  return json({ nodes, links } satisfies GraphData);
};
