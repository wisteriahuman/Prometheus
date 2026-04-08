export type ViewMode = "split" | "editor" | "preview";

export function isBareNoteContent(content: string, title?: string): boolean {
  const normalizedContent = content.replace(/\r\n/g, "\n").trim();
  if (!normalizedContent) return true;

  const normalizedTitle = title?.trim();
  if (!normalizedTitle) return false;

  return normalizedContent === `# ${normalizedTitle}`;
}
