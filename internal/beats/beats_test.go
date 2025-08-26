package beats

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/playstate"
	"github.com/chriserin/seq/internal/sequence"
	"github.com/stretchr/testify/assert"
)

func TestLoop(t *testing.T) {
	sequence, cursor := SimpleSequence()

	(*sequence.Parts)[0].Beats = 1
	sequence.Arrangement.ResetAllPlayCycles()

	beatsPlayed := PlayBeats(sequence, cursor, 4)
	assert.Equal(t, 1, beatsPlayed)
}

func PlayBeats(sequence sequence.Sequence, cursor arrangement.ArrCursor, limit int) int {
	testMessageChan := make(chan ModelPlayedMsg)
	beatsPlayedCounter := 0
	var update = ModelPlayedMsg{PlayState: playstate.PlayState{Playing: true}}
	sendFn := func(msg tea.Msg) {
		update = msg.(ModelPlayedMsg)
		testMessageChan <- update
	}

	updateChannel := GetUpdateChannel()
	beatChannel := GetBeatChannel()

	Loop(sendFn)

	updateChannel <- ModelMsg{
		PlayState: playstate.PlayState{Playing: true, LineStates: playstate.InitLineStates(1, []playstate.LineState{}, 0)},
		Sequence:  sequence,
		Cursor:    cursor,
	}

	for update.PlayState.Playing && beatsPlayedCounter < limit {
		beatChannel <- BeatMsg{Interval: 0}
		update = <-testMessageChan
		if update.PerformStop {
			break
		} else {
			beatsPlayedCounter++
		}
		updateChannel <- ModelMsg{PlayState: update.PlayState, Sequence: update.Definition, Cursor: update.Cursor}
	}
	doneChannel := GetDoneChannel()
	doneChannel <- struct{}{}
	return beatsPlayedCounter
}

func SimpleSequence() (sequence.Sequence, arrangement.ArrCursor) {
	var parts = sequence.InitParts()

	nodeA := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	root := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      make([]*arrangement.Arrangement, 0),
	}

	root.Nodes = append(root.Nodes, nodeA)

	testSequence := sequence.Sequence{
		Arrangement: root,
		Parts:       &parts,
		Keyline:     0,
		Lines:       make([]grid.LineDefinition, 1),
	}

	return testSequence, arrangement.ArrCursor{root, nodeA}
}
