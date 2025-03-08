package arrangement

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chriserin/seq/internal/colors"
	"github.com/chriserin/seq/internal/overlays"
)

type cursor struct {
	section   int
	attribute SectionAttribute
}

func (c cursor) Matches(section int, attribute SectionAttribute) bool {
	return c.section == section && c.attribute == attribute
}

type SectionAttribute int

const (
	SECTION_START_BEAT SectionAttribute = iota
	SECTION_START_CYCLE
	SECTION_CYCLES
)

type Model struct {
	Focus       bool
	cursor      cursor
	arrangement *[]SongSection
	parts       *[]Part
}

func InitModel(arrangement *[]SongSection, parts *[]Part) Model {
	return Model{
		arrangement: arrangement,
		parts:       parts,
	}
}

type keymap struct {
	CursorUp    key.Binding
	CursorDown  key.Binding
	CursorLeft  key.Binding
	CursorRight key.Binding
}

var keys = keymap{
	CursorUp:    Key("Up", "k"),
	CursorDown:  Key("Down", "j"),
	CursorLeft:  Key("Left", "h"),
	CursorRight: Key("Right", "l"),
}

func Key(help string, keyboardKey ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey...), key.WithHelp(keyboardKey[0], help))
}

func (m Model) Update(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case Is(msg, keys.CursorDown):
		if m.cursor.section+1 < len(*m.arrangement) {
			m.cursor.section++
		}
	case Is(msg, keys.CursorUp):
		if m.cursor.section > 0 {
			m.cursor.section--
		}
	case Is(msg, keys.CursorLeft):
		if m.cursor.attribute > 0 {
			m.cursor.attribute--
		}
	case Is(msg, keys.CursorRight):
		if m.cursor.attribute < 2 {
			m.cursor.attribute++
		}
	}
	return m, nil
}

func Is(msg tea.KeyMsg, k key.Binding) bool {
	return key.Matches(msg, k)
}

type SongSection struct {
	Part        int
	Cycles      int
	StartBeat   int
	StartCycles int
}

type Part struct {
	Overlays *overlays.Overlay
	Beats    uint8
	Name     string
}

func (m Model) View(currentSongSection int) string {
	var buf strings.Builder
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Part Name"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Start Beat"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Start Cycle"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Cycles"))
	buf.WriteString("\n")
	for i, songSection := range *m.arrangement {
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, m.SectionOutput(i, currentSongSection, songSection)))
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, m.StartBeatsOutput(i, m.cursor)))
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, m.StartCyclesOutput(i, m.cursor)))
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, m.CyclesOutput(i, m.cursor)))
		buf.WriteString("\n")
	}
	return buf.String()
}

func (m Model) CyclesOutput(index int, cursor cursor) string {
	cycles := (*m.arrangement)[index].Cycles
	if m.Focus && cursor.Matches(index, SECTION_CYCLES) {
		return colors.SelectedColor.Render(fmt.Sprintf("%d", cycles))
	} else {
		return colors.NumberColor.Render(fmt.Sprintf("%d", cycles))
	}
}

func (m Model) StartBeatsOutput(index int, cursor cursor) string {
	startBeat := (*m.arrangement)[index].StartBeat
	if m.Focus && cursor.Matches(index, SECTION_START_BEAT) {
		return colors.SelectedColor.Render(fmt.Sprintf("%d", startBeat))
	} else {
		return colors.NumberColor.Render(fmt.Sprintf("%d", startBeat))
	}
}

func (m Model) StartCyclesOutput(index int, cursor cursor) string {
	startCycle := (*m.arrangement)[index].StartBeat
	if m.Focus && cursor.Matches(index, SECTION_START_CYCLE) {
		return colors.SelectedColor.Render(fmt.Sprintf("%d", startCycle))
	} else {
		return colors.NumberColor.Render(fmt.Sprintf("%d", startCycle))
	}
}

func (m Model) SectionOutput(index int, currentSongSection int, songSection SongSection) string {
	section := fmt.Sprintf("%d) %s", index+1, (*m.parts)[songSection.Part].GetName())
	var sectionOutput string
	if currentSongSection == index {
		sectionOutput = fmt.Sprintf("   %s%s", colors.CurrentlyPlayingDot, section)
	} else {
		sectionOutput = fmt.Sprintf("%s %s", "   ", section)
	}
	return sectionOutput
}

func (p Part) GetName() string {
	return p.Name
}
