package main

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/operation"
	"github.com/stretchr/testify/assert"
)

func TestTempoChanges(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		initialTempo  int
		expectedTempo int
		description   string
	}{
		{
			name:          "Increase Tempo",
			commands:      []any{mappings.TempoInputSwitch, mappings.Increase},
			initialTempo:  120,
			expectedTempo: 121,
			description:   "Tempo should increase by 1",
		},
		{
			name:          "Increase Tempo At Boundary",
			commands:      []any{mappings.TempoInputSwitch, mappings.Increase},
			initialTempo:  300,
			expectedTempo: 300,
			description:   "Tempo should increase by 1",
		},
		{
			name:          "Decrease Tempo",
			commands:      []any{mappings.TempoInputSwitch, mappings.Decrease},
			initialTempo:  130,
			expectedTempo: 129,
			description:   "Tempo should decrease by 1",
		},
		{
			name:          "Decrease Tempo At Boundary",
			commands:      []any{mappings.TempoInputSwitch, mappings.Decrease},
			initialTempo:  30,
			expectedTempo: 30,
			description:   "Tempo should decrease by 1",
		},
		{
			name:          "Increase Tempo from arr view",
			commands:      []any{mappings.ToggleArrangementView, mappings.TempoInputSwitch, mappings.Increase},
			initialTempo:  120,
			expectedTempo: 121,
			description:   "Tempo should increase by 1",
		},
		{
			name:          "Decrease Tempo from arr view",
			commands:      []any{mappings.ToggleArrangementView, mappings.TempoInputSwitch, mappings.Decrease},
			initialTempo:  130,
			expectedTempo: 129,
			description:   "Tempo should decrease by 1",
		},
		{
			name:          "Increase Tempo by 5",
			commands:      []any{mappings.Increase},
			initialTempo:  120,
			expectedTempo: 125,
			description:   "Tempo should increase by 5",
		},
		{
			name:          "Increase Tempo by 5 at Boundary",
			commands:      []any{mappings.Increase},
			initialTempo:  297,
			expectedTempo: 300,
			description:   "Tempo should increase by 5",
		},
		{
			name:          "Decrease Tempo by 5",
			commands:      []any{mappings.Decrease},
			initialTempo:  130,
			expectedTempo: 125,
			description:   "Tempo should decrease by 5",
		},
		{
			name:          "Decrease Tempo by 5 at Boundary",
			commands:      []any{mappings.Decrease},
			initialTempo:  32,
			expectedTempo: 30,
			description:   "Tempo should decrease by 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.tempo = tt.initialTempo
					return *m
				},
			)
			assert.Equal(t, tt.initialTempo, m.definition.tempo, tt.description)
			m, _ = processCommands(tt.commands, m)
			assert.Equal(t, tt.expectedTempo, m.definition.tempo, tt.description)
		})
	}
}

func TestTempoInputSwitchEscapes(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedSelection operation.Selection
		expectedFocus     operation.Focus
		description       string
	}{
		{
			name:              "Tempo Input Switch Escape",
			commands:          []any{mappings.TempoInputSwitch, mappings.Escape},
			expectedSelection: operation.SelectGrid,
			expectedFocus:     operation.FocusGrid,
			description:       "Tempo input switch should select grid and escape should return to grid",
		},
		{
			name:              "Tempo Input Switch Escape from Arrangement View",
			commands:          []any{mappings.ToggleArrangementView, mappings.TempoInputSwitch},
			expectedSelection: operation.SelectTempo,
			expectedFocus:     operation.FocusArrangementEditor,
			description:       "Tempo input switch should select grid and escape should return to grid",
		},
		{
			name:              "Tempo Input Switch Escape from Arrangement View",
			commands:          []any{mappings.ToggleArrangementView, mappings.TempoInputSwitch, mappings.Escape},
			expectedSelection: operation.SelectGrid,
			expectedFocus:     operation.FocusArrangementEditor,
			description:       "Tempo input switch should select grid and escape should return to grid",
		},
		{
			name:              "Tempo Input Switch Enter",
			commands:          []any{mappings.TempoInputSwitch, mappings.Enter},
			expectedSelection: operation.SelectGrid,
			expectedFocus:     operation.FocusGrid,
			description:       "Tempo input switch should select grid and escape should return to grid",
		},
		{
			name:              "Tempo Input Switch Enter from Arrangement View",
			commands:          []any{mappings.ToggleArrangementView, mappings.TempoInputSwitch, mappings.Enter},
			expectedSelection: operation.SelectGrid,
			expectedFocus:     operation.FocusArrangementEditor,
			description:       "Tempo input switch should select grid and escape should return to grid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)
			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description)
		})
	}
}

func TestSubdivisionChanges(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []any
		initialSubdivisions  int
		expectedSubdivisions int
		description          string
	}{
		{
			name:                 "Increase Subdivisions",
			commands:             []any{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Increase},
			initialSubdivisions:  2,
			expectedSubdivisions: 3,
			description:          "Subdivisions should increase by 1",
		},
		{
			name:                 "Increase Subdivisions At Boundary",
			commands:             []any{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Increase},
			initialSubdivisions:  8,
			expectedSubdivisions: 8,
			description:          "Subdivisions should be at maximum",
		},
		{
			name:                 "Decrease Subdivisions",
			commands:             []any{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Decrease},
			initialSubdivisions:  3,
			expectedSubdivisions: 2,
			description:          "Subdivisions should decrease by 1",
		},
		{
			name:                 "Decrease Subdivisions At Boundary",
			commands:             []any{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Decrease},
			initialSubdivisions:  1,
			expectedSubdivisions: 1,
			description:          "Subdivisions should be at minimum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.subdivisions = tt.initialSubdivisions
					return *m
				},
			)
			assert.Equal(t, tt.initialSubdivisions, m.definition.subdivisions, tt.description)
			m, _ = processCommands(tt.commands, m)
			assert.Equal(t, tt.expectedSubdivisions, m.definition.subdivisions, tt.description)
		})
	}
}

func TestSetupInputSwitchWithIncrease(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []any
		initialSetupChannel  uint8
		expectedSetupChannel uint8
		description          string
	}{
		{
			name:                 "Setup Channel Increase",
			commands:             []any{mappings.SetupInputSwitch, mappings.Increase},
			initialSetupChannel:  0,
			expectedSetupChannel: 1,
			description:          "Setup input switch should select channel and increase should increment it",
		},
		{
			name:                 "Setup Channel Decrease",
			commands:             []any{mappings.SetupInputSwitch, mappings.Decrease},
			initialSetupChannel:  5,
			expectedSetupChannel: 4,
			description:          "Setup input switch should select channel and decrease should decrement it",
		},
		{
			name:                 "Setup Channel Increase At Upper Boundary",
			commands:             []any{mappings.SetupInputSwitch, mappings.Increase},
			initialSetupChannel:  16,
			expectedSetupChannel: 16,
			description:          "Setup input switch should select channel and increase should not go above 16",
		},
		{
			name:                 "Setup Channel Decrease At Lower Boundary",
			commands:             []any{mappings.SetupInputSwitch, mappings.Decrease},
			initialSetupChannel:  1,
			expectedSetupChannel: 1,
			description:          "Setup input switch should select channel and decrease should not go below 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.gridCursor.Line].Channel = tt.initialSetupChannel
					return *m
				},
			)

			assert.Equal(t, tt.initialSetupChannel, m.definition.lines[m.gridCursor.Line].Channel, "Initial setup channel should match")
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectSetupChannel, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedSetupChannel, m.definition.lines[m.gridCursor.Line].Channel, tt.description+" - channel value")
		})
	}
}

func TestSetupInputSwitchMessageTypeIncrease(t *testing.T) {
	tests := []struct {
		name                string
		commands            []any
		initialMessageType  grid.MessageType
		expectedMessageType grid.MessageType
		description         string
	}{
		{
			name:                "Message Type Increase from Note to CC",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialMessageType:  grid.MessageTypeNote,
			expectedMessageType: grid.MessageTypeCc,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Increase from CC to Program Change",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialMessageType:  grid.MessageTypeCc,
			expectedMessageType: grid.MessageTypeProgramChange,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Increase from Program Change to Note (wraparound)",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialMessageType:  grid.MessageTypeProgramChange,
			expectedMessageType: grid.MessageTypeNote,
			description:         "Two setup input switches should select message type and increase should wrap to Note",
		},
		{
			name:                "Message Type Decrease from Note to Program Change (wraparound)",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialMessageType:  grid.MessageTypeNote,
			expectedMessageType: grid.MessageTypeProgramChange,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Decrease from CC to Note",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialMessageType:  grid.MessageTypeCc,
			expectedMessageType: grid.MessageTypeNote,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Decrease from Program Change to CC",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialMessageType:  grid.MessageTypeProgramChange,
			expectedMessageType: grid.MessageTypeCc,
			description:         "Two setup input switches should select message type and increase should wrap to Note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.gridCursor.Line].MsgType = tt.initialMessageType
					return *m
				},
			)

			assert.Equal(t, tt.initialMessageType, m.definition.lines[m.gridCursor.Line].MsgType, "Initial message type should match")
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectSetupMessageType, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedMessageType, m.definition.lines[m.gridCursor.Line].MsgType, tt.description+" - message type value")
		})
	}
}

func TestSetupInputSwitchMessageTypeBackToGrid(t *testing.T) {
	tests := []struct {
		name                string
		commands            []any
		expectedMessageType grid.MessageType
		description         string
	}{
		{
			name:                "Message Type Increase from Note to Program Change and back to Grid",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.Increase, mappings.SetupInputSwitch},
			expectedMessageType: grid.MessageTypeProgramChange,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Increase from Note to Cc and back to Grid",
			commands:            []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.SetupInputSwitch},
			expectedMessageType: grid.MessageTypeCc,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Should select back to nothing")
			assert.Equal(t, tt.expectedMessageType, m.definition.lines[m.gridCursor.Line].MsgType, tt.description+" - message type value")
		})
	}
}

func TestSetupNoteChange(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		initialNote  uint8
		expectedNote uint8
		description  string
	}{
		{
			name:         "Note Increase",
			commands:     []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialNote:  60,
			expectedNote: 61,
			description:  "Note should increase by 1",
		},
		{
			name:         "Note Decrease",
			commands:     []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialNote:  60,
			expectedNote: 59,
			description:  "Note should decrease by 1",
		},
		{
			name:         "Note Increase At Boundary",
			commands:     []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialNote:  127,
			expectedNote: 127,
			description:  "Note should not increase beyond 127",
		},
		{
			name:         "Note Decrease at Boundary",
			commands:     []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialNote:  0,
			expectedNote: 0,
			description:  "Note should not decrease below 0",
		},
		{
			name:         "Note Decrease right above Boundary",
			commands:     []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialNote:  1,
			expectedNote: 0,
			description:  "Note should not decrease below 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.gridCursor.Line].Note = tt.initialNote
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedNote, m.definition.lines[m.gridCursor.Line].Note, tt.description)
		})
	}
}

func TestSetupCCChange(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		initialCc   uint8
		expectedCc  uint8
		description string
	}{
		{
			name:        "Should increase CC",
			commands:    []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Increase},
			initialCc:   0,
			expectedCc:  1,
			description: "Select CC and increment it",
		},
		{
			name:        "Should decrease CC",
			commands:    []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Decrease},
			initialCc:   1,
			expectedCc:  0,
			description: "Select CC and decrement it",
		},
		{
			name:        "Should skip unused CC on increase",
			commands:    []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Increase},
			initialCc:   2,
			expectedCc:  4,
			description: "Should skip over unused CCs on increase",
		},
		{
			name:        "Should skip unused CC on decrease",
			commands:    []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Decrease},
			initialCc:   4,
			expectedCc:  2,
			description: "Should skip over unused CCs on decrease",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.gridCursor.Line].Note = tt.initialCc
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedCc, m.definition.lines[m.gridCursor.Line].Note, tt.description)
		})
	}
}

func TestBeatInputSwitch(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		initialBeats  uint8
		expectedBeats uint8
		description   string
	}{
		{
			name:          "Beat Input Switch with Increase",
			commands:      []any{mappings.BeatInputSwitch, mappings.Increase},
			initialBeats:  16,
			expectedBeats: 17,
			description:   "Beat input switch should select beats and increase should increment it",
		},
		{
			name:          "Beat Input Switch with Decrease",
			commands:      []any{mappings.BeatInputSwitch, mappings.Decrease},
			initialBeats:  16,
			expectedBeats: 15,
			description:   "Beat input switch should select beats and decrease should decrement it",
		},
		{
			name:          "Beat Input Switch with Increase At Upper Boundary",
			commands:      []any{mappings.BeatInputSwitch, mappings.Increase},
			initialBeats:  127,
			expectedBeats: 127,
			description:   "Beat input switch should select beats and increase should not go above 127",
		},
		{
			name:          "Beat Input Switch with Decrease At Lower Boundary",
			commands:      []any{mappings.BeatInputSwitch, mappings.Decrease},
			initialBeats:  0,
			expectedBeats: 0,
			description:   "Beat input switch should select beats and decrease should not go below 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					(*m.definition.parts)[m.CurrentPartID()].Beats = tt.initialBeats
					return *m
				},
			)

			assert.Equal(t, tt.initialBeats, m.CurrentPart().Beats, "Initial beats should match")
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectBeats, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedBeats, m.CurrentPart().Beats, tt.description+" - beats value")
		})
	}
}

func TestBeatInputSwitchCyclesIncrease(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		initialCycles  int
		expectedCycles int
		description    string
	}{
		{
			name:           "Beat Input Switch Cycles with Increase",
			commands:       []any{mappings.CyclesInputSwitch, mappings.Increase},
			initialCycles:  4,
			expectedCycles: 5,
			description:    "Three beat input switches should select cycles and increase should increment it",
		},
		{
			name:           "Beat Input Switch Cycles with Decrease",
			commands:       []any{mappings.CyclesInputSwitch, mappings.Decrease},
			initialCycles:  4,
			expectedCycles: 3,
			description:    "Three beat input switches should select cycles and decrease should decrement it",
		},
		{
			name:           "Beat Input Switch Cycles with Increase At Upper Boundary",
			commands:       []any{mappings.CyclesInputSwitch, mappings.Increase},
			initialCycles:  127,
			expectedCycles: 127,
			description:    "Three beat input switches should select cycles and increase should not go above 127",
		},
		{
			name:           "Beat Input Switch Cycles with Decrease At Lower Boundary",
			commands:       []any{mappings.CyclesInputSwitch, mappings.Decrease},
			initialCycles:  0,
			expectedCycles: 0,
			description:    "Three beat input switches should select cycles and decrease should not go below 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					currentNode := m.arrangement.Cursor.GetCurrentNode()
					currentNode.Section.Cycles = tt.initialCycles
					return *m
				},
			)

			assert.Equal(t, tt.initialCycles, m.CurrentSongSection().Cycles, "Initial cycles should match")
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectCycles, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedCycles, m.CurrentSongSection().Cycles, tt.description+" - cycles value")
		})
	}
}

func TestBeatInputSwitchStartBeatsIncrease(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		initialStartBeats  int
		expectedStartBeats int
		description        string
	}{
		{
			name:               "Beat Input Switch StartBeats with Increase",
			commands:           []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Increase},
			initialStartBeats:  8,
			expectedStartBeats: 9,
			description:        "Two beat input switches should select start beats and increase should increment it",
		},
		{
			name:               "Beat Input Switch StartBeats with Decrease",
			commands:           []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Decrease},
			initialStartBeats:  8,
			expectedStartBeats: 7,
			description:        "Two beat input switches should select start beats and decrease should decrement it",
		},
		{
			name:     "Beat Input Switch StartBeats with Increase At Upper Boundary",
			commands: []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Increase},
			// NOTE: The upper boundary for start beats is the number of beats in a part
			initialStartBeats:  31,
			expectedStartBeats: 31,
			description:        "Two beat input switches should select start beats and increase should not go above 127",
		},
		{
			name:               "Beat Input Switch StartBeats with Decrease At Lower Boundary",
			commands:           []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Decrease},
			initialStartBeats:  0,
			expectedStartBeats: 0,
			description:        "Two beat input switches should select start beats and decrease should not go below 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					currentNode := m.arrangement.Cursor.GetCurrentNode()
					currentNode.Section.StartBeat = tt.initialStartBeats
					return *m
				},
			)

			assert.Equal(t, tt.initialStartBeats, m.CurrentSongSection().StartBeat, "Initial start beats should match")
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectStartBeats, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedStartBeats, m.CurrentSongSection().StartBeat, tt.description+" - start beats value")
		})
	}
}

func TestBeatInputSwitchStartCyclesIncrease(t *testing.T) {
	tests := []struct {
		name                string
		commands            []any
		initialStartCycles  int
		expectedStartCycles int
		description         string
	}{
		{
			name:                "Beat Input Switch StartCycles with Increase",
			commands:            []any{mappings.CyclesInputSwitch, mappings.CyclesInputSwitch, mappings.Increase},
			initialStartCycles:  4,
			expectedStartCycles: 5,
			description:         "Two beat input switches should select start cycles and increase should increment it",
		},
		{
			name:                "Beat Input Switch StartCycles with Decrease",
			commands:            []any{mappings.CyclesInputSwitch, mappings.CyclesInputSwitch, mappings.Decrease},
			initialStartCycles:  4,
			expectedStartCycles: 3,
			description:         "Two beat input switches should select start cycles and decrease should decrement it",
		},
		{
			name:                "Beat Input Switch StartCycles with Increase At Upper Boundary",
			commands:            []any{mappings.CyclesInputSwitch, mappings.CyclesInputSwitch, mappings.Increase},
			initialStartCycles:  127,
			expectedStartCycles: 127,
			description:         "Two beat input switches should select start cycles and increase should not go above 127",
		},
		{
			name:                "Beat Input Switch StartCycles with Decrease At Lower Boundary",
			commands:            []any{mappings.CyclesInputSwitch, mappings.CyclesInputSwitch, mappings.Decrease},
			initialStartCycles:  0,
			expectedStartCycles: 0,
			description:         "Two beat input switches should select start cycles and decrease should not go below 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					currentNode := m.arrangement.Cursor.GetCurrentNode()
					currentNode.Section.StartCycles = tt.initialStartCycles
					return *m
				},
			)

			assert.Equal(t, tt.initialStartCycles, m.CurrentSongSection().StartCycles, "Initial start cycles should match")
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectStartCycles, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedStartCycles, m.CurrentSongSection().StartCycles, tt.description+" - start cycles value")
		})
	}
}

func TestNewLine(t *testing.T) {
	tests := []struct {
		name                string
		commands            []any
		initialLineCount    int
		expectedLineCount   int
		expectedLastChannel uint8
		expectedLastNote    uint8
		description         string
	}{
		{
			name:                "Add new line with default values",
			commands:            []any{mappings.NewLine},
			initialLineCount:    8,
			expectedLineCount:   9,
			expectedLastChannel: 10,
			expectedLastNote:    68,
			description:         "Should add a new line with channel from last line and note+1",
		},
		{
			name:                "Add new line with custom last line values",
			commands:            []any{mappings.NewLine},
			initialLineCount:    8,
			expectedLineCount:   9,
			expectedLastChannel: 5,
			expectedLastNote:    73,
			description:         "Should add a new line with channel 5 and note 73",
		},
		{
			name:                "Add multiple new lines",
			commands:            []any{mappings.NewLine, mappings.NewLine, mappings.NewLine},
			initialLineCount:    8,
			expectedLineCount:   11,
			expectedLastChannel: 10,
			expectedLastNote:    70,
			description:         "Should add three new lines with incrementing note values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					// Set up the last line with specific values for the second test
					if tt.name == "Add new line with custom last line values" {
						m.definition.lines[len(m.definition.lines)-1].Channel = 5
						m.definition.lines[len(m.definition.lines)-1].Note = 72
					}
					return *m
				},
			)

			// Verify initial state
			assert.Equal(t, tt.initialLineCount, len(m.definition.lines), "Initial line count should match")

			// Execute commands
			m, _ = processCommands(tt.commands, m)

			// Verify final state
			assert.Equal(t, tt.expectedLineCount, len(m.definition.lines), tt.description+" - line count")

			if tt.expectedLineCount > tt.initialLineCount {
				// Check the last line properties
				lastLine := m.definition.lines[len(m.definition.lines)-1]
				assert.Equal(t, tt.expectedLastChannel, lastLine.Channel, tt.description+" - last line channel")
				assert.Equal(t, tt.expectedLastNote, lastLine.Note, tt.description+" - last line note")

				// Check that the new line has default message type
				assert.Equal(t, grid.MessageTypeNote, lastLine.MsgType, tt.description+" - last line message type should be default")
			}
		})
	}
}

func TestMuteAndSolo(t *testing.T) {
	tests := []struct {
		name                   string
		commands               []any
		cursorLine             uint8
		expectedGroupPlayState groupPlayState
		expectedHasSolo        bool
		description            string
	}{
		{
			name:                   "Mute line from play state",
			commands:               []any{mappings.Mute},
			cursorLine:             0,
			expectedGroupPlayState: PlayStateMute,
			expectedHasSolo:        false,
			description:            "Should mute line 0 from play state",
		},
		{
			name:                   "Unmute line back to play state",
			commands:               []any{mappings.Mute, mappings.Mute},
			cursorLine:             0,
			expectedGroupPlayState: PlayStatePlay,
			expectedHasSolo:        false,
			description:            "Should unmute line 0 back to play state",
		},
		{
			name:                   "Solo line from play state",
			commands:               []any{mappings.Solo},
			cursorLine:             0,
			expectedGroupPlayState: PlayStateSolo,
			expectedHasSolo:        true,
			description:            "Should solo line 0 from play state",
		},
		{
			name:                   "Unsolo line back to play state",
			commands:               []any{mappings.Solo, mappings.Solo},
			cursorLine:             0,
			expectedGroupPlayState: PlayStatePlay,
			expectedHasSolo:        false,
			description:            "Should unsolo line 0 back to play state",
		},
		{
			name:                   "Solo line from muted state",
			commands:               []any{mappings.Mute, mappings.Solo},
			cursorLine:             0,
			expectedGroupPlayState: PlayStateSolo,
			expectedHasSolo:        true,
			description:            "Should solo line 0 from muted state",
		},
		{
			name:                   "Mute line from solo state",
			commands:               []any{mappings.Solo, mappings.Mute},
			cursorLine:             0,
			expectedGroupPlayState: PlayStateMute,
			expectedHasSolo:        false,
			description:            "Should mute line 0 from solo state",
		},
		{
			name:                   "Mute different line",
			commands:               []any{mappings.CursorDown, mappings.Mute},
			cursorLine:             1,
			expectedGroupPlayState: PlayStateMute,
			expectedHasSolo:        false,
			description:            "Should mute line 1 after moving cursor down",
		},
		{
			name:                   "Solo different line",
			commands:               []any{mappings.CursorDown, mappings.CursorDown, mappings.Solo},
			cursorLine:             2,
			expectedGroupPlayState: PlayStateSolo,
			expectedHasSolo:        true,
			description:            "Should solo line 2 after moving cursor down twice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.playState = playState{lineStates: make([]linestate, len(m.definition.lines))}
					for i := range m.playState.lineStates {
						m.playState.lineStates[i] = InitLineState(PlayStatePlay, uint8(i), 0)
					}
					return *m
				},
			)

			assert.Equal(t, PlayStatePlay, m.playState.lineStates[tt.cursorLine].groupPlayState, "Initial play state should be PlayStatePlay")
			assert.False(t, m.playState.hasSolo, "Initial hasSolo should be false")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedGroupPlayState, m.playState.lineStates[tt.cursorLine].groupPlayState, tt.description+" - groupPlayState should match")
			assert.Equal(t, tt.expectedHasSolo, m.playState.hasSolo, tt.description+" - hasSolo should match")

			switch tt.expectedGroupPlayState {
			case PlayStateMute:
				assert.True(t, m.playState.lineStates[tt.cursorLine].IsMuted(), tt.description+" - IsMuted() should return true")
				assert.False(t, m.playState.lineStates[tt.cursorLine].IsSolo(), tt.description+" - IsSolo() should return false")
			case PlayStateSolo:
				assert.False(t, m.playState.lineStates[tt.cursorLine].IsMuted(), tt.description+" - IsMuted() should return false")
				assert.True(t, m.playState.lineStates[tt.cursorLine].IsSolo(), tt.description+" - IsSolo() should return true")
			default:
				assert.False(t, m.playState.lineStates[tt.cursorLine].IsMuted(), tt.description+" - IsMuted() should return false")
				assert.False(t, m.playState.lineStates[tt.cursorLine].IsSolo(), tt.description+" - IsSolo() should return false")
			}
		})
	}
}
