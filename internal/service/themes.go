package service

// ThemeColors defines CSS colors for a theme
type ThemeColors struct {
	BgDark      string
	BgCard      string
	BgCardHover string
	Primary     string
	PrimaryLight string
	Accent      string
	AccentLight string
	TextMain    string
	TextMuted   string
	TextDim     string
	Border      string
	Success     string
	Warning     string
	Error       string
}

// BuiltinThemes maps theme slugs to their color definitions.
// Must stay in sync with web/src/lib/stores/theme.ts.
var BuiltinThemes = map[string]ThemeColors{
	"prometheus": {
		BgDark: "#0f172a", BgCard: "#1e293b", BgCardHover: "#334155",
		Primary: "#6366f1", PrimaryLight: "#818cf8", Accent: "#38bdf8", AccentLight: "#7dd3fc",
		TextMain: "#e2e8f0", TextMuted: "#94a3b8", TextDim: "#64748b",
		Border: "#334155", Success: "#22c55e", Warning: "#f59e0b", Error: "#ef4444",
	},
	"ocean": {
		BgDark: "#0c1222", BgCard: "#162032", BgCardHover: "#1e2d45",
		Primary: "#0ea5e9", PrimaryLight: "#38bdf8", Accent: "#06b6d4", AccentLight: "#22d3ee",
		TextMain: "#e0f2fe", TextMuted: "#7dd3fc", TextDim: "#38bdf8",
		Border: "#1e3a5f", Success: "#34d399", Warning: "#fbbf24", Error: "#f87171",
	},
	"forest": {
		BgDark: "#0f1a0f", BgCard: "#1a2e1a", BgCardHover: "#243824",
		Primary: "#22c55e", PrimaryLight: "#4ade80", Accent: "#a3e635", AccentLight: "#bef264",
		TextMain: "#dcfce7", TextMuted: "#86efac", TextDim: "#4ade80",
		Border: "#1e3a1e", Success: "#34d399", Warning: "#fbbf24", Error: "#f87171",
	},
	"sakura": {
		BgDark: "#1a0f18", BgCard: "#2e1a2a", BgCardHover: "#3d2438",
		Primary: "#ec4899", PrimaryLight: "#f472b6", Accent: "#e879f9", AccentLight: "#f0abfc",
		TextMain: "#fce7f3", TextMuted: "#f9a8d4", TextDim: "#f472b6",
		Border: "#4a2040", Success: "#34d399", Warning: "#fbbf24", Error: "#f87171",
	},
	"nord": {
		BgDark: "#2e3440", BgCard: "#3b4252", BgCardHover: "#434c5e",
		Primary: "#88c0d0", PrimaryLight: "#8fbcbb", Accent: "#81a1c1", AccentLight: "#88c0d0",
		TextMain: "#eceff4", TextMuted: "#d8dee9", TextDim: "#a0aec0",
		Border: "#4c566a", Success: "#a3be8c", Warning: "#ebcb8b", Error: "#bf616a",
	},
	"solarized-dark": {
		BgDark: "#002b36", BgCard: "#073642", BgCardHover: "#0a4050",
		Primary: "#268bd2", PrimaryLight: "#2aa198", Accent: "#b58900", AccentLight: "#cb4b16",
		TextMain: "#fdf6e3", TextMuted: "#93a1a1", TextDim: "#657b83",
		Border: "#586e75", Success: "#859900", Warning: "#b58900", Error: "#dc322f",
	},
	"midnight": {
		BgDark: "#0a0a0f", BgCard: "#12121a", BgCardHover: "#1a1a25",
		Primary: "#a78bfa", PrimaryLight: "#c4b5fd", Accent: "#f472b6", AccentLight: "#f9a8d4",
		TextMain: "#e8e8f0", TextMuted: "#9898b0", TextDim: "#585870",
		Border: "#22222e", Success: "#34d399", Warning: "#fbbf24", Error: "#f87171",
	},
	"light": {
		BgDark: "#ffffff", BgCard: "#f8fafc", BgCardHover: "#f1f5f9",
		Primary: "#4f46e5", PrimaryLight: "#6366f1", Accent: "#0284c7", AccentLight: "#0ea5e9",
		TextMain: "#1e293b", TextMuted: "#64748b", TextDim: "#94a3b8",
		Border: "#e2e8f0", Success: "#16a34a", Warning: "#d97706", Error: "#dc2626",
	},
	"solarized-light": {
		BgDark: "#fdf6e3", BgCard: "#eee8d5", BgCardHover: "#e8e1cc",
		Primary: "#268bd2", PrimaryLight: "#2aa198", Accent: "#859900", AccentLight: "#b58900",
		TextMain: "#073642", TextMuted: "#586e75", TextDim: "#93a1a1",
		Border: "#d6cdb5", Success: "#859900", Warning: "#b58900", Error: "#dc322f",
	},
	"corporate": {
		BgDark: "#f9fafb", BgCard: "#ffffff", BgCardHover: "#f3f4f6",
		Primary: "#1f2937", PrimaryLight: "#374151", Accent: "#2563eb", AccentLight: "#3b82f6",
		TextMain: "#111827", TextMuted: "#6b7280", TextDim: "#9ca3af",
		Border: "#e5e7eb", Success: "#059669", Warning: "#d97706", Error: "#dc2626",
	},
	"slate": {
		BgDark: "#1e293b", BgCard: "#273548", BgCardHover: "#334155",
		Primary: "#3b82f6", PrimaryLight: "#60a5fa", Accent: "#f59e0b", AccentLight: "#fbbf24",
		TextMain: "#f1f5f9", TextMuted: "#cbd5e1", TextDim: "#94a3b8",
		Border: "#3d4f65", Success: "#10b981", Warning: "#f59e0b", Error: "#ef4444",
	},
	"executive": {
		BgDark: "#fafaf9", BgCard: "#ffffff", BgCardHover: "#f5f5f4",
		Primary: "#78716c", PrimaryLight: "#a8a29e", Accent: "#b45309", AccentLight: "#d97706",
		TextMain: "#1c1917", TextMuted: "#78716c", TextDim: "#a8a29e",
		Border: "#e7e5e4", Success: "#15803d", Warning: "#a16207", Error: "#b91c1c",
	},
}

// GetThemeColors returns colors for a theme slug, falling back to "light"
func GetThemeColors(slug string) ThemeColors {
	if t, ok := BuiltinThemes[slug]; ok {
		return t
	}
	return BuiltinThemes["light"]
}
