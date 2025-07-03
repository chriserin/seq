package main

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/overlaykey"
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
		commands             []any
		expectedAccentDiff   uint8
		expectedAccentValues []uint8
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase",
			commands:             []any{mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentDiff:   16,
			expectedAccentValues: []uint8{0, 120, 104, 88, 72, 56, 40, 24, 8},
			description:          "Should select accent input and increase should set accent values",
		},
		{
			name:                 "Accent Input Switch with Decrease",
			commands:             []any{mappings.AccentInputSwitch, mappings.Decrease},
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
		commands             []any
		expectedAccentTarget accentTarget
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase on Target",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentTarget: AccentTargetNote,
			description:          "Should select accent input and increase should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			expectedAccentTarget: AccentTargetNote,
			description:          "Should select accent input and decrease should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target Wraparound",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease, mappings.Decrease},
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
		commands            []any
		initialAccentStart  uint8
		expectedAccentStart uint8
		description         string
	}{
		{
			name:                "Accent Input Switch with Increase on Start Value",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  120,
			expectedAccentStart: 121,
			description:         "Should select accent input and increase should set accent start value to 0",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			initialAccentStart:  120,
			expectedAccentStart: 119,
			description:         "Should select accent input and decrease should set accent start value to 127",
		},
		{
			name:                "Accent Input Switch with Increase on Start Value at Boundary",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  127,
			expectedAccentStart: 127,
			description:         "Should select accent input and increase should not go above 127",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value at Boundary",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
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
		commands          []any
		expectedSelection Selection
		expectedFocus     focus
		description       string
	}{
		{
			name:              "Single OverlayInputSwitch",
			commands:          []any{mappings.OverlayInputSwitch},
			expectedSelection: SelectNothing,
			expectedFocus:     FocusOverlayKey,
			description:       "First overlay input switch should select overlay and set focus to overlay key",
		},
		{
			name:              "Double OverlayInputSwitch",
			commands:          []any{mappings.OverlayInputSwitch, mappings.OverlayInputSwitch},
			expectedSelection: SelectNothing,
			expectedFocus:     FocusGrid,
			description:       "Second overlay input switch should cycle back to nothing but keep overlay key focus",
		},
		{
			name:              "Escape Overlay Input state",
			commands:          []any{mappings.OverlayInputSwitch, mappings.Escape},
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

func TestRatchetInputSwitch(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedSelection Selection
		description       string
	}{
		{
			name:              "RatchetInputSwitch when cursor note no note",
			commands:          []any{mappings.RatchetInputSwitch},
			expectedSelection: SelectNothing,
			description:       "First ratchet input switch does nothing if not on a note",
		},
		{
			name:              "RatchetInputSwitch when cursor on note",
			commands:          []any{mappings.TriggerAdd, mappings.RatchetInputSwitch},
			expectedSelection: SelectRatchets,
			description:       "First ratchet input switch should select ratchet",
		},
		{
			name:              "Second RatchetInputSwitch when cursor on note",
			commands:          []any{mappings.TriggerAdd, mappings.RatchetInputSwitch, mappings.RatchetInputSwitch},
			expectedSelection: SelectRatchetSpan,
			description:       "Second ratchet input switch should select ratchet span",
		},
		{
			name:              "Third RatchetInputSwitch",
			commands:          []any{mappings.TriggerAdd, mappings.RatchetInputSwitch, mappings.RatchetInputSwitch, mappings.RatchetInputSwitch},
			expectedSelection: SelectNothing,
			description:       "Second ratchet input switch should select ratchet span",
		},
		{
			name:              "Escape Ratchet Input",
			commands:          []any{mappings.TriggerAdd, mappings.RatchetInputSwitch, mappings.Escape},
			expectedSelection: SelectNothing,
			description:       "Second ratchet input switch should select ratchet span",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description)
		})
	}
}

func TestRatchetInputValues(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedRatchet grid.Ratchet
		description     string
	}{
		{
			name: "RatchetInputSwitch with Mute Ratchet",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetInputSwitch,
				mappings.Mute,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   0,
				Length: 0,
			},
			description: "Should select ratchet input and mute should set all values to 0",
		},
		{
			name: "RatchetInputSwitch with Mute Ratchet and Remute",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetInputSwitch,
				mappings.Mute,
				mappings.Mute,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1,
				Length: 0,
			},
			description: "",
		},
		{
			name: "RatchetInputSwitch with Span Increase",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetInputSwitch,
				mappings.RatchetInputSwitch,
				mappings.Increase,
			},
			expectedRatchet: grid.Ratchet{
				Span:   1,
				Hits:   1,
				Length: 0,
			},
			description: "",
		},
		{
			name: "RatchetInputSwitch with Span Increase/Decrease",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetInputSwitch,
				mappings.RatchetInputSwitch,
				mappings.Increase,
				mappings.Decrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1,
				Length: 0,
			},
			description: "",
		},
		{
			name: "RatchetInputSwitch with Mute Second Hit",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
				mappings.RatchetInputSwitch,
				mappings.CursorRight,
				mappings.Mute,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1,
				Length: 1,
			},
			description: "",
		},
		{
			name: "RatchetInputSwitch with Mute First Hit after moving cursor",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
				mappings.RatchetInputSwitch,
				mappings.CursorRight,
				mappings.CursorLeft,
				mappings.Mute,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   2, // This is a bit mask so 0b10 = 2
				Length: 1,
			},
			description: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)
			currentNote, exists := m.CurrentNote()
			currentRatchet := currentNote.Ratchets
			assert.True(t, exists, tt.description+" - current note should exist")
			assert.Equal(t, tt.expectedRatchet.Span, currentRatchet.Span, tt.description+" - ratchet span")
			assert.Equal(t, tt.expectedRatchet.Hits, currentRatchet.Hits, tt.description+" - ratchet hits")
			assert.Equal(t, tt.expectedRatchet.Length, currentRatchet.Length, tt.description+" - ratchet length")
		})
	}
}

func TestBeatInputSwitchIncrease(t *testing.T) {
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
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectBeats, m.selectionIndicator, tt.description+" - selection state")
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
			commands:       []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Increase},
			initialCycles:  4,
			expectedCycles: 5,
			description:    "Three beat input switches should select cycles and increase should increment it",
		},
		{
			name:           "Beat Input Switch Cycles with Decrease",
			commands:       []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Decrease},
			initialCycles:  4,
			expectedCycles: 3,
			description:    "Three beat input switches should select cycles and decrease should decrement it",
		},
		{
			name:           "Beat Input Switch Cycles with Increase At Upper Boundary",
			commands:       []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Increase},
			initialCycles:  127,
			expectedCycles: 127,
			description:    "Three beat input switches should select cycles and increase should not go above 127",
		},
		{
			name:           "Beat Input Switch Cycles with Decrease At Lower Boundary",
			commands:       []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Decrease},
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
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectCycles, m.selectionIndicator, tt.description+" - selection state")
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
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectStartBeats, m.selectionIndicator, tt.description+" - selection state")
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
			commands:            []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Increase},
			initialStartCycles:  4,
			expectedStartCycles: 5,
			description:         "Two beat input switches should select start cycles and increase should increment it",
		},
		{
			name:                "Beat Input Switch StartCycles with Decrease",
			commands:            []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Decrease},
			initialStartCycles:  4,
			expectedStartCycles: 3,
			description:         "Two beat input switches should select start cycles and decrease should decrement it",
		},
		{
			name:                "Beat Input Switch StartCycles with Increase At Upper Boundary",
			commands:            []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Increase},
			initialStartCycles:  127,
			expectedStartCycles: 127,
			description:         "Two beat input switches should select start cycles and increase should not go above 127",
		},
		{
			name:                "Beat Input Switch StartCycles with Decrease At Lower Boundary",
			commands:            []any{mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.BeatInputSwitch, mappings.Decrease},
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
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectStartCycles, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedStartCycles, m.CurrentSongSection().StartCycles, tt.description+" - start cycles value")
		})
	}
}

func TestToggleArrangementViewSwitch(t *testing.T) {

	tests := []struct {
		name              string
		commands          []any
		expectedFocus     focus
		expectedSelection Selection
		expectedArrIsOpen bool
		description       string
	}{
		{
			name:              "Toggle Arrangement View",
			commands:          []any{mappings.ToggleArrangementView},
			expectedFocus:     FocusArrangementEditor,
			expectedSelection: SelectNothing,
			expectedArrIsOpen: true,
			description:       "First toggle should open arrangement view",
		},
		{
			name:              "Toggle Arrangement View Switch Back to Grid",
			commands:          []any{mappings.ToggleArrangementView, mappings.ToggleArrangementView},
			expectedFocus:     FocusGrid,
			expectedSelection: SelectNothing,
			expectedArrIsOpen: false,
			description:       "Second toggle should switch back to grid and close arrangement view",
		},
		{
			name:              "Toggle Arrangement View Switch to Grid keep Arrangement Open",
			commands:          []any{mappings.ToggleArrangementView, mappings.Escape},
			expectedFocus:     FocusGrid,
			expectedSelection: SelectNothing,
			expectedArrIsOpen: true,
			description:       "Escape should switch back to grid and keep arrangement view open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd tea.Cmd
			m := createTestModel()

			m, cmd = processCommands(tt.commands, m)
			if cmd != nil {
				updateModel, _ := m.Update(cmd())
				switch um := updateModel.(type) {
				case model:
					m = um
				}
			}
			assert.Equal(t, tt.expectedFocus, m.focus, tt.description+" - arrangement view state")
			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedArrIsOpen, m.showArrangementView, tt.description+" - arrangement open state")
		})
	}
}

func TestPatternModeGateIncrease(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		expectedGate uint8
		description  string
	}{
		{
			name: "Add note, switch to gate mode, increase gate by 1",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedGate: 1,
			description:  "Should add note, switch to gate mode, and increase gate by 1",
		},
		{
			name: "Add note, switch to gate mode, move to right no change",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleGateMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedGate: 0,
			description:  "Should add note, switch to gate mode, and increase gate by 1",
		},
		{
			name: "Add note, switch to gate mode, increase gate by 1, decrease gate by 1",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedGate: 0,
			description:  "Should add note, switch to gate mode, and decrease gate by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedGate, currentNote.GateIndex, tt.description+" - gate value")
		})
	}
}

func TestPatternModeAccentIncrease(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		expectedAccent uint8
		description    string
	}{
		{
			name: "Add note, switch to accent mode, increase accent by 1",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleAccentMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedAccent: 4,
			description:    "Should add note, switch to accent mode, and increase accent by 1",
		},
		{
			name: "Add note, switch to accent mode, move to right no change",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleAccentMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedAccent: 5,
			description:    "Should add note, switch to accent mode, move cursor right and accent should remain 0",
		},
		{
			name: "Add note, switch to accent mode, increase accent by 1, decrease accent by 1",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleAccentMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedAccent: 6,
			description:    "Should add note, switch to accent mode, and decrease accent by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// NOTE: Initial value for accent is 5
			// and "increasing" accent index means moving to 4
			// as accents are defined backwards
			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedAccent, currentNote.AccentIndex, tt.description+" - accent value")
		})
	}
}

func TestPatternModeRatchetIncrease(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedRatchet grid.Ratchet
		description     string
	}{
		{
			name: "Add note, switch to ratchet mode, increase ratchet by 1",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleRatchetMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   3, // hits are 0b11, so 3
				Length: 1,
			},
			description: "Should add note, switch to ratchet mode, and increase ratchet by 1",
		},
		{
			name: "Add note, switch to ratchet mode, move to right no change",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleRatchetMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1,
				Length: 0,
			},
			description: "Should add note, switch to ratchet mode, move cursor right and ratchet should remain default",
		},
		{
			name: "Add note, switch to ratchet mode, increase then decrease ratchet",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleRatchetMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1,
				Length: 0,
			},
			description: "Should add note, switch to ratchet mode, and set ratchet to 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedRatchet.Span, currentNote.Ratchets.Span, tt.description+" - ratchet span")
			assert.Equal(t, tt.expectedRatchet.Hits, currentNote.Ratchets.Hits, tt.description+" - ratchet hits")
			assert.Equal(t, tt.expectedRatchet.Length, currentNote.Ratchets.Length, tt.description+" - ratchet length")
		})
	}
}

func TestPatternModeWaitIncrease(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		expectedWait uint8
		description  string
	}{
		{
			name: "Add note, switch to wait mode, increase wait by 1",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleWaitMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedWait: 1,
			description:  "Should add note, switch to wait mode, and increase wait by 1",
		},
		{
			name: "Add note, switch to wait mode, move to right no change",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleWaitMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedWait: 0,
			description:  "Should add note, switch to wait mode, move cursor right and wait should remain 0",
		},
		{
			name: "Add note, switch to wait mode, increase wait by 1, decrease wait by 1",
			commands: []any{
				mappings.TriggerAdd,
				mappings.ToggleWaitMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedWait: 0,
			description:  "Should add note, switch to wait mode, and decrease wait by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedWait, currentNote.WaitIndex, tt.description+" - wait value")
		})
	}
}

func TestAccentIncrease(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		expectedAccent uint8
		description    string
	}{
		{
			name: "Add note and increase accent",
			commands: []any{
				mappings.TriggerAdd,
				mappings.AccentIncrease,
			},
			expectedAccent: 4,
			description:    "Should add note and increase accent by 1 (index 5 -> 4)",
		},
		{
			name: "Add note and increase accent twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.AccentIncrease,
				mappings.AccentIncrease,
			},
			expectedAccent: 3,
			description:    "Should add note and increase accent by 2 (index 5 -> 3)",
		},
		{
			name: "Add note and increase accent at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.AccentIncrease,
				mappings.AccentIncrease,
				mappings.AccentIncrease,
				mappings.AccentIncrease,
				mappings.AccentIncrease,
			},
			expectedAccent: 1,
			description:    "Should add note and increase accent to minimum index (1)",
		},
		{
			name: "Add note and decrease accent",
			commands: []any{
				mappings.TriggerAdd,
				mappings.AccentDecrease,
			},
			expectedAccent: 6,
			description:    "Should add note and decrease accent by 1 (index 5 -> 6)",
		},
		{
			name: "Add note and decrease accent twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.AccentDecrease,
				mappings.AccentDecrease,
			},
			expectedAccent: 7,
			description:    "Should add note and decrease accent by 2 (index 5 -> 7)",
		},
		{
			name: "Add note and decrease accent at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.AccentDecrease,
				mappings.AccentDecrease,
				mappings.AccentDecrease,
				mappings.AccentDecrease,
				mappings.AccentDecrease,
			},
			expectedAccent: 8,
			description:    "Should add note and decrease accent to maximum index (8)",
		},
		{
			name: "Add note, increase then decrease accent",
			commands: []any{
				mappings.TriggerAdd,
				mappings.AccentIncrease,
				mappings.AccentDecrease,
			},
			expectedAccent: 5,
			description:    "Should add note, increase then decrease accent (index 5 -> 4 -> 5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedAccent, currentNote.AccentIndex, tt.description+" - accent value")
		})
	}
}

func TestGateIncrease(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		expectedGate uint8
		description  string
	}{
		{
			name: "Add note and increase gate",
			commands: []any{
				mappings.TriggerAdd,
				mappings.GateIncrease,
			},
			expectedGate: 1,
			description:  "Should add note and increase gate by 1 (index 0 -> 1)",
		},
		{
			name: "Add note and increase gate twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.GateIncrease,
				mappings.GateIncrease,
			},
			expectedGate: 2,
			description:  "Should add note and increase gate by 2 (index 0 -> 2)",
		},
		{
			name: "Add note and increase gate at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateIncrease,
			},
			expectedGate: 7,
			description:  "Should add note and increase gate to maximum index (8)",
		},
		{
			name: "Add note and decrease gate",
			commands: []any{
				mappings.TriggerAdd,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateDecrease,
			},
			expectedGate: 1,
			description:  "Should add note and decrease gate by 1 (index 2 -> 1)",
		},
		{
			name: "Add note and decrease gate twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateIncrease,
				mappings.GateDecrease,
				mappings.GateDecrease,
			},
			expectedGate: 1,
			description:  "Should add note and decrease gate by 2 (index 3 -> 1)",
		},
		{
			name: "Add note and decrease gate at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.GateDecrease,
			},
			expectedGate: 0,
			description:  "Should add note and decrease gate stays at minimum index (0)",
		},
		{
			name: "Add note, increase then decrease gate",
			commands: []any{
				mappings.TriggerAdd,
				mappings.GateIncrease,
				mappings.GateDecrease,
			},
			expectedGate: 0,
			description:  "Should add note, increase then decrease gate (index 0 -> 1 -> 0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedGate, currentNote.GateIndex, tt.description+" - gate value")
		})
	}
}

func TestWaitIncrease(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		expectedWait uint8
		description  string
	}{
		{
			name: "Add note and increase wait",
			commands: []any{
				mappings.TriggerAdd,
				mappings.WaitIncrease,
			},
			expectedWait: 1,
			description:  "Should add note and increase wait by 1 (index 0 -> 1)",
		},
		{
			name: "Add note and increase wait twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
			},
			expectedWait: 2,
			description:  "Should add note and increase wait by 2 (index 0 -> 2)",
		},
		{
			name: "Add note and increase wait at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
			},
			expectedWait: 7,
			description:  "Should add note and increase wait to maximum index (8)",
		},
		{
			name: "Add note and decrease wait",
			commands: []any{
				mappings.TriggerAdd,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitDecrease,
			},
			expectedWait: 1,
			description:  "Should add note and decrease wait by 1 (index 2 -> 1)",
		},
		{
			name: "Add note and decrease wait twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
				mappings.WaitDecrease,
				mappings.WaitDecrease,
			},
			expectedWait: 1,
			description:  "Should add note and decrease wait by 2 (index 3 -> 1)",
		},
		{
			name: "Add note and decrease wait at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.WaitDecrease,
			},
			expectedWait: 0,
			description:  "Should add note and decrease wait stays at minimum index (0)",
		},
		{
			name: "Add note, increase then decrease wait",
			commands: []any{
				mappings.TriggerAdd,
				mappings.WaitIncrease,
				mappings.WaitDecrease,
			},
			expectedWait: 0,
			description:  "Should add note, increase then decrease wait (index 0 -> 1 -> 0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedWait, currentNote.WaitIndex, tt.description+" - wait value")
		})
	}
}

func TestRatchetIncrease(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedRatchet grid.Ratchet
		description     string
	}{
		{
			name: "Add note and increase ratchet",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   3, // 0b11 = 3 (two hits)
				Length: 1,
			},
			description: "Should add note and increase ratchet by 1 (1 -> 2 hits)",
		},
		{
			name: "Add note and increase ratchet twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   7, // 0b111 = 7 (three hits)
				Length: 2,
			},
			description: "Should add note and increase ratchet by 2 (1 -> 3 hits)",
		},
		{
			name: "Add note and increase ratchet at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   255, // 0b11111111 = 255 (eight hits)
				Length: 7,
			},
			description: "Should add note and increase ratchet to maximum (8 hits)",
		},
		{
			name: "Add note and decrease ratchet",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetDecrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   3, // 0b11 = 3 (two hits)
				Length: 1,
			},
			description: "Should add note and decrease ratchet by 1 (3 -> 2 hits)",
		},
		{
			name: "Add note and decrease ratchet twice",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetIncrease,
				mappings.RatchetDecrease,
				mappings.RatchetDecrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   3, // 0b11 = 3 (two hits)
				Length: 1,
			},
			description: "Should add note and decrease ratchet by 2 (4 -> 2 hits)",
		},
		{
			name: "Add note and decrease ratchet at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetDecrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1, // 0b1 = 1 (one hit - minimum)
				Length: 0,
			},
			description: "Should add note and decrease ratchet stays at minimum (1 hit)",
		},
		{
			name: "Add note, increase then decrease ratchet",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RatchetIncrease,
				mappings.RatchetDecrease,
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1, // 0b1 = 1 (one hit - back to default)
				Length: 0,
			},
			description: "Should add note, increase then decrease ratchet (1 -> 2 -> 1 hits)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedRatchet.Span, currentNote.Ratchets.Span, tt.description+" - ratchet span")
			assert.Equal(t, tt.expectedRatchet.Hits, currentNote.Ratchets.Hits, tt.description+" - ratchet hits")
			assert.Equal(t, tt.expectedRatchet.Length, currentNote.Ratchets.Length, tt.description+" - ratchet length")
		})
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		initialPos  grid.GridKey
		expectedPos grid.GridKey
		description string
	}{
		{
			name: "Rotate down",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RotateDown,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 4},
			expectedPos: grid.GridKey{Line: 3, Beat: 4},
			description: "Should rotate pattern down by one line",
		},
		{
			name: "Rotate down at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RotateDown,
			},
			initialPos:  grid.GridKey{Line: 7, Beat: 4},
			expectedPos: grid.GridKey{Line: 0, Beat: 4},
			description: "Should rotate pattern down wrapping to top",
		},
		{
			name: "Rotate up",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RotateUp,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 4},
			expectedPos: grid.GridKey{Line: 1, Beat: 4},
			description: "Should rotate pattern up by one line",
		},
		{
			name: "Rotate up at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RotateUp,
			},
			initialPos:  grid.GridKey{Line: 0, Beat: 4},
			expectedPos: grid.GridKey{Line: 7, Beat: 4},
			description: "Should rotate pattern up wrapping to bottom",
		},
		{
			name: "Rotate right",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RotateRight,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 4},
			expectedPos: grid.GridKey{Line: 2, Beat: 5},
			description: "Should rotate pattern right by one beat",
		},
		{
			name: "Rotate right at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.CursorLineStart,
				mappings.RotateRight,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 31},
			expectedPos: grid.GridKey{Line: 2, Beat: 0},
			description: "Should rotate pattern right wrapping to start",
		},
		{
			name: "Rotate left",
			commands: []any{
				mappings.TriggerAdd,
				mappings.CursorLeft,
				mappings.RotateLeft,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 4},
			expectedPos: grid.GridKey{Line: 2, Beat: 3},
			description: "Should rotate pattern left by one beat",
		},
		{
			name: "Rotate left at boundary",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RotateLeft,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 0},
			expectedPos: grid.GridKey{Line: 2, Beat: 31},
			description: "Should rotate pattern left wrapping to end",
		},
		{
			name: "Multiple rotations",
			commands: []any{
				mappings.TriggerAdd,
				mappings.RotateDown,
				mappings.CursorDown,
				mappings.RotateRight,
			},
			initialPos:  grid.GridKey{Line: 1, Beat: 2},
			expectedPos: grid.GridKey{Line: 2, Beat: 3},
			description: "Should rotate pattern down and right",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithCurosrPos(tt.initialPos),
			)

			m, _ = processCommand(mappings.TriggerAdd, m)
			initialNote, exists := m.currentOverlay.GetNote(tt.initialPos)
			assert.True(t, exists, tt.description+" - note should exist at initial position")

			rotateCommands := tt.commands[1:]
			m, _ = processCommands(rotateCommands, m)

			rotatedNote, exists := m.currentOverlay.GetNote(tt.expectedPos)
			assert.True(t, exists, tt.description+" - note should exist at expected position")
			assert.Equal(t, initialNote, rotatedNote, tt.description+" - note should be the same")

			if tt.initialPos != tt.expectedPos {
				_, stillExists := m.currentOverlay.GetNote(tt.initialPos)
				assert.False(t, stillExists, tt.description+" - note should not exist at initial position")
			}
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

func processCommands(commands []any, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	for _, command := range commands {
		switch c := command.(type) {
		case mappings.Command:
			m, cmd = processCommand(c, m)
		case mappings.Mapping:
			m, cmd = processMapping(c, m)
		}
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

func processMapping(mapping mappings.Mapping, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	var updateModel tea.Model
	switch mapping.Command {
	case mappings.NumberPattern:
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(mapping.LastValue)}
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

func TestNextOverlay(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		addedOverlayKeys   []overlaykey.OverlayPeriodicity
		expectedOverlayKey overlaykey.OverlayPeriodicity
		description        string
	}{
		{
			name: "Next Overlay",
			commands: []any{
				mappings.NextOverlay,
			},
			addedOverlayKeys: []overlaykey.OverlayPeriodicity{
				{Shift: 1, Interval: 2, Width: 0, StartCycle: 0},
				{Shift: 1, Interval: 3, Width: 0, StartCycle: 0},
			},
			expectedOverlayKey: overlaykey.OverlayPeriodicity{Shift: 1, Interval: 2, Width: 0, StartCycle: 0},
			description:        "Should switch to next overlay with key 1/2",
		},
		{
			name: "Next And Prev Overlay",
			commands: []any{
				mappings.NextOverlay,
				mappings.NextOverlay,
				mappings.PrevOverlay,
			},
			addedOverlayKeys: []overlaykey.OverlayPeriodicity{
				{Shift: 1, Interval: 2, Width: 0, StartCycle: 0},
				{Shift: 1, Interval: 3, Width: 0, StartCycle: 0},
			},
			expectedOverlayKey: overlaykey.OverlayPeriodicity{Shift: 1, Interval: 2, Width: 0, StartCycle: 0},
			description:        "Should switch to next overlay with key 1/2",
		},
		{
			name: "Back to Root",
			commands: []any{
				mappings.NextOverlay,
				mappings.NextOverlay,
				mappings.PrevOverlay,
				mappings.PrevOverlay,
			},
			addedOverlayKeys: []overlaykey.OverlayPeriodicity{
				{Shift: 1, Interval: 2, Width: 0, StartCycle: 0},
				{Shift: 1, Interval: 3, Width: 0, StartCycle: 0},
			},
			expectedOverlayKey: overlaykey.OverlayPeriodicity{Shift: 1, Interval: 1, Width: 0, StartCycle: 0},
			description:        "Should switch to next overlay with key 1/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			rootKey := overlaykey.OverlayPeriodicity{Shift: 1, Interval: 1, Width: 0, StartCycle: 0}

			for _, key := range tt.addedOverlayKeys {
				(*m.definition.parts)[0].Overlays = m.CurrentPart().Overlays.Add(key)
			}

			assert.Equal(t, rootKey, m.currentOverlay.Key, "Initial overlay key should be root")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedOverlayKey, m.currentOverlay.Key, tt.description)
		})
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name        string
		command     mappings.Command
		description string
	}{
		{
			name:        "Save With Filename",
			command:     mappings.Save,
			description: "Should save file when filename is set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary filename
			tempDir := t.TempDir()
			testFilename := filepath.Join(tempDir, "test_save.seq")

			m := createTestModel(
				func(m *model) model {
					m.filename = testFilename
					return *m
				},
			)

			_, err := os.Stat(testFilename)
			assert.True(t, os.IsNotExist(err), "File should not exist initially")

			processCommand(tt.command, m)

			_, err = os.Stat(testFilename)
			assert.NoError(t, err, tt.description+" - file should be created")

			fileInfo, err := os.Stat(testFilename)
			assert.NoError(t, err, "Should be able to get file info")
			assert.Greater(t, fileInfo.Size(), int64(0), "File should not be empty")
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "New Sequence Clears Notes",
			commands: []any{
				mappings.TriggerAdd,
				mappings.CursorRight,
				mappings.New,
				mappings.Enter,
			},
			description: "Should clear all notes when creating new sequence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			note, exists := m.CurrentNote()
			assert.False(t, exists, tt.description+" - note should not exist after new sequence")
			assert.Equal(t, zeronote, note, tt.description+" - note should be zero note")

			assert.Equal(t, grid.GridKey{Line: 0, Beat: 0}, m.cursorPos, "Cursor should be reset to origin")
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
			expectedLastNote:    61,
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
			expectedLastNote:    63,
			description:         "Should add three new lines with incrementing note values",
		},
		{
			name: "Cannot add more than 16 lines total",
			commands: []any{
				mappings.NewLine,
				mappings.NewLine,
				mappings.NewLine,
				mappings.NewLine,
				mappings.NewLine,
				mappings.NewLine,
				mappings.NewLine,
				mappings.NewLine,
				mappings.NewLine,
			},
			initialLineCount:    8,
			expectedLineCount:   16,
			expectedLastChannel: 10,
			expectedLastNote:    68,
			description:         "Should cap at 16 lines total",
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

func TestSectionCommands(t *testing.T) {
	tests := []struct {
		name                  string
		commands              []any
		expectedSelection     Selection
		expectedSideIndicator bool
		description           string
	}{
		{
			name:                  "NewSectionAfter sets selection and side indicator",
			commands:              []any{mappings.NewSectionAfter},
			expectedSelection:     SelectPart,
			expectedSideIndicator: true,
			description:           "Should set selection to SelectPart and sectionSideIndicator to true",
		},
		{
			name:                  "NewSectionBefore sets selection and side indicator",
			commands:              []any{mappings.NewSectionBefore},
			expectedSelection:     SelectPart,
			expectedSideIndicator: false,
			description:           "Should set selection to SelectPart and sectionSideIndicator to false",
		},
		{
			name:              "ChangePart sets selection",
			commands:          []any{mappings.ChangePart},
			expectedSelection: SelectChangePart,
			description:       "Should set selection to SelectChangePart",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Verify initial state
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			// Execute commands
			m, _ = processCommands(tt.commands, m)

			// Verify final state
			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description+" - selection state")

			if tt.name == "NewSectionAfter sets selection and side indicator" || tt.name == "NewSectionBefore sets selection and side indicator" {
				assert.Equal(t, tt.expectedSideIndicator, m.sectionSideIndicator, tt.description+" - side indicator")
			}
		})
	}
}

func TestSectionNavigation(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		expectedPartIndex  int
		expectedMoveResult bool
		description        string
	}{
		{
			name:               "NextSection moves cursor forward",
			commands:           []any{mappings.NewSectionAfter, mappings.Enter, mappings.NextSection},
			expectedMoveResult: true,
			description:        "Should move to next section successfully",
		},
		{
			name:               "PrevSection moves cursor backward",
			commands:           []any{mappings.NewSectionBefore, mappings.Enter, mappings.PrevSection},
			expectedMoveResult: true,
			description:        "Should move to previous section successfully",
		},
		{
			name:               "NextSection on single section",
			commands:           []any{mappings.NextSection},
			expectedMoveResult: false,
			description:        "Should not move when only one section exists",
		},
		{
			name:               "PrevSection on first section",
			commands:           []any{mappings.NewSectionAfter, mappings.PrevSection},
			expectedMoveResult: false,
			description:        "Should not move when already on first section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Store initial cursor state for comparison
			initialCursor := m.arrangement.Cursor

			// Execute commands
			m, _ = processCommands(tt.commands, m)

			// Verify cursor movement occurred or not as expected
			cursorMoved := !m.arrangement.Cursor.Equals(initialCursor)

			assert.Equal(t, tt.expectedMoveResult, cursorMoved, tt.description+" - cursor movement")
		})
	}
}

func TestSectionNavigationResetOverlay(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name:        "NextSection resets current overlay",
			commands:    []any{mappings.NewSectionAfter, mappings.Enter, mappings.NextSection},
			description: "Should reset current overlay after moving to next section",
		},
		{
			name:        "PrevSection resets current overlay",
			commands:    []any{mappings.NewSectionBefore, mappings.Enter, mappings.PrevSection},
			description: "Should reset current overlay after moving to previous section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					nonRootKey := overlaykey.OverlayPeriodicity{Shift: 1, Interval: 2, Width: 0, StartCycle: 0}
					(*m.definition.parts)[0].Overlays = m.CurrentPart().Overlays.Add(nonRootKey)
					m.currentOverlay = (*m.definition.parts)[0].Overlays.FindOverlay(nonRootKey)
					return *m
				},
			)

			rootKey := overlaykey.OverlayPeriodicity{Shift: 1, Interval: 1, Width: 0, StartCycle: 0}
			assert.NotEqual(t, rootKey, m.currentOverlay.Key, "Should start with non-root overlay")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, rootKey, m.currentOverlay.Key, tt.description+" - overlay should be reset to root")
		})
	}
}

func WithCurosrPos(pos grid.GridKey) modelFunc {
	return func(m *model) model {
		m.cursorPos = pos
		return *m
	}
}
