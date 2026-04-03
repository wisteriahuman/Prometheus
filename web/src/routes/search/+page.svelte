<script lang="ts">
  import { Search } from "lucide-svelte";

  interface SearchResult {
    id: string;
    path: string;
    title: string;
    snippet: string;
  }

  let query = $state("");
  let results: SearchResult[] = $state([]);
  let loading = $state(false);
  let debounceTimer: ReturnType<typeof setTimeout> | undefined;

  function handleInput() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => search(), 200);
  }

  async function search() {
    if (!query.trim()) {
      results = [];
      return;
    }

    loading = true;
    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
      results = await res.json();
    } catch {
      results = [];
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>検索 - Prometheus</title>
</svelte:head>

<div class="mx-auto max-w-2xl py-4">
  <h1 class="mb-5 font-mono text-2xl font-bold text-text-main">検索</h1>

  <div class="relative mb-6">
    <Search size={15} class="absolute left-3.5 top-1/2 -translate-y-1/2 text-text-dim" />
    <input
      bind:value={query}
      oninput={handleInput}
      placeholder="ノートを検索..."
      class="w-full rounded-lg border border-border bg-bg-card py-2.5 pl-10 pr-4 text-sm text-text-main outline-none placeholder:text-text-dim focus:border-primary"
    />
  </div>

  {#if loading}
    <p class="text-sm text-text-muted">検索中...</p>
  {:else if results.length > 0}
    <div class="space-y-2">
      {#each results as result}
        <a
          href="/note/{result.path}"
          class="block rounded-lg border border-border bg-bg-card p-4 transition-colors hover:border-primary/50 hover:bg-bg-card-hover"
        >
          <h2 class="text-sm font-medium text-text-main">{result.title}</h2>
          <p class="mt-0.5 text-xs text-text-dim">{result.path}</p>
          {#if result.snippet}
            <p class="mt-2 text-xs text-text-muted leading-relaxed">{result.snippet}</p>
          {/if}
        </a>
      {/each}
    </div>
  {:else if query.trim()}
    <p class="text-sm text-text-muted">「{query}」に一致するノートはありません</p>
  {/if}
</div>
