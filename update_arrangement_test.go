package main

import (
	"testing"

	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

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
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)
			assert.Equal(t, tt.expectedFocus, m.focus, tt.description+" - arrangement view state")
			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedArrIsOpen, m.showArrangementView, tt.description+" - arrangement open state")
		})
	}
}

func TestRenamePartCommand(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name:        "RenamePart changes part name",
			commands:    []any{mappings.ToggleArrangementView, TestKey{Keys: "R"}, TestKey{Keys: "XYZ"}, mappings.Enter},
			description: "RenamePart command should update part name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// Verify initial state
			assert.Equal(t, SelectNothing, m.selectionIndicator, "Initial selection should be nothing")

			// Execute commands
			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, "XYZ", m.CurrentPart().Name, tt.description+" - part name should be updated")
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
			expectedPartIndex:  1,
			description:        "Should move to next section successfully",
		},
		{
			name:               "PrevSection moves cursor backward",
			commands:           []any{mappings.NewSectionBefore, mappings.Enter, mappings.PrevSection},
			expectedMoveResult: true,
			expectedPartIndex:  1,
			description:        "Should move to previous section successfully",
		},
		{
			name:               "NextSection on single section",
			commands:           []any{mappings.NextSection},
			expectedMoveResult: false,
			expectedPartIndex:  0,
			description:        "Should not move when only one section exists",
		},
		{
			name:               "PrevSection on first section",
			commands:           []any{mappings.NewSectionAfter, mappings.PrevSection},
			expectedMoveResult: false,
			expectedPartIndex:  0,
			description:        "Should not move when already on first section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			initialArrangementCursor := m.arrangement.Cursor

			m, _ = processCommands(tt.commands, m)

			arrangementCursorMoved := !m.arrangement.Cursor.Equals(initialArrangementCursor)

			assert.Equal(t, tt.expectedPartIndex, m.CurrentSongSection().Part, tt.description+" - part index should match")
			assert.Equal(t, tt.expectedMoveResult, arrangementCursorMoved, tt.description+" - cursor movement")
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

func TestChangePartAfterNewSectionAfter(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedPartIndex int
		description       string
	}{
		{
			name:              "ChangePart after NewSectionAfter updates part mapping",
			commands:          []any{mappings.NewSectionAfter, mappings.Enter, mappings.ChangePart, mappings.Increase, mappings.Increase, mappings.Enter},
			expectedPartIndex: 1,
			description:       "Should change part mapping to index 1 after creating new section",
		},
		{
			name:              "ChangePart after NewSectionAfter updates part mapping, more increases stays at top part",
			commands:          []any{mappings.NewSectionAfter, mappings.Enter, mappings.ChangePart, mappings.Increase, mappings.Increase, mappings.Increase, mappings.Enter},
			expectedPartIndex: 1,
			description:       "Should change part mapping to index 1 after creating new section",
		},
		{
			name:              "ChangePart with decrease after NewSectionAfter",
			commands:          []any{mappings.NewSectionAfter, mappings.Enter, mappings.ChangePart, mappings.Increase, mappings.Increase, mappings.Decrease, mappings.Enter},
			expectedPartIndex: 0,
			description:       "Should change part mapping to index 0 after increase then decrease",
		},
		{
			name:              "ChangePart creates new part as default operation",
			commands:          []any{mappings.NewSectionAfter, mappings.Enter, mappings.ChangePart, mappings.Enter},
			expectedPartIndex: 2,
			description:       "Should create new part when selecting beyond existing parts",
		},
		{
			name:              "ChangePart creates new part as default operation, decrease stays on 'Create New Part'",
			commands:          []any{mappings.NewSectionAfter, mappings.Enter, mappings.ChangePart, mappings.Decrease, mappings.Enter},
			expectedPartIndex: 2,
			description:       "Should create new part when selecting beyond existing parts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedPartIndex, m.CurrentSongSection().Part, tt.description+" - part index should match")

			initialPartCount := 1
			if tt.expectedPartIndex >= initialPartCount {
				expectedPartCount := tt.expectedPartIndex + 1
				assert.Equal(t, expectedPartCount, len(*m.definition.parts), tt.description+" - new parts should be created")
			}
		})
	}
}
