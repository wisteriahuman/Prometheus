<script lang="ts">
  import { addCustomTheme, type PrometheusTheme, type ThemeColors, applyTheme, setTheme } from "$lib/stores/theme";
  import { Palette, Save, X, RotateCcw } from "lucide-svelte";

  interface Props {
    open: boolean;
    onclose: () => void;
  }

  let { open, onclose }: Props = $props();

  const defaultColors: ThemeColors = {
    bgDark: "#0f172a",
    bgCard: "#1e293b",
    bgCardHover: "#334155",
    primary: "#6366f1",
    primaryLight: "#818cf8",
    primaryDark: "#4f46e5",
    accent: "#38bdf8",
    accentLight: "#7dd3fc",
    textMain: "#e2e8f0",
    textMuted: "#94a3b8",
    textDim: "#64748b",
    border: "#334155",
    success: "#22c55e",
    warning: "#f59e0b",
    error: "#ef4444",
  };

  let name = $state("My Theme");
  let colors = $state<ThemeColors>({ ...defaultColors });

  const colorFields: { key: keyof ThemeColors; label: string; group: string }[] = [
    { key: "bgDark", label: "背景", group: "背景" },
    { key: "bgCard", label: "カード", group: "背景" },
    { key: "bgCardHover", label: "ホバー", group: "背景" },
    { key: "primary", label: "プライマリ", group: "アクセント" },
    { key: "primaryLight", label: "プライマリ明", group: "アクセント" },
    { key: "primaryDark", label: "プライマリ暗", group: "アクセント" },
    { key: "accent", label: "アクセント", group: "アクセント" },
    { key: "accentLight", label: "アクセント明", group: "アクセント" },
    { key: "textMain", label: "テキスト", group: "テキスト" },
    { key: "textMuted", label: "テキスト薄", group: "テキスト" },
    { key: "textDim", label: "テキスト暗", group: "テキスト" },
    { key: "border", label: "ボーダー", group: "その他" },
    { key: "success", label: "成功", group: "その他" },
    { key: "warning", label: "警告", group: "その他" },
    { key: "error", label: "エラー", group: "その他" },
  ];

  let groups = $derived(() => {
    const g = new Map<string, typeof colorFields>();
    for (const f of colorFields) {
      if (!g.has(f.group)) g.set(f.group, []);
      g.get(f.group)!.push(f);
    }
    return Array.from(g.entries());
  });

  function preview() {
    applyTheme({
      slug: "custom-preview",
      name,
      colors,
    });
  }

  async function save() {
    const slug = name.toLowerCase().replace(/\s+/g, "-").replace(/[^a-z0-9-]/g, "");
    if (!slug) return;

    const theme: PrometheusTheme = {
      slug: `custom-${slug}`,
      name,
      colors: { ...colors },
      isCustom: true,
    };

    await addCustomTheme(theme);
    await setTheme(theme.slug);
    onclose();
  }

  function reset() {
    colors = { ...defaultColors };
    name = "My Theme";
  }

  $effect(() => {
    if (open) {
      colors = { ...defaultColors };
      name = "My Theme";
    }
  });
</script>

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <button class="absolute inset-0 bg-black/60 backdrop-blur-sm" onclick={onclose} aria-label="閉じる"></button>

    <div class="relative z-10 w-full max-w-lg overflow-hidden rounded-xl border border-border bg-bg-card shadow-2xl">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-border px-5 py-3">
        <div class="flex items-center gap-2">
          <Palette size={15} />
          <h2 class="text-sm font-semibold text-text-main">カスタムテーマ作成</h2>
        </div>
        <button onclick={onclose} class="text-text-dim hover:text-text-main">
          <X size={16} />
        </button>
      </div>

      <!-- Name -->
      <div class="border-b border-border px-5 py-3">
        <label class="mb-1 block text-[11px] text-text-dim">テーマ名</label>
        <input
          bind:value={name}
          class="w-full rounded border border-border bg-bg-dark px-3 py-1.5 text-sm text-text-main outline-none focus:border-primary"
          placeholder="My Theme"
        />
      </div>

      <!-- Colors -->
      <div class="max-h-72 overflow-y-auto px-5 py-3">
        {#each groups() as [groupName, fields]}
          <div class="mb-3">
            <p class="mb-1.5 text-[10px] font-medium uppercase tracking-wider text-text-dim">{groupName}</p>
            <div class="grid grid-cols-3 gap-2">
              {#each fields as field}
                <label class="flex items-center gap-2">
                  <input
                    type="color"
                    bind:value={colors[field.key]}
                    oninput={preview}
                    class="h-6 w-6 cursor-pointer rounded border border-border bg-transparent"
                  />
                  <span class="text-[11px] text-text-muted">{field.label}</span>
                </label>
              {/each}
            </div>
          </div>
        {/each}
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between border-t border-border px-5 py-3">
        <button
          onclick={reset}
          class="flex items-center gap-1.5 rounded-md border border-border px-3 py-1.5 text-xs text-text-muted hover:text-text-main"
        >
          <RotateCcw size={12} />
          リセット
        </button>
        <div class="flex gap-2">
          <button
            onclick={preview}
            class="rounded-md border border-border px-3 py-1.5 text-xs text-text-muted hover:text-text-main"
          >
            プレビュー
          </button>
          <button
            onclick={save}
            class="flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-white hover:bg-primary-dark"
          >
            <Save size={12} />
            保存
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
