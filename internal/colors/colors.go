package colors

import "github.com/charmbracelet/lipgloss"

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

var ActiveRatchetColor lipgloss.Color = "#abfaa9"
var MutedRatchetColor lipgloss.Color = "#f34213"

var currentPlayingColor lipgloss.Color = "#abfaa9"
var activePlayingColor lipgloss.Color = "#f34213"

var ArrangementHeaderColor = lipgloss.Color("FAFAFA")
var ArrangementTitleColor = lipgloss.Color("FAFAFA")
var ArrangementGroupColor = lipgloss.Color("#F25D94")
var ArrangementIndentColor = lipgloss.Color("#4b4261")
var ArrangementSelectedLineColor = lipgloss.Color("#3b4261")

var ActiveStyle = lipgloss.NewStyle().Foreground(ActiveRatchetColor)
var MutedStyle = lipgloss.NewStyle().Foreground(MutedRatchetColor)

var HeartStyle = lipgloss.NewStyle().Foreground(Heart)
var SelectedStyle = lipgloss.NewStyle().Background(SelectedAttributeColor).Foreground(Black)
var NumberStyle = lipgloss.NewStyle().Foreground(NumberColor)

var CurrentlyPlayingSymbol = lipgloss.NewStyle().Foreground(currentPlayingColor).Render(" \u25CF ")
var OverlayCurrentlyPlayingSymbol = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(currentPlayingColor).Render(" \u25CF ")
var ActiveSymbol = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(activePlayingColor).Render(" \u25C9 ")

var AccentModeStyle = lipgloss.NewStyle().Background(Heart).Foreground(Black)

var BlackKeyStyle = lipgloss.NewStyle().Background(Black).Foreground(White)
var WhiteKeyStyle = lipgloss.NewStyle().Background(White).Foreground(Black)
