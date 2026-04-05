<script lang="ts">
  import { goto } from "$app/navigation";
  import { toggleSidebar } from "$lib/stores/sidebar";
  import { getAllThemes, setTheme } from "$lib/stores/theme";
  import { CalendarDays, Network, Search, FileText, PanelLeft, Palette, CheckSquare, Home, Keyboard, ArrowLeft } from "lucide-svelte";

  interface Command {
    id: string;
    title: string;
    description?: string;
    keybinding?: string;
    category: string;
    icon?: typeof Home;
    action: () => void;
  }

  interface Props {
    open: boolean;
    onclose: () => void;
  }

  let { open, onclose }: Props = $props();

  let query = $state("");
  let selectedIndex = $state(0);
  let inputEl: HTMLInputElement | undefined = $state(undefined);
  let mode = $state<"commands" | "themes">("commands");

  const commands: Command[] = [
    {
      id: "daily",
      title: "デイリーノートを開く",
      description: "今日のノートへ移動",
      keybinding: "⌘D",
      category: "移動",
      icon: CalendarDays,
      action: () => {
        const today = new Date().toISOString().slice(0, 10);
        goto(`/note/daily/${today}.md`);
      },
    },
    {
      id: "graph",
      title: "グラフビューを開く",
      description: "ノートの関係を可視化",
      keybinding: "⌘G",
      category: "移動",
      icon: Network,
      action: () => goto("/graph"),
    },
    {
      id: "search",
      title: "全文検索",
      description: "ノート内を検索",
      keybinding: "⌘⇧F",
      category: "検索",
      icon: Search,
      action: () => {
        onclose();
        window.dispatchEvent(new CustomEvent("prometheus:search"));
      },
    },
    {
      id: "switcher",
      title: "ファイルを開く",
      description: "ファイルを素早く切替",
      keybinding: "⌘P",
      category: "移動",
      icon: FileText,
      action: () => {
        onclose();
        window.dispatchEvent(new CustomEvent("prometheus:switcher"));
      },
    },
    {
      id: "sidebar",
      title: "サイドバー切替",
      keybinding: "⌘\\",
      category: "表示",
      icon: PanelLeft,
      action: () => toggleSidebar(),
    },
    {
      id: "theme",
      title: "テーマを変更",
      description: "カラーテーマを切替",
      keybinding: "⌘.",
      category: "外観",
      icon: Palette,
      action: () => {
        mode = "themes";
        query = "";
        selectedIndex = 0;
      },
    },
    {
      id: "tasks",
      title: "タスク一覧",
      description: "タスクを確認",
      category: "移動",
      icon: CheckSquare,
      action: () => goto("/tasks"),
    },
    {
      id: "home",
      title: "ホームへ戻る",
      category: "移動",
      icon: Home,
      action: () => goto("/"),
    },
    {
      id: "shortcuts",
      title: "ショートカット一覧",
      category: "ヘルプ",
      icon: Keyboard,
      action: () => {
        onclose();
        window.dispatchEvent(new CustomEvent("prometheus:shortcuts"));
      },
    },
    {
      id: "create-theme",
      title: "カスタムテーマを作成",
      description: "自分だけのテーマを作る",
      category: "外観",
      icon: Palette,
      action: () => {
        onclose();
        window.dispatchEvent(new CustomEvent("prometheus:create-theme"));
      },
    },
    {
      id: "tutorial",
      title: "チュートリアル",
      description: "使い方ガイドを表示",
      category: "ヘルプ",
      icon: Keyboard,
      action: () => {
        onclose();
        window.dispatchEvent(new CustomEvent("prometheus:tutorial"));
      },
    },
  ];

  let themeCommands = $derived(
    getAllThemes().map((t) => ({
      id: `theme-${t.slug}`,
      title: t.name,
      description: t.slug,
      color: t.colors.primary,
      category: "テーマ",
      action: () => setTheme(t.slug),
    })),
  );

  let filteredCommands = $derived(
    query.trim()
      ? (mode === "themes" ? themeCommands : commands).filter(
          (c) =>
            c.title.toLowerCase().includes(query.toLowerCase()) ||
            (c.description?.toLowerCase().includes(query.toLowerCase()) ?? false),
        )
      : mode === "themes"
        ? themeCommands
        : commands,
  );

  $effect(() => {
    if (open) {
      query = "";
      selectedIndex = 0;
      mode = "commands";
      setTimeout(() => inputEl?.focus(), 50);
    }
  });

  // Listen for theme-mode event (from Cmd+. shortcut)
  import { onMount } from "svelte";
  onMount(() => {
    const handleThemeMode = () => {
      if (open) {
        mode = "themes";
        query = "";
        selectedIndex = 0;
      }
    };
    window.addEventListener("prometheus:theme-mode", handleThemeMode);
    return () => window.removeEventListener("prometheus:theme-mode", handleThemeMode);
  });

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") {
      e.preventDefault();
      if (mode === "themes") {
        mode = "commands";
        query = "";
        selectedIndex = 0;
      } else {
        onclose();
      }
    } else if (e.key === "ArrowDown" || (e.ctrlKey && e.key === "n")) {
      e.preventDefault();
      selectedIndex = Math.min(selectedIndex + 1, filteredCommands.length - 1);
    } else if (e.key === "ArrowUp" || (e.ctrlKey && e.key === "p")) {
      e.preventDefault();
      selectedIndex = Math.max(selectedIndex - 1, 0);
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (filteredCommands[selectedIndex]) {
        executeCommand(filteredCommands[selectedIndex]);
      }
    }
  }

  function executeCommand(cmd: { action: () => void; id?: string }) {
    if (mode === "themes") {
      // Apply theme and close
      cmd.action();
      onclose();
    } else if ("id" in cmd && cmd.id === "theme") {
      // Switch to theme sub-menu without closing
      cmd.action();
    } else {
      onclose();
      cmd.action();
    }
  }
</script>

{#if open}
  <div class="fixed inset-0 z-50 flex items-start justify-center pt-[20vh]">
    <button
      class="absolute inset-0 bg-black/60 backdrop-blur-sm"
      onclick={onclose}
      aria-label="閉じる"
    ></button>

    <div class="relative z-10 w-full max-w-md overflow-hidden rounded-xl border border-border bg-bg-card shadow-2xl">
      <!-- Input -->
      <div class="flex items-center gap-2 border-b border-border px-4">
        {#if mode === "themes"}
          <button
            onclick={() => { mode = "commands"; query = ""; selectedIndex = 0; }}
            class="text-text-dim hover:text-text-muted"
          >
            <ArrowLeft size={14} />
          </button>
        {/if}
        <input
          bind:this={inputEl}
          bind:value={query}
          onkeydown={handleKeydown}
          placeholder={mode === "themes" ? "テーマを選択..." : "コマンドを入力..."}
          class="flex-1 bg-transparent py-3 text-sm text-text-main outline-none placeholder:text-text-dim"
        />
      </div>

      <!-- Results -->
      <div class="max-h-72 overflow-y-auto py-1">
        {#each filteredCommands as cmd, i}
          <button
            class="flex w-full items-center gap-3 px-4 py-2 text-left transition-colors
              {i === selectedIndex ? 'bg-primary/10 text-text-main' : 'text-text-muted hover:bg-bg-card-hover'}"
            onclick={() => executeCommand(cmd)}
            onmouseenter={() => (selectedIndex = i)}
          >
            {#if mode === "themes" && "color" in cmd}
              <span class="h-3 w-3 shrink-0 rounded-full" style="background-color: {cmd.color}"></span>
            {:else if "icon" in cmd && cmd.icon}
              {@const Icon = cmd.icon}
              <Icon size={15} class="shrink-0" />
            {/if}
            <div class="flex-1 min-w-0">
              <span class="text-sm">{cmd.title}</span>
              {#if mode !== "themes" && cmd.description}
                <span class="ml-2 text-xs text-text-dim">{cmd.description}</span>
              {/if}
            </div>
            {#if "keybinding" in cmd && cmd.keybinding}
              <kbd class="shrink-0 rounded border border-border bg-bg-dark px-1.5 py-0.5 font-mono text-[10px] text-text-dim">
                {cmd.keybinding}
              </kbd>
            {/if}
          </button>
        {/each}
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between border-t border-border px-4 py-1.5 text-[10px] text-text-dim">
        <span>
          <kbd class="rounded border border-border bg-bg-dark px-1 py-0.5">↑↓</kbd> 移動
          <kbd class="ml-1.5 rounded border border-border bg-bg-dark px-1 py-0.5">↵</kbd> {mode === "themes" ? "適用" : "実行"}
        </span>
        <span>
          <kbd class="rounded border border-border bg-bg-dark px-1 py-0.5">esc</kbd> {mode === "themes" ? "戻る" : "閉じる"}
        </span>
      </div>
    </div>
  </div>
{/if}
