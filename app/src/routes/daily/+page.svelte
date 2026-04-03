<script lang="ts">
  import EditorLayout from "$lib/components/editor/EditorLayout.svelte";
  import BacklinkPanel from "$lib/components/editor/BacklinkPanel.svelte";

  let { data } = $props();

  let previewHtml = $state("");
  let saving = $state(false);
  let lastSaved = $state<string | null>(null);
  let debounceTimer: ReturnType<typeof setTimeout> | undefined;

  $effect(() => {
    previewHtml = data.html;
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
    const current = new Date(data.currentDate);
    current.setDate(current.getDate() + offset);
    const dateStr = current.toISOString().slice(0, 10);
    window.location.href = `/daily?date=${dateStr}`;
  }
</script>

<svelte:head>
  <title>Daily - {data.currentDate} - Prometheus</title>
</svelte:head>

<div class="flex h-full flex-col -m-6">
  <!-- Daily header -->
  <div class="flex items-center gap-3 border-b border-border px-4 py-2">
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
    <span class="text-xs text-text-dim">{data.path}</span>
    <div class="flex-1"></div>
    {#if saving}
      <span class="text-xs text-text-muted">Saving...</span>
    {:else if lastSaved}
      <span class="text-xs text-text-dim">Saved {lastSaved}</span>
    {/if}

    <!-- Recent dates -->
    <div class="hidden items-center gap-1 md:flex">
      {#each data.recentNotes.slice(0, 7) as recent}
        <a
          href="/daily?date={recent.date}"
          class="h-2 w-2 rounded-full transition-colors {recent.date === data.currentDate
            ? 'bg-primary'
            : recent.exists
              ? 'bg-primary/40 hover:bg-primary/60'
              : 'bg-border hover:bg-text-dim'}"
          title={recent.date}
        ></a>
      {/each}
    </div>
  </div>

  <!-- Editor -->
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
