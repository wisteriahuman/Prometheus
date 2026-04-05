<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { EditorState, type Extension } from "@codemirror/state";
  import { EditorView, keymap, lineNumbers, highlightActiveLine, highlightActiveLineGutter, drawSelection } from "@codemirror/view";
  import { defaultKeymap, history, historyKeymap, indentWithTab } from "@codemirror/commands";
  import { markdown, markdownLanguage } from "@codemirror/lang-markdown";
  import { languages } from "@codemirror/language-data";
  import { syntaxHighlighting, defaultHighlightStyle, bracketMatching, indentOnInput } from "@codemirror/language";
  import { searchKeymap, highlightSelectionMatches } from "@codemirror/search";
  import { vim, Vim } from "@replit/codemirror-vim";
  import { prometheusTheme } from "./theme.js";

  interface Props {
    content: string;
    vimMode?: boolean;
    onchange?: (content: string) => void;
    onsave?: (content: string) => void;
  }

  let { content, vimMode = true, onchange, onsave }: Props = $props();

  let editorContainer: HTMLDivElement;
  let view: EditorView | undefined;
  let dragging = $state(false);
  let dragCounter = 0;

  // Mutable ref so closures always access the latest callback
  const saveRef: { fn: ((content: string) => void) | undefined } = { fn: undefined };
  $effect(() => {
    saveRef.fn = onsave;
  });

  async function uploadAndInsert(file: File) {
    if (!view) return;

    const formData = new FormData();
    formData.append("file", file);

    try {
      const res = await fetch("/api/assets", {
        method: "POST",
        body: formData,
      });
      if (!res.ok) return;

      const result = await res.json();
      const markdown = result.markdown as string;

      // Insert markdown at cursor position
      const pos = view.state.selection.main.head;
      view.dispatch({
        changes: { from: pos, insert: markdown + "\n" },
        selection: { anchor: pos + markdown.length + 1 },
      });
    } catch (e) {
      console.error("Upload failed:", e);
    }
  }

  function doSave() {
    if (saveRef.fn && view) {
      saveRef.fn(view.state.doc.toString());
    }
  }

  function registerVimCommands() {
    Vim.defineEx("write", "w", () => doSave());
    Vim.defineEx("wq", "wq", () => doSave());
  }

  function createExtensions(): Extension[] {
    const extensions: Extension[] = [];

    if (vimMode) {
      extensions.push(vim());
    }

    extensions.push(
      keymap.of([
        {
          key: "Mod-s",
          run: () => {
            doSave();
            return true;
          },
        },
      ]),
      lineNumbers(),
      highlightActiveLine(),
      highlightActiveLineGutter(),
      drawSelection(),
      indentOnInput(),
      bracketMatching(),
      highlightSelectionMatches(),
      history(),
      syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
      markdown({ base: markdownLanguage, codeLanguages: languages }),
      prometheusTheme,
      keymap.of([
        ...defaultKeymap,
        ...searchKeymap,
        ...historyKeymap,
        indentWithTab,
      ]),
      EditorView.updateListener.of((update) => {
        if (update.docChanged && onchange) {
          onchange(update.state.doc.toString());
        }
      }),
      EditorView.lineWrapping,
    );

    return extensions;
  }

  onMount(() => {
    if (vimMode) {
      registerVimCommands();
    }

    const state = EditorState.create({
      doc: content,
      extensions: createExtensions(),
    });

    view = new EditorView({
      state,
      parent: editorContainer,
    });

    // Listen for save events from UI buttons
    const handleSaveEvent = () => doSave();
    window.addEventListener("prometheus:save", handleSaveEvent);

    // File drop handler — use counter to handle child element enter/leave
    editorContainer.addEventListener("dragenter", (e) => {
      e.preventDefault();
      dragCounter++;
      dragging = true;
    });
    editorContainer.addEventListener("dragover", (e) => {
      e.preventDefault();
    });
    editorContainer.addEventListener("dragleave", () => {
      dragCounter--;
      if (dragCounter <= 0) {
        dragCounter = 0;
        dragging = false;
      }
    });
    editorContainer.addEventListener("drop", async (e) => {
      e.preventDefault();
      dragCounter = 0;
      dragging = false;
      const files = e.dataTransfer?.files;
      if (files) {
        for (const file of files) {
          await uploadAndInsert(file);
        }
      }
    });

    // Clipboard paste handler (images)
    editorContainer.addEventListener("paste", async (e) => {
      const items = e.clipboardData?.items;
      if (!items) return;
      for (const item of items) {
        if (item.type.startsWith("image/")) {
          e.preventDefault();
          const file = item.getAsFile();
          if (file) await uploadAndInsert(file);
        }
      }
    });

    return () => window.removeEventListener("prometheus:save", handleSaveEvent);
  });

  onDestroy(() => {
    view?.destroy();
  });

  $effect(() => {
    if (view && content !== view.state.doc.toString()) {
      view.dispatch({
        changes: {
          from: 0,
          to: view.state.doc.length,
          insert: content,
        },
      });
    }
  });

  export function getContent(): string {
    return view?.state.doc.toString() ?? content;
  }

  export function focus() {
    view?.focus();
  }
</script>

<div class="relative h-full w-full">
  <div
    bind:this={editorContainer}
    class="h-full w-full overflow-auto"
  ></div>

  {#if dragging}
    <div class="drop-overlay">
      <div class="drop-overlay-content">
        <span class="text-sm font-medium">ファイルをドロップしてアップロード</span>
      </div>
    </div>
  {/if}
</div>

<style>
  :global(.cm-editor) {
    height: 100%;
    font-family: var(--font-mono);
    font-size: 14px;
  }

  :global(.cm-editor .cm-scroller) {
    overflow: auto;
    padding: 8px 0;
  }

  :global(.cm-editor.cm-focused) {
    outline: none;
  }

  :global(.cm-editor .cm-content) {
    padding: 0 4px;
  }

  .drop-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: color-mix(in srgb, var(--color-primary) 8%, transparent);
    border: 2px dashed var(--color-primary);
    border-radius: 8px;
    pointer-events: none;
    z-index: 10;
  }

  .drop-overlay-content {
    padding: 12px 24px;
    border-radius: 8px;
    background: var(--color-bg-card);
    border: 1px solid var(--color-primary);
    color: var(--color-primary-light);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }
</style>
