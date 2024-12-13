package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	Quit          key.Binding
	Help          key.Binding
	CursorUp      key.Binding
	CursorDown    key.Binding
	CursorLeft    key.Binding
	CursorRight   key.Binding
	TriggerAdd    key.Binding
	TriggerRemove key.Binding
}

func Key(keyboardKey string, help string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey), key.WithHelp(keyboardKey, help))
}

var keys = keymap{
	Quit:          Key("q", "Quit"),
	Help:          Key("?", "Expand Help"),
	CursorUp:      Key("k", "Up"),
	CursorDown:    Key("j", "Down"),
	CursorLeft:    Key("h", "Left"),
	CursorRight:   Key("l", "Right"),
	TriggerAdd:    Key("f", "Add Trigger"),
	TriggerRemove: Key("d", "Remove Trigger"),
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help, k.Quit,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
		{k.CursorUp, k.CursorDown, k.CursorLeft, k.CursorRight},
		{k.TriggerAdd, k.TriggerRemove},
	}
}

const BLANK = " "
const TRIGGER = '■'

var zeroRune rune

type line []rune

type CursorPosition struct {
	lineNumber int
	beat       int
}

type model struct {
	keys      keymap
	beats     int
	tempo     int
	help      help.Model
	lines     []line
	cursorPos CursorPosition
	cursor    cursor.Model
}

func (m *model) AddTrigger() {
	m.lines[m.cursorPos.lineNumber][m.cursorPos.beat] = TRIGGER
}

func (m *model) RemoveTrigger() {
	m.lines[m.cursorPos.lineNumber][m.cursorPos.beat] = zeroRune
}

func InitSeq(lineNumber int, beatNumber int) []line {
	var lines = make([]line, 0, lineNumber)

	for i := 0; i < lineNumber; i++ {
		lines = append(lines, make([]rune, beatNumber))
	}
	return lines
}

func InitModel() model {
	newCursor := cursor.New()
	newCursor.Style = lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	return model{
		keys:      keys,
		beats:     32,
		tempo:     240,
		help:      help.New(),
		lines:     InitSeq(8, 32),
		cursorPos: CursorPosition{lineNumber: 0, beat: 0},
		cursor:    newCursor,
	}
}

func RunProgram() *tea.Program {
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	return p
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return tea.FocusMsg{} }
}

func Is(msg tea.KeyMsg, k key.Binding) bool {
	return key.Matches(msg, k)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case Is(msg, m.keys.Quit):
			return m, tea.Quit
		case Is(msg, m.keys.CursorDown):
			if m.cursorPos.lineNumber < len(m.lines)-1 {
				m.cursorPos.lineNumber++
			}
		case Is(msg, m.keys.CursorUp):
			if m.cursorPos.lineNumber > 0 {
				m.cursorPos.lineNumber--
			}
		case Is(msg, m.keys.CursorLeft):
			if m.cursorPos.beat > 0 {
				m.cursorPos.beat--
			}
		case Is(msg, m.keys.CursorRight):
			if m.cursorPos.beat < m.beats-1 {
				m.cursorPos.beat++
			}
		case Is(msg, m.keys.TriggerAdd):
			m.AddTrigger()
		case Is(msg, m.keys.TriggerRemove):
			m.RemoveTrigger()
		}
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor
	return m, cmd
}

func (m model) View() string {
	var buf strings.Builder
	buf.WriteString("   Seq - A sequencer for your cli\n")
	buf.WriteString("  ┌────────────────────────────────────\n")
	for i, line := range m.lines {
		buf.WriteString(line.View(i, m))
	}
	buf.WriteString(m.help.View(m.keys))
	buf.WriteString("\n")
	return buf.String()
}

var altSeqColor = lipgloss.NewStyle().Background(lipgloss.Color("#222222"))

func (line line) View(lineNumber int, m model) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d │", lineNumber))

	for i := 0; i < m.beats; i++ {

		var char string
		if line[i] == zeroRune {
			char = BLANK
		} else {
			char = string(line[i])
		}
		if m.cursorPos.lineNumber == lineNumber && m.cursorPos.beat == i {
			m.cursor.SetChar(char)
			char = m.cursor.View()
		}

		if i%8 > 3 {
			buf.WriteString(altSeqColor.Render(char))
		} else {
			buf.WriteString(char)
		}
	}

	buf.WriteString("\n")
	return buf.String()
}
