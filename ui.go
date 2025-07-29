package main

import (
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/notereg"
	"github.com/chriserin/seq/internal/operation"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
	"github.com/chriserin/seq/internal/seqmidi"
	themes "github.com/chriserin/seq/internal/themes"
	"github.com/chriserin/seq/internal/theory"
	midi "gitlab.com/gomidi/midi/v2"
)

func Key(help string, keyboardKey ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey...), key.WithHelp(keyboardKey[0], help))
}

type groupPlayState uint8

const (
	PlayStatePlay groupPlayState = iota
	PlayStateMute
	PlayStateSolo
)

type linestate struct {
	index               uint8
	currentBeat         uint8
	resetLocation       uint8
	resetActionLocation uint8
	resetAction         action
	groupPlayState      groupPlayState
	direction           int8
	resetDirection      int8
}

func (ls linestate) IsMuted() bool {
	return ls.groupPlayState == PlayStateMute
}

func (ls linestate) IsSolo() bool {
	return ls.groupPlayState == PlayStateSolo
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

var zeronote note

type Delayable interface {
	Delay() time.Duration
}

type noteMsg struct {
	channel   uint8
	noteValue uint8
	velocity  uint8
	midiType  midi.Type
	delay     time.Duration
	id        int
}

func (nm noteMsg) Delay() time.Duration {
	return nm.delay
}

type programChangeMsg struct {
	channel uint8
	pcValue uint8
	delay   time.Duration
}

func (pcm programChangeMsg) MidiMessage() midi.Message {
	return midi.ProgramChange(pcm.channel, pcm.pcValue)
}

func (pcm programChangeMsg) Delay() time.Duration {
	return pcm.delay
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

func (nm noteMsg) GetID() int {
	return nm.id
}

func (nm noteMsg) GetOnMidi() midi.Message {
	return midi.NoteOn(nm.channel, nm.noteValue, nm.velocity)
}

func (nm noteMsg) GetOffMidi() midi.Message {
	return midi.NoteOff(nm.channel, nm.noteValue)
}

func (nm noteMsg) OffMessage() midi.Message {
	return midi.NoteOff(nm.channel, nm.noteValue)
}

func NoteMessages(l grid.LineDefinition, accentValue uint8, gateLength time.Duration, accentTarget accentTarget, delay time.Duration) (noteMsg, noteMsg) {
	var noteValue uint8
	var velocityValue uint8

	switch accentTarget {
	case AccentTargetNote:
		noteValue = l.Note + accentValue
		velocityValue = 96
	case AccentTargetVelocity:
		noteValue = l.Note
		velocityValue = accentValue
	}

	id := rand.Int()
	onMsg := noteMsg{id: id, midiType: midi.NoteOnMsg, channel: l.Channel - 1, noteValue: noteValue, velocity: velocityValue, delay: delay}
	offMsg := noteMsg{id: id, midiType: midi.NoteOffMsg, channel: l.Channel - 1, noteValue: noteValue, velocity: 0, delay: delay + gateLength}

	return onMsg, offMsg
}

func CCMessage(l grid.LineDefinition, note note, accents []config.Accent, delay time.Duration, includeDelay bool, instrument string) controlChangeMsg {
	if note.Action == grid.ActionSpecificValue {
		return controlChangeMsg{l.Channel - 1, l.Note, note.AccentIndex, delay}
	} else {
		cc, _ := config.FindCC(l.Note, instrument)
		ccValue := uint8((float32((len(accents))-int(note.AccentIndex)) / float32(len(accents)-1)) * float32(cc.UpperLimit))
		if cc.UpperLimit == 1 && note.AccentIndex > 4 {
			ccValue = uint8(1)
		}

		return controlChangeMsg{l.Channel - 1, l.Note, ccValue, delay}
	}
}

func PCMessage(l grid.LineDefinition, note note, accents []config.Accent, delay time.Duration, includeDelay bool, instrument string) programChangeMsg {
	if note.Action == grid.ActionSpecificValue {
		return programChangeMsg{l.Channel - 1, note.AccentIndex, delay}
	} else {
		return programChangeMsg{l.Channel - 1, l.Note - 1, delay}
	}
}

type model struct {
	hasUIFocus            bool
	connected             bool
	visualMode            bool
	logFileAvailable      bool
	playEditing           bool
	hasSolo               bool
	showArrangementView   bool
	ratchetCursor         uint8
	focus                 operation.Focus
	sectionSideIndicator  SectionSide
	playing               PlayMode
	loopMode              LoopMode
	selectionIndicator    operation.Selection
	patternMode           operation.PatternMode
	midiLoopMode          MidiLoopMode
	gridCursor            gridKey
	visualAnchorCursor    gridKey
	partSelectorIndex     int
	needsWrite            int
	playState             []linestate
	currentOverlay        *overlays.Overlay
	logFile               *os.File
	undoStack             UndoStack
	redoStack             UndoStack
	yankBuffer            Buffer
	beatTime              time.Time
	programChannel        chan midiEventLoopMsg
	lockReceiverChannel   chan bool
	unlockReceiverChannel chan bool
	errChan               chan error
	currentError          error
	theme                 string
	filename              string
	help                  help.Model
	cursor                cursor.Model
	overlayKeyEdit        overlaykey.Model
	arrangement           arrangement.Model
	textInput             textinput.Model
	midiConnection        seqmidi.MidiConnection
	activeChord           overlays.OverlayChord
	temporaryState        temporaryState
	// save everything below here
	definition Definition
}

func (m *model) SetGridCursor(key gridKey) {
	m.gridCursor = key
	note, exists := m.CurrentNote()
	if exists && note.Action == grid.ActionSpecificValue {
		m.SetSelectionIndicator(operation.SelectSpecificValue)
	} else {
		m.SetSelectionIndicator(operation.SelectNothing)
	}
}

func (m *model) UnsetActiveChord() {
	m.activeChord = overlays.OverlayChord{}
}

func (m *model) SetCurrentError(err error) {
	m.currentError = err
	m.selectionIndicator = operation.SelectError
	m.LogError(err)
}

func (m *model) ResetCurrentOverlay() {
	if m.playing != PlayStopped && m.playEditing {
		return
	}
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	if currentNode != nil && currentNode.IsEndNode() {
		partID := currentNode.Section.Part
		if len(*m.definition.parts) > partID {
			m.currentOverlay = (*m.definition.parts)[partID].Overlays
			m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
		}
	}
}

type SectionSide uint8

const (
	SectionAfter SectionSide = iota
	SectionBefore
)

type PlayMode uint8

const (
	PlayStopped PlayMode = iota
	PlayStandard
	PlayReceiver
)

type LoopMode uint8

const (
	LoopSong LoopMode = iota
	LoopPart
	LoopOverlay
)

type GridNote struct {
	gridKey gridKey
	note    note
}

func (m *model) PushArrUndo(arrundo arrangement.Undo) {
	m.PushUndoables(UndoArrangement{arrUndo: arrundo.Undo}, UndoArrangement{arrUndo: arrundo.Redo})
}

func (m *model) PushUndoables(undo Undoable, redo Undoable) {
	if m.undoStack == EmptyStack {
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
	if m.undoStack == EmptyStack {
		undo.next = nil
		m.undoStack = undo
	} else {
		pusheddown := m.undoStack
		undo.next = &pusheddown
		m.undoStack = undo
	}
}

func (m *model) PushRedo(redo UndoStack) {
	if m.redoStack == EmptyStack {
		redo.next = nil
		m.redoStack = redo
	} else {
		pusheddown := m.redoStack
		redo.next = &pusheddown
		m.redoStack = redo
	}
}

func (m *model) ResetRedo() {
	m.redoStack = EmptyStack
}

func (m *model) PopUndo() UndoStack {
	firstout := m.undoStack
	if firstout != EmptyStack && firstout.next != nil {
		m.undoStack = *m.undoStack.next
	} else {
		m.undoStack = EmptyStack
	}
	return firstout
}

func (m *model) PopRedo() UndoStack {
	firstout := m.redoStack
	if firstout != EmptyStack && firstout.next != nil {
		m.redoStack = *m.redoStack.next
	} else {
		m.redoStack = EmptyStack
	}
	return firstout
}

func (m *model) Undo() UndoStack {
	undoStack := m.PopUndo()
	if undoStack != EmptyStack {
		location := undoStack.undo.ApplyUndo(m)
		if location.ApplyLocation {
			m.ApplyLocation(location)
		}
	}
	return undoStack
}

func (m *model) Redo() UndoStack {
	undoStack := m.PopRedo()
	if undoStack != EmptyStack {
		location := undoStack.redo.ApplyUndo(m)
		if location.ApplyLocation {
			m.ApplyLocation(location)
		}
	}
	return undoStack
}

func (m *model) ApplyLocation(location Location) {
	m.SetGridCursor(location.GridKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(location.OverlayKey)
	if overlay == nil {
		m.currentOverlay = m.CurrentPart().Overlays
	} else {
		m.currentOverlay = overlay
	}
	m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
	m.focus = operation.FocusGrid
	m.arrangement.Focus = false
}

type Definition struct {
	parts                 *[]arrangement.Part
	arrangement           *arrangement.Arrangement
	lines                 []grid.LineDefinition
	tempo                 int
	subdivisions          int
	keyline               uint8
	accents               patternAccents
	instrument            string
	template              string
	templateUIStyle       string
	templateSequencerType operation.SequencerMode
}

type temporaryState struct {
	lines        []grid.LineDefinition
	tempo        int
	subdivisions int
	accents      patternAccents
	beats        uint8
	active       bool
}

type patternAccents struct {
	Diff   uint8
	Data   []config.Accent
	Start  uint8
	Target accentTarget
}

type accentTarget uint8

const (
	AccentTargetNote accentTarget = iota
	AccentTargetVelocity
)

func (pa *patternAccents) ReCalc() {
	accents := make([]config.Accent, 9)
	for i, a := range pa.Data[1:] {
		a.Value = pa.Start - pa.Diff*uint8(i)
		accents[i+1] = a
	}
	pa.Data = accents
}

func (pa *patternAccents) Equal(other *patternAccents) bool {
	if pa.Diff != other.Diff {
		return false
	}
	if pa.Start != other.Start {
		return false
	}
	if pa.Target != other.Target {
		return false
	}
	if !slices.Equal(pa.Data, other.Data) {
		return false
	}
	return true
}

type StateDiff struct {
	LinesChanged        bool
	TempoChanged        bool
	SubdivisionsChanged bool
	AccentsChanged      bool
	OldTempo            int
	NewTempo            int
	OldSubdivisions     int
	NewSubdivisions     int
	OldAccents          patternAccents
	NewAccents          patternAccents
	OldLines            []grid.LineDefinition
	NewLines            []grid.LineDefinition
}

func (s StateDiff) Changed() bool {
	return s.LinesChanged || s.TempoChanged || s.SubdivisionsChanged || s.AccentsChanged
}

func (s StateDiff) Reverse() StateDiff {
	return StateDiff{
		LinesChanged:        s.LinesChanged,
		TempoChanged:        s.TempoChanged,
		SubdivisionsChanged: s.SubdivisionsChanged,
		AccentsChanged:      s.AccentsChanged,
		OldTempo:            s.NewTempo,
		NewTempo:            s.OldTempo,
		OldSubdivisions:     s.NewSubdivisions,
		NewSubdivisions:     s.OldSubdivisions,
		OldAccents:          s.NewAccents,
		NewAccents:          s.OldAccents,
		OldLines:            s.NewLines,
		NewLines:            s.OldLines,
	}
}

func (s StateDiff) Apply(m *model) {
	if s.AccentsChanged {
		m.definition.accents = s.NewAccents
	}
	if s.LinesChanged {
		m.definition.lines = s.NewLines
	}
	if s.TempoChanged {
		m.definition.tempo = s.NewTempo
		m.SyncTempo()
	}
	if s.SubdivisionsChanged {
		m.definition.subdivisions = s.NewSubdivisions
	}
}

func createStateDiff(definition *Definition, temporary *temporaryState) StateDiff {
	diff := StateDiff{}

	if !slices.Equal(definition.lines, temporary.lines) {
		diff.LinesChanged = true
		diff.OldLines = definition.lines
		diff.NewLines = temporary.lines
	}

	if definition.tempo != temporary.tempo {
		diff.TempoChanged = true
		diff.OldTempo = definition.tempo
		diff.NewTempo = temporary.tempo
	}

	if definition.subdivisions != temporary.subdivisions {
		diff.SubdivisionsChanged = true
		diff.OldSubdivisions = definition.subdivisions
		diff.NewSubdivisions = temporary.subdivisions
	}

	if !definition.accents.Equal(&temporary.accents) {
		diff.AccentsChanged = true
		diff.OldAccents = definition.accents
		diff.NewAccents = temporary.accents
	}

	return diff
}

type beatMsg struct {
	interval time.Duration
}
type uiStartMsg struct{}
type uiStopMsg struct{}
type uiConnectedMsg struct{}
type uiNotConnectedMsg struct{}
type errorMsg struct {
	error error
}

func (m model) TickInterval() time.Duration {
	return time.Minute / time.Duration(m.definition.tempo*m.definition.subdivisions)
}

func (m model) SyncTempo() {
	go func() {
		m.programChannel <- tempoMsg{
			tempo:        m.definition.tempo,
			subdivisions: m.definition.subdivisions,
		}
	}()
}

func (m model) ProcessRatchets(note grid.Note, beatInterval time.Duration, line grid.LineDefinition) {
	for i := range note.Ratchets.Length + 1 {
		if note.Ratchets.HitAt(i) {
			shortGateLength := 20 * time.Millisecond
			ratchetInterval := time.Duration(i) * note.Ratchets.Interval(beatInterval)
			onMessage, offMessage := NoteMessages(line, m.definition.accents.Data[note.AccentIndex].Value, shortGateLength, m.definition.accents.Target, ratchetInterval)
			err := m.ProcessNoteMsg(onMessage)
			if err != nil {
				m.SetCurrentError(fault.Wrap(err, fmsg.With("cannot turn on ratchet note")))
			}
			err = m.ProcessNoteMsg(offMessage)
			if err != nil {
				m.SetCurrentError(fault.Wrap(err, fmsg.With("cannot turn off ratchet note")))
			}
		}
	}
}

func (m model) PlayBeat(beatInterval time.Duration, pattern grid.Pattern) error {

	lines := m.definition.lines

	for gridKey, note := range pattern {
		line := lines[gridKey.Line]
		if note.Ratchets.Length > 0 {
			m.ProcessRatchets(note, beatInterval, line)
		} else if note != zeronote {
			accents := m.definition.accents

			delay := Delay(note.WaitIndex, beatInterval)
			gateLength := GateLength(note.GateIndex, beatInterval)

			switch line.MsgType {
			case grid.MessageTypeNote:
				onMessage, offMessage := NoteMessages(
					line,
					m.definition.accents.Data[note.AccentIndex].Value,
					gateLength,
					accents.Target,
					delay,
				)
				err := m.ProcessNoteMsg(onMessage)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process note on msg"))
				}
				err = m.ProcessNoteMsg(offMessage)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process note off msg"))
				}
			case grid.MessageTypeCc:
				ccMessage := CCMessage(line, note, accents.Data, delay, true, m.definition.instrument)
				err := m.ProcessNoteMsg(ccMessage)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process cc msg"))
				}
			case grid.MessageTypeProgramChange:
				pcMessage := PCMessage(line, note, accents.Data, delay, true, m.definition.instrument)
				err := m.ProcessNoteMsg(pcMessage)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process cc msg"))
				}
			}
		}
	}

	return nil
}

func Delay(waitIndex uint8, beatInterval time.Duration) time.Duration {
	var delay time.Duration
	if waitIndex != 0 {
		delay = time.Duration((float64(config.WaitPercentages[waitIndex])) / float64(100) * float64(beatInterval))
	} else {
		delay = 0
	}
	return delay
}

func GateLength(gateIndex uint8, beatInterval time.Duration) time.Duration {
	var delay time.Duration
	if gateIndex < 8 {
		var delay time.Duration
		var value = config.ShortGates[gateIndex].Value
		if value > 1 {
			delay = time.Duration(config.ShortGates[gateIndex].Value) * time.Millisecond
		} else {
			delay = time.Duration(config.ShortGates[gateIndex].Value * float32(beatInterval))
		}
		return delay
	} else if gateIndex >= 8 {
		return time.Duration(float64(config.LongGates[gateIndex].Value) * float64(beatInterval))
	}
	return delay
}

func PlayMessage(delay time.Duration, message midi.Message, sendFn seqmidi.SendFunc, errChan chan error) {
	time.AfterFunc(delay, func() {
		err := sendFn(message)
		if err != nil {
			errChan <- fault.Wrap(err, fmsg.With("cannot send play message"))
		}
	})
}

func PlayOffMessage(nm noteMsg, sendFn seqmidi.SendFunc, errChan chan error) {
	time.AfterFunc(nm.delay, func() {
		if notereg.RemoveID(nm) {
			err := sendFn(nm.GetOffMidi())
			if err != nil {
				errChan <- fault.Wrap(err, fmsg.With("cannot send off message"))
			}
		}
	})
}

func (m *model) EnsureOverlay() {
	m.EnsureOverlayWithKey(m.overlayKeyEdit.GetKey())
}

func (m model) FindOverlay(key overlayKey) *overlays.Overlay {
	return m.CurrentPart().Overlays.FindOverlay(key)
}

func (m *model) EnsureOverlayWithKey(key overlayKey) {
	partID := m.CurrentPartID()
	if m.FindOverlay(key) == nil {
		var newOverlay = m.CurrentPart().Overlays.Add(key)
		(*m.definition.parts)[partID].Overlays = newOverlay
		m.currentOverlay = newOverlay.FindOverlay(key)
	}
}

func InitBounds(cursorA, cursorB gridKey) grid.Bounds {
	return grid.Bounds{
		Top:    min(cursorA.Line, cursorB.Line),
		Right:  max(cursorA.Beat, cursorB.Beat),
		Bottom: max(cursorA.Line, cursorB.Line),
		Left:   min(cursorA.Beat, cursorB.Beat),
	}
}

func (m model) VisualSelectionBounds() grid.Bounds {
	return InitBounds(m.gridCursor, m.visualAnchorCursor)
}

func (m model) PatternBounds() grid.Bounds {
	return grid.Bounds{
		Top:    0,
		Right:  m.CurrentPart().Beats - 1,
		Bottom: uint8(len(m.definition.lines)),
		Left:   0,
	}
}

func (m model) PasteBounds() grid.Bounds {
	if m.visualMode {
		return m.VisualSelectionBounds()
	} else {
		return m.PatternBounds()
	}
}

func (m model) YankBounds() grid.Bounds {
	if m.visualMode {
		return m.VisualSelectionBounds()
	} else {
		return InitBounds(m.gridCursor, m.gridCursor)
	}
}

func (m model) InVisualSelection(key gridKey) bool {
	return m.VisualSelectionBounds().InBounds(key)
}

func (m model) VisualSelectedGridKeys() []gridKey {
	if m.visualMode {
		return InitBounds(m.visualAnchorCursor, m.gridCursor).GridKeys()
	} else {
		return []gridKey{m.gridCursor}
	}
}

func (m *model) AddNote() {
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

func (m *model) RemoveNote() {
	keys := m.VisualSelectedGridKeys()
	for _, k := range keys {
		m.currentOverlay.SetNote(k, zeronote)
	}
}

func (m *model) OverlayRemoveNote() {
	keys := m.VisualSelectedGridKeys()
	for _, gridKey := range keys {
		m.currentOverlay.RemoveNote(gridKey)
	}
}

func (m *model) RatchetModify(modifier int8) {
	modifyFunc := func(key gridKey, currentNote note) {
		m.currentOverlay.SetNote(key, currentNote.IncrementRatchet(modifier))
	}
	m.Modify(modifyFunc)
}

func (m *model) EnsureRatchetCursorVisible() {
	currentNote, _ := m.CurrentNote()
	if m.ratchetCursor > currentNote.Ratchets.Length {
		m.ratchetCursor = m.ratchetCursor - 1
	}
}

func (m *model) IncrementCC() {
	note := m.definition.lines[m.gridCursor.Line].Note
	for i := note + 1; i <= 127; i++ {
		_, exists := config.FindCC(i, m.definition.instrument)
		if exists {
			m.definition.lines[m.gridCursor.Line].Note = i
			return
		}
	}
}

func (m *model) DecrementCC() {
	note := m.definition.lines[m.gridCursor.Line].Note
	for i := note - 1; i != 255; i-- {
		_, exists := config.FindCC(i, m.definition.instrument)
		if exists {
			m.definition.lines[m.gridCursor.Line].Note = i
			return
		}
	}
}

func (m *model) IncreaseSpan() {
	currentNote, _ := m.CurrentNote()
	if currentNote != zeronote && currentNote.Action == grid.ActionNothing {
		span := currentNote.Ratchets.Span
		if span < 8 {
			currentNote.Ratchets.Span = span + 1
		}
		m.currentOverlay.SetNote(m.gridCursor, currentNote)
	}
}

func (m *model) DecreaseSpan() {
	currentNote, _ := m.CurrentNote()
	if currentNote != zeronote && currentNote.Action == grid.ActionNothing {
		span := currentNote.Ratchets.Span
		if span > 0 {
			currentNote.Ratchets.Span = span - 1
		}
		m.currentOverlay.SetNote(m.gridCursor, currentNote)
	}
}

func (m *model) IncreaseAccentDiff() {
	m.definition.accents.Diff = m.definition.accents.Diff + 1
	m.definition.accents.ReCalc()

	// NOTE: uint8 values cannot be less than 0, so instead we check
	// if the value wrapped around to a value greater than 127
	// and set it back to the original value if so
	if m.definition.accents.Data[len(m.definition.accents.Data)-1].Value > 127 {
		m.definition.accents.Diff -= 1
		m.definition.accents.ReCalc()
	}
}

func (m *model) DecreaseAccentDiff() {
	m.definition.accents.Diff = m.definition.accents.Diff - 1
	m.definition.accents.ReCalc()
}

func (m *model) DecreaseAccentTarget() {
	m.definition.accents.Target = (m.definition.accents.Target + 1) % 2
}

func (m *model) IncreaseTempo(amount int) {
	newAmount := m.definition.tempo + amount
	if newAmount <= 300 {
		m.definition.tempo = newAmount
		m.SyncTempo()
	} else if m.definition.tempo == 300 {
		// do nothing if already at 300
	} else if newAmount > 300 {
		m.definition.tempo = 300
		m.SyncTempo()
	}
}

func (m *model) DecreaseTempo(amount int) {
	newAmount := m.definition.tempo - amount
	if newAmount > 30 {
		m.definition.tempo = newAmount
		m.SyncTempo()
	} else if m.definition.tempo == 30 {
		// do nothing if already at 30
	} else if newAmount < 30 {
		m.definition.tempo = 30
		m.SyncTempo()
	}
}

func (m *model) IncreaseAccentStart() {
	if m.definition.accents.Start < 127 {
		m.definition.accents.Start = m.definition.accents.Start + 1
		m.definition.accents.ReCalc()
	}
}

func (m *model) DecreaseAccentStart() {
	if m.definition.accents.Data[len(m.definition.accents.Data)-1].Value > 0 {
		m.definition.accents.Start = m.definition.accents.Start - 1
		m.definition.accents.ReCalc()
	}
}

func (m *model) IncreaseBeats() {
	newBeats := m.CurrentPart().Beats + 1
	if newBeats < 128 {
		(*m.definition.parts)[m.CurrentPartID()].Beats = newBeats
	}
}

func (m *model) DecreaseBeats() {
	newBeats := int(m.CurrentPart().Beats) - 1
	if newBeats >= 1 {
		(*m.definition.parts)[m.CurrentPartID()].Beats = uint8(newBeats)
		if m.gridCursor.Beat >= uint8(newBeats) {
			m.SetGridCursor(gridKey{
				Line: m.gridCursor.Line,
				Beat: uint8(newBeats - 1),
			})
		}
	}
}

func (m *model) IncreaseStartBeats() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	partBeats := m.CurrentPart().Beats
	currentNode.Section.IncreaseStartBeats(int(partBeats))
}

func (m *model) DecreaseStartBeats() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.DecreaseStartBeats()
}

func (m *model) IncreaseStartCycles() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.IncreaseStartCycles()
}

func (m *model) DecreaseStartCycles() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.DecreaseStartCycles()
}

func (m *model) IncreaseCycles() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.IncreaseCycles()
}

func (m *model) DecreaseCycles() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.DecreaseCycles()
}

func (m *model) IncrementPlayCycles() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.IncrementPlayCycles()
}

func (m *model) SetPlayCycles(keyCycles int) {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.SetPlayCycles(keyCycles)
}

func (m *model) DuringPlayResetPlayCycles() {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	currentNode.Section.DuringPlayReset()
}

func (m *model) IncreasePartSelector() {
	newIndex := m.partSelectorIndex + 1
	if newIndex < len(*m.definition.parts) {
		m.partSelectorIndex = newIndex
	}
}

func (m *model) DecreasePartSelector() {
	newIndex := m.partSelectorIndex - 1
	if newIndex > -2 {
		m.partSelectorIndex = newIndex
	}
}

func (m *model) ToggleRatchetMute() {
	currentNote, _ := m.CurrentNote()
	currentNote.Ratchets.Toggle(m.ratchetCursor)
	m.currentOverlay.SetNote(m.gridCursor, currentNote)
}

func InitLineStates(lines int, previousPlayState []linestate, startBeat uint8) []linestate {
	linestates := make([]linestate, 0, lines)

	for i := range uint8(lines) {
		var previousGroupPlayState = PlayStatePlay
		if len(previousPlayState) > int(i) {
			previousState := previousPlayState[i]
			previousGroupPlayState = previousState.groupPlayState
		}

		linestates = append(linestates, InitLineState(previousGroupPlayState, i, startBeat))
	}
	return linestates
}

func InitLineState(previousGroupPlayState groupPlayState, index uint8, startBeat uint8) linestate {
	return linestate{
		index:               index,
		currentBeat:         startBeat,
		direction:           1,
		resetDirection:      1,
		resetLocation:       0,
		resetActionLocation: 0,
		resetAction:         0,
		groupPlayState:      previousGroupPlayState,
	}
}

func InitDefinition(template string, instrument string) Definition {
	gridTemplate, exists := config.GetTemplate(template)
	if !exists {
		gridTemplate = config.GetDefaultTemplate()
	}
	config.LongGates = config.GetGateLengths(gridTemplate.MaxGateLength)
	newLines := make([]grid.LineDefinition, len(gridTemplate.Lines))
	copy(newLines, gridTemplate.Lines)

	parts := InitParts()
	return Definition{
		parts:                 &parts,
		arrangement:           InitArrangement(parts),
		tempo:                 120,
		keyline:               0,
		subdivisions:          2,
		lines:                 newLines,
		accents:               patternAccents{Diff: 15, Data: config.Accents, Start: 120, Target: AccentTargetVelocity},
		template:              gridTemplate.Name,
		instrument:            instrument,
		templateUIStyle:       gridTemplate.UIStyle,
		templateSequencerType: gridTemplate.SequencerType,
	}
}

func InitArrangement(parts []arrangement.Part) *arrangement.Arrangement {
	root := arrangement.InitRoot(parts)

	for i := range parts {
		section := arrangement.InitSongSection(i)

		node := &arrangement.Arrangement{
			Section:    section,
			Iterations: 1,
		}

		root.Nodes = append(root.Nodes, node)
	}

	return root
}

func InitParts() []arrangement.Part {
	firstPart := arrangement.InitPart("Part 1")
	return []arrangement.Part{firstPart}
}

func LoadFile(filename string, template string) (Definition, error) {
	var definition Definition
	var fileErr error
	if filename != "" {
		definition, fileErr = Read(filename)
		gridTemplate, exists := config.GetTemplate(definition.template)
		if exists {
			config.LongGates = config.GetGateLengths(gridTemplate.MaxGateLength)
		} else {
			maxGateLength := config.GetDefaultTemplate().MaxGateLength
			config.LongGates = config.GetGateLengths(maxGateLength)
		}
	}

	if filename == "" || fileErr != nil {
		newDefinition := InitDefinition(template, instrument)
		definition = newDefinition
	}

	return definition, fileErr
}

func InitModel(filename string, midiConnection seqmidi.MidiConnection, template string, instrument string, midiLoopMode MidiLoopMode, theme string) model {
	logFile, logFileErr := tea.LogToFile("debug.log", "debug")

	newCursor := cursor.New()
	newCursor.BlinkSpeed = 600 * time.Millisecond
	newCursor.Style = lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	definition, err := LoadFile(filename, template)

	programChannel := make(chan midiEventLoopMsg)
	lockReceiverChannel := make(chan bool)
	unlockReceiverChannel := make(chan bool)
	errorChannel := make(chan error)

	themes.ChooseTheme(theme)

	if err == nil {
		// No capacity for multiple errors currently
		err = logFileErr
	}

	return model{
		currentError:          err,
		theme:                 theme,
		filename:              filename,
		textInput:             InitTextInput(),
		partSelectorIndex:     -1,
		midiLoopMode:          midiLoopMode,
		programChannel:        programChannel,
		lockReceiverChannel:   lockReceiverChannel,
		unlockReceiverChannel: unlockReceiverChannel,
		errChan:               errorChannel,
		help:                  help.New(),
		cursor:                newCursor,
		midiConnection:        midiConnection,
		logFile:               logFile,
		logFileAvailable:      logFileErr == nil,
		gridCursor:            GK(0, 0),
		currentOverlay:        (*definition.parts)[0].Overlays,
		overlayKeyEdit:        overlaykey.InitModel(),
		arrangement:           arrangement.InitModel(definition.arrangement, definition.parts),
		definition:            definition,
		playState:             InitLineStates(len(definition.lines), []linestate{}, 0),
	}
}

func InitTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "------"
	ti.Prompt = ""
	ti.CharLimit = 20
	ti.Width = 20
	ti.Focus()
	return ti
}

func (m model) LogTeaMsg(msg tea.Msg) {
	switch msg := msg.(type) {
	case beatMsg:
		m.LogString(fmt.Sprintf("beatMsg %d %d\n", msg.interval, m.definition.tempo))
	case tea.KeyMsg:
		m.LogString(fmt.Sprintf("keyMsg %s\n", msg.String()))
	case cursor.BlinkMsg:
	default:
		m.LogString(fmt.Sprintf("%T\n", msg))
	}
}

func (m *model) LogString(message string) {
	if m.logFileAvailable {
		_, err := m.logFile.WriteString(message + "\n")
		if err != nil {
			m.logFileAvailable = false
			panic("could not write to log file")
		}
	}
}

func (m *model) LogError(err error) {
	if m.logFileAvailable {
		m.LogString("ERROR --------------------- \n")

		_, writeErr := fmt.Fprintf(m.logFile, "%+v", err)
		if writeErr != nil {
			m.logFileAvailable = false
			panic("could not write to log file")
		}

		m.LogString("\n")
	}
}

func (m *model) LogFromBeatTime() {
	if m.logFileAvailable {
		_, err := fmt.Fprintf(m.logFile, "%d\n", time.Since(m.beatTime))
		if err != nil {
			m.logFileAvailable = false
			panic("could not write to log file")
		}
	}
}

func RunProgram(filename string, midiConnection seqmidi.MidiConnection, template string, instrument string, midiLoopMode MidiLoopMode, theme string) *tea.Program {
	config.Init()
	model := InitModel(filename, midiConnection, template, instrument, midiLoopMode, theme)
	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithReportFocus())
	err := MidiEventLoop(midiLoopMode, model.lockReceiverChannel, model.unlockReceiverChannel, model.programChannel, program)
	if err != nil {
		go func() {
			program.Send(errorMsg{fault.Wrap(err, fmsg.With("could not setup midi event loop"))})
			for {
				// NOTE: Eat the messages sent on the program channel so that we can quit and still send the stopMsg
				<-model.programChannel
			}
		}()
		return program
	} else {
		ErrorLoop(program, model.errChan)
		model.SyncTempo()
		return program
	}
}

func ErrorLoop(program *tea.Program, errChan chan error) {
	go func() {
		for {
			programError := <-errChan
			program.Send(errorMsg{programError})
		}
	}()
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return tea.FocusMsg{} }
}

func Is(msg tea.KeyMsg, k ...key.Binding) bool {
	return key.Matches(msg, k...)
}

func IsNot(msg tea.KeyMsg, k ...key.Binding) bool {
	return !key.Matches(msg, k...)
}

var skipArrangementMappings = []mappings.Command{
	mappings.ToggleArrangementView,
	mappings.Undo,
	mappings.Redo,
	mappings.PlayLoop,
	mappings.PlayStop,
	mappings.PlayPart,
	mappings.ChangePart,
	mappings.NewSectionAfter,
	mappings.NewSectionBefore,
	mappings.NextSection,
	mappings.PrevSection,
	mappings.HoldingKeys,
}

func (m model) ShouldSkipArrangement(mapping mappings.Mapping) bool {
	return slices.Contains(skipArrangementMappings, mapping.Command) ||
		((m.selectionIndicator == operation.SelectPart || m.selectionIndicator == operation.SelectChangePart) && slices.Contains([]mappings.Command{mappings.Decrease, mappings.Increase}, mapping.Command))
}

type panicMsg struct {
	message    string
	stacktrace []byte
}

func (m model) Update(msg tea.Msg) (rModel tea.Model, rCmd tea.Cmd) {
	defer func() {
		if r := recover(); r != nil {
			rModel = m
			stackTrace := make([]byte, 4096)
			n := runtime.Stack(stackTrace, false)
			rCmd = func() tea.Msg {
				return panicMsg{message: fmt.Sprintf("Caught Update Panic: %v", r), stacktrace: stackTrace[:n]}
			}
		}
	}()

	// NOTE: Processing enter/escape keys before anything else
	switch msg := msg.(type) {
	case errorMsg:
		m.SetCurrentError(msg.error)
	case panicMsg:
		m.SetCurrentError(errors.New(msg.message))
		m.LogString(fmt.Sprintf(" ------ Panic Message ------- \n%s\n", msg.message))
		m.LogString(fmt.Sprintf(" ------ Stacktrace ---------- \n%s\n", msg.stacktrace))
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch m.focus {
			case operation.FocusOverlayKey:
				m.EnsureOverlay()
			}
			switch m.selectionIndicator {
			case operation.SelectRenamePart:
				m.RenamePart(m.textInput.Value())
				m.textInput.Reset()
				m.selectionIndicator = operation.SelectNothing
				return m, nil
			case operation.SelectFileName:
				if m.textInput.Value() != "" {
					m.filename = fmt.Sprintf("%s.seq", m.textInput.Value())
					m.textInput.Reset()
					m.selectionIndicator = operation.SelectNothing
					m.Save()
					return m, nil
				} else {
					m.selectionIndicator = operation.SelectNothing
					return m, nil
				}
			case operation.SelectPart:
				_, cmd := m.arrangement.Update(arrangement.NewPart{Index: m.partSelectorIndex, After: m.sectionSideIndicator == SectionAfter, IsPlaying: m.playing != PlayStopped})
				if m.sectionSideIndicator == SectionAfter {
					m.arrangement.Cursor.MoveNext()
				} else {
					m.arrangement.Cursor.MovePrev()
				}
				m.ResetCurrentOverlay()
				m.selectionIndicator = operation.SelectNothing
				return m, cmd
			case operation.SelectChangePart:
				_, cmd := m.arrangement.Update(arrangement.ChangePart{Index: m.partSelectorIndex})
				m.currentOverlay = m.CurrentPart().Overlays
				m.selectionIndicator = operation.SelectNothing
				return m, cmd
			case operation.SelectConfirmNew:
				m.NewSequence()
				m.selectionIndicator = operation.SelectNothing
			case operation.SelectConfirmReload:
				m.ReloadFile()
				m.selectionIndicator = operation.SelectNothing
			case operation.SelectConfirmQuit:
				m.programChannel <- quitMsg{}
				err := m.logFile.Close()
				if err != nil {
					// NOTE: no good way to display this error when quitting, just panic
					panic("Unable to close logfile")
				}
				return m, tea.Quit
			default:
				m.Escape()
			}
		case "esc":
			if m.selectionIndicator != operation.SelectNothing {
				m.textInput.Reset()
				m.Escape()
				return m, nil
			}
		}

		// NOTE: Process text input before mappings if we have a text input selected
		if m.selectionIndicator == operation.SelectRenamePart || m.selectionIndicator == operation.SelectFileName {
			tiModel, cmd := m.textInput.Update(msg)
			m.textInput = tiModel
			return m, cmd
		}

		// NOTE: Overlay key edit has it's own key bindings
		if m.focus == operation.FocusOverlayKey {
			okModel, cmd := m.overlayKeyEdit.Update(msg)
			m.overlayKeyEdit = okModel
			return m, cmd
		}

		mapping := mappings.ProcessKey(msg, m.definition.templateSequencerType, m.patternMode != operation.PatternFill)

		if mapping.Command == mappings.Quit {
			m.SetSelectionIndicator(operation.SelectConfirmQuit)
		}

		// NOTE: Arrangement editor has its own key bindings, but some mappings we still want to route to this model
		if m.focus == operation.FocusArrangementEditor && !m.ShouldSkipArrangement(mapping) {
			mappings.ResetKeycombo()
			arrangmementModel, cmd := m.arrangement.Update(msg)
			m.arrangement = arrangmementModel
			m.ResetCurrentOverlay()
			return m, cmd
		}

		// NOTE: Finally process the mapping
		switch mapping.Command {
		case mappings.ReloadFile:
			m.SetSelectionIndicator(operation.SelectConfirmReload)
		case mappings.HoldingKeys:
			return m, nil
		case mappings.CursorDown:
			if slices.Contains([]operation.Selection{operation.SelectNothing, operation.SelectSetupChannel, operation.SelectSetupMessageType, operation.SelectSetupValue}, m.selectionIndicator) {
				m.CursorDown()
				m.UnsetActiveChord()
			}
		case mappings.CursorUp:
			if slices.Contains([]operation.Selection{operation.SelectNothing, operation.SelectSetupChannel, operation.SelectSetupMessageType, operation.SelectSetupValue}, m.selectionIndicator) {
				m.CursorUp()
				m.UnsetActiveChord()
			}
		case mappings.CursorLeft:
			if m.selectionIndicator == operation.SelectRatchets {
				if m.ratchetCursor > 0 {
					m.ratchetCursor--
				}
			} else if m.selectionIndicator > 0 && m.selectionIndicator != operation.SelectSpecificValue {
				// Do Nothing
			} else {
				m.CursorLeft()
				m.UnsetActiveChord()
			}
		case mappings.CursorRight:
			if m.selectionIndicator == operation.SelectRatchets {
				currentNote, _ := m.CurrentNote()
				if m.ratchetCursor < currentNote.Ratchets.Length {
					m.ratchetCursor++
				}
			} else if m.selectionIndicator > 0 && m.selectionIndicator != operation.SelectSpecificValue {
				// Do Nothing
			} else {
				m.CursorRight()
				m.UnsetActiveChord()
			}
		case mappings.CursorLineStart:
			m.SetGridCursor(gridKey{
				Line: m.gridCursor.Line,
				Beat: 0,
			})
		case mappings.CursorLineEnd:
			m.SetGridCursor(gridKey{
				Line: m.gridCursor.Line,
				Beat: m.CurrentPart().Beats - 1,
			})
		case mappings.CursorLastLine:
			m.SetGridCursor(gridKey{
				Line: uint8(len(m.definition.lines) - 1),
				Beat: m.gridCursor.Beat,
			})
		case mappings.CursorFirstLine:
			m.SetGridCursor(gridKey{
				Line: 0,
				Beat: m.gridCursor.Beat,
			})
		case mappings.Escape:
			m.Escape()
		case mappings.PlayStop:
			if m.playing == PlayStopped {
				m.loopMode = LoopSong
			}
			m.StartStop()
		case mappings.PlayPart:
			if m.playing == PlayStopped {
				m.loopMode = LoopPart
			}
			m.StartStop()
		case mappings.PlayLoop:
			if m.playing == PlayStopped {
				m.loopMode = LoopSong
			}
			m.arrangement.Root.SetInfinite()
			m.StartStop()
		case mappings.PlayOverlayLoop:
			if m.playing == PlayStopped {
				m.loopMode = LoopOverlay
			}
			m.StartStop()
			m.SetPlayCycles(m.currentOverlay.Key.GetMinimumKeyCycle())
		case mappings.PlayRecord:
			if m.playing == PlayStopped {
				err := seqmidi.SendRecordMessage()
				if err != nil {
					m.SetCurrentError(err)
				} else {
					m.StartStop()
				}
			} else {
				m.StartStop()
			}
		case mappings.OverlayInputSwitch:
			// NOTE: This component handles getting into the overlay key edit mode
			// the overlaykey component handles getting out of it
			m.SetSelectionIndicator(operation.SelectNothing)
			m.focus = operation.FocusOverlayKey
			m.overlayKeyEdit.Focus(true)
		case mappings.TempoInputSwitch:
			states := []operation.Selection{operation.SelectNothing, operation.SelectTempo, operation.SelectTempoSubdivision}
			m.SetSelectionIndicator(AdvanceSelectionState(states, m.selectionIndicator))
		case mappings.SetupInputSwitch:
			var states []operation.Selection
			if m.definition.lines[m.gridCursor.Line].MsgType == grid.MessageTypeProgramChange {
				states = []operation.Selection{operation.SelectNothing, operation.SelectSetupChannel, operation.SelectSetupMessageType}
			} else {
				states = []operation.Selection{operation.SelectNothing, operation.SelectSetupChannel, operation.SelectSetupMessageType, operation.SelectSetupValue}
			}
			m.SetSelectionIndicator(AdvanceSelectionState(states, m.selectionIndicator))
		case mappings.AccentInputSwitch:
			states := []operation.Selection{operation.SelectNothing, operation.SelectAccentDiff, operation.SelectAccentTarget, operation.SelectAccentStart}
			m.SetSelectionIndicator(AdvanceSelectionState(states, m.selectionIndicator))
		case mappings.RatchetInputSwitch:
			currentNote, _ := m.CurrentNote()
			if currentNote.AccentIndex > 0 {
				states := []operation.Selection{operation.SelectNothing, operation.SelectRatchets, operation.SelectRatchetSpan}
				m.SetSelectionIndicator(AdvanceSelectionState(states, m.selectionIndicator))
				m.ratchetCursor = 0
			}
		case mappings.BeatInputSwitch:
			states := []operation.Selection{operation.SelectNothing, operation.SelectBeats, operation.SelectStartBeats}
			m.SetSelectionIndicator(AdvanceSelectionState(states, m.selectionIndicator))
		case mappings.CyclesInputSwitch:
			states := []operation.Selection{operation.SelectNothing, operation.SelectCycles, operation.SelectStartCycles}
			m.SetSelectionIndicator(AdvanceSelectionState(states, m.selectionIndicator))
		case mappings.ToggleArrangementView:
			m.showArrangementView = !m.showArrangementView || m.focus != operation.FocusArrangementEditor
			if m.showArrangementView {
				m.SetSelectionIndicator(operation.SelectNothing)
				m.focus = operation.FocusArrangementEditor
				m.arrangement.Focus = true
			} else {
				m.Escape()
				m.focus = operation.FocusGrid
				model, cmd := m.arrangement.Update(tea.KeyMsg{Type: tea.KeyEsc})
				m.arrangement = model
				return m, cmd
			}
		case mappings.Increase:
			switch m.selectionIndicator {
			case operation.SelectTempo:
				m.IncreaseTempo(1)
			case operation.SelectTempoSubdivision:
				if m.definition.subdivisions < 8 {
					m.definition.subdivisions++
				}
				m.SyncTempo()
			case operation.SelectSetupChannel:
				m.definition.lines[m.gridCursor.Line].IncrementChannel()
			case operation.SelectSetupMessageType:
				m.definition.lines[m.gridCursor.Line].IncrementMessageType()
			case operation.SelectSetupValue:
				switch m.definition.lines[m.gridCursor.Line].MsgType {
				case grid.MessageTypeNote:
					m.definition.lines[m.gridCursor.Line].IncrementNote()
				case grid.MessageTypeCc:
					m.IncrementCC()
				}
			case operation.SelectRatchetSpan:
				m.IncreaseSpan()
			case operation.SelectAccentDiff:
				m.IncreaseAccentDiff()
			case operation.SelectAccentTarget:
				// Only two options right now, so increase and decrease would do the
				// same thing
				m.DecreaseAccentTarget()
			case operation.SelectAccentStart:
				m.IncreaseAccentStart()
			case operation.SelectBeats:
				m.IncreaseBeats()
			case operation.SelectCycles:
				m.IncreaseCycles()
			case operation.SelectStartBeats:
				m.IncreaseStartBeats()
			case operation.SelectStartCycles:
				m.IncreaseStartCycles()
			case operation.SelectPart:
				m.IncreasePartSelector()
			case operation.SelectChangePart:
				m.IncreasePartSelector()
			case operation.SelectSpecificValue:
				note, _ := m.CurrentNote()
				m.IncrementSpecificValue(note)
			default:
				m.IncreaseTempo(5)
			}
		case mappings.Decrease:
			switch m.selectionIndicator {
			case operation.SelectTempo:
				m.DecreaseTempo(1)
			case operation.SelectTempoSubdivision:
				if m.definition.subdivisions > 1 {
					m.definition.subdivisions--
				}
				m.SyncTempo()
			case operation.SelectSetupChannel:
				m.definition.lines[m.gridCursor.Line].DecrementChannel()
			case operation.SelectSetupMessageType:
				m.definition.lines[m.gridCursor.Line].DecrementMessageType()
			case operation.SelectSetupValue:
				switch m.definition.lines[m.gridCursor.Line].MsgType {
				case grid.MessageTypeNote:
					m.definition.lines[m.gridCursor.Line].DecrementNote()
				case grid.MessageTypeCc:
					m.DecrementCC()
				}
			case operation.SelectRatchetSpan:
				m.DecreaseSpan()
			case operation.SelectAccentDiff:
				m.DecreaseAccentDiff()
			case operation.SelectAccentTarget:
				m.DecreaseAccentTarget()
			case operation.SelectAccentStart:
				m.DecreaseAccentStart()
			case operation.SelectBeats:
				m.DecreaseBeats()
			case operation.SelectCycles:
				m.DecreaseCycles()
			case operation.SelectStartBeats:
				m.DecreaseStartBeats()
			case operation.SelectStartCycles:
				m.DecreaseStartCycles()
			case operation.SelectPart:
				m.DecreasePartSelector()
			case operation.SelectChangePart:
				m.DecreasePartSelector()
			case operation.SelectSpecificValue:
				note, _ := m.CurrentNote()
				m.DecrementSpecificValue(note)
			default:
				m.DecreaseTempo(5)
			}
		case mappings.ToggleGateMode:
			m.SetPatternMode(operation.PatternGate)
		case mappings.ToggleWaitMode:
			m.SetPatternMode(operation.PatternWait)
		case mappings.ToggleAccentMode:
			m.SetPatternMode(operation.PatternAccent)
		case mappings.ToggleRatchetMode:
			m.SetPatternMode(operation.PatternRatchet)
		case mappings.ToggleChordMode:
			if m.definition.templateSequencerType == operation.SeqModeChord {
				m.definition.templateSequencerType = operation.SeqModeLine
			} else {
				m.definition.templateSequencerType = operation.SeqModeChord
			}
		case mappings.PrevOverlay:
			m.NextOverlay(-1)
			m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
		case mappings.NextOverlay:
			m.NextOverlay(+1)
			m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
		case mappings.Save:
			if m.filename == "" {
				m.selectionIndicator = operation.SelectFileName
			} else {
				m.Save()
			}
		case mappings.Undo:
			tempo, subdiv := m.definition.tempo, m.definition.subdivisions
			undoStack := m.Undo()
			if undoStack != EmptyStack {
				m.PushRedo(undoStack)
			}
			if tempo != m.definition.tempo || subdiv != m.definition.subdivisions {
				m.SyncTempo()
			}
		case mappings.Redo:
			tempo, subdiv := m.definition.tempo, m.definition.subdivisions
			undoStack := m.Redo()
			if undoStack != EmptyStack {
				m.PushUndo(undoStack)
			}
			if tempo != m.definition.tempo || subdiv != m.definition.subdivisions {
				m.SyncTempo()
			}
		case mappings.New:
			m.selectionIndicator = operation.SelectConfirmNew
		case mappings.ToggleVisualMode:
			m.visualAnchorCursor = m.gridCursor
			m.visualMode = !m.visualMode
		case mappings.TogglePlayEdit:
			m.playEditing = !m.playEditing
		case mappings.NewLine:
			if len(m.definition.lines) < 100 {
				lastline := m.definition.lines[len(m.definition.lines)-1]
				m.definition.lines = append(m.definition.lines, grid.LineDefinition{
					Channel: lastline.Channel,
					Note:    lastline.Note + 1,
				})
				if m.playing != PlayStopped {
					m.playState = append(m.playState, InitLineState(PlayStatePlay, uint8(len(m.definition.lines)-1), 0))
				}
			}
		case mappings.NewSectionAfter:
			m.SetSelectionIndicator(operation.SelectPart)
			m.sectionSideIndicator = SectionAfter
		case mappings.NewSectionBefore:
			m.SetSelectionIndicator(operation.SelectPart)
			m.sectionSideIndicator = SectionBefore
		case mappings.ChangePart:
			m.SetSelectionIndicator(operation.SelectChangePart)
		case mappings.NextTheme:
			m.NextTheme()
		case mappings.PrevTheme:
			m.PrevTheme()
		case mappings.NextSection:
			m.NextSection()
		case mappings.PrevSection:
			m.PrevSection()
		case mappings.Yank:
			m.yankBuffer = m.Yank()
			m.SetGridCursor(m.YankBounds().TopLeft())
			m.visualMode = false
		case mappings.Mute:
			if m.IsRatchetSelector() {
				m.ToggleRatchetMute()
			} else {
				m.playState = Mute(m.playState, m.gridCursor.Line)
				m.hasSolo = m.HasSolo()
			}
		case mappings.Solo:
			m.playState = Solo(m.playState, m.gridCursor.Line)
			m.hasSolo = m.HasSolo()
		case mappings.Enter:
			// NOTE: Do nothing, we've reacted to Enter at the beginning of Update
		default:
			m = m.UpdateDefinition(mapping)
		}
	case tea.FocusMsg:
		m.hasUIFocus = true
	case tea.BlurMsg:
		m.hasUIFocus = false
	case overlaykey.UpdatedOverlayKey:
		if !msg.HasFocus {
			m.focus = operation.FocusGrid
			m.selectionIndicator = operation.SelectNothing
		}
	case uiStartMsg:
		if m.playing == PlayStopped {
			m.playing = PlayReceiver
		} else {
			m.SetCurrentError(errors.New("cannot start when already started"))
		}
		m.Start()
	case uiStopMsg:
		m.playing = PlayStopped
		m.Stop()
	case uiConnectedMsg:
		m.connected = true
	case uiNotConnectedMsg:
		m.connected = false
	case beatMsg:
		m.beatTime = time.Now()
		playingOverlay := m.CurrentPart().Overlays.HighestMatchingOverlay(m.CurrentSongSection().PlayCycles())
		if m.playing != PlayStopped {
			m.advanceCurrentBeat(playingOverlay)
			m.advanceKeyCycle()
		}
		if m.playing != PlayStopped {
			playingOverlay := m.CurrentPart().Overlays.HighestMatchingOverlay(m.CurrentSongSection().PlayCycles())
			gridKeys := make([]grid.GridKey, 0, len(m.playState))
			m.CurrentBeatGridKeys(&gridKeys)
			pattern := make(grid.Pattern)
			playingOverlay.CurrentBeatOverlayPattern(&pattern, m.CurrentSongSection().PlayCycles(), gridKeys)
			err := m.PlayBeat(msg.interval, pattern)
			if err != nil {
				m.SetCurrentError(fault.Wrap(err, fmsg.With("error when playing beat")))
			}
		}
	case arrangement.GiveBackFocus:
		m.selectionIndicator = operation.SelectNothing
		m.focus = operation.FocusGrid
	case arrangement.RenamePart:
		m.SetSelectionIndicator(operation.SelectRenamePart)
	case arrangement.Undo:
		m.PushArrUndo(msg)
	}
	var cmd tea.Cmd
	cursor, cmd := m.cursor.Update(msg)
	m.cursor = cursor
	return m, cmd
}

func (m *model) CursorDown() {
	if m.gridCursor.Line < uint8(len(m.definition.lines)-1) {
		m.SetGridCursor(gridKey{
			Line: m.gridCursor.Line + 1,
			Beat: m.gridCursor.Beat,
		})
	}
}

func (m *model) CursorUp() {
	if m.gridCursor.Line > 0 {
		m.SetGridCursor(gridKey{
			Line: m.gridCursor.Line - 1,
			Beat: m.gridCursor.Beat,
		})
	}
}

func (m *model) CursorRight() {
	if m.gridCursor.Beat < m.CurrentPart().Beats-1 {
		m.SetGridCursor(gridKey{
			Line: m.gridCursor.Line,
			Beat: m.gridCursor.Beat + 1,
		})
	}
}

func (m *model) CursorLeft() {
	if m.gridCursor.Beat > 0 {
		m.SetGridCursor(gridKey{
			Line: m.gridCursor.Line,
			Beat: m.gridCursor.Beat - 1,
		})
	}
}

func (m *model) SetSelectionIndicator(desiredIndicator operation.Selection) {
	if desiredIndicator != operation.SelectNothing && m.selectionIndicator == operation.SelectNothing {
		m.CaptureTemporaryState()
	}
	if desiredIndicator == operation.SelectNothing && m.selectionIndicator != operation.SelectNothing {
		diff := createStateDiff(&m.definition, &m.temporaryState)
		if diff.Changed() {
			m.PushUndoables(UndoStateDiff{diff}, UndoStateDiff{diff.Reverse()})
		}
		if m.CurrentPart().Beats != m.temporaryState.beats {
			m.PushUndoables(UndoBeats{m.temporaryState.beats, m.arrangement.Cursor}, UndoBeats{m.CurrentPart().Beats, m.arrangement.Cursor})
		}
	}
	m.patternMode = operation.PatternFill
	m.selectionIndicator = desiredIndicator
}

func (m *model) CaptureTemporaryState() {
	linesCopy := make([]grid.LineDefinition, len(m.definition.lines))
	for i, defLine := range m.definition.lines {
		newLine := grid.LineDefinition{
			Channel: defLine.Channel,
			Note:    defLine.Note,
			MsgType: defLine.MsgType,
			Name:    defLine.Name,
		}
		linesCopy[i] = newLine
	}
	accentDataCopy := make([]config.Accent, len(m.definition.accents.Data))
	copy(accentDataCopy, m.definition.accents.Data)

	m.temporaryState = temporaryState{
		lines:        linesCopy,
		tempo:        m.definition.tempo,
		subdivisions: m.definition.subdivisions,
		accents: patternAccents{
			Diff:   m.definition.accents.Diff,
			Start:  m.definition.accents.Start,
			Target: m.definition.accents.Target,
			Data:   accentDataCopy,
		},
		beats:  m.CurrentPart().Beats,
		active: true,
	}
}

func (m *model) SetPatternMode(mode operation.PatternMode) {
	m.patternMode = mode
	m.selectionIndicator = operation.SelectNothing
}

func IsDefinitionChangeSelection(indicator operation.Selection) bool {
	return slices.Contains([]operation.Selection{
		operation.SelectTempo,
		operation.SelectTempoSubdivision,
		operation.SelectSetupChannel,
		operation.SelectSetupMessageType,
		operation.SelectSetupValue,
		operation.SelectAccentDiff,
		operation.SelectAccentTarget,
		operation.SelectAccentStart,
	}, indicator)
}

func (m *model) Escape() {
	if IsPropertyIndicator(m.selectionIndicator) {
		m.RecordPropertyUndo()
	}
	if m.selectionIndicator == operation.SelectNothing && m.patternMode == operation.PatternFill {
		m.visualMode = false
	}
	m.patternMode = operation.PatternFill
	m.selectionIndicator = operation.SelectNothing
}

func IsPropertyIndicator(selection operation.Selection) bool {
	return slices.Contains([]operation.Selection{
		operation.SelectTempo,
		operation.SelectTempoSubdivision,
		operation.SelectSetupChannel,
		operation.SelectSetupMessageType,
		operation.SelectSetupValue,
		operation.SelectAccentDiff,
		operation.SelectAccentTarget,
		operation.SelectAccentStart,
		operation.SelectBeats,
		operation.SelectSpecificValue,
	}, selection)
}

func (m *model) RecordPropertyUndo() {
	if m.temporaryState.active {
		diff := createStateDiff(&m.definition, &m.temporaryState)
		if diff.Changed() {
			m.PushUndoables(UndoStateDiff{diff}, UndoStateDiff{diff.Reverse()})
		}
		if m.CurrentPart().Beats != m.temporaryState.beats {
			m.PushUndoables(UndoBeats{m.temporaryState.beats, m.arrangement.Cursor}, UndoBeats{m.CurrentPart().Beats, m.arrangement.Cursor})
		}
		m.temporaryState = temporaryState{}
	}
}

func (m *model) NewSequence() {
	m.filename = ""
	m.definition = InitDefinition(m.definition.template, m.definition.instrument)
	m.arrangement = arrangement.InitModel(m.definition.arrangement, m.definition.parts)
	m.SetGridCursor(GK(0, 0))
	m.ResetCurrentOverlay()
}

func (m model) NeedsWrite() bool {
	return m.needsWrite != m.undoStack.id
}

func (m *model) NextTheme() {
	index := slices.Index(themes.Themes, m.theme)
	if index+1 < len(themes.Themes) {
		m.theme = themes.Themes[index+1]
	} else {
		m.theme = themes.Themes[0]
	}
	themes.ChooseTheme(m.theme)
}

func (m *model) PrevTheme() {
	index := slices.Index(themes.Themes, m.theme)
	if index-1 >= 0 {
		m.theme = themes.Themes[index-1]
	} else {
		m.theme = themes.Themes[len(themes.Themes)-1]
	}
	themes.ChooseTheme(m.theme)
}

func (m *model) NextSection() {
	if m.arrangement.Cursor.MoveNext() {
		m.ResetCurrentOverlay()
		m.arrangement.ResetDepth()
	}
}

func (m *model) PrevSection() {
	if m.arrangement.Cursor.MovePrev() {
		m.ResetCurrentOverlay()
		m.arrangement.ResetDepth()
	}
}

func (m *model) Start() {
	if !m.midiConnection.IsReady() {
		err := m.midiConnection.ConnectAndOpen()
		if err != nil {
			m.SetCurrentError(fault.Wrap(err, fmsg.With("cannot open midi connection")))
			m.playing = PlayStopped
			return
		}
	}

	switch m.loopMode {
	case LoopSong:
		m.arrangement.Cursor = arrangement.ArrCursor{m.definition.arrangement}
		m.arrangement.Cursor.MoveNext()
		m.arrangement.Cursor.ResetIterations()
	case LoopPart:
		m.arrangement.SetCurrentNodeInfinite()
	}

	m.arrangement.Root.ResetAllPlayCycles()
	section := m.CurrentSongSection()
	section.ResetPlayCycles()
	m.playState = InitLineStates(len(m.definition.lines), m.playState, uint8(section.StartBeat))

	playingOverlay := m.CurrentPart().Overlays.HighestMatchingOverlay(section.PlayCycles())
	tickInterval := m.TickInterval()

	pattern := m.CombinedBeatPattern(playingOverlay)
	err := m.PlayBeat(tickInterval, pattern)
	if err != nil {
		m.SetCurrentError(fault.Wrap(err, fmsg.With("cannot play first beat")))
	}
	if m.playing == PlayStandard {
		m.programChannel <- startMsg{tempo: m.definition.tempo, subdivisions: m.definition.subdivisions}
	}
}

func (m *model) Stop() {
	if m.loopMode == LoopPart {
		m.arrangement.ResetDepth()
	}
	m.arrangement.Root.ResetCycles()
	m.arrangement.Root.ResetAllPlayCycles()
	m.arrangement.Root.ResetIterations()
	m.ResetCurrentOverlay()

	notes := notereg.Clear()
	sendFn, err := m.midiConnection.AcquireSendFunc()
	if err != nil {
		m.SetCurrentError(fault.Wrap(err))
	}
	for _, n := range notes {
		switch n := n.(type) {
		case noteMsg:
			PlayMessage(time.Duration(0), n.OffMessage(), sendFn, m.errChan)
		}
	}
}

func (m *model) StartStop() {
	m.playEditing = false
	if m.playing == PlayStopped {
		m.playing = PlayStandard
		if m.midiLoopMode == MlmReceiver {
			// NOTE: When instance is receiver, allow it to play alone and lock out transmitter messages
			m.lockReceiverChannel <- true
		}
		m.Start()
	} else {
		if m.playing == PlayStandard {
			m.programChannel <- stopMsg{}
			if m.midiLoopMode == MlmReceiver {
				// NOTE: Allow transmitter messages
				m.unlockReceiverChannel <- true
			}
		}
		m.playing = PlayStopped
		m.Stop()
	}
}

func (m model) ProcessNoteMsg(msg Delayable) error {
	// TODO: We don't need the redirection, we know what each type of note is when this function is called
	sendFn, err := m.midiConnection.AcquireSendFunc()
	if err != nil {
		return fault.Wrap(err, fmsg.With("could not acquire send func"))
	}
	switch msg := msg.(type) {
	case noteMsg:
		switch msg.midiType {
		case midi.NoteOnMsg:
			if notereg.Has(msg) {
				notereg.Remove(msg)
				PlayMessage(0, msg.OffMessage(), sendFn, m.errChan)
			}
			if err := notereg.Add(msg); err != nil {
				return fault.Wrap(err, fmsg.With("added a note already in registry"))
			}
			PlayMessage(msg.delay, msg.GetOnMidi(), sendFn, m.errChan)
		case midi.NoteOffMsg:
			PlayOffMessage(msg, sendFn, m.errChan)
		}
	case controlChangeMsg:
		PlayMessage(msg.delay, msg.MidiMessage(), sendFn, m.errChan)
	case programChangeMsg:
		PlayMessage(msg.delay, msg.MidiMessage(), sendFn, m.errChan)
	}
	return nil
}

func AdvanceSelectionState(states []operation.Selection, currentSelection operation.Selection) operation.Selection {
	index := slices.Index(states, currentSelection)
	var resultSelection operation.Selection
	if index < 0 {
		indexNothing := slices.Index(states, operation.SelectNothing)
		resultSelection = states[uint8(indexNothing+1)%uint8(len(states))]
	} else {
		resultSelection = states[uint8(index+1)%uint8(len(states))]
	}
	return resultSelection
}

func (m model) UpdateDefinitionKeys(mapping mappings.Mapping) model {
	if !m.activeChord.HasValue() {
		chord, exists := m.currentOverlay.FindChord(m.gridCursor)
		if exists {
			m.activeChord = chord
		}
	}
	switch mapping.Command {
	case mappings.NoteAdd:
		m.AddNote()
	case mappings.NoteRemove:
		m.yankBuffer = m.Yank()
		m.RemoveNote()
		m.visualMode = false
	case mappings.AccentIncrease:
		m.AccentModify(1)
	case mappings.AccentDecrease:
		m.AccentModify(-1)
	case mappings.GateIncrease:
		m.GateModify(1)
	case mappings.GateDecrease:
		m.GateModify(-1)
	case mappings.GateBigIncrease:
		m.GateModify(8)
	case mappings.GateBigDecrease:
		m.GateModify(-8)
	case mappings.WaitIncrease:
		m.WaitModify(1)
	case mappings.WaitDecrease:
		m.WaitModify(-1)
	case mappings.OverlayNoteRemove:
		m.OverlayRemoveNote()
	case mappings.ClearLine:
		if m.visualMode {
			m.RemoveNote()
		} else {
			m.ClearOverlayLine()
		}
	case mappings.RatchetIncrease:
		m.RatchetModify(1)
	case mappings.RatchetDecrease:
		m.RatchetModify(-1)
		m.EnsureRatchetCursorVisible()
	case mappings.ActionAddLineReset:
		m.AddAction(grid.ActionLineReset)
	case mappings.ActionAddLineReverse:
		m.AddAction(grid.ActionLineReverse)
	case mappings.ActionAddSkipBeat:
		m.AddAction(grid.ActionLineSkipBeat)
	case mappings.ActionAddSkipBeatAll:
		m.AddAction(grid.ActionLineSkipBeatAll)
	case mappings.ActionAddLineResetAll:
		m.AddAction(grid.ActionLineResetAll)
	case mappings.ActionAddLineBounce:
		m.AddAction(grid.ActionLineBounce)
	case mappings.ActionAddLineBounceAll:
		m.AddAction(grid.ActionLineBounceAll)
	case mappings.ActionAddLineDelay:
		m.AddAction(grid.ActionLineDelay)
	case mappings.ActionAddSpecificValue:
		if m.definition.lines[m.gridCursor.Line].MsgType != grid.MessageTypeNote {
			m.AddAction(grid.ActionSpecificValue)
			m.SetSelectionIndicator(operation.SelectSpecificValue)
		}
	case mappings.SelectKeyLine:
		m.definition.keyline = m.gridCursor.Line
	case mappings.OverlayStackToggle:
		m.currentOverlay.ToggleOverlayStackOptions()
	case mappings.ClearOverlay:
		m.ClearOverlay()
	case mappings.RotateRight:
		switch m.definition.templateSequencerType {
		case operation.SeqModeLine:
			m.RotateRight()
		case operation.SeqModeChord:
			m.EnsureChord()
			m.MoveChordRight()
			m.CursorRight()
		}
	case mappings.RotateLeft:
		switch m.definition.templateSequencerType {
		case operation.SeqModeLine:
			m.RotateLeft()
		case operation.SeqModeChord:
			m.EnsureChord()
			m.MoveChordLeft()
			m.CursorLeft()
		}
	case mappings.RotateUp:
		switch m.definition.templateSequencerType {
		case operation.SeqModeLine:
			m.RotateUp()
		case operation.SeqModeChord:
			m.EnsureChord()
			m.MoveChordUp()
			m.CursorUp()
		}
	case mappings.RotateDown:
		switch m.definition.templateSequencerType {
		case operation.SeqModeLine:
			m.RotateDown()
		case operation.SeqModeChord:
			m.EnsureChord()
			m.MoveChordDown()
			m.CursorDown()
		}
	case mappings.Paste:
		m.Paste()
	case mappings.MajorTriad:
		m.EnsureChord()
		m.ChordChange(theory.MajorTriad)
	case mappings.MinorTriad:
		m.EnsureChord()
		m.ChordChange(theory.MinorTriad)
	case mappings.AugmentedTriad:
		m.EnsureChord()
		m.ChordChange(theory.AugmentedTriad)
	case mappings.DiminishedTriad:
		m.EnsureChord()
		m.ChordChange(theory.DiminishedTriad)
	case mappings.MajorSeventh:
		m.EnsureChord()
		m.ChordChange(theory.MajorSeventh)
	case mappings.MinorSeventh:
		m.EnsureChord()
		m.ChordChange(theory.MinorSeventh)
	case mappings.AugFifth:
		m.EnsureChord()
		m.ChordChange(theory.Aug5)
	case mappings.DimFifth:
		m.EnsureChord()
		m.ChordChange(theory.Dim5)
	case mappings.PerfectFifth:
		m.EnsureChord()
		m.ChordChange(theory.Perfect5)
	case mappings.MinorSecond:
		m.EnsureChord()
		m.ChordChange(theory.Minor2)
	case mappings.MajorSecond:
		m.EnsureChord()
		m.ChordChange(theory.Major2)
	case mappings.MinorThird:
		m.EnsureChord()
		m.ChordChange(theory.Minor3)
	case mappings.MajorThird:
		m.EnsureChord()
		m.ChordChange(theory.Major3)
	case mappings.PerfectFourth:
		m.EnsureChord()
		m.ChordChange(theory.Perfect4)
	case mappings.MajorSixth:
		m.EnsureChord()
		m.ChordChange(theory.Major6)
	case mappings.Octave:
		m.EnsureChord()
		m.ChordChange(theory.Octave)
	case mappings.MinorNinth:
		m.EnsureChord()
		m.ChordChange(theory.Minor9)
	case mappings.MajorNinth:
		m.EnsureChord()
		m.ChordChange(theory.Major9)
	case mappings.IncreaseInversions:
		m.EnsureChord()
		m.NextInversion()
	case mappings.DecreaseInversions:
		m.EnsureChord()
		m.PreviousInversion()
	case mappings.NextArpeggio:
		m.EnsureChord()
		m.NextArpeggio()
	case mappings.PrevArpeggio:
		m.EnsureChord()
		m.PrevArpeggio()
	case mappings.NextDouble:
		m.EnsureChord()
		m.NextDouble()
	case mappings.PrevDouble:
		m.EnsureChord()
		m.PrevDouble()
	case mappings.OmitRoot:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitRoot)
	case mappings.OmitSecond:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitSecond)
	case mappings.OmitThird:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitThird)
	case mappings.OmitFourth:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitFourth)
	case mappings.OmitFifth:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitFifth)
	case mappings.OmitSixth:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitSixth)
	case mappings.OmitSeventh:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitSeventh)
	case mappings.OmitOctave:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitOctave)
	case mappings.OmitNinth:
		m.EnsureChord()
		m.OmitChordNote(theory.OmitNinth)
	case mappings.RemoveChord:
		m.RemoveChord()
	case mappings.ConvertToNotes:
		m.ConvertChordToNotes()
	}

	if mapping.LastValue >= "1" && mapping.LastValue <= "9" {
		beatInterval, _ := strconv.ParseInt(mapping.LastValue, 0, 8)
		switch m.patternMode {
		case operation.PatternFill:
			m.fill(uint8(beatInterval))
		case operation.PatternAccent:
			m.incrementAccent(uint8(beatInterval), -1)
		case operation.PatternGate:
			m.incrementGate(uint8(beatInterval), -1)
		case operation.PatternRatchet:
			m.incrementRatchet(uint8(beatInterval), -1)
		case operation.PatternWait:
			m.incrementWait(uint8(beatInterval), -1)
		}
	}

	if IsShiftSymbol(mapping.LastValue) {
		beatInterval := convertSymbolToInt(mapping.LastValue)
		switch m.patternMode {
		case operation.PatternFill:
			m.fill(uint8(beatInterval))
		case operation.PatternAccent:
			m.incrementAccent(uint8(beatInterval), 1)
		case operation.PatternGate:
			m.incrementGate(uint8(beatInterval), 1)
		case operation.PatternRatchet:
			m.incrementRatchet(uint8(beatInterval), 1)
		case operation.PatternWait:
			m.incrementWait(uint8(beatInterval), 1)
		}
	}

	return m
}

func (m *model) EnsureChord() {
	overlayChord := m.CurrentChord()
	if overlayChord.HasValue() {
		if !overlayChord.BelongsTo(m.currentOverlay) {
			newGridChord := m.currentOverlay.SetChord(overlayChord.GridChord)
			m.activeChord = overlays.OverlayChord{GridChord: newGridChord}
		}
	}
}

func (m model) CurrentChord() overlays.OverlayChord {
	if m.definition.templateSequencerType == operation.SeqModeLine {
		return overlays.OverlayChord{}
	}
	overlayChord, exists := m.currentOverlay.FindChord(m.gridCursor)
	if exists {
		return overlayChord
	} else if m.activeChord.GridChord != nil {
		return m.activeChord
	} else {
		return overlays.OverlayChord{}
	}
}

func (m *model) ChordChange(alteration uint32) {
	overlayChord := m.CurrentChord()
	if overlayChord.GridChord == nil {
		m.currentOverlay.CreateChord(m.gridCursor, alteration)
	} else {
		overlayChord.GridChord.ApplyAlteration(alteration)
	}
}

func (m *model) RemoveChord() {
	if m.activeChord.HasValue() {
		m.currentOverlay.RemoveChord(m.activeChord)
		m.UnsetActiveChord()
	}
}

func (m *model) ConvertChordToNotes() {
	if m.activeChord.HasValue() {
		pattern := make(grid.Pattern)
		m.activeChord.GridChord.ArpeggiatedPattern(&pattern)
		m.currentOverlay.RemoveChord(m.activeChord)
		for gk, note := range pattern {
			m.currentOverlay.SetNote(gk, note)
		}
		m.activeChord = overlays.OverlayChord{}
	}
}

func (m *model) OmitChordNote(omission uint32) {
	if m.activeChord.HasValue() {
		m.activeChord.GridChord.Chord.OmitNote(omission)
	}
}

func (m *model) NextInversion() {
	if m.activeChord.HasValue() {
		m.activeChord.GridChord.Chord.NextInversion()
	}
}

func (m *model) PreviousInversion() {
	if m.activeChord.HasValue() {
		m.activeChord.GridChord.Chord.PreviousInversion()
	}
}

func (m *model) NextArpeggio() {
	if m.activeChord.HasValue() {
		m.activeChord.GridChord.NextArp()
	}
}

func (m *model) PrevArpeggio() {
	if m.activeChord.HasValue() {
		m.activeChord.GridChord.PrevArp()
	}
}

func (m *model) NextDouble() {
	if m.activeChord.HasValue() {
		m.activeChord.GridChord.NextDouble()
	}
}

func (m *model) PrevDouble() {
	if m.activeChord.HasValue() {
		m.activeChord.GridChord.PrevDouble()
	}
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

func (m model) UpdateDefinition(mapping mappings.Mapping) model {

	deepCopy := overlays.DeepCopy(m.currentOverlay)
	m.EnsureOverlay()
	if m.playing != PlayStopped && !m.playEditing {
		m.playEditing = true
		playingOverlay := m.CurrentPart().Overlays.HighestMatchingOverlay(m.CurrentSongSection().PlayCycles())
		m.currentOverlay = playingOverlay
		m.overlayKeyEdit.SetOverlayKey(playingOverlay.Key)
	}
	m = m.UpdateDefinitionKeys(mapping)
	undoable := m.UndoableOverlay(m.currentOverlay, deepCopy)
	redoable := m.UndoableOverlay(deepCopy, m.currentOverlay)

	if !undoable.overlayDiff.IsEmpty() {
		m.PushUndoables(undoable, redoable)
		m.ResetRedo()
	}

	return m
}

func (m model) UndoableOverlay(overlayA, overlayB *overlays.Overlay) UndoOverlayDiff {
	diff := overlays.DiffOverlays(overlayA, overlayB)
	return UndoOverlayDiff{m.currentOverlay.Key, m.gridCursor, m.arrangement.Cursor, diff}
}

func (m *model) Save() {
	err := Write(m, m.filename)
	if err != nil {
		m.SetCurrentError(fault.Wrap(err, fmsg.With("cannot write file")))
	}
	m.needsWrite = m.undoStack.id
}

func (m model) CurrentPart() arrangement.Part {
	section := m.CurrentSongSection()
	partID := section.Part
	return (*m.definition.parts)[partID]
}

func (m model) RenamePart(value string) {
	section := m.CurrentSongSection()
	partID := section.Part
	(*m.definition.parts)[partID].Name = value
}

func (m model) CurrentPartID() int {
	section := m.CurrentSongSection()
	return section.Part
}

func (m model) CurrentSongSection() arrangement.SongSection {
	currentNode := m.arrangement.Cursor.GetCurrentNode()
	if currentNode != nil && currentNode.IsEndNode() {
		return currentNode.Section
	}

	panic("Cursor should always be at end node")
}

func (m model) CurrentNote() (note, bool) {
	note, exists := m.currentOverlay.GetNote(m.gridCursor)
	return note, exists
}

func (m *model) NextOverlay(direction int) {
	switch direction {
	case 1:
		overlay := m.CurrentPart().Overlays.FindAboveOverlay(m.currentOverlay.Key)
		m.currentOverlay = overlay
	case -1:
		if m.currentOverlay.Below != nil {
			m.currentOverlay = m.currentOverlay.Below
		}
	default:
	}
}

func (m *model) ClearOverlayLine() {
	start := m.gridCursor.Beat
	for i := start; i < m.CurrentPart().Beats; i++ {
		m.currentOverlay.RemoveNote(GK(m.gridCursor.Line, i))
	}
}

func (m *model) ClearOverlay() {
	newOverlay := (*m.definition.parts)[m.CurrentPartID()].Overlays.Remove(m.currentOverlay.Key)
	if newOverlay != nil {
		(*m.definition.parts)[m.CurrentPartID()].Overlays = newOverlay
		m.currentOverlay = (*m.definition.parts)[m.CurrentPartID()].Overlays
		m.overlayKeyEdit.SetOverlayKey(m.currentOverlay.Key)
	}
}

type MoveFunc func()

func (m *model) RotateRight() {
	combinedPattern := m.CombinedEditPattern(m.currentOverlay)

	lineStart, lineEnd := m.PatternActionLineBoundaries()
	start, end := m.PatternActionBeatBoundaries()

	for l := lineStart; l <= lineEnd; l++ {

		moves := make([]MoveFunc, 0, (end - start))
		lastKey := GK(l, end)
		firstKey := GK(l, start)
		lastNote := combinedPattern[lastKey]
		previousNote := zeronote
		var previousKey grid.GridKey
		for i := start; i <= end; i++ {
			currentKey := GK(l, i)
			currentNote := combinedPattern[currentKey]

			chord, chordExists := m.currentOverlay.Chords.FindChordWithNote(previousKey)
			if chordExists {
				fromKey := previousKey
				moves = append(moves, func() {
					chord.Move(fromKey, currentKey)
				})
			} else {
				noteToMove := previousNote
				moves = append(moves, func() {
					m.currentOverlay.MoveNoteTo(currentKey, noteToMove)
				})
			}
			previousNote = currentNote
			previousKey = currentKey
		}
		noteToMove := lastNote
		moves = append(moves, func() {
			m.currentOverlay.MoveNoteTo(firstKey, noteToMove)
		})

		for _, moveFn := range moves {
			if moveFn != nil {
				moveFn()
			}
		}
	}
}

func (m *model) RotateLeft() {
	combinedPattern := m.CombinedEditPattern(m.currentOverlay)

	lineStart, lineEnd := m.PatternActionLineBoundaries()
	start, end := m.PatternActionBeatBoundaries()
	for l := lineStart; l <= lineEnd; l++ {
		moves := make([]MoveFunc, 0, (end - start))
		firstNote := combinedPattern[GK(l, start)]
		previousNote := zeronote
		var previousKey grid.GridKey
		for i := int8(end); i >= int8(start); i-- {
			currentKey := GK(l, uint8(i))
			currentNote := combinedPattern[currentKey]

			chord, chordExists := m.currentOverlay.Chords.FindChordWithNote(previousKey)
			if chordExists {
				fromKey := previousKey
				moves = append(moves, func() {
					chord.Move(fromKey, currentKey)
				})
			} else {
				noteToMove := previousNote
				moves = append(moves, func() {
					m.currentOverlay.MoveNoteTo(currentKey, noteToMove)
				})
			}

			previousNote = currentNote
			previousKey = currentKey
		}

		moves = append(moves, func() {
			m.currentOverlay.MoveNoteTo(GK(l, end), firstNote)
		})

		for _, moveFn := range moves {
			if moveFn != nil {
				moveFn()
			}
		}
	}
}

func (m *model) RotateUp() {
	pattern := m.CombinedEditPattern(m.currentOverlay)
	beat := m.gridCursor.Beat
	for l := 0; l < len(m.definition.lines); l++ {
		key := GK(uint8(l), beat)
		_, exists := pattern[key]
		if exists {
			m.currentOverlay.RemoveNote(key)
		}
		if l != 0 {
			newKey := GK(uint8(l-1), beat)
			note, exists := pattern[key]
			if exists {
				m.currentOverlay.SetNote(newKey, note)
			}
		} else {
			newKey := GK(uint8(len(m.definition.lines)-1), beat)
			note, exists := pattern[key]
			if exists {
				m.currentOverlay.SetNote(newKey, note)
			}
		}
	}
}

func (m *model) RotateDown() {
	pattern := m.CombinedEditPattern(m.currentOverlay)
	beat := m.gridCursor.Beat
	for l := len(m.definition.lines); l >= 0; l-- {
		key := GK(uint8(l), beat)
		_, exists := pattern[key]
		if exists {
			m.currentOverlay.RemoveNote(key)
		}
		index := l + 1
		if int(index) < len(m.definition.lines) {
			newKey := GK(uint8(index), beat)
			note, exists := pattern[key]
			if exists {
				m.currentOverlay.SetNote(newKey, note)
			}
		} else {
			newKey := GK(0, beat)
			note, exists := pattern[key]
			if exists {
				m.currentOverlay.SetNote(newKey, note)
			}
		}
	}
}

func (m *model) MoveChordLeft() {
	if m.activeChord.HasValue() {
		gridChord := m.activeChord.GridChord
		newBeat := gridChord.Root.Beat - 1
		gridChord.Root.Beat = newBeat
	}
}

func (m *model) MoveChordRight() {
	if m.activeChord.HasValue() {
		gridChord := m.activeChord.GridChord
		newBeat := gridChord.Root.Beat + 1
		gridChord.Root.Beat = newBeat
	}
}

func (m *model) MoveChordUp() {
	if m.activeChord.HasValue() {
		gridChord := m.activeChord.GridChord
		newLine := gridChord.Root.Line - 1
		gridChord.Root.Line = newLine
	}
}

func (m *model) MoveChordDown() {
	if m.activeChord.HasValue() {
		gridChord := m.activeChord.GridChord
		newLine := gridChord.Root.Line + 1
		gridChord.Root.Line = newLine
	}
}

func Mute(playState []linestate, line uint8) []linestate {
	switch playState[line].groupPlayState {
	case PlayStatePlay:
		playState[line].groupPlayState = PlayStateMute
	case PlayStateMute:
		playState[line].groupPlayState = PlayStatePlay
	case PlayStateSolo:
		playState[line].groupPlayState = PlayStateMute
	}
	return playState
}

func Solo(playState []linestate, line uint8) []linestate {
	switch playState[line].groupPlayState {
	case PlayStatePlay:
		playState[line].groupPlayState = PlayStateSolo
	case PlayStateMute:
		playState[line].groupPlayState = PlayStateSolo
	case PlayStateSolo:
		playState[line].groupPlayState = PlayStatePlay
	}
	return playState
}

func (m model) HasSolo() bool {
	for _, state := range m.playState {
		if state.groupPlayState == PlayStateSolo {
			return true
		}
	}
	return false
}

type Buffer struct {
	bounds    grid.Bounds
	gridNotes []GridNote
	gridChord *overlays.GridChord
}

func (b Buffer) IsChord() bool {
	if b.gridChord != nil {
		return len(b.gridChord.Notes) > 0
	} else {
		return false
	}
}

func InitChordBuffer(gridChord *overlays.GridChord) Buffer {
	return Buffer{
		gridChord: gridChord,
	}
}

func InitBuffer(bounds grid.Bounds, notes []GridNote) Buffer {
	return Buffer{
		bounds:    bounds.Normalized(),
		gridNotes: notes,
	}
}

func (m model) Yank() Buffer {
	currentChord := m.CurrentChord()
	if !m.visualMode && currentChord.HasValue() {
		return InitChordBuffer(currentChord.GridChord)
	} else {
		bounds := m.YankBounds()
		combinedPattern := m.CombinedEditPattern(m.currentOverlay)
		capturedGridNotes := make([]GridNote, 0, len(combinedPattern))

		for key, note := range combinedPattern {
			if bounds.InBounds(key) {
				normalizedGridKey := GK(key.Line-bounds.Top, key.Beat-bounds.Left)
				capturedGridNotes = append(capturedGridNotes, GridNote{normalizedGridKey, note})
			}
		}

		return InitBuffer(bounds, capturedGridNotes)
	}

}

func (m *model) Paste() {
	if m.yankBuffer.IsChord() {
		m.currentOverlay.PasteChord(m.gridCursor, m.yankBuffer.gridChord)
	} else {
		bounds := m.PasteBounds()

		var keyModifier gridKey
		if m.visualMode {
			keyModifier = bounds.TopLeft()
		} else {
			keyModifier = m.gridCursor
		}
		gridNotes := m.yankBuffer.gridNotes

		for _, gridNote := range gridNotes {
			key := gridNote.gridKey
			newKey := GK(key.Line+keyModifier.Line, key.Beat+keyModifier.Beat)
			note := gridNote.note
			if bounds.InBounds(newKey) {
				m.currentOverlay.SetNote(newKey, note)
			}
		}
	}
}

func (m *model) advanceCurrentBeat(playingOverlay *overlays.Overlay) {
	pattern := make(grid.Pattern)
	playingOverlay.CombineActionPattern(&pattern, m.CurrentSongSection().PlayCycles())
	for i := range m.playState {
		doContinue := m.advancePlayState(pattern, i)
		if !doContinue {
			break
		}
	}
}

func (m model) CurrentBeatGridKeys(gridKeys *[]grid.GridKey) {
	for _, linestate := range m.playState {
		if linestate.IsSolo() || (!linestate.IsMuted() && !m.hasSolo) {
			*gridKeys = append(*gridKeys, linestate.GridKey())
		}
	}
}

var chordFindDistance int = 3

func (m model) ChordFindBeatGridKeys(gridKeys *[]grid.GridKey) {
	start := max(0, int(m.gridCursor.Line)-chordFindDistance)
	end := min(len(m.definition.lines)-1, int(m.gridCursor.Line)+chordFindDistance)
	for i := start; i <= end; i++ {
		*gridKeys = append(*gridKeys, GK(uint8(i), m.gridCursor.Beat))
	}
}

func (m *model) advancePlayState(combinedPattern grid.Pattern, lineIndex int) bool {
	currentState := m.playState[lineIndex]
	advancedBeat := int8(currentState.currentBeat) + currentState.direction

	if advancedBeat >= int8(m.CurrentPart().Beats) || advancedBeat < 0 {
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
	case grid.ActionNothing:
		return true
	case grid.ActionLineReset:
		m.playState[lineIndex].currentBeat = 0
	case grid.ActionLineReverse:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-2, 0))
		m.playState[lineIndex].direction = -1
		m.playState[lineIndex].resetLocation = uint8(max(advancedBeat-1, 0))
		m.playState[lineIndex].resetActionLocation = uint8(advancedBeat)
		m.playState[lineIndex].resetAction = grid.ActionLineReverse
	case grid.ActionLineBounce:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-1, 0))
		m.playState[lineIndex].direction = -1
	case grid.ActionLineSkipBeat:
		m.advancePlayState(combinedPattern, lineIndex)
	case grid.ActionLineDelay:
		m.playState[lineIndex].currentBeat = uint8(max(advancedBeat-1, 0))
	case grid.ActionLineResetAll:
		for i := range m.playState {
			m.playState[i].currentBeat = 0
			m.playState[i].direction = 1
			m.playState[i].resetLocation = 0
			m.playState[i].resetDirection = 1
		}
		return false
	case grid.ActionLineBounceAll:
		for i := range m.playState {
			if i <= lineIndex {
				m.playState[i].currentBeat = uint8(max(m.playState[i].currentBeat-1, 0))
			}
			m.playState[i].direction = -1
		}
		return false
	case grid.ActionLineSkipBeatAll:
		for i := range m.playState {
			if i <= lineIndex {
				m.advancePlayState(combinedPattern, i)
			} else {
				m.advancePlayState(combinedPattern, i)
				m.advancePlayState(combinedPattern, i)
			}
		}
		return false
	}

	return true
}

func (m *model) advanceKeyCycle() {
	if m.playState[m.definition.keyline].currentBeat == 0 && m.loopMode != LoopOverlay {
		m.IncrementPlayCycles()
		songSection := m.CurrentSongSection()

		if songSection.IsDone() {
			if m.PlayMove() {
				m.StartPart()
			}
		}
	}
}

func (m *model) StartPart() {
	m.DuringPlayResetPlayCycles()
	m.playState = InitLineStates(len(m.definition.lines), m.playState, uint8(m.CurrentSongSection().StartBeat))
	m.ResetCurrentOverlay()
}

func (m *model) PlayMove() bool {
	if m.arrangement.Cursor.IsRoot() {
		m.StartStop()
		m.arrangement.Cursor.MoveNext()
		m.ResetCurrentOverlay()
		return false
	} else if m.arrangement.Cursor.IsLastSibling() {
		m.arrangement.Cursor.GetParentNode().DrawDown()
		if m.arrangement.Cursor.HasParentIterations() {
			m.arrangement.Cursor.MoveToFirstSibling()
			if m.arrangement.Cursor.GetCurrentNode().IsGroup() {
				m.arrangement.Cursor.MoveNext()
			}
		} else {
			m.arrangement.Cursor.ResetIterations()
			m.arrangement.Cursor.Up()
			return m.PlayMove()
		}
	} else {
		m.arrangement.Cursor.MoveToSibling()
	}
	return true
}

func (m model) PlayingOverlayKeys() []overlayKey {
	keys := make([]overlayKey, 0, 10)
	m.CurrentPart().Overlays.GetMatchingOverlayKeys(&keys, m.CurrentSongSection().PlayCycles())
	return keys
}

func (m model) CombinedEditPattern(overlay *overlays.Overlay) grid.Pattern {
	pattern := make(grid.Pattern)
	overlay.CombinePattern(&pattern, overlay.Key.GetMinimumKeyCycle())
	return pattern
}

func (m model) CombinedBeatPattern(overlay *overlays.Overlay) grid.Pattern {
	pattern := make(grid.Pattern)
	gridKeys := make([]grid.GridKey, 0, len(m.playState))
	m.CurrentBeatGridKeys(&gridKeys)
	overlay.CurrentBeatOverlayPattern(&pattern, m.CurrentSongSection().PlayCycles(), gridKeys)
	return pattern
}

func (m model) CombinedOverlayPattern(overlay *overlays.Overlay) overlays.OverlayPattern {
	pattern := make(overlays.OverlayPattern)
	if m.playing != PlayStopped && !m.playEditing {
		m.CurrentPart().Overlays.CombineOverlayPattern(&pattern, m.CurrentSongSection().PlayCycles())
	} else {
		overlay.CombineOverlayPattern(&pattern, overlay.Key.GetMinimumKeyCycle())
	}
	return pattern
}

func (m *model) Every(every uint8, everyFn func(gridKey)) {
	if m.definition.templateSequencerType == operation.SeqModeChord {
		bounds := m.PasteBounds()
		combinedOverlay := m.CombinedEditPattern(m.currentOverlay)
		keys := slices.Collect(maps.Keys(combinedOverlay))
		slices.SortFunc(keys, grid.Compare)
		counter := 0
		for _, gk := range keys {
			if bounds.InBounds(gk) {
				if counter%int(every) == 0 {
					everyFn(gk)
				}
				counter++
			}
		}
	} else {
		lineStart, lineEnd := m.PatternActionLineBoundaries()
		start, end := m.PatternActionBeatBoundaries()

		for l := lineStart; l <= lineEnd; l++ {
			for i := start; i <= end; i += every {
				everyFn(GK(l, i))
			}
		}
	}
}

func (m *model) fill(every uint8) {
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	everyFn := func(gridKey gridKey) {
		currentNote, hasNote := combinedOverlay[gridKey]
		hasNote = hasNote && currentNote != zeronote

		if hasNote {
			m.currentOverlay.SetNote(gridKey, zeronote)
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
			m.currentOverlay.SetNote(gridKey, currentNote.IncrementGate(modifier, len(config.ShortGates)+len(config.LongGates)))
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

func (m *model) IncrementSpecificValue(note grid.Note) {
	note.AccentIndex = note.AccentIndex + 1
	if note.AccentIndex < 128 {

		m.currentOverlay.SetNote(m.gridCursor, note)
	}
}

func (m *model) DecrementSpecificValue(note grid.Note) {
	if note.AccentIndex > 0 {
		note.AccentIndex = note.AccentIndex - 1
		m.currentOverlay.SetNote(m.gridCursor, note)
	}
}

func (m *model) AccentModify(modifier int8) {
	modifyFunc := func(key gridKey, currentNote note) {
		m.currentOverlay.SetNote(key, currentNote.IncrementAccent(modifier, uint8(len(config.Accents))))
	}
	m.Modify(modifyFunc)
}

func (m *model) GateModify(modifier int8) {
	modifyFunc := func(key gridKey, currentNote note) {
		m.currentOverlay.SetNote(key, currentNote.IncrementGate(modifier, len(config.ShortGates)+len(config.LongGates)))
	}
	m.Modify(modifyFunc)
}

type Boundable interface {
	InBounds(key grid.GridKey) bool
}

func (m *model) Modify(modifyFunc func(gridKey, note)) {

	var bounds Boundable
	currentChord := m.CurrentChord()
	if !m.visualMode && currentChord.HasValue() {
		bounds = currentChord.GridChord
	} else {
		bounds = m.YankBounds()
	}
	combinedOverlay := m.CombinedEditPattern(m.currentOverlay)

	for key, currentNote := range combinedOverlay {
		if bounds.InBounds(key) {
			if currentNote != zeronote {
				modifyFunc(key, currentNote)
			}
		}
	}
}

func (m *model) WaitModify(modifier int8) {
	modifyFunc := func(key gridKey, currentNote note) {
		m.currentOverlay.SetNote(key, currentNote.IncrementWait(modifier))
	}

	m.Modify(modifyFunc)
}

func (m model) PatternActionBeatBoundaries() (uint8, uint8) {
	if m.visualMode {
		if m.visualAnchorCursor.Beat < m.gridCursor.Beat {
			return m.visualAnchorCursor.Beat, m.gridCursor.Beat
		} else {
			return m.gridCursor.Beat, m.visualAnchorCursor.Beat
		}
	} else {
		return m.gridCursor.Beat, m.CurrentPart().Beats - 1
	}
}

func (m model) PatternActionLineBoundaries() (uint8, uint8) {
	if m.visualMode {
		if m.visualAnchorCursor.Line < m.gridCursor.Line {
			return m.visualAnchorCursor.Line, m.gridCursor.Line
		} else {
			return m.gridCursor.Line, m.visualAnchorCursor.Line
		}
	} else {
		return m.gridCursor.Line, m.gridCursor.Line
	}
}

func (m *model) ReloadFile() {
	if m.filename != "" {
		var err error
		m.definition, err = LoadFile(m.filename, m.definition.template)
		if err != nil {
			m.SetCurrentError(fault.Wrap(err, fmsg.WithDesc("could not reload file", fmt.Sprintf("Could not reload file %s", m.filename))))
		}
		m.ResetCurrentOverlay()
	}
}

// ----------------------------- View Layer ----------------------------------------

func (m model) CyclesEditView() string {
	var buf strings.Builder
	cycles := m.CurrentSongSection().Cycles
	startCycles := m.CurrentSongSection().StartCycles
	cyclesInput := themes.NumberStyle.Render(strconv.Itoa(int(cycles)))
	startCyclesInput := themes.NumberStyle.Render(strconv.Itoa(int(startCycles)))
	switch m.selectionIndicator {
	case operation.SelectCycles:
		cyclesInput = themes.SelectedStyle.Render(strconv.Itoa(int(cycles)))
	case operation.SelectStartCycles:
		startCyclesInput = themes.SelectedStyle.Render(strconv.Itoa(int(startCycles)))
	}
	buf.WriteString(themes.AltArtStyle.Render(" ⟳ Amount "))
	buf.WriteString(cyclesInput)
	buf.WriteString(themes.AltArtStyle.Render("    ⟳ Start "))
	buf.WriteString(startCyclesInput)
	buf.WriteString("\n")
	return buf.String()
}

func (m model) BeatsEditView() string {
	var buf strings.Builder
	beats := m.CurrentPart().Beats
	startBeats := m.CurrentSongSection().StartBeat

	beatsInput := themes.NumberStyle.Render(strconv.Itoa(int(beats)))
	startBeatsInput := themes.NumberStyle.Render(strconv.Itoa(int(startBeats)))
	switch m.selectionIndicator {
	case operation.SelectBeats:
		beatsInput = themes.SelectedStyle.Render(strconv.Itoa(int(beats)))
	case operation.SelectStartBeats:
		startBeatsInput = themes.SelectedStyle.Render(strconv.Itoa(int(startBeats)))
	}
	buf.WriteString(themes.AltArtStyle.Render(" Beats "))
	buf.WriteString(beatsInput)
	buf.WriteString(themes.AltArtStyle.Render("  Start Beat "))
	buf.WriteString(startBeatsInput)
	buf.WriteString("\n")
	return buf.String()
}

func (m model) TempoEditView() string {
	var tempo, division string
	tempo = themes.NumberStyle.Render(strconv.Itoa(m.definition.tempo))
	division = themes.NumberStyle.Render(strconv.Itoa(m.definition.subdivisions))
	switch m.selectionIndicator {
	case operation.SelectTempo:
		tempo = themes.SelectedStyle.Render(strconv.Itoa(m.definition.tempo))
	case operation.SelectTempoSubdivision:
		division = themes.SelectedStyle.Render(strconv.Itoa(m.definition.subdivisions))
	}
	var buf strings.Builder
	buf.WriteString(themes.AltArtStyle.Render(" Tempo "))
	buf.WriteString(tempo)
	buf.WriteString(themes.AltArtStyle.Render("  Subdivisions "))
	buf.WriteString(division)
	buf.WriteString("\n")
	return buf.String()
}

func (m model) SpecificValueEditView(note grid.Note) string {
	var specificValue = themes.SelectedStyle.Render(fmt.Sprintf("%d", note.AccentIndex))
	var buf strings.Builder
	buf.WriteString(themes.AltArtStyle.Render(" Specific Value "))
	buf.WriteString(specificValue)
	buf.WriteString("\n")
	return buf.String()
}

func (m model) WriteView() string {
	if m.NeedsWrite() {
		return " [+]"
	} else {
		return "    "
	}
}

func (m model) IsAccentSelector() bool {
	states := []operation.Selection{operation.SelectAccentDiff, operation.SelectAccentTarget, operation.SelectAccentStart}
	return slices.Contains(states, m.selectionIndicator)
}

func (m model) IsRatchetSelector() bool {
	states := []operation.Selection{operation.SelectRatchets, operation.SelectRatchetSpan}
	return slices.Contains(states, m.selectionIndicator)
}

func (m model) View() string {
	var buf strings.Builder
	var sideView string

	if m.patternMode == operation.PatternAccent || m.IsAccentSelector() {
		sideView = m.AccentKeyView()
	} else if (m.CurrentPart().Overlays.Key == overlaykey.ROOT && m.CurrentPart().Overlays.IsFresh() && len(*m.definition.parts) == 1 && m.CurrentPartID() == 0) ||
		slices.Contains([]operation.Selection{operation.SelectSetupValue, operation.SelectSetupMessageType, operation.SelectSetupChannel}, m.selectionIndicator) {
		// NOTE: We want to show the setupView on the very initial screen, before any sequencing has begun OR a setup value is selected
		sideView = m.SetupView()
	} else {
		sideView = m.OverlaysView()

		var chordView string
		if m.definition.templateSequencerType == operation.SeqModeChord {
			currentChord := m.CurrentChord()
			chordView = m.ChordView(currentChord.GridChord)
		}
		sideView = lipgloss.JoinVertical(lipgloss.Left, sideView, chordView)
	}

	seqView := m.SeqView()
	buf.WriteString(lipgloss.JoinHorizontal(0, "  ", seqView, "  ", sideView))
	if m.currentError != nil && m.selectionIndicator == operation.SelectError {
		buf.WriteString("\n")
		style := lipgloss.NewStyle().Width(50)
		style = style.Border(lipgloss.NormalBorder())
		style = style.Padding(1)
		style = style.BorderForeground(lipgloss.Color("#880000"))
		style = style.MarginLeft(2)
		var errorBuf strings.Builder
		errorBuf.WriteString("ERROR: ")
		issue := fmsg.GetIssue(m.currentError)
		if issue != "" {
			errorBuf.WriteString(issue)
		} else {
			chain := fault.Flatten(m.currentError)
			errorBuf.WriteString(chain[0].Message)
		}
		buf.WriteString(style.Render(errorBuf.String()))
	} else {
		buf.WriteString("\n")
	}
	if m.showArrangementView {
		buf.WriteString("\n")
		buf.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, "  ", m.arrangement.View()))
	}
	return buf.String()
}

func (m model) AccentKeyView() string {
	var buf strings.Builder
	var accentDiffString string
	var accentDiff = m.definition.accents.Diff
	var accentStart = m.definition.accents.Start

	var accentTarget string
	switch m.definition.accents.Target {
	case AccentTargetNote:
		accentTarget = "N"
	case AccentTargetVelocity:
		accentTarget = "V"
	}

	if m.selectionIndicator == operation.SelectAccentDiff {
		accentDiffString = themes.SelectedStyle.Render(fmt.Sprintf("%2d", accentDiff))
	} else {
		accentDiffString = themes.NumberStyle.Render(fmt.Sprintf("%2d", accentDiff))
	}

	var accentTargetString string
	if m.selectionIndicator == operation.SelectAccentTarget {
		accentTargetString = themes.SelectedStyle.Render(fmt.Sprintf(" %s", accentTarget))
	} else {
		accentTargetString = themes.NumberStyle.Render(fmt.Sprintf(" %s", accentTarget))
	}

	title := themes.AppDescriptorStyle.Render("Accents")
	buf.WriteString(fmt.Sprintf(" %s %s %s\n", title, accentDiffString, accentTargetString))
	buf.WriteString(themes.SeqBorderStyle.Render("──────────────"))
	buf.WriteString("\n")

	var accentStartString string
	if m.selectionIndicator == operation.SelectAccentStart {
		accentStartString = themes.SelectedStyle.Render(fmt.Sprintf("%2d", accentStart))
	} else {
		accentStartString = themes.NumberStyle.Render(fmt.Sprintf("%2d", accentStart))
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(themes.AccentColors[1]))
	buf.WriteString(fmt.Sprintf("  %s  -  %s\n", style.Render(string(themes.AccentIcons[1])), accentStartString))
	for i, accent := range m.definition.accents.Data[2:] {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(themes.AccentColors[i+2]))
		buf.WriteString(fmt.Sprintf("  %s  -  %d\n", style.Render(string(themes.AccentIcons[i+2])), accent.Value))
	}
	return buf.String()
}

func (m model) SetupView() string {
	var buf strings.Builder
	buf.WriteString(themes.AppDescriptorStyle.Render("Setup"))
	buf.WriteString("\n")
	buf.WriteString(themes.SeqBorderStyle.Render("──────────────"))
	buf.WriteString("\n")
	for i, line := range m.definition.lines {

		buf.WriteString("CH ")
		if uint8(i) == m.gridCursor.Line && m.selectionIndicator == operation.SelectSetupChannel {
			buf.WriteString(themes.SelectedStyle.Render(fmt.Sprintf("%2d", line.Channel)))
		} else {
			buf.WriteString(themes.NumberStyle.Render(fmt.Sprintf("%2d", line.Channel)))
		}

		var messageType string
		switch line.MsgType {
		case grid.MessageTypeNote:
			messageType = "NOTE"
		case grid.MessageTypeCc:
			messageType = "CC"
		case grid.MessageTypeProgramChange:
			messageType = "Program Change"
		}

		if uint8(i) == m.gridCursor.Line && m.selectionIndicator == operation.SelectSetupMessageType {
			messageType = fmt.Sprintf(" %s ", themes.SelectedStyle.Render(messageType))
		} else {
			messageType = fmt.Sprintf(" %s ", messageType)
		}

		buf.WriteString(messageType)

		if line.MsgType == grid.MessageTypeProgramChange {
			buf.WriteString("")
		} else {
			if uint8(i) == m.gridCursor.Line && m.selectionIndicator == operation.SelectSetupValue {
				buf.WriteString(themes.SelectedStyle.Render(strconv.Itoa(int(line.Note))))
			} else {
				buf.WriteString(themes.NumberStyle.Render(strconv.Itoa(int(line.Note))))
			}
		}
		buf.WriteString(fmt.Sprintf(" %s\n", LineValueName(line, m.definition.instrument)))
	}
	return buf.String()
}

func NoteName(note uint8) string {
	return fmt.Sprintf("%s%d", strings.ReplaceAll(midi.Note(note).Name(), "b", "♭"), int(midi.Note(note).Octave())-2)
}

func LineValueName(ld grid.LineDefinition, instrument string) string {
	switch ld.MsgType {
	case grid.MessageTypeNote:
		return NoteName(ld.Note)
	case grid.MessageTypeCc:
		cc, _ := config.FindCC(ld.Note, instrument)
		return cc.Name
	}
	return ""
}

func (m model) ChordView(gridChord *overlays.GridChord) string {

	var buf strings.Builder
	buf.WriteString(themes.AppDescriptorStyle.Render("Chord"))
	if gridChord == nil {
		buf.WriteString("\n")
		buf.WriteString(themes.SeqBorderStyle.Render("──────────────"))
		return buf.String()
	}
	buf.WriteString(" - ")
	chord := gridChord.Chord
	pattern := make(grid.Pattern)
	gridChord.ChordNotes(&pattern)
	baseNote := m.definition.lines[0].Note
	buf.WriteString(NoteName(baseNote - gridChord.Root.Line))
	buf.WriteString("\n")
	buf.WriteString(themes.SeqBorderStyle.Render("──────────────"))
	buf.WriteString("\n")

	buf.WriteString(fmt.Sprintf("Inversions: %d", chord.Inversion))
	buf.WriteString("\n")

	intervals := chord.NamedIntervals()
	uninvertedNotes := chord.UninvertedNotes()
	slices.Reverse(uninvertedNotes)
	for i, n := range uninvertedNotes {
		buf.WriteString(fmt.Sprintf("%d - %s - %s", n, intervals[i], NoteName(baseNote-gridChord.Root.Line+n)))
		buf.WriteString("\n")
	}

	return buf.String()
}

func (m model) OverlaysView() string {
	var buf strings.Builder
	buf.WriteString(themes.AppDescriptorStyle.Render("Overlays"))
	buf.WriteString("\n")
	buf.WriteString(themes.SeqBorderStyle.Render("──────────────"))
	buf.WriteString("\n")
	playingStyle := lipgloss.NewStyle().Background(themes.SeqOverlayColor).Foreground(themes.AppDescriptorColor)
	notPlayingStyle := themes.AppDescriptorStyle
	var playingOverlayKeys = m.PlayingOverlayKeys()
	for currentOverlay := m.CurrentPart().Overlays; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		var playingSpacer = "   "
		var playing = ""
		if m.playing != PlayStopped && playingOverlayKeys[0] == currentOverlay.Key {
			playing = themes.OverlayCurrentlyPlayingSymbol
			buf.WriteString(playing)
			playingSpacer = ""
		} else if m.playing != PlayStopped && slices.Contains(playingOverlayKeys, currentOverlay.Key) {
			playing = themes.ActiveSymbol
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
		if m.playing != PlayStopped && slices.Contains(playingOverlayKeys, currentOverlay.Key) {
			buf.WriteString(playingStyle.Render(overlayLine))
		} else {
			buf.WriteString(notPlayingStyle.Render(overlayLine))
		}
		buf.WriteString(playing)
		buf.WriteString(playingSpacer)
		buf.WriteString("\n")
	}
	return buf.String()
}

func (m model) SeqView() string {
	var buf strings.Builder
	var mode string

	visualCombinedPattern := m.CombinedOverlayPattern(m.currentOverlay)
	currentNote := visualCombinedPattern[m.gridCursor]

	buf.WriteString(m.WriteView())
	if m.patternMode == operation.PatternAccent {
		mode = " Accent "
		buf.WriteString(fmt.Sprintf(" %s  %s\n", themes.AccentModeStyle.Render(" PATTERN MODE "), themes.AccentModeStyle.Render(mode)))
	} else if m.patternMode == operation.PatternGate {
		mode = " Gate "
		buf.WriteString(fmt.Sprintf(" %s  %s\n", themes.AccentModeStyle.Render(" PATTERN MODE "), themes.AccentModeStyle.Render(mode)))
	} else if m.patternMode == operation.PatternWait {
		mode = " Wait "
		buf.WriteString(fmt.Sprintf(" %s  %s\n", themes.AccentModeStyle.Render(" PATTERN MODE "), themes.AccentModeStyle.Render(mode)))
	} else if m.patternMode == operation.PatternRatchet {
		mode = " Ratchet "
		buf.WriteString(fmt.Sprintf(" %s  %s\n", themes.AccentModeStyle.Render(" PATTERN MODE "), themes.AccentModeStyle.Render(mode)))
	} else if m.selectionIndicator == operation.SelectRatchets || m.selectionIndicator == operation.SelectRatchetSpan {
		buf.WriteString(m.RatchetEditView())
	} else if m.selectionIndicator == operation.SelectTempo || m.selectionIndicator == operation.SelectTempoSubdivision {
		buf.WriteString(m.TempoEditView())
	} else if slices.Contains([]operation.Selection{operation.SelectBeats, operation.SelectStartBeats}, m.selectionIndicator) {
		buf.WriteString(m.BeatsEditView())
	} else if slices.Contains([]operation.Selection{operation.SelectCycles, operation.SelectStartCycles}, m.selectionIndicator) {
		buf.WriteString(m.CyclesEditView())
	} else if m.selectionIndicator == operation.SelectPart {
		buf.WriteString(m.ChoosePartView())
	} else if m.selectionIndicator == operation.SelectChangePart {
		buf.WriteString(m.ChoosePartView())
	} else if m.selectionIndicator == operation.SelectRenamePart {
		buf.WriteString(m.RenamePartView())
	} else if m.selectionIndicator == operation.SelectFileName {
		buf.WriteString(m.FileNameView())
	} else if m.selectionIndicator == operation.SelectConfirmNew {
		buf.WriteString(m.ConfirmNewSequenceView())
	} else if m.selectionIndicator == operation.SelectConfirmQuit {
		buf.WriteString(m.ConfirmQuitView())
	} else if m.selectionIndicator == operation.SelectConfirmReload {
		buf.WriteString(m.ConfirmReloadView())
	} else if m.selectionIndicator == operation.SelectSpecificValue {
		buf.WriteString(m.SpecificValueEditView(currentNote.Note))
	} else if m.playing != PlayStopped {
		buf.WriteString(m.arrangement.Cursor.PlayStateView(m.CurrentSongSection().PlayCycles()))
	} else if len(*m.definition.parts) > 1 {
		buf.WriteString(themes.AppTitleStyle.Render(" Seq "))
		buf.WriteString(themes.AppDescriptorStyle.Render(fmt.Sprintf("- %s", m.CurrentPart().GetName())))
		buf.WriteString("\n")
	} else {
		buf.WriteString(themes.AppTitleStyle.Render(" Seq "))
		buf.WriteString(themes.AppDescriptorStyle.Render("- A sequencer for your cli"))
		buf.WriteString("\n")
	}

	beats := m.CurrentPart().Beats
	topLine := strings.Repeat("─", max(32, int(beats)))
	buf.WriteString("   ")
	buf.WriteString(themes.SeqBorderStyle.Render(fmt.Sprintf(" ┌%s", topLine)))
	buf.WriteString("\n")

	for i := uint8(0); i < uint8(len(m.definition.lines)); i++ {
		buf.WriteString(lineView(i, m, visualCombinedPattern))
	}

	buf.WriteString(m.CurrentOverlayView())
	return buf.String()
}

func (m model) RenamePartView() string {
	var buf strings.Builder
	buf.WriteString(" Rename Part: ")
	buf.WriteString(m.textInput.View())
	buf.WriteString("\n")
	return buf.String()
}

func (m model) FileNameView() string {
	var buf strings.Builder
	buf.WriteString(" File Name: ")
	buf.WriteString(m.textInput.View())
	buf.WriteString("\n")
	return buf.String()
}

func (m model) ChoosePartView() string {
	var buf strings.Builder
	buf.WriteString(" Choose Part: ")
	var name string
	if m.partSelectorIndex < 0 {
		name = "New Part"
	} else {
		name = (*m.definition.parts)[m.partSelectorIndex].GetName()
	}
	buf.WriteString(themes.SelectedStyle.Render(name))
	buf.WriteString("\n")
	return buf.String()
}

func (m model) ConfirmNewSequenceView() string {
	var buf strings.Builder
	buf.WriteString(" New Sequence: ")
	buf.WriteString(themes.SelectedStyle.Render("Confirm"))
	buf.WriteString("\n")
	return buf.String()
}

func (m model) ConfirmQuitView() string {
	var buf strings.Builder
	buf.WriteString(" Quit: ")
	buf.WriteString(themes.SelectedStyle.Render("Confirm"))
	buf.WriteString("\n")
	return buf.String()
}

func (m model) ConfirmReloadView() string {
	var buf strings.Builder
	buf.WriteString(" Reload: ")
	buf.WriteString(themes.SelectedStyle.Render("Confirm"))
	buf.WriteString("\n")
	return buf.String()
}

func (m model) RatchetEditView() string {
	currentNote, _ := m.CurrentNote()

	var buf strings.Builder
	var ratchetsBuf strings.Builder
	buf.WriteString(" Ratchets ")
	for i := range uint8(8) {
		var backgroundColor lipgloss.Color
		if i <= currentNote.Ratchets.Length {
			if m.ratchetCursor == i && m.selectionIndicator == operation.SelectRatchets {
				backgroundColor = themes.SelectedAttributeColor
			}
			if currentNote.Ratchets.HitAt(i) {
				ratchetsBuf.WriteString(themes.ActiveStyle.Background(backgroundColor).Render("\u25CF"))
			} else {
				ratchetsBuf.WriteString(themes.MutedStyle.Background(backgroundColor).Render("\u25C9"))
			}
			ratchetsBuf.WriteString(" ")
		} else {

			ratchetsBuf.WriteString("  ")
		}
	}
	buf.WriteString(fmt.Sprintf("%*s", 32, ratchetsBuf.String()))
	if m.selectionIndicator == operation.SelectRatchetSpan {
		buf.WriteString(fmt.Sprintf(" Span %s ", themes.SelectedStyle.Render(strconv.Itoa(int(currentNote.Ratchets.GetSpan())))))
	} else {
		buf.WriteString(fmt.Sprintf(" Span %s ", themes.NumberStyle.Render(strconv.Itoa(int(currentNote.Ratchets.GetSpan())))))
	}
	buf.WriteString("\n")

	return buf.String()
}

func (m model) ViewOverlay() string {
	return m.overlayKeyEdit.ViewOverlay()
}

func (m model) CurrentOverlayView() string {
	var matchedKey overlayKey
	if m.playing != PlayStopped {
		matchedKey = m.CurrentPart().Overlays.HighestMatchingOverlay(m.CurrentSongSection().PlayCycles()).Key
	} else {
		matchedKey = overlaykey.ROOT
	}

	var editOverlayTitle string
	if m.playEditing {
		editOverlayTitle = lipgloss.NewStyle().Background(themes.SeqOverlayColor).Foreground(themes.AppTitleColor).Render("Edit")
	} else {
		editOverlayTitle = lipgloss.NewStyle().Foreground(themes.AppTitleColor).Render("Edit")
	}

	playOverlayTitle := lipgloss.NewStyle().Foreground(themes.AppTitleColor).Render("Play")

	editOverlay := fmt.Sprintf("%s %s", editOverlayTitle, lipgloss.PlaceHorizontal(11, 0, m.ViewOverlay()))
	playOverlay := fmt.Sprintf("%s %s", playOverlayTitle, lipgloss.PlaceHorizontal(11, 0, overlaykey.View(matchedKey)))
	return fmt.Sprintf("    %s  %s", editOverlay, playOverlay)
}

func KeyLineIndicator(k uint8, l uint8) string {
	if k == l {
		return themes.AltArtStyle.Render("K")
	} else {
		return " "
	}
}

var blackNotes = []uint8{1, 3, 6, 8, 10}

func ensureStringLengthWc(s string, length int) string {
	if ansi.StringWidthWc(s) <= length {
		padding := length - ansi.StringWidthWc(s)
		return strings.Repeat(" ", padding) + s
	}

	return ansi.CutWc(s, 0, length)
}

func (m model) LineIndicator(lineNumber uint8) string {
	indicator := themes.SeqBorderStyle.Render("│")
	if lineNumber == m.gridCursor.Line {
		indicator = themes.LineCursorStyle.Render("┤")
	}
	if len(m.playState) > int(lineNumber) && m.playState[lineNumber].groupPlayState == PlayStateMute {
		indicator = "M"
	}
	if len(m.playState) > int(lineNumber) && m.playState[lineNumber].groupPlayState == PlayStateSolo {
		indicator = "S"
	}

	var lineName string
	if m.definition.lines[lineNumber].Name != "" {
		lineName = themes.LineNumberStyle.Render(m.definition.lines[lineNumber].Name)
	} else if m.definition.templateUIStyle == "blackwhite" {
		notename := NoteName(m.definition.lines[lineNumber].Note)
		if slices.Contains(blackNotes, m.definition.lines[lineNumber].Note%12) {
			lineName = themes.BlackKeyStyle.Render(notename[0:4])
		} else {
			lineName = themes.WhiteKeyStyle.Render(notename)
		}
	} else {
		lineName = themes.LineNumberStyle.Render(fmt.Sprintf("%2d", lineNumber))
	}

	return fmt.Sprintf("%3s%s%s", ensureStringLengthWc(lineName, 3), KeyLineIndicator(m.definition.keyline, lineNumber), indicator)

}

type GateSpace struct {
	StringValue []rune
	Color       lipgloss.Color
}

func (gs GateSpace) HasMore() bool {
	return len(gs.StringValue) > 0
}
func (gs *GateSpace) ShiftString() string {
	if len(gs.StringValue) == 1 {
		v := gs.StringValue
		gs.StringValue = []rune{}
		return string(v)
	} else if len(gs.StringValue) > 1 {
		v := gs.StringValue[0]
		gs.StringValue = gs.StringValue[1:]
		return string(v)
	} else {
		return ""
	}
}

func lineView(lineNumber uint8, m model, visualCombinedPattern overlays.OverlayPattern) string {
	var buf strings.Builder
	buf.WriteString(m.LineIndicator(lineNumber))

	gateSpace := GateSpace{}
	currentChord := m.CurrentChord()
	for i := uint8(0); i < m.CurrentPart().Beats; i++ {
		currentGridKey := GK(uint8(lineNumber), i)
		overlayNote, hasNote := visualCombinedPattern[currentGridKey]

		var backgroundSeqColor lipgloss.Color
		if m.playing != PlayStopped && m.playState[lineNumber].currentBeat == i {
			backgroundSeqColor = themes.SeqCursorColor
		} else if m.visualMode && m.InVisualSelection(currentGridKey) {
			backgroundSeqColor = themes.SeqVisualColor
		} else if hasNote && overlayNote.HighestOverlay && overlayNote.OverlayKey != overlaykey.ROOT {
			backgroundSeqColor = themes.SeqOverlayColor
		} else if hasNote && !overlayNote.HighestOverlay && overlayNote.OverlayKey != overlaykey.ROOT {
			backgroundSeqColor = themes.SeqMiddleOverlayColor
		} else if i%8 > 3 {
			backgroundSeqColor = themes.AltSeqBackgroundColor
		} else {
			backgroundSeqColor = themes.SeqBackgroundColor
		}

		char, foregroundColor := ViewNoteComponents(overlayNote.Note)
		var hasGateTail = false
		if (!hasNote || overlayNote.Note == zeronote) && gateSpace.HasMore() {
			char = gateSpace.ShiftString()
			hasGateTail = true
		} else if gateSpace.HasMore() {
			gateSpace = GateSpace{}
		}

		style := lipgloss.NewStyle().Background(backgroundSeqColor)
		cursorMatch := m.gridCursor.Line == uint8(lineNumber) && m.gridCursor.Beat == i
		if cursorMatch {
			m.cursor.SetChar(char)
			char = m.cursor.View()
		} else if m.visualMode && m.InVisualSelection(currentGridKey) {
			style = style.Foreground(themes.Black)
		} else if hasGateTail {
			style = style.Foreground(gateSpace.Color)
		} else {
			style = style.Foreground(foregroundColor)
		}

		if overlayNote.Note.GateIndex > uint8(len(config.ShortGates))-1 && int(overlayNote.Note.GateIndex) < int(len(config.ShortGates)+len(config.LongGates)) {
			gateSpaceValue := config.LongGates[overlayNote.Note.GateIndex-8].Shape
			gateSpace.StringValue = []rune(gateSpaceValue)
			gateSpace.Color = foregroundColor
		}

		gridChord, exists := m.currentOverlay.FindChordWithNote(currentGridKey)

		if exists && gridChord == currentChord.GridChord && !cursorMatch {
			fg := style.GetForeground()
			bg := style.GetBackground()
			style = style.Background(fg).Foreground(bg)
		}

		buf.WriteString(style.Render(char))
	}

	buf.WriteString("\n")
	return buf.String()
}

func ViewNoteComponents(currentNote grid.Note) (string, lipgloss.Color) {
	currentAction := currentNote.Action
	var char string
	var foregroundColor lipgloss.Color
	var waitShape string

	if currentNote.WaitIndex > 0 {
		waitShape = "\u0320"
	}

	if currentAction == grid.ActionNothing && currentNote != zeronote {
		currentAccentShape := themes.AccentIcons[currentNote.AccentIndex]
		currentAccentColor := themes.AccentColors[currentNote.AccentIndex]
		char = string(currentAccentShape) +
			string(config.Ratchets[currentNote.Ratchets.Length]) +
			ShortGate(currentNote) +
			waitShape
		foregroundColor = lipgloss.Color(currentAccentColor)
	} else {
		lineaction := config.Lineactions[currentAction]
		lineActionColor := themes.ActionColors[currentAction]
		char = string(lineaction.Shape)
		foregroundColor = lipgloss.Color(lineActionColor)
	}

	return char, foregroundColor
}

func ShortGate(note note) string {
	if note.GateIndex < uint8(len(config.ShortGates)) {
		return string(config.ShortGates[note.GateIndex].Shape)
	} else {
		return ""
	}
}
