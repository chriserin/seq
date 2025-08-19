package main

import (
	"fmt"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/notereg"
	"github.com/chriserin/seq/internal/overlays"
	"github.com/chriserin/seq/internal/seqmidi"
	midi "gitlab.com/gomidi/midi/v2"
)

type modelMsg struct {
	playState  playState
	definition Definition
	cursor     arrangement.ArrCursor
	midiSendFn seqmidi.SendFunc
}

type modelPlayedMsg struct {
	performStop bool
	playState   playState
	definition  Definition
	cursor      arrangement.ArrCursor
}

func StartBeatLoop(updateChannel chan modelMsg, midiLoopChannel chan beatMsg, sendFn func(tea.Msg)) {
	go func() {
		var playState playState
		var definition Definition
		var cursor arrangement.ArrCursor
		var midiSendFn seqmidi.SendFunc
		var errChan chan error
		for {
			select {
			case modelMsg := <-updateChannel:
				playState = modelMsg.playState
				definition = modelMsg.definition
				cursor = modelMsg.cursor
				midiSendFn = modelMsg.midiSendFn
			case beatMsg := <-midiLoopChannel:
				Beat(beatMsg, playState, definition, cursor, midiSendFn, sendFn, errChan)
			case err := <-errChan:
				fmt.Println(err)
			}
		}
	}()
}

func Beat(msg beatMsg, playState playState, definition Definition, cursor arrangement.ArrCursor, midiSendFn seqmidi.SendFunc, sendFn func(tea.Msg), errChan chan error) {

	currentSection := &cursor[len(cursor)-1].Section
	partID := currentSection.Part
	currentPart := (*definition.parts)[partID]
	playingOverlay := currentPart.Overlays.HighestMatchingOverlay(currentSection.PlayCycles())

	if playState.playing && playState.recordPreRollBeats == 0 {
		// NOTE: Only advance if we've already played the first beat.
		if playState.allowAdvance {
			advanceCurrentBeat(*currentSection, playingOverlay, playState.lineStates, currentPart.Beats)
			advanceKeyCycle(definition.keyline, playState.lineStates, playState.loopMode, currentSection)
			if currentSection.IsDone() && playState.loopMode != LoopOverlay {
				if PlayMove(&cursor) {
					cursor[len(cursor)-1].Section.DuringPlayReset()
					currentSection = &cursor[len(cursor)-1].Section
					playState.lineStates = InitLineStates(len(definition.lines), playState.lineStates, uint8(cursor[len(cursor)-1].Section.StartBeat))
				} else {
					playState.lineStates = InitLineStates(len(definition.lines), playState.lineStates, 0)
					sendFn(modelPlayedMsg{performStop: true, playState: playState, definition: definition, cursor: cursor})
					return
				}
			}
		}
	}

	if playState.playing {
		partID = currentSection.Part
		currentPart = (*definition.parts)[partID]
		playingOverlay = currentPart.Overlays.HighestMatchingOverlay(currentSection.PlayCycles())
		gridKeys := make([]grid.GridKey, 0, len(playState.lineStates))
		CurrentBeatGridKeys(&gridKeys, playState.lineStates, playState.hasSolo)

		pattern := make(grid.Pattern)
		playingOverlay.CurrentBeatOverlayPattern(&pattern, currentSection.PlayCycles(), gridKeys)

		if playState.recordPreRollBeats > 0 {
			playState.recordPreRollBeats--
			return
		}

		err := PlayBeat(msg.interval, pattern, definition, midiSendFn, errChan)

		if !playState.allowAdvance {
			playState.allowAdvance = true
		}

		if err != nil {
			fault.Wrap(err, fmsg.With("error when playing beat"))
		}
		sendFn(modelPlayedMsg{playState: playState, definition: definition, cursor: cursor})
	}
}

func CurrentBeatGridKeys(gridKeys *[]grid.GridKey, lineStates []linestate, hasSolo bool) {
	for _, linestate := range lineStates {
		if linestate.IsSolo() || (!linestate.IsMuted() && !hasSolo) {
			*gridKeys = append(*gridKeys, linestate.GridKey())
		}
	}
}

func advanceCurrentBeat(currentSection arrangement.SongSection, playingOverlay *overlays.Overlay, lineStates []linestate, partBeats uint8) {
	pattern := make(grid.Pattern)
	playingOverlay.CombineActionPattern(&pattern, currentSection.PlayCycles())
	for i := range lineStates {
		doContinue := lineStates[i].advancePlayState(pattern, i, partBeats, lineStates)
		if !doContinue {
			break
		}
	}
}

func (ls *linestate) advancePlayState(combinedPattern grid.Pattern, lineIndex int, beats uint8, lineStates []linestate) bool {

	advancedBeat := int8(ls.currentBeat) + ls.direction
	currentState := *ls

	if advancedBeat >= int8(beats) || advancedBeat < 0 {
		// reset locations should be 1 time use.  Reset back to 0.
		if ls.resetLocation != 0 && combinedPattern[GK(uint8(lineIndex), currentState.resetActionLocation)].Action == currentState.resetAction {
			ls.currentBeat = currentState.resetLocation
			advancedBeat = int8(currentState.resetLocation)
		} else {
			ls.currentBeat = 0
			advancedBeat = int8(0)
		}
		ls.direction = currentState.resetDirection
		ls.resetLocation = 0
	} else {
		ls.currentBeat = uint8(advancedBeat)
	}

	switch combinedPattern[GK(uint8(lineIndex), uint8(advancedBeat))].Action {
	case grid.ActionNothing:
		return true
	case grid.ActionLineReset:
		ls.currentBeat = 0
	case grid.ActionLineReverse:
		ls.currentBeat = uint8(max(advancedBeat-2, 0))
		ls.direction = -1
		ls.resetLocation = uint8(max(advancedBeat-1, 0))
		ls.resetActionLocation = uint8(advancedBeat)
		ls.resetAction = grid.ActionLineReverse
	case grid.ActionLineBounce:
		ls.currentBeat = uint8(max(advancedBeat-1, 0))
		ls.direction = -1
	case grid.ActionLineSkipBeat:
		ls.advancePlayState(combinedPattern, lineIndex, beats, lineStates)
	case grid.ActionLineDelay:
		ls.currentBeat = uint8(max(advancedBeat-1, 0))
	case grid.ActionLineResetAll:
		for i := range lineStates {
			lineStates[i].currentBeat = 0
			lineStates[i].direction = 1
			lineStates[i].resetLocation = 0
			lineStates[i].resetDirection = 1
		}
		return false
	case grid.ActionLineBounceAll:
		for i := range lineStates {
			if i <= lineIndex {
				lineStates[i].currentBeat = uint8(max(lineStates[i].currentBeat-1, 0))
			}
			lineStates[i].direction = -1
		}
		return false
	case grid.ActionLineSkipBeatAll:
		for i := range lineStates {
			if i <= lineIndex {
				ls.advancePlayState(combinedPattern, i, beats, lineStates)
			} else {
				ls.advancePlayState(combinedPattern, i, beats, lineStates)
				ls.advancePlayState(combinedPattern, i, beats, lineStates)
			}
		}
		return false
	}

	return true
}

func advanceKeyCycle(keyline uint8, lineStates []linestate, loopMode LoopMode, section *arrangement.SongSection) {
	if lineStates[keyline].currentBeat == 0 && loopMode != LoopOverlay {
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

func PlayBeat(beatInterval time.Duration, pattern grid.Pattern, definition Definition, midiSendFn seqmidi.SendFunc, errChan chan error) error {
	lines := definition.lines

	for gridKey, note := range pattern {
		line := lines[gridKey.Line]
		if note.Ratchets.Length > 0 {
			ProcessRatchets(note, beatInterval, line, definition, midiSendFn, errChan)
		} else if note != zeronote {
			accents := definition.accents

			delay := Delay(note.WaitIndex, beatInterval)
			gateLength := GateLength(note.GateIndex, beatInterval)

			switch line.MsgType {
			case grid.MessageTypeNote:
				onMessage, offMessage := NoteMessages(
					line,
					definition.accents.Data[note.AccentIndex].Value,
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
				ccMessage := CCMessage(line, note, accents.Data, delay, true, definition.instrument)
				err := ProcessNoteMsg(ccMessage, midiSendFn, errChan)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process cc msg"))
				}
			case grid.MessageTypeProgramChange:
				pcMessage := PCMessage(line, note, accents.Data, delay, true, definition.instrument)
				err := ProcessNoteMsg(pcMessage, midiSendFn, errChan)
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot process cc msg"))
				}
			}
		}
	}

	return nil
}

func ProcessRatchets(note grid.Note, beatInterval time.Duration, line grid.LineDefinition, definition Definition, midiSendFn seqmidi.SendFunc, errChan chan error) error {
	for i := range note.Ratchets.Length + 1 {
		if note.Ratchets.HitAt(i) {
			shortGateLength := 20 * time.Millisecond
			ratchetInterval := time.Duration(i) * note.Ratchets.Interval(beatInterval)
			onMessage, offMessage := NoteMessages(line, definition.accents.Data[note.AccentIndex].Value, shortGateLength, definition.accents.Target, ratchetInterval)
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
	case noteMsg:
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
