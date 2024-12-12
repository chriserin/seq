package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

type model struct {
	keys  keymap
	beats int
	tempo int
	help  help.Model
}

func InitModel() model {
	return model{
		keys:  keys,
		beats: 32,
		tempo: 240,
		help:  help.New(),
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
	buf.WriteString(m.help.View(m.keys))
	buf.WriteString("\n")
	return buf.String()
}
