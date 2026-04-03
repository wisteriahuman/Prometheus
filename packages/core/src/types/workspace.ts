export type DbEngine = "sqlite" | "postgresql";

export interface WorkspaceConfig {
  name: string;
  engine: DbEngine;
  dbUrl?: string;
  theme?: string;
  dailyNoteTemplate?: string;
}

export interface Workspace {
  id: string;
  name: string;
  vaultPath: string;
  engine: DbEngine;
  dbUrl: string | null;
  createdAt: Date;
  ownerId: string;
}
