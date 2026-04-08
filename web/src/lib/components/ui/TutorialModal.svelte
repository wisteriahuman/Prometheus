<script lang="ts">
  import { X, FileText, Network, CalendarDays, Search, Palette, Command, Download, Keyboard } from "lucide-svelte";

  interface Props {
    open: boolean;
    onclose: () => void;
  }

  let { open, onclose }: Props = $props();

  let step = $state(0);

  const steps = [
    {
      title: "Prometheusへようこそ",
      icon: FileText,
      content: `Markdownベースの「第二の脳」ノートアプリです。
ノートは .md ファイルとして保存され、Neovimや好きなエディタで直接編集できます。`,
    },
    {
      title: "ノートの操作",
      icon: FileText,
      content: `• サイドバーの + からノートを作成
• [[wikilink]] でノート同士をリンク
• タイトルクリックで名前を変更
• 右クリックでリネーム・削除
• :w または ⌘S で保存`,
    },
    {
      title: "グラフビュー",
      icon: Network,
      content: `⌘G でノート間の関係をグラフとして可視化。
タグごとにクラスタ化され、フィルタリングも可能です。`,
    },
    {
      title: "デイリーノート",
      icon: CalendarDays,
      content: `⌘D で今日のノートを即座に開けます。
テンプレートは .prometheus/config.json でカスタマイズ可能。`,
    },
    {
      title: "検索とコマンド",
      icon: Search,
      content: `⌘P でファイル切替
⌘⇧F で全文検索
⌘K でコマンドパレット（全機能にアクセス）`,
    },
    {
      title: "テーマ",
      icon: Palette,
      content: `⌘. または コマンドパレットからテーマを変更。
13種のビルトインテーマ + カスタムテーマ作成が可能。
ノートごとに個別テーマも設定できます。`,
    },
    {
      title: "エクスポート",
      icon: Download,
      content: `ノートヘッダーの↓アイコンからエクスポート:
• HTML（テーマ付き、共有向け）
• Markdown（フロントマターなし）
• PDF（ブラウザ印刷）`,
    },
    {
      title: "ショートカット",
      icon: Keyboard,
      content: `⌘K  コマンドパレット
⌘P  ファイル切替
⌘D  デイリーノート
⌘G  グラフビュー
⌘.  テーマ切替
⌘\\  サイドバー切替
:w   保存（Vimモード）`,
    },
    {
      title: "AI / Agent連携",
      icon: Command,
      content: `vaultディレクトリでClaude Code や agent 対応ツールを起動すると、
ノートの構造を理解した上でAIが支援します。

cd ~/my-notes && claude
cd ~/my-notes && codex

生成される設定:
• CLAUDE.md / .claude/
• AGENTS.md / .agents/

.mdファイルはそのままAIへの入力になります。`,
    },
  ];

  function next() {
    if (step < steps.length - 1) {
      step++;
    } else {
      finish();
    }
  }

  function prev() {
    if (step > 0) step--;
  }

  async function finish() {
    // Mark tutorial as shown
    try {
      await fetch("/api/config", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ tutorialShown: true }),
      });
    } catch {}
    onclose();
  }

  $effect(() => {
    if (open) step = 0;
  });
</script>

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <button class="absolute inset-0 bg-black/60 backdrop-blur-sm" onclick={finish} aria-label="閉じる"></button>

    <div class="relative z-10 w-full max-w-md overflow-hidden rounded-xl border border-border bg-bg-card shadow-2xl">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-border px-5 py-3">
        <div class="flex items-center gap-2">
          {#each [steps[step]] as s}
            {@const Icon = s.icon}
            <Icon size={16} class="text-primary" />
            <h2 class="text-sm font-semibold text-text-main">{s.title}</h2>
          {/each}
        </div>
        <button onclick={finish} class="text-text-dim hover:text-text-main">
          <X size={16} />
        </button>
      </div>

      <!-- Content -->
      <div class="px-5 py-5">
        <pre class="whitespace-pre-wrap font-sans text-sm leading-relaxed text-text-muted">{steps[step].content}</pre>
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between border-t border-border px-5 py-3">
        <!-- Progress dots -->
        <div class="flex gap-1.5">
          {#each steps as _, i}
            <button
              onclick={() => (step = i)}
              aria-label={`${i + 1}番目のステップへ移動`}
              title={`${i + 1} / ${steps.length}`}
              class="h-1.5 w-1.5 rounded-full transition-colors {i === step ? 'bg-primary' : 'bg-border hover:bg-text-dim'}"
            ></button>
          {/each}
        </div>

        <div class="flex gap-2">
          {#if step > 0}
            <button
              onclick={prev}
              class="rounded-md border border-border px-3 py-1.5 text-xs text-text-muted hover:text-text-main"
            >
              戻る
            </button>
          {/if}
          <button
            onclick={next}
            class="rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-white hover:bg-primary-dark"
          >
            {step < steps.length - 1 ? "次へ" : "はじめる"}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
