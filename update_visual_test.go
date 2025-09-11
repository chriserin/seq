package main

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/operation"
	"github.com/stretchr/testify/assert"
)

func TestToggleVisualMode(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		initialVisualMode  operation.VisualMode
		expectedVisualMode operation.VisualMode
		description        string
	}{
		{
			name:               "Toggle from VisualNone to VisualBlock",
			commands:           []any{mappings.ToggleVisualMode},
			initialVisualMode:  operation.VisualNone,
			expectedVisualMode: operation.VisualBlock,
			description:        "Should toggle from VisualNone to VisualBlock",
		},
		{
			name:               "Toggle from VisualBlock to VisualNone",
			commands:           []any{mappings.ToggleVisualMode},
			initialVisualMode:  operation.VisualBlock,
			expectedVisualMode: operation.VisualNone,
			description:        "Should toggle from VisualBlock to VisualNone",
		},
		{
			name:               "Toggle from VisualLine to VisualBlock",
			commands:           []any{mappings.ToggleVisualMode},
			initialVisualMode:  operation.VisualLine,
			expectedVisualMode: operation.VisualBlock,
			description:        "Should toggle from VisualLine to VisualBlock",
		},
		{
			name:               "Multiple toggles cycle correctly",
			commands:           []any{mappings.ToggleVisualMode, mappings.ToggleVisualMode},
			initialVisualMode:  operation.VisualNone,
			expectedVisualMode: operation.VisualNone,
			description:        "Should return to VisualNone after two toggles",
		},
		{
			name:               "Toggle VisualMode sets anchor and lead to current cursor",
			commands:           []any{mappings.ToggleVisualMode},
			initialVisualMode:  operation.VisualNone,
			expectedVisualMode: operation.VisualBlock,
			description:        "Should set anchor and lead to current cursor position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridCursor(GK(2, 3)),
				func(m *model) model {
					m.visualSelection.visualMode = tt.initialVisualMode
					return *m
				},
			)

			assert.Equal(t, tt.initialVisualMode, m.visualSelection.visualMode, "Initial visual mode should match")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedVisualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match expected value")

			if tt.expectedVisualMode != operation.VisualNone {
				assert.Equal(t, GK(2, 3), m.visualSelection.anchor, tt.description+" - anchor should be set to cursor position")
				assert.Equal(t, GK(2, 3), m.visualSelection.lead, tt.description+" - lead should be set to cursor position")
			}
		})
	}
}

func TestToggleVisualLineMode(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		initialVisualMode  operation.VisualMode
		expectedVisualMode operation.VisualMode
		description        string
	}{
		{
			name:               "Toggle from VisualNone to VisualLine",
			commands:           []any{mappings.ToggleVisualLineMode},
			initialVisualMode:  operation.VisualNone,
			expectedVisualMode: operation.VisualLine,
			description:        "Should toggle from VisualNone to VisualLine",
		},
		{
			name:               "Toggle from VisualLine to VisualNone",
			commands:           []any{mappings.ToggleVisualLineMode},
			initialVisualMode:  operation.VisualLine,
			expectedVisualMode: operation.VisualNone,
			description:        "Should toggle from VisualLine to VisualNone",
		},
		{
			name:               "Toggle from VisualBlock to VisualLine",
			commands:           []any{mappings.ToggleVisualLineMode},
			initialVisualMode:  operation.VisualBlock,
			expectedVisualMode: operation.VisualLine,
			description:        "Should toggle from VisualBlock to VisualLine",
		},
		{
			name:               "Multiple toggles cycle correctly",
			commands:           []any{mappings.ToggleVisualLineMode, mappings.ToggleVisualLineMode},
			initialVisualMode:  operation.VisualNone,
			expectedVisualMode: operation.VisualNone,
			description:        "Should return to VisualNone after two toggles",
		},
		{
			name:               "Toggle VisualLineMode sets anchor and lead to current cursor",
			commands:           []any{mappings.ToggleVisualLineMode},
			initialVisualMode:  operation.VisualNone,
			expectedVisualMode: operation.VisualLine,
			description:        "Should set anchor and lead to current cursor position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridCursor(GK(1, 4)),
				func(m *model) model {
					m.visualSelection.visualMode = tt.initialVisualMode
					return *m
				},
			)

			assert.Equal(t, tt.initialVisualMode, m.visualSelection.visualMode, "Initial visual mode should match")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedVisualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match expected value")

			if tt.expectedVisualMode != operation.VisualNone {
				assert.Equal(t, GK(1, 4), m.visualSelection.anchor, tt.description+" - anchor should be set to cursor position")
				assert.Equal(t, GK(1, 4), m.visualSelection.lead, tt.description+" - lead should be set to cursor position")
			}
		})
	}
}

func TestVisualModeInteractions(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		expectedVisualMode operation.VisualMode
		description        string
	}{
		{
			name:               "ToggleVisualMode then ToggleVisualLineMode",
			commands:           []any{mappings.ToggleVisualMode, mappings.ToggleVisualLineMode},
			expectedVisualMode: operation.VisualLine,
			description:        "Should end up in VisualLine mode after toggling both",
		},
		{
			name:               "ToggleVisualLineMode then ToggleVisualMode",
			commands:           []any{mappings.ToggleVisualLineMode, mappings.ToggleVisualMode},
			expectedVisualMode: operation.VisualBlock,
			description:        "Should end up in VisualBlock mode after toggling both",
		},
		{
			name:               "Complex visual mode cycling",
			commands:           []any{mappings.ToggleVisualMode, mappings.ToggleVisualLineMode, mappings.ToggleVisualMode, mappings.ToggleVisualLineMode},
			expectedVisualMode: operation.VisualNone,
			description:        "Should end up in VisualNone mode after complex cycling",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridCursor(GK(0, 0)),
			)

			assert.Equal(t, operation.VisualNone, m.visualSelection.visualMode, "Initial visual mode should be VisualNone")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedVisualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match expected value")
		})
	}
}

func TestVisualSelectionBounds(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		cursorPos      gridKey
		expectedBounds grid.Bounds
		visualMode     operation.VisualMode
		description    string
	}{
		{
			name:           "VisualBlock bounds match cursor position initially",
			commands:       []any{mappings.ToggleVisualMode},
			cursorPos:      GK(2, 5),
			expectedBounds: grid.Bounds{Top: 2, Left: 5, Bottom: 2, Right: 5},
			visualMode:     operation.VisualBlock,
			description:    "Should create single-cell bounds at cursor position for VisualBlock",
		},
		{
			name:           "VisualLine bounds span entire line",
			commands:       []any{mappings.ToggleVisualLineMode},
			cursorPos:      GK(3, 7),
			expectedBounds: grid.Bounds{Top: 3, Left: 0, Bottom: 3, Right: 15},
			visualMode:     operation.VisualLine,
			description:    "Should create line-spanning bounds for VisualLine",
		},
		{
			name:           "VisualLine bounds span entire lines after move",
			commands:       []any{mappings.ToggleVisualLineMode, mappings.CursorDown, mappings.CursorDown},
			cursorPos:      GK(3, 7),
			expectedBounds: grid.Bounds{Top: 3, Left: 0, Bottom: 5, Right: 15},
			visualMode:     operation.VisualLine,
			description:    "Should create line-spanning bounds for VisualLine",
		},
		{
			name:           "VisualBlock with cursor movement expands selection",
			commands:       []any{mappings.ToggleVisualMode, mappings.CursorRight, mappings.CursorDown},
			cursorPos:      GK(1, 2),
			expectedBounds: grid.Bounds{Top: 1, Left: 2, Bottom: 2, Right: 3},
			visualMode:     operation.VisualBlock,
			description:    "Should expand VisualBlock selection with cursor movement",
		},
		{
			name:           "Switch from visual block to visual line",
			commands:       []any{mappings.ToggleVisualMode, mappings.CursorRight, mappings.CursorDown, mappings.ToggleVisualLineMode},
			cursorPos:      GK(1, 2),
			expectedBounds: grid.Bounds{Top: 1, Left: 0, Bottom: 2, Right: 15},
			visualMode:     operation.VisualLine,
			description:    "Should expand VisualBlock selection with cursor movement",
		},
		{
			name:           "Switch from visual line to visual block",
			commands:       []any{mappings.ToggleVisualLineMode, mappings.CursorRight, mappings.CursorDown, mappings.ToggleVisualMode},
			cursorPos:      GK(1, 2),
			expectedBounds: grid.Bounds{Top: 1, Left: 2, Bottom: 2, Right: 3},
			visualMode:     operation.VisualBlock,
			description:    "Should expand VisualBlock selection with cursor movement",
		},
		{
			name:           "Switch from visual block to visual none",
			commands:       []any{mappings.ToggleVisualMode, mappings.CursorRight, mappings.CursorDown, mappings.ToggleVisualMode},
			cursorPos:      GK(1, 2),
			expectedBounds: grid.Bounds{Top: 1, Left: 0, Bottom: 2, Right: 15},
			visualMode:     operation.VisualNone,
			description:    "Should expand VisualBlock selection with cursor movement",
		},
		{
			name:           "Switch from visual line to visual none",
			commands:       []any{mappings.ToggleVisualLineMode, mappings.CursorRight, mappings.CursorDown, mappings.ToggleVisualLineMode},
			cursorPos:      GK(1, 2),
			expectedBounds: grid.Bounds{Top: 1, Left: 2, Bottom: 2, Right: 3},
			visualMode:     operation.VisualNone,
			description:    "Should expand VisualBlock selection with cursor movement",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridCursor(tt.cursorPos),
				WithGridSize(16, 8),
			)

			assert.Equal(t, operation.VisualNone, m.visualSelection.visualMode, "Initial visual mode should be VisualNone")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.visualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match expected value")

			if m.visualSelection.visualMode != operation.VisualNone {
				assert.Equal(t, tt.expectedBounds, m.VisualSelectionBounds(), tt.description+" - visual selection bounds should match expected")
			}
		})
	}
}

func TestVisualModeAndYank(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		expectedVisualMode operation.VisualMode
		description        string
	}{
		{
			name:               "Yank resets visual mode to VisualNone",
			commands:           []any{mappings.ToggleVisualMode, mappings.Yank},
			expectedVisualMode: operation.VisualNone,
			description:        "Should reset visual mode to VisualNone after yank",
		},
		{
			name:               "Yank from VisualLine resets to VisualNone",
			commands:           []any{mappings.ToggleVisualLineMode, mappings.Yank},
			expectedVisualMode: operation.VisualNone,
			description:        "Should reset visual mode to VisualNone after yank from VisualLine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridCursor(GK(0, 0)),
			)

			assert.Equal(t, operation.VisualNone, m.visualSelection.visualMode, "Initial visual mode should be VisualNone")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedVisualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match expected value")
		})
	}
}

func TestVisualModeAndNoteRemove(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		expectedVisualMode operation.VisualMode
		description        string
	}{
		{
			name:               "NoteRemove resets visual mode to VisualNone",
			commands:           []any{mappings.ToggleVisualMode, mappings.NoteRemove},
			expectedVisualMode: operation.VisualNone,
			description:        "Should reset visual mode to VisualNone after NoteRemove",
		},
		{
			name:               "NoteRemove from VisualLine resets to VisualNone",
			commands:           []any{mappings.ToggleVisualLineMode, mappings.NoteRemove},
			expectedVisualMode: operation.VisualNone,
			description:        "Should reset visual mode to VisualNone after NoteRemove from VisualLine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridCursor(GK(0, 0)),
			)

			assert.Equal(t, operation.VisualNone, m.visualSelection.visualMode, "Initial visual mode should be VisualNone")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedVisualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match expected value")
		})
	}
}

