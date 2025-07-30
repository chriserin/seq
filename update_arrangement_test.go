package main

import (
	"testing"

	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/operation"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

func TestToggleArrangementViewSwitch(t *testing.T) {

	tests := []struct {
		name              string
		commands          []any
		expectedFocus     operation.Focus
		expectedSelection operation.Selection
		expectedArrIsOpen bool
		description       string
	}{
		{
			name:              "Toggle Arrangement View",
			commands:          []any{mappings.ToggleArrangementView},
			expectedFocus:     operation.FocusArrangementEditor,
			expectedSelection: operation.SelectGrid,
			expectedArrIsOpen: true,
			description:       "First toggle should open arrangement view",
		},
		{
			name:              "Toggle Arrangement View Switch Back to Grid",
			commands:          []any{mappings.ToggleArrangementView, mappings.ToggleArrangementView},
			expectedFocus:     operation.FocusGrid,
			expectedSelection: operation.SelectGrid,
			expectedArrIsOpen: false,
			description:       "Second toggle should switch back to grid and close arrangement view",
		},
		{
			name:              "Toggle Arrangement View Switch to Grid keep Arrangement Open",
			commands:          []any{mappings.ToggleArrangementView, mappings.Escape},
			expectedFocus:     operation.FocusGrid,
			expectedSelection: operation.SelectGrid,
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
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			// Execute commands
			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, "XYZ", m.CurrentPart().Name, tt.description+" - part name should be updated")
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
			commands:    []any{mappings.NewSectionAfter, mappings.Enter},
			description: "Should reset current overlay after moving to next section",
		},
		{
			name:        "PrevSection resets current overlay",
			commands:    []any{mappings.NewSectionBefore, mappings.Enter},
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

func TestGroupKeybinding(t *testing.T) {
	tests := []struct {
		name           string
		setupCommands  []any
		commands       []any
		expectedGroups int
		expectedNodes  int
		description    string
	}{
		{
			name:           "Group keybinding creates group with current and next node",
			setupCommands:  []any{mappings.ToggleArrangementView, mappings.NewSectionAfter, mappings.Enter, mappings.NewSectionAfter, mappings.Enter, mappings.CursorUp, mappings.CursorUp},
			commands:       []any{TestKey{Keys: "g"}},
			expectedGroups: 1,
			expectedNodes:  2,
			description:    "Should create a group containing current and next node",
		},
		{
			name:           "Group keybinding at end of list groups with itself",
			setupCommands:  []any{mappings.ToggleArrangementView, mappings.NewSectionAfter, mappings.Enter, mappings.NewSectionAfter, mappings.Enter},
			commands:       []any{TestKey{Keys: "g"}},
			expectedGroups: 1,
			expectedNodes:  1,
			description:    "Should create a group containing the last node with itself",
		},
		{
			name:           "Group keybinding with single node",
			setupCommands:  []any{mappings.ToggleArrangementView},
			commands:       []any{TestKey{Keys: "g"}},
			expectedGroups: 1,
			expectedNodes:  1,
			description:    "Should create a group containing single node with itself",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.setupCommands, m)

			m, _ = processCommands(tt.commands, m)

			groupCount := m.arrangement.Root.CountGroupNodes()
			currentNode := m.arrangement.CurrentNode()

			assert.True(t, currentNode.IsGroup(), tt.description+" - current node should be a group")
			assert.Equal(t, tt.expectedGroups, groupCount, tt.description+" - group count")
			assert.Equal(t, tt.expectedNodes, currentNode.CountEndNodes(), tt.description+" - group node should contain expected number of children")
		})
	}
}
