<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import EditorLayout from "$lib/components/editor/EditorLayout.svelte";
  import BacklinkPanel from "$lib/components/editor/BacklinkPanel.svelte";

  interface DailyData {
    path: string;
    title: string;
    content: string;
    frontmatter: Record<string, unknown>;
    html: string;
    currentDate: string;
    recentNotes: { date: string; exists: boolean }[];
  }

  let data = $state<DailyData | null>(null);
  let loading = $state(true);
  let previewHtml = $state("");
  let saving = $state(false);
  let lastSaved = $state<string | null>(null);
  let debounceTimer: ReturnType<typeof setTimeout> | undefined;

  async function loadDaily(dateStr?: string) {
    loading = true;
    try {
      const url = dateStr ? `/api/daily?date=${dateStr}` : "/api/daily";
      const res = await fetch(url);
      data = await res.json();
      if (data) previewHtml = data.html;
    } catch (e) {
      console.error("Failed to load daily:", e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    const dateParam = $page.url.searchParams.get("date");
    loadDaily(dateParam ?? undefined);
  });

  async function updatePreview(newContent: string) {
    try {
      const res = await fetch("/api/preview", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content: newContent }),
      });
      const result = await res.json();
      previewHtml = result.html;
    } catch {}
  }

  function handleChange(newContent: string) {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => updatePreview(newContent), 300);
  }

  async function handleSave(saveContent: string) {
    if (!data) return;
    saving = true;
    try {
      const res = await fetch(`/api/notes/${data.path}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: saveContent,
          frontmatter: data.frontmatter,
        }),
      });
      if (res.ok) {
        lastSaved = new Date().toLocaleTimeString("ja-JP");
      }
    } catch (e) {
      console.error("Failed to save:", e);
    } finally {
      saving = false;
    }
  }

  function navigateDate(offset: number) {
    if (!data) return;
    const current = new Date(data.currentDate);
    current.setDate(current.getDate() + offset);
    const dateStr = current.toISOString().slice(0, 10);
    loadDaily(dateStr);
  }
</script>

<svelte:head>
  <title>デイリー{data ? ` - ${data.currentDate}` : ""} - Prometheus</title>
</svelte:head>

{#if loading}
  <div class="flex h-full items-center justify-center">
    <span class="text-text-muted">読み込み中...</span>
  </div>
{:else if data}
  <div class="flex h-full flex-col -m-5">
    <div class="flex items-center gap-3 border-b border-border px-5 py-2.5">
      <button
        onclick={() => navigateDate(-1)}
        class="rounded-md px-2 py-1 text-text-muted transition-colors hover:bg-bg-card-hover hover:text-text-main"
      >
        ←
      </button>
      <h1 class="font-mono text-sm font-medium text-text-main">{data.currentDate}</h1>
      <button
        onclick={() => navigateDate(1)}
        class="rounded-md px-2 py-1 text-text-muted transition-colors hover:bg-bg-card-hover hover:text-text-main"
      >
        →
      </button>
      <span class="text-[11px] text-text-dim">{data.path}</span>
      <div class="flex-1"></div>
      {#if saving}
        <span class="text-[11px] text-text-muted">保存中...</span>
      {:else if lastSaved}
        <span class="text-[11px] text-text-dim">保存済み {lastSaved}</span>
      {/if}

      <div class="hidden items-center gap-1 md:flex">
        {#each data.recentNotes.slice(0, 7) as recent}
          <button
            onclick={() => loadDaily(recent.date)}
            class="h-2 w-2 rounded-full transition-colors {recent.date === data?.currentDate
              ? 'bg-primary'
              : recent.exists
                ? 'bg-primary/40 hover:bg-primary/60'
                : 'bg-border hover:bg-text-dim'}"
            title={recent.date}
          ></button>
        {/each}
      </div>
    </div>

    <div class="flex-1 overflow-hidden">
      <EditorLayout
        content={data.content}
        {previewHtml}
        vimMode={true}
        onchange={handleChange}
        onsave={handleSave}
      />
    </div>

    <BacklinkPanel notePath={data.path} />
  </div>
{/if}
