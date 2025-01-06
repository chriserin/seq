package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
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
	SelectKeyLine        key.Binding
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
	ClearSeq:             Key("C", "Clear Overlay"),
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
	SelectKeyLine:        Key("K", "Select Key Line"),
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

func (n note) IncrementAccent(modifier int8) note {
	var newAccent = int8(n.accentIndex) - modifier
	if newAccent >= 1 && newAccent < int8(len(accents)) {
		n.accentIndex = uint8(newAccent)
	}
	return n
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

type line []note

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

func OverlayKeySort(x overlayKey, y overlayKey) int {
	var result int
	if x.num == y.num {
		result = int(y.denom) - int(x.denom)
	} else {
		result = int(y.num) - int(x.num)
	}
	return result
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

type Notable interface {
	SetNote(gridKey, note)
}
type pattern []line

func (p *pattern) SetNote(gridKey gridKey, note note) {
	(*p)[gridKey.line][gridKey.beat] = note
}

type overlays map[overlayKey]overlay

func (ol *overlay) SetNote(gridKey gridKey, note note) {
	(*ol)[gridKey] = note
}

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
	rootPattern               pattern
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
	keyline                   uint8
	keyCycles                 int
	playingMatchedOverlay     overlayKey
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

func PlayBeat(beatInterval time.Duration, lines pattern, currentBeat []linestate, sendFn SendFunc) tea.Cmd {
	return func() tea.Msg {
		Play(beatInterval, lines, currentBeat, sendFn)
		return nil
	}
}

type SendFunc func(msg midi.Message) error

func Play(beatInterval time.Duration, lines pattern, currentBeat []linestate, sendFn SendFunc) {
	for i, line := range lines {
		spot := line[currentBeat[i].currentBeat]
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

func (m *model) EnsureOverlay() {
	if m.overlayKey != ROOT_OVERLAY && len(m.overlays[m.overlayKey]) == 0 {
		m.overlays[m.overlayKey] = make(overlay)
	}
}

func (m *model) CurrentNotable() Notable {
	var notable Notable
	if m.overlayKey == ROOT_OVERLAY {
		notable = &m.rootPattern
	} else {
		overlay := m.overlays[m.overlayKey]
		notable = &overlay
	}
	return notable
}

func (m *model) AddTrigger() {
	m.CurrentNotable().SetNote(m.cursorPos, note{5, 0, ACTION_NOTHING})
}

func (m *model) AddAction(act action) {
	m.CurrentNotable().SetNote(m.cursorPos, note{0, 0, act})
}

func (m *model) RemoveTrigger() {
	m.CurrentNotable().SetNote(m.cursorPos, zeronote)
}

func (m *model) IncreaseRatchet() {
	line := m.CombinedLine(m.cursorPos.line, m.overlayKey)
	currentNote := line[m.cursorPos.beat]
	currentRatchet := currentNote.ratchetIndex

	if currentRatchet+1 < uint8(len(ratchets)) {
		currentNote.ratchetIndex = currentNote.ratchetIndex + 1
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) DecreaseRatchet() {
	line := m.CombinedLine(m.cursorPos.line, m.overlayKey)
	currentNote := line[m.cursorPos.beat]
	currentRatchet := currentNote.ratchetIndex

	if currentRatchet > 0 {
		currentNote.ratchetIndex = currentNote.ratchetIndex - 1
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func InitSeq(lineNumber int, beatNumber int) pattern {
	var lines = make([]line, 0, lineNumber)

	for i := 0; i < lineNumber; i++ {
		lines = append(lines, InitLine(beatNumber))
	}
	return lines
}

func InitLine(beatNumber int) line {
	return make([]note, beatNumber)
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
		rootPattern:    InitSeq(8, 32),
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

func GetMatchingOverlays(keyCycles int, keys []overlayKey) []overlayKey {
	var matchedKeys = make([]overlayKey, 0, 5)

	for _, k := range keys {
		additional := k.num / k.denom
		over := k.num % k.denom
		if over > 0 {
			rem := keyCycles % (int(k.denom) * (1 + int(additional)))
			if rem == int(k.num) {
				matchedKeys = append(matchedKeys, k)
			}
		} else {
			if keyCycles%int(k.num) == 0 {
				matchedKeys = append(matchedKeys, k)
			}
		}
	}

	return matchedKeys
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
			if m.cursorPos.line < uint8(len(m.rootPattern)-1) {
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
			m.EnsureOverlay()
			m.AddTrigger()
		case Is(msg, m.keys.TriggerRemove):
			m.EnsureOverlay()
			m.RemoveTrigger()
		case Is(msg, m.keys.ClearLine):
			m.EnsureOverlay()
			m.ClearOverlayLine()
		case Is(msg, m.keys.ClearSeq):
			m.ClearOverlay()
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
				m.keyCycles = 0
				m.totalBeats = 0
				m.playState = InitPlayState(len(m.rootPattern))
				m.advanceKeyCycle()
				m.trackTime = time.Duration(0)
				sendFn, err := midi.SendTo(m.outport)
				if err != nil {
					panic("sendFn is broken")
				}
				beatInterval := m.BeatInterval()
				return m, tea.Batch(PlayBeat(beatInterval, m.CombinedPattern(), m.playState, sendFn), BeatTick(beatInterval))
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
			m.EnsureOverlay()
			m.IncreaseRatchet()
		case Is(msg, m.keys.RatchetDecrease):
			m.EnsureOverlay()
			m.DecreaseRatchet()
		case Is(msg, m.keys.ActionAddLineReset):
			m.EnsureOverlay()
			m.AddAction(ACTION_LINE_RESET)
		case Is(msg, m.keys.ActionAddLineReverse):
			m.EnsureOverlay()
			m.AddAction(ACTION_LINE_REVERSE)
		case Is(msg, m.keys.SelectKeyLine):
			m.keyline = m.cursorPos.line
		}
		if msg.String() >= "1" && msg.String() <= "9" {
			m.EnsureOverlay()
			beatInterval, _ := strconv.Atoi(msg.String())
			if m.accentMode {
				m.incrementAccent(beatInterval)
			} else {
				m.fill(beatInterval)
			}
		}
	case beatMsg:
		if m.playing {
			m.advanceCurrentBeat()
			m.advanceKeyCycle()
			m.totalBeats++
			sendFn, err := midi.SendTo(m.outport)
			if err != nil {
				panic("sendFn is broken")
			}
			beatInterval := m.BeatInterval()
			return m, tea.Batch(PlayBeat(beatInterval, m.CombinedPattern(), m.playState, sendFn), BeatTick(beatInterval))
		}
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor
	return m, cmd
}

func (m *model) ClearOverlayLine() {
	if m.overlayKey == ROOT_OVERLAY {
		zeroLine(m.rootPattern[m.cursorPos.line])
	} else {
		for i := uint8(0); i < m.beats; i++ {
			key := gridKey{m.cursorPos.line, i}
			delete(m.overlays[m.overlayKey], key)
		}
	}
}

func (m *model) ClearOverlay() {
	if m.overlayKey == ROOT_OVERLAY {
		m.rootPattern = InitSeq(8, 32)
	} else {

		delete(m.overlays, m.overlayKey)

	}
}

func (m *model) advanceCurrentBeat() {
	combinedPattern := m.CombinedPattern()
	for i, currentState := range m.playState {
		advancedBeat := int8(currentState.currentBeat) + currentState.direction
		if advancedBeat < 0 || advancedBeat >= int8(m.beats) {
			m.playState[i].currentBeat = currentState.resetLocation
			advancedBeat = int8(currentState.resetLocation)
			m.playState[i].direction = currentState.resetDirection
		} else {
			m.playState[i].currentBeat = uint8(advancedBeat)
		}

		switch combinedPattern[i][advancedBeat].action {
		case ACTION_LINE_RESET:
			m.playState[i].currentBeat = 0
		case ACTION_LINE_REVERSE:
			m.playState[i].currentBeat = uint8(max(advancedBeat-1, 0))
			m.playState[i].direction = -1
		}
	}
}

func (m *model) advanceKeyCycle() {
	if m.playState[m.keyline].currentBeat == 0 {
		m.keyCycles++
		m.playingMatchedOverlay = m.MatchingOverlay()
	}
}

func zeroLine(beatline line) {
	for i := range beatline {
		beatline[i] = zeronote
	}
}

func (m model) CombinedPattern() pattern {
	newPattern := make(pattern, 0, len(m.rootPattern))
	for i, _ := range m.rootPattern {
		newPattern = append(newPattern, m.CombinedLine(uint8(i), m.playingMatchedOverlay))
	}
	return newPattern
}

func (m model) CombinedLine(lineNo uint8, key overlayKey) line {
	var rootLine line = m.rootPattern[lineNo]
	if key == ROOT_OVERLAY || key == (overlayKey{0, 0}) {
		return rootLine
	}
	var resultLine line = make([]note, len(rootLine))
	copy(resultLine, rootLine)
	for i := range rootLine {
		gk := gridKey{lineNo, uint8(i)}
		note, hasNote := m.overlays[key][gk]
		if hasNote {
			resultLine[i] = note
		}
	}
	return resultLine
}

func (m *model) fill(every int) {
	combinedLine := m.CombinedLine(m.cursorPos.line, m.overlayKey)
	start := m.cursorPos.beat

	for i := range combinedLine[start:] {
		if i%every == 0 {
			currentBeat := start + uint8(i)
			gridKey := gridKey{m.cursorPos.line, currentBeat}
			if combinedLine[currentBeat] != zeronote {
				m.CurrentNotable().SetNote(gridKey, zeronote)
			} else {
				m.CurrentNotable().SetNote(gridKey, note{5, 0, 0})
			}
		}
	}
}

func (m *model) incrementAccent(every int) {
	combinedLine := m.CombinedLine(m.cursorPos.line, m.overlayKey)
	start := m.cursorPos.beat

	for i := range combinedLine[start:] {
		if i%every == 0 {
			currentBeat := start + uint8(i)
			gridKey := gridKey{m.cursorPos.line, currentBeat}
			if combinedLine[currentBeat] != zeronote {
				n := combinedLine[currentBeat]
				m.CurrentNotable().SetNote(gridKey, n.IncrementAccent(m.accentModifier))
			}
		}
	}
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
	} else if m.playing {
		buf.WriteString(fmt.Sprintf("   Seq - Playing - %d - %s\n", m.keyCycles, m.CurrentOverlayView()))
	} else {
		buf.WriteString("   Seq - A sequencer for your cli\n")
	}
	buf.WriteString("  ┌─────────────────────────────────\n")
	for i, line := range m.rootPattern {
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

func (m model) CurrentOverlayView() string {
	matchedKey := m.playingMatchedOverlay
	if matchedKey != (overlayKey{0, 0}) {
		return fmt.Sprintf("%d/%d", matchedKey.num, matchedKey.denom)
	}
	return " "
}

func (m model) MatchingOverlay() overlayKey {
	keys := make([]overlayKey, 0, 5)
	for k := range maps.Keys(m.overlays) {
		keys = append(keys, k)
	}
	matchingKeys := GetMatchingOverlays(m.keyCycles, keys)
	if len(matchingKeys) > 0 {
		slices.SortFunc(matchingKeys, OverlayKeySort)
		return matchingKeys[0]
	} else {
		return overlayKey{0, 0}
	}
}

var altSeqColor = lipgloss.Color("#222222")
var seqColor = lipgloss.Color("#000000")
var seqCursorColor = lipgloss.Color("#444444")
var seqOverlayColor = lipgloss.Color("#333388")

func KeyLineIndicator(k uint8, l uint8) string {
	if k == l {
		return "K"
	} else {
		return " "
	}
}

func (line line) View(lineNumber int, m model) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d%s│", lineNumber, KeyLineIndicator(m.keyline, uint8(lineNumber))))

	var backgroundSeqColor lipgloss.Color

	var key overlayKey
	if m.playing {
		key = m.playingMatchedOverlay
	} else {
		key = m.overlayKey
	}

	for i := uint8(0); i < uint8(m.beats); i++ {
		overlayNote, hasOverlayNote := m.overlays[key][gridKey{uint8(lineNumber), i}]
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
			currentNote = line[i]
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
