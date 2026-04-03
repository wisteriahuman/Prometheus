<script lang="ts">
  import Editor from "./Editor.svelte";
  import Preview from "./Preview.svelte";
  import { PenLine, Columns2, Eye, Save } from "lucide-svelte";

  interface Props {
    content: string;
    previewHtml: string;
    vimMode?: boolean;
    onchange?: (content: string) => void;
    onsave?: (content: string) => void;
  }

  let { content, previewHtml, vimMode = true, onchange, onsave }: Props = $props();

  type ViewMode = "split" | "editor" | "preview";
  let viewMode: ViewMode = $state("split");

  function cycleView() {
    const modes: ViewMode[] = ["split", "editor", "preview"];
    const idx = modes.indexOf(viewMode);
    viewMode = modes[(idx + 1) % modes.length];
  }

  function handleKeydown(e: KeyboardEvent) {
    if ((e.metaKey || e.ctrlKey) && e.key === "Enter") {
      e.preventDefault();
      cycleView();
    }
  }

  function handleSave() {
    // Dispatch event to let Editor component handle save with current content
    window.dispatchEvent(new CustomEvent("prometheus:save"));
  }

  const VIEW_TABS = [
    { mode: "editor" as const, icon: PenLine, label: "編集" },
    { mode: "split" as const, icon: Columns2, label: "分割" },
    { mode: "preview" as const, icon: Eye, label: "プレビュー" },
  ];
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="flex h-full flex-col">
  <!-- Toolbar -->
  <div class="flex h-9 shrink-0 items-center gap-2 border-b border-border bg-bg-card px-3">
    <div class="flex rounded-md border border-border">
      {#each VIEW_TABS as tab}
        <button
          onclick={() => (viewMode = tab.mode)}
          class="flex items-center gap-1 px-2 py-1 text-[11px] transition-colors
            {tab.mode !== 'editor' ? 'border-l border-border' : ''}
            {viewMode === tab.mode
              ? 'bg-primary text-white'
              : 'text-text-muted hover:text-text-main'}"
        >
          <tab.icon size={12} />
          <span class="hidden sm:inline">{tab.label}</span>
        </button>
      {/each}
    </div>

    <div class="flex-1"></div>

    <div class="flex items-center gap-2 text-[10px] text-text-dim">
      <button
        onclick={handleSave}
        class="flex items-center gap-1 rounded border border-border px-1.5 py-0.5 transition-colors hover:border-primary hover:text-primary-light"
        title="保存 (:w / ⌘S)"
      >
        <Save size={10} />
        保存
      </button>
      {#if vimMode}
        <span class="rounded border border-border bg-bg-dark px-1.5 py-0.5 font-mono">VIM</span>
      {/if}
      <span>
        <kbd class="rounded border border-border bg-bg-dark px-1 py-0.5 font-mono">⌘↵</kbd>
        表示切替
      </span>
    </div>
  </div>

  <!-- Content -->
  <div class="flex flex-1 overflow-hidden bg-bg-dark">
    {#if viewMode === "editor"}
      <div class="h-full w-full">
        <Editor {content} {vimMode} {onchange} {onsave} />
      </div>
    {:else if viewMode === "split"}
      <div class="h-full w-1/2 border-r border-border">
        <Editor {content} {vimMode} {onchange} {onsave} />
      </div>
      <div class="h-full w-1/2">
        <Preview html={previewHtml} />
      </div>
    {:else}
      <div class="flex h-full w-full justify-center">
        <div class="h-full w-full max-w-3xl">
          <Preview html={previewHtml} />
        </div>
      </div>
    {/if}
  </div>
</div>
