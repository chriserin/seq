package main

import (
	"testing"

	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUndoGridNote(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedCursorPos grid.GridKey
		description       string
	}{
		{
			name: "Undo note addition",
			commands: []any{
				mappings.NoteAdd,
				mappings.Undo,
			},
			expectedCursorPos: grid.GridKey{Line: 0, Beat: 0},
			description:       "Should remove note after adding and undoing",
		},
		{
			name: "Undo note removal",
			commands: []any{
				mappings.NoteAdd,
				mappings.NoteRemove,
				mappings.Undo,
			},
			expectedCursorPos: grid.GridKey{Line: 0, Beat: 0},
			description:       "Should restore note after removing and undoing",
		},
		{
			name: "Undo note modification",
			commands: []any{
				mappings.NoteAdd,
				mappings.AccentIncrease,
				mappings.Undo,
			},
			expectedCursorPos: grid.GridKey{Line: 0, Beat: 0},
			description:       "Should restore original note after modification and undo",
		},
		{
			name: "Undo with cursor movement",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLeft,
				mappings.Undo,
			},
			expectedCursorPos: grid.GridKey{Line: 0, Beat: 1},
			description:       "Should undo note addition and restore cursor position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			// Verify cursor position is restored
			assert.Equal(t, tt.expectedCursorPos, m.cursorPos, tt.description+" - cursor position should be restored")
		})
	}
}

func TestUpdateUndoLineGridNotes(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		cursorPos   grid.GridKey
		description string
	}{
		{
			name: "Undo line modification with multiple notes",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLineStart,
				mappings.ClearLine,
				mappings.Undo,
			},
			cursorPos:   grid.GridKey{Line: 0, Beat: 0},
			description: "Should restore original line state after modification and undo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithCurosrPos(tt.cursorPos),
			)

			m, _ = processCommands(tt.commands, m)

			// Verify cursor position is restored
			assert.Equal(t, tt.cursorPos, m.cursorPos, tt.description+" - cursor position should be restored")

			// Verify notes are restored on the line
			for beat := range uint8(3) {
				_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: tt.cursorPos.Line, Beat: beat})
				assert.True(t, exists, tt.description+" - note should exist at beat %d after undo", beat)
			}
		})
	}
}

func TestUpdateUndoBounds(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		cursorPos         grid.GridKey
		expectedCursorPos grid.GridKey
		description       string
	}{
		{
			name: "Undo bounds operation",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLineStart,
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.NoteRemove,
				mappings.Undo,
			},
			cursorPos:         grid.GridKey{Line: 0, Beat: 0},
			expectedCursorPos: grid.GridKey{Line: 0, Beat: 2},
			description:       "Should restore notes in bounds after removal and undo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithCurosrPos(tt.cursorPos),
			)

			m, _ = processCommands(tt.commands, m)

			// Verify cursor position is restored
			assert.Equal(t, tt.expectedCursorPos, m.cursorPos, tt.description+" - cursor position should be restored")
		})
	}
}

func TestUpdateUndoBeats(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		initialBeats  uint8
		expectedBeats uint8
		description   string
	}{
		{
			name: "Undo beats change",
			commands: []any{
				mappings.BeatInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			initialBeats:  32,
			expectedBeats: 32,
			description:   "Should restore original beats count after modification and undo",
		},
		{
			name: "Undo beats change cycling through selections",
			commands: []any{
				mappings.BeatInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.BeatInputSwitch,
				mappings.BeatInputSwitch,
				mappings.BeatInputSwitch,
				mappings.BeatInputSwitch,
				mappings.Undo,
			},
			initialBeats:  32,
			expectedBeats: 32,
			description:   "Should restore original beats count after modification and undo",
		},
		{
			name: "Escape from arrangement mode does not undo beats",
			commands: []any{
				mappings.ToggleArrangementView,
				mappings.Enter,
				mappings.Undo,
			},
			initialBeats:  32,
			expectedBeats: 32,
			description:   "Should restore original beats count after modification and undo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Verify initial beats
			assert.Equal(t, tt.initialBeats, m.CurrentPart().Beats, "Initial beats should match")

			m, _ = processCommands(tt.commands, m)

			// Verify beats are restored
			assert.Equal(t, tt.expectedBeats, m.CurrentPart().Beats, tt.description+" - beats should be restored")
		})
	}
}

func TestUpdateUndoTempo(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		initialTempo int
		description  string
	}{
		{
			name: "Undo tempo change",
			commands: []any{
				mappings.TempoInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			initialTempo: 120,
			description:  "Should restore original tempo after modification and undo",
		},
		{
			name: "Undo tempo change cycling through selections",
			commands: []any{
				mappings.TempoInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.TempoInputSwitch,
				mappings.TempoInputSwitch,
				mappings.Undo,
			},
			initialTempo: 120,
			description:  "Should restore original tempo after modification and undo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Verify initial tempo
			assert.Equal(t, tt.initialTempo, m.definition.tempo, "Initial tempo should match")

			m, _ = processCommands(tt.commands, m)

			// Verify tempo is restored
			assert.Equal(t, tt.initialTempo, m.definition.tempo, tt.description+" - tempo should be restored")

			assert.Equal(t, m.selectionIndicator, SelectTempo, tt.description+" - selection should be SelectTempo after undo")
		})
	}
}

func TestUpdateUndoSubdivisions(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		initialSubdivision int
		description        string
	}{
		{
			name: "Undo subdivisions change",
			commands: []any{
				mappings.TempoInputSwitch,
				mappings.TempoInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			initialSubdivision: 2,
			description:        "Should restore original subdivisions after modification and undo",
		},
		{
			name: "Undo subdivisions change cycling through selections",
			commands: []any{
				mappings.TempoInputSwitch,
				mappings.TempoInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.TempoInputSwitch,
				mappings.Undo,
			},
			initialSubdivision: 2,
			description:        "Should restore original subdivisions after modification and undo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Verify initial subdivisions
			assert.Equal(t, tt.initialSubdivision, m.definition.subdivisions, "Initial subdivisions should match")

			m, _ = processCommands(tt.commands, m)

			// Verify subdivisions are restored
			assert.Equal(t, tt.initialSubdivision, m.definition.subdivisions, tt.description+" - subdivisions should be restored")
			assert.Equal(t, m.selectionIndicator, SelectTempoSubdivision, tt.description+" - selection should be SelectTempo after undo")
		})
	}
}

func TestUpdateUndoNewOverlay(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "Undo new overlay creation",
			commands: []any{
				mappings.NextOverlay,
				mappings.Undo,
			},
			description: "Should remove newly created overlay after undo",
		},
	}

	overlayKey := overlaykey.OverlayPeriodicity{
		Shift:      2,
		Interval:   4,
		Width:      0,
		StartCycle: 0,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithNonRootOverlay(overlayKey),
			)

			keys := make([]overlaykey.OverlayPeriodicity, 0)
			m.CurrentPart().Overlays.CollectKeys(&keys)
			initialOverlayCount := len(keys)

			m, _ = processCommands(tt.commands, m)

			keys = make([]overlaykey.OverlayPeriodicity, 0)
			m.CurrentPart().Overlays.CollectKeys(&keys)
			finalOverlayCount := len(keys)

			assert.Equal(t, initialOverlayCount, finalOverlayCount, tt.description+" - overlay count should be restored")
		})
	}
}

func TestUndoArrangement(t *testing.T) {
	tests := []struct {
		name          string
		setupCommands []any
		commands      []any
		description   string
	}{
		{
			name:          "Undo arrangement operation",
			setupCommands: []any{},
			commands: []any{
				mappings.NewSectionAfter,
				mappings.Enter,
				mappings.NextSection,
			},
			description: "Should restore arrangement state after modification and undo",
		},
		{
			name:          "Undo arrangement operation new section before",
			setupCommands: []any{},
			commands: []any{
				mappings.NewSectionBefore,
				mappings.Enter,
				mappings.PrevSection,
			},
			description: "Should restore arrangement state",
		},
		{
			name:          "Undo arrangement operation group",
			setupCommands: []any{mappings.ToggleArrangementView},
			commands: []any{
				TestKey{"g"},
				mappings.Enter,
			},
			description: "Should restore arrangement state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.setupCommands, m)

			initialCursor := m.arrangement.Cursor
			initialArrangementNodesCount := m.arrangement.Root.CountEndNodes()

			m, _ = processCommands(tt.commands, m)

			cursorBeforeUndo := m.arrangement.Cursor
			endNodesCountBeforeUndo := m.arrangement.Root.CountEndNodes()

			m, _ = processCommand(mappings.Undo, m)

			assert.Equal(t, initialCursor, m.arrangement.Cursor, tt.description+" - arrangement cursor should be restored")
			assert.Equal(t, initialArrangementNodesCount, m.arrangement.Root.CountEndNodes(), tt.description+" - arrangement nodes count should be restored")
			assert.Equal(t, m.focus, FocusArrangementEditor, tt.description+" - focus should be on arrangement editor after undo")
			assert.Equal(t, m.arrangement.Focus, true, tt.description+" - arrangement should be focused after undo")

			m, _ = processCommand(mappings.Redo, m)

			assert.Equal(t, cursorBeforeUndo, m.arrangement.Cursor, tt.description+" - arrangement cursor should be restored after redo")
			assert.Equal(t, endNodesCountBeforeUndo, m.arrangement.Root.CountEndNodes(), tt.description+" - arrangement nodes count should match after redo")
		})
	}
}

func TestUndoMultipleOperations(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "Multiple undo operations",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.Undo,
				mappings.Undo,
				mappings.Undo,
			},
			description: "Should undo multiple operations in reverse order",
		},
		{
			name: "Undo after complex operations",
			commands: []any{
				mappings.NoteAdd,
				mappings.AccentIncrease,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLineStart,
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.Yank,
				mappings.CursorDown,
				mappings.Paste,
				mappings.Undo,
				mappings.Undo,
				mappings.Undo,
				mappings.Redo,
				mappings.Undo,
			},
			description: "Should undo complex operations correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			// Verify model is in a consistent state after multiple undos
			assert.Equal(t, grid.GridKey{Line: 0, Beat: 0}, m.cursorPos, tt.description+" - cursor should be at origin after all undos")
		})
	}
}

func TestUndoWithArrangementCursor(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "Undo restores arrangement cursor",
			commands: []any{
				mappings.NewSectionAfter,
				mappings.Enter,
				mappings.NoteAdd,
				mappings.NewSectionBefore,
				mappings.Enter,
				mappings.NoteAdd,
				mappings.Undo,
			},
			description: "Should restore arrangement cursor position after undo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			initialArrangementCursor := m.arrangement.Cursor

			m, _ = processCommands(tt.commands, m)

			// Verify arrangement cursor is restored to the position when the undoable operation was performed
			assert.Equal(t, initialArrangementCursor, m.arrangement.Cursor, tt.description+" - arrangement cursor should be restored")
		})
	}
}

func TestUndoEmptyStack(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "Undo with empty stack",
			commands: []any{
				mappings.Undo,
			},
			description: "Should handle undo with empty stack gracefully",
		},
		{
			name: "Multiple undos beyond stack",
			commands: []any{
				mappings.NoteAdd,
				mappings.Undo,
				mappings.Undo,
				mappings.Undo,
			},
			description: "Should handle multiple undos beyond stack size gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			// Verify model is in a consistent state
			assert.Equal(t, grid.GridKey{Line: 0, Beat: 0}, m.cursorPos, tt.description+" - cursor should remain at origin")
		})
	}
}

func TestUndoOverlayDiff(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "Undo overlay diff operation",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLineStart,
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.Yank,
				mappings.CursorDown,
				mappings.Paste,
				mappings.Undo,
			},
			description: "Should restore overlay state after diff operation and undo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			// Verify overlay state is restored
			_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: 1, Beat: 0})
			assert.False(t, exists, tt.description+" - pasted note should not exist after undo")
		})
	}
}

func TestUndoAccentInputSwitch(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		expectedDiff   uint8
		expectedData   []config.Accent
		expectedStart  uint8
		expectedTarget accentTarget
		description    string
	}{
		{
			name: "Undo accent diff change",
			commands: []any{
				mappings.AccentInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			expectedDiff:   15,
			expectedData:   []config.Accent{{Value: 0}, {Value: 120}, {Value: 105}, {Value: 90}, {Value: 75}, {Value: 60}, {Value: 45}, {Value: 30}, {Value: 15}},
			expectedStart:  120,
			expectedTarget: AccentTargetVelocity,
			description:    "Should restore original accent diff after modification and undo",
		},
		{
			name: "Undo accent start change",
			commands: []any{
				mappings.AccentInputSwitch,
				mappings.AccentInputSwitch,
				mappings.AccentInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			expectedDiff:   15,
			expectedData:   []config.Accent{{Value: 0}, {Value: 120}, {Value: 105}, {Value: 90}, {Value: 75}, {Value: 60}, {Value: 45}, {Value: 30}, {Value: 15}},
			expectedStart:  120,
			expectedTarget: AccentTargetVelocity,
			description:    "Should restore original accent start after modification and undo",
		},
		{
			name: "Undo accent target change",
			commands: []any{
				mappings.AccentInputSwitch,
				mappings.AccentInputSwitch,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			expectedDiff:   15,
			expectedData:   []config.Accent{{Value: 0}, {Value: 120}, {Value: 105}, {Value: 90}, {Value: 75}, {Value: 60}, {Value: 45}, {Value: 30}, {Value: 15}},
			expectedStart:  120,
			expectedTarget: AccentTargetVelocity,
			description:    "Should restore original accent target after modification and undo",
		},
		{
			name: "Undo multiple accent changes",
			commands: []any{
				mappings.AccentInputSwitch,
				mappings.Increase,
				mappings.AccentInputSwitch,
				mappings.Increase,
				mappings.AccentInputSwitch,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			expectedDiff:   15,
			expectedData:   []config.Accent{{Value: 0}, {Value: 120}, {Value: 105}, {Value: 90}, {Value: 75}, {Value: 60}, {Value: 45}, {Value: 30}, {Value: 15}},
			expectedStart:  120,
			expectedTarget: AccentTargetVelocity,
			description:    "Should restore original accent values after multiple modifications and undos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Verify initial accent values
			assert.Equal(t, tt.expectedDiff, m.definition.accents.Diff, "Initial accent diff should match")
			assert.Equal(t, tt.expectedStart, m.definition.accents.Start, "Initial accent start should match")
			assert.Equal(t, tt.expectedTarget, m.definition.accents.Target, "Initial accent target should match")

			m, _ = processCommands(tt.commands, m)

			// Verify accent values are restored
			assert.Equal(t, tt.expectedDiff, m.definition.accents.Diff, tt.description+" - accent diff should be restored")
			assert.Equal(t, tt.expectedStart, m.definition.accents.Start, tt.description+" - accent start should be restored")
			assert.Equal(t, tt.expectedTarget, m.definition.accents.Target, tt.description+" - accent target should be restored")
			assert.Equal(t, tt.expectedData, m.definition.accents.Data, tt.description+" - accent data should be restored")
			assert.Equal(t, m.selectionIndicator, SelectAccentDiff, tt.description+" - selection should be SelectAccentDiff after undo")
		})
	}
}

func TestUndoSetupInputSwitch(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedChannel uint8
		expectedNote    uint8
		expectedMsgType grid.MessageType
		description     string
	}{
		{
			name: "Undo channel change",
			commands: []any{
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			expectedChannel: 10,
			expectedNote:    60,
			expectedMsgType: grid.MessageTypeNote,
			description:     "Should restore original channel after modification and undo",
		},
		{
			name: "Undo message type change",
			commands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			expectedChannel: 10,
			expectedNote:    60,
			expectedMsgType: grid.MessageTypeNote,
			description:     "Should restore original message type after modification and undo",
		},
		{
			name: "Undo note value change",
			commands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Increase,
				mappings.SetupInputSwitch,
				mappings.Undo,
			},
			expectedChannel: 10,
			expectedNote:    60,
			expectedMsgType: grid.MessageTypeNote,
			description:     "Should restore original note value after modification and undo",
		},
		{
			name: "Undo multiple setup changes",
			commands: []any{
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
				mappings.Undo,
			},
			expectedChannel: 10,
			expectedNote:    60,
			expectedMsgType: grid.MessageTypeNote,
			description:     "Should restore original setup values after multiple modifications and undos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Verify initial setup values
			currentLine := m.definition.lines[m.cursorPos.Line]
			assert.Equal(t, tt.expectedChannel, currentLine.Channel, "Initial channel should match")
			assert.Equal(t, tt.expectedNote, currentLine.Note, "Initial note should match")
			assert.Equal(t, tt.expectedMsgType, currentLine.MsgType, "Initial message type should match")

			m, _ = processCommands(tt.commands, m)

			// Verify setup values are restored
			restoredLine := m.definition.lines[m.cursorPos.Line]
			assert.Equal(t, tt.expectedChannel, restoredLine.Channel, tt.description+" - channel should be restored")
			assert.Equal(t, tt.expectedNote, restoredLine.Note, tt.description+" - note should be restored")
			assert.Equal(t, tt.expectedMsgType, restoredLine.MsgType, tt.description+" - message type should be restored")
			assert.Equal(t, m.selectionIndicator, SelectSetupChannel, tt.description+" - selection should be SetupInputSwitch after undo")
		})
	}
}
