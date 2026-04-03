<script lang="ts">
  import { tick } from "svelte";

  interface Props {
    html: string;
  }

  let { html }: Props = $props();
  let previewEl: HTMLDivElement | undefined = $state(undefined);

  let mermaidModule: typeof import("mermaid").default | null = null;
  let renderCounter = 0;

  async function loadMermaid() {
    if (mermaidModule) return mermaidModule;
    const m = (await import("mermaid")).default;
    m.initialize({
      startOnLoad: false,
      theme: "base",
      themeVariables: {
        background: "#1e293b",
        primaryColor: "#334155",
        primaryTextColor: "#e2e8f0",
        primaryBorderColor: "#475569",
        lineColor: "#94a3b8",
        secondaryColor: "#1e293b",
        tertiaryColor: "#0f172a",
        fontFamily: "Inter Variable, Noto Sans JP Variable, sans-serif",
      },
    });
    mermaidModule = m;
    return m;
  }

  async function renderMermaid() {
    if (!previewEl) return;

    const codeBlocks = previewEl.querySelectorAll("pre > code.language-mermaid");
    if (codeBlocks.length === 0) return;

    const mermaid = await loadMermaid();

    for (const block of codeBlocks) {
      const pre = block.parentElement;
      if (!pre) continue;

      const mermaidCode = block.textContent?.trim() ?? "";
      if (!mermaidCode) continue;

      try {
        renderCounter++;
        const id = `mermaid-${renderCounter}`;
        const { svg } = await mermaid.render(id, mermaidCode);

        const container = document.createElement("div");
        container.className = "mermaid-container";
        container.innerHTML = svg;
        pre.replaceWith(container);
      } catch {
        // Leave code block as-is if mermaid fails
      }
    }
  }

  $effect(() => {
    if (html && previewEl) {
      // Wait for DOM update then render mermaid
      tick().then(() => renderMermaid());
    }
  });
</script>

<div bind:this={previewEl} class="prose h-full overflow-y-auto bg-bg-dark px-10 py-8">
  {@html html}
</div>

<style>
  .prose {
    color: var(--color-text-main);
    font-size: 15px;
    line-height: 1.8;
    max-width: 72ch;
  }

  .prose :global(h1) {
    color: var(--color-primary-light);
    font-size: 1.75em;
    font-weight: 700;
    margin: 0 0 0.6em;
    padding-bottom: 0.4em;
    border-bottom: 1px solid var(--color-border);
    line-height: 1.3;
  }

  .prose :global(h2) {
    color: var(--color-primary-light);
    font-size: 1.35em;
    font-weight: 600;
    margin: 1.8em 0 0.6em;
    line-height: 1.3;
  }

  .prose :global(h3) {
    color: var(--color-text-main);
    font-size: 1.15em;
    font-weight: 600;
    margin: 1.5em 0 0.4em;
    line-height: 1.4;
  }

  .prose :global(h4),
  .prose :global(h5),
  .prose :global(h6) {
    color: var(--color-text-main);
    font-weight: 600;
    margin: 1.2em 0 0.4em;
  }

  .prose :global(p) {
    margin: 0.8em 0;
  }

  .prose :global(a) {
    color: var(--color-accent);
    text-decoration: underline;
    text-underline-offset: 2px;
    text-decoration-thickness: 1px;
  }

  .prose :global(a:hover) {
    color: var(--color-accent-light);
  }

  .prose :global(ul) {
    list-style-type: disc;
    padding-left: 1.75em;
    margin: 0.8em 0;
  }

  .prose :global(ol) {
    list-style-type: decimal;
    padding-left: 1.75em;
    margin: 0.8em 0;
  }

  .prose :global(ul ul) { list-style-type: circle; margin: 0.2em 0; }
  .prose :global(ul ul ul) { list-style-type: square; }
  .prose :global(ol ol) { list-style-type: lower-alpha; margin: 0.2em 0; }

  .prose :global(li) {
    margin: 0.25em 0;
    padding-left: 0.25em;
  }

  .prose :global(li > p) { margin: 0.3em 0; }

  .prose :global(.contains-task-list) {
    list-style-type: none;
    padding-left: 0.5em;
  }

  .prose :global(.task-list-item) {
    display: flex;
    align-items: baseline;
    gap: 0.5em;
    padding-left: 0;
  }

  .prose :global(li input[type="checkbox"]) {
    margin: 0;
    accent-color: var(--color-primary);
    position: relative;
    top: 1px;
  }

  .prose :global(code) {
    color: var(--color-accent);
    background-color: var(--color-bg-card);
    padding: 0.15em 0.4em;
    border-radius: 4px;
    font-family: var(--font-mono);
    font-size: 0.875em;
  }

  .prose :global(pre) {
    background-color: var(--color-bg-card);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 1em 1.25em;
    overflow-x: auto;
    margin: 1.2em 0;
    line-height: 1.6;
  }

  .prose :global(pre code) {
    background: none;
    padding: 0;
    color: var(--color-text-main);
    font-size: 0.85em;
  }

  .prose :global(blockquote) {
    border-left: 3px solid var(--color-primary);
    padding: 0.5em 0 0.5em 1.25em;
    color: var(--color-text-muted);
    font-style: italic;
    margin: 1.2em 0;
  }

  .prose :global(blockquote p) { margin: 0.3em 0; }

  .prose :global(hr) {
    border: none;
    border-top: 1px solid var(--color-border);
    margin: 2em 0;
  }

  .prose :global(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 1.2em 0;
    font-size: 0.9em;
  }

  .prose :global(th) {
    background-color: var(--color-bg-card);
    border: 1px solid var(--color-border);
    padding: 0.5em 0.75em;
    text-align: left;
    font-weight: 600;
  }

  .prose :global(td) {
    border: 1px solid var(--color-border);
    padding: 0.5em 0.75em;
  }

  .prose :global(img) {
    max-width: 100%;
    border-radius: 8px;
    margin: 1em 0;
  }

  .prose :global(strong) {
    color: var(--color-text-main);
    font-weight: 600;
  }

  .prose :global(em) {
    color: var(--color-text-muted);
  }

  /* Mermaid diagrams */
  .prose :global(.mermaid-container) {
    display: flex;
    justify-content: center;
    background-color: var(--color-bg-card);
    padding: 1.5em;
    border-radius: 8px;
    margin: 1.2em 0;
    overflow-x: auto;
    border: 1px solid var(--color-border);
  }

  .prose :global(.mermaid-container svg) {
    max-width: 100%;
  }
</style>
