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
	"AppTitleColor":                "#fafafa",
	"AppDescriptorColor":           "#fafafa",
	"LineNumberColor":              "#fafafa",
	"RightSideTitleColor":          "#fafafa",
	"AltSeqBackgroundColor":        "#2c2c2c",
	"SeqBackgroundColor":           "#000000",
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
	"SeqBorderLineColor":           "#fafafa",
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
	"AppTitleColor":                "#e6f2ef",
	"AppDescriptorColor":           "#cce7e0",
	"LineNumberColor":              "#a6d0c7",
	"RightSideTitleColor":          "#8abdb2",
	"AltSeqBackgroundColor":        "#013112",
	"SeqBackgroundColor":           "#0a2622",
	"SeqCursorColor":               "#3d8a7d",
	"SeqVisualColor":               "#88c3b8",
	"SeqOverlayColor":              "#336699",
	"SeqMiddleOverlayColor":        "#4d7a68",
	"SelectedAttributeColor":       "#ff9f45",
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
	"SeqBorderLineColor":           "#e6f2ef",
	"PatternModeColor":             "#f39c6b",
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
	"AppTitleColor":                "#fff8f0",
	"AppDescriptorColor":           "#ffe6d1",
	"LineNumberColor":              "#ffcbad",
	"RightSideTitleColor":          "#ffb38f",
	"AltSeqBackgroundColor":        "#3d0d0d",
	"SeqBackgroundColor":           "#1c0404",
	"SeqCursorColor":               "#7a1414",
	"SeqVisualColor":               "#ff5252",
	"SeqOverlayColor":              "#9e2424",
	"SeqMiddleOverlayColor":        "#c91a1a",
	"SelectedAttributeColor":       "#00e5ff",
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
	"SeqBorderLineColor":           "#fff8f0",
	"PatternModeColor":             "#ff914d",
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
	"AppTitleColor":                "#ffdb58",
	"AppDescriptorColor":           "#ff9eb3",
	"LineNumberColor":              "#617852",
	"RightSideTitleColor":          "#769164",
	"AltSeqBackgroundColor":        "#e5f0d8",
	"SeqBackgroundColor":           "#f9fdf5",
	"SeqCursorColor":               "#b7d49e",
	"SeqVisualColor":               "#8eb656",
	"SeqOverlayColor":              "#b3daff",
	"SeqMiddleOverlayColor":        "#a6d173",
	"SelectedAttributeColor":       "#ff5e8a",
	"NumberColor":                  "#ff9eb3",
	"Black":                        "#070215",
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
	"SeqBorderLineColor":           "#3c4f2f",
	"PatternModeColor":             "#ffdb58",
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
	"AppTitleColor":                "#fff4e6",
	"AppDescriptorColor":           "#ffe8cc",
	"LineNumberColor":              "#ffdbb3",
	"RightSideTitleColor":          "#ffc999",
	"AltSeqBackgroundColor":        "#3a2e1e",
	"SeqBackgroundColor":           "#1a1410",
	"SeqCursorColor":               "#6a5433",
	"SeqVisualColor":               "#cc9966",
	"SeqOverlayColor":              "#ff8c42",
	"SeqMiddleOverlayColor":        "#dd7733",
	"SelectedAttributeColor":       "#00c8ff",
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
	"SeqBorderLineColor":           "#fff4e6",
	"PatternModeColor":             "#ff8c42",
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
		'⊙',
		'⊚',
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
	"AppTitleColor":                "#edfdfd",
	"AppDescriptorColor":           "#c0f0f0",
	"LineNumberColor":              "#a0e0e0",
	"RightSideTitleColor":          "#80c0c0",
	"AltSeqBackgroundColor":        "#13032e",
	"SeqBackgroundColor":           "#070215",
	"SeqCursorColor":               "#2d555a",
	"SeqVisualColor":               "#8e44ad",
	"SeqOverlayColor":              "#3e1f6d",
	"SeqMiddleOverlayColor":        "#2e0f50",
	"SelectedAttributeColor":       "#fcee21",
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
	"SeqBorderLineColor":           "#edfdfd",
	"PatternModeColor":             "#ff00ff",
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
	"AppTitleColor":                "#e8f1f2",
	"AppDescriptorColor":           "#c4dfe6",
	"LineNumberColor":              "#a0c7d3",
	"RightSideTitleColor":          "#7cafbf",
	"AltSeqBackgroundColor":        "#142233",
	"SeqBackgroundColor":           "#002b00",
	"SeqCursorColor":               "#006b82",
	"SeqVisualColor":               "#07889b",
	"SeqOverlayColor":              "#007888",
	"SeqMiddleOverlayColor":        "#006677",
	"SelectedAttributeColor":       "#f8b24b",
	"NumberColor":                  "#c4dfe6",
	"Black":                        "#002b00",
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
	"SeqBorderLineColor":           "#e8f1f2",
	"PatternModeColor":             "#4fb99f",
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
		'∂',
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
	"AppTitleColor":                "#e6e6ff",
	"AppDescriptorColor":           "#c0c0f0",
	"LineNumberColor":              "#a0a0d6",
	"RightSideTitleColor":          "#8080c0",
	"AltSeqBackgroundColor":        "#14142a",
	"SeqBackgroundColor":           "#020210",
	"SeqCursorColor":               "#2a6655",
	"SeqVisualColor":               "#3a3a75",
	"SeqOverlayColor":              "#2d2d50",
	"SeqMiddleOverlayColor":        "#40406f",
	"SelectedAttributeColor":       "#ffcf00",
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
	"SeqBorderLineColor":           "#e6e6ff",
	"PatternModeColor":             "#9090e0",
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
		"#dd80c0",
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

var nineteenfiftyeightColors = map[string]string{
	"AppTitleColor":                "#ecdfce",
	"AppDescriptorColor":           "#d9c8af",
	"LineNumberColor":              "#b2a58c",
	"RightSideTitleColor":          "#a18f74",
	"AltSeqBackgroundColor":        "#2e2418",
	"SeqBackgroundColor":           "#1a1610",
	"SeqCursorColor":               "#5a4a3c",
	"SeqVisualColor":               "#917a64",
	"SeqOverlayColor":              "#6b5744",
	"SeqMiddleOverlayColor":        "#4a3c30",
	"SelectedAttributeColor":       "#4587be",
	"NumberColor":                  "#d09554",
	"Black":                        "#1a1610",
	"White":                        "#ecdfce",
	"Heart":                        "#c94a35",
	"ActiveRatchetColor":           "#b1b85a",
	"MutedRatchetColor":            "#c94a35",
	"CurrentPlayingColor":          "#b1b85a",
	"ActivePlayingColor":           "#c94a35",
	"ArrangementHeaderColor":       "#ecdfce",
	"ArrangementTitleColor":        "#ecdfce",
	"ArrangementGroupColor":        "#c94a35",
	"ArrangementIndentColor":       "#6b5744",
	"ArrangementSelectedLineColor": "#2e2418",
	"SeqBorderLineColor":           "#ecdfce",
	"PatternModeColor":             "#d09554",
}

var nineteenfiftyeightTheme = Theme{
	colors: nineteenfiftyeightColors,
	accentColors: []string{
		"#1a1610",
		"#c94a35",
		"#d09554",
		"#b1b85a",
		"#8ba353",
		"#4587be",
		"#376996",
		"#6b5744",
		"#483b2f",
	},
	accentIcons: []rune{
		' ',
		'◆',
		'■',
		'●',
		'▲',
		'◉',
		'◍',
		'◇',
		'◠',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#1a1610",
		grid.ACTION_LINE_RESET:     "#c94a35",
		grid.ACTION_LINE_REVERSE:   "#d09554",
		grid.ACTION_LINE_SKIP_BEAT: "#4587be",
		grid.ACTION_RESET:          "#b1b85a",
		grid.ACTION_LINE_BOUNCE:    "#8ba353",
		grid.ACTION_LINE_DELAY:     "#6b5744",
	},
}

var appleiiplusColors = map[string]string{
	"AppTitleColor":                "#33ff33",
	"AppDescriptorColor":           "#33dd33",
	"LineNumberColor":              "#33cc33",
	"RightSideTitleColor":          "#33bb33",
	"AltSeqBackgroundColor":        "#002200",
	"SeqBackgroundColor":           "#000000",
	"SeqCursorColor":               "#007700",
	"SeqVisualColor":               "#00aa00",
	"SeqOverlayColor":              "#006600",
	"SeqMiddleOverlayColor":        "#004400",
	"SelectedAttributeColor":       "#9933cc",
	"NumberColor":                  "#cc9933",
	"Black":                        "#000000",
	"White":                        "#33ff33",
	"Heart":                        "#cc3333",
	"ActiveRatchetColor":           "#66ff66",
	"MutedRatchetColor":            "#cc3333",
	"CurrentPlayingColor":          "#66ff66",
	"ActivePlayingColor":           "#cc3333",
	"ArrangementHeaderColor":       "#33ff33",
	"ArrangementTitleColor":        "#33ff33",
	"ArrangementGroupColor":        "#cc3333",
	"ArrangementIndentColor":       "#007700",
	"ArrangementSelectedLineColor": "#003300",
	"SeqBorderLineColor":           "#33ff33",
	"PatternModeColor":             "#cc9933",
}

var appleiiplusTheme = Theme{
	colors: appleiiplusColors,
	accentColors: []string{
		"#000000",
		"#cc3333",
		"#cc9933",
		"#66ff66",
		"#33cc33",
		"#9933cc",
		"#3333cc",
		"#007700",
		"#004400",
	},
	accentIcons: []rune{
		' ',
		'○',
		'◎',
		'●',
		'□',
		'◆',
		'▲',
		'△',
		'☐',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#000000",
		grid.ACTION_LINE_RESET:     "#cc3333",
		grid.ACTION_LINE_REVERSE:   "#cc9933",
		grid.ACTION_LINE_SKIP_BEAT: "#3333cc",
		grid.ACTION_RESET:          "#66ff66",
		grid.ACTION_LINE_BOUNCE:    "#33cc33",
		grid.ACTION_LINE_DELAY:     "#9933cc",
	},
}

var matrixColors = map[string]string{
	"AppTitleColor":                "#00ff00",
	"AppDescriptorColor":           "#00ee00",
	"LineNumberColor":              "#00dd00",
	"RightSideTitleColor":          "#00cc00",
	"AltSeqBackgroundColor":        "#002200",
	"SeqBackgroundColor":           "#000000",
	"SeqCursorColor":               "#00aa00",
	"SeqVisualColor":               "#008800",
	"SeqOverlayColor":              "#005500",
	"SeqMiddleOverlayColor":        "#002200",
	"SelectedAttributeColor":       "#ffffff",
	"NumberColor":                  "#88ff88",
	"Black":                        "#000000",
	"White":                        "#00ff00",
	"Heart":                        "#55ff55",
	"ActiveRatchetColor":           "#00ff00",
	"MutedRatchetColor":            "#008800",
	"CurrentPlayingColor":          "#00ff00",
	"ActivePlayingColor":           "#008800",
	"ArrangementHeaderColor":       "#00ff00",
	"ArrangementTitleColor":        "#00ff00",
	"ArrangementGroupColor":        "#55ff55",
	"ArrangementIndentColor":       "#004400",
	"ArrangementSelectedLineColor": "#003300",
	"SeqBorderLineColor":           "#00ff00",
	"PatternModeColor":             "#88ff88",
}

var matrixTheme = Theme{
	colors: matrixColors,
	accentColors: []string{
		"#000000",
		"#00ff00",
		"#00cc00",
		"#00aa00",
		"#008800",
		"#88ff88",
		"#55ff55",
		"#003300",
		"#001100",
	},
	accentIcons: []rune{
		' ',
		'0',
		'1',
		'Φ',
		'Ψ',
		'Δ',
		'Ω',
		'∑',
		'π',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#000000",
		grid.ACTION_LINE_RESET:     "#008800",
		grid.ACTION_LINE_REVERSE:   "#00aa00",
		grid.ACTION_LINE_SKIP_BEAT: "#55ff55",
		grid.ACTION_RESET:          "#00ff00",
		grid.ACTION_LINE_BOUNCE:    "#88ff88",
		grid.ACTION_LINE_DELAY:     "#00cc00",
	},
}

var herbieColors = map[string]string{
	"AppTitleColor":                "#ffe01b",
	"AppDescriptorColor":           "#ffd000",
	"LineNumberColor":              "#ffc000",
	"RightSideTitleColor":          "#ffb000",
	"AltSeqBackgroundColor":        "#1c1c1c",
	"SeqBackgroundColor":           "#0c0c0c",
	"SeqCursorColor":               "#3c3c3c",
	"SeqVisualColor":               "#666666",
	"SeqOverlayColor":              "#2c2c2c",
	"SeqMiddleOverlayColor":        "#242424",
	"SelectedAttributeColor":       "#1b95e0",
	"NumberColor":                  "#ffe01b",
	"Black":                        "#0c0c0c",
	"White":                        "#ffe01b",
	"Heart":                        "#e84a5f",
	"ActiveRatchetColor":           "#1b95e0",
	"MutedRatchetColor":            "#e84a5f",
	"CurrentPlayingColor":          "#1b95e0",
	"ActivePlayingColor":           "#e84a5f",
	"ArrangementHeaderColor":       "#ffe01b",
	"ArrangementTitleColor":        "#ffe01b",
	"ArrangementGroupColor":        "#e84a5f",
	"ArrangementIndentColor":       "#3c3c3c",
	"ArrangementSelectedLineColor": "#1c1c1c",
	"SeqBorderLineColor":           "#ffe01b",
	"PatternModeColor":             "#e84a5f",
}

var herbieTheme = Theme{
	colors: herbieColors,
	accentColors: []string{
		"#0c0c0c",
		"#e84a5f",
		"#ffe01b",
		"#feae5a",
		"#f9a03f",
		"#1b95e0",
		"#55b9f3",
		"#4c7f9e",
		"#2c3e50",
	},
	accentIcons: []rune{
		' ',
		'◉',
		'◈',
		'◇',
		'◎',
		'◔',
		'◑',
		'◕',
		'●',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#0c0c0c",
		grid.ACTION_LINE_RESET:     "#e84a5f",
		grid.ACTION_LINE_REVERSE:   "#feae5a",
		grid.ACTION_LINE_SKIP_BEAT: "#55b9f3",
		grid.ACTION_RESET:          "#ffe01b",
		grid.ACTION_LINE_BOUNCE:    "#f9a03f",
		grid.ACTION_LINE_DELAY:     "#4c7f9e",
	},
}

var milesColors = map[string]string{
	"AppTitleColor":                "#3498db",
	"AppDescriptorColor":           "#2980b9",
	"LineNumberColor":              "#1f6aa7",
	"RightSideTitleColor":          "#155a91",
	"AltSeqBackgroundColor":        "#0c2233",
	"SeqBackgroundColor":           "#000e1a",
	"SeqCursorColor":               "#1f4662",
	"SeqVisualColor":               "#2a4d6a",
	"SeqOverlayColor":              "#17313d",
	"SeqMiddleOverlayColor":        "#102837",
	"SelectedAttributeColor":       "#3498db",
	"NumberColor":                  "#f39c12",
	"Black":                        "#000e1a",
	"White":                        "#ecf0f1",
	"Heart":                        "#e74c3c",
	"ActiveRatchetColor":           "#2ecc71",
	"MutedRatchetColor":            "#e74c3c",
	"CurrentPlayingColor":          "#2ecc71",
	"ActivePlayingColor":           "#e74c3c",
	"ArrangementHeaderColor":       "#ecf0f1",
	"ArrangementTitleColor":        "#ecf0f1",
	"ArrangementGroupColor":        "#e74c3c",
	"ArrangementIndentColor":       "#17313d",
	"ArrangementSelectedLineColor": "#0c2233",
	"SeqBorderLineColor":           "#3498db",
	"PatternModeColor":             "#8e44ad",
}

var milesTheme = Theme{
	colors: milesColors,
	accentColors: []string{
		"#000e1a",
		"#e74c3c",
		"#8e44ad",
		"#f39c12",
		"#2ecc71",
		"#1abc9c",
		"#3498db",
		"#2980b9",
		"#0c2233",
	},
	accentIcons: []rune{
		' ',
		'▣',
		'▢',
		'▤',
		'▥',
		'▧',
		'▨',
		'▩',
		'◆',
	},
	lineActionColors: map[grid.Action]string{
		grid.ACTION_NOTHING:        "#000e1a",
		grid.ACTION_LINE_RESET:     "#e74c3c",
		grid.ACTION_LINE_REVERSE:   "#f39c12",
		grid.ACTION_LINE_SKIP_BEAT: "#3498db",
		grid.ACTION_RESET:          "#2ecc71",
		grid.ACTION_LINE_BOUNCE:    "#1abc9c",
		grid.ACTION_LINE_DELAY:     "#8e44ad",
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
	LineNumberColor,
	RightSideTitleColor,
	AltSeqBackgroundColor,
	SeqBackgroundColor,
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
	LineNumberStyle,
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
	"1958",
	"appleiiplus",
	"matrix",
	"herbie",
	"miles",
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
	case "1958":
		ApplyTheme(nineteenfiftyeightTheme)
	case "appleiiplus":
		ApplyTheme(appleiiplusTheme)
	case "matrix":
		ApplyTheme(matrixTheme)
	case "herbie":
		ApplyTheme(herbieTheme)
	case "miles":
		ApplyTheme(milesTheme)
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
		case "LineNumberColor":
			LineNumberColor = newColor
		case "RightSideTitleColor":
			RightSideTitleColor = newColor
		case "AltSeqBackgroundColor":
			AltSeqBackgroundColor = newColor
		case "SeqBackgroundColor":
			SeqBackgroundColor = newColor
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
	LineNumberStyle = lipgloss.NewStyle().Foreground(LineNumberColor)

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
	AppTitleStyle = lipgloss.NewStyle().Bold(true).Background(Black).Foreground(AppTitleColor)
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
