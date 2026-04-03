<script lang="ts">
  interface Props {
    open: boolean;
    onclose: () => void;
  }

  let { open, onclose }: Props = $props();

  const shortcuts = [
    { category: "General", items: [
      { key: "⌘K", desc: "コマンドパレット" },
      { key: "⌘P", desc: "ファイル切替" },
      { key: "⌘⇧F", desc: "グローバル検索" },
      { key: "⌘\\", desc: "サイドバー切替" },
      { key: "⌘K → ?", desc: "ショートカット一覧" },
    ]},
    { category: "Navigation", items: [
      { key: "⌘D", desc: "デイリーノート" },
      { key: "⌘G", desc: "グラフビュー" },
      { key: "⌘.", desc: "テーマ切替" },
    ]},
    { category: "Editor", items: [
      { key: ":w / ⌘S", desc: "保存" },
      { key: "⌘↵", desc: "エディタ/プレビュー切替" },
      { key: "Vim", desc: "エディタ内でVimキーバインド" },
    ]},
  ];

  function handleKeydown(e: KeyboardEvent) {
    if (!open) return;
    if (e.key === "Escape" || e.key === "?") {
      e.preventDefault();
      onclose();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <button
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={onclose}
      aria-label="Close"
    ></button>

    <div
      class="relative z-10 w-full max-w-md overflow-hidden rounded-xl border border-border bg-bg-card shadow-2xl"
      role="dialog"
    >
      <div class="flex items-center justify-between border-b border-border px-5 py-3">
        <h2 class="font-mono text-sm font-semibold text-text-main">Keyboard Shortcuts</h2>
        <button
          onclick={onclose}
          class="text-text-dim hover:text-text-main"
        >✕</button>
      </div>

      <div class="max-h-96 overflow-y-auto p-5">
        {#each shortcuts as group}
          <div class="mb-4 last:mb-0">
            <h3 class="mb-2 text-xs font-medium uppercase tracking-wider text-text-dim">
              {group.category}
            </h3>
            <div class="space-y-1.5">
              {#each group.items as shortcut}
                <div class="flex items-center justify-between">
                  <span class="text-sm text-text-muted">{shortcut.desc}</span>
                  <kbd class="rounded border border-border bg-bg-dark px-2 py-0.5 font-mono text-xs text-text-dim">
                    {shortcut.key}
                  </kbd>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}
