package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keymap struct {
	Quit key.Binding
	Help key.Binding
}

var keys = keymap{
	Quit: key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit")),
	Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "Expand Help")),
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help, k.Quit,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
	}
}

const BLANK = " "

var zeroRune rune

type line []rune

type model struct {
	keys  keymap
	beats int
	tempo int
	help  help.Model
	lines []line
}

func InitSeq(lineNumber int, beatNumber int) []line {
	var lines = make([]line, 0, lineNumber)

	for i := 0; i < lineNumber; i++ {
		lines = append(lines, make([]rune, beatNumber))
	}
	return lines
}

func InitModel() model {
	return model{
		keys:  keys,
		beats: 32,
		tempo: 240,
		help:  help.New(),
		lines: InitSeq(8, 32),
	}
}

func RunProgram() *tea.Program {
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	return p
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var buf strings.Builder
	buf.WriteString("   Seq - A sequencer for your cli\n")
	buf.WriteString("  ┌────────────────────────────────────\n")
	for i, line := range m.lines {
		buf.WriteString(Line(i, line, m))
	}
	buf.WriteString(m.help.View(m.keys))
	buf.WriteString("\n")
	return buf.String()
}

var altSeqColor = lipgloss.NewStyle().Background(lipgloss.Color("#222222"))

func Line(lineNumber int, line line, m model) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d │", lineNumber))

	for i := 0; i < m.beats; i++ {

		var char string
		if line[i] == zeroRune {
			char = BLANK
		} else {
			char = string(line[i])
		}
		// if m.cursorPos.lineNumber == lineNumber && m.cursorPos.beat == i {
		// 	m.cursor.SetChar(char)
		// 	char = m.cursor.View()
		// }

		if i%8 > 3 {
			buf.WriteString(altSeqColor.Render(char))
		} else {
			buf.WriteString(char)
		}
	}

	buf.WriteString("\n")
	return buf.String()
}
