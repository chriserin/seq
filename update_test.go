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
			m = processCommand(tt.command, m)
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
			m = processCommands(tt.commands, m)
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
			m = processCommands(tt.commands, m)
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

			m = processCommands(tt.commands, m)

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

			m = processCommands(tt.commands, m)

			assert.Equal(t, SelectSetupMessageType, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedMessageType, m.definition.lines[m.cursorPos.Line].MsgType, tt.description+" - message type value")
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

func processCommands(commands []mappings.Command, m model) model {
	for _, command := range commands {
		m = processCommand(command, m)
	}
	return m
}

func processCommand(command mappings.Command, m model) model {
	keyMsgs := getKeyMsgs(command)
	for _, keyMsg := range keyMsgs {
		updateModel, _ := m.Update(keyMsg)
		switch um := updateModel.(type) {
		case model:
			m = um
		}
	}
	return m
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
