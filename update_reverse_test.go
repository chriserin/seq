package main

import (
	"testing"

	"github.com/chriserin/sq/internal/grid"
	"github.com/chriserin/sq/internal/mappings"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name          string
		setupNotes    []grid.GridKey        // positions to add notes at
		cursorPos     grid.GridKey          // where to position cursor before reverse
		commands      []any                 // commands to execute (including Reverse)
		expectedNotes map[grid.GridKey]bool // positions that should have notes after reverse
		description   string
	}{
		{
			name: "Reverse three notes from cursor to end",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 1},
				{Line: 0, Beat: 2},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 29}: true, // last note moved to end
				{Line: 0, Beat: 30}: true, // middle note
				{Line: 0, Beat: 31}: true, // first note moved to end
			},
			description: "Should reverse three notes from cursor to end of line",
		},
		{
			name: "Reverse with visual selection",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 1},
				{Line: 0, Beat: 2},
				{Line: 0, Beat: 3},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 0}: true, // was at beat 2
				{Line: 0, Beat: 1}: true, // was at beat 1
				{Line: 0, Beat: 2}: true, // was at beat 0
				{Line: 0, Beat: 3}: true, // unchanged (outside visual selection)
			},
			description: "Should reverse notes within visual selection only",
		},
		{
			name: "Reverse with gaps",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 2},
				{Line: 0, Beat: 4},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 27}: true,
				{Line: 0, Beat: 29}: true,
				{Line: 0, Beat: 31}: true,
			},
			description: "Should reverse notes preserving gaps",
		},
		{
			name: "Reverse from middle of line",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 5},
				{Line: 0, Beat: 6},
				{Line: 0, Beat: 7},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 5},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 29}: true,
				{Line: 0, Beat: 30}: true,
				{Line: 0, Beat: 31}: true,
			},
			description: "Should reverse notes from cursor to end, starting mid-line",
		},
		{
			name: "Reverse single note",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 31}: true,
			},
			description: "Should move single note to end of line",
		},
		{
			name:       "Reverse with no notes",
			setupNotes: []grid.GridKey{},
			cursorPos:  grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{},
			description:   "Should do nothing when no notes present",
		},
		{
			name: "Reverse multiple lines with visual line mode",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 1},
				{Line: 1, Beat: 0},
				{Line: 1, Beat: 1},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.ToggleVisualLineMode,
				mappings.CursorDown,
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 30}: true,
				{Line: 0, Beat: 31}: true,
				{Line: 1, Beat: 30}: true,
				{Line: 1, Beat: 31}: true,
			},
			description: "Should reverse notes on multiple lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.cursorPos))

			// Set up notes at specified positions
			for _, pos := range tt.setupNotes {
				m.SetGridCursor(pos)
				m, _ = processCommand(mappings.NoteAdd, m)
			}

			// Position cursor and execute commands
			m.SetGridCursor(tt.cursorPos)
			m, _ = processCommands(tt.commands, m)

			// Verify expected notes exist
			for pos := range tt.expectedNotes {
				note, exists := m.currentOverlay.GetNote(pos)
				assert.True(t, exists && note != zeronote, tt.description+" - note should exist at beat %d", pos.Beat)
			}

			// Verify no unexpected notes exist
			for beat := uint8(0); beat < 32; beat++ {
				pos := grid.GridKey{Line: tt.cursorPos.Line, Beat: beat}
				note, exists := m.currentOverlay.GetNote(pos)
				hasNote := exists && note != zeronote
				shouldHaveNote := tt.expectedNotes[pos]

				if shouldHaveNote {
					assert.True(t, hasNote, tt.description+" - note should exist at beat %d", beat)
				} else {
					assert.False(t, hasNote, tt.description+" - note should not exist at beat %d", beat)
				}
			}
		})
	}
}

func TestReverseWithChords(t *testing.T) {
	tests := []struct {
		name             string
		setupNotes       []grid.GridKey
		chordPos         grid.GridKey
		cursorPos        grid.GridKey
		commands         []any
		expectedNotes    map[grid.GridKey]bool
		chordStillExists bool
		description      string
	}{
		{
			name: "Reverse notes with chord present - chord unchanged",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 2},
				{Line: 0, Beat: 4},
			},
			chordPos:  grid.GridKey{Line: 0, Beat: 1},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 1}:  true, // chord position unchanged
				{Line: 0, Beat: 27}: true, // note from beat 4 mirrored
				{Line: 0, Beat: 29}: true, // note from beat 2 mirrored
				{Line: 0, Beat: 31}: true, // note from beat 0 mirrored
			},
			chordStillExists: true,
			description:      "Should reverse standalone notes but not affect chord",
		},
		{
			name:       "Reverse with only chords - no change",
			setupNotes: []grid.GridKey{},
			chordPos:   grid.GridKey{Line: 0, Beat: 0},
			cursorPos:  grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 0}: true, // chord unchanged
			},
			chordStillExists: true,
			description:      "Should not move chords",
		},
		{
			name: "Reverse notes before and after chord",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 1},
				{Line: 0, Beat: 5},
				{Line: 0, Beat: 6},
			},
			chordPos:  grid.GridKey{Line: 0, Beat: 3},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 3}:  true, // chord unchanged
				{Line: 0, Beat: 25}: true, // note from beat 6 mirrored: 31-(6-0)=25
				{Line: 0, Beat: 26}: true, // note from beat 5 mirrored: 31-(5-0)=26
				{Line: 0, Beat: 30}: true, // note from beat 1 mirrored: 31-(1-0)=30
				{Line: 0, Beat: 31}: true, // note from beat 0 mirrored: 31-(0-0)=31
			},
			chordStillExists: true,
			description:      "Should reverse all standalone notes while keeping chord in place",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.cursorPos))

			// Set up notes at specified positions
			for _, pos := range tt.setupNotes {
				m.SetGridCursor(pos)
				m, _ = processCommand(mappings.NoteAdd, m)
			}

			// Add chord
			m.SetGridCursor(tt.chordPos)
			m, _ = processCommand(mappings.ToggleChordMode, m)
			m, _ = processCommand(mappings.MajorTriad, m)
			m, _ = processCommand(mappings.ToggleChordMode, m)

			// Position cursor and execute commands
			m.SetGridCursor(tt.cursorPos)
			m, _ = processCommands(tt.commands, m)

			// Verify chord still exists
			chord, chordExists := m.currentOverlay.Chords.FindChordWithNote(tt.chordPos)
			if tt.chordStillExists {
				assert.True(t, chordExists, tt.description+" - chord should still exist")
				assert.Equal(t, tt.chordPos, chord.Root, tt.description+" - chord should be at original position")
			}

			// Verify expected notes/chords exist
			for pos := range tt.expectedNotes {
				// Check if it's a chord position
				if pos == tt.chordPos {
					_, exists := m.currentOverlay.Chords.FindChordWithNote(pos)
					assert.True(t, exists, tt.description+" - chord should exist at beat %d", pos.Beat)
				} else {
					note, exists := m.currentOverlay.GetNote(pos)
					assert.True(t, exists && note != zeronote, tt.description+" - note should exist at beat %d", pos.Beat)
				}
			}
		})
	}
}

func TestReversePreservesNoteAttributes(t *testing.T) {
	m := createTestModel(WithGridCursor(grid.GridKey{Line: 0, Beat: 0}))

	// Add note with specific attributes
	m, _ = processCommand(mappings.NoteAdd, m)
	m, _ = processCommand(mappings.AccentIncrease, m)
	m, _ = processCommand(mappings.AccentIncrease, m)
	m, _ = processCommand(mappings.GateIncrease, m)
	m, _ = processCommand(mappings.RatchetIncrease, m)

	originalNote, _ := m.currentOverlay.GetNote(grid.GridKey{Line: 0, Beat: 0})

	// Add another note
	m.SetGridCursor(grid.GridKey{Line: 0, Beat: 1})
	m, _ = processCommand(mappings.NoteAdd, m)

	// Reverse
	m.SetGridCursor(grid.GridKey{Line: 0, Beat: 0})
	m, _ = processCommand(mappings.Reverse, m)

	// Check that the note with attributes is now at the end
	reversedNote, exists := m.currentOverlay.GetNote(grid.GridKey{Line: 0, Beat: 31})
	assert.True(t, exists, "Reversed note should exist at beat 31")
	assert.Equal(t, originalNote.AccentIndex, reversedNote.AccentIndex, "Should preserve accent")
	assert.Equal(t, originalNote.GateIndex, reversedNote.GateIndex, "Should preserve gate")
	assert.Equal(t, originalNote.Ratchets, reversedNote.Ratchets, "Should preserve ratchets")
}

func TestReverseWithNonRootOverlay(t *testing.T) {
	tests := []struct {
		name          string
		setupCommands []any
		rootNotes     []grid.GridKey
		overlayNotes  []grid.GridKey
		cursorPos     grid.GridKey
		commands      []any
		expectedNotes map[grid.GridKey]bool
		description   string
	}{
		{
			name: "Reverse notes on non-root overlay",
			setupCommands: []any{
				mappings.OverlayInputSwitch,
				TestKey{Keys: "2"},
				mappings.Enter,
			},
			rootNotes: []grid.GridKey{
				{Line: 0, Beat: 10},
				{Line: 0, Beat: 11},
			},
			overlayNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 1},
				{Line: 0, Beat: 2},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 20}: true, // root note from beat 11 reversed: 31-(11-0)=20
				{Line: 0, Beat: 21}: true, // root note from beat 10 reversed: 31-(10-0)=21
				{Line: 0, Beat: 29}: true, // overlay note from beat 2 reversed: 31-(2-0)=29
				{Line: 0, Beat: 30}: true, // overlay note from beat 1 reversed: 31-(1-0)=30
				{Line: 0, Beat: 31}: true, // overlay note from beat 0 reversed: 31-(0-0)=31
			},
			description: "Should reverse all visible notes including root overlay notes",
		},
		{
			name: "Reverse with visual selection on non-root overlay",
			setupCommands: []any{
				mappings.OverlayInputSwitch,
				TestKey{Keys: "3"},
				mappings.Enter,
			},
			rootNotes: []grid.GridKey{
				{Line: 0, Beat: 5},
			},
			overlayNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 1},
				{Line: 0, Beat: 2},
				{Line: 0, Beat: 3},
			},
			cursorPos: grid.GridKey{Line: 0, Beat: 1},
			commands: []any{
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 0}: true, // outside visual selection, unchanged
				{Line: 0, Beat: 1}: true, // was at beat 3
				{Line: 0, Beat: 2}: true, // was at beat 2
				{Line: 0, Beat: 3}: true, // was at beat 1
				{Line: 0, Beat: 5}: true, // root overlay note unchanged
			},
			description: "Should reverse notes in visual selection on non-root overlay",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(grid.GridKey{Line: 0, Beat: 0}))

			// Add notes on root overlay
			for _, pos := range tt.rootNotes {
				m.SetGridCursor(pos)
				m, _ = processCommand(mappings.NoteAdd, m)
			}

			// Switch to non-root overlay
			m, _ = processCommands(tt.setupCommands, m)

			// Add notes on current overlay
			for _, pos := range tt.overlayNotes {
				m.SetGridCursor(pos)
				m, _ = processCommand(mappings.NoteAdd, m)
			}

			// Position cursor and execute reverse
			m.SetGridCursor(tt.cursorPos)
			m, _ = processCommands(tt.commands, m)

			// Get combined pattern to check all visible notes
			combinedPattern := m.CombinedEditPattern(m.currentOverlay)

			// Verify expected notes exist in combined pattern
			for pos := range tt.expectedNotes {
				note := combinedPattern[pos]
				assert.NotEqual(t, zeronote, note, tt.description+" - note should exist at beat %d", pos.Beat)
			}

			// Verify no unexpected notes exist
			for beat := uint8(0); beat < 32; beat++ {
				pos := grid.GridKey{Line: 0, Beat: beat}
				note := combinedPattern[pos]
				hasNote := note != zeronote
				shouldHaveNote := tt.expectedNotes[pos]

				if shouldHaveNote {
					assert.True(t, hasNote, tt.description+" - note should exist at beat %d", beat)
				} else {
					assert.False(t, hasNote, tt.description+" - note should not exist at beat %d", beat)
				}
			}
		})
	}
}

func TestReverseMultipleTimes(t *testing.T) {
	tests := []struct {
		name          string
		setupNotes    []grid.GridKey
		cursorPos     grid.GridKey
		reverseCount  int
		expectedNotes map[grid.GridKey]bool
		description   string
	}{
		{
			name: "Reverse twice returns to original positions",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 1},
				{Line: 0, Beat: 2},
			},
			cursorPos:    grid.GridKey{Line: 0, Beat: 0},
			reverseCount: 2,
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 0}: true,
				{Line: 0, Beat: 1}: true,
				{Line: 0, Beat: 2}: true,
			},
			description: "Two consecutive reverses should return notes to original positions",
		},
		{
			name: "Reverse three times same as one reverse",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 5},
				{Line: 0, Beat: 10},
			},
			cursorPos:    grid.GridKey{Line: 0, Beat: 0},
			reverseCount: 3,
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 21}: true, // 31-(10-0) = 21
				{Line: 0, Beat: 26}: true, // 31-(5-0) = 26
			},
			description: "Three consecutive reverses should have same result as one reverse",
		},
		{
			name: "Reverse four times returns to original",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 1},
				{Line: 0, Beat: 3},
				{Line: 0, Beat: 7},
			},
			cursorPos:    grid.GridKey{Line: 0, Beat: 0},
			reverseCount: 4,
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 1}: true,
				{Line: 0, Beat: 3}: true,
				{Line: 0, Beat: 7}: true,
			},
			description: "Four consecutive reverses should return to original positions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.cursorPos))

			// Set up notes
			for _, pos := range tt.setupNotes {
				m.SetGridCursor(pos)
				m, _ = processCommand(mappings.NoteAdd, m)
			}

			// Position cursor
			m.SetGridCursor(tt.cursorPos)

			// Apply reverse multiple times
			for i := 0; i < tt.reverseCount; i++ {
				m, _ = processCommand(mappings.Reverse, m)
			}

			// Verify expected notes
			for pos := range tt.expectedNotes {
				note, exists := m.currentOverlay.GetNote(pos)
				assert.True(t, exists && note != zeronote, tt.description+" - note should exist at beat %d", pos.Beat)
			}

			// Verify no unexpected notes
			for beat := uint8(0); beat < 32; beat++ {
				pos := grid.GridKey{Line: tt.cursorPos.Line, Beat: beat}
				note, exists := m.currentOverlay.GetNote(pos)
				hasNote := exists && note != zeronote
				shouldHaveNote := tt.expectedNotes[pos]

				if shouldHaveNote {
					assert.True(t, hasNote, tt.description+" - note should exist at beat %d", beat)
				} else {
					assert.False(t, hasNote, tt.description+" - note should not exist at beat %d", beat)
				}
			}
		})
	}
}

func TestReverseWithActions(t *testing.T) {
	tests := []struct {
		name           string
		setupNotes     []grid.GridKey
		actionPos      grid.GridKey
		actionCommand  mappings.Command
		cursorPos      grid.GridKey
		commands       []any
		expectedNotes  map[grid.GridKey]bool
		expectedAction grid.GridKey
		description    string
	}{
		{
			name: "Reverse notes with action - action reversed with notes",
			setupNotes: []grid.GridKey{
				{Line: 0, Beat: 0},
				{Line: 0, Beat: 2},
			},
			actionPos:     grid.GridKey{Line: 0, Beat: 1},
			actionCommand: mappings.ActionAddLineReset,
			cursorPos:     grid.GridKey{Line: 0, Beat: 0},
			commands: []any{
				mappings.Reverse,
			},
			expectedNotes: map[grid.GridKey]bool{
				{Line: 0, Beat: 29}: true, // note from beat 2 mirrored: 31-(2-0)=29
				{Line: 0, Beat: 30}: true, // action from beat 1 mirrored: 31-(1-0)=30
				{Line: 0, Beat: 31}: true, // note from beat 0 mirrored: 31-(0-0)=31
			},
			expectedAction: grid.GridKey{Line: 0, Beat: 30},
			description:    "Should reverse action along with notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.cursorPos))

			// Set up notes
			for _, pos := range tt.setupNotes {
				m.SetGridCursor(pos)
				m, _ = processCommand(mappings.NoteAdd, m)
			}

			// Add action
			m.SetGridCursor(tt.actionPos)
			m, _ = processCommand(tt.actionCommand, m)

			// Position cursor and execute reverse
			m.SetGridCursor(tt.cursorPos)
			m, _ = processCommands(tt.commands, m)

			// Verify expected notes exist
			for pos := range tt.expectedNotes {
				note, exists := m.currentOverlay.GetNote(pos)
				assert.True(t, exists && note != zeronote, tt.description+" - note should exist at beat %d", pos.Beat)
			}

			// Verify action is at expected position
			actionNote, exists := m.currentOverlay.GetNote(tt.expectedAction)
			assert.True(t, exists, tt.description+" - action note should exist")
			assert.NotEqual(t, grid.ActionNothing, actionNote.Action, tt.description+" - should have an action")
		})
	}
}
