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

type transitiveKeyMap struct {
	Quit                 key.Binding
	Help                 key.Binding
	CursorUp             key.Binding
	CursorDown           key.Binding
	CursorLeft           key.Binding
	CursorRight          key.Binding
	PlayStop             key.Binding
	TempoInputSwitch     key.Binding
	Increase             key.Binding
	Decrease             key.Binding
	ToggleAccentMode     key.Binding
	ToggleAccentModifier key.Binding
	OverlayInputSwitch   key.Binding
	SetupInputSwitch     key.Binding
	NextOverlay          key.Binding
	PrevOverlay          key.Binding
	Save                 key.Binding
	Undo                 key.Binding
	Redo                 key.Binding
	New                  key.Binding
	ToggleRatchetMode    key.Binding
}

type definitionKeyMap struct {
	TriggerAdd           key.Binding
	TriggerRemove        key.Binding
	OverlayTriggerRemove key.Binding
	ClearLine            key.Binding
	ClearSeq             key.Binding
	RatchetIncrease      key.Binding
	RatchetDecrease      key.Binding
	ToggleRatchetMute    key.Binding
	ActionAddLineReset   key.Binding
	ActionAddLineReverse key.Binding
	SelectKeyLine        key.Binding
	PressDownOverlay     key.Binding
	NumberPattern        key.Binding
}

var noteWiseKeys = []key.Binding{
	definitionKeys.TriggerAdd,
	definitionKeys.TriggerRemove,
	definitionKeys.OverlayTriggerRemove,
	definitionKeys.RatchetIncrease,
	definitionKeys.RatchetDecrease,
	definitionKeys.ToggleRatchetMute,
	definitionKeys.ActionAddLineReset,
	definitionKeys.ActionAddLineReverse,
}

var lineWiseKeys = []key.Binding{
	definitionKeys.ClearLine,
	definitionKeys.NumberPattern,
}

var overlayWiseKeys = []key.Binding{
	definitionKeys.ClearSeq,
}

func (dkm definitionKeyMap) IsNoteWiseKey(keyMsg tea.KeyMsg) bool {
	for _, kb := range noteWiseKeys {
		if key.Matches(keyMsg, kb) {
			return true
		}
	}
	return false
}

func (dkm definitionKeyMap) IsLineWiseKey(keyMsg tea.KeyMsg) bool {
	for _, kb := range lineWiseKeys {
		if key.Matches(keyMsg, kb) {
			return true
		}
	}
	return false
}

func (dkm definitionKeyMap) IsOverlayWiseKey(keyMsg tea.KeyMsg) bool {
	for _, kb := range overlayWiseKeys {
		if key.Matches(keyMsg, kb) {
			return true
		}
	}
	return false
}

func Key(help string, keyboardKey ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey...), key.WithHelp(keyboardKey[0], help))
}

var transitiveKeys = transitiveKeyMap{
	Quit:                 Key("Quit", "q"),
	Help:                 Key("Expand Help", "?"),
	CursorUp:             Key("Up", "k"),
	CursorDown:           Key("Down", "j"),
	CursorLeft:           Key("Left", "h"),
	CursorRight:          Key("Right", "l"),
	PlayStop:             Key("Play/Stop", " "),
	TempoInputSwitch:     Key("Select Tempo Indicator", "ctrl+t"),
	Increase:             Key("Tempo Increase", "+", "="),
	Decrease:             Key("Tempo Decrease", "-"),
	ToggleAccentMode:     Key("Toggle Accent Mode", "A"),
	ToggleAccentModifier: Key("Toggle Accent Modifier", "a"),
	OverlayInputSwitch:   Key("Select Overlay Indicator", "ctrl+o"),
	SetupInputSwitch:     Key("Setup Input Indicator", "ctrl+s"),
	NextOverlay:          Key("Next Overlay", "{"),
	PrevOverlay:          Key("Prev Overlay", "}"),
	Save:                 Key("Save", "ctrl+w"),
	Undo:                 Key("Undo", "u"),
	Redo:                 Key("Redo", "ctrl+r"),
	New:                  Key("New", "ctrl+n"),
	ToggleRatchetMode:    Key("Toggle Ratchet Mode", "ctrl+r"),
}

var definitionKeys = definitionKeyMap{
	TriggerAdd:           Key("Add Trigger", "f"),
	TriggerRemove:        Key("Remove Trigger", "d"),
	OverlayTriggerRemove: Key("Remove Overlay Note", "x"),
	ClearLine:            Key("Clear Line", "c"),
	ClearSeq:             Key("Clear Overlay", "C"),
	RatchetIncrease:      Key("Increase Ratchet", "R"),
	RatchetDecrease:      Key("Decrease Ratchet", "r"),
	ToggleRatchetMute:    Key("Toggle Ratchet Mute", "m"),
	ActionAddLineReset:   Key("Add Line Reset Action", "s"),
	ActionAddLineReverse: Key("Add Line Reverse Action", "S"),
	SelectKeyLine:        Key("Select Key Line", "K"),
	PressDownOverlay:     Key("Press Down Overlay", "ctrl+p"),
	NumberPattern:        Key("Number Pattern", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
}

// func (k keymap) ShortHelp() []key.Binding {
// 	return []key.Binding{
// 		k.Help, k.Quit,
// 	}
// }
//
// func (k keymap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{k.Help, k.Quit},
// 		{k.CursorUp, k.CursorDown, k.CursorLeft, k.CursorRight},
// 		{k.TriggerAdd, k.TriggerRemove},
// 	}
// }

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

type ratchet struct {
	hits   [9]bool
	length uint8
	span   uint8
}

func (r ratchet) Span() uint8 {
	return r.span + 1
}

func (r ratchet) Interval(beatInterval time.Duration) time.Duration {
	return (beatInterval * time.Duration(r.Span())) / time.Duration(r.length+1)
}

func InitRatchet() ratchet {
	return ratchet{
		hits: [9]bool{true, false, false, false, false, false, false, false, false},
	}
}

type note struct {
	AccentIndex uint8
	Ratchets    ratchet
	Action      action
}

func (n note) IncrementAccent(modifier int8) note {
	var newAccent = int8(n.AccentIndex) - modifier
	if newAccent >= 1 && newAccent < int8(len(accents)) {
		n.AccentIndex = uint8(newAccent)
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

type ratchetDiacritical string

var ratchets = []ratchetDiacritical{
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
	return fmt.Sprintf("Overlay-%d/%d", o.num, o.denom)
}

func (o overlayKey) MarshalTOML() ([]byte, error) {
	return []byte(fmt.Sprintf("%d/%d", o.num, o.denom)), nil
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

func (gk gridKey) String() string {
	return fmt.Sprintf("Grid-%0.2d-%0.2d", gk.line, gk.beat)
}

type overlay map[gridKey]note

type accentTarget uint8

const (
	ACCENT_TARGET_VELOCITY accentTarget = iota
	ACCENT_TARGET_NOTE
)

type lineDefinition struct {
	Channel uint8
	Note    uint8
	Target  accentTarget
}

func (l *lineDefinition) IncrementChannel() {
	if l.Channel < 16 {
		l.Channel++
	}
}

func (l *lineDefinition) DecrementChannel() {
	if l.Channel > 1 {
		l.Channel--
	}
}

func (l *lineDefinition) IncrementNote() {
	if l.Note < 128 {
		l.Note++
	}
}

func (l *lineDefinition) DecrementNote() {
	if l.Note > 1 {
		l.Note--
	}
}

type model struct {
	transitiveStatekeys       transitiveKeyMap
	definitionKeys            definitionKeyMap
	help                      help.Model
	cursor                    cursor.Model
	cursorPos                 gridKey
	outport                   drivers.Out
	logFile                   *os.File
	playing                   bool
	playTime                  time.Time
	trackTime                 time.Duration
	totalBeats                int
	playState                 []linestate
	tempoSelectionIndicator   uint8
	overlaySelectionIndicator uint8
	setupSelectionIndicator   uint8
	ratchetSelectionIndicator uint8
	accentMode                bool
	accentModifier            int8
	ratchetCursor             uint8
	overlayKey                overlayKey
	keyCycles                 int
	playingMatchedOverlays    []overlayKey
	undoStack                 UndoStack
	redoStack                 UndoStack
	// save everything below here
	definition Definition
}

type Undoable interface {
	ApplyUndo(m *model)
}

type UndoStack struct {
	undo Undoable
	redo Undoable
	next *UndoStack
}

var NIL_STACK = UndoStack{}

type UndoKeyline struct {
	keyline uint8
}

func (ukl UndoKeyline) ApplyUndo(m *model) {
	m.definition.keyline = ukl.keyline
}

type UndoBeats struct {
	beats uint8
}

func (ukl UndoBeats) ApplyUndo(m *model) {
	m.definition.beats = ukl.beats
}

type UndoTempo struct {
	tempo int
}

func (ukl UndoTempo) ApplyUndo(m *model) {
	m.definition.tempo = ukl.tempo
}

type UndoSubdivisions struct {
	subdivisions int
}

func (ukl UndoSubdivisions) ApplyUndo(m *model) {
	m.definition.subdivisions = ukl.subdivisions
}

type GridNote struct {
	gridKey gridKey
	note    note
}

type UndoGridNote struct {
	overlayKey
	gridNote GridNote
}

func (ugn UndoGridNote) ApplyUndo(m *model) {
	m.EnsureOverlayWithKey(ugn.overlayKey)
	overlay := m.definition.overlays[ugn.overlayKey]
	overlay[ugn.gridNote.gridKey] = ugn.gridNote.note
}

type UndoLineGridNotes struct {
	overlayKey
	line      uint8
	gridNotes []GridNote
}

func (ugn UndoLineGridNotes) ApplyUndo(m *model) {
	overlay := m.definition.overlays[ugn.overlayKey]
	for i := range m.definition.beats {
		delete(overlay, gridKey{ugn.line, i})
	}
	for _, gridNote := range ugn.gridNotes {
		overlay[gridNote.gridKey] = gridNote.note
	}
}

type UndoGridNotes struct {
	overlayKey
	gridNotes []GridNote
}

func (ugn UndoGridNotes) ApplyUndo(m *model) {
	m.EnsureOverlay()
	overlay := m.definition.overlays[ugn.overlayKey]
	for _, gridNote := range ugn.gridNotes {
		overlay[gridNote.gridKey] = gridNote.note
	}
}

type UndoToNothing struct {
	overlayKey overlayKey
	location   gridKey
}

func (utn UndoToNothing) ApplyUndo(m *model) {
	overlay := m.definition.overlays[utn.overlayKey]
	delete(overlay, utn.location)
}

type UndoLineToNothing struct {
	overlayKey overlayKey
	line       uint8
}

func (ultn UndoLineToNothing) ApplyUndo(m *model) {
	overlay := m.definition.overlays[ultn.overlayKey]
	for i := range m.definition.beats {
		delete(overlay, gridKey{ultn.line, i})
	}
}

type UndoNewOverlay struct {
	overlayKey overlayKey
}

func (uno UndoNewOverlay) ApplyUndo(m *model) {
	delete(m.definition.overlays, uno.overlayKey)
}

func (m *model) PushUndo(undo Undoable, redo Undoable) {
	if m.undoStack == NIL_STACK {
		m.undoStack = UndoStack{
			undo: undo,
			redo: redo,
			next: nil,
		}
	} else {
		pusheddown := m.undoStack
		lastin := UndoStack{
			undo: undo,
			redo: redo,
			next: &pusheddown,
		}
		m.undoStack = lastin
	}
}

func (m *model) PushRedo(undo Undoable, redo Undoable) {
	if m.redoStack == NIL_STACK {
		m.redoStack = UndoStack{
			undo: undo,
			redo: redo,
			next: nil,
		}
	} else {
		pusheddown := m.redoStack
		lastin := UndoStack{
			undo: undo,
			redo: redo,
			next: &pusheddown,
		}
		m.redoStack = lastin
	}
}

func (m *model) ResetRedo() {
	m.redoStack = NIL_STACK
}

func (m *model) PopUndo() UndoStack {
	firstout := m.undoStack
	if firstout != NIL_STACK && firstout.next != nil {
		m.undoStack = *m.undoStack.next
	} else {
		m.undoStack = NIL_STACK
	}
	return firstout
}

func (m *model) PopRedo() UndoStack {
	firstout := m.redoStack
	if firstout != NIL_STACK && firstout.next != nil {
		m.redoStack = *m.redoStack.next
	} else {
		m.redoStack = NIL_STACK
	}
	return firstout
}

func (m *model) Undo() UndoStack {
	undoStack := m.PopUndo()
	if undoStack != NIL_STACK {
		undoStack.undo.ApplyUndo(m)
	}
	return undoStack
}

func (m *model) Redo() UndoStack {
	undoStack := m.PopRedo()
	if undoStack != NIL_STACK {
		undoStack.redo.ApplyUndo(m)
	}
	return undoStack
}

type Definition struct {
	overlays     overlays
	lines        []lineDefinition
	beats        uint8
	tempo        int
	subdivisions int
	keyline      uint8
	metaOverlays map[overlayKey]metaOverlay
}

type metaOverlay struct {
	PressUp   bool
	PressDown bool
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

type ratchetMsg struct {
	lineNote
	iterations   uint8
	beatInterval time.Duration
}

func RatchetTick(ratchet lineNote, times uint8, beatInterval time.Duration) tea.Cmd {
	ratchetInterval := ratchet.Ratchets.Interval(beatInterval)
	return tea.Tick(
		ratchetInterval,
		func(t time.Time) tea.Msg { return ratchetMsg{ratchet, times, beatInterval} },
	)
}

func (m *model) BeatInterval() time.Duration {
	tickInterval := time.Minute / time.Duration(m.definition.tempo*m.definition.subdivisions)
	adjuster := time.Since(m.playTime) - m.trackTime
	m.trackTime = m.trackTime + tickInterval
	next := tickInterval - adjuster
	return next
}

type lineNote struct {
	note
	line lineDefinition
}

func PlayBeat(beatInterval time.Duration, lines []lineDefinition, pattern overlay, currentBeat []linestate, sendFn SendFunc) []tea.Cmd {
	notes := make([]lineNote, 0, len(lines))
	ratchetNotes := make([]lineNote, 0, len(lines))

	for i, line := range lines {
		currentGridKey := gridKey{uint8(i), currentBeat[i].currentBeat}
		note, hasNote := pattern[currentGridKey]
		if hasNote && note.Ratchets.length > 0 {
			ratchetNotes = append(ratchetNotes, lineNote{note, line})
		} else if hasNote && note != zeronote {
			notes = append(notes, lineNote{note, line})
		}
	}

	playCmds := make([]tea.Cmd, 0, len(lines))

	playNotes := func() tea.Msg {
		Play(notes, sendFn)
		return nil
	}

	playCmds = append(playCmds, playNotes)

	for _, ratchetNote := range ratchetNotes {
		playCmds = append(playCmds, func() tea.Msg {
			return ratchetMsg{ratchetNote, 0, beatInterval}
		})
	}

	return playCmds
}

func PlayRatchets(lineNote lineNote, beatInterval time.Duration, sendFn SendFunc) tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

type SendFunc func(msg midi.Message) error

func Play(lineNotes []lineNote, sendFn SendFunc) {
	for _, lineNote := range lineNotes {
		onMessage := midi.NoteOn(lineNote.line.Channel, lineNote.line.Note, accents[lineNote.note.AccentIndex].value)
		offMessage := midi.NoteOff(lineNote.line.Channel, lineNote.line.Note)
		err := sendFn(onMessage)
		if err != nil {
			panic("note on failed")
		}
		err = sendFn(offMessage)
		if err != nil {
			panic("note off failed")
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
	m.EnsureOverlayWithKey(m.overlayKey)
}

func (m *model) EnsureOverlayWithKey(key overlayKey) {
	if len(m.definition.overlays[key]) == 0 {
		m.definition.overlays[key] = make(overlay)
		if m.playing {
			m.determineMatachedOverlays()
		}
	}
}

func (m *model) CurrentNotable() Notable {
	var notable Notable
	overlay := m.definition.overlays[m.overlayKey]
	notable = &overlay
	return notable
}

func (m *model) AddTrigger() {
	m.CurrentNotable().SetNote(m.cursorPos, note{5, InitRatchet(), ACTION_NOTHING})
}

func (m *model) AddAction(act action) {
	m.CurrentNotable().SetNote(m.cursorPos, note{0, InitRatchet(), act})
}

func (m *model) RemoveTrigger() {
	m.CurrentNotable().SetNote(m.cursorPos, zeronote)
}

func (m *model) OverlayRemoveTrigger() {
	delete(m.definition.overlays[m.overlayKey], m.cursorPos)
}

func (m *model) IncreaseRatchet() {
	currentNote := m.CurrentNote()
	currentRatchet := currentNote.Ratchets.length
	if currentNote.Action == ACTION_NOTHING && currentRatchet+1 < uint8(len(ratchets)) {
		currentNote.Ratchets.length = currentRatchet + 1
		currentNote.Ratchets.hits[currentNote.Ratchets.length] = true
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) IncreaseSpan() {
	currentNote := m.CurrentNote()
	if currentNote != zeronote && currentNote.Action == ACTION_NOTHING {
		span := currentNote.Ratchets.span
		if span < 8 {
			currentNote.Ratchets.span = span + 1
		}
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) DecreaseSpan() {
	currentNote := m.CurrentNote()
	if currentNote != zeronote && currentNote.Action == ACTION_NOTHING {
		span := currentNote.Ratchets.span
		if span > 0 {
			currentNote.Ratchets.span = span - 1
		}
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) DecreaseRatchet() {
	currentNote := m.CurrentNote()
	currentRatchet := currentNote.Ratchets.length

	if currentNote.Action == ACTION_NOTHING && currentRatchet > 0 {
		currentNote.Ratchets.length = currentRatchet - 1
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) ToggleRatchetMute() {
	currentNote := m.CurrentNote()
	currentNote.Ratchets.hits[m.ratchetCursor] = !currentNote.Ratchets.hits[m.ratchetCursor]
	m.CurrentNotable().SetNote(m.cursorPos, currentNote)
}

func InitLines(n uint8) []lineDefinition {
	var lines = make([]lineDefinition, n)
	for i := range n {
		lines[i] = lineDefinition{
			Channel: 10,
			Note:    C1 + i,
			Target:  ACCENT_TARGET_VELOCITY,
		}
	}
	return lines
}

func InitPlayState(lines uint8) []linestate {
	linestates := make([]linestate, lines)
	for i := range linestates {
		linestates[i].direction = 1
		linestates[i].resetDirection = 1
		linestates[i].resetLocation = 0
	}
	return linestates
}

func InitDefinition() Definition {
	return Definition{
		overlays:     make(overlays),
		beats:        32,
		tempo:        120,
		keyline:      0,
		subdivisions: 2,
		lines:        InitLines(8),
		metaOverlays: make(map[overlayKey]metaOverlay),
	}
}

func InitModel(midiOutport drivers.Out) model {
	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic("could not open log file")
	}

	newCursor := cursor.New()
	newCursor.Style = lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	definition, hasDefinition := Read()

	if !hasDefinition {
		definition = InitDefinition()
	}

	return model{
		transitiveStatekeys: transitiveKeys,
		definitionKeys:      definitionKeys,
		help:                help.New(),
		cursor:              newCursor,
		outport:             midiOutport,
		logFile:             logFile,
		cursorPos:           gridKey{0, 0},
		accentModifier:      1,
		overlayKey:          ROOT_OVERLAY,
		definition:          definition,
	}
}

func (m model) LogTeaMsg(msg tea.Msg) {
	switch msg := msg.(type) {
	case beatMsg:
		m.LogString(fmt.Sprintf("beatMsg %d %d %d\n", msg.interval, m.totalBeats, m.definition.tempo))
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

func GetMinimumKeyCycle(key overlayKey) int {
	for i := 0; i < 100; i++ {
		if DoesKeyMatch(i, key) {
			return i
		}
	}
	return 100
}

func (d Definition) GetMatchingOverlays(keyCycles int, keys []overlayKey) []overlayKey {
	var matchedKeys = make([]overlayKey, 0, 5)

	slices.SortFunc(keys, OverlayKeySort)
	var pressNext = false

	for _, key := range keys {
		matches := DoesKeyMatch(keyCycles, key)
		if (matches && len(matchedKeys) == 0) || pressNext {
			matchedKeys = append(matchedKeys, key)
			if d.metaOverlays[key].PressDown {
				pressNext = true
			} else {
				pressNext = false
			}
		} else if matches && len(matchedKeys) != 0 {
			if d.metaOverlays[key].PressUp {
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
	keys := transitiveKeys

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case Is(msg, keys.Quit):
			m.logFile.Close()
			return m, tea.Quit
		case Is(msg, keys.CursorDown):
			if m.ratchetSelectionIndicator == 0 {
				if m.cursorPos.line < uint8(len(m.definition.lines)-1) {
					m.cursorPos.line++
				}
			}
		case Is(msg, keys.CursorUp):
			if m.ratchetSelectionIndicator == 0 {
				if m.cursorPos.line > 0 {
					m.cursorPos.line--
				}
			}
		case Is(msg, keys.CursorLeft):
			if m.ratchetSelectionIndicator > 0 {
				if m.ratchetCursor > 0 {
					m.ratchetCursor--
				}
			} else {
				if m.cursorPos.beat > 0 {
					m.cursorPos.beat--
				}
			}
		case Is(msg, keys.CursorRight):
			if m.ratchetSelectionIndicator > 0 {
				currentNote := m.CurrentNote()
				if m.ratchetCursor < currentNote.Ratchets.length {
					m.ratchetCursor++
				}
			} else {
				if m.cursorPos.beat < m.definition.beats-1 {
					m.cursorPos.beat++
				}
			}
		case Is(msg, keys.PlayStop):
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
				m.playState = InitPlayState(uint8(len(m.definition.lines)))
				m.advanceKeyCycle()
				m.trackTime = time.Duration(0)
				sendFn, err := midi.SendTo(m.outport)
				if err != nil {
					panic("sendFn is broken")
				}
				beatInterval := m.BeatInterval()

				cmds := make([]tea.Cmd, 0, 10)
				cmds = append(cmds, PlayBeat(beatInterval, m.definition.lines, m.definition.CombinedPattern(m.playingMatchedOverlays), m.playState, sendFn)...)
				cmds = append(cmds, BeatTick(beatInterval))
				return m, tea.Batch(cmds...)
			} else {
				m.keyCycles = 0
				m.playingMatchedOverlays = []overlayKey{}
			}
		case Is(msg, keys.OverlayInputSwitch):
			m.overlaySelectionIndicator = (m.overlaySelectionIndicator + 1) % 3
			m.tempoSelectionIndicator = 0
			m.setupSelectionIndicator = 0
			m.ratchetSelectionIndicator = 0
		case Is(msg, keys.TempoInputSwitch):
			m.tempoSelectionIndicator = (m.tempoSelectionIndicator + 1) % 3
			m.overlaySelectionIndicator = 0
			m.setupSelectionIndicator = 0
			m.ratchetSelectionIndicator = 0
		case Is(msg, keys.SetupInputSwitch):
			m.setupSelectionIndicator = (m.setupSelectionIndicator + 1) % 3
			m.overlaySelectionIndicator = 0
			m.tempoSelectionIndicator = 0
			m.ratchetSelectionIndicator = 0
		case Is(msg, keys.Increase):
			switch {
			case m.tempoSelectionIndicator == 1:
				if m.definition.tempo < 300 {
					m.definition.tempo++
				}
			case m.tempoSelectionIndicator == 2:
				if m.definition.subdivisions < 8 {
					m.definition.subdivisions++
				}
			case m.overlaySelectionIndicator == 1:
				m.overlayKey.IncrementNumerator()
			case m.overlaySelectionIndicator == 2:
				m.overlayKey.IncrementDenominator()
			case m.setupSelectionIndicator == 1:
				m.definition.lines[m.cursorPos.line].IncrementChannel()
			case m.setupSelectionIndicator == 2:
				m.definition.lines[m.cursorPos.line].IncrementNote()
			case m.ratchetSelectionIndicator == 2:
				m.IncreaseSpan()
			}
		case Is(msg, keys.Decrease):
			switch {
			case m.tempoSelectionIndicator == 1:
				if m.definition.tempo > 30 {
					m.definition.tempo--
				}
			case m.tempoSelectionIndicator == 2:
				if m.definition.subdivisions > 1 {
					m.definition.subdivisions--
				}
			case m.overlaySelectionIndicator == 1:
				m.overlayKey.DecrementNumerator()
			case m.overlaySelectionIndicator == 2:
				m.overlayKey.DecrementDenominator()
			case m.setupSelectionIndicator == 1:
				m.definition.lines[m.cursorPos.line].DecrementChannel()
			case m.setupSelectionIndicator == 2:
				m.definition.lines[m.cursorPos.line].DecrementNote()
			case m.ratchetSelectionIndicator == 2:
				m.DecreaseSpan()
			}
		case Is(msg, keys.ToggleAccentMode):
			m.accentMode = !m.accentMode
		case Is(msg, keys.ToggleRatchetMode):
			m.ratchetCursor = 0
			if m.CurrentNote() != zeronote {
				m.ratchetSelectionIndicator = (m.ratchetSelectionIndicator + 1) % 3
			}
			m.tempoSelectionIndicator = 0
			m.setupSelectionIndicator = 0
			m.overlaySelectionIndicator = 0
		case Is(msg, keys.ToggleAccentModifier):
			m.accentModifier = -1 * m.accentModifier
		case Is(msg, keys.PrevOverlay):
			m.NextOverlay(-1)
		case Is(msg, keys.NextOverlay):
			m.NextOverlay(+1)
		case Is(msg, keys.Save):
			m.Save()
		case Is(msg, keys.Undo):
			undoStack := m.Undo()
			if undoStack != NIL_STACK {
				m.PushRedo(undoStack.undo, undoStack.redo)
			}
		case Is(msg, keys.Redo):
			undoStack := m.Redo()
			if undoStack != NIL_STACK {
				m.PushUndo(undoStack.undo, undoStack.redo)
			}
		case Is(msg, keys.New):
			m.cursorPos = gridKey{0, 0}
			m.overlayKey = ROOT_OVERLAY
			m.accentModifier = 1
			m.tempoSelectionIndicator = 0
			m.setupSelectionIndicator = 0
			m.overlaySelectionIndicator = 0
			m.ratchetSelectionIndicator = 0
			m.definition = InitDefinition()
		default:
			m = m.UpdateDefinition(msg)
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
			cmds := make([]tea.Cmd, 0, 10)
			cmds = append(cmds, PlayBeat(beatInterval, m.definition.lines, m.definition.CombinedPattern(m.playingMatchedOverlays), m.playState, sendFn)...)
			cmds = append(cmds, BeatTick(beatInterval))
			return m, tea.Batch(
				cmds...,
			)
		}
	case ratchetMsg:
		if m.playing && msg.iterations < (msg.Ratchets.length+1) {
			var playCmd tea.Cmd
			var ratchetTickCmd tea.Cmd
			sendFn, err := midi.SendTo(m.outport)
			if err != nil {
				panic("sendFn is broken")
			}
			if msg.Ratchets.hits[msg.iterations] {
				playCmd = func() tea.Msg {
					Play([]lineNote{msg.lineNote}, sendFn)
					return nil
				}
			}
			if msg.iterations+1 < (msg.Ratchets.length + 1) {
				ratchetTickCmd = RatchetTick(msg.lineNote, msg.iterations+1, msg.beatInterval)
			}
			return m, tea.Batch(playCmd, ratchetTickCmd)
		}
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor
	return m, cmd
}

func (m model) UpdateDefinitionKeys(msg tea.KeyMsg) model {
	keys := definitionKeys
	switch {
	case Is(msg, keys.TriggerAdd):
		m.AddTrigger()
	case Is(msg, keys.TriggerRemove):
		m.RemoveTrigger()
	case Is(msg, keys.OverlayTriggerRemove):
		m.OverlayRemoveTrigger()
	case Is(msg, keys.ClearLine):
		m.ClearOverlayLine()
	case Is(msg, keys.RatchetIncrease):
		m.IncreaseRatchet()
	case Is(msg, keys.RatchetDecrease):
		m.DecreaseRatchet()
	case Is(msg, keys.ToggleRatchetMute):
		m.ToggleRatchetMute()
	case Is(msg, keys.ActionAddLineReset):
		m.AddAction(ACTION_LINE_RESET)
	case Is(msg, keys.ActionAddLineReverse):
		m.AddAction(ACTION_LINE_REVERSE)
	case Is(msg, keys.SelectKeyLine):
		undoable := UndoKeyline{m.definition.keyline}
		m.definition.keyline = m.cursorPos.line
		redoable := UndoKeyline{m.definition.keyline}
		m.PushUndo(undoable, redoable)
	case Is(msg, keys.PressDownOverlay):
		m.ToggleOverlayStackOptions(m.overlayKey)
	case Is(msg, keys.ClearSeq):
		m.ClearOverlay()
	}
	if msg.String() >= "1" && msg.String() <= "9" {
		beatInterval, _ := strconv.ParseInt(msg.String(), 0, 8)
		if m.accentMode {
			m.incrementAccent(uint8(beatInterval))
		} else {
			m.fill(uint8(beatInterval))
		}
	}
	return m
}

func (m model) UpdateDefinition(msg tea.KeyMsg) model {
	keys := definitionKeys
	if keys.IsNoteWiseKey(msg) {
		undoable := m.UndoableNote()
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableNote()
		m.PushUndo(undoable, redoable)
	} else if keys.IsLineWiseKey(msg) {
		undoable := m.UndoableLine()
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableLine()
		m.PushUndo(undoable, redoable)
	} else if keys.IsOverlayWiseKey(msg) {
		undoable := m.UndoableOverlay()
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableOverlay()
		m.PushUndo(undoable, redoable)
	} else {
		m = m.UpdateDefinitionKeys(msg)
	}
	m.ResetRedo()
	return m
}

func (m model) UndoableNote() Undoable {
	overlay, hasOverlay := m.definition.overlays[m.overlayKey]
	if !hasOverlay {
		return UndoNewOverlay{m.overlayKey}
	}
	currentNote, hasNote := overlay[m.cursorPos]
	if hasNote {
		return UndoGridNote{m.overlayKey, GridNote{m.cursorPos, currentNote}}
	} else {
		return UndoToNothing{m.overlayKey, m.cursorPos}
	}
}

func (m model) UndoableLine() Undoable {
	overlay, hasOverlay := m.definition.overlays[m.overlayKey]
	if !hasOverlay {
		return UndoNewOverlay{m.overlayKey}
	}
	notesToUndo := make([]GridNote, 0, m.definition.beats)
	for i := range m.definition.beats {
		key := gridKey{m.cursorPos.line, i}
		currentNote, hasNote := overlay[key]
		if hasNote {
			notesToUndo = append(notesToUndo, GridNote{key, currentNote})
		}
	}
	if len(notesToUndo) == 0 {
		return UndoLineToNothing{m.overlayKey, m.cursorPos.line}
	}
	return UndoLineGridNotes{m.overlayKey, m.cursorPos.line, notesToUndo}
}

func (m model) UndoableOverlay() Undoable {
	_, hasOverlay := m.definition.overlays[m.overlayKey]
	if !hasOverlay {
		return UndoNewOverlay{m.overlayKey}
	}
	notesToUndo := make([]GridNote, 0, m.definition.beats)
	for key, note := range m.definition.overlays[m.overlayKey] {
		notesToUndo = append(notesToUndo, GridNote{key, note})
	}
	return UndoGridNotes{m.overlayKey, notesToUndo}
}

func (m model) Save() {
	Write(m.definition)
}

func (m model) SerializeLines() string {
	// buf strings.Builder
	// for line := range m.lines {
	// 	buf.WriteString()
	// }
	return ""
}
func (m model) SerializeBeat() string         { return "" }
func (m model) SerializeTempo() string        { return "" }
func (m model) SerializeKeyLine() string      { return "" }
func (m model) SerializeOverlays() string     { return "" }
func (m model) SerializeMetaOverlays() string { return "" }

func (m *model) ToggleOverlayStackOptions(key overlayKey) {
	meta, hasMeta := m.definition.metaOverlays[m.overlayKey]

	if !hasMeta {
		m.definition.metaOverlays[m.overlayKey] = metaOverlay{PressUp: true, PressDown: false}
	} else if meta.PressUp {
		m.definition.metaOverlays[m.overlayKey] = metaOverlay{PressUp: false, PressDown: true}
	} else {
		m.definition.metaOverlays[m.overlayKey] = metaOverlay{PressUp: false, PressDown: false}
	}
}

func (m model) CurrentNote() note {
	matchedKeys := m.EditKeys()
	combinedOverlay := m.definition.CombinedPattern(matchedKeys)
	return combinedOverlay[m.cursorPos]
}

func RemoveRootKey(keys []overlayKey) []overlayKey {
	index := slices.Index(keys, ROOT_OVERLAY)
	if index >= 0 {
		return append(keys[:index], keys[index+1:]...)
	}
	return keys
}

func (m *model) NextOverlay(direction int) {
	keys := m.OverlayKeys()
	slices.SortFunc(keys, OverlayKeySort)
	index := slices.Index(keys, m.overlayKey)
	if index+direction < len(keys) && index+direction >= 0 {
		m.overlayKey = keys[index+direction]
	}
}

func (m *model) ClearOverlayLine() {
	for i := uint8(0); i < m.definition.beats; i++ {
		key := gridKey{m.cursorPos.line, i}
		delete(m.definition.overlays[m.overlayKey], key)
	}
}

func (m *model) ClearOverlay() {
	delete(m.definition.overlays, m.overlayKey)
}

func (m *model) advanceCurrentBeat() {
	combinedPattern := m.definition.CombinedPattern(m.playingMatchedOverlays)
	for i, currentState := range m.playState {
		advancedBeat := int8(currentState.currentBeat) + currentState.direction

		if advancedBeat < 0 || advancedBeat >= int8(m.definition.beats) {
			m.playState[i].currentBeat = currentState.resetLocation
			advancedBeat = int8(currentState.resetLocation)
			m.playState[i].direction = currentState.resetDirection
		} else {
			m.playState[i].currentBeat = uint8(advancedBeat)
		}

		switch combinedPattern[gridKey{uint8(i), uint8(advancedBeat)}].Action {
		case ACTION_LINE_RESET:
			m.playState[i].currentBeat = 0
		case ACTION_LINE_REVERSE:
			m.playState[i].currentBeat = uint8(max(advancedBeat-1, 0))
			m.playState[i].direction = -1
		}
	}
}

func (m *model) advanceKeyCycle() {
	if m.playState[m.definition.keyline].currentBeat == 0 {
		m.keyCycles++
		m.determineMatachedOverlays()
	}
}

func (m *model) determineMatachedOverlays() {
	keys := m.OverlayKeys()
	m.playingMatchedOverlays = m.definition.GetMatchingOverlays(m.keyCycles, keys)
}

func (d Definition) CombinedPattern(keys []overlayKey) overlay {
	var combinedOverlay = make(overlay)

	for _, key := range slices.Backward(keys) {
		for gridKey, note := range d.overlays[key] {
			combinedOverlay[gridKey] = note
		}
	}
	return combinedOverlay
}

type VisualNote struct {
	origin overlayKey
	note
}

type VisualOverlay map[gridKey]VisualNote

func (d Definition) VisualCombinedPattern(keys []overlayKey) VisualOverlay {
	var combinedOverlay = make(VisualOverlay)

	for _, key := range slices.Backward(keys) {
		for gridKey, note := range d.overlays[key] {
			combinedOverlay[gridKey] = VisualNote{key, note}
		}
	}
	return combinedOverlay
}

func (m *model) KeysBelowCurrent() []overlayKey {
	keys := m.OverlayKeys()
	slices.SortFunc(keys, OverlayKeySort)
	slices.Reverse(keys)
	indexOfCurrent := slices.Index(keys, m.overlayKey)
	if indexOfCurrent >= 0 {
		return keys[:indexOfCurrent]
	} else {
		return []overlayKey{}
	}
}

func (m *model) fill(every uint8) {
	start := m.cursorPos.beat

	matchedKeys := m.EditKeys()
	combinedOverlay := m.definition.CombinedPattern(matchedKeys)

	for i := uint8(0); i < m.definition.beats; i++ {
		if i%every == 0 {
			currentBeat := start + i
			if currentBeat >= m.definition.beats {
				return
			}
			gridKey := gridKey{m.cursorPos.line, currentBeat}
			currentNote, hasNote := combinedOverlay[gridKey]
			hasNote = hasNote && currentNote != zeronote

			if m.overlayKey != ROOT_OVERLAY && hasNote {
				m.CurrentNotable().SetNote(gridKey, zeronote)
			} else if m.overlayKey == ROOT_OVERLAY && hasNote {
				delete(m.definition.overlays[ROOT_OVERLAY], gridKey)
			} else {
				m.CurrentNotable().SetNote(gridKey, note{5, InitRatchet(), 0})
			}
		}
	}
}

func (m *model) EditKeys() []overlayKey {
	keysBelowCurrent := m.KeysBelowCurrent()
	matchedKeys := m.definition.GetMatchingOverlays(GetMinimumKeyCycle(m.overlayKey), keysBelowCurrent)
	matchedKeys = append(matchedKeys, m.overlayKey)
	slices.SortFunc(matchedKeys, OverlayKeySort)
	return matchedKeys
}

func (m *model) CurrentKeys() []overlayKey {
	if m.playing {
		return m.playingMatchedOverlays
	} else {
		return m.EditKeys()
	}
}

func (m *model) incrementAccent(every uint8) {
	start := m.cursorPos.beat
	rootOverlay := m.definition.overlays[ROOT_OVERLAY]
	currentOverlay := m.definition.overlays[m.overlayKey]

	for i := uint8(0); i < m.definition.beats; i++ {
		if i%every == 0 {
			currentBeat := start + i
			if currentBeat >= m.definition.beats {
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
	return slices.AppendSeq(keys, maps.Keys(m.definition.overlays))
}

var heartColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed3902"))
var selectedColor = lipgloss.NewStyle().Background(lipgloss.Color("#5cdffb")).Foreground(lipgloss.Color("#000000"))
var numberColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcbd15"))

func (m model) TempoView() string {
	var buf strings.Builder
	var tempo, division string
	switch m.tempoSelectionIndicator {
	case 0:
		tempo = numberColor.Render(strconv.Itoa(m.definition.tempo))
		division = numberColor.Render(strconv.Itoa(m.definition.subdivisions))
	case 1:
		tempo = selectedColor.Render(strconv.Itoa(m.definition.tempo))
		division = numberColor.Render(strconv.Itoa(m.definition.subdivisions))
	case 2:
		tempo = numberColor.Render(strconv.Itoa(m.definition.tempo))
		division = selectedColor.Render(strconv.Itoa(m.definition.subdivisions))
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
	} else if len(m.definition.overlays) > 0 && m.setupSelectionIndicator == 0 {
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
	for i, line := range m.definition.lines {
		buf.WriteString("CH ")
		if uint8(i) == m.cursorPos.line && m.setupSelectionIndicator == 1 {
			buf.WriteString(selectedColor.Render(fmt.Sprintf("%2d", line.Channel)))
		} else {
			buf.WriteString(numberColor.Render(fmt.Sprintf("%2d", line.Channel)))
		}
		buf.WriteString(" NOTE ")
		if uint8(i) == m.cursorPos.line && m.setupSelectionIndicator == 2 {
			buf.WriteString(selectedColor.Render(strconv.Itoa(int(line.Note))))
		} else {

			buf.WriteString(numberColor.Render(strconv.Itoa(int(line.Note))))
		}
		buf.WriteString("\n")
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
		if m.definition.metaOverlays[k].PressDown {
			stackModifier = " \u2193\u0332"
		} else if m.definition.metaOverlays[k].PressUp {
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
	visualCombinedPattern := m.definition.VisualCombinedPattern(m.CurrentKeys())

	if m.accentMode {
		if m.accentModifier > 0 {
			mode = " Accent Mode \u2191 "
		} else {
			mode = " Accent Mode \u2193 "
		}
		buf.WriteString(fmt.Sprintf("   Seq - %s\n", accentModeStyle.Render(mode)))
	} else if m.ratchetSelectionIndicator > 0 {
		buf.WriteString(m.RatchetModeView())
	} else if m.playing {
		buf.WriteString(fmt.Sprintf("   Seq - Playing - %d\n", m.keyCycles))
	} else {
		buf.WriteString("   Seq - A sequencer for your cli\n")
	}
	buf.WriteString("  ┌─────────────────────────────────\n")
	for i := uint8(0); i < uint8(len(m.definition.lines)); i++ {
		buf.WriteString(lineView(i, m, visualCombinedPattern))
	}
	buf.WriteString(m.CurrentOverlayView())
	buf.WriteString("\n")
	// buf.WriteString(m.help.View(m.keys))
	// buf.WriteString("\n")
	return buf.String()
}

var activeRatchetColor lipgloss.Color = "#abfaa9"
var mutedRatchetColor lipgloss.Color = "#f34213"

func (m model) RatchetModeView() string {
	activeStyle := lipgloss.NewStyle().Foreground(activeRatchetColor)
	mutedStyle := lipgloss.NewStyle().Foreground(mutedRatchetColor)

	currentNote := m.CurrentNote()

	var buf strings.Builder
	var ratchetsBuf strings.Builder
	buf.WriteString("   Ratchets ")
	for i := range uint8(8) {
		var backgroundColor lipgloss.Color
		if i <= currentNote.Ratchets.length {
			if m.ratchetCursor == i && m.ratchetSelectionIndicator == 1 {
				backgroundColor = lipgloss.Color("#5cdffb")
			}
			if currentNote.Ratchets.hits[i] {
				ratchetsBuf.WriteString(activeStyle.Background(backgroundColor).Render("\u25CF"))
			} else {
				ratchetsBuf.WriteString(mutedStyle.Background(backgroundColor).Render("\u25C9"))
			}
			ratchetsBuf.WriteString(" ")
		} else {

			ratchetsBuf.WriteString("  ")
		}
	}
	buf.WriteString(fmt.Sprintf("%*s", 32, ratchetsBuf.String()))
	if m.ratchetSelectionIndicator == 2 {
		buf.WriteString(fmt.Sprintf(" Span %s ", selectedColor.Render(strconv.Itoa(int(currentNote.Ratchets.Span())))))
	} else {
		buf.WriteString(fmt.Sprintf(" Span %s ", numberColor.Render(strconv.Itoa(int(currentNote.Ratchets.Span())))))
	}
	buf.WriteString("\n")

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

func lineView(lineNumber uint8, m model, visualCombinedPattern VisualOverlay) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d%s│", lineNumber, KeyLineIndicator(m.definition.keyline, lineNumber)))

	for i := uint8(0); i < m.definition.beats; i++ {
		currentGridKey := gridKey{uint8(lineNumber), i}
		overlayNote, hasNote := visualCombinedPattern[currentGridKey]

		var backgroundSeqColor lipgloss.Color
		if m.playing && m.playState[lineNumber].currentBeat == i {
			backgroundSeqColor = seqCursorColor
		} else if hasNote && overlayNote.origin != ROOT_OVERLAY {
			backgroundSeqColor = seqOverlayColor
		} else if i%8 > 3 {
			backgroundSeqColor = altSeqColor
		} else {
			backgroundSeqColor = seqColor
		}

		char, foregroundColor := overlayNote.ViewComponents()

		style := lipgloss.NewStyle().Background(backgroundSeqColor)
		if m.cursorPos.line == uint8(lineNumber) && m.cursorPos.beat == i {
			m.cursor.SetChar(char)
			char = m.cursor.View()
		} else {
			style = style.Foreground(foregroundColor)
		}
		buf.WriteString(style.Render(char))
	}

	buf.WriteString("\n")
	return buf.String()
}

func (n note) ViewComponents() (string, lipgloss.Color) {

	currentNote := n
	currentAccent := accents[currentNote.AccentIndex]
	currentAction := currentNote.Action
	var char string
	var foregroundColor lipgloss.Color
	if currentAction == ACTION_NOTHING {
		char = string(currentAccent.shape) + string(ratchets[currentNote.Ratchets.length])
		foregroundColor = currentAccent.color
	} else {
		lineaction := lineactions[currentAction]
		char = string(lineaction.shape)
		foregroundColor = lineaction.color
	}

	return char, foregroundColor
}
