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
}

export const DEFAULT_THEME: PrometheusTheme = {
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
};
