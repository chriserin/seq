package main

import (
	"fmt"
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
	colors "github.com/chriserin/seq/internal/colors"
	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/notereg"
	overlaykey "github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
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
	OverlayInputSwitch key.Binding
	SetupInputSwitch   key.Binding
	AccentInputSwitch  key.Binding
	RatchetInputSwitch key.Binding
	Increase           key.Binding
	Decrease           key.Binding
	ToggleGateMode     key.Binding
	ToggleWaitMode     key.Binding
	ToggleAccentMode   key.Binding
	ToggleRatchetMode  key.Binding
	NextOverlay        key.Binding
	PrevOverlay        key.Binding
	Save               key.Binding
	Undo               key.Binding
	Redo               key.Binding
	New                key.Binding
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
	GateIncrease         key.Binding
	GateDecrease         key.Binding
	WaitIncrease         key.Binding
	WaitDecrease         key.Binding
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
	definitionKeys.GateIncrease,
	definitionKeys.GateDecrease,
	definitionKeys.WaitIncrease,
	definitionKeys.WaitDecrease,
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
	Increase:           Key("Tempo Increase", "+", "="),
	Decrease:           Key("Tempo Decrease", "-"),
	TempoInputSwitch:   Key("Select Tempo Indicator", "ctrl+t"),
	OverlayInputSwitch: Key("Select Overlay Indicator", "ctrl+o"),
	SetupInputSwitch:   Key("Setup Input Indicator", "ctrl+s"),
	AccentInputSwitch:  Key("Accent Input Indicator", "ctrl+e"),
	RatchetInputSwitch: Key("Ratchet Input Indicator", "ctrl+h"),
	ToggleRatchetMode:  Key("Toggle Ratchet Mode", "ctrl+r"),
	ToggleGateMode:     Key("Toggle Gate Mode", "ctrl+g"),
	ToggleWaitMode:     Key("Toggle Wait Mode", "ctrl+w"),
	ToggleAccentMode:   Key("Toggle Wait Mode", "ctrl+a"),
	NextOverlay:        Key("Next Overlay", "{"),
	PrevOverlay:        Key("Prev Overlay", "}"),
	Save:               Key("Save", "ctrl+v"),
	Undo:               Key("Undo", "u"),
	Redo:               Key("Redo", "U"),
	ToggleVisualMode:   Key("Toggle Visual Mode", "v"),
	New:                Key("New", "ctrl+n"),
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
	GateIncrease:         Key("Gate Increase", "G"),
	GateDecrease:         Key("Gate Increase", "g"),
	WaitIncrease:         Key("Wait Increase", "W"),
	WaitDecrease:         Key("Wait Increase", "w"),
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

type groupPlayState uint

const (
	PLAY_STATE_PLAY groupPlayState = iota
	PLAY_STATE_MUTE
	PLAY_STATE_SOLO
)

type linestate struct {
	index               uint8
	currentBeat         uint8
	direction           int8
	resetDirection      int8
	resetLocation       uint8
	resetActionLocation uint8
	resetAction         action
	groupPlayState      groupPlayState
}

func (ls linestate) IsMuted() bool {
	return ls.groupPlayState == PLAY_STATE_MUTE
}

func (ls linestate) IsSolo() bool {
	return ls.groupPlayState == PLAY_STATE_SOLO
}

func (ls linestate) GridKey() grid.GridKey {
	return grid.GridKey{Line: ls.index, Beat: ls.currentBeat}
}

type overlayKey = overlaykey.OverlayPeriodicity

type gridKey = grid.GridKey

func GK(line uint8, beat uint8) gridKey {
	return gridKey{
		Line: line,
		Beat: beat,
	}
}

type note = grid.Note
type action = grid.Action
type overlay = overlays.Overlay

var zeronote note

type Delayable interface {
	Delay() time.Duration
}

type noteMsg struct {
	midiType  midi.Type
	channel   uint8
	noteValue uint8
	velocity  uint8
	delay     time.Duration
}

func (n noteMsg) Delay() time.Duration {
	return n.delay
}

type controlChangeMsg struct {
	channel uint8
	control uint8
	ccValue uint8
	delay   time.Duration
}

func (ccm controlChangeMsg) MidiMessage() midi.Message {
	return midi.ControlChange(ccm.channel, ccm.control, ccm.ccValue)
}

func (ccm controlChangeMsg) Delay() time.Duration {
	return ccm.delay
}

func (nm noteMsg) GetKey() notereg.NoteRegKey {
	return notereg.NoteRegKey{
		Channel: nm.channel,
		Note:    nm.noteValue,
	}
}

func (nm noteMsg) GetMidi() midi.Message {
	switch nm.midiType {
	case midi.NoteOnMsg:
		return midi.NoteOn(nm.channel, nm.noteValue, nm.velocity)
	case midi.NoteOffMsg:
		return midi.NoteOff(nm.channel, nm.noteValue)
	}
	panic("No message matching midiType")
}

func (nm noteMsg) OffMessage() midi.Message {
	return midi.NoteOff(nm.channel, nm.noteValue)
}

func NoteMessages(l grid.LineDefinition, accentValue uint8, gateIndex uint8, accentTarget accentTarget, delay time.Duration) (noteMsg, noteMsg) {
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
	duration := time.Duration(config.Gates[gateIndex].Value) * time.Millisecond

	return noteMsg{midi.NoteOnMsg, l.Channel, noteValue, velocityValue, delay},
		noteMsg{midi.NoteOffMsg, l.Channel, noteValue, 0, delay + duration}
}

func CCMessage(l grid.LineDefinition, note note, accents []config.Accent, beatInterval time.Duration, includeDelay bool, instrument string) controlChangeMsg {
	var delay time.Duration
	if includeDelay && note.WaitIndex != 0 {
		delay = time.Duration((float64(config.WaitPercentages[note.WaitIndex])) / float64(100) * float64(beatInterval))
	} else {
		delay = 0
	}

	ccValue := uint8((float32(note.AccentIndex) / float32(len(accents))) * float32(config.FindCC(l.Note, instrument).UpperLimit))

	return controlChangeMsg{l.Channel, l.Note, ccValue, delay}
}

type model struct {
	transitiveStatekeys transitiveKeyMap
	definitionKeys      definitionKeyMap
	help                help.Model
	cursor              cursor.Model
	overlayKeyEdit      overlaykey.Model
	cursorPos           gridKey
	visualAnchorCursor  gridKey
	visualMode          bool
	midiConnection      MidiConnection
	logFile             *os.File
	playing             bool
	beatTime            time.Time
	playEditing         bool
	playTime            time.Time
	trackTime           time.Duration
	totalBeats          int
	playState           []linestate
	hasSolo             bool
	selectionIndicator  Selection
	focus               focus
	patternMode         PatternMode
	ratchetCursor       uint8
	currentOverlay      *overlays.Overlay
	keyCycles           int
	undoStack           UndoStack
	redoStack           UndoStack
	yankBuffer          Buffer
	needsWrite          int
	// save everything below here
	definition Definition
}

type focus int

const (
	FOCUS_GRID focus = iota
	FOCUS_OVERLAY_KEY
)

type Selection uint8

const (
	SELECT_NOTHING Selection = iota
	SELECT_TEMPO
	SELECT_TEMPO_SUBDIVISION
	SELECT_OVERLAY
	SELECT_SETUP_CHANNEL
	SELECT_SETUP_MESSAGE_TYPE
	SELECT_SETUP_VALUE
	SELECT_RATCHETS
	SELECT_RATCHET_SPAN
	SELECT_ACCENT_DIFF
	SELECT_ACCENT_TARGET
	SELECT_ACCENT_START
)

type PatternMode uint8

const (
	PATTERN_FILL PatternMode = iota
	PATTERN_ACCENT
	PATTERN_GATE
	PATTERN_WAIT
	PATTERN_RATCHET
)

type Undoable interface {
	ApplyUndo(m *model) (overlayKey, gridKey)
}

type UndoStack struct {
	undo Undoable
	redo Undoable
	next *UndoStack
	id   int
}

var NIL_STACK = UndoStack{}

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

func (ugn UndoGridNote) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(ugn.overlayKey)
	overlay := m.definition.overlays.FindOverlay(ugn.overlayKey)
	overlay.SetNote(ugn.gridNote.gridKey, ugn.gridNote.note)
	return ugn.overlayKey, ugn.gridNote.gridKey
}

type UndoLineGridNotes struct {
	overlayKey
	cursorPosition gridKey
	line           uint8
	gridNotes      []GridNote
}

func (ulgn UndoLineGridNotes) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(ulgn.overlayKey)
	overlay := m.definition.overlays.FindOverlay(ulgn.overlayKey)
	for i := range m.definition.beats {
		overlay.RemoveNote(GK(ulgn.line, i))
	}
	for _, gridNote := range ulgn.gridNotes {
		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return ulgn.overlayKey, ulgn.cursorPosition
}

type UndoBounds struct {
	overlayKey
	cursorPosition gridKey
	bounds         Bounds
	gridNotes      []GridNote
}

func (uvs UndoBounds) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(uvs.overlayKey)
	overlay := m.definition.overlays.FindOverlay(uvs.overlayKey)
	for _, k := range uvs.bounds.GridKeys() {
		overlay.RemoveNote(k)
	}
	for _, gridNote := range uvs.gridNotes {
		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return uvs.overlayKey, uvs.cursorPosition
}

type UndoGridNotes struct {
	overlayKey
	gridNotes []GridNote
}

func (ugn UndoGridNotes) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(ugn.overlayKey)
	overlay := m.definition.overlays.FindOverlay(ugn.overlayKey)
	for _, gridNote := range ugn.gridNotes {

		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return ugn.overlayKey, ugn.gridNotes[0].gridKey
}

type UndoToNothing struct {
	overlayKey overlayKey
	location   gridKey
}

func (utn UndoToNothing) ApplyUndo(m *model) (overlayKey, gridKey) {
	overlay := m.definition.overlays.FindOverlay(utn.overlayKey)
	overlay.RemoveNote(utn.location)
	return utn.overlayKey, utn.location
}

type UndoLineToNothing struct {
	overlayKey     overlayKey
	cursorPosition gridKey
	line           uint8
}

func (ultn UndoLineToNothing) ApplyUndo(m *model) (overlayKey, gridKey) {
	overlay := m.definition.overlays.FindOverlay(ultn.overlayKey)
	for i := range m.definition.beats {
		overlay.RemoveNote(GK(ultn.line, i))
	}

	return ultn.overlayKey, ultn.cursorPosition
}

type UndoNewOverlay struct {
	overlayKey     overlayKey
	cursorPosition gridKey
}

func (uno UndoNewOverlay) ApplyUndo(m *model) (overlayKey, gridKey) {
	newOverlay := m.definition.overlays.Remove(uno.overlayKey)
	m.definition.overlays = newOverlay
	return uno.overlayKey, uno.cursorPosition
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
		ok, gk := undoStack.undo.ApplyUndo(m)
		m.cursorPos = gk
		overlay := m.definition.overlays.FindOverlay(ok)
		if overlay == nil {
			m.currentOverlay = m.definition.overlays
		} else {
			m.currentOverlay = overlay
		}
		m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
	}
	return undoStack
}

func (m *model) Redo() UndoStack {
	undoStack := m.PopRedo()
	if undoStack != NIL_STACK {
		ok, gk := undoStack.redo.ApplyUndo(m)
		m.cursorPos = gk
		m.currentOverlay = m.definition.overlays.FindOverlay(ok)
		m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
	}
	return undoStack
}

type Definition struct {
	overlays        *overlays.Overlay
	lines           []grid.LineDefinition
	beats           uint8
	tempo           int
	subdivisions    int
	keyline         uint8
	accents         patternAccents
	instrument      string
	template        string
	templateUIStyle string
}

type patternAccents struct {
	Diff   uint8
	Data   []config.Accent
	Start  uint8
	Target accentTarget
}

type accentTarget uint8

const (
	ACCENT_TARGET_NOTE accentTarget = iota
	ACCENT_TARGET_VELOCITY
)

func (pa *patternAccents) ReCalc() {
	accents := make([]config.Accent, 9)
	for i, a := range pa.Data[1:] {
		a.Value = pa.Start - pa.Diff*uint8(i)
		accents[i+1] = a
	}
	pa.Data = accents
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
	tickInterval := m.TickInterval()
	adjuster := time.Since(m.playTime) - m.trackTime
	m.trackTime = m.trackTime + tickInterval
	next := tickInterval - adjuster
	return next
}

func (m model) TickInterval() time.Duration {
	return time.Minute / time.Duration(m.definition.tempo*m.definition.subdivisions)
}

type lineNote struct {
	note
	line grid.LineDefinition
}

func (m model) PlayBeat(beatInterval time.Duration, pattern grid.Pattern, cmds *[]tea.Cmd) {

	lines := m.definition.lines
	ratchetNotes := make([]lineNote, 0, len(lines))

	for gridKey, note := range pattern {
		line := lines[gridKey.Line]
		if note.Ratchets.Length > 0 {
			ratchetNotes = append(ratchetNotes, lineNote{note, line})
		} else if note != zeronote {
			accents := m.definition.accents
			switch line.MsgType {
			case grid.MESSAGE_TYPE_NOTE:
				var delay time.Duration
				if note.WaitIndex != 0 {
					delay = time.Duration((float64(config.WaitPercentages[note.WaitIndex])) / float64(100) * float64(beatInterval))
				} else {
					delay = 0
				}
				onMessage, offMessage := NoteMessages(line, m.definition.accents.Data[note.AccentIndex].Value, note.GateIndex, accents.Target, delay)
				m.ProcessNoteMsg(onMessage)
				m.ProcessNoteMsg(offMessage)
			case grid.MESSAGE_TYPE_CC:
				ccMessage := CCMessage(line, note, accents.Data, beatInterval, true, m.definition.instrument)
				m.ProcessNoteMsg(ccMessage)
			}
		}
	}

	for _, ratchetNote := range ratchetNotes {
		*cmds = append(*cmds, func() tea.Msg {
			return ratchetMsg{ratchetNote, 0, beatInterval}
		})
	}
}

func PlayMessage(delay time.Duration, message midi.Message, sendFn SendFunc) {
	time.AfterFunc(delay, func() {
		err := sendFn(message)
		if err != nil {
			panic("midi message send failed")
		}
	})
}

func (m *model) EnsureOverlay() {
	m.EnsureOverlayWithKey(m.overlayKeyEdit.GetKey())
}

func (m *model) EnsureOverlayWithKey(key overlayKey) {
	if m.definition.overlays.FindOverlay(key) == nil {
		newOverlay := m.definition.overlays.Add(key)
		m.definition.overlays = newOverlay
		m.currentOverlay = newOverlay.FindOverlay(key)
	}
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
		m.currentOverlay.SetNote(k, grid.InitNote())
	}
}

func (m *model) AddAction(act action) {
	keys := m.VisualSelectedGridKeys()
	for _, k := range keys {
		m.currentOverlay.SetNote(k, grid.InitActionNote(act))
	}
}

func (m *model) RemoveTrigger() {
	keys := m.VisualSelectedGridKeys()
	for _, k := range keys {
		m.currentOverlay.SetNote(k, zeronote)
	}
}

func (m *model) OverlayRemoveTrigger() {
	keys := m.VisualSelectedGridKeys()
	for _, gridKey := range keys {
		m.currentOverlay.RemoveNote(gridKey)
	}
}

func (m *model) IncreaseRatchet() {
	combinedPattern := m.CombinedEditPattern(m.currentOverlay)
	bounds := m.YankBounds()

	for key, currentNote := range combinedPattern {
		if bounds.InBounds(key) {
			m.currentOverlay.SetNote(key, currentNote.IncrementRatchet(1))
		}
	}
}

func (m *model) DecreaseRatchet() {
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)
	bounds := m.YankBounds()

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			m.currentOverlay.SetNote(key, currentNote.IncrementRatchet(-1))
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
	if currentNote != zeronote && currentNote.Action == grid.ACTION_NOTHING {
		span := currentNote.Ratchets.Span
		if span < 8 {
			currentNote.Ratchets.Span = span + 1
		}
		m.currentOverlay.SetNote(m.cursorPos, currentNote)
	}
}

func (m *model) DecreaseSpan() {
	currentNote := m.CurrentNote()
	if currentNote != zeronote && currentNote.Action == grid.ACTION_NOTHING {
		span := currentNote.Ratchets.Span
		if span > 0 {
			currentNote.Ratchets.Span = span - 1
		}
		m.currentOverlay.SetNote(m.cursorPos, currentNote)
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
	currentNote.Ratchets.Toggle(m.ratchetCursor)
	m.currentOverlay.SetNote(m.cursorPos, currentNote)
}

func InitLines(template string) []grid.LineDefinition {
	gridTemplate := config.GetTemplate(template)
	newLines := make([]grid.LineDefinition, len(gridTemplate.Lines))
	copy(newLines, gridTemplate.Lines)
	return newLines
}

func InitLineStates(lines int, previousPlayState []linestate) []linestate {
	linestates := make([]linestate, 0, lines)

	for i := range uint8(lines) {
		var previousGroupPlayState = PLAY_STATE_PLAY
		if len(previousPlayState) > int(i) {
			previousState := previousPlayState[i]
			previousGroupPlayState = previousState.groupPlayState
		}

		linestates = append(linestates, InitLineState(previousGroupPlayState, i))
	}
	return linestates
}

func InitLineState(previousGroupPlayState groupPlayState, index uint8) linestate {
	return linestate{index, 0, 1, 1, 0, 0, 0, previousGroupPlayState}
}

func InitDefinition(template string, instrument string) Definition {
	gridTemplate := config.GetTemplate(template)
	newLines := make([]grid.LineDefinition, len(gridTemplate.Lines))
	copy(newLines, gridTemplate.Lines)

	return Definition{
		overlays:        overlays.InitOverlay(overlaykey.ROOT, nil),
		beats:           32,
		tempo:           120,
		keyline:         0,
		subdivisions:    2,
		lines:           newLines,
		accents:         patternAccents{Diff: 15, Data: config.Accents, Start: 120, Target: ACCENT_TARGET_VELOCITY},
		template:        template,
		instrument:      instrument,
		templateUIStyle: gridTemplate.UIStyle,
	}
}

func InitModel(midiConnection MidiConnection, template string, instrument string) model {
	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic("could not open log file")
	}

	newCursor := cursor.New()
	newCursor.Style = lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	definition, hasDefinition := Definition{}, false // Read()

	if !hasDefinition {
		definition = InitDefinition(template, instrument)
	}

	return model{
		transitiveStatekeys: transitiveKeys,
		definitionKeys:      definitionKeys,
		help:                help.New(),
		cursor:              newCursor,
		midiConnection:      midiConnection,
		logFile:             logFile,
		cursorPos:           GK(0, 0),
		//TODO: Initial overlay key should be read from file before setting to ROOT
		currentOverlay: definition.overlays,
		overlayKeyEdit: overlaykey.InitModel(),
		definition:     definition,
		playState:      InitLineStates(len(definition.lines), []linestate{}),
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

func (m model) LogFromBeatTime() {
	_, err := fmt.Fprintf(m.logFile, "%d\n", time.Since(m.beatTime))
	if err != nil {
		panic("could not write to log file")
	}
}

func RunProgram(midiConnection MidiConnection, template string, instrument string) *tea.Program {
	config.ProcessConfig("./config/init.lua")
	p := tea.NewProgram(InitModel(midiConnection, template, instrument), tea.WithAltScreen())
	return p
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return tea.FocusMsg{} }
}

func Is(msg tea.KeyMsg, k key.Binding) bool {
	return key.Matches(msg, k)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := transitiveKeys

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if Is(msg, keys.Quit) {
			m.logFile.Close()
			return m, tea.Quit
		}
		if m.focus == FOCUS_OVERLAY_KEY {
			okModel, cmd := m.overlayKeyEdit.Update(msg)
			m.overlayKeyEdit = okModel
			return m, cmd
		}
		switch {
		case Is(msg, keys.CursorDown):
			if slices.Contains([]Selection{SELECT_NOTHING, SELECT_SETUP_CHANNEL, SELECT_SETUP_MESSAGE_TYPE, SELECT_SETUP_VALUE}, m.selectionIndicator) {
				if m.cursorPos.Line < uint8(len(m.definition.lines)-1) {
					m.cursorPos.Line++
				}
			}
		case Is(msg, keys.CursorUp):
			if slices.Contains([]Selection{SELECT_NOTHING, SELECT_SETUP_CHANNEL, SELECT_SETUP_MESSAGE_TYPE, SELECT_SETUP_VALUE}, m.selectionIndicator) {
				if m.cursorPos.Line > 0 {
					m.cursorPos.Line--
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
				if m.cursorPos.Beat > 0 {
					m.cursorPos.Beat--
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
				if m.cursorPos.Beat < m.definition.beats-1 {
					m.cursorPos.Beat++
				}
			}
		case Is(msg, keys.CursorLineStart):
			m.cursorPos.Beat = 0
		case Is(msg, keys.CursorLineEnd):
			m.cursorPos.Beat = m.definition.beats - 1
		case Is(msg, keys.Escape):
			m.selectionIndicator = 0
			m.patternMode = PATTERN_FILL
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
			m.playEditing = false
			m.playTime = time.Now()
			if m.playing {
				m.keyCycles = 0
				m.totalBeats = 0
				m.playState = InitLineStates(len(m.definition.lines), m.playState)
				m.advanceKeyCycle()
				m.trackTime = time.Duration(0)
				playingOverlay := m.definition.overlays.HighestMatchingOverlay(m.keyCycles)
				beatInterval := m.BeatInterval()

				pattern := m.CombinedBeatPattern(playingOverlay)
				cmds := make([]tea.Cmd, 0, len(pattern))
				m.PlayBeat(beatInterval, pattern, &cmds)
				cmds = append(cmds, BeatTick(beatInterval))
				return m, tea.Batch(cmds...)
			} else {
				m.keyCycles = 0
				notes := notereg.Clear()
				sendFn := m.midiConnection.AcquireSendFunc()
				for _, n := range notes {
					switch n := n.(type) {
					case noteMsg:
						PlayMessage(time.Duration(0), n.OffMessage(), sendFn)
					}
				}
			}
		case Is(msg, keys.OverlayInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_OVERLAY}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
			m.focus = FOCUS_OVERLAY_KEY
			m.overlayKeyEdit.Focus(m.selectionIndicator == SELECT_OVERLAY)
		case Is(msg, keys.TempoInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_TEMPO, SELECT_TEMPO_SUBDIVISION}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
		case Is(msg, keys.SetupInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_SETUP_CHANNEL, SELECT_SETUP_MESSAGE_TYPE, SELECT_SETUP_VALUE}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
		case Is(msg, keys.AccentInputSwitch):
			states := []Selection{SELECT_NOTHING, SELECT_ACCENT_DIFF, SELECT_ACCENT_TARGET, SELECT_ACCENT_START}
			m.selectionIndicator = AdvanceSelectionState(states, m.selectionIndicator)
		case Is(msg, keys.RatchetInputSwitch):
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
			case SELECT_SETUP_CHANNEL:
				m.definition.lines[m.cursorPos.Line].IncrementChannel()
			case SELECT_SETUP_MESSAGE_TYPE:
				m.definition.lines[m.cursorPos.Line].IncrementMessageType()
			case SELECT_SETUP_VALUE:
				m.definition.lines[m.cursorPos.Line].IncrementNote()
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
			case SELECT_SETUP_CHANNEL:
				m.definition.lines[m.cursorPos.Line].DecrementChannel()
			case SELECT_SETUP_MESSAGE_TYPE:
				m.definition.lines[m.cursorPos.Line].DecrementMessageType()
			case SELECT_SETUP_VALUE:
				m.definition.lines[m.cursorPos.Line].DecrementNote()
			case SELECT_RATCHET_SPAN:
				m.DecreaseSpan()
			case SELECT_ACCENT_DIFF:
				m.DecreaseAccent()
			case SELECT_ACCENT_TARGET:
				m.DecreaseAccentTarget()
			case SELECT_ACCENT_START:
				m.DecreaseAccentStart()
			}
		case Is(msg, keys.ToggleGateMode):
			m.patternMode = PATTERN_GATE
		case Is(msg, keys.ToggleWaitMode):
			m.patternMode = PATTERN_WAIT
		case Is(msg, keys.ToggleAccentMode):
			m.patternMode = PATTERN_ACCENT
		case Is(msg, keys.ToggleRatchetMode):
			m.patternMode = PATTERN_RATCHET
		case Is(msg, keys.PrevOverlay):
			m.NextOverlay(-1)
			m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
		case Is(msg, keys.NextOverlay):
			m.NextOverlay(+1)
			m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
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
			m.cursorPos = GK(0, 0)
			m.definition = InitDefinition(m.definition.template, m.definition.instrument)
			m.currentOverlay = m.definition.overlays
			m.selectionIndicator = SELECT_NOTHING
		case Is(msg, keys.ToggleVisualMode):
			m.visualAnchorCursor = m.cursorPos
			m.visualMode = !m.visualMode
		case Is(msg, keys.NewLine):
			if len(m.definition.lines) < 16 {
				lastline := m.definition.lines[len(m.definition.lines)-1]
				m.definition.lines = append(m.definition.lines, grid.LineDefinition{
					Channel: lastline.Channel,
					Note:    lastline.Note + 1,
				})
				if m.playing {
					m.playState = append(m.playState, InitLineState(PLAY_STATE_PLAY, uint8(len(m.definition.lines)-1)))
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
				m.playState = Mute(m.playState, m.cursorPos.Line)
				m.hasSolo = m.HasSolo()
			}
		case Is(msg, keys.Solo):
			m.playState = Solo(m.playState, m.cursorPos.Line)
			m.hasSolo = m.HasSolo()
		default:
			m = m.UpdateDefinition(msg)
		}
	case overlaykey.UpdatedOverlayKey:
		if !msg.HasFocus {
			m.focus = FOCUS_GRID
			m.selectionIndicator = SELECT_NOTHING
		}
	case beatMsg:
		m.beatTime = time.Now()
		if m.playing {
			playingOverlay := m.definition.overlays.HighestMatchingOverlay(m.keyCycles)
			m.advanceCurrentBeat(playingOverlay)
			m.advanceKeyCycle()
			beatInterval := m.BeatInterval()

			pattern := make(grid.Pattern)
			playingOverlay.CurrentBeatOverlayPattern(&pattern, m.keyCycles, m.CurrentBeatGridKeys())
			cmds := make([]tea.Cmd, 0, len(pattern)+1)
			m.PlayBeat(beatInterval, pattern, &cmds)
			if len(cmds) > 0 {
				cmds = append(cmds, BeatTick(beatInterval))
				return m, tea.Batch(
					cmds...,
				)
			} else {
				return m, BeatTick(beatInterval)
			}
		}
	case ratchetMsg:
		if m.playing && msg.iterations < (msg.Ratchets.Length+1) {
			if msg.Ratchets.HitAt(msg.iterations) {
				onMessage, offMessage := NoteMessages(msg.line, m.definition.accents.Data[msg.AccentIndex].Value, msg.GateIndex, m.definition.accents.Target, 0)
				m.ProcessNoteMsg(onMessage)
				m.ProcessNoteMsg(offMessage)
			}
			if msg.iterations+1 < (msg.Ratchets.Length + 1) {
				ratchetTickCmd := RatchetTick(msg.lineNote, msg.iterations+1, msg.beatInterval)
				return m, ratchetTickCmd
			}
		}
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor

	return m, cmd
}

func (m model) ProcessNoteMsg(msg Delayable) {
	sendFn := m.midiConnection.AcquireSendFunc()
	switch msg := msg.(type) {
	case noteMsg:
		switch msg.midiType {
		case midi.NoteOnMsg:
			if notereg.Has(msg) {
				PlayMessage(msg.delay, msg.OffMessage(), sendFn)
			} else {
				if err := notereg.Add(msg); err != nil {
					panic("Added a note that was already there")
				}
			}
			m.LogFromBeatTime()
		case midi.NoteOffMsg:
			notereg.Remove(msg)
		}
		PlayMessage(msg.delay, msg.GetMidi(), sendFn)
	case controlChangeMsg:
		PlayMessage(msg.delay, msg.MidiMessage(), sendFn)
	}
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
	case Is(msg, keys.GateIncrease):
		m.GateModify(1)
	case Is(msg, keys.GateDecrease):
		m.GateModify(-1)
	case Is(msg, keys.WaitIncrease):
		m.WaitModify(1)
	case Is(msg, keys.WaitDecrease):
		m.WaitModify(-1)
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
		m.AddAction(grid.ACTION_LINE_RESET)
	case Is(msg, keys.ActionAddLineReverse):
		m.AddAction(grid.ACTION_LINE_REVERSE)
	case Is(msg, keys.ActionAddSkipBeat):
		m.AddAction(grid.ACTION_LINE_SKIP_BEAT)
	case Is(msg, keys.ActionAddReset):
		m.AddAction(grid.ACTION_RESET)
	case Is(msg, keys.ActionAddLineBounce):
		m.AddAction(grid.ACTION_LINE_BOUNCE)
	case Is(msg, keys.ActionAddLineDelay):
		m.AddAction(grid.ACTION_LINE_DELAY)
	case Is(msg, keys.SelectKeyLine):
		m.definition.keyline = m.cursorPos.Line
	case Is(msg, keys.PressDownOverlay):
		m.currentOverlay.ToggleOverlayStackOptions()
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
		switch m.patternMode {
		case PATTERN_FILL:
			m.fill(uint8(beatInterval))
		case PATTERN_ACCENT:
			m.incrementAccent(uint8(beatInterval), -1)
		case PATTERN_GATE:
			m.incrementGate(uint8(beatInterval), -1)
		case PATTERN_RATCHET:
			m.incrementRatchet(uint8(beatInterval), -1)
		case PATTERN_WAIT:
			m.incrementWait(uint8(beatInterval), -1)
		}
	}
	if IsShiftSymbol(msg.String()) {
		beatInterval := convertSymbolToInt(msg.String())
		switch m.patternMode {
		case PATTERN_FILL:
			m.fill(uint8(beatInterval))
		case PATTERN_ACCENT:
			m.incrementAccent(uint8(beatInterval), 1)
		case PATTERN_GATE:
			m.incrementGate(uint8(beatInterval), 1)
		case PATTERN_RATCHET:
			m.incrementRatchet(uint8(beatInterval), 1)
		case PATTERN_WAIT:
			m.incrementWait(uint8(beatInterval), 1)
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
	if m.playing {
		m.playEditing = true
	}
	m.ResetRedo()
	return m
}

func (m model) UndoableNote() Undoable {
	overlay := m.definition.overlays.FindOverlay(m.overlayKeyEdit.GetKey())
	if overlay == nil {
		return UndoNewOverlay{m.overlayKeyEdit.GetKey(), m.cursorPos}
	}
	currentNote, hasNote := overlay.Notes[m.cursorPos]
	if hasNote {
		return UndoGridNote{m.currentOverlay.Key, m.cursorPos, GridNote{m.cursorPos, currentNote}}
	} else {
		return UndoToNothing{m.currentOverlay.Key, m.cursorPos}
	}
}

func (m model) UndoableBounds(pointA, pointB gridKey) Undoable {
	overlay := m.definition.overlays.FindOverlay(m.overlayKeyEdit.GetKey())
	if overlay == nil {
		return UndoNewOverlay{m.overlayKeyEdit.GetKey(), m.cursorPos}
	}
	bounds := InitBounds(pointA, pointB)
	gridKeys := bounds.GridKeys()
	gridNotes := make([]GridNote, 0, len(gridKeys))
	for _, k := range gridKeys {
		currentNote, hasNote := overlay.Notes[k]
		if hasNote {
			gridNotes = append(gridNotes, GridNote{k, currentNote})
		}
	}
	return UndoBounds{m.currentOverlay.Key, m.cursorPos, bounds, gridNotes}
}

func (m model) UndoableLine() Undoable {
	overlay := m.definition.overlays.FindOverlay(m.overlayKeyEdit.GetKey())
	if overlay == nil {
		return UndoNewOverlay{m.overlayKeyEdit.GetKey(), m.cursorPos}
	}
	notesToUndo := make([]GridNote, 0, m.definition.beats)
	for i := range m.definition.beats {
		key := GK(m.cursorPos.Line, i)
		currentNote, hasNote := overlay.Notes[key]
		if hasNote {
			notesToUndo = append(notesToUndo, GridNote{key, currentNote})
		}
	}
	if len(notesToUndo) == 0 {
		return UndoLineToNothing{m.currentOverlay.Key, m.cursorPos, m.cursorPos.Line}
	}
	return UndoLineGridNotes{m.currentOverlay.Key, m.cursorPos, m.cursorPos.Line, notesToUndo}
}

func (m model) UndoableOverlay() Undoable {
	overlay := m.definition.overlays.FindOverlay(m.overlayKeyEdit.GetKey())
	if overlay == nil {
		return UndoNewOverlay{m.overlayKeyEdit.GetKey(), m.cursorPos}
	}
	notesToUndo := make([]GridNote, 0, m.definition.beats)
	for key, note := range overlay.Notes {
		notesToUndo = append(notesToUndo, GridNote{key, note})
	}
	return UndoGridNotes{m.currentOverlay.Key, notesToUndo}
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

func (m model) CurrentNote() note {
	note, _ := m.currentOverlay.GetNote(m.cursorPos)
	return note
}

func RemoveRootKey(keys []overlayKey) []overlayKey {
	index := slices.Index(keys, overlaykey.ROOT)
	if index >= 0 {
		return append(keys[:index], keys[index+1:]...)
	}
	return keys
}

func (m *model) NextOverlay(direction int) {
	switch direction {
	case 1:
		overlay := m.definition.overlays.FindAboveOverlay(m.currentOverlay.Key)
		m.currentOverlay = overlay
	case -1:
		if m.currentOverlay.Below != nil {
			m.currentOverlay = m.currentOverlay.Below
		}
	default:
	}
}

func (m *model) ClearOverlayLine() {
	for i := uint8(0); i < m.definition.beats; i++ {
		m.currentOverlay.RemoveNote(GK(m.cursorPos.Line, i))
	}
}

func (m *model) ClearOverlay() {
	m.definition.overlays.Remove(m.currentOverlay.Key)
}

func (m *model) RotateRight() {
	combinedPattern := m.CombinedEditPattern(m.currentOverlay)

	lineStart, lineEnd := m.PatternActionLineBoundaries()
	start, end := m.PatternActionBeatBoundaries()

	for l := lineStart; l <= lineEnd; l++ {
		lastNote := combinedPattern[GK(l, end)]
		previousNote := zeronote
		for i := start; i <= end; i++ {
			gridKey := GK(l, i)
			currentNote := combinedPattern[gridKey]

			m.currentOverlay.SetNote(gridKey, previousNote)
			previousNote = currentNote
		}
		m.currentOverlay.SetNote(GK(l, start), lastNote)
	}

}

func (m *model) RotateLeft() {
	combinedPattern := m.CombinedEditPattern(m.currentOverlay)

	lineStart, lineEnd := m.PatternActionLineBoundaries()
	start, end := m.PatternActionBeatBoundaries()
	for l := lineStart; l <= lineEnd; l++ {
		firstNote := combinedPattern[GK(l, start)]
		previousNote := zeronote
		for i := int8(end); i >= int8(start); i-- {
			gridKey := GK(l, uint8(i))
			currentNote := combinedPattern[gridKey]

			m.currentOverlay.SetNote(gridKey, previousNote)
			previousNote = currentNote
		}
		m.currentOverlay.SetNote(GK(l, end), firstNote)
	}

}

func Mute(playState []linestate, line uint8) []linestate {
	switch playState[line].groupPlayState {
	case PLAY_STATE_PLAY:
		playState[line].groupPlayState = PLAY_STATE_MUTE
	case PLAY_STATE_MUTE:
		playState[line].groupPlayState = PLAY_STATE_PLAY
	case PLAY_STATE_SOLO:
		playState[line].groupPlayState = PLAY_STATE_MUTE
	}
	return playState
}

func Solo(playState []linestate, line uint8) []linestate {
	switch playState[line].groupPlayState {
	case PLAY_STATE_PLAY:
		playState[line].groupPlayState = PLAY_STATE_SOLO
	case PLAY_STATE_MUTE:
		playState[line].groupPlayState = PLAY_STATE_SOLO
	case PLAY_STATE_SOLO:
		playState[line].groupPlayState = PLAY_STATE_PLAY
	}
	return playState
}

func (m model) HasSolo() bool {
	for _, state := range m.playState {
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
	combinedPattern := m.CombinedEditPattern(m.currentOverlay)
	bounds := m.YankBounds()
	capturedGridNotes := make([]GridNote, 0, len(combinedPattern))

	for key, note := range combinedPattern {
		if bounds.InBounds(key) {
			normalizedGridKey := GK(key.Line-bounds.top, key.Beat-bounds.left)
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
		newKey := GK(key.Line+keyModifier.Line, key.Beat+keyModifier.Beat)
		if bounds.InBounds(newKey) {
			m.currentOverlay.SetNote(newKey, gridNote.note)
		}
	}
}

func (m *model) advanceCurrentBeat(playingOverlay *overlays.Overlay) {
	pattern := make(grid.Pattern)
	playingOverlay.CombineActionPattern(&pattern, m.keyCycles)
	for i := range m.playState {
		doContinue := m.advancePlayState(pattern, i)
		if !doContinue {
			break
		}
	}
}

func (m *model) advancePlayState(combinedPattern grid.Pattern, lineIndex int) bool {
	currentState := m.playState[lineIndex]
	advancedBeat := int8(currentState.currentBeat) + currentState.direction

	if advancedBeat >= int8(m.definition.beats) || advancedBeat < 0 {
		// reset locations should be 1 time use.  Reset back to 0.
		if m.playState[lineIndex].resetLocation != 0 && combinedPattern[GK(uint8(lineIndex), currentState.resetActionLocation)].Action == currentState.resetAction {
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

	switch combinedPattern[GK(uint8(lineIndex), uint8(advancedBeat))].Action {
	case grid.ACTION_NOTHING:
		return true
	case grid.ACTION_LINE_RESET:
		m.playState[lineIndex].currentBeat = 0
	case grid.ACTION_LINE_REVERSE:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-2, 0))
		m.playState[lineIndex].direction = -1
		m.playState[lineIndex].resetLocation = uint8(max(advancedBeat-1, 0))
		m.playState[lineIndex].resetActionLocation = uint8(advancedBeat)
		m.playState[lineIndex].resetAction = grid.ACTION_LINE_REVERSE
	case grid.ACTION_LINE_BOUNCE:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-1, 0))
		m.playState[lineIndex].direction = -1
	case grid.ACTION_LINE_SKIP_BEAT:
		m.advancePlayState(combinedPattern, lineIndex)
	case grid.ACTION_LINE_DELAY:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-1, 0))
	case grid.ACTION_RESET:
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
	}
}

func (m model) PlayingOverlayKeys() []overlayKey {
	keys := make([]overlayKey, 0, 10)
	m.definition.overlays.GetMatchingOverlayKeys(&keys, m.keyCycles)
	return keys
}

func (m model) CombinedEditPattern(overlay *overlays.Overlay) grid.Pattern {
	pattern := make(grid.Pattern)
	overlay.CombinePattern(&pattern, overlay.Key.GetMinimumKeyCycle())
	return pattern
}

func (m model) CurrentBeatGridKeys() []grid.GridKey {
	result := make([]grid.GridKey, 0, len(m.playState))
	for _, linestate := range m.playState {
		if linestate.IsSolo() || (!linestate.IsMuted() && !m.hasSolo) {
			result = append(result, linestate.GridKey())
		}
	}
	return result
}

func (m model) CombinedBeatPattern(overlay *overlays.Overlay) grid.Pattern {
	pattern := make(grid.Pattern)
	overlay.CurrentBeatOverlayPattern(&pattern, m.keyCycles, m.CurrentBeatGridKeys())
	return pattern
}

func (m model) CombinedOverlayPattern(overlay *overlays.Overlay) overlays.OverlayPattern {
	pattern := make(overlays.OverlayPattern)
	if m.playing && !m.playEditing {
		m.definition.overlays.CombineOverlayPattern(&pattern, m.keyCycles)
	} else {
		overlay.CombineOverlayPattern(&pattern, overlay.Key.GetMinimumKeyCycle())
	}
	return pattern
}

func (m *model) Every(every uint8, everyFn func(gridKey)) {
	lineStart, lineEnd := m.PatternActionLineBoundaries()
	start, end := m.PatternActionBeatBoundaries()

	for l := lineStart; l <= lineEnd; l++ {
		for i := start; i <= end; i += every {
			everyFn(GK(l, i))
		}
	}
}

func (m *model) fill(every uint8) {
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	everyFn := func(gridKey gridKey) {
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.RemoveNote(gridKey)
		} else {
			m.currentOverlay.SetNote(gridKey, grid.InitNote())
		}
	}

	m.Every(every, everyFn)
}

func (m *model) incrementAccent(every uint8, modifier int8) {
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	everyFn := func(gridKey gridKey) {
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.currentOverlay.SetNote(gridKey, currentNote.IncrementAccent(modifier, uint8(len(config.Accents))))
		}
	}
	m.Every(every, everyFn)
}

func (m *model) incrementGate(every uint8, modifier int8) {
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	everyFn := func(gridKey gridKey) {
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.currentOverlay.SetNote(gridKey, currentNote.IncrementGate(modifier))
		}
	}
	m.Every(every, everyFn)
}

func (m *model) incrementRatchet(every uint8, modifier int8) {
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	everyFn := func(gridKey gridKey) {
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.currentOverlay.SetNote(gridKey, currentNote.IncrementRatchet(modifier))
		}
	}
	m.Every(every, everyFn)
}

func (m *model) incrementWait(every uint8, modifier int8) {
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	everyFn := func(gridKey gridKey) {
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.currentOverlay.SetNote(gridKey, currentNote.IncrementWait(modifier))
		}
	}
	m.Every(every, everyFn)
}

func (m *model) AccentModify(modifier int8) {
	bounds := m.YankBounds()
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			if currentNote != zeronote {
				m.currentOverlay.SetNote(key, currentNote.IncrementAccent(modifier, uint8(len(config.Accents))))
			}
		}
	}
}

func (m *model) GateModify(modifier int8) {
	bounds := m.YankBounds()
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			if currentNote != zeronote {
				m.currentOverlay.SetNote(key, currentNote.IncrementGate(modifier))
			}
		}
	}
}

func (m *model) WaitModify(modifier int8) {
	bounds := m.YankBounds()
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			if currentNote != zeronote {
				m.currentOverlay.SetNote(key, currentNote.IncrementWait(modifier))
			}
		}
	}
}

func (m model) PatternActionBeatBoundaries() (uint8, uint8) {
	if m.visualMode {
		if m.visualAnchorCursor.Beat < m.cursorPos.Beat {
			return m.visualAnchorCursor.Beat, m.cursorPos.Beat
		} else {
			return m.cursorPos.Beat, m.visualAnchorCursor.Beat
		}
	} else {
		return m.cursorPos.Beat, m.definition.beats - 1
	}
}

func (m model) PatternActionLineBoundaries() (uint8, uint8) {
	if m.visualMode {
		if m.visualAnchorCursor.Line < m.cursorPos.Line {
			return m.visualAnchorCursor.Line, m.cursorPos.Line
		} else {
			return m.cursorPos.Line, m.visualAnchorCursor.Line
		}
	} else {
		return m.cursorPos.Line, m.cursorPos.Line
	}
}

func (m *model) RemoveNote(gridKey gridKey) {
	if m.currentOverlay.Key == overlaykey.ROOT {
		m.currentOverlay.RemoveNote(gridKey)
	} else {
		m.currentOverlay.SetNote(gridKey, zeronote)
	}
}

func (m model) OverlayKeys() []overlayKey {
	keys := make([]overlayKey, 0, 10)
	m.definition.overlays.CollectKeys(&keys)
	return keys
}

func (m model) TempoView() string {
	var buf strings.Builder
	var tempo, division string
	switch m.selectionIndicator {
	case SELECT_TEMPO:
		tempo = colors.SelectedColor.Render(strconv.Itoa(m.definition.tempo))
		division = colors.NumberColor.Render(strconv.Itoa(m.definition.subdivisions))
	case SELECT_TEMPO_SUBDIVISION:
		tempo = colors.NumberColor.Render(strconv.Itoa(m.definition.tempo))
		division = colors.SelectedColor.Render(strconv.Itoa(m.definition.subdivisions))
	default:
		tempo = colors.NumberColor.Render(strconv.Itoa(m.definition.tempo))
		division = colors.NumberColor.Render(strconv.Itoa(m.definition.subdivisions))
	}
	heart := colors.HeartColor.Render("♡")
	buf.WriteString("             \n")
	buf.WriteString(colors.HeartColor.Render("   ♡♡♡☆ ☆♡♡♡ ") + "\n")
	buf.WriteString(colors.HeartColor.Render("  ♡    ◊    ♡") + "\n")
	buf.WriteString(colors.HeartColor.Render("  ♡  TEMPO  ♡") + "\n")
	buf.WriteString(fmt.Sprintf("  %s   %s   %s\n", heart, tempo, heart))
	buf.WriteString(colors.HeartColor.Render("   ♡ BEATS ♡") + "\n")
	buf.WriteString(fmt.Sprintf("    %s  %s  %s  \n", heart, division, heart))
	buf.WriteString(colors.HeartColor.Render("     ♡   ♡   ") + "\n")
	buf.WriteString(colors.HeartColor.Render("      ♡ ♡    ") + "\n")
	buf.WriteString(colors.HeartColor.Render("       †     ") + "\n")
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

	if m.patternMode == PATTERN_ACCENT || m.IsAccentSelector() {
		sideView = m.AccentKeyView()
	} else if (m.definition.overlays.Key == overlaykey.ROOT && len(m.definition.overlays.Notes) == 0) ||
		m.selectionIndicator == SELECT_SETUP_VALUE ||
		m.selectionIndicator == SELECT_SETUP_MESSAGE_TYPE ||
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
		accentDiffString = colors.SelectedColor.Render(fmt.Sprintf("%2d", accentDiff))
	} else {
		accentDiffString = colors.NumberColor.Render(fmt.Sprintf("%2d", accentDiff))
	}

	var accentTargetString string
	if m.selectionIndicator == SELECT_ACCENT_TARGET {
		accentTargetString = colors.SelectedColor.Render(fmt.Sprintf(" %s", accentTarget))
	} else {
		accentTargetString = colors.NumberColor.Render(fmt.Sprintf(" %s", accentTarget))
	}

	buf.WriteString(fmt.Sprintf(" ACCENTS %s %s\n", accentDiffString, accentTargetString))
	buf.WriteString("———————————————\n")
	startAccent := m.definition.accents.Data[1]

	var accentStartString string
	if m.selectionIndicator == SELECT_ACCENT_START {
		accentStartString = colors.SelectedColor.Render(fmt.Sprintf("%2d", accentStart))
	} else {
		accentStartString = colors.NumberColor.Render(fmt.Sprintf("%2d", accentStart))
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
		if uint8(i) == m.cursorPos.Line && m.selectionIndicator == SELECT_SETUP_CHANNEL {
			buf.WriteString(colors.SelectedColor.Render(fmt.Sprintf("%2d", line.Channel)))
		} else {
			buf.WriteString(colors.NumberColor.Render(fmt.Sprintf("%2d", line.Channel)))
		}

		var messageType string
		switch line.MsgType {
		case grid.MESSAGE_TYPE_NOTE:
			messageType = "NOTE"
		case grid.MESSAGE_TYPE_CC:
			messageType = "CC"
		}

		if uint8(i) == m.cursorPos.Line && m.selectionIndicator == SELECT_SETUP_MESSAGE_TYPE {
			messageType = fmt.Sprintf(" %s ", colors.SelectedColor.Render(messageType))
		} else {
			messageType = fmt.Sprintf(" %s ", messageType)
		}

		buf.WriteString(messageType)

		if uint8(i) == m.cursorPos.Line && m.selectionIndicator == SELECT_SETUP_VALUE {
			buf.WriteString(colors.SelectedColor.Render(strconv.Itoa(int(line.Note))))
		} else {
			buf.WriteString(colors.NumberColor.Render(strconv.Itoa(int(line.Note))))
		}
		buf.WriteString(fmt.Sprintf(" %s\n", LineValueName(line, m.definition.instrument)))
	}
	return buf.String()
}

func NoteName(note uint8) string {
	return fmt.Sprintf("%s%d", strings.ReplaceAll(midi.Note(note).Name(), "b", "♭"), midi.Note(note).Octave()-2)
}

func LineValueName(ld grid.LineDefinition, instrument string) string {
	switch ld.MsgType {
	case grid.MESSAGE_TYPE_NOTE:
		return NoteName(ld.Note)
	case grid.MESSAGE_TYPE_CC:
		return config.FindCC(ld.Note, instrument).Name
	}
	return ""
}

const SELECTED_OVERLAY_ARROW = "\u2192"

var currentPlayingColor lipgloss.Color = "#abfaa9"
var activePlayingColor lipgloss.Color = "#f34213"

func (m model) OverlaysView() string {
	var buf strings.Builder
	buf.WriteString("Overlays\n")
	buf.WriteString("——————————————\n")
	style := lipgloss.NewStyle().Background(seqOverlayColor)
	var playingOverlayKeys = m.PlayingOverlayKeys()
	for currentOverlay := m.definition.overlays; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		var playingSpacer = "   "
		var playing = ""
		if m.playing && playingOverlayKeys[0] == currentOverlay.Key {
			playing = lipgloss.NewStyle().Background(seqOverlayColor).Foreground(currentPlayingColor).Render(" \u25CF ")
			buf.WriteString(playing)
			playingSpacer = ""
		} else if m.playing && slices.Contains(playingOverlayKeys, currentOverlay.Key) {
			playing = lipgloss.NewStyle().Background(seqOverlayColor).Foreground(activePlayingColor).Render(" \u25C9 ")
			buf.WriteString(playing)
			playingSpacer = ""
		}
		var editing = ""
		if m.currentOverlay.Key == currentOverlay.Key {
			editing = " E"
		}
		var stackModifier = ""
		if currentOverlay.PressDown {
			stackModifier = " \u2193\u0332"
		} else if currentOverlay.PressUp {
			stackModifier = " \u2191\u0305"
		}

		overlayLine := fmt.Sprintf("%s%2s%2s", overlaykey.View(currentOverlay.Key), stackModifier, editing)

		buf.WriteString(playingSpacer)
		if m.playing && slices.Contains(playingOverlayKeys, currentOverlay.Key) {
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

var accentModeStyle = lipgloss.NewStyle().Background(config.Accents[1].Color).Foreground(lipgloss.Color("#000000"))

func (m model) ViewTriggerSeq() string {
	var buf strings.Builder
	var mode string
	visualCombinedPattern := m.CombinedOverlayPattern(m.currentOverlay)

	if m.patternMode == PATTERN_ACCENT {
		mode = " Accent Mode "
		buf.WriteString(fmt.Sprintf("    %s\n", accentModeStyle.Render(mode)))
	} else if m.patternMode == PATTERN_GATE {
		mode = " Gate Mode "
		buf.WriteString(fmt.Sprintf("    %s\n", accentModeStyle.Render(mode)))
	} else if m.patternMode == PATTERN_WAIT {
		mode = " Wait Mode "
		buf.WriteString(fmt.Sprintf("    %s\n", accentModeStyle.Render(mode)))
	} else if m.patternMode == PATTERN_RATCHET {
		mode = " Ratchet Mode "
		buf.WriteString(fmt.Sprintf("    %s\n", accentModeStyle.Render(mode)))
	} else if m.selectionIndicator == SELECT_RATCHETS || m.selectionIndicator == SELECT_RATCHET_SPAN {
		buf.WriteString(m.RatchetEditView())
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

func (m model) RatchetEditView() string {
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
			if currentNote.Ratchets.HitAt(i) {
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
		buf.WriteString(fmt.Sprintf(" Span %s ", colors.SelectedColor.Render(strconv.Itoa(int(currentNote.Ratchets.GetSpan())))))
	} else {
		buf.WriteString(fmt.Sprintf(" Span %s ", colors.NumberColor.Render(strconv.Itoa(int(currentNote.Ratchets.GetSpan())))))
	}
	buf.WriteString("\n")

	return buf.String()
}

func (m model) ViewOverlay() string {
	return m.overlayKeyEdit.ViewOverlay()
}

func (m model) CurrentOverlayView() string {
	var matchedKey overlayKey
	if m.playing {
		matchedKey = m.definition.overlays.HighestMatchingOverlay(m.keyCycles).Key
	} else {
		matchedKey = overlaykey.ROOT
	}

	var editOverlayTitle string
	if m.playEditing {
		editOverlayTitle = lipgloss.NewStyle().Background(seqOverlayColor).Render("Edit")
	} else {
		editOverlayTitle = "Edit"
	}

	editOverlay := fmt.Sprintf("%s %s", editOverlayTitle, lipgloss.PlaceHorizontal(11, 0, m.ViewOverlay()))
	playOverlay := fmt.Sprintf("Play %s", lipgloss.PlaceHorizontal(11, 0, overlaykey.View(matchedKey)))
	return fmt.Sprintf("   %s  %s", editOverlay, playOverlay)
}

var altSeqColor = lipgloss.Color("#222222")
var seqColor = lipgloss.Color("#000000")
var seqCursorColor = lipgloss.Color("#444444")
var seqVisualColor = lipgloss.Color("#aaaaaa")
var seqOverlayColor = lipgloss.Color("#333388")
var seqMiddleOverlayColor = lipgloss.Color("#405810")

func KeyLineIndicator(k uint8, l uint8) string {
	if k == l {
		return "K"
	} else {
		return " "
	}
}

var blackKeyStyle = lipgloss.NewStyle().Background(lipgloss.Color("#000")).Foreground(lipgloss.Color("#fff"))
var whiteKeyStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ccc")).Foreground(lipgloss.Color("#000"))

var blackNotes = []uint8{1, 3, 6, 8, 10}

func (m model) LineIndicator(lineNumber uint8) string {
	indicator := "│"
	if lineNumber == m.cursorPos.Line {
		indicator = colors.SelectedColor.Render("┤")
	}
	if len(m.playState) > int(lineNumber) && m.playState[lineNumber].groupPlayState == PLAY_STATE_MUTE {
		indicator = "M"
	}
	if len(m.playState) > int(lineNumber) && m.playState[lineNumber].groupPlayState == PLAY_STATE_SOLO {
		indicator = "S"
	}

	var lineName string
	if m.definition.templateUIStyle == "blackwhite" {
		notename := NoteName(m.definition.lines[lineNumber].Note)
		if slices.Contains(blackNotes, m.definition.lines[lineNumber].Note%12) {
			lineName = blackKeyStyle.Render(notename[0:4])
		} else {
			lineName = whiteKeyStyle.Render(notename)
		}
	} else {
		lineName = fmt.Sprintf("%d", lineNumber)
	}

	return fmt.Sprintf("%2s%s%s", lineName, KeyLineIndicator(m.definition.keyline, lineNumber), indicator)
}

func lineView(lineNumber uint8, m model, visualCombinedPattern overlays.OverlayPattern) string {
	var buf strings.Builder
	buf.WriteString(m.LineIndicator(lineNumber))

	for i := uint8(0); i < m.definition.beats; i++ {
		currentGridKey := GK(uint8(lineNumber), i)
		overlayNote, hasNote := visualCombinedPattern[currentGridKey]

		var backgroundSeqColor lipgloss.Color
		if m.playing && m.playState[lineNumber].currentBeat == i {
			backgroundSeqColor = seqCursorColor
		} else if m.visualMode && m.InVisualSelection(currentGridKey) {
			backgroundSeqColor = seqVisualColor
		} else if hasNote && overlayNote.HighestOverlay && overlayNote.OverlayKey != overlaykey.ROOT {
			backgroundSeqColor = seqOverlayColor
		} else if hasNote && !overlayNote.HighestOverlay && overlayNote.OverlayKey != overlaykey.ROOT {
			backgroundSeqColor = seqMiddleOverlayColor
		} else if i%8 > 3 {
			backgroundSeqColor = altSeqColor
		} else {
			backgroundSeqColor = seqColor
		}

		char, foregroundColor := ViewNoteComponents(overlayNote.Note)

		style := lipgloss.NewStyle().Background(backgroundSeqColor)
		if m.cursorPos.Line == uint8(lineNumber) && m.cursorPos.Beat == i {
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
		top:    min(cursorA.Line, cursorB.Line),
		right:  max(cursorA.Beat, cursorB.Beat),
		bottom: max(cursorA.Line, cursorB.Line),
		left:   min(cursorA.Beat, cursorB.Beat),
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
			keys = append(keys, GK(i, j))
		}
	}
	return keys
}

func (b Bounds) InBounds(key gridKey) bool {
	return key.Line >= b.top &&
		key.Line <= b.bottom &&
		key.Beat >= b.left &&
		key.Beat <= b.right
}

func (b Bounds) Normalized() Bounds {
	return Bounds{top: 0, right: b.right - b.left, bottom: b.bottom - b.top, left: 0}
}

func (b Bounds) BottomRightFrom(key gridKey) gridKey {
	return GK(key.Line+b.bottom, key.Beat+b.right)
}

func (b Bounds) TopLeft() gridKey {
	return GK(b.top, b.left)
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

func ViewNoteComponents(currentNote grid.Note) (string, lipgloss.Color) {

	currentAccent := config.Accents[currentNote.AccentIndex]
	currentAction := currentNote.Action
	var char string
	var foregroundColor lipgloss.Color
	var waitShape string
	if currentNote.WaitIndex > 0 {
		waitShape = "\u0320"
	}
	if currentAction == grid.ACTION_NOTHING && currentNote != zeronote {
		char = string(currentAccent.Shape) + string(config.Ratchets[currentNote.Ratchets.Length]) + string(config.Gates[currentNote.GateIndex].Shape) + waitShape
		foregroundColor = currentAccent.Color
	} else {
		lineaction := config.Lineactions[currentAction]
		char = string(lineaction.Shape)
		foregroundColor = lineaction.Color
	}

	return char, foregroundColor
}
