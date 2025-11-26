package main

import (
	"testing"

	"github.com/chriserin/sq/internal/grid"
	"github.com/chriserin/sq/internal/mappings"
	"github.com/stretchr/testify/assert"
)

func TestToggleBoundedLoop(t *testing.T) {
	tests := []struct {
		name           string
		cursorPos      grid.GridKey
		commands       []any
		expectedActive bool
		expectedLeft   uint8
		expectedRight  uint8
		description    string
	}{
		{
			name:           "Toggle on at cursor position",
			cursorPos:      grid.GridKey{Line: 0, Beat: 5},
			commands:       []any{mappings.ToggleBoundedLoop},
			expectedActive: true,
			expectedLeft:   5,
			expectedRight:  5,
			description:    "Should activate bounded loop with both bounds at cursor position",
		},
		{
			name:           "Toggle off from active state",
			cursorPos:      grid.GridKey{Line: 0, Beat: 10},
			commands:       []any{mappings.ToggleBoundedLoop, mappings.ToggleBoundedLoop},
			expectedActive: false,
			expectedLeft:   10,
			expectedRight:  10,
			description:    "Should deactivate bounded loop but keep bounds",
		},
		{
			name:           "Toggle on at beat 0",
			cursorPos:      grid.GridKey{Line: 0, Beat: 0},
			commands:       []any{mappings.ToggleBoundedLoop},
			expectedActive: true,
			expectedLeft:   0,
			expectedRight:  0,
			description:    "Should activate bounded loop at start of line",
		},
		{
			name:           "Toggle on at last beat",
			cursorPos:      grid.GridKey{Line: 0, Beat: 31},
			commands:       []any{mappings.ToggleBoundedLoop},
			expectedActive: true,
			expectedLeft:   31,
			expectedRight:  31,
			description:    "Should activate bounded loop at end of line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.cursorPos))

			assert.False(t, m.playState.BoundedLoop.Active, "Initial state should be inactive")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedActive, m.playState.BoundedLoop.Active, tt.description+" - active state")
			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}

func TestToggleBoundedLoopWithVisualSelection(t *testing.T) {
	tests := []struct {
		name           string
		cursorStart    grid.GridKey
		commands       []any
		expectedActive bool
		expectedLeft   uint8
		expectedRight  uint8
		description    string
	}{
		{
			name:        "Toggle with visual selection from beat 5 to 10",
			cursorStart: grid.GridKey{Line: 0, Beat: 5},
			commands: []any{
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.ToggleBoundedLoop,
			},
			expectedActive: true,
			expectedLeft:   5,
			expectedRight:  10,
			description:    "Should set bounds to visual selection range",
		},
		{
			name:        "Toggle with visual selection from beat 0 to 3",
			cursorStart: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.ToggleBoundedLoop,
			},
			expectedActive: true,
			expectedLeft:   0,
			expectedRight:  3,
			description:    "Should set bounds to visual selection from start",
		},
		{
			name:        "Toggle with visual selection backwards",
			cursorStart: grid.GridKey{Line: 0, Beat: 10},
			commands: []any{
				mappings.ToggleVisualMode,
				mappings.CursorLeft,
				mappings.CursorLeft,
				mappings.CursorLeft,
				mappings.ToggleBoundedLoop,
			},
			expectedActive: true,
			expectedLeft:   7,
			expectedRight:  10,
			description:    "Should set bounds correctly even with backwards selection",
		},
		{
			name:        "Toggle with single beat visual selection",
			cursorStart: grid.GridKey{Line: 0, Beat: 15},
			commands: []any{
				mappings.ToggleVisualMode,
				mappings.ToggleBoundedLoop,
			},
			expectedActive: true,
			expectedLeft:   15,
			expectedRight:  15,
			description:    "Should handle single beat visual selection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.cursorStart))

			assert.False(t, m.playState.BoundedLoop.Active, "Initial state should be inactive")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedActive, m.playState.BoundedLoop.Active, tt.description+" - active state")
			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}

func TestExpandLeftLoopBound(t *testing.T) {
	tests := []struct {
		name          string
		initialLeft   uint8
		initialRight  uint8
		commands      []any
		expectedLeft  uint8
		expectedRight uint8
		description   string
	}{
		{
			name:          "Expand left bound by one",
			initialLeft:   10,
			initialRight:  15,
			commands:      []any{mappings.ExpandLeftLoopBound},
			expectedLeft:  9,
			expectedRight: 15,
			description:   "Should decrease left bound by 1",
		},
		{
			name:          "Expand left bound multiple times",
			initialLeft:   10,
			initialRight:  15,
			commands:      []any{mappings.ExpandLeftLoopBound, mappings.ExpandLeftLoopBound, mappings.ExpandLeftLoopBound},
			expectedLeft:  7,
			expectedRight: 15,
			description:   "Should decrease left bound by 3",
		},
		{
			name:          "Expand left bound to minimum (0)",
			initialLeft:   2,
			initialRight:  10,
			commands:      []any{mappings.ExpandLeftLoopBound, mappings.ExpandLeftLoopBound, mappings.ExpandLeftLoopBound},
			expectedLeft:  0,
			expectedRight: 10,
			description:   "Should stop at 0 and not go negative",
		},
		{
			name:          "Cannot expand left bound below 0",
			initialLeft:   0,
			initialRight:  10,
			commands:      []any{mappings.ExpandLeftLoopBound, mappings.ExpandLeftLoopBound},
			expectedLeft:  0,
			expectedRight: 10,
			description:   "Should remain at 0 when already at minimum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.playState.BoundedLoop.Active = true
				m.playState.BoundedLoop.LeftBound = tt.initialLeft
				m.playState.BoundedLoop.RightBound = tt.initialRight
				return *m
			})

			assert.Equal(t, tt.initialLeft, m.playState.BoundedLoop.LeftBound, "Initial left bound")
			assert.Equal(t, tt.initialRight, m.playState.BoundedLoop.RightBound, "Initial right bound")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}

func TestExpandRightLoopBound(t *testing.T) {
	tests := []struct {
		name          string
		initialLeft   uint8
		initialRight  uint8
		beats         uint8
		commands      []any
		expectedLeft  uint8
		expectedRight uint8
		description   string
	}{
		{
			name:          "Expand right bound by one",
			initialLeft:   10,
			initialRight:  15,
			beats:         32,
			commands:      []any{mappings.ExpandRightLoopBound},
			expectedLeft:  10,
			expectedRight: 16,
			description:   "Should increase right bound by 1",
		},
		{
			name:          "Expand right bound multiple times",
			initialLeft:   10,
			initialRight:  15,
			beats:         32,
			commands:      []any{mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound},
			expectedLeft:  10,
			expectedRight: 18,
			description:   "Should increase right bound by 3",
		},
		{
			name:          "Expand right bound to maximum",
			initialLeft:   10,
			initialRight:  29,
			beats:         32,
			commands:      []any{mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound},
			expectedLeft:  10,
			expectedRight: 32,
			description:   "Should stop at beats limit",
		},
		{
			name:          "Cannot expand right bound beyond beats",
			initialLeft:   10,
			initialRight:  32,
			beats:         32,
			commands:      []any{mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound},
			expectedLeft:  10,
			expectedRight: 32,
			description:   "Should remain at beats when already at maximum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridSize(int(tt.beats), 8),
				func(m *model) model {
					m.playState.BoundedLoop.Active = true
					m.playState.BoundedLoop.LeftBound = tt.initialLeft
					m.playState.BoundedLoop.RightBound = tt.initialRight
					return *m
				},
			)

			assert.Equal(t, tt.initialLeft, m.playState.BoundedLoop.LeftBound, "Initial left bound")
			assert.Equal(t, tt.initialRight, m.playState.BoundedLoop.RightBound, "Initial right bound")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}

func TestContractLeftLoopBound(t *testing.T) {
	tests := []struct {
		name          string
		initialLeft   uint8
		initialRight  uint8
		commands      []any
		expectedLeft  uint8
		expectedRight uint8
		description   string
	}{
		{
			name:          "Contract left bound by one",
			initialLeft:   10,
			initialRight:  15,
			commands:      []any{mappings.ContractLeftLoopBound},
			expectedLeft:  11,
			expectedRight: 15,
			description:   "Should increase left bound by 1",
		},
		{
			name:          "Contract left bound multiple times",
			initialLeft:   10,
			initialRight:  15,
			commands:      []any{mappings.ContractLeftLoopBound, mappings.ContractLeftLoopBound, mappings.ContractLeftLoopBound},
			expectedLeft:  13,
			expectedRight: 15,
			description:   "Should increase left bound by 3",
		},
		{
			name:          "Contract left bound to meet right bound",
			initialLeft:   10,
			initialRight:  13,
			commands:      []any{mappings.ContractLeftLoopBound, mappings.ContractLeftLoopBound, mappings.ContractLeftLoopBound},
			expectedLeft:  13,
			expectedRight: 13,
			description:   "Should stop one before right bound",
		},
		{
			name:          "Cannot contract left bound past right bound",
			initialLeft:   14,
			initialRight:  15,
			commands:      []any{mappings.ContractLeftLoopBound, mappings.ContractLeftLoopBound},
			expectedLeft:  15,
			expectedRight: 15,
			description:   "Should remain at right-1 when already at maximum contraction",
		},
		{
			name:          "Contract when bounds are equal",
			initialLeft:   10,
			initialRight:  10,
			commands:      []any{mappings.ContractLeftLoopBound},
			expectedLeft:  10,
			expectedRight: 10,
			description:   "Should not change when left equals right",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.playState.BoundedLoop.Active = true
				m.playState.BoundedLoop.LeftBound = tt.initialLeft
				m.playState.BoundedLoop.RightBound = tt.initialRight
				return *m
			})

			assert.Equal(t, tt.initialLeft, m.playState.BoundedLoop.LeftBound, "Initial left bound")
			assert.Equal(t, tt.initialRight, m.playState.BoundedLoop.RightBound, "Initial right bound")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}

func TestContractRightLoopBound(t *testing.T) {
	tests := []struct {
		name          string
		initialLeft   uint8
		initialRight  uint8
		commands      []any
		expectedLeft  uint8
		expectedRight uint8
		description   string
	}{
		{
			name:          "Contract right bound by one",
			initialLeft:   10,
			initialRight:  15,
			commands:      []any{mappings.ContractRightLoopBound},
			expectedLeft:  10,
			expectedRight: 14,
			description:   "Should decrease right bound by 1",
		},
		{
			name:          "Contract right bound multiple times",
			initialLeft:   10,
			initialRight:  15,
			commands:      []any{mappings.ContractRightLoopBound, mappings.ContractRightLoopBound, mappings.ContractRightLoopBound},
			expectedLeft:  10,
			expectedRight: 12,
			description:   "Should decrease right bound by 3",
		},
		{
			name:          "Contract right bound to meet left bound",
			initialLeft:   10,
			initialRight:  13,
			commands:      []any{mappings.ContractRightLoopBound, mappings.ContractRightLoopBound, mappings.ContractRightLoopBound},
			expectedLeft:  10,
			expectedRight: 10,
			description:   "Should stop one after left bound",
		},
		{
			name:          "Cannot contract right bound past left bound",
			initialLeft:   10,
			initialRight:  11,
			commands:      []any{mappings.ContractRightLoopBound, mappings.ContractRightLoopBound},
			expectedLeft:  10,
			expectedRight: 10,
			description:   "Should remain at left+1 when already at maximum contraction",
		},
		{
			name:          "Contract when bounds are equal",
			initialLeft:   10,
			initialRight:  10,
			commands:      []any{mappings.ContractRightLoopBound},
			expectedLeft:  10,
			expectedRight: 10,
			description:   "Should not change when left equals right",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.playState.BoundedLoop.Active = true
				m.playState.BoundedLoop.LeftBound = tt.initialLeft
				m.playState.BoundedLoop.RightBound = tt.initialRight
				return *m
			})

			assert.Equal(t, tt.initialLeft, m.playState.BoundedLoop.LeftBound, "Initial left bound")
			assert.Equal(t, tt.initialRight, m.playState.BoundedLoop.RightBound, "Initial right bound")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}

func TestBoundedLoopComplexScenarios(t *testing.T) {
	tests := []struct {
		name          string
		initialLeft   uint8
		initialRight  uint8
		commands      []any
		expectedLeft  uint8
		expectedRight uint8
		description   string
	}{
		{
			name:         "Expand both bounds from center",
			initialLeft:  15,
			initialRight: 15,
			commands: []any{
				mappings.ExpandLeftLoopBound,
				mappings.ExpandLeftLoopBound,
				mappings.ExpandLeftLoopBound,
				mappings.ExpandRightLoopBound,
				mappings.ExpandRightLoopBound,
				mappings.ExpandRightLoopBound,
			},
			expectedLeft:  12,
			expectedRight: 18,
			description:   "Should expand in both directions symmetrically",
		},
		{
			name:         "Expand then contract left bound",
			initialLeft:  10,
			initialRight: 20,
			commands: []any{
				mappings.ExpandLeftLoopBound,
				mappings.ExpandLeftLoopBound,
				mappings.ContractLeftLoopBound,
			},
			expectedLeft:  9,
			expectedRight: 20,
			description:   "Should expand left by 2 then contract by 1",
		},
		{
			name:         "Expand then contract right bound",
			initialLeft:  10,
			initialRight: 20,
			commands: []any{
				mappings.ExpandRightLoopBound,
				mappings.ExpandRightLoopBound,
				mappings.ContractRightLoopBound,
			},
			expectedLeft:  10,
			expectedRight: 21,
			description:   "Should expand right by 2 then contract by 1",
		},
		{
			name:         "Contract to single beat then expand",
			initialLeft:  10,
			initialRight: 15,
			commands: []any{
				mappings.ContractLeftLoopBound,
				mappings.ContractLeftLoopBound,
				mappings.ContractLeftLoopBound,
				mappings.ContractLeftLoopBound,
				mappings.ContractLeftLoopBound,
				mappings.ExpandLeftLoopBound,
				mappings.ExpandRightLoopBound,
			},
			expectedLeft:  14,
			expectedRight: 16,
			description:   "Should contract to near-single beat then expand outward",
		},
		{
			name:         "All operations combined",
			initialLeft:  16,
			initialRight: 16,
			commands: []any{
				mappings.ExpandLeftLoopBound,
				mappings.ExpandLeftLoopBound,
				mappings.ExpandRightLoopBound,
				mappings.ExpandRightLoopBound,
				mappings.ExpandRightLoopBound,
				mappings.ContractLeftLoopBound,
				mappings.ContractRightLoopBound,
			},
			expectedLeft:  15,
			expectedRight: 18,
			description:   "Should handle complex sequence of operations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.playState.BoundedLoop.Active = true
				m.playState.BoundedLoop.LeftBound = tt.initialLeft
				m.playState.BoundedLoop.RightBound = tt.initialRight
				return *m
			})

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}

func TestBoundedLoopWithDifferentGridSizes(t *testing.T) {
	tests := []struct {
		name          string
		beats         int
		initialLeft   uint8
		initialRight  uint8
		commands      []any
		expectedLeft  uint8
		expectedRight uint8
		description   string
	}{
		{
			name:          "16 beat grid - expand right to limit",
			beats:         16,
			initialLeft:   0,
			initialRight:  10,
			commands:      []any{mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound, mappings.ExpandRightLoopBound},
			expectedLeft:  0,
			expectedRight: 16,
			description:   "Should expand right to 16 beat limit",
		},
		{
			name:          "64 beat grid - expand operations",
			beats:         64,
			initialLeft:   30,
			initialRight:  35,
			commands:      []any{mappings.ExpandLeftLoopBound, mappings.ExpandRightLoopBound},
			expectedLeft:  29,
			expectedRight: 36,
			description:   "Should work correctly with larger grid",
		},
		{
			name:          "8 beat grid - full range",
			beats:         8,
			initialLeft:   0,
			initialRight:  8,
			commands:      []any{mappings.ExpandLeftLoopBound, mappings.ExpandRightLoopBound},
			expectedLeft:  0,
			expectedRight: 8,
			description:   "Should respect boundaries in smaller grid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithGridSize(tt.beats, 8),
				func(m *model) model {
					m.playState.BoundedLoop.Active = true
					m.playState.BoundedLoop.LeftBound = tt.initialLeft
					m.playState.BoundedLoop.RightBound = tt.initialRight
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedLeft, m.playState.BoundedLoop.LeftBound, tt.description+" - left bound")
			assert.Equal(t, tt.expectedRight, m.playState.BoundedLoop.RightBound, tt.description+" - right bound")
		})
	}
}
