package colors

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/chriserin/seq/internal/grid"
)

type Theme struct {
	colors           map[string]string
	accentColors     []string
	accentIcons      []rune
	lineActionColors map[grid.Action]string
}

var defaultColors = map[string]string{
	"AltSeqColor":                  "#222222",
	"SeqColor":                     "#000000",
	"SeqCursorColor":               "#444444",
	"SeqVisualColor":               "#aaaaaa",
	"SeqOverlayColor":              "#333388",
	"SeqMiddleOverlayColor":        "#405810",
	"SelectedAttributeColor":       "#5cdffb",
	"NumberColor":                  "#fcbd15",
	"Black":                        "#000000",
	"White":                        "#ffffff",
	"Heart":                        "#ed3902",
	"ActiveRatchetColor":           "#abfaa9",
	"MutedRatchetColor":            "#f34213",
	"CurrentPlayingColor":          "#abfaa9",
	"ActivePlayingColor":           "#f34213",
	"ArrangementHeaderColor":       "FAFAFA",
	"ArrangementTitleColor":        "FAFAFA",
	"ArrangementGroupColor":        "#F25D94",
	"ArrangementIndentColor":       "#4b4261",
	"ArrangementSelectedLineColor": "#3b4261",
	"SeqBorderLineColor":           "fafafa",
	"PatternModeColor":             "#ed3902",
}

var defaultTheme = Theme{
	colors: defaultColors,
	accentColors: []string{
		"#000000",
		"#ed3902",
		"#f564a9",
		"#f8730e",
		"#fcc05c",
		"#5cdffb",
		"#1e89ef",
		"#164de5",
		"#0246a7",
	},
	accentIcons: []rune{
		' ',
		'✤',
		'⎈',
		'⚙',
		'⊚',
		'✦',
		'❖',
		'✥',
		'❄',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#000000",
		grid.ACTION_LINE_RESET:     "#cf142b",
		grid.ACTION_LINE_REVERSE:   "#f8730e",
		grid.ACTION_LINE_SKIP_BEAT: "#a9e5bb",
		grid.ACTION_RESET:          "#fcf6b1",
		grid.ACTION_LINE_BOUNCE:    "#fcf6b1",
		grid.ACTION_LINE_DELAY:     "#cc4bc2",
	},
}

var seafoamColors = map[string]string{
	"AltSeqColor":                  "#1a3632",
	"SeqColor":                     "#0a2622",
	"SeqCursorColor":               "#2c5d54",
	"SeqVisualColor":               "#88c3b8",
	"SeqOverlayColor":              "#336699",
	"SeqMiddleOverlayColor":        "#4d7a68",
	"SelectedAttributeColor":       "#4fd1bf",
	"NumberColor":                  "#c5d86d",
	"Black":                        "#0a2622",
	"White":                        "#e6f2ef",
	"Heart":                        "#f76c5e",
	"ActiveRatchetColor":           "#90e0c9",
	"MutedRatchetColor":            "#f76c5e",
	"CurrentPlayingColor":          "#90e0c9",
	"ActivePlayingColor":           "#f76c5e",
	"ArrangementHeaderColor":       "#e6f2ef",
	"ArrangementTitleColor":        "#e6f2ef",
	"ArrangementGroupColor":        "#f76c5e",
	"ArrangementIndentColor":       "#4d7a68",
	"ArrangementSelectedLineColor": "#2c5d54",
}

var seafoamTheme = Theme{
	colors: seafoamColors,
	accentColors: []string{
		"#0a2622",
		"#f76c5e",
		"#f39c6b",
		"#ffd166",
		"#c5d86d",
		"#4fd1bf",
		"#2ec4b6",
		"#1a8fe3",
		"#3c73a8",
	},
	accentIcons: []rune{
		' ',
		'✤',
		'⎈',
		'⚙',
		'⊚',
		'✦',
		'❖',
		'✥',
		'❄',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#0a2622",
		grid.ACTION_LINE_RESET:     "#f76c5e",
		grid.ACTION_LINE_REVERSE:   "#f39c6b",
		grid.ACTION_LINE_SKIP_BEAT: "#90e0c9",
		grid.ACTION_RESET:          "#c5d86d",
		grid.ACTION_LINE_BOUNCE:    "#c5d86d",
		grid.ACTION_LINE_DELAY:     "#a682ff",
	},
}

var dynamiteColors = map[string]string{
	"AltSeqColor":                  "#300808",
	"SeqColor":                     "#1c0404",
	"SeqCursorColor":               "#580f0f",
	"SeqVisualColor":               "#ff5252",
	"SeqOverlayColor":              "#9e2424",
	"SeqMiddleOverlayColor":        "#c91a1a",
	"SelectedAttributeColor":       "#ffd700",
	"NumberColor":                  "#ff914d",
	"Black":                        "#1c0404",
	"White":                        "#fff8f0",
	"Heart":                        "#f44336",
	"ActiveRatchetColor":           "#ffeb3b",
	"MutedRatchetColor":            "#f44336",
	"CurrentPlayingColor":          "#ffeb3b",
	"ActivePlayingColor":           "#f44336",
	"ArrangementHeaderColor":       "#fff8f0",
	"ArrangementTitleColor":        "#fff8f0",
	"ArrangementGroupColor":        "#e91e63",
	"ArrangementIndentColor":       "#740505",
	"ArrangementSelectedLineColor": "#580f0f",
}

var dynamiteTheme = Theme{
	colors: dynamiteColors,
	accentColors: []string{
		"#1c0404",
		"#f44336",
		"#e91e63",
		"#ff914d",
		"#ffeb3b",
		"#ffd700",
		"#ff5722",
		"#d81b60",
		"#b71c1c",
	},
	accentIcons: []rune{
		' ',
		'✸',
		'✯',
		'☄',
		'✧',
		'⊚',
		'☢',
		'⚔',
		'❄',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#1c0404",
		grid.ACTION_LINE_RESET:     "#f44336",
		grid.ACTION_LINE_REVERSE:   "#ff914d",
		grid.ACTION_LINE_SKIP_BEAT: "#ffa000",
		grid.ACTION_RESET:          "#ffeb3b",
		grid.ACTION_LINE_BOUNCE:    "#ffd700",
		grid.ACTION_LINE_DELAY:     "#e91e63",
	},
}

var springtimeColors = map[string]string{
	"AltSeqColor":                  "#eef7e4",
	"SeqColor":                     "#f9fdf5",
	"SeqCursorColor":               "#d6e9c6",
	"SeqVisualColor":               "#8eb656",
	"SeqOverlayColor":              "#b3daff",
	"SeqMiddleOverlayColor":        "#a6d173",
	"SelectedAttributeColor":       "#9ed36a",
	"NumberColor":                  "#ff9eb3",
	"Black":                        "#3c4f2f",
	"White":                        "#f9fdf5",
	"Heart":                        "#fc6c85",
	"ActiveRatchetColor":           "#c3f584",
	"MutedRatchetColor":            "#fc6c85",
	"CurrentPlayingColor":          "#c3f584",
	"ActivePlayingColor":           "#fc6c85",
	"ArrangementHeaderColor":       "#f9fdf5",
	"ArrangementTitleColor":        "#f9fdf5",
	"ArrangementGroupColor":        "#ff9eb3",
	"ArrangementIndentColor":       "#a6d173",
	"ArrangementSelectedLineColor": "#d6e9c6",
}

var springtimeTheme = Theme{
	colors: springtimeColors,
	accentColors: []string{
		"#3c4f2f",
		"#fc6c85",
		"#ff9eb3",
		"#ffdb58",
		"#c3f584",
		"#9ed36a",
		"#7dcfb6",
		"#b3daff",
		"#93c4e6",
	},
	accentIcons: []rune{
		' ',
		'❀',
		'✿',
		'❁',
		'✾',
		'❃',
		'✤',
		'✽',
		'✻',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#3c4f2f",
		grid.ACTION_LINE_RESET:     "#fc6c85",
		grid.ACTION_LINE_REVERSE:   "#ffdb58",
		grid.ACTION_LINE_SKIP_BEAT: "#7dcfb6",
		grid.ACTION_RESET:          "#c3f584",
		grid.ACTION_LINE_BOUNCE:    "#9ed36a",
		grid.ACTION_LINE_DELAY:     "#ff9eb3",
	},
}

var orangegroveColors = map[string]string{
	"AltSeqColor":                  "#2d2418",
	"SeqColor":                     "#1a1410",
	"SeqCursorColor":               "#4f3f29",
	"SeqVisualColor":               "#cc9966",
	"SeqOverlayColor":              "#ff8c42",
	"SeqMiddleOverlayColor":        "#dd7733",
	"SelectedAttributeColor":       "#ffb347",
	"NumberColor":                  "#ffcb69",
	"Black":                        "#1a1410",
	"White":                        "#fff4e6",
	"Heart":                        "#f45d48",
	"ActiveRatchetColor":           "#f8c537",
	"MutedRatchetColor":            "#f45d48",
	"CurrentPlayingColor":          "#f8c537",
	"ActivePlayingColor":           "#f45d48",
	"ArrangementHeaderColor":       "#fff4e6",
	"ArrangementTitleColor":        "#fff4e6",
	"ArrangementGroupColor":        "#f45d48",
	"ArrangementIndentColor":       "#8c6d46",
	"ArrangementSelectedLineColor": "#4f3f29",
}

var orangegroveTheme = Theme{
	colors: orangegroveColors,
	accentColors: []string{
		"#1a1410",
		"#f45d48",
		"#ff8c42",
		"#ffb347",
		"#ffcb69",
		"#f8c537",
		"#e8871e",
		"#c05746",
		"#95533c",
	},
	accentIcons: []rune{
		' ',
		'☀',
		'♠',
		'♣',
		'❂',
		'✺',
		'♦',
		'♥',
		'♨',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#1a1410",
		grid.ACTION_LINE_RESET:     "#f45d48",
		grid.ACTION_LINE_REVERSE:   "#ff8c42",
		grid.ACTION_LINE_SKIP_BEAT: "#ffb347",
		grid.ACTION_RESET:          "#f8c537",
		grid.ACTION_LINE_BOUNCE:    "#ffcb69",
		grid.ACTION_LINE_DELAY:     "#c05746",
	},
}

var cyberpunkColors = map[string]string{
	"AltSeqColor":                  "#0d0221",
	"SeqColor":                     "#070215",
	"SeqCursorColor":               "#1a0933",
	"SeqVisualColor":               "#8e44ad",
	"SeqOverlayColor":              "#3e1f6d",
	"SeqMiddleOverlayColor":        "#2e0f50",
	"SelectedAttributeColor":       "#00ff9f",
	"NumberColor":                  "#ff00ff",
	"Black":                        "#070215",
	"White":                        "#edfdfd",
	"Heart":                        "#ff003c",
	"ActiveRatchetColor":           "#00ff9f",
	"MutedRatchetColor":            "#ff003c",
	"CurrentPlayingColor":          "#00ff9f",
	"ActivePlayingColor":           "#ff003c",
	"ArrangementHeaderColor":       "#edfdfd",
	"ArrangementTitleColor":        "#edfdfd",
	"ArrangementGroupColor":        "#ff00ff",
	"ArrangementIndentColor":       "#3e1f6d",
	"ArrangementSelectedLineColor": "#1a0933",
}

var cyberpunkTheme = Theme{
	colors: cyberpunkColors,
	accentColors: []string{
		"#070215",
		"#ff003c",
		"#ff00ff",
		"#f706cf",
		"#c100fd",
		"#00ff9f",
		"#00f0ff",
		"#0575e6",
		"#6023f8",
	},
	accentIcons: []rune{
		' ',
		'◉',
		'◎',
		'◍',
		'○',
		'Ω',
		'λ',
		'⚠',
		'☣',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#070215",
		grid.ACTION_LINE_RESET:     "#ff003c",
		grid.ACTION_LINE_REVERSE:   "#f706cf",
		grid.ACTION_LINE_SKIP_BEAT: "#00f0ff",
		grid.ACTION_RESET:          "#00ff9f",
		grid.ACTION_LINE_BOUNCE:    "#c100fd",
		grid.ACTION_LINE_DELAY:     "#ff00ff",
	},
}

var brainiacColors = map[string]string{
	"AltSeqColor":                  "#003b46",
	"SeqColor":                     "#002b33",
	"SeqCursorColor":               "#004d5d",
	"SeqVisualColor":               "#07889b",
	"SeqOverlayColor":              "#007888",
	"SeqMiddleOverlayColor":        "#006677",
	"SelectedAttributeColor":       "#66b3ba",
	"NumberColor":                  "#c4dfe6",
	"Black":                        "#002b33",
	"White":                        "#e8f1f2",
	"Heart":                        "#4fb99f",
	"ActiveRatchetColor":           "#66b3ba",
	"MutedRatchetColor":            "#4fb99f",
	"CurrentPlayingColor":          "#66b3ba",
	"ActivePlayingColor":           "#4fb99f",
	"ArrangementHeaderColor":       "#e8f1f2",
	"ArrangementTitleColor":        "#e8f1f2",
	"ArrangementGroupColor":        "#66b3ba",
	"ArrangementIndentColor":       "#007888",
	"ArrangementSelectedLineColor": "#004d5d",
}

var brainiacTheme = Theme{
	colors: brainiacColors,
	accentColors: []string{
		"#002b33",
		"#4fb99f",
		"#66b3ba",
		"#07889b",
		"#c4dfe6",
		"#8bd7d2",
		"#1b98a2",
		"#006494",
		"#065a82",
	},
	accentIcons: []rune{
		' ',
		'⌘',
		'⌥',
		'Φ',
		'Σ',
		'Π',
		'≡',
		'≈',
		'∞',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#002b33",
		grid.ACTION_LINE_RESET:     "#4fb99f",
		grid.ACTION_LINE_REVERSE:   "#07889b",
		grid.ACTION_LINE_SKIP_BEAT: "#1b98a2",
		grid.ACTION_RESET:          "#66b3ba",
		grid.ACTION_LINE_BOUNCE:    "#8bd7d2",
		grid.ACTION_LINE_DELAY:     "#c4dfe6",
	},
}

var spaceodysseyColors = map[string]string{
	"AltSeqColor":                  "#0c0c1d",
	"SeqColor":                     "#020210",
	"SeqCursorColor":               "#191936",
	"SeqVisualColor":               "#3a3a75",
	"SeqOverlayColor":              "#2d2d50",
	"SeqMiddleOverlayColor":        "#40406f",
	"SelectedAttributeColor":       "#c0c0f0",
	"NumberColor":                  "#f0c0c0",
	"Black":                        "#020210",
	"White":                        "#e6e6ff",
	"Heart":                        "#c050c0",
	"ActiveRatchetColor":           "#c0c0f0",
	"MutedRatchetColor":            "#c050c0",
	"CurrentPlayingColor":          "#c0c0f0",
	"ActivePlayingColor":           "#c050c0",
	"ArrangementHeaderColor":       "#e6e6ff",
	"ArrangementTitleColor":        "#e6e6ff",
	"ArrangementGroupColor":        "#c050c0",
	"ArrangementIndentColor":       "#40406f",
	"ArrangementSelectedLineColor": "#191936",
}

var spaceodysseyTheme = Theme{
	colors: spaceodysseyColors,
	accentColors: []string{
		"#020210",
		"#c050c0",
		"#9090e0",
		"#f0c0c0",
		"#c0c0f0",
		"#a0a0ff",
		"#8080c0",
		"#606090",
		"#404060",
	},
	accentIcons: []rune{
		' ',
		'★',
		'☽',
		'☼',
		'⋆',
		'✧',
		'☄',
		'✪',
		'⊛',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#020210",
		grid.ACTION_LINE_RESET:     "#c050c0",
		grid.ACTION_LINE_REVERSE:   "#9090e0",
		grid.ACTION_LINE_SKIP_BEAT: "#8080c0",
		grid.ACTION_RESET:          "#c0c0f0",
		grid.ACTION_LINE_BOUNCE:    "#a0a0ff",
		grid.ACTION_LINE_DELAY:     "#f0c0c0",
	},
}

var AccentColors = []string{
	"#000000",
	"#ed3902",
	"#f564a9",
	"#f8730e",
	"#fcc05c",
	"#5cdffb",
	"#1e89ef",
	"#164de5",
	"#0246a7",
}

var AccentIcons = []rune{
	' ',
	'✤',
	'⎈',
	'⚙',
	'⊚',
	'✦',
	'❖',
	'✥',
	'❄',
}

var ActionColors = map[grid.Action]string{
	grid.ACTION_NOTHING:        "#000000",
	grid.ACTION_LINE_RESET:     "#cf142b",
	grid.ACTION_LINE_REVERSE:   "#f8730e",
	grid.ACTION_LINE_SKIP_BEAT: "#a9e5bb",
	grid.ACTION_RESET:          "#fcf6b1",
	grid.ACTION_LINE_BOUNCE:    "#fcf6b1",
	grid.ACTION_LINE_DELAY:     "#cc4bc2",
}

// Colors

var AppTitleColor,
	AppDescriptorColor,
	AltSeqColor,
	SeqColor,
	SeqCursorColor,
	SeqVisualColor,
	SeqOverlayColor,
	SeqMiddleOverlayColor,
	SeqBorderLineColor,
	PatternModeColor,
	SelectedAttributeColor,
	NumberColor,
	Black,
	White,
	Heart,
	ActiveRatchetColor,
	MutedRatchetColor,
	CurrentPlayingColor,
	ActivePlayingColor,
	ArrangementHeaderColor,
	ArrangementTitleColor,
	ArrangementGroupColor,
	ArrangementIndentColor,
	ArrangementSelectedLineColor lipgloss.Color

// Styles
var AppTitleStyle,
	AppDescriptorStyle,
	ActiveStyle,
	MutedStyle,
	HeartStyle,
	SelectedStyle,
	NumberStyle,
	AccentModeStyle,
	BlackKeyStyle,
	WhiteKeyStyle,
	GroupStyle,
	IndentStyle,
	NodeRowStyle,
	SectionNameStyle,
	SeqBorderStyle lipgloss.Style

// Symbols
var CurrentlyPlayingSymbol,
	OverlayCurrentlyPlayingSymbol,
	ActiveSymbol string

var Themes = []string{
	"default",
	"seafoam",
	"dynamite",
	"springtime",
	"orangegrove",
	"cyberpunk",
	"brainiac",
	"spaceodyssey",
}

func ChooseTheme(colorscheme string) {
	switch colorscheme {
	case "default":
		ApplyTheme(defaultTheme)
	case "seafoam":
		ApplyTheme(seafoamTheme)
	case "dynamite":
		ApplyTheme(dynamiteTheme)
	case "springtime":
		ApplyTheme(springtimeTheme)
	case "orangegrove":
		ApplyTheme(orangegroveTheme)
	case "cyberpunk":
		ApplyTheme(cyberpunkTheme)
	case "brainiac":
		ApplyTheme(brainiacTheme)
	case "spaceodyssey":
		ApplyTheme(spaceodysseyTheme)
	}

	EvaluateStyles()
	EvaluateSymbols()
}

func ApplyTheme(theme Theme) {
	SetColors(theme.colors)
	SetAccentColors(theme.accentColors)
	SetAccentIcons(theme.accentIcons)
	SetActionColors(theme.lineActionColors)
}

func SetColors(newColors map[string]string) {
	for key, value := range newColors {
		newColor := lipgloss.Color(value)
		switch key {
		case "AppTitleColor":
			AppTitleColor = newColor
		case "AppDescriptorColor":
			AppDescriptorColor = newColor
		case "AltSeqColor":
			AltSeqColor = newColor
		case "SeqColor":
			SeqColor = newColor
		case "SeqCursorColor":
			SeqCursorColor = newColor
		case "SeqVisualColor":
			SeqVisualColor = newColor
		case "SeqOverlayColor":
			SeqOverlayColor = newColor
		case "SeqMiddleOverlayColor":
			SeqMiddleOverlayColor = newColor
		case "SelectedAttributeColor":
			SelectedAttributeColor = newColor
		case "NumberColor":
			NumberColor = newColor
		case "Black":
			Black = newColor
		case "White":
			White = newColor
		case "Heart":
			Heart = newColor
		case "ActiveRatchetColor":
			ActiveRatchetColor = newColor
		case "MutedRatchetColor":
			MutedRatchetColor = newColor
		case "CurrentPlayingColor":
			CurrentPlayingColor = newColor
		case "ActivePlayingColor":
			ActivePlayingColor = newColor
		case "ArrangementHeaderColor":
			ArrangementHeaderColor = newColor
		case "ArrangementTitleColor":
			ArrangementTitleColor = newColor
		case "ArrangementGroupColor":
			ArrangementGroupColor = newColor
		case "ArrangementIndentColor":
			ArrangementIndentColor = newColor
		case "ArrangementSelectedLineColor":
			ArrangementSelectedLineColor = newColor
		case "SeqBorderLineColor":
			SeqBorderLineColor = newColor
		case "PatternModeColor":
			PatternModeColor = newColor
		}
	}
}

func EvaluateSymbols() {
	CurrentlyPlayingSymbol = lipgloss.NewStyle().Foreground(CurrentPlayingColor).Render(" \u25CF ")
	OverlayCurrentlyPlayingSymbol = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(CurrentPlayingColor).Render(" \u25CF ")
	ActiveSymbol = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(ActivePlayingColor).Render(" \u25C9 ")
}

func EvaluateStyles() {

	ActiveStyle = lipgloss.NewStyle().Foreground(ActiveRatchetColor)
	MutedStyle = lipgloss.NewStyle().Foreground(MutedRatchetColor)
	HeartStyle = lipgloss.NewStyle().Foreground(Heart)
	SelectedStyle = lipgloss.NewStyle().Background(SelectedAttributeColor).Foreground(Black)
	NumberStyle = lipgloss.NewStyle().Foreground(NumberColor)
	AccentModeStyle = lipgloss.NewStyle().Background(PatternModeColor).Foreground(Black)
	BlackKeyStyle = lipgloss.NewStyle().Background(Black).Foreground(White)
	WhiteKeyStyle = lipgloss.NewStyle().Background(White).Foreground(Black)

	GroupStyle = lipgloss.NewStyle().
		Foreground(ArrangementGroupColor).
		MarginRight(1).
		Bold(true)

	IndentStyle = lipgloss.NewStyle().
		Foreground(ArrangementIndentColor)

	NodeRowStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		MarginBottom(0)

	SectionNameStyle = lipgloss.NewStyle().
		Foreground(ArrangementHeaderColor)

	SeqBorderStyle = lipgloss.NewStyle().Background(Black).Foreground(SeqBorderLineColor)
	AppTitleStyle = lipgloss.NewStyle().Background(Black).Foreground(AppTitleColor)
	AppDescriptorStyle = lipgloss.NewStyle().Background(Black).Foreground(AppDescriptorColor)
}

func SetAccentColors(accentColors []string) {
	for i := range len(accentColors) {
		if i == 0 {
			continue
		}
		AccentColors[i] = accentColors[i]
	}
}

func SetAccentIcons(accentIcons []rune) {
	for i := range len(accentIcons) {
		if i == 0 {
			continue
		}
		AccentIcons[i] = accentIcons[i]
	}
}

func SetActionColors(actionColors map[grid.Action]string) {
	for k := range actionColors {
		if k == grid.ACTION_NOTHING {
			continue
		}
		ActionColors[k] = actionColors[k]
	}
}
