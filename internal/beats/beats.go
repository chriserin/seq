package beats

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/notereg"
	"github.com/chriserin/seq/internal/overlays"
	"github.com/chriserin/seq/internal/playstate"
	"github.com/chriserin/seq/internal/seqmidi"
	"github.com/chriserin/seq/internal/sequence"
	midi "gitlab.com/gomidi/midi/v2"
)

type ModelMsg struct {
	PlayState      playstate.PlayState
	Definition     sequence.Sequence
	Cursor         arrangement.ArrCursor
	MidiConnection seqmidi.MidiConnection
}

type ModelPlayedMsg struct {
	PerformStop bool
	PlayState   playstate.PlayState
	Definition  sequence.Sequence
	Cursor      arrangement.ArrCursor
}

var beatChannel chan BeatMsg
var updateChannel chan ModelMsg

func init() {
	beatChannel = make(chan BeatMsg)
	updateChannel = make(chan ModelMsg)
}

func GetBeatChannel() chan BeatMsg {
	return beatChannel
}

func GetUpdateChannel() chan ModelMsg {
	return updateChannel
}

func Loop(sendFn func(tea.Msg)) {
	go func() {
		var playState playstate.PlayState
		var definition sequence.Sequence
		var cursor arrangement.ArrCursor
		var midiConn seqmidi.MidiConnection
		var errChan = make(chan error)
		for {
			if !playState.Playing {
				// NOTE: Wait for a model update that puts us into a playing state.
				modelMsg := <-updateChannel
				playState = modelMsg.PlayState
				definition = modelMsg.Definition
				cursor = modelMsg.Cursor
				midiConn = modelMsg.MidiConnection
			} else {
				// NOTE: In a plyaing state, respond to beat messages
				select {
				case modelMsg := <-updateChannel:
					playState = modelMsg.PlayState
					definition = modelMsg.Definition
					cursor = modelMsg.Cursor
					midiConn = modelMsg.MidiConnection
				case BeatMsg := <-beatChannel:
					midiSendFn, err := midiConn.AcquireSendFunc()
					if err != nil {
						errChan <- err
					} else {
						Beat(BeatMsg, playState, definition, cursor, midiSendFn, sendFn, errChan)
					}
				case err := <-errChan:
					fmt.Println(err)
				}
			}
		}
	}()
}

func Beat(msg BeatMsg, playState playstate.PlayState, definition sequence.Sequence, cursor arrangement.ArrCursor, midiSendFn seqmidi.SendFunc, sendFn func(tea.Msg), errChan chan error) {

	currentSection := &cursor[len(cursor)-1].Section
	partID := currentSection.Part
	currentPart := (*definition.Parts)[partID]
	playingOverlay := currentPart.Overlays.HighestMatchingOverlay(currentSection.PlayCycles())

	if playState.Playing && playState.RecordPreRollBeats == 0 {
		// NOTE: Only advance if we've already played the first beat.
		if playState.AllowAdvance {
			advanceCurrentBeat(*currentSection, playingOverlay, playState.LineStates, currentPart.Beats)
			advanceKeyCycle(definition.Keyline, playState.LineStates, playState.LoopMode, currentSection)
			if currentSection.IsDone() && playState.LoopMode != playstate.LoopOverlay {
				if PlayMove(&cursor) {
					cursor[len(cursor)-1].Section.DuringPlayReset()
					currentSection = &cursor[len(cursor)-1].Section
					playState.LineStates = playstate.InitLineStates(len(definition.Lines), playState.LineStates, uint8(cursor[len(cursor)-1].Section.StartBeat))
				} else {
					playState.LineStates = playstate.InitLineStates(len(definition.Lines), playState.LineStates, 0)
					sendFn(ModelPlayedMsg{PerformStop: true, PlayState: playState, Definition: definition, Cursor: cursor})
					return
				}
			}
		}
	}

	if playState.Playing {
		partID = currentSection.Part
		currentPart = (*definition.Parts)[partID]
		playingOverlay = currentPart.Overlays.HighestMatchingOverlay(currentSection.PlayCycles())
		gridKeys := make([]grid.GridKey, 0, len(playState.LineStates))
		CurrentBeatGridKeys(&gridKeys, playState.LineStates, playState.HasSolo)

		pattern := make(grid.Pattern)
		playingOverlay.CurrentBeatOverlayPattern(&pattern, currentSection.PlayCycles(), gridKeys)

		if playState.RecordPreRollBeats > 0 {
			playState.RecordPreRollBeats--
			return
		}

		err := PlayBeat(msg.Interval, pattern, definition, midiSendFn, errChan)

		if !playState.AllowAdvance {
			playState.AllowAdvance = true
		}

		if err != nil {
			errChan <- fault.Wrap(err, fmsg.With("error when playing beat"))
		}
		sendFn(ModelPlayedMsg{PlayState: playState, Definition: definition, Cursor: cursor})
	}
}

func CurrentBeatGridKeys(gridKeys *[]grid.GridKey, lineStates []playstate.LineState, hasSolo bool) {
	for _, linestate := range lineStates {
		if linestate.IsSolo() || (!linestate.IsMuted() && !hasSolo) {
			*gridKeys = append(*gridKeys, linestate.GridKey())
		}
	}
}

func advanceCurrentBeat(currentSection arrangement.SongSection, playingOverlay *overlays.Overlay, lineStates []playstate.LineState, partBeats uint8) {
	pattern := make(grid.Pattern)
	playingOverlay.CombineActionPattern(&pattern, currentSection.PlayCycles())
	for i := range lineStates {
		doContinue := lineStates[i].AdvancePlayState(pattern, i, partBeats, lineStates)
		if !doContinue {
			break
		}
	}
}

func advanceKeyCycle(keyline uint8, lineStates []playstate.LineState, loopMode playstate.LoopMode, section *arrangement.SongSection) {
	if lineStates[keyline].CurrentBeat == 0 && loopMode != playstate.LoopOverlay {
		section.IncrementPlayCycles()
	}
}

func PlayMove(cursor *arrangement.ArrCursor) bool {
	if cursor.IsRoot() {
		cursor.MoveNext()
		return false
	} else if cursor.IsLastSibling() {
		cursor.GetParentNode().DrawDown()
		if cursor.HasParentIterations() {
			cursor.MoveToFirstSibling()
			if cursor.GetCurrentNode().IsGroup() {
				cursor.MoveNext()
			}
		} else {
			cursor.ResetIterations()
			cursor.Up()
			return PlayMove(cursor)
		}
	} else {
		cursor.MoveToSibling()
		cursor.ResetIterations()
	}
	return true
}

func PlayBeat(beatInterval time.Duration, pattern grid.Pattern, definition sequence.Sequence, midiSendFn seqmidi.SendFunc, errChan chan error) error {
	lines := definition.Lines

	for gridKey, note := range pattern {
		line := lines[gridKey.Line]
		if note.Ratchets.Length > 0 {
			ProcessRatchets(note, beatInterval, line, definition, midiSendFn, errChan)
		} else if note != grid.ZeroNote {
			accents := definition.Accents

			delay := Delay(note.WaitIndex, beatInterval)
			gateLength := GateLength(note.GateIndex, beatInterval)

			switch line.MsgType {
			case grid.MessageTypeNote:
				onMessage, offMessage := NoteMessages(
					line,
					definition.Accents.Data[note.AccentIndex].Value,
					gateLength,
					accents.Target,
					delay,
				)
				err := ProcessNoteMsg(onMessage, midiSendFn, errChan)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process note on msg"))
				}
				err = ProcessNoteMsg(offMessage, midiSendFn, errChan)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process note off msg"))
				}
			case grid.MessageTypeCc:
				ccMessage := CCMessage(line, note, accents.Data, delay, true, definition.Instrument)
				err := ProcessNoteMsg(ccMessage, midiSendFn, errChan)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process cc msg"))
				}
			case grid.MessageTypeProgramChange:
				pcMessage := PCMessage(line, note, accents.Data, delay, true, definition.Instrument)
				err := ProcessNoteMsg(pcMessage, midiSendFn, errChan)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process cc msg"))
				}
			}
		}
	}

	return nil
}

func ProcessRatchets(note grid.Note, beatInterval time.Duration, line grid.LineDefinition, definition sequence.Sequence, midiSendFn seqmidi.SendFunc, errChan chan error) error {
	for i := range note.Ratchets.Length + 1 {
		if note.Ratchets.HitAt(i) {
			shortGateLength := 20 * time.Millisecond
			ratchetInterval := time.Duration(i) * note.Ratchets.Interval(beatInterval)
			onMessage, offMessage := NoteMessages(line, definition.Accents.Data[note.AccentIndex].Value, shortGateLength, definition.Accents.Target, ratchetInterval)
			err := ProcessNoteMsg(onMessage, midiSendFn, errChan)
			if err != nil {
				return fault.Wrap(err, fmsg.With("cannot turn on ratchet note"))
			}
			err = ProcessNoteMsg(offMessage, midiSendFn, errChan)
			if err != nil {
				return fault.Wrap(err, fmsg.With("cannot turn off ratchet note"))
			}
		}
	}
	return nil
}

func ProcessNoteMsg(msg Delayable, midiSendFn seqmidi.SendFunc, errChan chan error) error {
	// TODO: We don't need the redirection, we know what each type of note is when this function is called
	switch msg := msg.(type) {
	case NoteMsg:
		switch msg.midiType {
		case midi.NoteOnMsg:
			if notereg.Has(msg) {
				notereg.Remove(msg)
				PlayMessage(0, msg.OffMessage(), midiSendFn, errChan)
			}
			if err := notereg.Add(msg); err != nil {
				return fault.Wrap(err, fmsg.With("added a note already in registry"))
			}
			PlayMessage(msg.delay, msg.GetOnMidi(), midiSendFn, errChan)
		case midi.NoteOffMsg:
			PlayOffMessage(msg, midiSendFn, errChan)
		}
	case controlChangeMsg:
		PlayMessage(msg.delay, msg.MidiMessage(), midiSendFn, errChan)
	case programChangeMsg:
		PlayMessage(msg.delay, msg.MidiMessage(), midiSendFn, errChan)
	}
	return nil
}

func PlayMessage(delay time.Duration, message midi.Message, sendFn seqmidi.SendFunc, errChan chan error) {
	time.AfterFunc(delay, func() {
		err := sendFn(message)
		if err != nil {
			errChan <- fault.Wrap(err, fmsg.With("cannot send play message"))
		}
	})
}

func PlayOffMessage(nm NoteMsg, sendFn seqmidi.SendFunc, errChan chan error) {
	time.AfterFunc(nm.delay, func() {
		if notereg.RemoveID(nm) {
			err := sendFn(nm.GetOffMidi())
			if err != nil {
				errChan <- fault.Wrap(err, fmsg.With("cannot send off message"))
			}
		}
	})
}

type BeatMsg struct {
	Interval time.Duration
}

func NoteMessages(l grid.LineDefinition, accentValue uint8, gateLength time.Duration, accentTarget sequence.AccentTarget, delay time.Duration) (NoteMsg, NoteMsg) {
	var noteValue uint8
	var velocityValue uint8

	switch accentTarget {
	case sequence.AccentTargetNote:
		noteValue = l.Note + accentValue
		velocityValue = 96
	case sequence.AccentTargetVelocity:
		noteValue = l.Note
		velocityValue = accentValue
	}

	id := rand.Int()
	onMsg := NoteMsg{id: id, midiType: midi.NoteOnMsg, channel: l.Channel - 1, noteValue: noteValue, velocity: velocityValue, delay: delay}
	offMsg := NoteMsg{id: id, midiType: midi.NoteOffMsg, channel: l.Channel - 1, noteValue: noteValue, velocity: 0, delay: delay + gateLength}

	return onMsg, offMsg
}

func CCMessage(l grid.LineDefinition, note grid.Note, accents []config.Accent, delay time.Duration, includeDelay bool, instrument string) controlChangeMsg {
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

func PCMessage(l grid.LineDefinition, note grid.Note, accents []config.Accent, delay time.Duration, includeDelay bool, instrument string) programChangeMsg {
	if note.Action == grid.ActionSpecificValue {
		return programChangeMsg{l.Channel - 1, note.AccentIndex, delay}
	} else {
		return programChangeMsg{l.Channel - 1, l.Note - 1, delay}
	}
}

type Delayable interface {
	Delay() time.Duration
}

type NoteMsg struct {
	channel   uint8
	noteValue uint8
	velocity  uint8
	midiType  midi.Type
	delay     time.Duration
	id        int
}

func (nm NoteMsg) Delay() time.Duration {
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

func (nm NoteMsg) GetKey() notereg.NoteRegKey {
	return notereg.NoteRegKey{
		Channel: nm.channel,
		Note:    nm.noteValue,
	}
}

func (nm NoteMsg) GetID() int {
	return nm.id
}

func (nm NoteMsg) GetOnMidi() midi.Message {
	return midi.NoteOn(nm.channel, nm.noteValue, nm.velocity)
}

func (nm NoteMsg) GetOffMidi() midi.Message {
	return midi.NoteOff(nm.channel, nm.noteValue)
}

func (nm NoteMsg) OffMessage() midi.Message {
	return midi.NoteOff(nm.channel, nm.noteValue)
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
