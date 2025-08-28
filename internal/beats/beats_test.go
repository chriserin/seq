package beats

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/playstate"
	"github.com/chriserin/seq/internal/seqmidi"
	"github.com/chriserin/seq/internal/sequence"
	"github.com/stretchr/testify/assert"
)

func TestSimpleSequenceBeats(t *testing.T) {
	tests := []struct {
		name                string
		partBeats           uint8
		expectedBeatsPlayed int
	}{
		{"Part with 1 beat", 1, 1},
		{"Part with 3 beats", 3, 3},
		{"Part with 7 beats", 7, 7},
		{"Part with 13 beats", 13, 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sequence, cursor := SimpleSequence()

			(*sequence.Parts)[0].Beats = tt.partBeats

			beatsPlayed := PlayBeats(sequence, cursor, int(tt.partBeats)+3)
			assert.Equal(t, tt.expectedBeatsPlayed, beatsPlayed)
		})
	}
}

func TestGroupedSequenceBeats(t *testing.T) {
	tests := []struct {
		name                string
		partBeats           uint8
		groupIterations     int
		expectedBeatsPlayed int
	}{
		{"Part with 1 beat", 1, 1, 1},
		{"Part with 2 beat", 2, 2, 4},
		{"Part with 3 beats", 3, 3, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sequence, cursor := SimpleGroupedSequence()

			(*sequence.Parts)[0].Beats = tt.partBeats
			cursor[1].Iterations = tt.groupIterations

			beatsPlayed := PlayBeats(sequence, cursor, tt.expectedBeatsPlayed+3)
			assert.Equal(t, tt.expectedBeatsPlayed, beatsPlayed)
		})
	}
}

func TestSiblingSections(t *testing.T) {
	tests := []struct {
		name                string
		partABeats          uint8
		partBBeats          uint8
		expectedBeatsPlayed int
	}{
		{"Parts with 1 beat", 1, 1, 2},
		{"Parts with 2 beats", 2, 2, 4},
		{"Parts with different beats 1/2", 1, 2, 3},
		{"Parts with different beats 2/1", 2, 1, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sequence, cursor := SiblingSectionSequence()

			(*sequence.Parts)[0].Beats = tt.partABeats
			(*sequence.Parts)[1].Beats = tt.partBBeats

			beatsPlayed := PlayBeats(sequence, cursor, tt.expectedBeatsPlayed+3)
			assert.Equal(t, tt.expectedBeatsPlayed, beatsPlayed)
		})
	}
}

func TestNestedGroups(t *testing.T) {
	tests := []struct {
		name                string
		partBeats           uint8
		groupAIterations    int
		groupBIterations    int
		expectedBeatsPlayed int
	}{
		{"Part with 1 beat", 1, 2, 2, 4},
		{"Part with 2 beats", 2, 2, 2, 8},
		{"Part with 3 beats", 3, 2, 2, 12},
		{"Part with 3 beats and different iterations", 3, 2, 3, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sequence, cursor := NestedGroupsSequence()

			(*sequence.Parts)[0].Beats = tt.partBeats
			cursor[1].Iterations = tt.groupAIterations
			cursor[2].Iterations = tt.groupBIterations

			beatsPlayed := PlayBeats(sequence, cursor, tt.expectedBeatsPlayed+3)
			assert.Equal(t, tt.expectedBeatsPlayed, beatsPlayed)
		})
	}
}

func TestGroupPartSiblingSequence(t *testing.T) {
	tests := []struct {
		name                string
		partABeats          uint8
		partBBeats          uint8
		groupIterations     int
		expectedBeatsPlayed int
	}{
		{"Parts with 1 beat", 1, 1, 1, 2},
		{"Parts with 2 beats", 2, 2, 2, 6},
		{"Parts with different beats 1/2", 1, 2, 2, 4},
		{"Parts with different beats 2/1", 2, 1, 2, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sequence, cursor := GroupPartSiblingSequence()

			(*sequence.Parts)[0].Beats = tt.partABeats
			(*sequence.Parts)[1].Beats = tt.partBBeats
			sequence.Arrangement.Nodes[0].Iterations = tt.groupIterations

			beatsPlayed := PlayBeats(sequence, cursor, tt.expectedBeatsPlayed+3)
			assert.Equal(t, tt.expectedBeatsPlayed, beatsPlayed)
		})
	}
}

func TestPartGroupSiblingSequence(t *testing.T) {
	tests := []struct {
		name                string
		partABeats          uint8
		partBBeats          uint8
		groupIterations     int
		expectedBeatsPlayed int
	}{
		{"Parts with 1 beat", 1, 1, 1, 2},
		{"Parts with 2 beats", 1, 1, 2, 3},
		{"Parts with different beats 1/2", 1, 2, 2, 4},
		{"Parts with different beats 2/1", 2, 1, 2, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sequence, cursor := PartGroupSiblingSequence()

			(*sequence.Parts)[0].Beats = tt.partABeats
			(*sequence.Parts)[1].Beats = tt.partBBeats
			sequence.Arrangement.Nodes[1].Iterations = tt.groupIterations

			beatsPlayed := PlayBeats(sequence, cursor, tt.expectedBeatsPlayed+3)
			assert.Equal(t, tt.expectedBeatsPlayed, beatsPlayed)
		})
	}
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
	iterations := make(map[*arrangement.Arrangement]int)
	playstate.BuildIterationsMap(sequence.Arrangement, &iterations)
	playState := playstate.PlayState{Playing: true, LineStates: playstate.InitLineStates(1, []playstate.LineState{}, 0), Iterations: &iterations}

	updateChannel <- ModelMsg{
		PlayState:      playState,
		Sequence:       sequence,
		Cursor:         cursor,
		MidiConnection: seqmidi.MidiConnection{Test: true},
	}

	for update.PlayState.Playing && beatsPlayedCounter < limit {
		beatChannel <- BeatMsg{Interval: 0}
		update = <-testMessageChan
		if update.PerformStop {
			break
		} else {
			beatsPlayedCounter++
		}
		updateChannel <- ModelMsg{PlayState: update.PlayState, Sequence: update.Definition, Cursor: update.Cursor, MidiConnection: seqmidi.MidiConnection{Test: true}}
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

func SiblingSectionSequence() (sequence.Sequence, arrangement.ArrCursor) {
	var parts = sequence.InitParts()
	parts = append(parts, arrangement.InitPart("Part 2"))

	nodeA := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 1, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	root := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      make([]*arrangement.Arrangement, 0),
	}

	root.Nodes = append(root.Nodes, nodeA, nodeB)

	testSequence := sequence.Sequence{
		Arrangement: root,
		Parts:       &parts,
		Keyline:     0,
		Lines:       make([]grid.LineDefinition, 1),
	}

	return testSequence, arrangement.ArrCursor{root, nodeA}
}

func SimpleGroupedSequence() (sequence.Sequence, arrangement.ArrCursor) {
	var parts = sequence.InitParts()

	nodeA := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	groupA := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      []*arrangement.Arrangement{nodeA},
	}

	root := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      make([]*arrangement.Arrangement, 0),
	}

	root.Nodes = append(root.Nodes, groupA)

	testSequence := sequence.Sequence{
		Arrangement: root,
		Parts:       &parts,
		Keyline:     0,
		Lines:       make([]grid.LineDefinition, 1),
	}

	return testSequence, arrangement.ArrCursor{root, groupA, nodeA}
}

func NestedGroupsSequence() (sequence.Sequence, arrangement.ArrCursor) {
	var parts = sequence.InitParts()

	nodeA := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	groupA := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      []*arrangement.Arrangement{nodeA},
	}

	groupB := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      []*arrangement.Arrangement{groupA},
	}

	root := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      make([]*arrangement.Arrangement, 0),
	}

	root.Nodes = append(root.Nodes, groupB)

	testSequence := sequence.Sequence{
		Arrangement: root,
		Parts:       &parts,
		Keyline:     0,
		Lines:       make([]grid.LineDefinition, 1),
	}

	return testSequence, arrangement.ArrCursor{root, groupB, groupA, nodeA}
}

func GroupPartSiblingSequence() (sequence.Sequence, arrangement.ArrCursor) {
	var parts = sequence.InitParts()
	parts = append(parts, arrangement.InitPart("Part 2"))

	nodeA := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 1, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	groupA := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      []*arrangement.Arrangement{nodeA},
	}

	root := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      make([]*arrangement.Arrangement, 0),
	}

	root.Nodes = append(root.Nodes, groupA, nodeB)

	testSequence := sequence.Sequence{
		Arrangement: root,
		Parts:       &parts,
		Keyline:     0,
		Lines:       make([]grid.LineDefinition, 1),
	}

	return testSequence, arrangement.ArrCursor{root, groupA, nodeA}
}

func PartGroupSiblingSequence() (sequence.Sequence, arrangement.ArrCursor) {
	var parts = sequence.InitParts()
	parts = append(parts, arrangement.InitPart("Part 2"))

	nodeA := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &arrangement.Arrangement{
		Section:    arrangement.SongSection{Part: 1, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	groupA := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      []*arrangement.Arrangement{nodeA},
	}

	root := &arrangement.Arrangement{
		Iterations: 1,
		Nodes:      make([]*arrangement.Arrangement, 0),
	}

	root.Nodes = append(root.Nodes, nodeB, groupA)

	testSequence := sequence.Sequence{
		Arrangement: root,
		Parts:       &parts,
		Keyline:     0,
		Lines:       make([]grid.LineDefinition, 1),
	}

	return testSequence, arrangement.ArrCursor{root, nodeB}
}
