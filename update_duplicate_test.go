package main

import (
	"slices"
	"testing"

	"github.com/chriserin/sq/internal/grid"
	"github.com/chriserin/sq/internal/mappings"
	"github.com/chriserin/sq/internal/operation"
	"github.com/stretchr/testify/assert"
)

func TestDuplicateSingleNote(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		description       string
		expectedNoteBeats []uint8
		cursorLine        uint8
		cursorBeat        uint8
	}{
		{
			name: "Duplicate single note to next beat",
			commands: []any{
				mappings.NoteAdd,
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{0, 1},
			cursorLine:        0,
			cursorBeat:        1,
			description:       "Should duplicate note to next beat and move cursor right",
		},
		{
			name: "Duplicate at end of grid should not crash",
			commands: []any{
				mappings.CursorLineEnd,
				mappings.NoteAdd,
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{31},
			cursorLine:        0,
			cursorBeat:        31,
			description:       "Should not duplicate when at last beat",
		},
		{
			name: "Duplicate multiple times creates chain",
			commands: []any{
				mappings.NoteAdd,
				mappings.Duplicate,
				mappings.Duplicate,
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{0, 1, 2, 3},
			cursorLine:        0,
			cursorBeat:        3,
			description:       "Should allow repeated duplication",
		},
		{
			name: "Duplicate note with GateIndex 16 (2 beats)",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease, // +8
				mappings.GateBigIncrease, // +8, total 16
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{0, 2},
			cursorLine:        0,
			cursorBeat:        2,
			description:       "Should duplicate note 2 beats away when GateIndex is 16",
		},
		{
			name: "Duplicate note with GateIndex 24 (3 beats)",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease, // +8
				mappings.GateBigIncrease, // +8
				mappings.GateBigIncrease, // +8, total 24
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{0, 3},
			cursorLine:        0,
			cursorBeat:        3,
			description:       "Should duplicate note 3 beats away when GateIndex is 24",
		},
		{
			name: "Duplicate note with GateIndex 12 (1.5 beats, rounds up to 2)",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease, // +8
				mappings.GateIncrease,    // +1
				mappings.GateIncrease,    // +1
				mappings.GateIncrease,    // +1
				mappings.GateIncrease,    // +1, total 12
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{0, 2},
			cursorLine:        0,
			cursorBeat:        2,
			description:       "Should duplicate note 2 beats away when GateIndex is 12 (rounds up)",
		},
		{
			name: "Duplicate note with GateIndex 9 (rounds up to 2 beats)",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease, // +7 for first note
				mappings.GateIncrease,    // +1, total 8
				mappings.GateIncrease,    // +1, total 9
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{0, 2},
			cursorLine:        0,
			cursorBeat:        2,
			description:       "Should duplicate note 2 beats away when GateIndex is 9",
		},
		{
			name: "Duplicate long note chain",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease, // +8
				mappings.GateBigIncrease, // +8, total 16 (2 beats)
				mappings.Duplicate,
				mappings.Duplicate,
			},
			expectedNoteBeats: []uint8{0, 2, 4},
			cursorLine:        0,
			cursorBeat:        4,
			description:       "Should allow repeated duplication of long notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			// Check cursor position
			assert.Equal(t, tt.cursorLine, m.gridCursor.Line, tt.description+" - cursor line should match")
			assert.Equal(t, tt.cursorBeat, m.gridCursor.Beat, tt.description+" - cursor beat should match")

			// Check notes on the current line
			for beat := range uint8(32) {
				_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: tt.cursorLine, Beat: beat})
				shouldExist := slices.Contains(tt.expectedNoteBeats, beat)
				if shouldExist {
					assert.True(t, exists, tt.description+" - note should exist at beat %d", beat)
				} else {
					assert.False(t, exists, tt.description+" - note should not exist at beat %d", beat)
				}
			}
		})
	}
}

func TestDuplicateVisualSelection(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		description        string
		expectedNotesLine0 []uint8
		expectedNotesLine1 []uint8
		cursorLine         uint8
		cursorBeat         uint8
		visualMode         operation.VisualMode
		visualAnchorLine   uint8
		visualAnchorBeat   uint8
		visualLeadLine     uint8
		visualLeadBeat     uint8
	}{
		{
			name: "Duplicate visual block selection",
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
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0, 1, 2, 3, 4, 5},
			cursorLine:         0,
			cursorBeat:         3,
			visualMode:         operation.VisualBlock,
			visualAnchorLine:   0,
			visualAnchorBeat:   3,
			visualLeadLine:     0,
			visualLeadBeat:     5,
			description:        "Should duplicate visual block to area immediately after",
		},
		{
			name: "Duplicate visual block with multiple lines",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorDown,
				mappings.CursorLineStart,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorUp,
				mappings.CursorLineStart,
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorDown,
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0, 1, 2, 3},
			expectedNotesLine1: []uint8{0, 1, 2, 3},
			cursorLine:         0,
			cursorBeat:         2,
			visualMode:         operation.VisualBlock,
			visualAnchorLine:   0,
			visualAnchorBeat:   2,
			visualLeadLine:     1,
			visualLeadBeat:     3,
			description:        "Should duplicate visual block with multiple lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			// Check cursor position
			assert.Equal(t, tt.cursorLine, m.gridCursor.Line, tt.description+" - cursor line should match")
			assert.Equal(t, tt.cursorBeat, m.gridCursor.Beat, tt.description+" - cursor beat should match")

			// Check visual selection
			assert.Equal(t, tt.visualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match")
			assert.Equal(t, tt.visualAnchorLine, m.visualSelection.anchor.Line, tt.description+" - visual anchor line should match")
			assert.Equal(t, tt.visualAnchorBeat, m.visualSelection.anchor.Beat, tt.description+" - visual anchor beat should match")
			assert.Equal(t, tt.visualLeadLine, m.visualSelection.lead.Line, tt.description+" - visual lead line should match")
			assert.Equal(t, tt.visualLeadBeat, m.visualSelection.lead.Beat, tt.description+" - visual lead beat should match")

			// Check notes on line 0
			for beat := range uint8(32) {
				_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: 0, Beat: beat})
				shouldExist := slices.Contains(tt.expectedNotesLine0, beat)
				if shouldExist {
					assert.True(t, exists, tt.description+" - note should exist on line 0 at beat %d", beat)
				} else {
					assert.False(t, exists, tt.description+" - note should not exist on line 0 at beat %d", beat)
				}
			}

			// Check notes on line 1 if expected
			if tt.expectedNotesLine1 != nil {
				for beat := range uint8(32) {
					_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: 1, Beat: beat})
					shouldExist := slices.Contains(tt.expectedNotesLine1, beat)
					if shouldExist {
						assert.True(t, exists, tt.description+" - note should exist on line 1 at beat %d", beat)
					} else {
						assert.False(t, exists, tt.description+" - note should not exist on line 1 at beat %d", beat)
					}
				}
			}
		})
	}
}

func TestDuplicateSingleLine(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		description        string
		expectedNotesLine0 []uint8
		expectedNotesLine1 []uint8
		cursorLine         uint8
		cursorBeat         uint8
		visualMode         operation.VisualMode
		visualAnchorLine   uint8
		visualLeadLine     uint8
	}{
		{
			name: "Duplicate single line to next line",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.ToggleVisualLineMode,
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0, 1, 2},
			expectedNotesLine1: []uint8{0, 1, 2},
			cursorLine:         1,
			cursorBeat:         0,
			visualMode:         operation.VisualLine,
			visualAnchorLine:   1,
			visualLeadLine:     1,
			description:        "Should duplicate entire line to next line",
		},
		{
			name: "Duplicate line with gaps",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.ToggleVisualLineMode,
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0, 2, 4},
			expectedNotesLine1: []uint8{0, 2, 4},
			cursorLine:         1,
			cursorBeat:         0,
			visualMode:         operation.VisualLine,
			visualAnchorLine:   1,
			visualLeadLine:     1,
			description:        "Should duplicate line preserving gaps between notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			// Check cursor position
			assert.Equal(t, tt.cursorLine, m.gridCursor.Line, tt.description+" - cursor line should match")
			assert.Equal(t, tt.cursorBeat, m.gridCursor.Beat, tt.description+" - cursor beat should match")

			// Check visual selection
			assert.Equal(t, tt.visualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match")
			assert.Equal(t, tt.visualAnchorLine, m.visualSelection.anchor.Line, tt.description+" - visual anchor line should match")
			assert.Equal(t, tt.visualLeadLine, m.visualSelection.lead.Line, tt.description+" - visual lead line should match")

			// Check notes on line 0
			for beat := range uint8(32) {
				_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: 0, Beat: beat})
				shouldExist := slices.Contains(tt.expectedNotesLine0, beat)
				if shouldExist {
					assert.True(t, exists, tt.description+" - note should exist on line 0 at beat %d", beat)
				} else {
					assert.False(t, exists, tt.description+" - note should not exist on line 0 at beat %d", beat)
				}
			}

			// Check notes on line 1
			for beat := range uint8(32) {
				_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: 1, Beat: beat})
				shouldExist := slices.Contains(tt.expectedNotesLine1, beat)
				if shouldExist {
					assert.True(t, exists, tt.description+" - note should exist on line 1 at beat %d", beat)
				} else {
					assert.False(t, exists, tt.description+" - note should not exist on line 1 at beat %d", beat)
				}
			}
		})
	}
}

func TestDuplicateMultipleLines(t *testing.T) {
	tests := []struct {
		name               string
		setupFunc          modelFunc
		commands           []any
		description        string
		expectedNotesLine0 []uint8
		expectedNotesLine1 []uint8
		expectedNotesLine2 []uint8
		expectedNotesLine3 []uint8
		cursorLine         uint8
		cursorBeat         uint8
		visualMode         operation.VisualMode
		visualAnchorLine   uint8
		visualLeadLine     uint8
	}{
		{
			name:      "Duplicate two lines to next two lines",
			setupFunc: WithGridSize(32, 10),
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorDown,
				mappings.CursorLineStart,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorUp,
				mappings.ToggleVisualLineMode,
				mappings.CursorDown,
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0, 1},
			expectedNotesLine1: []uint8{0, 2},
			expectedNotesLine2: []uint8{0, 1},
			expectedNotesLine3: []uint8{0, 2},
			cursorLine:         2,
			cursorBeat:         0,
			visualMode:         operation.VisualLine,
			visualAnchorLine:   2,
			visualLeadLine:     3,
			description:        "Should duplicate two lines below the selection",
		},
		{
			name:      "Duplicate three lines to next three lines",
			setupFunc: WithGridSize(32, 10),
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorDown,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorDown,
				mappings.CursorLineStart,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorUp,
				mappings.CursorUp,
				mappings.ToggleVisualLineMode,
				mappings.CursorDown,
				mappings.CursorDown,
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0},
			expectedNotesLine1: []uint8{0, 1},
			expectedNotesLine2: []uint8{0, 2},
			expectedNotesLine3: []uint8{0},
			cursorLine:         3,
			cursorBeat:         0,
			visualMode:         operation.VisualLine,
			visualAnchorLine:   3,
			visualLeadLine:     5,
			description:        "Should duplicate three lines below the selection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(tt.setupFunc)

			m, _ = processCommands(tt.commands, m)

			// Check cursor position
			assert.Equal(t, tt.cursorLine, m.gridCursor.Line, tt.description+" - cursor line should match")
			assert.Equal(t, tt.cursorBeat, m.gridCursor.Beat, tt.description+" - cursor beat should match")

			// Check visual selection
			assert.Equal(t, tt.visualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match")
			assert.Equal(t, tt.visualAnchorLine, m.visualSelection.anchor.Line, tt.description+" - visual anchor line should match")
			assert.Equal(t, tt.visualLeadLine, m.visualSelection.lead.Line, tt.description+" - visual lead line should match")

			// Check notes on each line
			lines := []struct {
				lineNum       uint8
				expectedNotes []uint8
			}{
				{0, tt.expectedNotesLine0},
				{1, tt.expectedNotesLine1},
				{2, tt.expectedNotesLine2},
				{3, tt.expectedNotesLine3},
			}

			for _, line := range lines {
				if line.expectedNotes != nil {
					for beat := range uint8(32) {
						_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: line.lineNum, Beat: beat})
						shouldExist := slices.Contains(line.expectedNotes, beat)
						if shouldExist {
							assert.True(t, exists, tt.description+" - note should exist on line %d at beat %d", line.lineNum, beat)
						} else {
							assert.False(t, exists, tt.description+" - note should not exist on line %d at beat %d", line.lineNum, beat)
						}
					}
				}
			}
		})
	}
}

func TestDuplicateLinesWithLimitedSpace(t *testing.T) {
	tests := []struct {
		name               string
		setupFunc          modelFunc
		commands           []any
		description        string
		expectedNotesLine0 []uint8
		expectedNotesLine1 []uint8
		expectedNotesLine2 []uint8
		expectedNotesLine3 []uint8
		cursorLine         uint8
		cursorBeat         uint8
		visualMode         operation.VisualMode
		visualAnchorLine   uint8
		visualLeadLine     uint8
	}{
		{
			name:      "Duplicate two lines when only one line available",
			setupFunc: WithGridSize(32, 3),
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorDown,
				mappings.CursorLineStart,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorUp,
				mappings.ToggleVisualLineMode,
				mappings.CursorDown,
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0, 1},
			expectedNotesLine1: []uint8{0, 2},
			expectedNotesLine2: []uint8{0, 1},
			expectedNotesLine3: nil,
			cursorLine:         2,
			cursorBeat:         0,
			visualMode:         operation.VisualLine,
			visualAnchorLine:   2,
			visualLeadLine:     2,
			description:        "Should duplicate only one line when space is limited",
		},
		{
			name:      "Duplicate three lines when only two lines available",
			setupFunc: WithGridSize(32, 5),
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorDown,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorDown,
				mappings.CursorLineStart,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorUp,
				mappings.CursorUp,
				mappings.ToggleVisualLineMode,
				mappings.CursorDown,
				mappings.CursorDown,
				mappings.Duplicate,
			},
			expectedNotesLine0: []uint8{0},
			expectedNotesLine1: []uint8{0, 1},
			expectedNotesLine2: []uint8{0, 2},
			expectedNotesLine3: []uint8{0},
			cursorLine:         3,
			cursorBeat:         0,
			visualMode:         operation.VisualLine,
			visualAnchorLine:   3,
			visualLeadLine:     4,
			description:        "Should duplicate only two lines when space is limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(tt.setupFunc)

			m, _ = processCommands(tt.commands, m)

			// Check cursor position
			assert.Equal(t, tt.cursorLine, m.gridCursor.Line, tt.description+" - cursor line should match")
			assert.Equal(t, tt.cursorBeat, m.gridCursor.Beat, tt.description+" - cursor beat should match")

			// Check visual selection
			assert.Equal(t, tt.visualMode, m.visualSelection.visualMode, tt.description+" - visual mode should match")
			assert.Equal(t, tt.visualAnchorLine, m.visualSelection.anchor.Line, tt.description+" - visual anchor line should match")
			assert.Equal(t, tt.visualLeadLine, m.visualSelection.lead.Line, tt.description+" - visual lead line should match")

			// Check notes on each line
			lines := []struct {
				lineNum       uint8
				expectedNotes []uint8
			}{
				{0, tt.expectedNotesLine0},
				{1, tt.expectedNotesLine1},
				{2, tt.expectedNotesLine2},
				{3, tt.expectedNotesLine3},
			}

			for _, line := range lines {
				if line.expectedNotes != nil {
					for beat := range uint8(32) {
						_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: line.lineNum, Beat: beat})
						shouldExist := slices.Contains(line.expectedNotes, beat)
						if shouldExist {
							assert.True(t, exists, tt.description+" - note should exist on line %d at beat %d", line.lineNum, beat)
						} else {
							assert.False(t, exists, tt.description+" - note should not exist on line %d at beat %d", line.lineNum, beat)
						}
					}
				}
			}
		})
	}
}
