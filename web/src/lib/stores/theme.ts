import { writable } from "svelte/store";

export interface ThemeColors {
  bgDark: string;
  bgCard: string;
  bgCardHover: string;
  primary: string;
  primaryLight: string;
  primaryDark: string;
  accent: string;
  accentLight: string;
  textMain: string;
  textMuted: string;
  textDim: string;
  border: string;
  success: string;
  warning: string;
  error: string;
}

export interface PrometheusTheme {
  slug: string;
  name: string;
  colors: ThemeColors;
  isCustom?: boolean;
}

export const BUILTIN_THEMES: PrometheusTheme[] = [
  // --- Dark themes ---
  {
    slug: "prometheus",
    name: "Prometheus",
    colors: {
      bgDark: "#0f172a", bgCard: "#1e293b", bgCardHover: "#334155",
      primary: "#6366f1", primaryLight: "#818cf8", primaryDark: "#4f46e5",
      accent: "#38bdf8", accentLight: "#7dd3fc",
      textMain: "#e2e8f0", textMuted: "#94a3b8", textDim: "#64748b",
      border: "#334155", success: "#22c55e", warning: "#f59e0b", error: "#ef4444",
    },
  },
  {
    slug: "ocean",
    name: "Ocean",
    colors: {
      bgDark: "#0c1222", bgCard: "#162032", bgCardHover: "#1e2d45",
      primary: "#0ea5e9", primaryLight: "#38bdf8", primaryDark: "#0284c7",
      accent: "#06b6d4", accentLight: "#22d3ee",
      textMain: "#e0f2fe", textMuted: "#7dd3fc", textDim: "#38bdf8",
      border: "#1e3a5f", success: "#34d399", warning: "#fbbf24", error: "#f87171",
    },
  },
  {
    slug: "forest",
    name: "Forest",
    colors: {
      bgDark: "#0f1a0f", bgCard: "#1a2e1a", bgCardHover: "#243824",
      primary: "#22c55e", primaryLight: "#4ade80", primaryDark: "#16a34a",
      accent: "#a3e635", accentLight: "#bef264",
      textMain: "#dcfce7", textMuted: "#86efac", textDim: "#4ade80",
      border: "#1e3a1e", success: "#34d399", warning: "#fbbf24", error: "#f87171",
    },
  },
  {
    slug: "sakura",
    name: "Sakura",
    colors: {
      bgDark: "#1a0f18", bgCard: "#2e1a2a", bgCardHover: "#3d2438",
      primary: "#ec4899", primaryLight: "#f472b6", primaryDark: "#db2777",
      accent: "#e879f9", accentLight: "#f0abfc",
      textMain: "#fce7f3", textMuted: "#f9a8d4", textDim: "#f472b6",
      border: "#4a2040", success: "#34d399", warning: "#fbbf24", error: "#f87171",
    },
  },
  {
    slug: "nord",
    name: "Nord",
    colors: {
      bgDark: "#2e3440", bgCard: "#3b4252", bgCardHover: "#434c5e",
      primary: "#88c0d0", primaryLight: "#8fbcbb", primaryDark: "#5e81ac",
      accent: "#81a1c1", accentLight: "#88c0d0",
      textMain: "#eceff4", textMuted: "#d8dee9", textDim: "#a0aec0",
      border: "#4c566a", success: "#a3be8c", warning: "#ebcb8b", error: "#bf616a",
    },
  },
  {
    slug: "solarized-dark",
    name: "Solarized Dark",
    colors: {
      bgDark: "#002b36", bgCard: "#073642", bgCardHover: "#0a4050",
      primary: "#268bd2", primaryLight: "#2aa198", primaryDark: "#1a6fb5",
      accent: "#b58900", accentLight: "#cb4b16",
      textMain: "#fdf6e3", textMuted: "#93a1a1", textDim: "#657b83",
      border: "#586e75", success: "#859900", warning: "#b58900", error: "#dc322f",
    },
  },
  {
    slug: "midnight",
    name: "Midnight",
    colors: {
      bgDark: "#0a0a0f", bgCard: "#12121a", bgCardHover: "#1a1a25",
      primary: "#a78bfa", primaryLight: "#c4b5fd", primaryDark: "#7c3aed",
      accent: "#f472b6", accentLight: "#f9a8d4",
      textMain: "#e8e8f0", textMuted: "#9898b0", textDim: "#5858708",
      border: "#22222e", success: "#34d399", warning: "#fbbf24", error: "#f87171",
    },
  },

  // --- Light themes ---
  {
    slug: "light",
    name: "Light",
    colors: {
      bgDark: "#ffffff", bgCard: "#f8fafc", bgCardHover: "#f1f5f9",
      primary: "#4f46e5", primaryLight: "#6366f1", primaryDark: "#4338ca",
      accent: "#0284c7", accentLight: "#0ea5e9",
      textMain: "#1e293b", textMuted: "#64748b", textDim: "#94a3b8",
      border: "#e2e8f0", success: "#16a34a", warning: "#d97706", error: "#dc2626",
    },
  },
  {
    slug: "solarized-light",
    name: "Solarized Light",
    colors: {
      bgDark: "#fdf6e3", bgCard: "#eee8d5", bgCardHover: "#e8e1cc",
      primary: "#268bd2", primaryLight: "#2aa198", primaryDark: "#1a6fb5",
      accent: "#859900", accentLight: "#b58900",
      textMain: "#073642", textMuted: "#586e75", textDim: "#93a1a1",
      border: "#d6cdb5", success: "#859900", warning: "#b58900", error: "#dc322f",
    },
  },

  // --- Business themes ---
  {
    slug: "corporate",
    name: "Corporate",
    colors: {
      bgDark: "#f9fafb", bgCard: "#ffffff", bgCardHover: "#f3f4f6",
      primary: "#1f2937", primaryLight: "#374151", primaryDark: "#111827",
      accent: "#2563eb", accentLight: "#3b82f6",
      textMain: "#111827", textMuted: "#6b7280", textDim: "#9ca3af",
      border: "#e5e7eb", success: "#059669", warning: "#d97706", error: "#dc2626",
    },
  },
  {
    slug: "slate",
    name: "Slate Business",
    colors: {
      bgDark: "#1e293b", bgCard: "#273548", bgCardHover: "#334155",
      primary: "#3b82f6", primaryLight: "#60a5fa", primaryDark: "#2563eb",
      accent: "#f59e0b", accentLight: "#fbbf24",
      textMain: "#f1f5f9", textMuted: "#cbd5e1", textDim: "#94a3b8",
      border: "#3d4f65", success: "#10b981", warning: "#f59e0b", error: "#ef4444",
    },
  },
  {
    slug: "executive",
    name: "Executive",
    colors: {
      bgDark: "#fafaf9", bgCard: "#ffffff", bgCardHover: "#f5f5f4",
      primary: "#78716c", primaryLight: "#a8a29e", primaryDark: "#57534e",
      accent: "#b45309", accentLight: "#d97706",
      textMain: "#1c1917", textMuted: "#78716c", textDim: "#a8a29e",
      border: "#e7e5e4", success: "#15803d", warning: "#a16207", error: "#b91c1c",
    },
  },
];

const STORAGE_KEY = "prometheus-theme";
const CUSTOM_THEMES_KEY = "prometheus-custom-themes";

function loadCustomThemes(): PrometheusTheme[] {
  if (typeof window === "undefined") return [];
  try {
    const data = localStorage.getItem(CUSTOM_THEMES_KEY);
    if (!data) return [];
    const themes = JSON.parse(data) as PrometheusTheme[];
    return themes.map((t) => ({ ...t, isCustom: true }));
  } catch {
    return [];
  }
}

function saveCustomThemes(themes: PrometheusTheme[]) {
  if (typeof window === "undefined") return;
  localStorage.setItem(CUSTOM_THEMES_KEY, JSON.stringify(themes));
}

export function getAllThemes(): PrometheusTheme[] {
  return [...BUILTIN_THEMES, ...loadCustomThemes()];
}

export function addCustomTheme(theme: PrometheusTheme): void {
  const customs = loadCustomThemes();
  const existing = customs.findIndex((t) => t.slug === theme.slug);
  if (existing >= 0) {
    customs[existing] = { ...theme, isCustom: true };
  } else {
    customs.push({ ...theme, isCustom: true });
  }
  saveCustomThemes(customs);
}

export function removeCustomTheme(slug: string): void {
  const customs = loadCustomThemes().filter((t) => t.slug !== slug);
  saveCustomThemes(customs);
}

export const currentTheme = writable<PrometheusTheme>(BUILTIN_THEMES[0]);

function getInitialTheme(): PrometheusTheme {
  if (typeof window !== "undefined") {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      const all = getAllThemes();
      const found = all.find((t) => t.slug === saved);
      if (found) return found;
    }
  }
  return BUILTIN_THEMES[0];
}

export function initTheme() {
  const theme = getInitialTheme();
  currentTheme.set(theme);
  applyTheme(theme);
}

export function setTheme(slug: string) {
  const all = getAllThemes();
  const theme = all.find((t) => t.slug === slug);
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
