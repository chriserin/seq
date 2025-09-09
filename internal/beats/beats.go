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
	PlayState playstate.PlayState
	Sequence  sequence.Sequence
	Cursor    arrangement.ArrCursor
}

type ModelPlayedMsg struct {
	PerformStop bool
	PlayState   playstate.PlayState
	Cursor      arrangement.ArrCursor
}

type AnticipatoryStop struct{}

var beatChannel chan BeatMsg
var updateChannel chan ModelMsg
var doneChannel chan struct{}
var playQueue chan midi.Message
var errChan chan error

func init() {
	beatChannel = make(chan BeatMsg)
	updateChannel = make(chan ModelMsg)
	doneChannel = make(chan struct{})
	playQueue = make(chan midi.Message)
	errChan = make(chan error)
}

func GetBeatChannel() chan BeatMsg {
	return beatChannel
}

func GetUpdateChannel() chan ModelMsg {
	return updateChannel
}

func GetDoneChannel() chan struct{} {
	return doneChannel
}

func Loop(sendFn func(tea.Msg), midiConn seqmidi.MidiConnection) {
	logFile, _ := tea.LogToFile("debug.log", "debug")

	go func() {
		var playState playstate.PlayState
		var definition sequence.Sequence
		var cursor arrangement.ArrCursor
		for {
			if !playState.Playing {
				// NOTE: Wait for a model update that puts us into a playing state.
				select {
				case modelMsg := <-updateChannel:
					playState = modelMsg.PlayState
					definition = modelMsg.Sequence
					cursor = modelMsg.Cursor
				case <-doneChannel:
					return
				}
			} else {
				// NOTE: In a plyaing state, respond to beat messages
				select {
				case modelMsg := <-updateChannel:
					playState = modelMsg.PlayState
					definition = modelMsg.Sequence
					cursor = modelMsg.Cursor
				case BeatMsg := <-beatChannel:
					Beat(BeatMsg, playState, definition, cursor, sendFn, errChan)
				case <-doneChannel:
					return
				case err := <-errChan:
					_, logErr := fmt.Fprintf(logFile, "Error: %v", err)
					if logErr != nil {
						fmt.Println("An error occurred while writing the original error to the log file", err, logErr)
					}
				}
			}
		}
	}()
	go func() {
		midiSendFn, err := midiConn.AcquireSendFunc()
		if err != nil {
			errChan <- err
			return
		}
		for {
			midiMessage := <-playQueue
			err := midiSendFn(midiMessage)
			if err != nil {
				errChan <- err
			}
		}
	}()
}

func IsDone(playState playstate.PlayState, currentNode *arrangement.Arrangement, currentSection arrangement.SongSection) bool {
	return playState.LoopedArrangement != currentNode && currentSection.Cycles+currentSection.StartCycles <= (*playState.Iterations)[currentNode]
}

func Beat(msg BeatMsg, playState playstate.PlayState, definition sequence.Sequence, cursor arrangement.ArrCursor, sendFn func(tea.Msg), errChan chan error) {

	if playState.Playing {
		AdvancePlayState(&playState, definition, &cursor)
	}

	if !playState.Playing {
		sendFn(ModelPlayedMsg{PerformStop: true, PlayState: playState, Cursor: cursor})
		return
	} else {
		PlaySequence(&playState, definition, cursor, msg, errChan)
		sendFn(ModelPlayedMsg{PlayState: playState, Cursor: cursor})
	}

	copiedPlayState := playstate.Copy(playState)
	copiedCursor := make(arrangement.ArrCursor, len(cursor))
	copy(copiedCursor, cursor)
	AdvancePlayState(&copiedPlayState, definition, &copiedCursor)
	if !copiedPlayState.Playing {
		sendFn(AnticipatoryStop{})
	}
}

func PlaySequence(playState *playstate.PlayState, definition sequence.Sequence, cursor arrangement.ArrCursor, msg BeatMsg, errChan chan error) {

	currentNode := cursor[len(cursor)-1]
	currentSection := cursor[len(cursor)-1].Section
	var partID int
	var currentCycles int
	var currentPart arrangement.Part
	var playingOverlay *overlays.Overlay

	partID = currentSection.Part
	currentPart = (*definition.Parts)[partID]
	currentCycles = (*playState.Iterations)[currentNode]
	playingOverlay = currentPart.Overlays.HighestMatchingOverlay(currentCycles)

	noteLineStates := make([]playstate.LineState, 0, len(playState.LineStates))
	metaLineStates := make([]playstate.LineState, 0, len(playState.LineStates))
	for i, ls := range playState.LineStates {
		if definition.Lines[i].MsgType == grid.MessageTypeNote {
			noteLineStates = append(noteLineStates, ls)
		} else {
			metaLineStates = append(metaLineStates, ls)
		}
	}

	// Play the CC/PC Messages
	gridKeys := make([]grid.GridKey, 0, len(playState.LineStates))
	CurrentBeatGridKeys(&gridKeys, metaLineStates, playState.HasSolo)

	pattern := make(grid.Pattern)
	playingOverlay.CurrentBeatOverlayPattern(&pattern, currentCycles, gridKeys)

	err := PlayBeat(msg.Interval, pattern, definition, errChan)
	if err != nil {
		errChan <- fault.Wrap(err, fmsg.With("error when playing beat"))
		return
	}

	// Play the Note Messages
	gridKeys = make([]grid.GridKey, 0, len(playState.LineStates))
	CurrentBeatGridKeys(&gridKeys, noteLineStates, playState.HasSolo)

	pattern = make(grid.Pattern)
	playingOverlay.CurrentBeatOverlayPattern(&pattern, currentCycles, gridKeys)

	err = PlayBeat(msg.Interval, pattern, definition, errChan)

	if !playState.AllowAdvance {
		playState.AllowAdvance = true
	}

	if err != nil {
		errChan <- fault.Wrap(err, fmsg.With("error when playing beat"))
	}
}

func AdvancePlayState(playState *playstate.PlayState, definition sequence.Sequence, cursor *arrangement.ArrCursor) {
	currentNode := (*cursor)[len(*cursor)-1]
	currentSection := (*cursor)[len(*cursor)-1].Section
	partID := currentSection.Part
	currentPart := (*definition.Parts)[partID]
	currentCycles := (*playState.Iterations)[currentNode]
	playingOverlay := currentPart.Overlays.HighestMatchingOverlay(currentCycles)

	if playState.Playing {
		// NOTE: Only advance if we've already played the first beat.
		if playState.AllowAdvance {
			advanceCurrentBeat(currentCycles, *playingOverlay, playState.LineStates, currentPart.Beats)
			advanceKeyCycle(definition.Keyline, playState.LineStates, playState.LoopMode, currentNode, playState.Iterations)
			if IsDone(*playState, currentNode, currentSection) && playState.LoopMode != playstate.LoopOverlay {
				if PlayMove(cursor, playState.Iterations, playState.LoopedArrangement) || playState.PlayMode == playstate.PlayReceiver {
					currentSection = (*cursor)[len(*cursor)-1].Section
					currentNode = (*cursor)[len(*cursor)-1]
					if !currentSection.KeepCycles {
						(*playState.Iterations)[currentNode] = currentSection.StartCycles
					}
					playState.LineStates = playstate.InitLineStates(len(definition.Lines), playState.LineStates, uint8((*cursor)[len(*cursor)-1].Section.StartBeat))
				} else {
					playState.Playing = false
					return
				}
			}
		}
	}

}

func CurrentBeatGridKeys(gridKeys *[]grid.GridKey, lineStates []playstate.LineState, hasSolo bool) {
	for _, linestate := range lineStates {
		if linestate.IsSolo() || (!linestate.IsMuted() && !hasSolo) {
			*gridKeys = append(*gridKeys, linestate.GridKey())
		}
	}
}

func advanceCurrentBeat(keyCycles int, playingOverlay overlays.Overlay, lineStates []playstate.LineState, partBeats uint8) {
	pattern := make(grid.Pattern)
	playingOverlay.CombineActionPattern(&pattern, keyCycles)
	for i := range lineStates {
		doContinue := lineStates[i].AdvancePlayState(pattern, i, partBeats, lineStates)
		if !doContinue {
			break
		}
	}
}

func advanceKeyCycle(keyline uint8, lineStates []playstate.LineState, loopMode playstate.LoopMode, node *arrangement.Arrangement, iterations *map[*arrangement.Arrangement]int) {
	if lineStates[keyline].CurrentBeat == 0 && loopMode != playstate.LoopOverlay {
		(*iterations)[node]++
	}
}

func PlayMove(cursor *arrangement.ArrCursor, iterations *map[*arrangement.Arrangement]int, loopNode *arrangement.Arrangement) bool {
	if cursor.IsRoot() {
		cursor.MoveNext()
		return false
	} else if cursor.IsLastSibling() {
		(*iterations)[cursor.GetParentNode()]++
		hasParentIterations := (*iterations)[cursor.GetParentNode()] < cursor.GetParentNode().Iterations
		if hasParentIterations || loopNode == cursor.GetParentNode() {
			cursor.MoveToFirstSibling()
			if cursor.GetCurrentNode().IsGroup() {
				cursor.MoveNext()
			}
		} else {
			cursor.ResetIterations(iterations)
			cursor.Up()
			return PlayMove(cursor, iterations, loopNode)
		}
	} else {
		cursor.MoveToSibling()
		cursor.ResetIterations(iterations)
	}
	return true
}

func PlayBeat(beatInterval time.Duration, pattern grid.Pattern, definition sequence.Sequence, errChan chan error) error {
	lines := definition.Lines

	for gridKey, note := range pattern {
		line := lines[gridKey.Line]
		if note.Ratchets.Length > 0 {
			ProcessRatchets(note, beatInterval, line, definition, errChan)
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
				PlayOnMessage(onMessage)
				PlayOffMessage(offMessage)
			case grid.MessageTypeCc:
				ccMessage := CCMessage(line, note, accents.Data, delay, true, definition.Instrument)

				PlayMessage(ccMessage.delay, ccMessage.MidiMessage())
			case grid.MessageTypeProgramChange:
				pcMessage := PCMessage(line, note, accents.Data, delay, true, definition.Instrument)
				PlayMessage(pcMessage.delay, pcMessage.MidiMessage())
			}
		}
	}

	return nil
}

func ProcessRatchets(note grid.Note, beatInterval time.Duration, line grid.LineDefinition, definition sequence.Sequence, errChan chan error) {
	for i := range note.Ratchets.Length + 1 {
		if note.Ratchets.HitAt(i) {
			shortGateLength := 20 * time.Millisecond
			ratchetInterval := time.Duration(i) * note.Ratchets.Interval(beatInterval)
			onMessage, offMessage := NoteMessages(line, definition.Accents.Data[note.AccentIndex].Value, shortGateLength, definition.Accents.Target, ratchetInterval)
			PlayOnMessage(onMessage)
			PlayOffMessage(offMessage)
		}
	}
}

func PlayMessage(delay time.Duration, message midi.Message) {
	if delay == 0 {
		playQueue <- message
	} else {
		time.AfterFunc(delay, func() {
			playQueue <- message
		})
	}
}

func PlayOnMessage(nm NoteMsg) {
	time.AfterFunc(nm.delay, func() {
		err := notereg.Add(nm)
		if err != nil {
			errChan <- err
		}
		playQueue <- nm.GetOffMidi()
	})
}

func PlayOffMessage(nm NoteMsg) {
	time.AfterFunc(nm.delay, func() {
		if notereg.RemoveID(nm) {
			playQueue <- nm.GetOffMidi()
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
		var ccValue uint8
		if cc.UpperLimit == 1 && note.AccentIndex > 4 {
			ccValue = 0
		} else if cc.UpperLimit == 1 {
			ccValue = 1
		} else {
			ccValue = uint8((float32((len(accents))-int(note.AccentIndex)) / float32(len(accents)-1)) * float32(cc.UpperLimit))
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

func GateLength(gateIndex int16, beatInterval time.Duration) time.Duration {
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
		shortGatesLen := int16(len(config.ShortGates))
		return time.Duration(float64(config.LongGates[gateIndex-shortGatesLen].Value) * float64(beatInterval))
	}
	return delay
}
