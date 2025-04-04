package colors

import "github.com/charmbracelet/lipgloss"

var AltSeqColor = lipgloss.Color("#222222")
var SeqColor = lipgloss.Color("#000000")
var SeqCursorColor = lipgloss.Color("#444444")
var SeqVisualColor = lipgloss.Color("#aaaaaa")
var SeqOverlayColor = lipgloss.Color("#333388")
var SeqMiddleOverlayColor = lipgloss.Color("#405810")
var SelectedAttributeColor = lipgloss.Color("#5cdffb")

var ActiveRatchetColor lipgloss.Color = "#abfaa9"
var MutedRatchetColor lipgloss.Color = "#f34213"

var ActiveStyle = lipgloss.NewStyle().Foreground(ActiveRatchetColor)
var MutedStyle = lipgloss.NewStyle().Foreground(MutedRatchetColor)

var HeartColor = lipgloss.NewStyle().Foreground(Heart)
var SelectedColor = lipgloss.NewStyle().Background(lipgloss.Color("#5cdffb")).Foreground(lipgloss.Color("#000000"))
var NumberColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcbd15"))

var currentPlayingColor lipgloss.Color = "#abfaa9"
var activePlayingColor lipgloss.Color = "#f34213"
var CurrentlyPlayingDot = lipgloss.NewStyle().Foreground(currentPlayingColor).Render(" \u25CF ")
var OverlayCurrentlyPlayingDot = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(currentPlayingColor).Render(" \u25CF ")
var ActiveDot = lipgloss.NewStyle().Background(SeqOverlayColor).Foreground(activePlayingColor).Render(" \u25C9 ")

var AccentModeStyle = lipgloss.NewStyle().Background(Heart).Foreground(lipgloss.Color("#000000"))

var BlackKeyStyle = lipgloss.NewStyle().Background(Black).Foreground(White)
var WhiteKeyStyle = lipgloss.NewStyle().Background(White).Foreground(Black)

var Black = lipgloss.Color("#000000")
var White = lipgloss.Color("#ffffff")
var Heart = lipgloss.Color("#ed3902")
