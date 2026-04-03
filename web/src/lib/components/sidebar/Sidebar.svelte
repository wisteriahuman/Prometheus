<script lang="ts">
  import { sidebarOpen, toggleSidebar } from "$lib/stores/sidebar";
  import FileTree from "./FileTree.svelte";
  import TagList from "./TagList.svelte";
  import { Home, CalendarDays, Network, CheckSquare, Search, Plus, X, FolderPlus, HardDrive } from "lucide-svelte";
  interface VaultFileEntry {
    name: string;
    path: string;
    isDirectory: boolean;
    children?: VaultFileEntry[];
  }

  interface Props {
    vaultPath?: string;
  }

  let { vaultPath = "" }: Props = $props();

  let fileTree: VaultFileEntry[] = $state([]);
  let noteThemes: Map<string, string> = $state(new Map());
  let showNewNote = $state(false);
  let showNewFolder = $state(false);
  let newNoteName = $state("");
  let newFolderName = $state("");

  const NAV_ITEMS = [
    { href: "/", label: "ホーム", icon: Home },
    { href: "/daily", label: "デイリー", icon: CalendarDays },
    { href: "/graph", label: "グラフ", icon: Network },
    { href: "/tasks", label: "タスク", icon: CheckSquare },
    { href: "/search", label: "検索", icon: Search },
  ] as const;

  async function loadFileTree() {
    try {
      const [treeRes, notesRes] = await Promise.all([
        fetch("/api/tree"),
        fetch("/api/notes"),
      ]);
      fileTree = await treeRes.json();

      const notes: { path: string; theme: string | null }[] = await notesRes.json();
      const themes = new Map<string, string>();
      for (const n of notes) {
        if (n.theme) themes.set(n.path, n.theme);
      }
      noteThemes = themes;
    } catch {
      fileTree = [];
    }
  }

  async function createNote() {
    if (!newNoteName.trim()) return;

    const path = newNoteName.endsWith(".md") ? newNoteName : `${newNoteName}.md`;

    try {
      const res = await fetch("/api/notes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path, title: newNoteName.replace(/\.md$/, "") }),
      });

      if (res.ok) {
        newNoteName = "";
        showNewNote = false;
        await loadFileTree();
        window.location.href = `/note/${path}`;
      }
    } catch (e) {
      console.error("Failed to create note:", e);
    }
  }

  async function createFolder() {
    if (!newFolderName.trim()) return;

    try {
      const res = await fetch(`/api/folders/${newFolderName.trim()}`, {
        method: "POST",
      });

      if (res.ok) {
        newFolderName = "";
        showNewFolder = false;
        await loadFileTree();
      }
    } catch (e) {
      console.error("Failed to create folder:", e);
    }
  }

  import { onMount } from "svelte";

  onMount(() => {
    loadFileTree();

    const handleRefresh = () => loadFileTree();
    window.addEventListener("prometheus:refresh-tree", handleRefresh);
    return () => window.removeEventListener("prometheus:refresh-tree", handleRefresh);
  });
</script>

<aside
  class="fixed left-0 top-0 z-40 flex h-screen flex-col border-r border-border bg-bg-card transition-all duration-200 ease-in-out
    {$sidebarOpen ? 'w-60 translate-x-0' : 'w-0 -translate-x-full overflow-hidden'}
    md:relative"
>
  <!-- Header -->
  <div class="flex h-12 shrink-0 items-center justify-between border-b border-border px-4">
    <a href="/" class="font-mono text-base font-semibold tracking-tight text-primary-light">
      Prometheus
    </a>
    <button
      onclick={toggleSidebar}
      class="rounded-md p-1 text-text-muted transition-colors hover:bg-bg-card-hover hover:text-text-main"
      aria-label="閉じる"
    >
      <X size={16} />
    </button>
  </div>

  <!-- Navigation -->
  <nav class="px-2 py-2">
    {#each NAV_ITEMS as item}
      <a
        href={item.href}
        class="flex items-center gap-2.5 rounded-md px-3 py-1.5 text-[13px] text-text-muted transition-colors hover:bg-bg-card-hover hover:text-text-main"
      >
        <item.icon size={15} strokeWidth={1.8} />
        <span>{item.label}</span>
      </a>
    {/each}
  </nav>

  <!-- File Tree -->
  <div class="flex-1 overflow-y-auto border-t border-border px-2 py-2">
    <div class="mb-1.5 flex items-center justify-between px-3">
      <p class="text-[11px] font-medium uppercase tracking-wider text-text-dim">
        ノート
      </p>
      <div class="flex items-center gap-0.5">
        <button
          onclick={() => { showNewFolder = !showNewFolder; showNewNote = false; }}
          class="rounded p-0.5 text-text-dim transition-colors hover:bg-bg-card-hover hover:text-text-main"
          aria-label="新規フォルダ"
          title="新規フォルダ"
        >
          <FolderPlus size={13} />
        </button>
        <button
          onclick={() => { showNewNote = !showNewNote; showNewFolder = false; }}
          class="rounded p-0.5 text-text-dim transition-colors hover:bg-bg-card-hover hover:text-text-main"
          aria-label="新規ノート"
          title="新規ノート"
        >
          <Plus size={14} />
        </button>
      </div>
    </div>

    {#if showNewFolder}
      <div class="mb-2 px-2">
        <form onsubmit={(e) => { e.preventDefault(); createFolder(); }}>
          <input
            bind:value={newFolderName}
            placeholder="フォルダ名"
            class="w-full rounded border border-border bg-bg-input px-2.5 py-1.5 text-xs text-text-main placeholder-text-dim outline-none focus:border-primary"
          />
        </form>
      </div>
    {/if}

    {#if showNewNote}
      <div class="mb-2 px-2">
        <form onsubmit={(e) => { e.preventDefault(); createNote(); }}>
          <input
            bind:value={newNoteName}
            placeholder="ファイル名.md（folder/name.md も可）"
            class="w-full rounded border border-border bg-bg-input px-2.5 py-1.5 text-xs text-text-main placeholder-text-dim outline-none focus:border-primary"
          />
        </form>
      </div>
    {/if}

    {#if fileTree.length > 0}
      <FileTree entries={fileTree} {noteThemes} onrefresh={loadFileTree} />
    {:else}
      <p class="px-3 py-2 text-xs text-text-dim">ノートがありません</p>
    {/if}
  </div>

  <!-- Tags -->
  <TagList />

  <!-- Vault path -->
  {#if vaultPath}
    <div class="shrink-0 border-t border-border px-3 py-2">
      <div class="flex items-center gap-1.5 text-[10px] text-text-dim" title={vaultPath}>
        <HardDrive size={10} class="shrink-0" />
        <span class="truncate">{vaultPath}</span>
      </div>
    </div>
  {/if}
</aside>

<!-- Overlay for mobile -->
{#if $sidebarOpen}
  <button
    class="fixed inset-0 z-30 bg-black/50 md:hidden"
    onclick={toggleSidebar}
    aria-label="閉じる"
  ></button>
{/if}
