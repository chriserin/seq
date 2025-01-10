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
	OverlayTriggerRemove key.Binding
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
	NextOverlay          key.Binding
	PrevOverlay          key.Binding
	StackUpOverlay       key.Binding
	PressDownOverlay     key.Binding
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
	OverlayTriggerRemove: Key("x", "Remove Overlay Note"),
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
	NextOverlay:          Key("{", "Next Overlay"),
	PrevOverlay:          Key("}", "Prev Overlay"),
	StackUpOverlay:       Key("ctrl+s", "StackUp Overlay"),
	PressDownOverlay:     Key("ctrl+p", "PressDown Overlay"),
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

func (o overlayKey) String() string {
	return fmt.Sprintf("%d/%d", o.num, o.denom)
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
	lines                     uint8
	beats                     uint8
	tempo                     int
	subdivisions              int
	help                      help.Model
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
	playingMatchedOverlays    []overlayKey
	stackedupKeys             []overlayKey
	pressedDownKeys           []overlayKey
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

func PlayBeat(beatInterval time.Duration, lines uint8, pattern overlay, currentBeat []linestate, sendFn SendFunc) tea.Cmd {
	return func() tea.Msg {
		Play(beatInterval, lines, pattern, currentBeat, sendFn)
		return nil
	}
}

type SendFunc func(msg midi.Message) error

func Play(beatInterval time.Duration, lines uint8, pattern overlay, currentBeat []linestate, sendFn SendFunc) {
	for i := range lines {
		currentGridKey := gridKey{i, currentBeat[i].currentBeat}
		note, hasNote := pattern[currentGridKey]
		if hasNote && note != zeronote {
			onMessage := midi.NoteOn(10, C1+uint8(i), accents[note.accentIndex].value)
			offMessage := midi.NoteOff(10, C1+uint8(i))
			err := sendFn(onMessage)
			if err != nil {
				panic("note on failed")
			}
			err = sendFn(offMessage)
			if err != nil {
				panic("note off failed")
			}
			if note.ratchetIndex > 0 {
				ratchetInterval := beatInterval / time.Duration(note.ratchetIndex+1)
				PlayRatchet(note.ratchetIndex-1, ratchetInterval, onMessage, offMessage, sendFn)
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
	if len(m.overlays[m.overlayKey]) == 0 {
		m.overlays[m.overlayKey] = make(overlay)
		if m.playing {
			m.determineMatachedOverlays()
		}
	}
}

func (m *model) CurrentNotable() Notable {
	var notable Notable
	overlay := m.overlays[m.overlayKey]
	notable = &overlay
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

func (m *model) OverlayRemoveTrigger() {
	delete(m.overlays[m.overlayKey], m.cursorPos)
}

func (m *model) IncreaseRatchet() {
	rootOverlay := m.overlays[ROOT_OVERLAY]
	currentOverlay := m.overlays[m.overlayKey]
	rootNote, rootHasNote := rootOverlay[m.cursorPos]
	currentOverlayNote, currentHasNote := currentOverlay[m.cursorPos]
	var currentRatchet uint8
	var currentNote note
	if currentHasNote {
		currentRatchet = currentOverlayNote.ratchetIndex
		currentNote = currentOverlayNote
	} else if rootHasNote {
		currentRatchet = rootNote.ratchetIndex
		currentNote = rootNote
	}

	if currentNote.action == ACTION_NOTHING && currentRatchet+1 < uint8(len(ratchets)) {
		currentNote.ratchetIndex = currentNote.ratchetIndex + 1
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) DecreaseRatchet() {
	rootOverlay := m.overlays[ROOT_OVERLAY]
	currentOverlay := m.overlays[m.overlayKey]
	rootNote, rootHasNote := rootOverlay[m.cursorPos]
	currentOverlayNote, currentHasNote := currentOverlay[m.cursorPos]
	var currentRatchet uint8
	var currentNote note
	if currentHasNote {
		currentRatchet = currentOverlayNote.ratchetIndex
		currentNote = currentOverlayNote
	} else if rootHasNote {
		currentRatchet = rootNote.ratchetIndex
		currentNote = rootNote
	}

	if currentNote.action == ACTION_NOTHING && currentRatchet > 0 {
		currentNote.ratchetIndex = currentNote.ratchetIndex - 1
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func InitPlayState(lines uint8) []linestate {
	linestates := make([]linestate, lines)
	for i, _ := range linestates {
		linestates[i].direction = 1
		linestates[i].resetDirection = 1
		linestates[i].resetLocation = 0
	}
	return linestates
}

func InitModel(midiOutport drivers.Out) model {
	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic("could not open log file")
	}

	newCursor := cursor.New()
	newCursor.Style = lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	return model{
		keys:            keys,
		lines:           8,
		beats:           32,
		tempo:           120,
		subdivisions:    2,
		help:            help.New(),
		cursorPos:       gridKey{line: 0, beat: 0},
		cursor:          newCursor,
		outport:         midiOutport,
		accentModifier:  1,
		logFile:         logFile,
		overlayKey:      ROOT_OVERLAY,
		overlays:        make(overlays),
		stackedupKeys:   []overlayKey{ROOT_OVERLAY},
		pressedDownKeys: []overlayKey{},
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
	_, err := m.logFile.WriteString(message + "\n")
	if err != nil {
		panic("could not write to log file")
	}
}

func RunProgram(midiOutport drivers.Out) *tea.Program {
	p := tea.NewProgram(InitModel(midiOutport), tea.WithAltScreen())
	return p
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return tea.FocusMsg{} }
}

func Is(msg tea.KeyMsg, k key.Binding) bool {
	return key.Matches(msg, k)
}

func (m model) GetMatchingOverlays(keyCycles int, keys []overlayKey) []overlayKey {
	var matchedKeys = make([]overlayKey, 0, 5)

	slices.SortFunc(keys, OverlayKeySort)
	var pressNext = false

	for _, key := range keys {
		matches := DoesKeyMatch(keyCycles, key)
		if (matches && len(matchedKeys) == 0) || pressNext {
			matchedKeys = append(matchedKeys, key)
			if slices.Index(m.pressedDownKeys, key) >= 0 {
				pressNext = true
			} else {
				pressNext = false
			}
		} else if matches && len(matchedKeys) != 0 {
			if slices.Index(m.stackedupKeys, key) >= 0 {
				matchedKeys = append(matchedKeys, key)
			}
		}
	}

	return matchedKeys
}

func DoesKeyMatch(keyCycles int, key overlayKey) bool {
	additional := key.num / key.denom
	over := key.num % key.denom
	if over > 0 {
		rem := keyCycles % (int(key.denom) * (1 + int(additional)))
		if rem == int(key.num) {
			return true
		}
	} else {
		if keyCycles%int(key.num) == 0 {
			return true
		}
	}
	return false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case Is(msg, m.keys.Quit):
			m.logFile.Close()
			return m, tea.Quit
		case Is(msg, m.keys.CursorDown):
			if m.cursorPos.line < m.lines-1 {
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
		case Is(msg, m.keys.OverlayTriggerRemove):
			m.EnsureOverlay()
			m.OverlayRemoveTrigger()
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

			// if m.playing && m.outport.IsOpen() {
			// 	m.outport.Close()
			// }

			m.playing = !m.playing
			m.playTime = time.Now()
			if m.playing {
				m.keyCycles = 0
				m.totalBeats = 0
				m.playState = InitPlayState(m.lines)
				m.advanceKeyCycle()
				m.trackTime = time.Duration(0)
				sendFn, err := midi.SendTo(m.outport)
				if err != nil {
					panic("sendFn is broken")
				}
				beatInterval := m.BeatInterval()
				return m, tea.Batch(PlayBeat(beatInterval, m.lines, m.CombinedPattern(m.playingMatchedOverlays), m.playState, sendFn), BeatTick(beatInterval))
			} else {
				m.keyCycles = 0
				m.playingMatchedOverlays = []overlayKey{}
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
		case Is(msg, m.keys.PrevOverlay):
			m.PrevOverlay()
		case Is(msg, m.keys.NextOverlay):
			m.NextOverlay()
		case Is(msg, m.keys.StackUpOverlay):
			m.ToggleStackupOverlay(m.overlayKey)
		case Is(msg, m.keys.PressDownOverlay):
			m.TogglePressdownOverlay(m.overlayKey)
		}
		if msg.String() >= "1" && msg.String() <= "9" {
			m.EnsureOverlay()
			beatInterval, _ := strconv.ParseInt(msg.String(), 0, 8)
			if m.accentMode {
				m.incrementAccent(uint8(beatInterval))
			} else {
				m.fill(uint8(beatInterval))
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
			return m, tea.Batch(PlayBeat(beatInterval, m.lines, m.CombinedPattern(m.playingMatchedOverlays), m.playState, sendFn), BeatTick(beatInterval))
		}
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor
	return m, cmd
}

func (m *model) ToggleStackupOverlay(key overlayKey) {
	index := slices.Index(m.stackedupKeys, m.overlayKey)
	if index < 0 {
		m.stackedupKeys = append(m.stackedupKeys, m.overlayKey)
	} else {
		m.stackedupKeys = append(m.stackedupKeys[:index], m.stackedupKeys[index+1:]...)
	}
}

func (m *model) TogglePressdownOverlay(key overlayKey) {
	index := slices.Index(m.pressedDownKeys, m.overlayKey)
	if index < 0 {
		m.pressedDownKeys = append(m.pressedDownKeys, m.overlayKey)
	} else {
		m.pressedDownKeys = append(m.pressedDownKeys[:index], m.pressedDownKeys[index+1:]...)
	}
}

func RemoveRootKey(keys []overlayKey) []overlayKey {
	index := slices.Index(keys, ROOT_OVERLAY)
	if index >= 0 {
		return append(keys[:index], keys[index+1:]...)
	}
	return keys
}

func (m *model) PrevOverlay() {
	keys := m.OverlayKeys()
	slices.SortFunc(keys, OverlayKeySort)
	index := slices.Index(keys, m.overlayKey)
	if index+1 < len(keys) {
		m.overlayKey = keys[index+1]
	}
}

func (m *model) NextOverlay() {
	keys := m.OverlayKeys()
	slices.SortFunc(keys, OverlayKeySort)
	index := slices.Index(keys, m.overlayKey)
	if index-1 >= 0 {
		m.overlayKey = keys[index-1]
	}
}

func (m *model) ClearOverlayLine() {
	for i := uint8(0); i < m.beats; i++ {
		key := gridKey{m.cursorPos.line, i}
		delete(m.overlays[m.overlayKey], key)
	}
}

func (m *model) ClearOverlay() {
	delete(m.overlays, m.overlayKey)
}

func (m *model) advanceCurrentBeat() {
	combinedPattern := m.CombinedPattern(m.playingMatchedOverlays)
	for i, currentState := range m.playState {
		advancedBeat := int8(currentState.currentBeat) + currentState.direction

		if advancedBeat < 0 || advancedBeat >= int8(m.beats) {
			m.playState[i].currentBeat = currentState.resetLocation
			advancedBeat = int8(currentState.resetLocation)
			m.playState[i].direction = currentState.resetDirection
		} else {
			m.playState[i].currentBeat = uint8(advancedBeat)
		}

		switch combinedPattern[gridKey{uint8(i), uint8(advancedBeat)}].action {
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
		m.determineMatachedOverlays()
	}
}

func (m *model) determineMatachedOverlays() {
	keys := m.OverlayKeys()
	m.playingMatchedOverlays = m.GetMatchingOverlays(m.keyCycles, keys)
}

func (m model) CombinedPattern(keys []overlayKey) overlay {
	var combinedOverlay = make(overlay)

	for _, key := range slices.Backward(keys) {
		for gridKey, note := range m.overlays[key] {
			combinedOverlay[gridKey] = note
		}
	}
	return combinedOverlay
}

func (m *model) fill(every uint8) {
	start := m.cursorPos.beat

	matchedKeys := m.GetMatchingOverlays(0, m.OverlayKeys())
	matchedKeys = append(matchedKeys, m.overlayKey)
	combinedOverlay := m.CombinedPattern(matchedKeys)

	for i := uint8(0); i < m.beats; i++ {
		if i%every == 0 {
			currentBeat := start + i
			if currentBeat >= m.beats {
				return
			}
			gridKey := gridKey{m.cursorPos.line, currentBeat}
			currentNote, hasNote := combinedOverlay[gridKey]
			hasNote = hasNote && currentNote != zeronote

			if m.overlayKey != ROOT_OVERLAY && hasNote {
				m.CurrentNotable().SetNote(gridKey, zeronote)
			} else if m.overlayKey == ROOT_OVERLAY && hasNote {
				delete(m.overlays[ROOT_OVERLAY], gridKey)
			} else {
				m.CurrentNotable().SetNote(gridKey, note{5, 0, 0})
			}
		}
	}
}

func (m *model) incrementAccent(every uint8) {
	start := m.cursorPos.beat
	rootOverlay := m.overlays[ROOT_OVERLAY]
	currentOverlay := m.overlays[m.overlayKey]

	for i := uint8(0); i < m.beats; i++ {
		if i%every == 0 {
			currentBeat := start + i
			if currentBeat >= m.beats {
				return
			}
			gridKey := gridKey{m.cursorPos.line, currentBeat}
			rootNote, rootHasNote := rootOverlay[gridKey]
			currentNote, currentHasNote := currentOverlay[gridKey]
			if currentHasNote && currentNote != zeronote {
				m.CurrentNotable().SetNote(gridKey, currentNote.IncrementAccent(m.accentModifier))
			} else if rootHasNote && rootNote != zeronote {
				m.CurrentNotable().SetNote(gridKey, rootNote.IncrementAccent(m.accentModifier))
			}
		}
	}
}

func (m model) OverlayKeys() []overlayKey {
	keys := make([]overlayKey, 0, 5)
	return slices.AppendSeq(keys, maps.Keys(m.overlays))
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
	var sideView string

	if m.accentMode {
		sideView = AccentKeyView()
	} else if len(m.overlays) > 0 {
		sideView = m.OverlaysView()
	} else {
		sideView = m.SetupView()
	}

	buf.WriteString(lipgloss.JoinHorizontal(0, m.TempoView(), "  ", m.ViewTriggerSeq(), "  ", sideView))
	return buf.String()
}

func AccentKeyView() string {
	var buf strings.Builder
	buf.WriteString("    ACCENTS\n")
	buf.WriteString("———————————————\n")
	for _, accent := range accents[1:] {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(accent.color))
		buf.WriteString(fmt.Sprintf("  %s  -  %d\n", style.Render(string(accent.shape)), accent.value))
	}
	return buf.String()
}

func (m model) SetupView() string {
	var buf strings.Builder
	buf.WriteString("    Setup\n")
	buf.WriteString("———————————————\n")
	for range m.lines {
		buf.WriteString("CH 10 Note C1\n")
	}
	return buf.String()
}

const SELECTED_OVERLAY_ARROW = "\u2192"

var currentPlayingColor lipgloss.Color = "#abfaa9"
var activePlayingColor lipgloss.Color = "#f34213"

func (m model) OverlaysView() string {
	var buf strings.Builder
	buf.WriteString("Overlays\n")
	buf.WriteString("——————————————\n")
	keys := m.OverlayKeys()
	slices.SortFunc(keys, OverlayKeySort)
	style := lipgloss.NewStyle().Background(seqOverlayColor)
	for _, k := range keys {
		var playingSpacer = "   "
		var playing = ""
		if m.playing && m.playingMatchedOverlays[0] == k {
			playing = lipgloss.NewStyle().Background(seqOverlayColor).Foreground(activePlayingColor).Render(" \u25CF ")
			buf.WriteString(playing)
			playingSpacer = ""
		} else if m.playing && slices.Contains(m.playingMatchedOverlays, k) {
			playing = lipgloss.NewStyle().Background(seqOverlayColor).Foreground(currentPlayingColor).Render(" \u25C9 ")
			buf.WriteString(playing)
			playingSpacer = ""
		}
		var editing = ""
		if m.overlayKey == k {
			editing = " E"
		}
		var stackModifier = ""
		if slices.Index(m.pressedDownKeys, k) >= 0 {
			stackModifier = " \u2193\u0332"
		} else if slices.Index(m.stackedupKeys, k) >= 0 {
			stackModifier = " \u2191\u0305"
		}
		overlayLine := fmt.Sprintf("%d/%d%2s%2s", k.num, k.denom, stackModifier, editing)

		buf.WriteString(playingSpacer)
		if slices.Contains(m.playingMatchedOverlays, k) {
			buf.WriteString(style.Render(overlayLine))
		} else {
			buf.WriteString(overlayLine)
		}
		buf.WriteString(playing)
		buf.WriteString(playingSpacer)
		buf.WriteString("\n")
	}
	return buf.String()
}

var accentModeStyle = lipgloss.NewStyle().Background(accents[1].color).Foreground(lipgloss.Color("#000000"))

func (m model) ViewTriggerSeq() string {
	var buf strings.Builder
	var mode string
	if m.accentMode {
		if m.accentModifier > 0 {
			mode = " Accent Mode \u2191 "
		} else {
			mode = " Accent Mode \u2193 "
		}
		buf.WriteString(fmt.Sprintf("   Seq - %s\n", accentModeStyle.Render(mode)))
	} else if m.playing {
		buf.WriteString(fmt.Sprintf("   Seq - Playing - %d\n", m.keyCycles))
	} else {
		buf.WriteString("   Seq - A sequencer for your cli\n")
	}
	buf.WriteString("  ┌─────────────────────────────────\n")
	for i := uint8(0); i < m.lines; i++ {
		buf.WriteString(lineView(i, m))
	}
	buf.WriteString(m.CurrentOverlayView())
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
	var matchedKey overlayKey
	if len(m.playingMatchedOverlays) > 0 {
		matchedKey = m.playingMatchedOverlays[0]
	} else {
		matchedKey = overlayKey{1, 1}
	}
	return fmt.Sprintf("   Editing - %s     Playing - %d/%d", m.ViewOverlay(), matchedKey.num, matchedKey.denom)
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

func lineView(lineNumber uint8, m model) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d%s│", lineNumber, KeyLineIndicator(m.keyline, lineNumber)))

	var backgroundSeqColor lipgloss.Color

	var currentKeys []overlayKey
	if m.playing {
		currentKeys = m.playingMatchedOverlays
	} else {
		currentKeys = []overlayKey{m.overlayKey}
	}

	hasRoot := slices.Index(m.stackedupKeys, overlayKey{1, 1}) >= 0 || slices.Index(currentKeys, overlayKey{1, 1}) >= 0
	combinedOverlay := m.CombinedPattern(RemoveRootKey(currentKeys))

	for i := uint8(0); i < m.beats; i++ {
		currentGridKey := gridKey{uint8(lineNumber), i}
		overlayNote, hasOverlayNote := combinedOverlay[currentGridKey]
		rootNote := m.overlays[ROOT_OVERLAY][currentGridKey]

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
		if hasOverlayNote || !hasRoot {
			currentNote = overlayNote
		} else {
			currentNote = rootNote
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
