<script lang="ts">
  import { goto } from "$app/navigation";

  interface SearchResult {
    id: string;
    path: string;
    title: string;
    snippet: string;
  }

  interface Props {
    open: boolean;
    mode?: "search" | "switcher";
    onclose: () => void;
  }

  let { open, mode = "search", onclose }: Props = $props();

  let query = $state("");
  let results: SearchResult[] = $state([]);
  let selectedIndex = $state(0);
  let loading = $state(false);
  let debounceTimer: ReturnType<typeof setTimeout> | undefined;
  let inputEl: HTMLInputElement;

  $effect(() => {
    if (open) {
      query = "";
      results = [];
      selectedIndex = 0;
      // Load all notes for switcher mode
      if (mode === "switcher") {
        loadAllNotes();
      }
      setTimeout(() => inputEl?.focus(), 50);
    }
  });

  async function loadAllNotes() {
    try {
      const res = await fetch("/api/notes");
      const notes = await res.json();
      results = notes.map((n: any) => ({
        id: n.path,
        path: n.path,
        title: n.title,
        snippet: "",
      }));
    } catch {
      results = [];
    }
  }

  async function search(q: string) {
    if (!q.trim()) {
      if (mode === "switcher") {
        await loadAllNotes();
      } else {
        results = [];
      }
      return;
    }

    loading = true;
    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(q)}`);
      results = await res.json();
    } catch {
      results = [];
    } finally {
      loading = false;
    }
    selectedIndex = 0;
  }

  function handleInput() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => search(query), 150);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      e.preventDefault();
      onclose();
    } else if (e.key === "ArrowDown" || (e.ctrlKey && e.key === "n")) {
      e.preventDefault();
      selectedIndex = Math.min(selectedIndex + 1, results.length - 1);
    } else if (e.key === "ArrowUp" || (e.ctrlKey && e.key === "p")) {
      e.preventDefault();
      selectedIndex = Math.max(selectedIndex - 1, 0);
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (results[selectedIndex]) {
        navigateTo(results[selectedIndex].path);
      } else if (query.trim() && mode === "switcher") {
        // Create new note
        createAndNavigate(query.trim());
      }
    }
  }

  function navigateTo(path: string) {
    onclose();
    goto(`/note/${path}`);
  }

  async function createAndNavigate(name: string) {
    const path = name.endsWith(".md") ? name : `${name}.md`;
    try {
      const res = await fetch("/api/notes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path }),
      });
      if (res.ok) {
        onclose();
        goto(`/note/${path}`);
      }
    } catch {
      // failed
    }
  }
</script>

{#if open}
  <!-- Backdrop -->
  <div class="fixed inset-0 z-50 flex items-start justify-center pt-[15vh]">
    <button
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={onclose}
      aria-label="閉じる"
    ></button>

    <div class="relative z-10 w-full max-w-md overflow-hidden rounded-xl border border-border bg-bg-card shadow-2xl">
      <div class="flex items-center gap-2 border-b border-border px-4">
        <input
          bind:this={inputEl}
          bind:value={query}
          oninput={handleInput}
          onkeydown={handleKeydown}
          placeholder={mode === "search" ? "ノートを検索..." : "ノートを開く... (Enterで新規作成)"}
          class="flex-1 bg-transparent py-3 text-sm text-text-main outline-none placeholder:text-text-dim"
        />
        {#if loading}
          <span class="text-xs text-text-dim">...</span>
        {/if}
      </div>

      {#if results.length > 0}
        <div class="max-h-72 overflow-y-auto py-1">
          {#each results as result, i}
            <button
              class="flex w-full flex-col px-4 py-2 text-left transition-colors
                {i === selectedIndex ? 'bg-primary/10 text-text-main' : 'text-text-muted hover:bg-bg-card-hover'}"
              onclick={() => navigateTo(result.path)}
              onmouseenter={() => (selectedIndex = i)}
            >
              <span class="text-sm font-medium">{result.title}</span>
              <span class="text-xs text-text-dim">{result.path}</span>
              {#if result.snippet}
                <span class="mt-1 line-clamp-2 text-xs text-text-dim">{result.snippet}</span>
              {/if}
            </button>
          {/each}
        </div>
      {:else if query.trim() && !loading}
        <div class="px-4 py-6 text-center text-sm text-text-muted">
          {#if mode === "switcher"}
            Enterで「{query}」を作成
          {:else}
            見つかりませんでした
          {/if}
        </div>
      {/if}

      <div class="flex items-center justify-between border-t border-border px-4 py-1.5 text-[10px] text-text-dim">
        <span>
          <kbd class="rounded border border-border bg-bg-dark px-1 py-0.5">↑↓</kbd> 移動
          <kbd class="ml-1.5 rounded border border-border bg-bg-dark px-1 py-0.5">↵</kbd> 開く
        </span>
        <span>
          <kbd class="rounded border border-border bg-bg-dark px-1 py-0.5">esc</kbd> 閉じる
        </span>
      </div>
    </div>
  </div>
{/if}
