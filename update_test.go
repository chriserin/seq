package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/seqmidi"
	"github.com/stretchr/testify/assert"
)

// TestUpdateCursorMovements tests all cursor movement keybindings
func TestUpdateCursorMovements(t *testing.T) {
	tests := []struct {
		name        string
		command     mappings.Command
		initialPos  grid.GridKey
		expectedPos grid.GridKey
		description string
	}{
		{
			name:        "Cursor Up",
			command:     mappings.CursorUp,
			initialPos:  grid.GridKey{Line: 5, Beat: 0},
			expectedPos: grid.GridKey{Line: 4, Beat: 0},
			description: "Cursor should move up one line",
		},
		{
			name:        "Cursor Up At Boundary",
			command:     mappings.CursorUp,
			initialPos:  grid.GridKey{Line: 0, Beat: 0},
			expectedPos: grid.GridKey{Line: 0, Beat: 0},
			description: "Cursor should move up one line",
		},
		{
			name:        "Cursor Down",
			command:     mappings.CursorDown,
			initialPos:  grid.GridKey{Line: 2, Beat: 0},
			expectedPos: grid.GridKey{Line: 3, Beat: 0},
			description: "Cursor should move down one line",
		},
		{
			name:        "Cursor Down At Boundary",
			command:     mappings.CursorDown,
			initialPos:  grid.GridKey{Line: 7, Beat: 0},
			expectedPos: grid.GridKey{Line: 7, Beat: 0},
			description: "Cursor should move down one line",
		},
		{
			name:        "Cursor Left",
			command:     mappings.CursorLeft,
			initialPos:  grid.GridKey{Line: 0, Beat: 5},
			expectedPos: grid.GridKey{Line: 0, Beat: 4},
			description: "Cursor should move left one beat",
		},
		{
			name:        "Cursor Left At Boundary",
			command:     mappings.CursorLeft,
			initialPos:  grid.GridKey{Line: 0, Beat: 0},
			expectedPos: grid.GridKey{Line: 0, Beat: 0},
			description: "Cursor should move left one beat",
		},
		{
			name:        "Cursor Right",
			command:     mappings.CursorRight,
			initialPos:  grid.GridKey{Line: 0, Beat: 3},
			expectedPos: grid.GridKey{Line: 0, Beat: 4},
			description: "Cursor should move right one beat",
		},
		{
			name:        "Cursor Right At Boundary",
			command:     mappings.CursorRight,
			initialPos:  grid.GridKey{Line: 0, Beat: 31},
			expectedPos: grid.GridKey{Line: 0, Beat: 31},
			description: "Cursor should move right one beat",
		},
		{
			name:        "Cursor Line Start",
			command:     mappings.CursorLineStart,
			initialPos:  grid.GridKey{Line: 2, Beat: 8},
			expectedPos: grid.GridKey{Line: 2, Beat: 0},
			description: "Cursor should move to start of line",
		},
		{
			name:        "Cursor Line End",
			command:     mappings.CursorLineEnd,
			initialPos:  grid.GridKey{Line: 1, Beat: 2},
			expectedPos: grid.GridKey{Line: 1, Beat: 31},
			description: "Cursor should move to end of line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithCurosrPos(tt.initialPos),
			)
			m, _ = processCommand(tt.command, m)
			assert.Equal(t, tt.expectedPos.Line, m.cursorPos.Line, tt.description+" - line position")
			assert.Equal(t, tt.expectedPos.Beat, m.cursorPos.Beat, tt.description+" - beat position")
		})
	}
}

func TestTempoChanges(t *testing.T) {
	tests := []struct {
		name          string
		commands      []mappings.Command
		initialTempo  int
		expectedTempo int
		description   string
	}{
		{
			name:          "Increase Tempo",
			commands:      []mappings.Command{mappings.TempoInputSwitch, mappings.Increase},
			initialTempo:  120,
			expectedTempo: 121,
			description:   "Tempo should increase by 1",
		},
		{
			name:          "Increase Tempo At Boundary",
			commands:      []mappings.Command{mappings.TempoInputSwitch, mappings.Increase},
			initialTempo:  300,
			expectedTempo: 300,
			description:   "Tempo should increase by 1",
		},
		{
			name:          "Decrease Tempo",
			commands:      []mappings.Command{mappings.TempoInputSwitch, mappings.Decrease},
			initialTempo:  130,
			expectedTempo: 129,
			description:   "Tempo should decrease by 1",
		},
		{
			name:          "Decrease Tempo At Boundary",
			commands:      []mappings.Command{mappings.TempoInputSwitch, mappings.Decrease},
			initialTempo:  30,
			expectedTempo: 30,
			description:   "Tempo should decrease by 1",
		},
		{
			name:          "Increase Tempo by 5",
			commands:      []mappings.Command{mappings.Increase},
			initialTempo:  120,
			expectedTempo: 125,
			description:   "Tempo should increase by 5",
		},
		{
			name:          "Increase Tempo by 5 at Boundary",
			commands:      []mappings.Command{mappings.Increase},
			initialTempo:  297,
			expectedTempo: 300,
			description:   "Tempo should increase by 5",
		},
		{
			name:          "Decrease Tempo by 5",
			commands:      []mappings.Command{mappings.Decrease},
			initialTempo:  130,
			expectedTempo: 125,
			description:   "Tempo should decrease by 5",
		},
		{
			name:          "Decrease Tempo by 5 at Boundary",
			commands:      []mappings.Command{mappings.Decrease},
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

func TestSubdivisionChanges(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []mappings.Command
		initialSubdivisions  int
		expectedSubdivisions int
		description          string
	}{
		{
			name:                 "Increase Subdivisions",
			commands:             []mappings.Command{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Increase},
			initialSubdivisions:  2,
			expectedSubdivisions: 3,
			description:          "Subdivisions should increase by 1",
		},
		{
			name:                 "Increase Subdivisions At Boundary",
			commands:             []mappings.Command{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Increase},
			initialSubdivisions:  8,
			expectedSubdivisions: 8,
			description:          "Subdivisions should be at maximum",
		},
		{
			name:                 "Decrease Subdivisions",
			commands:             []mappings.Command{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Decrease},
			initialSubdivisions:  3,
			expectedSubdivisions: 2,
			description:          "Subdivisions should decrease by 1",
		},
		{
			name:                 "Decrease Subdivisions At Boundary",
			commands:             []mappings.Command{mappings.TempoInputSwitch, mappings.TempoInputSwitch, mappings.Decrease},
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
		commands             []mappings.Command
		initialSetupChannel  uint8
		expectedSetupChannel uint8
		description          string
	}{
		{
			name:                 "Setup Channel Increase",
			commands:             []mappings.Command{mappings.SetupInputSwitch, mappings.Increase},
			initialSetupChannel:  0,
			expectedSetupChannel: 1,
			description:          "Setup input switch should select channel and increase should increment it",
		},
		{
			name:                 "Setup Channel Decrease",
			commands:             []mappings.Command{mappings.SetupInputSwitch, mappings.Decrease},
			initialSetupChannel:  5,
			expectedSetupChannel: 4,
			description:          "Setup input switch should select channel and decrease should decrement it",
		},
		{
			name:                 "Setup Channel Increase At Upper Boundary",
			commands:             []mappings.Command{mappings.SetupInputSwitch, mappings.Increase},
			initialSetupChannel:  16,
			expectedSetupChannel: 16,
			description:          "Setup input switch should select channel and increase should not go above 16",
		},
		{
			name:                 "Setup Channel Decrease At Lower Boundary",
			commands:             []mappings.Command{mappings.SetupInputSwitch, mappings.Decrease},
			initialSetupChannel:  1,
			expectedSetupChannel: 1,
			description:          "Setup input switch should select channel and decrease should not go below 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.cursorPos.Line].Channel = tt.initialSetupChannel
					return *m
				},
			)

			assert.Equal(t, tt.initialSetupChannel, m.definition.lines[m.cursorPos.Line].Channel, "Initial setup channel should match")
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectSetupChannel, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedSetupChannel, m.definition.lines[m.cursorPos.Line].Channel, tt.description+" - channel value")
		})
	}
}

func TestSetupInputSwitchMessageTypeIncrease(t *testing.T) {
	tests := []struct {
		name                string
		commands            []mappings.Command
		initialMessageType  grid.MessageType
		expectedMessageType grid.MessageType
		description         string
	}{
		{
			name:                "Message Type Increase from Note to CC",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialMessageType:  grid.MessageTypeNote,
			expectedMessageType: grid.MessageTypeCc,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Increase from CC to Program Change",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialMessageType:  grid.MessageTypeCc,
			expectedMessageType: grid.MessageTypeProgramChange,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Increase from Program Change to Note (wraparound)",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialMessageType:  grid.MessageTypeProgramChange,
			expectedMessageType: grid.MessageTypeNote,
			description:         "Two setup input switches should select message type and increase should wrap to Note",
		},
		{
			name:                "Message Type Decrease from Note to Program Change (wraparound)",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialMessageType:  grid.MessageTypeNote,
			expectedMessageType: grid.MessageTypeProgramChange,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Decrease from CC to Note",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialMessageType:  grid.MessageTypeCc,
			expectedMessageType: grid.MessageTypeNote,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Decrease from Program Change to CC",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialMessageType:  grid.MessageTypeProgramChange,
			expectedMessageType: grid.MessageTypeCc,
			description:         "Two setup input switches should select message type and increase should wrap to Note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.cursorPos.Line].MsgType = tt.initialMessageType
					return *m
				},
			)

			assert.Equal(t, tt.initialMessageType, m.definition.lines[m.cursorPos.Line].MsgType, "Initial message type should match")
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectSetupMessageType, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedMessageType, m.definition.lines[m.cursorPos.Line].MsgType, tt.description+" - message type value")
		})
	}
}

func TestSetupInputSwitchMessageTypeBackToGrid(t *testing.T) {
	tests := []struct {
		name                string
		commands            []mappings.Command
		expectedMessageType grid.MessageType
		description         string
	}{
		{
			name:                "Message Type Increase from Note to Program Change and back to Grid",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.Increase, mappings.SetupInputSwitch},
			expectedMessageType: grid.MessageTypeProgramChange,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
		{
			name:                "Message Type Increase from Note to Cc and back to Grid",
			commands:            []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.SetupInputSwitch},
			expectedMessageType: grid.MessageTypeCc,
			description:         "Two setup input switches should select message type and increase should increment it",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectNothing, m.selectionIndicator, "Should select back to nothing")
			assert.Equal(t, tt.expectedMessageType, m.definition.lines[m.cursorPos.Line].MsgType, tt.description+" - message type value")
		})
	}
}

func TestSetupNoteChange(t *testing.T) {
	tests := []struct {
		name         string
		commands     []mappings.Command
		initialNote  uint8
		expectedNote uint8
		description  string
	}{
		{
			name:         "Note Increase",
			commands:     []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialNote:  60,
			expectedNote: 61,
			description:  "Note should increase by 1",
		},
		{
			name:         "Note Decrease",
			commands:     []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialNote:  60,
			expectedNote: 59,
			description:  "Note should decrease by 1",
		},
		{
			name:         "Note Increase At Boundary",
			commands:     []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase},
			initialNote:  127,
			expectedNote: 127,
			description:  "Note should not increase beyond 127",
		},
		{
			name:         "Note Decrease at Boundary",
			commands:     []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialNote:  0,
			expectedNote: 0,
			description:  "Note should not decrease below 0",
		},
		{
			name:         "Note Decrease right above Boundary",
			commands:     []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Decrease},
			initialNote:  1,
			expectedNote: 0,
			description:  "Note should not decrease below 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.cursorPos.Line].Note = tt.initialNote
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedNote, m.definition.lines[m.cursorPos.Line].Note, tt.description)
		})
	}
}

func TestAccentInputSwitchDiffAndData(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []mappings.Command
		expectedAccentDiff   uint8
		expectedAccentValues []uint8
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase",
			commands:             []mappings.Command{mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentDiff:   16,
			expectedAccentValues: []uint8{0, 120, 104, 88, 72, 56, 40, 24, 8},
			description:          "Should select accent input and increase should set accent values",
		},
		{
			name:                 "Accent Input Switch with Decrease",
			commands:             []mappings.Command{mappings.AccentInputSwitch, mappings.Decrease},
			expectedAccentDiff:   14,
			expectedAccentValues: []uint8{0, 120, 106, 92, 78, 64, 50, 36, 22},
			description:          "Should select accent input and decrease should set accent values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectAccentDiff, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentDiff, m.definition.accents.Diff, tt.description+" - accent diff value")
			for i, value := range tt.expectedAccentValues {
				assert.Equalf(t, value, m.definition.accents.Data[i].Value, "accent values should match expected values - %d == %d", value, m.definition.accents.Data[0].Value)
			}
		})
	}
}

func TestAccentInputSwitchTarget(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []mappings.Command
		expectedAccentTarget accentTarget
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase on Target",
			commands:             []mappings.Command{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentTarget: AccentTargetNote,
			description:          "Should select accent input and increase should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target",
			commands:             []mappings.Command{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			expectedAccentTarget: AccentTargetNote,
			description:          "Should select accent input and decrease should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target Wraparound",
			commands:             []mappings.Command{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease, mappings.Decrease},
			expectedAccentTarget: AccentTargetVelocity,
			description:          "Should select accent input and decrease should set accent target to Velocity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectAccentTarget, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentTarget, m.definition.accents.Target, tt.description)
		})
	}
}

func TestAccentInputSwitchStartValue(t *testing.T) {
	tests := []struct {
		name                string
		commands            []mappings.Command
		initialAccentStart  uint8
		expectedAccentStart uint8
		description         string
	}{
		{
			name:                "Accent Input Switch with Increase on Start Value",
			commands:            []mappings.Command{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  120,
			expectedAccentStart: 121,
			description:         "Should select accent input and increase should set accent start value to 0",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value",
			commands:            []mappings.Command{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			initialAccentStart:  120,
			expectedAccentStart: 119,
			description:         "Should select accent input and decrease should set accent start value to 127",
		},
		{
			name:                "Accent Input Switch with Increase on Start Value at Boundary",
			commands:            []mappings.Command{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  127,
			expectedAccentStart: 127,
			description:         "Should select accent input and increase should not go above 127",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value at Boundary",
			commands:            []mappings.Command{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			initialAccentStart:  105,
			expectedAccentStart: 105,
			description:         "Should select accent input and decrease should not go below 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.accents.Start = tt.initialAccentStart
					m.definition.accents.ReCalc()
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectAccentStart, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentStart, m.definition.accents.Start, tt.description+" - accent start value")
		})
	}
}

func TestSetupCCChange(t *testing.T) {
	tests := []struct {
		name        string
		commands    []mappings.Command
		initialCc   uint8
		expectedCc  uint8
		description string
	}{
		{
			name:        "Should increase CC",
			commands:    []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Increase},
			initialCc:   0,
			expectedCc:  1,
			description: "Select CC and increment it",
		},
		{
			name:        "Should decrease CC",
			commands:    []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Decrease},
			initialCc:   1,
			expectedCc:  0,
			description: "Select CC and decrement it",
		},
		{
			name:        "Should skip unused CC on increase",
			commands:    []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Increase},
			initialCc:   2,
			expectedCc:  4,
			description: "Should skip over unused CCs on increase",
		},
		{
			name:        "Should skip unused CC on decrease",
			commands:    []mappings.Command{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.SetupInputSwitch, mappings.Decrease},
			initialCc:   4,
			expectedCc:  2,
			description: "Should skip over unused CCs on decrease",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.lines[m.cursorPos.Line].Note = tt.initialCc
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedCc, m.definition.lines[m.cursorPos.Line].Note, tt.description)
		})
	}
}

func TestOverlayInputSwitch(t *testing.T) {
	tests := []struct {
		name              string
		commands          []mappings.Command
		expectedSelection Selection
		expectedFocus     focus
		description       string
	}{
		{
			name:              "Single OverlayInputSwitch",
			commands:          []mappings.Command{mappings.OverlayInputSwitch},
			expectedSelection: SelectOverlay,
			expectedFocus:     FocusOverlayKey,
			description:       "First overlay input switch should select overlay and set focus to overlay key",
		},
		{
			name:              "Double OverlayInputSwitch",
			commands:          []mappings.Command{mappings.OverlayInputSwitch, mappings.OverlayInputSwitch},
			expectedSelection: SelectNothing,
			expectedFocus:     FocusGrid,
			description:       "Second overlay input switch should cycle back to nothing but keep overlay key focus",
		},
		{
			name:              "Escape Overlay Input state",
			commands:          []mappings.Command{mappings.OverlayInputSwitch, mappings.Escape},
			expectedSelection: SelectNothing,
			expectedFocus:     FocusGrid,
			description:       "Escape should set focus and selection back to an initial state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd tea.Cmd
			m := createTestModel()

			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")
			assert.Equal(t, FocusGrid, m.focus, "Initial focus should be grid")

			m, cmd = processCommands(tt.commands, m)
			if cmd != nil {
				updateModel, _ := m.Update(cmd())
				switch um := updateModel.(type) {
				case model:
					m = um
				}
			}

			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedFocus, m.focus, tt.description+" - focus state")
		})
	}
}

func createTestModel(modelFns ...modelFunc) model {

	m := InitModel("", seqmidi.MidiConnection{}, "", "", MlmStandAlone, "default")

	for _, fn := range modelFns {
		m = fn(&m)
	}

	return m
}

func processCommands(commands []mappings.Command, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	for _, command := range commands {
		m, cmd = processCommand(command, m)
	}
	return m, cmd
}

func processCommand(command mappings.Command, m model) (model, tea.Cmd) {
	keyMsgs := getKeyMsgs(command)
	var cmd tea.Cmd
	for _, keyMsg := range keyMsgs {
		var updateModel tea.Model
		updateModel, cmd = m.Update(keyMsg)
		switch um := updateModel.(type) {
		case model:
			m = um
		}
	}
	return m, cmd
}

func getKeyMsgs(command mappings.Command) []tea.KeyMsg {
	keys := mappings.KeysForCommand(command)
	var keyMsgs []tea.KeyMsg
	for _, key := range keys {
		keyMsgs = append(keyMsgs, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	}
	return keyMsgs
}

type modelFunc func(m *model) model

func WithCurosrPos(pos grid.GridKey) modelFunc {
	return func(m *model) model {
		m.cursorPos = pos
		return *m
	}
}
