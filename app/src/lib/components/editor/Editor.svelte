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

  // Mutable ref so closures always access the latest callback
  const saveRef: { fn: ((content: string) => void) | undefined } = { fn: undefined };
  $effect(() => {
    saveRef.fn = onsave;
  });

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

<div
  bind:this={editorContainer}
  class="h-full w-full overflow-auto"
></div>

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
</style>
