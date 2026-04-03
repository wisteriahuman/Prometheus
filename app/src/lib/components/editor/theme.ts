import { EditorView } from "@codemirror/view";
import { HighlightStyle, syntaxHighlighting } from "@codemirror/language";
import { tags } from "@lezer/highlight";

const editorTheme = EditorView.theme(
  {
    "&": {
      color: "var(--color-text-main)",
      backgroundColor: "var(--color-bg-dark)",
    },
    ".cm-content": {
      caretColor: "var(--color-accent)",
      padding: "16px 0",
    },
    ".cm-cursor, .cm-dropCursor": {
      borderLeftColor: "var(--color-accent)",
      borderLeftWidth: "2px",
    },
    "&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection":
      {
        backgroundColor: "var(--color-bg-card-hover)",
      },
    ".cm-panels": {
      backgroundColor: "var(--color-bg-card)",
      color: "var(--color-text-main)",
    },
    ".cm-panels.cm-panels-top": {
      borderBottom: "1px solid var(--color-border)",
    },
    ".cm-panels.cm-panels-bottom": {
      borderTop: "1px solid var(--color-border)",
    },
    ".cm-activeLine": {
      backgroundColor: "color-mix(in srgb, var(--color-bg-card) 50%, transparent)",
    },
    ".cm-selectionMatch": {
      backgroundColor: "color-mix(in srgb, var(--color-accent) 15%, transparent)",
    },
    "&.cm-focused .cm-matchingBracket, &.cm-focused .cm-nonmatchingBracket": {
      backgroundColor: "color-mix(in srgb, var(--color-primary) 30%, transparent)",
    },
    ".cm-gutters": {
      backgroundColor: "var(--color-bg-dark)",
      color: "var(--color-text-dim)",
      border: "none",
      paddingRight: "8px",
    },
    ".cm-activeLineGutter": {
      backgroundColor: "color-mix(in srgb, var(--color-bg-card) 50%, transparent)",
      color: "var(--color-text-muted)",
    },
    ".cm-foldPlaceholder": {
      backgroundColor: "var(--color-bg-card)",
      border: "1px solid var(--color-border)",
      color: "var(--color-text-muted)",
    },
    ".cm-tooltip": {
      backgroundColor: "var(--color-bg-card)",
      border: "1px solid var(--color-border)",
      color: "var(--color-text-main)",
    },
    ".cm-vim-panel": {
      backgroundColor: "var(--color-bg-card)",
      color: "var(--color-text-muted)",
      padding: "2px 8px",
      fontFamily: "var(--font-mono)",
      fontSize: "12px",
    },
    // Syntax highlighting via CSS classes — these use CSS variables
    // so they change when the theme changes
    ".cm-s-keyword": { color: "var(--color-primary-light)" },
    ".cm-s-heading": { color: "var(--color-primary-light)", fontWeight: "bold" },
    ".cm-s-link": { color: "var(--color-accent)", textDecoration: "underline" },
    ".cm-s-string": { color: "var(--color-success)" },
    ".cm-s-comment": { color: "var(--color-text-dim)", fontStyle: "italic" },
    ".cm-s-quote": { color: "var(--color-text-muted)", fontStyle: "italic" },
    ".cm-s-emphasis": { fontStyle: "italic", color: "var(--color-text-muted)" },
    ".cm-s-strong": { fontWeight: "bold", color: "var(--color-text-main)" },
    ".cm-s-monospace": { color: "var(--color-accent)", fontFamily: "var(--font-mono)" },
    ".cm-s-meta": { color: "var(--color-text-dim)" },
  },
  { dark: true },
);

// Use className-based highlighting so CSS variables work for theme switching
const highlightStyles = HighlightStyle.define([
  { tag: tags.keyword, class: "cm-s-keyword" },
  { tag: tags.operator, color: "var(--color-text-muted)" },
  { tag: tags.special(tags.variableName), color: "var(--color-warning)" },
  { tag: tags.typeName, color: "var(--color-warning)" },
  { tag: tags.atom, color: "var(--color-accent)" },
  { tag: tags.number, color: "var(--color-warning)" },
  { tag: tags.bool, color: "var(--color-warning)" },
  { tag: tags.string, class: "cm-s-string" },
  { tag: tags.regexp, color: "var(--color-warning)" },
  { tag: tags.escape, color: "var(--color-warning)" },
  { tag: tags.definition(tags.variableName), color: "var(--color-accent)" },
  { tag: tags.function(tags.variableName), color: "var(--color-primary-light)" },
  { tag: tags.labelName, color: "var(--color-text-muted)" },
  { tag: tags.comment, class: "cm-s-comment" },
  { tag: tags.blockComment, class: "cm-s-comment" },
  { tag: tags.meta, class: "cm-s-meta" },
  { tag: tags.link, class: "cm-s-link" },
  { tag: tags.url, class: "cm-s-link" },
  { tag: tags.heading, class: "cm-s-heading" },
  { tag: tags.heading1, class: "cm-s-heading", fontSize: "1.4em" },
  { tag: tags.heading2, class: "cm-s-heading", fontSize: "1.2em" },
  { tag: tags.heading3, class: "cm-s-heading", fontSize: "1.1em" },
  { tag: tags.emphasis, class: "cm-s-emphasis" },
  { tag: tags.strong, class: "cm-s-strong" },
  { tag: tags.strikethrough, textDecoration: "line-through" },
  { tag: tags.processingInstruction, class: "cm-s-meta" },
  { tag: tags.monospace, class: "cm-s-monospace" },
  { tag: tags.quote, class: "cm-s-quote" },
  { tag: tags.invalid, color: "var(--color-error)" },
]);

export const prometheusTheme = [
  editorTheme,
  syntaxHighlighting(highlightStyles),
];
