package main

import (
	"fmt"
	"os"
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
	Quit                 key.Binding
	Help                 key.Binding
	CursorUp             key.Binding
	CursorDown           key.Binding
	CursorLeft           key.Binding
	CursorRight          key.Binding
	TriggerAdd           key.Binding
	TriggerRemove        key.Binding
	PlayStop             key.Binding
	ClearLine            key.Binding
	ClearSeq             key.Binding
	TempoInputSwitch     key.Binding
	TempoIncrease        key.Binding
	TempoDecrease        key.Binding
	ToggleAccentMode     key.Binding
	ToggleAccentModifier key.Binding
	RatchetIncrease      key.Binding
	RatchetDecrease      key.Binding
	ActionAddLineReset   key.Binding
	ActionAddLineReverse key.Binding
	OverlayInputSwitch   key.Binding
}

func Key(keyboardKey string, help string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey), key.WithHelp(keyboardKey, help))
}

var keys = keymap{
	Quit:                 Key("q", "Quit"),
	Help:                 Key("?", "Expand Help"),
	CursorUp:             Key("k", "Up"),
	CursorDown:           Key("j", "Down"),
	CursorLeft:           Key("h", "Left"),
	CursorRight:          Key("l", "Right"),
	TriggerAdd:           Key("f", "Add Trigger"),
	TriggerRemove:        Key("d", "Remove Trigger"),
	ClearLine:            Key("c", "Clear Line"),
	ClearSeq:             Key("C", "Clear Seq"),
	PlayStop:             Key(" ", "Play/Stop"),
	TempoInputSwitch:     Key("T", "Select Tempo Indicator"),
	TempoIncrease:        Key("+", "Tempo Increase"),
	TempoDecrease:        Key("-", "Tempo Decrease"),
	ToggleAccentMode:     Key("A", "Toggle Accent Mode"),
	ToggleAccentModifier: Key("a", "Toggle Accent Modifier"),
	RatchetIncrease:      Key("R", "Increase Ratchet"),
	RatchetDecrease:      Key("r", "Decrease Ratchet"),
	ActionAddLineReset:   Key("s", "Add Line Reset Action"),
	ActionAddLineReverse: Key("S", "Add Line Reverse Action"),
	OverlayInputSwitch:   Key("O", "Select Overlay Indicator"),
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

type Accent struct {
	shape rune
	color lipgloss.Color
	value uint8
}

var accents = []Accent{
	{' ', "#000000", 0},
	{'✤', "#ed3902", 120},
	{'⎈', "#f564a9", 105},
	{'⚙', "#f8730e", 90},
	{'⊚', "#fcc05c", 75},
	{'✦', "#5cdffb", 60},
	{'❖', "#1e89ef", 45},
	{'✥', "#164de5", 30},
	{'❄', "#0246a7", 15},
}

const C1 = 36

type note struct {
	accentIndex  uint8
	ratchetIndex uint8
	action       action
}

type action uint8

type lineaction struct {
	shape rune
	color lipgloss.Color
}

const (
	ACTION_NOTHING action = iota
	ACTION_LINE_RESET
	ACTION_LINE_REVERSE
)

var lineactions = map[action]lineaction{
	ACTION_NOTHING:      {' ', "#000000"},
	ACTION_LINE_RESET:   {'↔', "#cf142b"},
	ACTION_LINE_REVERSE: {'←', "#f8730e"},
}

type ratchet string

var ratchets = []ratchet{
	"",
	"\u0307",
	"\u030A",
	"\u030B",
	"\u030C",
	"\u0312",
	"\u0313",
	"\u0344",
}

var zeronote note

type line struct {
	notes []note
}

type ZeroBehavior uint8

const (
	NORMAL ZeroBehavior = iota
	LOOP
)

type linestate struct {
	currentBeat    uint8
	direction      int8
	resetDirection int8
	resetLocation  uint8
}

type overlayKey struct {
	num   uint8
	denom uint8
}

var ROOT_OVERLAY = overlayKey{1, 1}

func (o *overlayKey) IncrementNumerator() {
	if o.num < 9 {
		o.num++
	}
}

func (o *overlayKey) IncrementDenominator() {
	if o.denom < 9 {
		o.denom++
	}
}

func (o *overlayKey) DecrementNumerator() {
	if o.num > 1 {
		o.num--
	}
}

func (o *overlayKey) DecrementDenominator() {
	if o.denom > 1 {
		o.denom--
	}
}

type overlays map[overlayKey]overlay

type gridKey struct {
	line uint8
	beat uint8
}

type overlay map[gridKey]note

type model struct {
	keys                      keymap
	beats                     uint8
	tempo                     int
	subdivisions              int
	help                      help.Model
	lines                     []line
	cursorPos                 gridKey
	cursor                    cursor.Model
	outport                   drivers.Out
	playing                   bool
	playTime                  time.Time
	trackTime                 time.Duration
	totalBeats                int
	playState                 []linestate
	tempoSelectionIndicator   uint8
	overlaySelectionIndicator uint8
	accentMode                bool
	accentModifier            int8
	logFile                   *os.File
	overlayKey                overlayKey
	overlays                  overlays
}

type beatMsg struct {
	interval time.Duration
}

func BeatTick(beatInterval time.Duration) tea.Cmd {
	return tea.Tick(
		beatInterval,
		func(t time.Time) tea.Msg { return beatMsg{beatInterval} },
	)
}

func (m *model) BeatInterval() time.Duration {
	tickInterval := time.Minute / time.Duration(m.tempo*m.subdivisions)
	adjuster := time.Since(m.playTime) - m.trackTime
	m.trackTime = m.trackTime + tickInterval
	next := tickInterval - adjuster
	return next
}

func PlayBeat(beatInterval time.Duration, lines []line, currentBeat []linestate, sendFn SendFunc) tea.Cmd {
	return func() tea.Msg {
		Play(beatInterval, lines, currentBeat, sendFn)
		return nil
	}
}

type SendFunc func(msg midi.Message) error

func Play(beatInterval time.Duration, lines []line, currentBeat []linestate, sendFn SendFunc) {
	for i, line := range lines {
		spot := line.notes[currentBeat[i].currentBeat]
		if spot != zeronote {
			onMessage := midi.NoteOn(10, C1+uint8(i), accents[spot.accentIndex].value)
			offMessage := midi.NoteOff(10, C1+uint8(i))
			err := sendFn(onMessage)
			if err != nil {
				panic("note on failed")
			}
			err = sendFn(offMessage)
			if err != nil {
				panic("note off failed")
			}
			if spot.ratchetIndex > 0 {
				ratchetInterval := beatInterval / time.Duration(spot.ratchetIndex+1)
				PlayRatchet(spot.ratchetIndex-1, ratchetInterval, onMessage, offMessage, sendFn)
			}
		}
	}
}

func PlayRatchet(number uint8, timeInterval time.Duration, onMessage, offMessage midi.Message, sendFn SendFunc) {
	fn := func() {
		err := sendFn(onMessage)
		if err != nil {
			panic("ratchet note on failed")
		}
		err = sendFn(offMessage)
		if err != nil {
			panic("ratchet note off failed")
		}
		if number > 0 {
			PlayRatchet(number-1, timeInterval, onMessage, offMessage, sendFn)
		}
	}
	time.AfterFunc(timeInterval, fn)
}

func (m *model) AddTrigger() {
	if m.overlayKey == ROOT_OVERLAY {
		m.lines[m.cursorPos.line].notes[m.cursorPos.beat] = note{5, 0, ACTION_NOTHING}
	} else {
		if len(m.overlays[m.overlayKey]) == 0 {
			m.overlays[m.overlayKey] = make(overlay)
		}
		m.overlays[m.overlayKey][m.cursorPos] = note{5, 0, ACTION_NOTHING}
	}
}

func (m *model) AddAction(act action) {

	if m.overlayKey == ROOT_OVERLAY {
		m.lines[m.cursorPos.line].notes[m.cursorPos.beat] = note{0, 0, act}
	} else {
		if len(m.overlays[m.overlayKey]) == 0 {
			m.overlays[m.overlayKey] = make(overlay)
		}
		m.overlays[m.overlayKey][m.cursorPos] = note{0, 0, act}
	}
}

func (m *model) RemoveTrigger() {
	if m.overlayKey == ROOT_OVERLAY {
		m.lines[m.cursorPos.line].notes[m.cursorPos.beat] = zeronote
	} else {
		if len(m.overlays[m.overlayKey]) == 0 {
			m.overlays[m.overlayKey] = make(overlay)
		}
		m.overlays[m.overlayKey][m.cursorPos] = zeronote
	}
}

func (m *model) IncreaseRatchet() {
	ratchetIndex := m.lines[m.cursorPos.line].notes[m.cursorPos.beat].ratchetIndex
	if ratchetIndex+1 < uint8(len(ratchets)) {
		m.lines[m.cursorPos.line].notes[m.cursorPos.beat].ratchetIndex++
	}
}

func (m *model) DecreaseRatchet() {
	ratchetIndex := m.lines[m.cursorPos.line].notes[m.cursorPos.beat].ratchetIndex

	if ratchetIndex > 0 {
		m.lines[m.cursorPos.line].notes[m.cursorPos.beat].ratchetIndex--
	}
}

func InitSeq(lineNumber int, beatNumber int) []line {
	var lines = make([]line, 0, lineNumber)

	for i := 0; i < lineNumber; i++ {
		lines = append(lines, InitLine(beatNumber))
	}
	return lines
}

func InitLine(beatNumber int) line {
	return line{
		make([]note, beatNumber),
	}
}

func InitPlayState(lines int) []linestate {
	linestates := make([]linestate, lines)
	for i, _ := range linestates {
		linestates[i].direction = 1
		linestates[i].resetDirection = 1
		linestates[i].resetLocation = 0
	}
	return linestates
}

func InitModel() model {
	outport, err := midi.OutPort(0)
	if err != nil {
		panic("Did not get midi outport")
	}
	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic("could not open log file")
	}

	newCursor := cursor.New()
	newCursor.Style = lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	return model{
		keys:           keys,
		beats:          32,
		tempo:          120,
		subdivisions:   2,
		help:           help.New(),
		lines:          InitSeq(8, 32),
		cursorPos:      gridKey{line: 0, beat: 0},
		cursor:         newCursor,
		outport:        outport,
		accentModifier: 1,
		logFile:        logFile,
		overlayKey:     ROOT_OVERLAY,
		overlays:       make(overlays),
	}
}

func (m model) LogTeaMsg(msg tea.Msg) {
	switch msg := msg.(type) {
	case beatMsg:
		m.LogString(fmt.Sprintf("beatMsg %d %d %d\n", msg.interval, m.totalBeats, m.tempo))
	case tea.KeyMsg:
		m.LogString(fmt.Sprintf("keyMsg %s\n", msg.String()))
	default:
		m.LogString(fmt.Sprintf("%T\n", msg))
	}
}

func (m model) LogString(message string) {
	_, err := m.logFile.WriteString(message)
	if err != nil {
		panic("could not write to log file")
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

	//m.LogTeaMsg(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case Is(msg, m.keys.Quit):
			m.logFile.Close()
			return m, tea.Quit
		case Is(msg, m.keys.CursorDown):
			if m.cursorPos.line < uint8(len(m.lines)-1) {
				m.cursorPos.line++
			}
		case Is(msg, m.keys.CursorUp):
			if m.cursorPos.line > 0 {
				m.cursorPos.line--
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
			zeroLine(m.lines[m.cursorPos.line])
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
				m.playState = InitPlayState(len(m.lines))
				m.trackTime = time.Duration(0)
				sendFn, err := midi.SendTo(m.outport)
				if err != nil {
					panic("sendFn is broken")
				}
				beatInterval := m.BeatInterval()
				return m, tea.Batch(PlayBeat(beatInterval, m.lines, m.playState, sendFn), BeatTick(beatInterval))
			}
		case Is(msg, m.keys.OverlayInputSwitch):
			m.overlaySelectionIndicator = (m.overlaySelectionIndicator + 1) % 3
			m.tempoSelectionIndicator = 0
		case Is(msg, m.keys.TempoInputSwitch):
			m.tempoSelectionIndicator = (m.tempoSelectionIndicator + 1) % 3
			m.overlaySelectionIndicator = 0
		case Is(msg, m.keys.TempoIncrease):
			switch {
			case m.tempoSelectionIndicator == 1:
				if m.tempo < 300 {
					m.tempo++
				}
			case m.tempoSelectionIndicator == 2:
				if m.subdivisions < 8 {
					m.subdivisions++
				}
			case m.overlaySelectionIndicator == 1:
				m.overlayKey.IncrementNumerator()
			case m.overlaySelectionIndicator == 2:
				m.overlayKey.IncrementDenominator()
			}
		case Is(msg, m.keys.TempoDecrease):
			switch {
			case m.tempoSelectionIndicator == 1:
				if m.tempo > 30 {
					m.tempo--
				}
			case m.tempoSelectionIndicator == 2:
				if m.subdivisions > 1 {
					m.subdivisions--
				}
			case m.overlaySelectionIndicator == 1:
				m.overlayKey.DecrementNumerator()
			case m.overlaySelectionIndicator == 2:
				m.overlayKey.DecrementDenominator()
			}
		case Is(msg, m.keys.ToggleAccentMode):
			m.accentMode = !m.accentMode
		case Is(msg, m.keys.ToggleAccentModifier):
			m.accentModifier = -1 * m.accentModifier
		case Is(msg, m.keys.RatchetIncrease):
			m.IncreaseRatchet()
		case Is(msg, m.keys.RatchetDecrease):
			m.DecreaseRatchet()
		case Is(msg, m.keys.ActionAddLineReset):
			m.AddAction(ACTION_LINE_RESET)
		case Is(msg, m.keys.ActionAddLineReverse):
			m.AddAction(ACTION_LINE_REVERSE)
		}
		if msg.String() >= "1" && msg.String() <= "9" {
			beatInterval, _ := strconv.Atoi(msg.String())
			line := m.lines[m.cursorPos.line]
			if m.accentMode {
				m.lines[m.cursorPos.line] = incrementAccent(line, m.cursorPos.beat, beatInterval, m.accentModifier)
			} else {
				m.lines[m.cursorPos.line] = fill(line, m.cursorPos.beat, beatInterval)
			}
		}
	case beatMsg:
		if m.playing {
			m.advanceCurrentBeat()
			m.totalBeats++
			sendFn, err := midi.SendTo(m.outport)
			if err != nil {
				panic("sendFn is broken")
			}
			beatInterval := m.BeatInterval()
			return m, tea.Batch(PlayBeat(beatInterval, m.lines, m.playState, sendFn), BeatTick(beatInterval))
		}
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor
	return m, cmd
}

func (m *model) advanceCurrentBeat() {
	for i, currentState := range m.playState {
		advancedBeat := int8(currentState.currentBeat) + currentState.direction
		if advancedBeat < 0 || advancedBeat >= int8(m.beats) {
			m.playState[i].currentBeat = currentState.resetLocation
			advancedBeat = int8(currentState.resetLocation)
			m.playState[i].direction = currentState.resetDirection
		} else {
			m.playState[i].currentBeat = uint8(advancedBeat)
		}

		switch m.lines[i].notes[advancedBeat].action {
		case ACTION_LINE_RESET:
			m.playState[i].currentBeat = 0
		case ACTION_LINE_REVERSE:
			m.playState[i].currentBeat = uint8(max(advancedBeat-1, 0))
			m.playState[i].direction = -1
		}
	}
}

func zeroLine(beatline line) {
	for i := range beatline.notes {
		beatline.notes[i] = zeronote
	}
}

func fill(beatline line, start uint8, every int) line {
	for i := range beatline.notes[start:] {
		if i%every == 0 {
			gridLocation := start + uint8(i)
			if beatline.notes[gridLocation] != zeronote {
				beatline.notes[gridLocation] = zeronote
			} else {
				beatline.notes[gridLocation] = note{5, 0, 0}
			}
		}
	}
	return beatline
}

func incrementAccent(beatline line, start uint8, every int, modifier int8) line {
	for i := range beatline.notes[start:] {
		if i%every == 0 {
			gridLocation := start + uint8(i)
			if beatline.notes[gridLocation] != zeronote {
				currentAccentIndex := beatline.notes[gridLocation].accentIndex
				nextAccentIndex := uint8(currentAccentIndex - uint8(1*modifier))
				if nextAccentIndex >= 1 && nextAccentIndex < uint8(len(accents)) {
					beatline.notes[gridLocation].accentIndex = nextAccentIndex
				}
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
	buf.WriteString(lipgloss.JoinHorizontal(0, m.TempoView(), "  ", m.ViewTriggerSeq(), " ", AccentKeyView()))
	return buf.String()
}

func AccentKeyView() string {
	var buf strings.Builder
	buf.WriteString("      ACCENTS\n")
	buf.WriteString("  ----------------\n")
	for _, accent := range accents[1:] {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(accent.color))
		buf.WriteString(fmt.Sprintf("  %s  \n", style.Render(string(accent.shape))))
	}
	return buf.String()
}

var accentModeStyle = lipgloss.NewStyle().Background(accents[1].color).Foreground(lipgloss.Color("#000000"))

func (m model) ViewTriggerSeq() string {
	var buf strings.Builder
	var overlay = m.ViewOverlay()
	var mode string
	if m.accentMode {
		if m.accentModifier > 0 {
			mode = " Accent Mode \u2191 "
		} else {
			mode = " Accent Mode \u2193 "
		}
		buf.WriteString(fmt.Sprintf("   Seq - %s - %s\n", accentModeStyle.Render(mode), overlay))
	} else if m.overlaySelectionIndicator > 0 {
		buf.WriteString(fmt.Sprintf("   Seq - %s - %s\n", "Overlay", overlay))
	} else {
		buf.WriteString("   Seq - A sequencer for your cli\n")
	}
	buf.WriteString("  ┌─────────────────────────────────\n")
	for i, line := range m.lines {
		buf.WriteString(line.View(i, m))
	}
	if m.playing {
		buf.WriteString(fmt.Sprintf("   %*s%s", m.playState[0].currentBeat, "", "█"))
	}
	buf.WriteString("\n")
	// buf.WriteString(m.help.View(m.keys))
	// buf.WriteString("\n")
	return buf.String()
}

func (m model) ViewOverlay() string {
	var numerator, denominator string
	switch m.overlaySelectionIndicator {
	case 0:
		numerator = numberColor.Render(strconv.Itoa(int(m.overlayKey.num)))
		denominator = numberColor.Render(strconv.Itoa(int(m.overlayKey.denom)))
	case 1:
		numerator = selectedColor.Render(strconv.Itoa(int(m.overlayKey.num)))
		denominator = numberColor.Render(strconv.Itoa(int(m.overlayKey.denom)))
	case 2:
		numerator = numberColor.Render(strconv.Itoa(int(m.overlayKey.num)))
		denominator = selectedColor.Render(strconv.Itoa(int(m.overlayKey.denom)))
	}
	return fmt.Sprintf("%s/%s", numerator, denominator)
}

var altSeqColor = lipgloss.Color("#222222")
var seqColor = lipgloss.Color("#000000")
var seqCursorColor = lipgloss.Color("#444444")
var seqOverlayColor = lipgloss.Color("#dddddd")

func (line line) View(lineNumber int, m model) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d │", lineNumber))

	var backgroundSeqColor lipgloss.Color

	for i := uint8(0); i < uint8(m.beats); i++ {
		overlayNote, hasOverlayNote := m.overlays[m.overlayKey][gridKey{uint8(lineNumber), i}]
		if hasOverlayNote {
			backgroundSeqColor = seqOverlayColor
		} else if m.playing && m.playState[lineNumber].currentBeat == i {
			backgroundSeqColor = seqCursorColor
		} else if i%8 > 3 {
			backgroundSeqColor = altSeqColor
		} else {
			backgroundSeqColor = seqColor
		}

		var char string
		var foregroundColor lipgloss.Color
		var currentNote note
		if hasOverlayNote {
			currentNote = overlayNote
		} else {
			currentNote = line.notes[i]
		}
		currentAccent := accents[currentNote.accentIndex]
		currentAction := currentNote.action

		if currentAction == ACTION_NOTHING {
			char = string(currentAccent.shape) + string(ratchets[currentNote.ratchetIndex])
			foregroundColor = currentAccent.color
		} else {
			lineaction := lineactions[currentAction]
			char = string(lineaction.shape)
			foregroundColor = lineaction.color
		}

		if m.cursorPos.line == uint8(lineNumber) && m.cursorPos.beat == i {
			m.cursor.SetChar(char)
			char = m.cursor.View()
			buf.WriteString(lipgloss.NewStyle().Background(backgroundSeqColor).Render(char))
		} else {
			style := lipgloss.NewStyle().Background(backgroundSeqColor).Foreground(foregroundColor)
			buf.WriteString(style.Render(char))
		}
	}

	buf.WriteString("\n")
	return buf.String()
}
