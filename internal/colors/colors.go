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

var AltSeqColor = lipgloss.Color("#222222")
var SeqColor = lipgloss.Color("#000000")
var SeqCursorColor = lipgloss.Color("#444444")
var SeqVisualColor = lipgloss.Color("#aaaaaa")
var SeqOverlayColor = lipgloss.Color("#333388")
var SeqMiddleOverlayColor = lipgloss.Color("#405810")
var SelectedAttributeColor = lipgloss.Color("#5cdffb")
var NumberColor = lipgloss.Color("#fcbd15")
var Black = lipgloss.Color("#000000")
var White = lipgloss.Color("#ffffff")
var Heart = lipgloss.Color("#ed3902")
var ActiveRatchetColor = lipgloss.Color("#abfaa9")
var MutedRatchetColor = lipgloss.Color("#f34213")
var CurrentPlayingColor = lipgloss.Color("#abfaa9")
var ActivePlayingColor = lipgloss.Color("#f34213")
var ArrangementHeaderColor = lipgloss.Color("FAFAFA")
var ArrangementTitleColor = lipgloss.Color("FAFAFA")
var ArrangementGroupColor = lipgloss.Color("#F25D94")
var ArrangementIndentColor = lipgloss.Color("#4b4261")
var ArrangementSelectedLineColor = lipgloss.Color("#3b4261")

// Styles
var ActiveStyle = lipgloss.NewStyle().Foreground(ActiveRatchetColor)
var MutedStyle = lipgloss.NewStyle().Foreground(MutedRatchetColor)
var HeartStyle = lipgloss.NewStyle().Foreground(Heart)
var SelectedStyle = lipgloss.NewStyle().Background(SelectedAttributeColor).Foreground(Black)
var NumberStyle = lipgloss.NewStyle().Foreground(NumberColor)
var AccentModeStyle = lipgloss.NewStyle().Background(Heart).Foreground(Black)
var BlackKeyStyle = lipgloss.NewStyle().Background(Black).Foreground(White)
var WhiteKeyStyle = lipgloss.NewStyle().Background(White).Foreground(Black)

// Symbols
var CurrentlyPlayingSymbol = lipgloss.NewStyle().Foreground(CurrentPlayingColor).Render(" \u25CF ")
var OverlayCurrentlyPlayingSymbol = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(CurrentPlayingColor).Render(" \u25CF ")
var ActiveSymbol = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(ActivePlayingColor).Render(" \u25C9 ")

func ChooseColorScheme(colorscheme string) {
	switch colorscheme {
	case "default":
		ApplyTheme(defaultTheme)
	case "seafoam":
	case "dynamite":
	case "springtime":
	case "orangegrove":
	case "cyberpunk":
	case "brainiac":
	case "spaceodyssey":
	}
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
		}
	}
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
