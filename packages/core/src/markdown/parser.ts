import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkFrontmatter from "remark-frontmatter";
import remarkRehype from "remark-rehype";
import rehypeStringify from "rehype-stringify";
import rehypeSanitize, { defaultSchema } from "rehype-sanitize";
import type { Root } from "mdast";
import { remarkWikilink } from "./remark-wikilink.js";
import { extractWikilinks, type WikilinkInfo } from "./wikilink.js";
import { extractTasks, type TaskInfo } from "./task.js";

const sanitizeSchema = {
  ...defaultSchema,
  attributes: {
    ...defaultSchema.attributes,
    code: [...(defaultSchema.attributes?.code ?? []), "className"],
    span: [...(defaultSchema.attributes?.span ?? []), "className", "style"],
    a: [...(defaultSchema.attributes?.a ?? []), "className", "dataSlug"],
  },
  protocols: {
    ...defaultSchema.protocols,
    href: [...(defaultSchema.protocols?.href ?? []), "data"],
  },
};

const htmlProcessor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkFrontmatter, ["yaml"])
  .use(remarkWikilink)
  .use(remarkRehype, { allowDangerousHtml: true })
  .use(rehypeSanitize, sanitizeSchema)
  .use(rehypeStringify);

const astProcessor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkFrontmatter, ["yaml"]);

export async function markdownToHtml(markdown: string): Promise<string> {
  const result = await htmlProcessor.process(markdown);
  return String(result);
}

export function parseMarkdownAst(markdown: string): Root {
  return astProcessor.parse(markdown) as Root;
}

export interface ParsedMarkdown {
  html: string;
  wikilinks: WikilinkInfo[];
  tasks: TaskInfo[];
}

export async function parseMarkdown(markdown: string): Promise<ParsedMarkdown> {
  const [html, ast] = await Promise.all([
    markdownToHtml(markdown),
    Promise.resolve(parseMarkdownAst(markdown)),
  ]);

  const wikilinks = extractWikilinks(ast);
  const tasks = extractTasks(ast);

  return { html, wikilinks, tasks };
}
