import type { Root, Text } from "mdast";
import { visit } from "unist-util-visit";

export interface WikilinkInfo {
  slug: string;
  displayText: string | null;
  position: {
    line: number;
    column: number;
  } | null;
}

const WIKILINK_REGEX = /\[\[([^\]|]+)(?:\|([^\]]+))?\]\]/g;

export function extractWikilinks(ast: Root): WikilinkInfo[] {
  const wikilinks: WikilinkInfo[] = [];

  visit(ast, "text", (node: Text) => {
    let match;
    WIKILINK_REGEX.lastIndex = 0;

    while ((match = WIKILINK_REGEX.exec(node.value)) !== null) {
      const slug = match[1].trim();
      const displayText = match[2]?.trim() ?? null;

      wikilinks.push({
        slug,
        displayText,
        position: node.position
          ? {
              line: node.position.start.line,
              column: node.position.start.column + match.index,
            }
          : null,
      });
    }
  });

  return wikilinks;
}

export function extractWikilinksFromText(text: string): WikilinkInfo[] {
  const wikilinks: WikilinkInfo[] = [];
  let match;
  WIKILINK_REGEX.lastIndex = 0;

  while ((match = WIKILINK_REGEX.exec(text)) !== null) {
    wikilinks.push({
      slug: match[1].trim(),
      displayText: match[2]?.trim() ?? null,
      position: null,
    });
  }

  return wikilinks;
}

export function slugToPath(slug: string): string {
  return slug.toLowerCase().replace(/\s+/g, "-") + ".md";
}

export function pathToSlug(path: string): string {
  return path.replace(/\.md$/, "").replace(/\//g, "/");
}
