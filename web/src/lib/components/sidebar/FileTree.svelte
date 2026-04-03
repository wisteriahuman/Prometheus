<script lang="ts">
  import FileTree from "./FileTree.svelte";
  import { ChevronRight, ChevronDown, FileText, Folder, Trash2, Pencil, FolderX } from "lucide-svelte";
  import { goto } from "$app/navigation";
  import { BUILTIN_THEMES } from "$lib/stores/theme";
  import type { VaultFileEntry } from "$lib/server/fs/vault.js";

  interface Props {
    entries: VaultFileEntry[];
    noteThemes?: Map<string, string>;
    depth?: number;
    onrefresh?: () => void;
  }

  let { entries, noteThemes = new Map(), depth = 0, onrefresh }: Props = $props();

  let expanded: Record<string, boolean> = $state({});
  let contextMenu = $state<{ x: number; y: number; entry: VaultFileEntry } | null>(null);
  let renaming = $state<string | null>(null);
  let renameValue = $state("");

  function getThemeColor(path: string): string | null {
    const slug = noteThemes.get(path);
    if (!slug) return null;
    const theme = BUILTIN_THEMES.find((t) => t.slug === slug);
    return theme?.colors.primary ?? null;
  }

  function toggleDir(path: string) {
    expanded[path] = !expanded[path];
  }

  function showContextMenu(e: MouseEvent, entry: VaultFileEntry) {
    e.preventDefault();
    contextMenu = { x: e.clientX, y: e.clientY, entry };
  }

  function closeContextMenu() {
    contextMenu = null;
  }

  async function deleteNote(path: string) {
    closeContextMenu();
    if (!confirm(`「${path}」を削除しますか？`)) return;

    try {
      await fetch(`/api/notes/${path}`, { method: "DELETE" });
      onrefresh?.();
      if (window.location.pathname === `/note/${path}`) {
        goto("/");
      }
    } catch (e) {
      console.error("Failed to delete:", e);
    }
  }

  async function deleteFolder(path: string) {
    closeContextMenu();

    // First try without force
    let res = await fetch(`/api/folders/${path}`, { method: "DELETE" });
    if (!res.ok) {
      // Folder not empty — ask user
      if (!confirm(`「${path}」フォルダとその中のノートを全て削除しますか？`)) return;
      res = await fetch(`/api/folders/${path}?force=true`, { method: "DELETE" });
    }

    if (res.ok) {
      onrefresh?.();
    }
  }

  function startRename(entry: VaultFileEntry) {
    closeContextMenu();
    renaming = entry.path;
    renameValue = entry.name.replace(/\.md$/, "");
  }

  async function finishRename(oldPath: string) {
    if (!renameValue.trim()) {
      renaming = null;
      return;
    }

    const dir = oldPath.includes("/") ? oldPath.substring(0, oldPath.lastIndexOf("/") + 1) : "";
    const newPath = `${dir}${renameValue}.md`;

    if (newPath === oldPath) {
      renaming = null;
      return;
    }

    try {
      const res = await fetch(`/api/notes/${oldPath}`);
      if (!res.ok) return;
      const note = await res.json();

      await fetch("/api/notes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path: newPath, title: renameValue }),
      });

      await fetch(`/api/notes/${newPath}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: note.content,
          frontmatter: { ...note.frontmatter, title: renameValue },
        }),
      });

      await fetch(`/api/notes/${oldPath}`, { method: "DELETE" });
      onrefresh?.();
      goto(`/note/${newPath}`);
    } catch (e) {
      console.error("Failed to rename:", e);
    } finally {
      renaming = null;
    }
  }
</script>

<svelte:window onclick={closeContextMenu} />

{#each entries as entry}
  {#if entry.isDirectory}
    <div>
      <button
        onclick={() => toggleDir(entry.path)}
        oncontextmenu={(e) => showContextMenu(e, entry)}
        class="flex w-full items-center gap-1.5 rounded-md px-2 py-1 text-[13px] text-text-muted transition-colors hover:bg-bg-card-hover hover:text-text-main"
        style="padding-left: {depth * 12 + 8}px"
      >
        {#if expanded[entry.path]}
          <ChevronDown size={12} />
        {:else}
          <ChevronRight size={12} />
        {/if}
        <Folder size={13} strokeWidth={1.8} />
        <span>{entry.name}</span>
      </button>

      {#if expanded[entry.path] && entry.children}
        <FileTree entries={entry.children} {noteThemes} depth={depth + 1} {onrefresh} />
      {/if}
    </div>
  {:else}
    {#if renaming === entry.path}
      <div
        class="flex items-center gap-1.5 px-2 py-1"
        style="padding-left: {depth * 12 + 22}px"
      >
        <FileText size={13} strokeWidth={1.8} class="shrink-0 text-text-dim" />
        <form onsubmit={(e) => { e.preventDefault(); finishRename(entry.path); }} class="flex-1">
          <input
            bind:value={renameValue}
            class="w-full rounded border border-primary bg-bg-input px-1.5 py-0.5 text-xs text-text-main outline-none"
            onblur={() => finishRename(entry.path)}
          />
        </form>
      </div>
    {:else}
      <a
        href="/note/{entry.path}"
        oncontextmenu={(e) => showContextMenu(e, entry)}
        class="flex items-center gap-1.5 rounded-md px-2 py-1 text-[13px] text-text-muted transition-colors hover:bg-bg-card-hover hover:text-text-main"
        style="padding-left: {depth * 12 + 22}px"
      >
        <FileText size={13} strokeWidth={1.8} />
        <span class="flex-1">{entry.name.replace(/\.md$/, '')}</span>
        {#if getThemeColor(entry.path)}
          <span
            class="h-2 w-2 shrink-0 rounded-full"
            style="background-color: {getThemeColor(entry.path)}"
            title="テーマ: {noteThemes.get(entry.path)}"
          ></span>
        {/if}
      </a>
    {/if}
  {/if}
{/each}

<!-- Context Menu -->
{#if contextMenu}
  <div
    class="fixed z-50 min-w-[140px] rounded-lg border border-border bg-bg-card py-1 shadow-xl"
    style="left: {contextMenu.x}px; top: {contextMenu.y}px"
  >
    {#if contextMenu.entry.isDirectory}
      <button
        onclick={() => deleteFolder(contextMenu!.entry.path)}
        class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-error hover:bg-error/10"
      >
        <FolderX size={12} />
        フォルダを削除
      </button>
    {:else}
      <button
        onclick={() => startRename(contextMenu!.entry)}
        class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-text-muted hover:bg-bg-card-hover hover:text-text-main"
      >
        <Pencil size={12} />
        名前を変更
      </button>
      <button
        onclick={() => deleteNote(contextMenu!.entry.path)}
        class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-error hover:bg-error/10"
      >
        <Trash2 size={12} />
        削除
      </button>
    {/if}
  </div>
{/if}
