<script lang="ts">
  import { X, Plus, Tag } from "lucide-svelte";

  interface Props {
    tags: string[];
    notePath: string;
    frontmatter: Record<string, unknown>;
  }

  let { tags, notePath, frontmatter }: Props = $props();

  let showInput = $state(false);
  let newTag = $state("");
  let inputEl: HTMLInputElement | undefined = $state(undefined);
  let localTags = $state<string[]>([]);

  // Sync from props
  $effect(() => {
    localTags = [...tags];
  });

  async function saveTagsToNote(updatedTags: string[]) {
    try {
      const res = await fetch(`/api/notes/${notePath}`);
      if (!res.ok) return;
      const note = await res.json();

      await fetch(`/api/notes/${notePath}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: note.content,
          frontmatter: { ...note.frontmatter, tags: updatedTags },
        }),
      });
    } catch (e) {
      console.error("Failed to save tags:", e);
    }
  }

  async function addTag() {
    const tag = newTag.trim().toLowerCase().replace(/\s+/g, "-");
    if (!tag || localTags.includes(tag)) {
      newTag = "";
      return;
    }

    localTags = [...localTags, tag];
    newTag = "";
    showInput = false;
    await saveTagsToNote(localTags);
  }

  async function removeTag(tag: string) {
    localTags = localTags.filter((t) => t !== tag);
    await saveTagsToNote(localTags);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();
      addTag();
    } else if (e.key === "Escape") {
      showInput = false;
      newTag = "";
    }
  }
</script>

<div class="flex flex-wrap items-center gap-1.5">
  {#each localTags as tag}
    <span class="group inline-flex items-center gap-1 rounded-full border border-border bg-bg-dark px-2 py-0.5 text-[10px] text-primary-light">
      <Tag size={9} />
      {tag}
      <button
        onclick={() => removeTag(tag)}
        class="hidden rounded-full p-0.5 text-text-dim transition-colors hover:bg-error/20 hover:text-error group-hover:inline-flex"
        title="タグを削除"
      >
        <X size={8} />
      </button>
    </span>
  {/each}

  {#if showInput}
    <form onsubmit={(e) => { e.preventDefault(); addTag(); }} class="inline-flex">
      <input
        bind:this={inputEl}
        bind:value={newTag}
        onkeydown={handleKeydown}
        onblur={() => { if (!newTag.trim()) showInput = false; }}
        placeholder="タグ名"
        class="w-20 rounded border border-primary bg-bg-input px-1.5 py-0.5 text-[10px] text-text-main outline-none placeholder:text-text-dim"
      />
    </form>
  {:else}
    <button
      onclick={() => { showInput = true; setTimeout(() => inputEl?.focus(), 50); }}
      class="inline-flex items-center rounded-full border border-dashed border-border px-1.5 py-0.5 text-[10px] text-text-dim transition-colors hover:border-primary hover:text-primary-light"
      title="タグを追加"
    >
      <Plus size={10} />
    </button>
  {/if}
</div>
