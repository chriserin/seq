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

type Model struct {
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
		buf.WriteString(fmt.Sprintf("%*s", 15, m.SectionOutput(i, currentSongSection, songSection)))
		buf.WriteString(fmt.Sprintf("%*s", 15, m.StartBeatsOutput(i)))
		buf.WriteString(fmt.Sprintf("%*s", 15, m.StartCyclesOutput(i)))
		buf.WriteString(fmt.Sprintf("%*s", 15, m.CyclesOutput(i)))
		buf.WriteString("\n")
	}
	return buf.String()
}

func (m Model) CyclesOutput(index int) string {
	return fmt.Sprintf("%d", (*m.arrangement)[index].Cycles)
}

func (m Model) StartBeatsOutput(index int) string {
	return fmt.Sprintf("%d", (*m.arrangement)[index].StartBeat)
}

func (m Model) StartCyclesOutput(index int) string {
	return fmt.Sprintf("%d", (*m.arrangement)[index].StartCycles)
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
