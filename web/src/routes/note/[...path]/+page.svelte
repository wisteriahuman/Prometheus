<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import EditorLayout from "$lib/components/editor/EditorLayout.svelte";
  import BacklinkPanel from "$lib/components/editor/BacklinkPanel.svelte";
  import TagEditor from "$lib/components/editor/TagEditor.svelte";
  import { getAllThemes, currentTheme } from "$lib/stores/theme";
  import { isBareNoteContent, type ViewMode } from "$lib/utils/noteViewMode";
  import { get } from "svelte/store";
  import { FileText, Palette, Download } from "lucide-svelte";

  interface NoteData {
    path: string;
    title: string;
    content: string;
    frontmatter: {
      id: string;
      title: string;
      created: string;
      modified: string;
      tags: string[];
      theme?: string;
    };
    html: string;
  }

  let noteData = $state<NoteData | null>(null);
  let loading = $state(true);
  let previewHtml = $state("");
  let saving = $state(false);
  let lastSaved = $state<string | null>(null);
  let debounceTimer: ReturnType<typeof setTimeout> | undefined;
  let editingTitle = $state(false);
  let showExportMenu = $state(false);
  let autoCreatedNote = $state(false);

  function exportNote(format: string) {
    showExportMenu = false;
    if (!noteData) return;

    // Include theme: note theme > app theme
    const theme = noteThemeSlug || get(currentTheme).slug;
    const themeParam = `&theme=${theme}`;

    if (format === "pdf") {
      // Open themed HTML inline, user adjusts print settings manually
      const printWindow = window.open(
        `/api/export/${noteData.path}?format=html&inline=true${themeParam}`,
        "_blank",
      );
      if (printWindow) {
        printWindow.addEventListener("load", () => {
          setTimeout(() => printWindow.print(), 500);
        });
      }
      return;
    }

    window.open(`/api/export/${noteData.path}?format=${format}${themeParam}`, "_blank");
  }
  let titleInput = $state("");
  let showThemePicker = $state(false);
  let noteThemeSlug = $state<string | undefined>(undefined);

  // Scoped note theme CSS
  let noteThemeStyle = $derived(() => {
    if (!noteThemeSlug) return "";
    const theme = getAllThemes().find((t) => t.slug === noteThemeSlug);
    if (!theme) return "";
    const c = theme.colors;
    return [
      `--color-bg-dark: ${c.bgDark}`,
      `--color-bg-card: ${c.bgCard}`,
      `--color-bg-card-hover: ${c.bgCardHover}`,
      `--color-bg-input: ${c.bgDark}`,
      `--color-primary: ${c.primary}`,
      `--color-primary-light: ${c.primaryLight}`,
      `--color-primary-dark: ${c.primaryDark}`,
      `--color-accent: ${c.accent}`,
      `--color-accent-light: ${c.accentLight}`,
      `--color-text-main: ${c.textMain}`,
      `--color-text-muted: ${c.textMuted}`,
      `--color-text-dim: ${c.textDim}`,
      `--color-border: ${c.border}`,
      `--color-success: ${c.success}`,
      `--color-warning: ${c.warning}`,
      `--color-error: ${c.error}`,
    ].join("; ");
  });

  let notePath = $derived($page.params.path);

  async function loadNote() {
    const currentPath = notePath;
    if (!currentPath) return;

    loading = true;
    autoCreatedNote = false;
    try {
      const res = await fetch(`/api/notes/${currentPath}`);
      if (!res.ok) {
        // Auto-create if it doesn't exist and path ends with .md
        if (res.status === 404 && currentPath.endsWith(".md")) {
          await fetch("/api/notes", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ path: currentPath }),
          });
          const retry = await fetch(`/api/notes/${currentPath}`);
          if (retry.ok) {
            autoCreatedNote = true;
            noteData = await retry.json();
          }
        }
      } else {
        noteData = await res.json();
      }

      if (noteData) {
        // Get preview HTML
        const previewRes = await fetch("/api/preview", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ content: noteData.content }),
        });
        const previewResult = await previewRes.json();
        previewHtml = previewResult.html;
        noteThemeSlug = noteData.frontmatter.theme;
      }
    } catch (e) {
      console.error("Failed to load note:", e);
    } finally {
      loading = false;
    }
  }

  // Load on mount and when path changes
  onMount(() => loadNote());

  $effect(() => {
    // Re-load when navigating to a different note
    if (notePath) {
      loadNote();
    }
  });

  async function setNoteTheme(themeSlug: string | undefined) {
    showThemePicker = false;
    noteThemeSlug = themeSlug;
    if (!noteData) return;

    try {
      const fm = { ...noteData.frontmatter };
      fm.theme = themeSlug ?? (null as any);

      await fetch(`/api/notes/${noteData.path}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ content: noteData.content, frontmatter: fm }),
      });

      window.dispatchEvent(new CustomEvent("prometheus:refresh-tree"));
    } catch (e) {
      console.error("Failed to set theme:", e);
    }
  }

  async function saveTitle() {
    const newTitle = titleInput.trim();
    if (!newTitle || !noteData || newTitle === noteData.frontmatter.title) {
      editingTitle = false;
      return;
    }

    try {
      await fetch(`/api/notes/${noteData.path}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: noteData.content,
          frontmatter: { ...noteData.frontmatter, title: newTitle },
        }),
      });
      noteData.frontmatter.title = newTitle;
    } catch (e) {
      console.error("Failed to save title:", e);
    }
    editingTitle = false;
  }

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
    if (noteData) noteData.content = newContent;
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => updatePreview(newContent), 300);
  }

  async function handleSave(saveContent: string) {
    if (!noteData) return;
    saving = true;
    try {
      const res = await fetch(`/api/notes/${noteData.path}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: saveContent,
          frontmatter: noteData.frontmatter,
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

  function initialViewMode(): ViewMode {
    if (!noteData) return "preview";
    if ($page.url.searchParams.get("new") === "1" || autoCreatedNote) return "split";
    return isBareNoteContent(noteData.content, noteData.frontmatter.title) ? "split" : "preview";
  }
</script>

<svelte:head>
  <title>{noteData?.frontmatter.title ?? "読み込み中..."} - Prometheus</title>
</svelte:head>

{#if loading}
  <div class="flex h-full items-center justify-center">
    <span class="text-text-muted">読み込み中...</span>
  </div>
{:else if noteData}
  <div class="flex h-full flex-col -m-5">
    <!-- Note header -->
    <div class="flex items-center gap-3 border-b border-border px-5 py-2.5">
      <FileText size={14} class="shrink-0 text-text-dim" />
      {#if editingTitle}
        <form onsubmit={(e) => { e.preventDefault(); saveTitle(); }} class="flex items-center">
          <input
            bind:value={titleInput}
            onblur={saveTitle}
            onkeydown={(e) => { if (e.key === "Escape") { editingTitle = false; } }}
            class="rounded border border-primary bg-bg-input px-2 py-0.5 text-sm font-medium text-text-main outline-none"
          />
        </form>
      {:else}
        <button
          onclick={() => { editingTitle = true; titleInput = noteData?.frontmatter.title ?? ""; }}
          class="text-sm font-medium text-text-main hover:text-primary-light"
          title="クリックでタイトルを編集"
        >
          {noteData.frontmatter.title}
        </button>
      {/if}
      <span class="text-[11px] text-text-dim">{noteData.path}</span>
      <div class="flex-1"></div>
      {#if saving}
        <span class="text-[11px] text-text-muted">保存中...</span>
      {:else if lastSaved}
        <span class="text-[11px] text-text-dim">保存済み {lastSaved}</span>
      {/if}
      <TagEditor
        tags={noteData.frontmatter.tags}
        notePath={noteData.path}
        frontmatter={noteData.frontmatter}
      />

      <!-- Export -->
      <div class="relative">
        <button
          onclick={() => (showExportMenu = !showExportMenu)}
          class="flex items-center gap-1 rounded-md border border-border px-1.5 py-0.5 text-[10px] text-text-dim transition-colors hover:border-primary hover:text-text-muted"
          title="エクスポート"
        >
          <Download size={10} />
        </button>

        {#if showExportMenu}
          <div class="absolute right-0 top-full z-50 mt-1 min-w-[120px] rounded-lg border border-border bg-bg-card py-1 shadow-xl">
            <button
              onclick={() => exportNote("html")}
              class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-text-muted hover:bg-bg-card-hover"
            >
              HTML
            </button>
            <button
              onclick={() => exportNote("md")}
              class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-text-muted hover:bg-bg-card-hover"
            >
              Markdown
            </button>
            <button
              onclick={() => exportNote("pdf")}
              class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-text-muted hover:bg-bg-card-hover"
            >
              PDF（印刷）
            </button>
          </div>
        {/if}
      </div>

      <!-- Note theme picker -->
      <div class="relative">
        <button
          onclick={() => (showThemePicker = !showThemePicker)}
          class="flex items-center gap-1 rounded-md border border-border px-1.5 py-0.5 text-[10px] text-text-dim transition-colors hover:border-primary hover:text-text-muted"
          title="ノートのテーマ"
        >
          <Palette size={10} />
          {noteThemeSlug ?? "デフォルト"}
        </button>

        {#if showThemePicker}
          <div class="absolute right-0 top-full z-50 mt-1 min-w-[120px] rounded-lg border border-border bg-bg-card py-1 shadow-xl">
            <button
              onclick={() => setNoteTheme(undefined)}
              class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-text-muted hover:bg-bg-card-hover"
            >
              デフォルト {!noteThemeSlug ? "✓" : ""}
            </button>
            {#each getAllThemes() as theme}
              <button
                onclick={() => setNoteTheme(theme.slug)}
                class="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-text-muted hover:bg-bg-card-hover"
              >
                <span class="h-2.5 w-2.5 rounded-full" style="background-color: {theme.colors.primary}"></span>
                {theme.name}
                {#if noteThemeSlug === theme.slug}
                  <span class="ml-auto text-primary">✓</span>
                {/if}
              </button>
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <!-- Editor area: scoped to note theme -->
    <div class="flex-1 overflow-hidden" style={noteThemeStyle()}>
      {#key `${noteData.path}:${$page.url.searchParams.get("new") ?? ""}`}
        <EditorLayout
          content={noteData.content}
          {previewHtml}
          vimMode={true}
          initialViewMode={initialViewMode()}
          onchange={handleChange}
          onsave={handleSave}
        />
      {/key}
    </div>

    <BacklinkPanel notePath={noteData.path} />
  </div>
{:else}
  <div class="flex h-full items-center justify-center">
    <span class="text-text-muted">ノートが見つかりません</span>
  </div>
{/if}
