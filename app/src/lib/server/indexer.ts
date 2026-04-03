import { db, schema } from "./db/index.js";
import { eq } from "drizzle-orm";
import { readNote, listNotes } from "./fs/vault.js";
import { parseMarkdownAst } from "@prometheus/core/markdown";
import { extractWikilinks } from "@prometheus/core/markdown";
import { extractTasks } from "@prometheus/core/markdown";

export async function indexNote(notePath: string): Promise<void> {
  const note = await readNote(notePath);

  // Check if note already indexed with same checksum
  const existing = await db
    .select({ checksum: schema.notes.checksum })
    .from(schema.notes)
    .where(eq(schema.notes.id, note.frontmatter.id))
    .limit(1);

  if (existing.length > 0 && existing[0].checksum === note.checksum) {
    return; // No changes
  }

  const ast = parseMarkdownAst(note.content);
  const wikilinks = extractWikilinks(ast);
  const tasks = extractTasks(ast);

  // Delete existing entry for this path (handles ID changes)
  await db.delete(schema.notes).where(eq(schema.notes.path, note.path));

  // Insert note
  await db
    .insert(schema.notes)
    .values({
      id: note.frontmatter.id,
      path: note.path,
      title: note.frontmatter.title,
      content: note.content,
      createdAt: new Date(note.frontmatter.created),
      modifiedAt: new Date(note.frontmatter.modified),
      theme: note.frontmatter.theme ?? null,
      checksum: note.checksum,
    })
    .onConflictDoUpdate({
      target: schema.notes.id,
      set: {
        path: note.path,
        title: note.frontmatter.title,
        content: note.content,
        modifiedAt: new Date(note.frontmatter.modified),
        theme: note.frontmatter.theme ?? null,
        checksum: note.checksum,
      },
    });

  // Re-index links
  await db
    .delete(schema.links)
    .where(eq(schema.links.sourceId, note.frontmatter.id));

  if (wikilinks.length > 0) {
    await db.insert(schema.links).values(
      wikilinks.map((wl) => ({
        sourceId: note.frontmatter.id,
        targetId: null,
        targetSlug: wl.slug,
        context: wl.displayText,
      })),
    );
  }

  // Re-index tasks
  await db
    .delete(schema.tasks)
    .where(eq(schema.tasks.noteId, note.frontmatter.id));

  if (tasks.length > 0) {
    await db.insert(schema.tasks).values(
      tasks.map((t) => ({
        noteId: note.frontmatter.id,
        content: t.content,
        completed: t.completed,
        lineNumber: t.lineNumber,
      })),
    );
  }

  // Re-index tags
  await db
    .delete(schema.noteTags)
    .where(eq(schema.noteTags.noteId, note.frontmatter.id));

  if (note.frontmatter.tags.length > 0) {
    for (const tagName of note.frontmatter.tags) {
      // Upsert tag
      await db
        .insert(schema.tags)
        .values({ name: tagName })
        .onConflictDoNothing();

      const [tag] = await db
        .select({ id: schema.tags.id })
        .from(schema.tags)
        .where(eq(schema.tags.name, tagName))
        .limit(1);

      if (tag) {
        await db
          .insert(schema.noteTags)
          .values({ noteId: note.frontmatter.id, tagId: tag.id })
          .onConflictDoNothing();
      }
    }
  }
}

export async function indexAllNotes(): Promise<number> {
  const notes = await listNotes();
  const vaultPaths = new Set(notes);
  let count = 0;

  // Index all existing files
  for (const notePath of notes) {
    try {
      await indexNote(notePath);
      count++;
    } catch (e) {
      console.error(`Failed to index ${notePath}:`, e);
    }
  }

  // Remove stale entries (files that no longer exist in vault)
  const dbNotes = await db
    .select({ id: schema.notes.id, path: schema.notes.path })
    .from(schema.notes);

  for (const dbNote of dbNotes) {
    if (!vaultPaths.has(dbNote.path)) {
      await db.delete(schema.notes).where(eq(schema.notes.id, dbNote.id));
    }
  }

  return count;
}

export async function removeNoteFromIndex(notePath: string): Promise<void> {
  const [note] = await db
    .select({ id: schema.notes.id })
    .from(schema.notes)
    .where(eq(schema.notes.path, notePath))
    .limit(1);

  if (note) {
    await db.delete(schema.notes).where(eq(schema.notes.id, note.id));
  }
}
