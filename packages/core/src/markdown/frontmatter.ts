import type { NoteFrontmatter } from "../types/note.js";

export function validateFrontmatter(data: Record<string, unknown>): NoteFrontmatter {
  return {
    id: typeof data.id === "string" ? data.id : "",
    title: typeof data.title === "string" ? data.title : "Untitled",
    created: typeof data.created === "string" ? data.created : new Date().toISOString(),
    modified: typeof data.modified === "string" ? data.modified : new Date().toISOString(),
    tags: Array.isArray(data.tags) ? data.tags.filter((t): t is string => typeof t === "string") : [],
    theme: typeof data.theme === "string" ? data.theme : undefined,
  };
}
