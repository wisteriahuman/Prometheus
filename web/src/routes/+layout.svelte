<script lang="ts">
  import "../app.css";
  import { onMount } from "svelte";
  import Sidebar from "$lib/components/sidebar/Sidebar.svelte";
  import SearchModal from "$lib/components/search/SearchModal.svelte";
  import CommandPalette from "$lib/components/command/CommandPalette.svelte";
  import ShortcutHelp from "$lib/components/ui/ShortcutHelp.svelte";
  import { sidebarOpen, toggleSidebar } from "$lib/stores/sidebar";
  import { initTheme } from "$lib/stores/theme";
  import { goto } from "$app/navigation";
  import { Menu } from "lucide-svelte";

  let { children } = $props();

  let searchOpen = $state(false);
  let searchMode = $state<"search" | "switcher">("search");
  let commandPaletteOpen = $state(false);
  let shortcutHelpOpen = $state(false);
  let vaultPath = $state("");

  onMount(async () => {
    initTheme();

    // Fetch vault path from API
    try {
      const res = await fetch("/api/config");
      const config = await res.json();
      vaultPath = config.vaultPath;
    } catch {}

    const handleSearch = () => { searchMode = "search"; searchOpen = true; };
    const handleSwitcher = () => { searchMode = "switcher"; searchOpen = true; };
    const handleShortcuts = () => { shortcutHelpOpen = true; };

    window.addEventListener("prometheus:search", handleSearch);
    window.addEventListener("prometheus:switcher", handleSwitcher);
    window.addEventListener("prometheus:shortcuts", handleShortcuts);

    return () => {
      window.removeEventListener("prometheus:search", handleSearch);
      window.removeEventListener("prometheus:switcher", handleSwitcher);
      window.removeEventListener("prometheus:shortcuts", handleShortcuts);
    };
  });

  function handleKeydown(e: KeyboardEvent) {
    const mod = e.metaKey || e.ctrlKey;

    if (mod && e.key === "\\") {
      e.preventDefault();
      toggleSidebar();
    } else if (mod && e.key === ".") {
      e.preventDefault();
      commandPaletteOpen = true;
      setTimeout(() => window.dispatchEvent(new CustomEvent("prometheus:theme-mode")), 50);
    } else if (mod && e.shiftKey && e.key === "f") {
      e.preventDefault();
      searchMode = "search";
      searchOpen = true;
    } else if (mod && e.key === "p") {
      e.preventDefault();
      searchMode = "switcher";
      searchOpen = true;
    } else if (mod && e.key === "k") {
      e.preventDefault();
      commandPaletteOpen = true;
    } else if (mod && e.key === "d") {
      e.preventDefault();
      const today = new Date().toISOString().slice(0, 10);
      goto(`/note/daily/${today}.md`);
    } else if (mod && e.key === "g") {
      e.preventDefault();
      goto("/graph");
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="flex h-screen overflow-hidden">
  <Sidebar {vaultPath} />

  <main class="flex flex-1 flex-col overflow-hidden">
    <header class="flex h-12 shrink-0 items-center border-b border-border px-4">
      <button
        onclick={toggleSidebar}
        class="mr-3 rounded-md p-1.5 text-text-muted transition-colors hover:bg-bg-card-hover hover:text-text-main md:hidden"
        aria-label="サイドバー切替"
      >
        <Menu size={16} />
      </button>
      <div class="flex-1"></div>
      <div class="flex items-center gap-2 text-xs text-text-dim">
        <button
          onclick={() => { searchMode = "switcher"; searchOpen = true; }}
          class="flex items-center gap-1.5 rounded-md border border-border px-2 py-1 transition-colors hover:border-text-dim hover:text-text-muted"
        >
          <kbd class="font-mono text-[10px]">⌘P</kbd>
          <span class="hidden sm:inline">ファイル切替</span>
        </button>
        <button
          onclick={() => { commandPaletteOpen = true; }}
          class="flex items-center gap-1.5 rounded-md border border-border px-2 py-1 transition-colors hover:border-text-dim hover:text-text-muted"
        >
          <kbd class="font-mono text-[10px]">⌘K</kbd>
          <span class="hidden sm:inline">コマンド</span>
        </button>
      </div>
    </header>

    <div class="flex-1 overflow-y-auto p-5">
      {@render children()}
    </div>
  </main>
</div>

<SearchModal
  open={searchOpen}
  mode={searchMode}
  onclose={() => (searchOpen = false)}
/>

<CommandPalette
  open={commandPaletteOpen}
  onclose={() => (commandPaletteOpen = false)}
/>

<ShortcutHelp
  open={shortcutHelpOpen}
  onclose={() => (shortcutHelpOpen = false)}
/>
