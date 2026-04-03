<script lang="ts">
  import { ChevronRight, ChevronDown, Tag } from "lucide-svelte";

  interface TagItem {
    id: number;
    name: string;
    count: number;
  }

  let tags: TagItem[] = $state([]);
  let expanded = $state(false);

  $effect(() => {
    loadTags();
  });

  async function loadTags() {
    try {
      const res = await fetch("/api/tags");
      tags = await res.json();
    } catch {
      tags = [];
    }
  }
</script>

{#if tags.length > 0}
  <div class="border-t border-border px-2 py-2">
    <button
      onclick={() => (expanded = !expanded)}
      class="mb-1 flex w-full items-center gap-1.5 px-3 text-[11px] font-medium uppercase tracking-wider text-text-dim"
    >
      {#if expanded}
        <ChevronDown size={11} />
      {:else}
        <ChevronRight size={11} />
      {/if}
      タグ
    </button>

    {#if expanded}
      <div class="flex flex-wrap gap-1 px-3 pt-1">
        {#each tags as tag}
          <a
            href="/search?q=%23{tag.name}"
            class="inline-flex items-center gap-1 rounded-full border border-border bg-bg-dark px-2 py-0.5 text-[11px] text-primary-light transition-colors hover:border-primary hover:bg-primary/10"
          >
            <Tag size={10} />
            {tag.name}
            <span class="text-text-dim">{tag.count}</span>
          </a>
        {/each}
      </div>
    {/if}
  </div>
{/if}
