package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv"
)

type keymap struct {
	Quit             key.Binding
	Help             key.Binding
	CursorUp         key.Binding
	CursorDown       key.Binding
	CursorLeft       key.Binding
	CursorRight      key.Binding
	TriggerAdd       key.Binding
	TriggerRemove    key.Binding
	PlayStop         key.Binding
	ClearLine        key.Binding
	ClearSeq         key.Binding
	TempoInputSwitch key.Binding
	TempoIncrease    key.Binding
	TempoDecrease    key.Binding
}

func Key(keyboardKey string, help string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey), key.WithHelp(keyboardKey, help))
}

var keys = keymap{
	Quit:             Key("q", "Quit"),
	Help:             Key("?", "Expand Help"),
	CursorUp:         Key("k", "Up"),
	CursorDown:       Key("j", "Down"),
	CursorLeft:       Key("h", "Left"),
	CursorRight:      Key("l", "Right"),
	TriggerAdd:       Key("f", "Add Trigger"),
	TriggerRemove:    Key("d", "Remove Trigger"),
	ClearLine:        Key("c", "Clear Line"),
	ClearSeq:         Key("C", "Clear Seq"),
	PlayStop:         Key(" ", "Play/Stop"),
	TempoInputSwitch: Key("T", "Select Tempo Indicator"),
	TempoIncrease:    Key("+", "Tempo Increase"),
	TempoDecrease:    Key("-", "Tempo Decrease"),
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

const C1 = 36
const BLANK = " "
const TRIGGER = '■'

var zeroRune rune

type line []rune

type CursorPosition struct {
	lineNumber int
	beat       int
}

type model struct {
	keys                    keymap
	beats                   int
	tempo                   int
	subdivisions            int
	help                    help.Model
	lines                   []line
	cursorPos               CursorPosition
	cursor                  cursor.Model
	outport                 drivers.Out
	playing                 bool
	playTime                time.Time
	totalBeats              int
	currentBeat             int
	tempoSelectionIndicator uint8
}

type beatMsg struct{}

func BeatTick(playTime time.Time, totalBeats int, tempo int, subdivisions int) tea.Cmd {
	tickInterval := time.Minute / time.Duration(tempo*subdivisions)
	adjuster := time.Since(playTime) - (time.Duration(totalBeats) * tickInterval)
	next := tickInterval - adjuster
	return tea.Tick(
		next,
		func(t time.Time) tea.Msg { return beatMsg{} },
	)
}

func PlayBeat(lines []line, currentBeat int, sendFn func(msg midi.Message) error) tea.Cmd {
	return func() tea.Msg {
		Play(lines, currentBeat, sendFn)
		return nil
	}
}

func Play(lines []line, currentBeat int, sendFn func(msg midi.Message) error) {
	for i, line := range lines {
		spot := line[currentBeat]
		if spot != zeroRune {
			err := sendFn(midi.NoteOn(10, C1+uint8(i), 100))
			if err != nil {
				panic("note on failed")
			}
			err = sendFn(midi.NoteOff(10, C1+uint8(i)))
			if err != nil {
				panic("note off failed")
			}
		}
	}
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

	outport, err := midi.OutPort(0)
	if err != nil {
		panic("Did not get midi outport")
	}

	return model{
		keys:         keys,
		beats:        32,
		tempo:        120,
		subdivisions: 2,
		help:         help.New(),
		lines:        InitSeq(8, 32),
		cursorPos:    CursorPosition{lineNumber: 0, beat: 0},
		cursor:       newCursor,
		outport:      outport,
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
		case Is(msg, m.keys.ClearLine):
			zeroLine(m.lines[m.cursorPos.lineNumber])
		case Is(msg, m.keys.ClearSeq):
			m.lines = InitSeq(8, 32)

		case Is(msg, m.keys.PlayStop):
			if !m.playing && !m.outport.IsOpen() {
				err := m.outport.Open()
				if err != nil {
					panic("It's not open!")
				}
			}

			if m.playing && m.outport.IsOpen() {
				m.outport.Close()
			}

			m.playing = !m.playing
			m.playTime = time.Now()
			if m.playing {
				m.totalBeats = 0
				m.currentBeat = 0
				sendFn, err := midi.SendTo(m.outport)
				if err != nil {
					panic("sendFn is broken")
				}
				return m, tea.Batch(PlayBeat(m.lines, m.currentBeat, sendFn), BeatTick(m.playTime, m.totalBeats, m.tempo, m.subdivisions))
			}
		case Is(msg, m.keys.TempoInputSwitch):
			m.tempoSelectionIndicator = (m.tempoSelectionIndicator + 1) % 3
		case Is(msg, m.keys.TempoIncrease):
			if m.tempoSelectionIndicator == 1 {
				if m.tempo < 300 {
					m.tempo++
				}
			} else if m.tempoSelectionIndicator == 2 {
				if m.subdivisions < 8 {
					m.subdivisions++
				}
			}
		case Is(msg, m.keys.TempoDecrease):
			if m.tempoSelectionIndicator == 1 {
				if m.tempo > 30 {
					m.tempo--
				}
			} else if m.tempoSelectionIndicator == 2 {
				if m.subdivisions > 1 {
					m.subdivisions--
				}
			}
		}
		if msg.String() >= "1" && msg.String() <= "9" {
			beatInterval, _ := strconv.Atoi(msg.String())
			m.lines[m.cursorPos.lineNumber] = fill(m.lines[m.cursorPos.lineNumber], m.cursorPos.beat, beatInterval)
		}
	case beatMsg:
		if m.playing {
			m.currentBeat = (m.currentBeat + 1) % m.beats
			m.totalBeats++
			sendFn, err := midi.SendTo(m.outport)
			if err != nil {
				panic("sendFn is broken")
			}
			return m, tea.Batch(PlayBeat(m.lines, m.currentBeat, sendFn), BeatTick(m.playTime, m.totalBeats, m.tempo, m.subdivisions))
		}
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor
	return m, cmd
}

func zeroLine(beatline line) {
	for i := range beatline {
		beatline[i] = zeroRune
	}
}

func fill(beatline line, start int, every int) line {
	for i := range beatline[start:] {
		if i%every == 0 {
			if beatline[start+i] == TRIGGER {
				beatline[start+i] = zeroRune
			} else {
				beatline[start+i] = TRIGGER
			}
		}
	}
	return beatline
}

var heartColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed3902"))
var selectedColor = lipgloss.NewStyle().Background(lipgloss.Color("#5cdffb")).Foreground(lipgloss.Color("#000000"))
var numberColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcbd15"))

func (m model) TempoView() string {
	var buf strings.Builder
	var tempo, division string
	switch m.tempoSelectionIndicator {
	case 0:
		tempo = numberColor.Render(strconv.Itoa(m.tempo))
		division = numberColor.Render(strconv.Itoa(m.subdivisions))
	case 1:
		tempo = selectedColor.Render(strconv.Itoa(m.tempo))
		division = numberColor.Render(strconv.Itoa(m.subdivisions))
	case 2:
		tempo = numberColor.Render(strconv.Itoa(m.tempo))
		division = selectedColor.Render(strconv.Itoa(m.subdivisions))
	}
	heart := heartColor.Render("♡")
	buf.WriteString(heartColor.Render("             ") + "\n")
	buf.WriteString(heartColor.Render("   ♡♡♡☆ ☆♡♡♡ ") + "\n")
	buf.WriteString(heartColor.Render("  ♡    ◊    ♡") + "\n")
	buf.WriteString(heartColor.Render("  ♡  TEMPO  ♡") + "\n")
	buf.WriteString(fmt.Sprintf("  %s   %s   %s\n", heart, tempo, heart))
	buf.WriteString(heartColor.Render("   ♡ BEATS ♡") + "\n")
	buf.WriteString(fmt.Sprintf("    %s  %s  %s  \n", heart, division, heart))
	buf.WriteString(heartColor.Render("     ♡   ♡   ") + "\n")
	buf.WriteString(heartColor.Render("      ♡ ♡    ") + "\n")
	buf.WriteString(heartColor.Render("       †     ") + "\n")
	return buf.String()
}

func (m model) View() string {
	var buf strings.Builder
	buf.WriteString(lipgloss.JoinHorizontal(0, m.TempoView(), "  ", m.ViewTriggerSeq()))
	return buf.String()
}

func (m model) ViewTriggerSeq() string {
	var buf strings.Builder
	buf.WriteString("   Seq - A sequencer for your cli\n")
	buf.WriteString("  ┌────────────────────────────────────\n")
	for i, line := range m.lines {
		buf.WriteString(line.View(i, m))
	}
	if m.playing {
		buf.WriteString(fmt.Sprintf("   %*s%s", m.currentBeat, "", "█"))
	}
	buf.WriteString("\n")
	// buf.WriteString(m.help.View(m.keys))
	// buf.WriteString("\n")
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
