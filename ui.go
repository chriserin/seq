package main

import (
	"fmt"
	"maps"
	"math/rand"
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
	overlaykey "github.com/chriserin/seq/overlayKey"
	midi "gitlab.com/gomidi/midi/v2"
)

type transitiveKeyMap struct {
	Quit               key.Binding
	Help               key.Binding
	CursorUp           key.Binding
	CursorDown         key.Binding
	CursorLeft         key.Binding
	CursorRight        key.Binding
	CursorLineStart    key.Binding
	CursorLineEnd      key.Binding
	Escape             key.Binding
	PlayStop           key.Binding
	TempoInputSwitch   key.Binding
	Increase           key.Binding
	Decrease           key.Binding
	ToggleAccentMode   key.Binding
	OverlayInputSwitch key.Binding
	SetupInputSwitch   key.Binding
	AccentInputSwitch  key.Binding
	NextOverlay        key.Binding
	PrevOverlay        key.Binding
	Save               key.Binding
	Undo               key.Binding
	Redo               key.Binding
	New                key.Binding
	ToggleRatchetMode  key.Binding
	ToggleVisualMode   key.Binding
	NewLine            key.Binding
	Yank               key.Binding
	Mute               key.Binding
	Solo               key.Binding
}

type definitionKeyMap struct {
	TriggerAdd           key.Binding
	TriggerRemove        key.Binding
	AccentIncrease       key.Binding
	AccentDecrease       key.Binding
	OverlayTriggerRemove key.Binding
	ClearLine            key.Binding
	ClearSeq             key.Binding
	RatchetIncrease      key.Binding
	RatchetDecrease      key.Binding
	ActionAddLineReset   key.Binding
	ActionAddLineReverse key.Binding
	ActionAddSkipBeat    key.Binding
	ActionAddReset       key.Binding
	ActionAddLineBounce  key.Binding
	ActionAddLineDelay   key.Binding
	SelectKeyLine        key.Binding
	PressDownOverlay     key.Binding
	NumberPattern        key.Binding
	RotateRight          key.Binding
	RotateLeft           key.Binding
	Paste                key.Binding
}

var noteWiseKeys = []key.Binding{
	definitionKeys.TriggerAdd,
	definitionKeys.TriggerRemove,
	definitionKeys.AccentIncrease,
	definitionKeys.AccentDecrease,
	definitionKeys.OverlayTriggerRemove,
	definitionKeys.RatchetIncrease,
	definitionKeys.RatchetDecrease,
	definitionKeys.ActionAddLineReset,
	definitionKeys.ActionAddLineReverse,
	definitionKeys.ActionAddSkipBeat,
	definitionKeys.ActionAddReset,
	definitionKeys.ActionAddLineBounce,
	definitionKeys.ActionAddLineDelay,
}

var lineWiseKeys = []key.Binding{
	definitionKeys.ClearLine,
	definitionKeys.NumberPattern,
	definitionKeys.RotateRight,
	definitionKeys.RotateLeft,
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
	Quit:               Key("Quit", "q"),
	Help:               Key("Expand Help", "?"),
	CursorUp:           Key("Up", "k"),
	CursorDown:         Key("Down", "j"),
	CursorLeft:         Key("Left", "h"),
	CursorRight:        Key("Right", "l"),
	CursorLineStart:    Key("Line Start", "<"),
	CursorLineEnd:      Key("Line End", ">"),
	Escape:             Key("Escape", "esc", "enter"),
	PlayStop:           Key("Play/Stop", " "),
	TempoInputSwitch:   Key("Select Tempo Indicator", "ctrl+t"),
	Increase:           Key("Tempo Increase", "+", "="),
	Decrease:           Key("Tempo Decrease", "-"),
	ToggleAccentMode:   Key("Toggle Accent Mode", "g"),
	OverlayInputSwitch: Key("Select Overlay Indicator", "ctrl+o"),
	SetupInputSwitch:   Key("Setup Input Indicator", "ctrl+s"),
	AccentInputSwitch:  Key("Accent Input Indicator", "ctrl+a"),
	NextOverlay:        Key("Next Overlay", "{"),
	PrevOverlay:        Key("Prev Overlay", "}"),
	Save:               Key("Save", "ctrl+w"),
	Undo:               Key("Undo", "u"),
	Redo:               Key("Redo", "U"),
	New:                Key("New", "ctrl+n"),
	ToggleRatchetMode:  Key("Toggle Ratchet Mode", "ctrl+r"),
	ToggleVisualMode:   Key("Toggle Visual Mode", "v"),
	NewLine:            Key("New Line", "ctrl+l"),
	Yank:               Key("Yank", "y"),
	Mute:               Key("Mute", "m"),
	Solo:               Key("Solo", "M"),
}

var definitionKeys = definitionKeyMap{
	TriggerAdd:           Key("Add Trigger", "f"),
	TriggerRemove:        Key("Remove Trigger", "d"),
	AccentIncrease:       Key("Accent Increase", "A"),
	AccentDecrease:       Key("Accent Increase", "a"),
	OverlayTriggerRemove: Key("Remove Overlay Note", "x"),
	ClearLine:            Key("Clear Line", "c"),
	ClearSeq:             Key("Clear Overlay", "C"),
	RatchetIncrease:      Key("Increase Ratchet", "R"),
	RatchetDecrease:      Key("Decrease Ratchet", "r"),
	ActionAddLineReset:   Key("Add Line Reset Action", "s"),
	ActionAddLineReverse: Key("Add Line Reverse Action", "S"),
	ActionAddLineBounce:  Key("Add Line Bounce", "B"),
	ActionAddLineDelay:   Key("Add Line Delay", "z"),
	ActionAddSkipBeat:    Key("Add Skip Beat", "b"),
	ActionAddReset:       Key("Add Pattern Reset", "T"),
	SelectKeyLine:        Key("Select Key Line", "K"),
	PressDownOverlay:     Key("Press Down Overlay", "ctrl+p"),
	NumberPattern:        Key("Number Pattern", "1", "2", "3", "4", "5", "6", "7", "8", "9", "!", "@", "#", "$", "%", "^", "&", "*", "("),
	RotateRight:          Key("Right Right", "L"),
	RotateLeft:           Key("Right Left", "H"),
	Paste:                Key("Paste", "p"),
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
	Shape rune
	Color lipgloss.Color
	Value uint8
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
	Hits   [9]bool
	Length uint8
	Span   uint8
}

func (r ratchet) GetSpan() uint8 {
	return r.Span + 1
}

func (r ratchet) Interval(beatInterval time.Duration) time.Duration {
	return (beatInterval * time.Duration(r.GetSpan())) / time.Duration(r.Length+1)
}

func InitRatchet() ratchet {
	return ratchet{
		Hits: [9]bool{true, false, false, false, false, false, false, false, false},
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
	ACTION_LINE_SKIP_BEAT
	ACTION_RESET
	ACTION_LINE_BOUNCE
	ACTION_LINE_DELAY
)

var lineactions = map[action]lineaction{
	ACTION_NOTHING:        {' ', "#000000"},
	ACTION_LINE_RESET:     {'↔', "#cf142b"},
	ACTION_LINE_REVERSE:   {'←', "#f8730e"},
	ACTION_LINE_SKIP_BEAT: {'⇒', "#a9e5bb"},
	ACTION_RESET:          {'⇚', "#fcf6b1"},
	ACTION_LINE_BOUNCE:    {'↨', "#fcf6b1"},
	ACTION_LINE_DELAY:     {'ℤ', "#cc4bc2"},
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

type groupPlayState uint

const (
	PLAY_STATE_PLAY groupPlayState = iota
	PLAY_STATE_MUTE
	PLAY_STATE_SOLO
	PLAY_STATE_MUTED_BY_SOLO
)

type linestate struct {
	currentBeat         uint8
	direction           int8
	resetDirection      int8
	resetLocation       uint8
	resetActionLocation uint8
	resetAction         action
	groupPlayState      groupPlayState
}

type overlayKey = overlaykey.OverlayPeriodicity

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

type lineDefinition struct {
	Channel uint8
	Note    uint8
}

type noteMessage struct {
	channel  uint8
	note     uint8
	velocity uint8
}

func (l lineDefinition) Message(note note, accentValue uint8, accentTarget accentTarget) noteMessage {
	var noteValue uint8
	var velocityValue uint8
	switch accentTarget {
	case ACCENT_TARGET_NOTE:
		noteValue = l.Note + accentValue
		velocityValue = 96
	case ACCENT_TARGET_VELOCITY:
		noteValue = l.Note
		velocityValue = accentValue
	}

	return noteMessage{l.Channel, noteValue, velocityValue}
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
	transitiveStatekeys    transitiveKeyMap
	definitionKeys         definitionKeyMap
	help                   help.Model
	cursor                 cursor.Model
	cursorPos              gridKey
	visualAnchorCursor     gridKey
	visualMode             bool
	midiConnection         MidiConnection
	logFile                *os.File
	playing                bool
	playTime               time.Time
	trackTime              time.Duration
	totalBeats             int
	playState              []linestate
	selectionIndicator     Selection
	accentMode             bool
	ratchetCursor          uint8
	overlayKey             overlayKey
	keyCycles              int
	playingMatchedOverlays []overlayKey
	undoStack              UndoStack
	redoStack              UndoStack
	yankBuffer             Buffer
	needsWrite             int
	// save everything below here
	definition Definition
}

type Selection uint8

const (
	SELECT_NOTHING Selection = iota
	SELECT_TEMPO
	SELECT_TEMPO_SUBDIVISION
	SELECT_OVERLAY_NUM
	SELECT_OVERLAY_DENOM
	SELECT_SETUP_CHANNEL
	SELECT_SETUP_NOTE
	SELECT_RATCHETS
	SELECT_RATCHET_SPAN
	SELECT_ACCENT_DIFF
	SELECT_ACCENT_TARGET
	SELECT_ACCENT_START
)

type Undoable interface {
	ApplyUndo(m *model)
}

type UndoStack struct {
	undo Undoable
	redo Undoable
	next *UndoStack
	id   int
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
	cursorPosition gridKey
	gridNote       GridNote
}

func (ugn UndoGridNote) ApplyUndo(m *model) {
	m.EnsureOverlayWithKey(ugn.overlayKey)
	m.cursorPos = ugn.cursorPosition
	overlay := m.definition.overlays[ugn.overlayKey]
	overlay[ugn.gridNote.gridKey] = ugn.gridNote.note
}

type UndoLineGridNotes struct {
	overlayKey
	cursorPosition gridKey
	line           uint8
	gridNotes      []GridNote
}

func (ulgn UndoLineGridNotes) ApplyUndo(m *model) {
	m.EnsureOverlayWithKey(ulgn.overlayKey)
	m.cursorPos = ulgn.cursorPosition
	overlay := m.definition.overlays[ulgn.overlayKey]
	for i := range m.definition.beats {
		delete(overlay, gridKey{ulgn.line, i})
	}
	for _, gridNote := range ulgn.gridNotes {
		overlay[gridNote.gridKey] = gridNote.note
	}
}

type UndoBounds struct {
	overlayKey
	cursorPosition gridKey
	bounds         Bounds
	gridNotes      []GridNote
}

func (uvs UndoBounds) ApplyUndo(m *model) {
	m.EnsureOverlayWithKey(uvs.overlayKey)
	m.cursorPos = uvs.cursorPosition
	overlay := m.definition.overlays[uvs.overlayKey]
	for _, k := range uvs.bounds.GridKeys() {
		delete(overlay, k)
	}
	for _, gridNote := range uvs.gridNotes {
		overlay[gridNote.gridKey] = gridNote.note
	}
}

type UndoGridNotes struct {
	overlayKey
	gridNotes []GridNote
}

func (ugn UndoGridNotes) ApplyUndo(m *model) {
	m.EnsureOverlayWithKey(ugn.overlayKey)
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
	m.cursorPos = utn.location
	delete(overlay, utn.location)
}

type UndoLineToNothing struct {
	overlayKey     overlayKey
	cursorPosition gridKey
	line           uint8
}

func (ultn UndoLineToNothing) ApplyUndo(m *model) {
	overlay := m.definition.overlays[ultn.overlayKey]
	m.cursorPos = ultn.cursorPosition
	for i := range m.definition.beats {
		delete(overlay, gridKey{ultn.line, i})
	}
}

type UndoNewOverlay struct {
	overlayKey     overlayKey
	cursorPosition gridKey
}

func (uno UndoNewOverlay) ApplyUndo(m *model) {
	delete(m.definition.overlays, uno.overlayKey)
	m.cursorPos = uno.cursorPosition
}

func (m *model) PushUndoables(undo Undoable, redo Undoable) {
	if m.undoStack == NIL_STACK {
		m.undoStack = UndoStack{
			undo: undo,
			redo: redo,
			next: nil,
			id:   rand.Int(),
		}
	} else {
		pusheddown := m.undoStack
		lastin := UndoStack{
			undo: undo,
			redo: redo,
			next: &pusheddown,
			id:   rand.Int(),
		}
		m.undoStack = lastin
	}
}

func (m *model) PushUndo(undo UndoStack) {
	if m.undoStack == NIL_STACK {
		undo.next = nil
		m.undoStack = undo
	} else {
		pusheddown := m.undoStack
		undo.next = &pusheddown
		m.undoStack = undo
	}
}

func (m *model) PushRedo(redo UndoStack) {
	if m.redoStack == NIL_STACK {
		redo.next = nil
		m.redoStack = redo
	} else {
		pusheddown := m.redoStack
		redo.next = &pusheddown
		m.redoStack = redo
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
	accents      patternAccents
}

type patternAccents struct {
	Diff   uint8
	Data   []Accent
	Start  uint8
	Target accentTarget
}

type accentTarget uint8

const (
	ACCENT_TARGET_NOTE accentTarget = iota
	ACCENT_TARGET_VELOCITY
)

func (pa *patternAccents) ReCalc() {
	accents := make([]Accent, 9)
	for i, a := range pa.Data[1:] {
		a.Value = pa.Start - pa.Diff*uint8(i)
		accents[i+1] = a
	}
	pa.Data = accents
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

func PlayBeat(accents patternAccents, beatInterval time.Duration, lines []lineDefinition, pattern overlay, currentBeat []linestate, sendFn SendFunc) []tea.Cmd {
	messages := make([]noteMessage, 0, len(lines))
	ratchetNotes := make([]lineNote, 0, len(lines))

	for i, line := range lines {
		if currentBeat[i].groupPlayState != PLAY_STATE_MUTE && currentBeat[i].groupPlayState != PLAY_STATE_MUTED_BY_SOLO {
			currentGridKey := gridKey{uint8(i), currentBeat[i].currentBeat}
			note, hasNote := pattern[currentGridKey]
			if hasNote && note.Ratchets.Length > 0 {
				ratchetNotes = append(ratchetNotes, lineNote{note, line})
			} else if hasNote && note != zeronote {
				messages = append(messages, line.Message(note, accents.Data[note.AccentIndex].Value, accents.Target))
			}
		}
	}

	playCmds := make([]tea.Cmd, 0, len(lines))

	playNotes := func() tea.Msg {
		Play(messages, sendFn)
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

func Play(messages []noteMessage, sendFn SendFunc) {
	for _, message := range messages {
		PlayMessage(message, sendFn)
	}
}

func PlayMessage(message noteMessage, sendFn SendFunc) {
	onMessage := midi.NoteOn(message.channel, message.note, message.velocity)
	offMessage := midi.NoteOff(message.channel, message.note)
	err := sendFn(onMessage)
	if err != nil {
		panic("note on failed")
	}
	err = sendFn(offMessage)
	if err != nil {
		panic("note off failed")
	}
}

func (m *model) EnsureOverlay() {
	m.EnsureOverlayWithKey(m.overlayKey)
}

func (m *model) EnsureOverlayWithKey(key overlayKey) {
	if len(m.definition.overlays[key]) == 0 {
		m.definition.overlays[key] = make(overlay)
		if key == overlaykey.ROOT_OVERLAY {
			m.definition.metaOverlays[key] = metaOverlay{PressUp: true, PressDown: false}
		} else {
			m.definition.metaOverlays[key] = metaOverlay{PressUp: false, PressDown: false}
		}
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

func absdiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}

func (m model) VisualSelectedGridKeys() []gridKey {
	if m.visualMode {
		return InitBounds(m.visualAnchorCursor, m.cursorPos).GridKeys()
	} else {
		return []gridKey{m.cursorPos}
	}
}

func (m *model) AddTrigger() {
	keys := m.VisualSelectedGridKeys()
	for _, k := range keys {
		m.CurrentNotable().SetNote(k, note{5, InitRatchet(), ACTION_NOTHING})
	}
}

func (m *model) AddAction(act action) {
	keys := m.VisualSelectedGridKeys()
	for _, k := range keys {
		m.CurrentNotable().SetNote(k, note{0, InitRatchet(), act})
	}
}

func (m *model) RemoveTrigger() {
	keys := m.VisualSelectedGridKeys()
	for _, k := range keys {
		m.CurrentNotable().SetNote(k, zeronote)
	}
}

func (m *model) OverlayRemoveTrigger() {
	keys := m.VisualSelectedGridKeys()
	for _, k := range keys {
		delete(m.definition.overlays[m.overlayKey], k)
	}
}

func (m *model) IncreaseRatchet() {
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())
	bounds := m.YankBounds()

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			currentRatchet := currentNote.Ratchets.Length
			if currentNote.AccentIndex > 0 && currentNote.Action == ACTION_NOTHING && currentRatchet+1 < uint8(len(ratchets)) {
				currentNote.Ratchets.Length = currentRatchet + 1
				currentNote.Ratchets.Hits[currentNote.Ratchets.Length] = true
				m.CurrentNotable().SetNote(key, currentNote)
			}
		}
	}
}

func (m *model) DecreaseRatchet() {
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())
	bounds := m.YankBounds()

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			currentRatchet := currentNote.Ratchets.Length
			if currentNote.AccentIndex > 0 && currentNote.Action == ACTION_NOTHING && currentRatchet > 0 {
				currentNote.Ratchets.Length = currentRatchet - 1
				m.CurrentNotable().SetNote(key, currentNote)
			}
		}
	}
}

func (m *model) EnsureRatchetCursorVisisble() {
	currentNote := m.CurrentNote()
	if m.ratchetCursor > currentNote.Ratchets.Length {
		m.ratchetCursor = m.ratchetCursor - 1
	}
}

func (m *model) IncreaseSpan() {
	currentNote := m.CurrentNote()
	if currentNote != zeronote && currentNote.Action == ACTION_NOTHING {
		span := currentNote.Ratchets.Span
		if span < 8 {
			currentNote.Ratchets.Span = span + 1
		}
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) DecreaseSpan() {
	currentNote := m.CurrentNote()
	if currentNote != zeronote && currentNote.Action == ACTION_NOTHING {
		span := currentNote.Ratchets.Span
		if span > 0 {
			currentNote.Ratchets.Span = span - 1
		}
		m.CurrentNotable().SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) IncreaseAccent() {
	m.definition.accents.Diff = m.definition.accents.Diff + 1
	m.definition.accents.ReCalc()
}

func (m *model) DecreaseAccent() {
	m.definition.accents.Diff = m.definition.accents.Diff - 1
	m.definition.accents.ReCalc()
}

func (m *model) DecreaseAccentTarget() {
	m.definition.accents.Target = (m.definition.accents.Target + 1) % 2
}

func (m *model) IncreaseAccentStart() {
	m.definition.accents.Start = m.definition.accents.Start + 1
	m.definition.accents.ReCalc()
}

func (m *model) DecreaseAccentStart() {
	m.definition.accents.Start = m.definition.accents.Start - 1
	m.definition.accents.ReCalc()
}

func (m *model) ToggleRatchetMute() {
	currentNote := m.CurrentNote()
	currentNote.Ratchets.Hits[m.ratchetCursor] = !currentNote.Ratchets.Hits[m.ratchetCursor]
	m.CurrentNotable().SetNote(m.cursorPos, currentNote)
}

func InitLines(n uint8) []lineDefinition {
	var lines = make([]lineDefinition, n)
	for i := range n {
		lines[i] = lineDefinition{
			Channel: 10,
			Note:    C1 + i,
		}
	}
	return lines
}

func InitLineStates(lines int, previousPlayState []linestate) []linestate {
	linestates := make([]linestate, 0, lines)

	for i := range lines {

		var previousGroupPlayState = PLAY_STATE_PLAY
		if len(previousPlayState) > int(i) {
			previousState := previousPlayState[i]
			previousGroupPlayState = previousState.groupPlayState
		}

		linestates = append(linestates, InitLineState(previousGroupPlayState))
	}
	return linestates
}

func InitLineState(previousGroupPlayState groupPlayState) linestate {
	return linestate{0, 1, 1, 0, 0, 0, previousGroupPlayState}
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
		accents:      patternAccents{Diff: 15, Data: accents, Start: 120, Target: ACCENT_TARGET_VELOCITY},
	}
}

func InitModel(midiConnection MidiConnection) model {
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
		midiConnection:      midiConnection,
		logFile:             logFile,
		cursorPos:           gridKey{0, 0},
		overlayKey:          overlaykey.ROOT_OVERLAY,
		definition:          definition,
		playState:           InitLineStates(len(definition.lines), []linestate{}),
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

func RunProgram(midiConnection MidiConnection) *tea.Program {
	p := tea.NewProgram(InitModel(midiConnection), tea.WithAltScreen())
	return p
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return tea.FocusMsg{} }
}

func Is(msg tea.KeyMsg, k key.Binding) bool {
	return key.Matches(msg, k)
}

func GetMinimumKeyCycle(key overlayKey) int {
	for i := 1; i < 100; i++ {
		if key.DoesMatch(i) {
			return i
		}
	}
	return 100
}

func (d Definition) GetMatchingOverlays(keyCycles int, keys []overlayKey) []overlayKey {
	var matchedKeys = make([]overlayKey, 0, 5)

	slices.SortFunc(keys, overlaykey.Sort)
	var pressNext = false

	for _, key := range keys {
		matches := key.DoesMatch(keyCycles)
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := transitiveKeys

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case Is(msg, keys.Quit):
			m.logFile.Close()
			return m, tea.Quit
		case Is(msg, keys.CursorDown):
			if slices.Contains([]Selection{SELECT_NOTHING, SELECT_SETUP_CHANNEL, SELECT_SETUP_NOTE}, m.selectionIndicator) {
				if m.cursorPos.line < uint8(len(m.definition.lines)-1) {
					m.cursorPos.line++
				}
			}
		case Is(msg, keys.CursorUp):
			if slices.Contains([]Selection{SELECT_NOTHING, SELECT_SETUP_CHANNEL, SELECT_SETUP_NOTE}, m.selectionIndicator) {
				if m.cursorPos.line > 0 {
					m.cursorPos.line--
				}
			}
		case Is(msg, keys.CursorLeft):
			if m.selectionIndicator == SELECT_RATCHETS {
				if m.ratchetCursor > 0 {
					m.ratchetCursor--
				}
			} else if m.selectionIndicator > 0 {
				// Do Nothing
			} else {
				if m.cursorPos.beat > 0 {
					m.cursorPos.beat--
				}
			}
		case Is(msg, keys.CursorRight):
			if m.selectionIndicator == SELECT_RATCHETS {
				currentNote := m.CurrentNote()
				if m.ratchetCursor < currentNote.Ratchets.Length {
					m.ratchetCursor++
				}
			} else if m.selectionIndicator > 0 {
				// Do Nothing
			} else {
				if m.cursorPos.beat < m.definition.beats-1 {
					m.cursorPos.beat++
				}
			}
		case Is(msg, keys.CursorLineStart):
			m.cursorPos.beat = 0
		case Is(msg, keys.CursorLineEnd):
			m.cursorPos.beat = m.definition.beats - 1
		case Is(msg, keys.Escape):
			m.selectionIndicator = 0
			m.accentMode = false
			m.visualMode = false
		case Is(msg, keys.PlayStop):
			if !m.playing && !m.midiConnection.IsOpen() {
				err := m.midiConnection.ConnectAndOpen()
				if err != nil {
					panic("No Open Connection")
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
				m.playState = InitLineStates(len(m.definition.lines), m.playState)
				m.advanceKeyCycle()
				m.trackTime = time.Duration(0)
				sendFn := m.midiConnection.AcquireSendFunc()
				beatInterval := m.BeatInterval()

				cmds := make([]tea.Cmd, 0, 10)
				cmds = append(cmds, PlayBeat(m.definition.accents, beatInterval, m.definition.lines, m.definition.CombinedPattern(m.playingMatchedOverlays), m.playState, sendFn)...)
				cmds = append(cmds, BeatTick(beatInterval))
				return m, tea.Batch(cmds...)
			} else {
				m.keyCycles = 0
				m.playingMatchedOverlays = []overlayKey{}
			}
		case Is(msg, keys.OverlayInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_OVERLAY_NUM, SELECT_OVERLAY_DENOM}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
		case Is(msg, keys.TempoInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_TEMPO, SELECT_TEMPO_SUBDIVISION}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
		case Is(msg, keys.SetupInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_SETUP_CHANNEL, SELECT_SETUP_NOTE}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
		case Is(msg, keys.AccentInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_ACCENT_DIFF, SELECT_ACCENT_TARGET, SELECT_ACCENT_START}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
		case Is(msg, keys.ToggleRatchetMode):
			currentNote := m.CurrentNote()
			if currentNote.AccentIndex > 0 {
				states := []Selection{SELECT_NOTHING, SELECT_RATCHETS, SELECT_RATCHET_SPAN}
				m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
				m.ratchetCursor = 0
			}
		case Is(msg, keys.Increase):
			switch m.selectionIndicator {
			case SELECT_TEMPO:
				if m.definition.tempo < 300 {
					m.definition.tempo++
				}
			case SELECT_TEMPO_SUBDIVISION:
				if m.definition.subdivisions < 8 {
					m.definition.subdivisions++
				}
			case SELECT_OVERLAY_NUM:
				m.overlayKey.IncrementShift()
			case SELECT_OVERLAY_DENOM:
				m.overlayKey.IncrementInterval()
			case SELECT_SETUP_CHANNEL:
				m.definition.lines[m.cursorPos.line].IncrementChannel()
			case SELECT_SETUP_NOTE:
				m.definition.lines[m.cursorPos.line].IncrementNote()
			case SELECT_RATCHET_SPAN:
				m.IncreaseSpan()
			case SELECT_ACCENT_DIFF:
				m.IncreaseAccent()
			case SELECT_ACCENT_TARGET:
				// Only two options right now, so increase and decrease would do the
				// same thing
				m.DecreaseAccentTarget()
			case SELECT_ACCENT_START:
				m.IncreaseAccentStart()
			}
		case Is(msg, keys.Decrease):
			switch m.selectionIndicator {
			case SELECT_TEMPO:
				if m.definition.tempo > 30 {
					m.definition.tempo--
				}
			case SELECT_TEMPO_SUBDIVISION:
				if m.definition.subdivisions > 1 {
					m.definition.subdivisions--
				}
			case SELECT_OVERLAY_NUM:
				m.overlayKey.DecrementShift()
			case SELECT_OVERLAY_DENOM:
				m.overlayKey.DecrementInterval()
			case SELECT_SETUP_CHANNEL:
				m.definition.lines[m.cursorPos.line].DecrementChannel()
			case SELECT_SETUP_NOTE:
				m.definition.lines[m.cursorPos.line].DecrementNote()
			case SELECT_RATCHET_SPAN:
				m.DecreaseSpan()
			case SELECT_ACCENT_DIFF:
				m.DecreaseAccent()
			case SELECT_ACCENT_TARGET:
				m.DecreaseAccentTarget()
			case SELECT_ACCENT_START:
				m.DecreaseAccentStart()
			}
		case Is(msg, keys.ToggleAccentMode):
			m.accentMode = !m.accentMode
		case Is(msg, keys.PrevOverlay):
			m.NextOverlay(-1)
		case Is(msg, keys.NextOverlay):
			m.NextOverlay(+1)
		case Is(msg, keys.Save):
			m.Save()
			m.needsWrite = m.undoStack.id
		case Is(msg, keys.Undo):
			undoStack := m.Undo()
			if undoStack != NIL_STACK {
				m.PushRedo(undoStack)
			}
		case Is(msg, keys.Redo):
			undoStack := m.Redo()
			if undoStack != NIL_STACK {
				m.PushUndo(undoStack)
			}
		case Is(msg, keys.New):
			m.cursorPos = gridKey{0, 0}
			m.overlayKey = overlaykey.ROOT_OVERLAY
			m.selectionIndicator = SELECT_NOTHING
			m.definition = InitDefinition()
		case Is(msg, keys.ToggleVisualMode):
			m.visualAnchorCursor = m.cursorPos
			m.visualMode = !m.visualMode
		case Is(msg, keys.NewLine):
			if len(m.definition.lines) < 16 {
				lastline := m.definition.lines[len(m.definition.lines)-1]
				m.definition.lines = append(m.definition.lines, lineDefinition{
					Channel: lastline.Channel,
					Note:    lastline.Note + 1,
				})
				if m.playing {
					m.playState = append(m.playState, InitLineState(m.GroupPlayStateForNewLine()))
				}
			}
		case Is(msg, keys.Yank):
			m.yankBuffer = m.Yank()
			m.cursorPos = m.YankBounds().TopLeft()
			m.visualMode = false
		case Is(msg, keys.Mute):
			if m.IsRatchetSelector() {
				m.ToggleRatchetMute()
			} else {
				m.Mute()
			}
		case Is(msg, keys.Solo):
			m.Solo()
		default:
			m = m.UpdateDefinition(msg)
		}
	case beatMsg:
		if m.playing {
			m.advanceCurrentBeat()
			m.advanceKeyCycle()
			m.totalBeats++
			sendFn := m.midiConnection.AcquireSendFunc()
			beatInterval := m.BeatInterval()
			cmds := make([]tea.Cmd, 0, 10)
			cmds = append(cmds, PlayBeat(m.definition.accents, beatInterval, m.definition.lines, m.definition.CombinedPattern(m.playingMatchedOverlays), m.playState, sendFn)...)
			cmds = append(cmds, BeatTick(beatInterval))
			return m, tea.Batch(
				cmds...,
			)
		}
	case ratchetMsg:
		if m.playing && msg.iterations < (msg.Ratchets.Length+1) {
			var playCmd tea.Cmd
			var ratchetTickCmd tea.Cmd

			sendFn := m.midiConnection.AcquireSendFunc()

			if msg.Ratchets.Hits[msg.iterations] {
				note := msg.note
				message := msg.line.Message(msg.note, m.definition.accents.Data[note.AccentIndex].Value, m.definition.accents.Target)
				playCmd = func() tea.Msg {
					PlayMessage(message, sendFn)
					return nil
				}
			}
			if msg.iterations+1 < (msg.Ratchets.Length + 1) {
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

func AdvanceSelectionState(states []Selection, currentSelection Selection) Selection {
	index := slices.Index(states, currentSelection)
	var resultSelection Selection
	if index < 0 {
		indexNothing := slices.Index(states, SELECT_NOTHING)
		resultSelection = states[uint8(indexNothing+1)%uint8(len(states))]
	} else {
		resultSelection = states[uint8(index+1)%uint8(len(states))]
	}
	return resultSelection
}

func (m model) UpdateDefinitionKeys(msg tea.KeyMsg) model {
	keys := definitionKeys
	switch {
	case Is(msg, keys.TriggerAdd):
		m.AddTrigger()
	case Is(msg, keys.TriggerRemove):
		m.yankBuffer = m.Yank()
		m.RemoveTrigger()
		m.visualMode = false
	case Is(msg, keys.AccentIncrease):
		m.AccentModify(1)
	case Is(msg, keys.AccentDecrease):
		m.AccentModify(-1)
	case Is(msg, keys.OverlayTriggerRemove):
		m.OverlayRemoveTrigger()
	case Is(msg, keys.ClearLine):
		m.ClearOverlayLine()
	case Is(msg, keys.RatchetIncrease):
		m.IncreaseRatchet()
	case Is(msg, keys.RatchetDecrease):
		m.DecreaseRatchet()
		m.EnsureRatchetCursorVisisble()
	case Is(msg, keys.ActionAddLineReset):
		m.AddAction(ACTION_LINE_RESET)
	case Is(msg, keys.ActionAddLineReverse):
		m.AddAction(ACTION_LINE_REVERSE)
	case Is(msg, keys.ActionAddSkipBeat):
		m.AddAction(ACTION_LINE_SKIP_BEAT)
	case Is(msg, keys.ActionAddReset):
		m.AddAction(ACTION_RESET)
	case Is(msg, keys.ActionAddLineBounce):
		m.AddAction(ACTION_LINE_BOUNCE)
	case Is(msg, keys.ActionAddLineDelay):
		m.AddAction(ACTION_LINE_DELAY)
	case Is(msg, keys.SelectKeyLine):
		undoable := UndoKeyline{m.definition.keyline}
		m.definition.keyline = m.cursorPos.line
		redoable := UndoKeyline{m.definition.keyline}
		m.PushUndoables(undoable, redoable)
	case Is(msg, keys.PressDownOverlay):
		m.ToggleOverlayStackOptions(m.overlayKey)
	case Is(msg, keys.ClearSeq):
		m.ClearOverlay()
	case Is(msg, keys.RotateRight):
		m.RotateRight()
	case Is(msg, keys.RotateLeft):
		m.RotateLeft()
	case Is(msg, keys.Paste):
		m.Paste()
	}
	if msg.String() >= "1" && msg.String() <= "9" {
		beatInterval, _ := strconv.ParseInt(msg.String(), 0, 8)
		if m.accentMode {
			m.incrementAccent(uint8(beatInterval), -1)
		} else {
			m.fill(uint8(beatInterval))
		}
	}
	if IsShiftSymbol(msg.String()) {
		beatInterval := convertSymbolToInt(msg.String())
		if m.accentMode {
			m.incrementAccent(uint8(beatInterval), 1)
		} else {
			m.fill(uint8(beatInterval))
		}
	}
	return m
}

func IsShiftSymbol(symbol string) bool {
	return slices.Contains([]string{"!", "@", "#", "$", "%", "^", "&", "*", "("}, symbol)
}

func convertSymbolToInt(symbol string) int64 {
	switch symbol[0] {
	case '!':
		return 1
	case '@':
		return 2
	case '#':
		return 3
	case '$':
		return 4
	case '%':
		return 5
	case '^':
		return 6
	case '&':
		return 7
	case '*':
		return 8
	case '(':
		return 9
	}
	return 0
}

func (m model) UpdateDefinition(msg tea.KeyMsg) model {
	keys := definitionKeys
	if m.visualMode && (keys.IsLineWiseKey(msg) || keys.IsNoteWiseKey(msg) || Is(msg, keys.Paste)) {
		undoable := m.UndoableBounds(m.visualAnchorCursor, m.cursorPos)
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableBounds(m.visualAnchorCursor, m.cursorPos)
		m.PushUndoables(undoable, redoable)
	} else if keys.IsNoteWiseKey(msg) {
		undoable := m.UndoableNote()
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableNote()
		m.PushUndoables(undoable, redoable)
	} else if keys.IsLineWiseKey(msg) {
		undoable := m.UndoableLine()
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableLine()
		m.PushUndoables(undoable, redoable)
	} else if keys.IsOverlayWiseKey(msg) {
		undoable := m.UndoableOverlay()
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableOverlay()
		m.PushUndoables(undoable, redoable)
	} else if Is(msg, keys.Paste) {
		undoable := m.UndoableBounds(m.cursorPos, m.yankBuffer.bounds.BottomRightFrom(m.cursorPos))
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
		redoable := m.UndoableBounds(m.cursorPos, m.yankBuffer.bounds.BottomRightFrom(m.cursorPos))
		m.PushUndoables(undoable, redoable)
	} else {
		m.EnsureOverlay()
		m = m.UpdateDefinitionKeys(msg)
	}
	m.ResetRedo()
	return m
}

func (m model) UndoableNote() Undoable {
	overlay, hasOverlay := m.definition.overlays[m.overlayKey]
	if !hasOverlay {
		return UndoNewOverlay{m.overlayKey, m.cursorPos}
	}
	currentNote, hasNote := overlay[m.cursorPos]
	if hasNote {
		return UndoGridNote{m.overlayKey, m.cursorPos, GridNote{m.cursorPos, currentNote}}
	} else {
		return UndoToNothing{m.overlayKey, m.cursorPos}
	}
}

func (m model) UndoableBounds(pointA, pointB gridKey) Undoable {
	overlay, hasOverlay := m.definition.overlays[m.overlayKey]
	if !hasOverlay {
		return UndoNewOverlay{m.overlayKey, m.cursorPos}
	}
	bounds := InitBounds(pointA, pointB)
	gridKeys := bounds.GridKeys()
	gridNotes := make([]GridNote, 0, len(gridKeys))
	for _, k := range gridKeys {
		currentNote, hasNote := overlay[k]
		if hasNote {
			gridNotes = append(gridNotes, GridNote{k, currentNote})
		}
	}
	return UndoBounds{m.overlayKey, m.cursorPos, bounds, gridNotes}
}

func (m model) UndoableLine() Undoable {
	overlay, hasOverlay := m.definition.overlays[m.overlayKey]
	if !hasOverlay {
		return UndoNewOverlay{m.overlayKey, m.cursorPos}
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
		return UndoLineToNothing{m.overlayKey, m.cursorPos, m.cursorPos.line}
	}
	return UndoLineGridNotes{m.overlayKey, m.cursorPos, m.cursorPos.line, notesToUndo}
}

func (m model) UndoableOverlay() Undoable {
	_, hasOverlay := m.definition.overlays[m.overlayKey]
	if !hasOverlay {
		return UndoNewOverlay{m.overlayKey, m.cursorPos}
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
	} else if !meta.PressDown && !meta.PressUp {
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
	index := slices.Index(keys, overlaykey.ROOT_OVERLAY)
	if index >= 0 {
		return append(keys[:index], keys[index+1:]...)
	}
	return keys
}

func (m *model) NextOverlay(direction int) {
	keys := m.OverlayKeys()
	slices.SortFunc(keys, overlaykey.Sort)
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

func (m *model) RotateRight() {
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())

	start, end := m.PatternActionBoundaries()
	lastNote := combinedOverlay[gridKey{m.cursorPos.line, end}]
	previousNote := zeronote
	for i := uint8(start); i <= end; i++ {
		gridKey := gridKey{m.cursorPos.line, i}
		currentNote := combinedOverlay[gridKey]

		m.CurrentNotable().SetNote(gridKey, previousNote)
		previousNote = currentNote
	}

	m.CurrentNotable().SetNote(gridKey{m.cursorPos.line, start}, lastNote)
}

func (m *model) RotateLeft() {
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())

	start, end := m.PatternActionBoundaries()
	firstNote := combinedOverlay[gridKey{m.cursorPos.line, start}]
	previousNote := zeronote
	for i := int8(end); i >= int8(start); i-- {
		gridKey := gridKey{m.cursorPos.line, uint8(i)}
		currentNote := combinedOverlay[gridKey]

		m.CurrentNotable().SetNote(gridKey, previousNote)
		previousNote = currentNote
	}

	m.CurrentNotable().SetNote(gridKey{m.cursorPos.line, end}, firstNote)
}

func (m model) GroupPlayStateForNewLine() groupPlayState {
	for _, state := range m.playState {
		if state.groupPlayState == PLAY_STATE_SOLO {
			return PLAY_STATE_MUTED_BY_SOLO
		}
	}
	return PLAY_STATE_PLAY
}

func (m *model) Mute() {
	var hasOtherSolo = m.hasOtherSolo(m.cursorPos.line)
	switch m.playState[m.cursorPos.line].groupPlayState {
	case PLAY_STATE_PLAY:
		m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_MUTE
	case PLAY_STATE_MUTED_BY_SOLO:
		m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_MUTE
	case PLAY_STATE_MUTE:
		if hasOtherSolo {
			m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_MUTED_BY_SOLO
		} else {
			m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_PLAY
		}
	case PLAY_STATE_SOLO:
		m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_MUTE
	}
}

func (m *model) Solo() {
	var hasOtherSolo = m.hasOtherSolo(m.cursorPos.line)
	for i, state := range m.playState {
		if uint8(i) != m.cursorPos.line {
			switch state.groupPlayState {
			case PLAY_STATE_PLAY:
				m.playState[i].groupPlayState = PLAY_STATE_MUTED_BY_SOLO
			case PLAY_STATE_MUTED_BY_SOLO:
				if hasOtherSolo {
					m.playState[i].groupPlayState = PLAY_STATE_MUTED_BY_SOLO
				} else {
					m.playState[i].groupPlayState = PLAY_STATE_PLAY
				}
			case PLAY_STATE_MUTE:
				m.playState[i].groupPlayState = PLAY_STATE_MUTE
			case PLAY_STATE_SOLO:
				m.playState[i].groupPlayState = PLAY_STATE_SOLO
			}
		}
	}
	switch m.playState[m.cursorPos.line].groupPlayState {
	case PLAY_STATE_PLAY:
		m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_SOLO
	case PLAY_STATE_MUTED_BY_SOLO:
		m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_SOLO
	case PLAY_STATE_MUTE:
		m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_SOLO
	case PLAY_STATE_SOLO:
		if hasOtherSolo {
			m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_MUTED_BY_SOLO
		} else {
			m.playState[m.cursorPos.line].groupPlayState = PLAY_STATE_PLAY
		}
	}
}

func (m model) hasOtherSolo(than uint8) bool {
	for i, state := range m.playState {
		if i == int(than) {
			continue
		}
		if state.groupPlayState == PLAY_STATE_SOLO {
			return true
		}
	}
	return false
}

type Buffer struct {
	bounds    Bounds
	gridNotes []GridNote
}

func (m model) Yank() Buffer {
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())
	bounds := m.YankBounds()
	capturedGridNotes := make([]GridNote, 0, len(combinedOverlay))

	for key, note := range combinedOverlay {
		if bounds.InBounds(key) {
			normalizedGridKey := gridKey{key.line - bounds.top, key.beat - bounds.left}
			capturedGridNotes = append(capturedGridNotes, GridNote{normalizedGridKey, note})
		}
	}

	return Buffer{
		bounds:    bounds.Normalized(),
		gridNotes: capturedGridNotes,
	}
}

func (m *model) Paste() {
	bounds := m.PasteBounds()

	var keyModifier gridKey
	if m.visualMode {
		keyModifier = bounds.TopLeft()
	} else {
		keyModifier = m.cursorPos
	}

	for _, gridNote := range m.yankBuffer.gridNotes {
		key := gridNote.gridKey
		newKey := gridKey{key.line + keyModifier.line, key.beat + keyModifier.beat}
		if bounds.InBounds(newKey) {
			m.CurrentNotable().SetNote(newKey, gridNote.note)
		}
	}
}

func (m *model) advanceCurrentBeat() {
	combinedPattern := m.definition.CombinedPattern(m.playingMatchedOverlays)
	for i := range m.playState {
		doContinue := m.advancePlayState(combinedPattern, i)
		if !doContinue {
			break
		}
	}
}

func (m *model) advancePlayState(combinedPattern overlay, lineIndex int) bool {
	currentState := m.playState[lineIndex]
	advancedBeat := int8(currentState.currentBeat) + currentState.direction

	if advancedBeat < 0 || advancedBeat >= int8(m.definition.beats) {
		// reset locations should be 1 time use.  Reset back to 0.
		if m.playState[lineIndex].resetLocation != 0 && combinedPattern[gridKey{uint8(lineIndex), currentState.resetActionLocation}].Action == currentState.resetAction {
			m.playState[lineIndex].currentBeat = currentState.resetLocation
			advancedBeat = int8(currentState.resetLocation)
		} else {
			m.playState[lineIndex].currentBeat = 0
			advancedBeat = int8(0)
		}
		m.playState[lineIndex].direction = currentState.resetDirection
		m.playState[lineIndex].resetLocation = 0
	} else {
		m.playState[lineIndex].currentBeat = uint8(advancedBeat)
	}

	switch combinedPattern[gridKey{uint8(lineIndex), uint8(advancedBeat)}].Action {
	case ACTION_LINE_RESET:
		m.playState[lineIndex].currentBeat = 0
	case ACTION_LINE_REVERSE:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-2, 0))
		m.playState[lineIndex].direction = -1
		m.playState[lineIndex].resetLocation = uint8(max(advancedBeat-1, 0))
		m.playState[lineIndex].resetActionLocation = uint8(advancedBeat)
		m.playState[lineIndex].resetAction = ACTION_LINE_REVERSE
	case ACTION_LINE_BOUNCE:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-1, 0))
		m.playState[lineIndex].direction = -1
	case ACTION_LINE_SKIP_BEAT:
		m.advancePlayState(combinedPattern, lineIndex)
	case ACTION_LINE_DELAY:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-1, 0))
	case ACTION_RESET:
		for i := range m.playState {
			m.playState[i].currentBeat = 0
			m.playState[i].direction = 1
			m.playState[i].resetLocation = 0
			m.playState[i].resetDirection = 1
		}
		return false
	}

	return true
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
	if !slices.Contains(keys, m.overlayKey) {
		keys = append(keys, m.overlayKey)
	}
	slices.SortFunc(keys, overlaykey.Sort)
	slices.Reverse(keys)
	indexOfCurrent := slices.Index(keys, m.overlayKey)
	if indexOfCurrent >= 0 {
		return keys[:indexOfCurrent]
	} else {
		return []overlayKey{}
	}
}

func (m *model) fill(every uint8) {
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())

	start, end := m.PatternActionBoundaries()

	for i := uint8(start); i <= end; i += every {
		gridKey := gridKey{m.cursorPos.line, i}
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.RemoveNote(gridKey)
		} else {
			m.CurrentNotable().SetNote(gridKey, note{5, InitRatchet(), 0})
		}
	}
}

func (m *model) incrementAccent(every uint8, modifier int8) {
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())

	start, end := m.PatternActionBoundaries()

	for i := uint8(start); i <= end; i += every {
		gridKey := gridKey{m.cursorPos.line, i}
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.CurrentNotable().SetNote(gridKey, currentNote.IncrementAccent(modifier))
		}
	}
}

func (m *model) AccentModify(modifier int8) {
	bounds := m.YankBounds()
	combinedOverlay := m.definition.CombinedPattern(m.EditKeys())

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			if currentNote != zeronote {
				m.CurrentNotable().SetNote(key, currentNote.IncrementAccent(modifier))
			}
		}
	}
}

func (m model) PatternActionBoundaries() (uint8, uint8) {
	if m.visualMode {
		if m.visualAnchorCursor.beat < m.cursorPos.beat {
			return m.visualAnchorCursor.beat, m.cursorPos.beat
		} else {
			return m.cursorPos.beat, m.visualAnchorCursor.beat
		}
	} else {
		return m.cursorPos.beat, m.definition.beats - 1
	}
}

func (m *model) RemoveNote(gridKey gridKey) {
	if m.overlayKey == overlaykey.ROOT_OVERLAY {
		delete(m.definition.overlays[overlaykey.ROOT_OVERLAY], gridKey)
	} else {
		m.CurrentNotable().SetNote(gridKey, zeronote)
	}
}

func (m model) EditKeys() []overlayKey {
	keysBelowCurrent := m.KeysBelowCurrent()
	keysBelowCurrent = append(keysBelowCurrent, m.overlayKey)
	matchedKeys := m.definition.GetMatchingOverlays(GetMinimumKeyCycle(m.overlayKey), keysBelowCurrent)
	slices.SortFunc(matchedKeys, overlaykey.Sort)
	return matchedKeys
}

func (m *model) CurrentKeys() []overlayKey {
	if m.playing {
		return m.playingMatchedOverlays
	} else {
		return m.EditKeys()
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
	switch m.selectionIndicator {
	case SELECT_TEMPO:
		tempo = selectedColor.Render(strconv.Itoa(m.definition.tempo))
		division = numberColor.Render(strconv.Itoa(m.definition.subdivisions))
	case SELECT_TEMPO_SUBDIVISION:
		tempo = numberColor.Render(strconv.Itoa(m.definition.tempo))
		division = selectedColor.Render(strconv.Itoa(m.definition.subdivisions))
	default:
		tempo = numberColor.Render(strconv.Itoa(m.definition.tempo))
		division = numberColor.Render(strconv.Itoa(m.definition.subdivisions))
	}
	heart := heartColor.Render("♡")
	buf.WriteString("             \n")
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

func (m model) WriteView() string {
	if m.needsWrite != m.undoStack.id {
		return " [+]"
	} else {
		return "    "
	}
}

func (m model) IsAccentSelector() bool {
	states := []Selection{SELECT_ACCENT_DIFF, SELECT_ACCENT_TARGET, SELECT_ACCENT_START}
	return slices.Contains(states, m.selectionIndicator)
}

func (m model) IsRatchetSelector() bool {
	states := []Selection{SELECT_RATCHETS, SELECT_RATCHET_SPAN}
	return slices.Contains(states, m.selectionIndicator)
}

func (m model) View() string {
	var buf strings.Builder
	var sideView string

	if m.accentMode || m.IsAccentSelector() {
		sideView = m.AccentKeyView()
	} else if len(m.definition.overlays) == 0 ||
		m.selectionIndicator == SELECT_SETUP_NOTE ||
		m.selectionIndicator == SELECT_SETUP_CHANNEL {
		sideView = m.SetupView()
	} else {
		sideView = m.OverlaysView()
	}

	buf.WriteString(lipgloss.JoinHorizontal(0, m.TempoView(), "  ", m.ViewTriggerSeq(), "  ", sideView))
	return buf.String()
}

func (m model) AccentKeyView() string {
	var buf strings.Builder
	var accentDiffString string
	var accentDiff = m.definition.accents.Diff
	var accentStart = m.definition.accents.Start

	var accentTarget string
	switch m.definition.accents.Target {
	case ACCENT_TARGET_NOTE:
		accentTarget = "N"
	case ACCENT_TARGET_VELOCITY:
		accentTarget = "V"
	}

	if m.selectionIndicator == SELECT_ACCENT_DIFF {
		accentDiffString = selectedColor.Render(fmt.Sprintf("%2d", accentDiff))
	} else {
		accentDiffString = numberColor.Render(fmt.Sprintf("%2d", accentDiff))
	}

	var accentTargetString string
	if m.selectionIndicator == SELECT_ACCENT_TARGET {
		accentTargetString = selectedColor.Render(fmt.Sprintf(" %s", accentTarget))
	} else {
		accentTargetString = numberColor.Render(fmt.Sprintf(" %s", accentTarget))
	}

	buf.WriteString(fmt.Sprintf(" ACCENTS %s %s\n", accentDiffString, accentTargetString))
	buf.WriteString("———————————————\n")
	startAccent := m.definition.accents.Data[1]

	var accentStartString string
	if m.selectionIndicator == SELECT_ACCENT_START {
		accentStartString = selectedColor.Render(fmt.Sprintf("%2d", accentStart))
	} else {
		accentStartString = numberColor.Render(fmt.Sprintf("%2d", accentStart))
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(startAccent.Color))
	buf.WriteString(fmt.Sprintf("  %s  -  %s\n", style.Render(string(startAccent.Shape)), accentStartString))
	for _, accent := range m.definition.accents.Data[2:] {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(accent.Color))
		buf.WriteString(fmt.Sprintf("  %s  -  %d\n", style.Render(string(accent.Shape)), accent.Value))
	}
	return buf.String()
}

func (m model) SetupView() string {
	var buf strings.Builder
	buf.WriteString("    Setup\n")
	buf.WriteString("———————————————\n")
	for i, line := range m.definition.lines {

		buf.WriteString("CH ")
		if uint8(i) == m.cursorPos.line && m.selectionIndicator == SELECT_SETUP_CHANNEL {
			buf.WriteString(selectedColor.Render(fmt.Sprintf("%2d", line.Channel)))
		} else {
			buf.WriteString(numberColor.Render(fmt.Sprintf("%2d", line.Channel)))
		}

		buf.WriteString(" NOTE ")
		if uint8(i) == m.cursorPos.line && m.selectionIndicator == SELECT_SETUP_NOTE {
			buf.WriteString(selectedColor.Render(strconv.Itoa(int(line.Note))))
		} else {
			buf.WriteString(numberColor.Render(strconv.Itoa(int(line.Note))))
		}
		buf.WriteString(fmt.Sprintf(" %s%d\n", strings.ReplaceAll(midi.Note(line.Note).Name(), "b", "♭"), midi.Note(line.Note).Octave()-2))
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
	slices.SortFunc(keys, overlaykey.Sort)
	style := lipgloss.NewStyle().Background(seqOverlayColor)
	for _, k := range keys {
		var playingSpacer = "   "
		var playing = ""
		if m.playing && m.playingMatchedOverlays[0] == k {
			playing = lipgloss.NewStyle().Background(seqOverlayColor).Foreground(currentPlayingColor).Render(" \u25CF ")
			buf.WriteString(playing)
			playingSpacer = ""
		} else if m.playing && slices.Contains(m.playingMatchedOverlays, k) {
			playing = lipgloss.NewStyle().Background(seqOverlayColor).Foreground(activePlayingColor).Render(" \u25C9 ")
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
		overlayLine := fmt.Sprintf("%d/%d%2s%2s", k.Shift, k.Interval, stackModifier, editing)

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

var accentModeStyle = lipgloss.NewStyle().Background(accents[1].Color).Foreground(lipgloss.Color("#000000"))

func (m model) ViewTriggerSeq() string {
	var buf strings.Builder
	var mode string
	visualCombinedPattern := m.definition.VisualCombinedPattern(m.CurrentKeys())

	if m.accentMode {
		mode = " Accent Mode "
		buf.WriteString(fmt.Sprintf("    %s\n", accentModeStyle.Render(mode)))
	} else if m.selectionIndicator == SELECT_RATCHETS || m.selectionIndicator == SELECT_RATCHET_SPAN {
		buf.WriteString(m.RatchetModeView())
	} else if m.playing {
		buf.WriteString(fmt.Sprintf("    Seq - Playing - %d\n", m.keyCycles))
	} else {
		buf.WriteString(m.WriteView())
		buf.WriteString("Seq - A sequencer for your cli\n")
	}
	buf.WriteString("   ┌─────────────────────────────────\n")
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
		if i <= currentNote.Ratchets.Length {
			if m.ratchetCursor == i && m.selectionIndicator == SELECT_RATCHETS {
				backgroundColor = lipgloss.Color("#5cdffb")
			}
			if currentNote.Ratchets.Hits[i] {
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
	if m.selectionIndicator == SELECT_RATCHET_SPAN {
		buf.WriteString(fmt.Sprintf(" Span %s ", selectedColor.Render(strconv.Itoa(int(currentNote.Ratchets.GetSpan())))))
	} else {
		buf.WriteString(fmt.Sprintf(" Span %s ", numberColor.Render(strconv.Itoa(int(currentNote.Ratchets.GetSpan())))))
	}
	buf.WriteString("\n")

	return buf.String()
}

func (m model) ViewOverlay() string {
	var numerator, denominator string

	switch m.selectionIndicator {
	case SELECT_OVERLAY_NUM:
		numerator = selectedColor.Render(strconv.Itoa(int(m.overlayKey.Shift)))
		denominator = numberColor.Render(strconv.Itoa(int(m.overlayKey.Interval)))
	case SELECT_OVERLAY_DENOM:
		numerator = numberColor.Render(strconv.Itoa(int(m.overlayKey.Shift)))
		denominator = selectedColor.Render(strconv.Itoa(int(m.overlayKey.Interval)))
	default:
		numerator = numberColor.Render(strconv.Itoa(int(m.overlayKey.Shift)))
		denominator = numberColor.Render(strconv.Itoa(int(m.overlayKey.Interval)))
	}

	return fmt.Sprintf("%s/%s", numerator, denominator)
}

func (m model) CurrentOverlayView() string {
	var matchedKey overlayKey
	if len(m.playingMatchedOverlays) > 0 {
		matchedKey = m.playingMatchedOverlays[0]
	} else {
		matchedKey = overlaykey.ROOT_OVERLAY
	}
	return fmt.Sprintf("   Editing - %s     Playing - %d/%d", m.ViewOverlay(), matchedKey.Shift, matchedKey.Interval)
}

var altSeqColor = lipgloss.Color("#222222")
var seqColor = lipgloss.Color("#000000")
var seqCursorColor = lipgloss.Color("#444444")
var seqVisualColor = lipgloss.Color("#aaaaaa")
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
	indicator := "│"
	if len(m.playState) > int(lineNumber) && m.playState[lineNumber].groupPlayState == PLAY_STATE_MUTE {
		indicator = "M"
	}
	if len(m.playState) > int(lineNumber) && m.playState[lineNumber].groupPlayState == PLAY_STATE_SOLO {
		indicator = "S"
	}
	buf.WriteString(fmt.Sprintf("%2d%s%s", lineNumber, KeyLineIndicator(m.definition.keyline, lineNumber), indicator))

	for i := uint8(0); i < m.definition.beats; i++ {
		currentGridKey := gridKey{uint8(lineNumber), i}
		overlayNote, hasNote := visualCombinedPattern[currentGridKey]

		var backgroundSeqColor lipgloss.Color
		if m.playing && m.playState[lineNumber].currentBeat == i {
			backgroundSeqColor = seqCursorColor
		} else if m.visualMode && m.InVisualSelection(currentGridKey) {
			backgroundSeqColor = seqVisualColor
		} else if hasNote && overlayNote.origin != overlaykey.ROOT_OVERLAY {
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
		} else if m.visualMode && m.InVisualSelection(currentGridKey) {
			style = style.Foreground(lipgloss.Color("#000000"))
		} else {
			style = style.Foreground(foregroundColor)
		}
		buf.WriteString(style.Render(char))
	}

	buf.WriteString("\n")
	return buf.String()
}

func InitBounds(cursorA, cursorB gridKey) Bounds {
	return Bounds{
		top:    min(cursorA.line, cursorB.line),
		right:  max(cursorA.beat, cursorB.beat),
		bottom: max(cursorA.line, cursorB.line),
		left:   min(cursorA.beat, cursorB.beat),
	}
}

type Bounds struct {
	top    uint8
	right  uint8
	bottom uint8
	left   uint8
}

func (b Bounds) Area() int {
	return int(absdiff(b.top, b.bottom) * absdiff(b.left, b.right))
}

func (bounds Bounds) GridKeys() []gridKey {
	keys := make([]gridKey, 0, bounds.Area())
	for i := bounds.top; i <= bounds.bottom; i++ {
		for j := bounds.left; j <= bounds.right; j++ {
			keys = append(keys, gridKey{i, j})
		}
	}
	return keys
}

func (b Bounds) InBounds(key gridKey) bool {
	return key.line >= b.top &&
		key.line <= b.bottom &&
		key.beat >= b.left &&
		key.beat <= b.right
}

func (b Bounds) Normalized() Bounds {
	return Bounds{top: 0, right: b.right - b.left, bottom: b.bottom - b.top, left: 0}
}

func (b Bounds) BottomRightFrom(key gridKey) gridKey {
	return gridKey{key.line + b.bottom, key.beat + b.right}
}

func (b Bounds) TopLeft() gridKey {
	return gridKey{b.top, b.left}
}

func (m model) VisualSelectionBounds() Bounds {
	return InitBounds(m.cursorPos, m.visualAnchorCursor)
}

func (m model) PatternBounds() Bounds {
	return Bounds{0, m.definition.beats - 1, uint8(len(m.definition.lines)), 0}
}

func (m model) PasteBounds() Bounds {
	if m.visualMode {
		return m.VisualSelectionBounds()
	} else {
		return m.PatternBounds()
	}
}

func (m model) YankBounds() Bounds {
	if m.visualMode {
		return m.VisualSelectionBounds()
	} else {
		return InitBounds(m.cursorPos, m.cursorPos)
	}
}

func (m model) InVisualSelection(key gridKey) bool {
	return m.VisualSelectionBounds().InBounds(key)
}

func (n note) ViewComponents() (string, lipgloss.Color) {

	currentNote := n
	currentAccent := accents[currentNote.AccentIndex]
	currentAction := currentNote.Action
	var char string
	var foregroundColor lipgloss.Color
	if currentAction == ACTION_NOTHING {
		char = string(currentAccent.Shape) + string(ratchets[currentNote.Ratchets.Length])
		foregroundColor = currentAccent.Color
	} else {
		lineaction := lineactions[currentAction]
		char = string(lineaction.shape)
		foregroundColor = lineaction.color
	}

	return char, foregroundColor
}
