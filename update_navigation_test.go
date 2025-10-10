package main

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/operation"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

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
		{
			name:        "Cursor Last Line",
			command:     mappings.CursorLastLine,
			initialPos:  grid.GridKey{Line: 3, Beat: 5},
			expectedPos: grid.GridKey{Line: 7, Beat: 5},
			description: "Cursor should move to last line",
		},
		{
			name:        "Cursor Last Line From Last Line",
			command:     mappings.CursorLastLine,
			initialPos:  grid.GridKey{Line: 7, Beat: 10},
			expectedPos: grid.GridKey{Line: 7, Beat: 10},
			description: "Cursor should stay on last line when already there",
		},
		{
			name:        "Cursor First Line",
			command:     mappings.CursorFirstLine,
			initialPos:  grid.GridKey{Line: 4, Beat: 8},
			expectedPos: grid.GridKey{Line: 0, Beat: 8},
			description: "Cursor should move to first line",
		},
		{
			name:        "Cursor First Line From First Line",
			command:     mappings.CursorFirstLine,
			initialPos:  grid.GridKey{Line: 0, Beat: 12},
			expectedPos: grid.GridKey{Line: 0, Beat: 12},
			description: "Cursor should stay on first line when already there",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.initialPos))
			m, _ = processCommand(tt.command, m)
			assert.Equal(t, tt.expectedPos.Line, m.gridCursor.Line, tt.description+" - line position")
			assert.Equal(t, tt.expectedPos.Beat, m.gridCursor.Beat, tt.description+" - beat position")
		})
	}
}

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
				{Shift: 1, Interval: 2, Width: 1, StartCycle: 0},
				{Shift: 1, Interval: 3, Width: 1, StartCycle: 0},
			},
			expectedOverlayKey: overlaykey.OverlayPeriodicity{Shift: 1, Interval: 2, Width: 1, StartCycle: 0},
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
				{Shift: 1, Interval: 2, Width: 1, StartCycle: 0},
				{Shift: 1, Interval: 3, Width: 1, StartCycle: 0},
			},
			expectedOverlayKey: overlaykey.OverlayPeriodicity{Shift: 1, Interval: 2, Width: 1, StartCycle: 0},
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
				{Shift: 1, Interval: 2, Width: 1, StartCycle: 0},
				{Shift: 1, Interval: 3, Width: 1, StartCycle: 0},
			},
			expectedOverlayKey: overlaykey.OverlayPeriodicity{Shift: 1, Interval: 1, Width: 1, StartCycle: 0},
			description:        "Should switch to next overlay with key 1/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			rootKey := overlaykey.ROOT

			for _, key := range tt.addedOverlayKeys {
				(*m.definition.Parts)[0].Overlays = m.CurrentPart().Overlays.Add(key)
			}

			assert.Equal(t, rootKey, m.currentOverlay.Key, "Initial overlay key should be root")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedOverlayKey, m.currentOverlay.Key, tt.description)
		})
	}
}

func TestOverlayInputSwitch(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedSelection operation.Selection
		expectedFocus     operation.Focus
		description       string
	}{
		{
			name:              "Single OverlayInputSwitch",
			commands:          []any{mappings.OverlayInputSwitch},
			expectedSelection: operation.SelectGrid,
			expectedFocus:     operation.FocusOverlayKey,
			description:       "First overlay input switch should select overlay and set focus to overlay key",
		},
		{
			name:              "Double OverlayInputSwitch",
			commands:          []any{mappings.OverlayInputSwitch, mappings.OverlayInputSwitch},
			expectedSelection: operation.SelectGrid,
			expectedFocus:     operation.FocusGrid,
			description:       "Second overlay input switch should cycle back to nothing but keep overlay key focus",
		},
		{
			name:              "Escape Overlay Input state",
			commands:          []any{mappings.OverlayInputSwitch, mappings.Escape},
			expectedSelection: operation.SelectGrid,
			expectedFocus:     operation.FocusGrid,
			description:       "Escape should set focus and selection back to an initial state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")
			assert.Equal(t, operation.FocusGrid, m.focus, "Initial focus should be grid")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedFocus, m.focus, tt.description+" - focus state")
		})
	}
}

func TestNavigateWithHiddenLines(t *testing.T) {
	tests := []struct {
		name           string
		initialPos     grid.GridKey
		notesOnLines   []uint8
		commands       []any
		expectedPos    grid.GridKey
		expectedHidden bool
		description    string
	}{
		{
			name:           "Cursor Down With Hidden Lines",
			initialPos:     grid.GridKey{Line: 0, Beat: 0},
			notesOnLines:   []uint8{0, 2, 5},
			commands:       []any{mappings.ToggleHideLines, mappings.CursorDown},
			expectedPos:    grid.GridKey{Line: 2, Beat: 0},
			expectedHidden: true,
			description:    "Should skip hidden line 1 and move to line 2",
		},
		{
			name:           "Cursor Up With Hidden Lines",
			initialPos:     grid.GridKey{Line: 5, Beat: 0},
			notesOnLines:   []uint8{0, 2, 5},
			commands:       []any{mappings.ToggleHideLines, mappings.CursorUp},
			expectedPos:    grid.GridKey{Line: 2, Beat: 0},
			expectedHidden: true,
			description:    "Should skip hidden lines 3 and 4 and move to line 2",
		},
		{
			name:           "Multiple Cursor Down With Hidden Lines",
			initialPos:     grid.GridKey{Line: 0, Beat: 0},
			notesOnLines:   []uint8{0, 2, 5},
			commands:       []any{mappings.ToggleHideLines, mappings.CursorDown, mappings.CursorDown},
			expectedPos:    grid.GridKey{Line: 5, Beat: 0},
			expectedHidden: true,
			description:    "Should skip all hidden lines and move to line 5",
		},
		{
			name:           "Cursor First Line With Hidden Lines",
			initialPos:     grid.GridKey{Line: 5, Beat: 0},
			notesOnLines:   []uint8{2, 5},
			commands:       []any{mappings.ToggleHideLines, mappings.CursorFirstLine},
			expectedPos:    grid.GridKey{Line: 2, Beat: 0},
			expectedHidden: true,
			description:    "Should move to first visible line (2)",
		},
		{
			name:           "Cursor Last Line With Hidden Lines",
			initialPos:     grid.GridKey{Line: 0, Beat: 0},
			notesOnLines:   []uint8{0, 2, 5},
			commands:       []any{mappings.ToggleHideLines, mappings.CursorLastLine},
			expectedPos:    grid.GridKey{Line: 5, Beat: 0},
			expectedHidden: true,
			description:    "Should move to last visible line (5)",
		},
		{
			name:           "Toggle Hide Lines On Empty Line",
			initialPos:     grid.GridKey{Line: 3, Beat: 0},
			notesOnLines:   []uint8{0, 2, 5},
			commands:       []any{mappings.ToggleHideLines},
			expectedPos:    grid.GridKey{Line: 2, Beat: 0},
			expectedHidden: true,
			description:    "Should move cursor to nearest visible line when toggling hide on empty line",
		},
		{
			name:           "Toggle Hide Lines Off",
			initialPos:     grid.GridKey{Line: 2, Beat: 0},
			notesOnLines:   []uint8{0, 2, 5},
			commands:       []any{mappings.ToggleHideLines, mappings.ToggleHideLines},
			expectedPos:    grid.GridKey{Line: 2, Beat: 0},
			expectedHidden: false,
			description:    "Should restore all lines when toggling hide off",
		},
		{
			name:           "Toggle Hide Line Moves Cursor",
			initialPos:     grid.GridKey{Line: 0, Beat: 0},
			notesOnLines:   []uint8{2},
			commands:       []any{mappings.ToggleHideLines},
			expectedPos:    grid.GridKey{Line: 2, Beat: 0},
			expectedHidden: true,
			description:    "Should restore all lines when toggling hide off",
		},
		{
			name:           "Toggle Hide Line Moves Cursor note on 5",
			initialPos:     grid.GridKey{Line: 0, Beat: 0},
			notesOnLines:   []uint8{5},
			commands:       []any{mappings.ToggleHideLines},
			expectedPos:    grid.GridKey{Line: 5, Beat: 0},
			expectedHidden: true,
			description:    "Should restore all lines when toggling hide off",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.initialPos))

			// Add notes to specific lines to make them visible
			for _, lineNum := range tt.notesOnLines {
				key := grid.GridKey{Line: lineNum, Beat: 0}
				m.currentOverlay.AddNote(key, grid.InitNote())
			}

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedPos.Line, m.gridCursor.Line, tt.description+" - line position")
			assert.Equal(t, tt.expectedPos.Beat, m.gridCursor.Beat, tt.description+" - beat position")
			assert.Equal(t, tt.expectedHidden, m.hideEmptyLines, tt.description+" - hide state")
		})
	}
}
