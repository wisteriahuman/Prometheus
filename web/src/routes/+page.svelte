<script lang="ts">
  import { onMount } from "svelte";
  import { CalendarDays, Network, CheckSquare, Search, Command, FileText, PanelLeft, Keyboard } from "lucide-svelte";

  let vaultPath = $state("");

  onMount(async () => {
    try {
      const res = await fetch("/api/config");
      const config = await res.json();
      vaultPath = config.vaultPath;
    } catch {}
  });

  const today = new Date().toLocaleDateString("ja-JP", {
    year: "numeric",
    month: "long",
    day: "numeric",
    weekday: "long",
  });
</script>

<div class="mx-auto max-w-2xl py-4">
  <div class="mb-10">
    <h1 class="font-mono text-3xl font-bold tracking-tight text-text-main">Prometheus</h1>
    <p class="mt-2 text-sm text-text-muted">{today}</p>
  </div>

  <div class="grid gap-3 sm:grid-cols-2">
    <a
      href="/daily"
      class="group flex items-start gap-4 rounded-xl border border-border bg-bg-card p-5 transition-all hover:border-primary/50 hover:bg-bg-card-hover"
    >
      <div class="rounded-lg bg-primary/10 p-2.5 text-primary transition-colors group-hover:bg-primary/20">
        <CalendarDays size={20} />
      </div>
      <div>
        <h2 class="text-sm font-semibold text-text-main">デイリーノート</h2>
        <p class="mt-0.5 text-xs text-text-muted">今日のノートを開く</p>
      </div>
    </a>

    <a
      href="/graph"
      class="group flex items-start gap-4 rounded-xl border border-border bg-bg-card p-5 transition-all hover:border-primary/50 hover:bg-bg-card-hover"
    >
      <div class="rounded-lg bg-primary/10 p-2.5 text-primary transition-colors group-hover:bg-primary/20">
        <Network size={20} />
      </div>
      <div>
        <h2 class="text-sm font-semibold text-text-main">グラフビュー</h2>
        <p class="mt-0.5 text-xs text-text-muted">ノートの関係を可視化</p>
      </div>
    </a>

    <a
      href="/tasks"
      class="group flex items-start gap-4 rounded-xl border border-border bg-bg-card p-5 transition-all hover:border-primary/50 hover:bg-bg-card-hover"
    >
      <div class="rounded-lg bg-primary/10 p-2.5 text-primary transition-colors group-hover:bg-primary/20">
        <CheckSquare size={20} />
      </div>
      <div>
        <h2 class="text-sm font-semibold text-text-main">タスク</h2>
        <p class="mt-0.5 text-xs text-text-muted">タスク一覧を確認</p>
      </div>
    </a>

    <a
      href="/search"
      class="group flex items-start gap-4 rounded-xl border border-border bg-bg-card p-5 transition-all hover:border-primary/50 hover:bg-bg-card-hover"
    >
      <div class="rounded-lg bg-primary/10 p-2.5 text-primary transition-colors group-hover:bg-primary/20">
        <Search size={20} />
      </div>
      <div>
        <h2 class="text-sm font-semibold text-text-main">検索</h2>
        <p class="mt-0.5 text-xs text-text-muted">全文検索</p>
      </div>
    </a>
  </div>

  <div class="mt-8 rounded-xl border border-border bg-bg-card p-5">
    <h2 class="mb-4 flex items-center gap-2 text-sm font-semibold text-text-main">
      <Keyboard size={15} />
      ショートカット
    </h2>
    <div class="grid gap-2 sm:grid-cols-2">
      {#each [
        { key: "⌘K", label: "コマンドパレット" },
        { key: "⌘P", label: "ファイル切替" },
        { key: "⌘⇧F", label: "全文検索" },
        { key: "⌘\\", label: "サイドバー切替" },
        { key: "⌘D", label: "デイリーノート" },
        { key: "⌘G", label: "グラフビュー" },
        { key: "⌘.", label: "テーマ切替" },
        { key: "⌘K → ?", label: "ショートカット一覧" },
      ] as item}
        <div class="flex items-center gap-3">
          <kbd class="inline-flex min-w-[2.5rem] items-center justify-center rounded border border-border bg-bg-dark px-1.5 py-0.5 font-mono text-[11px] text-text-dim">
            {item.key}
          </kbd>
          <span class="text-xs text-text-muted">{item.label}</span>
        </div>
      {/each}
    </div>
  </div>

  <div class="mt-6 rounded-xl border border-border bg-bg-card p-5">
    <h2 class="mb-3 flex items-center gap-2 text-sm font-semibold text-text-main">
      ノートの保存先
    </h2>
    <div class="space-y-2 text-xs text-text-muted leading-relaxed">
      <p>ノートは <code class="rounded bg-bg-dark px-1.5 py-0.5 font-mono text-accent">{vaultPath}</code> に <code class="rounded bg-bg-dark px-1.5 py-0.5 font-mono text-accent">.md</code> ファイルとして保存されます。</p>
      <p>Neovimや好きなエディタで直接編集可能です。変更は自動的にPrometheusに反映されます。</p>
      <p>右クリックでノートの<strong class="text-text-main">リネーム</strong>・<strong class="text-text-main">削除</strong>ができます。</p>
    </div>
  </div>
</div>
