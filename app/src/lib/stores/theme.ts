import { writable } from "svelte/store";
import type { PrometheusTheme } from "@prometheus/core";

export const BUILTIN_THEMES: PrometheusTheme[] = [
  {
    slug: "prometheus",
    name: "Prometheus",
    colors: {
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
    },
  },
  {
    slug: "ocean",
    name: "Ocean",
    colors: {
      bgDark: "#0c1222",
      bgCard: "#162032",
      bgCardHover: "#1e2d45",
      primary: "#0ea5e9",
      primaryLight: "#38bdf8",
      primaryDark: "#0284c7",
      accent: "#06b6d4",
      accentLight: "#22d3ee",
      textMain: "#e0f2fe",
      textMuted: "#7dd3fc",
      textDim: "#38bdf8",
      border: "#1e3a5f",
      success: "#34d399",
      warning: "#fbbf24",
      error: "#f87171",
    },
  },
  {
    slug: "forest",
    name: "Forest",
    colors: {
      bgDark: "#0f1a0f",
      bgCard: "#1a2e1a",
      bgCardHover: "#243824",
      primary: "#22c55e",
      primaryLight: "#4ade80",
      primaryDark: "#16a34a",
      accent: "#a3e635",
      accentLight: "#bef264",
      textMain: "#dcfce7",
      textMuted: "#86efac",
      textDim: "#4ade80",
      border: "#1e3a1e",
      success: "#34d399",
      warning: "#fbbf24",
      error: "#f87171",
    },
  },
  {
    slug: "sakura",
    name: "Sakura",
    colors: {
      bgDark: "#1a0f18",
      bgCard: "#2e1a2a",
      bgCardHover: "#3d2438",
      primary: "#ec4899",
      primaryLight: "#f472b6",
      primaryDark: "#db2777",
      accent: "#e879f9",
      accentLight: "#f0abfc",
      textMain: "#fce7f3",
      textMuted: "#f9a8d4",
      textDim: "#f472b6",
      border: "#4a2040",
      success: "#34d399",
      warning: "#fbbf24",
      error: "#f87171",
    },
  },
  {
    slug: "nord",
    name: "Nord",
    colors: {
      bgDark: "#2e3440",
      bgCard: "#3b4252",
      bgCardHover: "#434c5e",
      primary: "#88c0d0",
      primaryLight: "#8fbcbb",
      primaryDark: "#5e81ac",
      accent: "#81a1c1",
      accentLight: "#88c0d0",
      textMain: "#eceff4",
      textMuted: "#d8dee9",
      textDim: "#a0aec0",
      border: "#4c566a",
      success: "#a3be8c",
      warning: "#ebcb8b",
      error: "#bf616a",
    },
  },
  {
    slug: "solarized",
    name: "Solarized Dark",
    colors: {
      bgDark: "#002b36",
      bgCard: "#073642",
      bgCardHover: "#0a4050",
      primary: "#268bd2",
      primaryLight: "#2aa198",
      primaryDark: "#1a6fb5",
      accent: "#b58900",
      accentLight: "#cb4b16",
      textMain: "#fdf6e3",
      textMuted: "#93a1a1",
      textDim: "#657b83",
      border: "#586e75",
      success: "#859900",
      warning: "#b58900",
      error: "#dc322f",
    },
  },
];

const STORAGE_KEY = "prometheus-theme";

function getInitialTheme(): PrometheusTheme {
  if (typeof window !== "undefined") {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      const found = BUILTIN_THEMES.find((t) => t.slug === saved);
      if (found) return found;
    }
  }
  return BUILTIN_THEMES[0];
}

export const currentTheme = writable<PrometheusTheme>(BUILTIN_THEMES[0]);

export function initTheme() {
  const theme = getInitialTheme();
  currentTheme.set(theme);
  applyTheme(theme);
}

export function setTheme(slug: string) {
  const theme = BUILTIN_THEMES.find((t) => t.slug === slug);
  if (!theme) return;

  currentTheme.set(theme);
  applyTheme(theme);

  if (typeof window !== "undefined") {
    localStorage.setItem(STORAGE_KEY, slug);
  }
}

export function applyTheme(theme: PrometheusTheme) {
  if (typeof document === "undefined") return;

  const root = document.documentElement;
  const c = theme.colors;

  // Set ALL CSS variables used throughout the app
  root.style.setProperty("--color-bg-dark", c.bgDark);
  root.style.setProperty("--color-bg-card", c.bgCard);
  root.style.setProperty("--color-bg-card-hover", c.bgCardHover);
  root.style.setProperty("--color-bg-input", c.bgDark);
  root.style.setProperty("--color-primary", c.primary);
  root.style.setProperty("--color-primary-light", c.primaryLight);
  root.style.setProperty("--color-primary-dark", c.primaryDark);
  root.style.setProperty("--color-accent", c.accent);
  root.style.setProperty("--color-accent-light", c.accentLight);
  root.style.setProperty("--color-text-main", c.textMain);
  root.style.setProperty("--color-text-muted", c.textMuted);
  root.style.setProperty("--color-text-dim", c.textDim);
  root.style.setProperty("--color-border", c.border);
  root.style.setProperty("--color-border-focus", c.primary);
  root.style.setProperty("--color-success", c.success);
  root.style.setProperty("--color-warning", c.warning);
  root.style.setProperty("--color-error", c.error);
}

// Note-level theme is now applied via scoped CSS variables on the editor container,
// not by overriding :root. See note/[...path]/+page.svelte.
