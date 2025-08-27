package playstate

import (
	"time"

	"github.com/chriserin/seq/internal/grid"
)

type PlayState struct {
	Playing            bool
	AllowAdvance       bool
	HasSolo            bool
	RecordPreRollBeats uint8
	PlayMode           PlayMode
	LoopMode           LoopMode
	BeatTime           time.Time
	LineStates         []LineState
}

type PlayMode uint8

const (
	PlayStandard PlayMode = iota
	PlayReceiver
)

type LoopMode uint8

const (
	OneTimeWholeSequence LoopMode = iota
	LoopWholeSequence
	LoopPart
	LoopOverlay
)

type GroupPlayState uint8

const (
	PlayStatePlay GroupPlayState = iota
	PlayStateMute
	PlayStateSolo
)

type LineState struct {
	Index               uint8
	CurrentBeat         uint8
	ResetLocation       uint8
	ResetActionLocation uint8
	ResetAction         grid.Action
	GroupPlayState      GroupPlayState
	Direction           int8
	ResetDirection      int8
}

func (ls LineState) IsMuted() bool {
	return ls.GroupPlayState == PlayStateMute
}

func (ls LineState) IsSolo() bool {
	return ls.GroupPlayState == PlayStateSolo
}

func (ls LineState) GridKey() grid.GridKey {
	return grid.GridKey{Line: ls.Index, Beat: ls.CurrentBeat}
}

func InitLineStates(lines int, previousPlayState []LineState, startBeat uint8) []LineState {
	linestates := make([]LineState, 0, lines)

	for i := range uint8(lines) {
		var previousGroupPlayState = PlayStatePlay
		if len(previousPlayState) > int(i) {
			previousState := previousPlayState[i]
			previousGroupPlayState = previousState.GroupPlayState
		}

		linestates = append(linestates, InitLineState(previousGroupPlayState, i, startBeat))
	}
	return linestates
}

func InitLineState(previousGroupPlayState GroupPlayState, index uint8, startBeat uint8) LineState {
	return LineState{
		Index:               index,
		CurrentBeat:         startBeat,
		Direction:           1,
		ResetDirection:      1,
		ResetLocation:       0,
		ResetActionLocation: 0,
		ResetAction:         0,
		GroupPlayState:      previousGroupPlayState,
	}
}

func (ls *LineState) AdvancePlayState(combinedPattern grid.Pattern, lineIndex int, beats uint8, lineStates []LineState) bool {

	advancedBeat := int8(ls.CurrentBeat) + ls.Direction
	currentState := *ls

	if advancedBeat >= int8(beats) || advancedBeat < 0 {
		// reset locations should be 1 time use.  Reset back to 0.
		if ls.ResetLocation != 0 && combinedPattern[grid.GK(uint8(lineIndex), currentState.ResetActionLocation)].Action == currentState.ResetAction {
			ls.CurrentBeat = currentState.ResetLocation
			advancedBeat = int8(currentState.ResetLocation)
		} else {
			ls.CurrentBeat = 0
			advancedBeat = int8(0)
		}
		ls.Direction = currentState.ResetDirection
		ls.ResetLocation = 0
	} else {
		ls.CurrentBeat = uint8(advancedBeat)
	}

	switch combinedPattern[grid.GK(uint8(lineIndex), uint8(advancedBeat))].Action {
	case grid.ActionNothing:
		return true
	case grid.ActionLineReset:
		ls.CurrentBeat = 0
	case grid.ActionLineReverse:
		ls.CurrentBeat = uint8(max(advancedBeat-2, 0))
		ls.Direction = -1
		ls.ResetLocation = uint8(max(advancedBeat-1, 0))
		ls.ResetActionLocation = uint8(advancedBeat)
		ls.ResetAction = grid.ActionLineReverse
	case grid.ActionLineBounce:
		ls.CurrentBeat = uint8(max(advancedBeat-1, 0))
		ls.Direction = -1
	case grid.ActionLineSkipBeat:
		ls.AdvancePlayState(combinedPattern, lineIndex, beats, lineStates)
	case grid.ActionLineDelay:
		ls.CurrentBeat = uint8(max(advancedBeat-1, 0))
	case grid.ActionLineResetAll:
		for i := range lineStates {
			lineStates[i].CurrentBeat = 0
			lineStates[i].Direction = 1
			lineStates[i].ResetLocation = 0
			lineStates[i].ResetDirection = 1
		}
		return false
	case grid.ActionLineBounceAll:
		for i := range lineStates {
			if i <= lineIndex {
				lineStates[i].CurrentBeat = uint8(max(lineStates[i].CurrentBeat-1, 0))
			}
			lineStates[i].Direction = -1
		}
		return false
	case grid.ActionLineSkipBeatAll:
		for i := range lineStates {
			if i <= lineIndex {
				ls.AdvancePlayState(combinedPattern, i, beats, lineStates)
			} else {
				ls.AdvancePlayState(combinedPattern, i, beats, lineStates)
				ls.AdvancePlayState(combinedPattern, i, beats, lineStates)
			}
		}
		return false
	}

	return true
}
