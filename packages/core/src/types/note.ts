export interface NoteFrontmatter {
  id: string;
  title: string;
  created: string;
  modified: string;
  tags: string[];
  theme?: string;
}

export interface Note {
  id: string;
  path: string;
  title: string;
  content: string;
  frontmatter: NoteFrontmatter;
  createdAt: Date;
  modifiedAt: Date;
  checksum: string;
}

export interface NoteLink {
  sourceId: string;
  targetId: string | null;
  targetSlug: string;
  context: string | null;
}

export interface NoteTask {
  id: number;
  noteId: string;
  content: string;
  completed: boolean;
  lineNumber: number;
  dueDate: string | null;
}
