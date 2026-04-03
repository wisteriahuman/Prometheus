import type { PageServerLoad } from "./$types.js";
import { ensureDailyNote, getDailyNotePath, getRecentDailyDates } from "$lib/server/fs/daily.js";
import { noteExists } from "$lib/server/fs/vault.js";
import { markdownToHtml } from "@prometheus/core/markdown";
import { indexNote } from "$lib/server/indexer.js";

export const load: PageServerLoad = async ({ url }) => {
  const dateParam = url.searchParams.get("date");
  const date = dateParam ? new Date(dateParam) : new Date();

  const note = await ensureDailyNote(date);
  const path = await getDailyNotePath(date);
  await indexNote(path);

  const html = await markdownToHtml(note.content);

  // Check which recent dates have notes
  const recentDates = getRecentDailyDates(14);
  const recentNotes = await Promise.all(
    recentDates.map(async (d) => {
      const p = await getDailyNotePath(d);
      return {
        date: d.toISOString().slice(0, 10),
        exists: await noteExists(p),
      };
    }),
  );

  return {
    path: note.path,
    title: note.frontmatter.title,
    content: note.content,
    frontmatter: note.frontmatter,
    html,
    currentDate: date.toISOString().slice(0, 10),
    recentNotes,
  };
};
