<script lang="ts">
  import { Link2 } from "lucide-svelte";

  interface Backlink {
    sourceId: string;
    sourcePath: string;
    sourceTitle: string;
    context: string | null;
  }

  interface Props {
    notePath: string;
  }

  let { notePath }: Props = $props();
  let backlinks: Backlink[] = $state([]);
  let loading = $state(true);

  $effect(() => {
    loadBacklinks();
  });

  async function loadBacklinks() {
    loading = true;
    try {
      const res = await fetch(`/api/backlinks/${notePath}`);
      backlinks = await res.json();
    } catch {
      backlinks = [];
    } finally {
      loading = false;
    }
  }
</script>

<div class="shrink-0 border-t border-border px-5 py-2.5">
  {#if loading}
    <p class="text-[11px] text-text-dim">バックリンクを読み込み中...</p>
  {:else if backlinks.length > 0}
    <h3 class="mb-1.5 flex items-center gap-1.5 text-[11px] font-medium uppercase tracking-wider text-text-dim">
      <Link2 size={11} />
      バックリンク ({backlinks.length})
    </h3>
    <div class="flex flex-wrap gap-1.5">
      {#each backlinks as bl}
        <a
          href="/note/{bl.sourcePath}"
          class="inline-flex items-center rounded-md border border-border bg-bg-card px-2.5 py-1 text-xs text-text-muted transition-colors hover:border-primary/50 hover:text-primary-light"
        >
          {bl.sourceTitle}
        </a>
      {/each}
    </div>
  {:else}
    <p class="text-[11px] text-text-dim">バックリンクなし</p>
  {/if}
</div>
