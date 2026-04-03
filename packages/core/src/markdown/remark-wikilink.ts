import type { Root, Text, Link } from "mdast";
import type { Plugin } from "unified";
import { visit } from "unist-util-visit";

const WIKILINK_REGEX = /\[\[([^\]|]+)(?:\|([^\]]+))?\]\]/g;

/**
 * Remark plugin that transforms [[wikilink]] syntax into markdown links.
 * [[Page Name]] becomes [Page Name](/note/page-name.md)
 * [[Page Name|Display Text]] becomes [Display Text](/note/page-name.md)
 */
export const remarkWikilink: Plugin<[], Root> = () => {
  return (tree: Root) => {
    visit(tree, "text", (node: Text, index, parent) => {
      if (!parent || index === undefined) return;

      const value = node.value;
      WIKILINK_REGEX.lastIndex = 0;

      if (!WIKILINK_REGEX.test(value)) return;
      WIKILINK_REGEX.lastIndex = 0;

      const children: (Text | Link)[] = [];
      let lastIndex = 0;
      let match;

      while ((match = WIKILINK_REGEX.exec(value)) !== null) {
        // Text before the wikilink
        if (match.index > lastIndex) {
          children.push({
            type: "text",
            value: value.slice(lastIndex, match.index),
          });
        }

        const slug = match[1].trim();
        const displayText = match[2]?.trim() ?? slug;
        const href = `/note/${slugToPath(slug)}`;

        children.push({
          type: "link",
          url: href,
          children: [{ type: "text", value: displayText }],
          data: {
            hProperties: {
              className: "wikilink",
              "data-slug": slug,
            },
          },
        } as Link);

        lastIndex = match.index + match[0].length;
      }

      // Remaining text after last wikilink
      if (lastIndex < value.length) {
        children.push({
          type: "text",
          value: value.slice(lastIndex),
        });
      }

      if (children.length > 0) {
        parent.children.splice(index, 1, ...children);
      }
    });
  };
};

function slugToPath(slug: string): string {
  return slug.toLowerCase().replace(/\s+/g, "-") + ".md";
}
