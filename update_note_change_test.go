package main

import (
	"testing"

	"github.com/chriserin/sq/internal/config"
	"github.com/chriserin/sq/internal/grid"
	"github.com/chriserin/sq/internal/mappings"
	"github.com/chriserin/sq/internal/operation"
	"github.com/chriserin/sq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

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
				mappings.NoteAdd,
				mappings.AccentIncrease,
			},
			expectedAccent: 4,
			description:    "Should add note and increase accent by 1 (index 5 -> 4)",
		},
		{
			name: "Add note and increase accent twice",
			commands: []any{
				mappings.NoteAdd,
				mappings.AccentIncrease,
				mappings.AccentIncrease,
			},
			expectedAccent: 3,
			description:    "Should add note and increase accent by 2 (index 5 -> 3)",
		},
		{
			name: "Add note and increase accent at boundary",
			commands: []any{
				mappings.NoteAdd,
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
				mappings.NoteAdd,
				mappings.AccentDecrease,
			},
			expectedAccent: 6,
			description:    "Should add note and decrease accent by 1 (index 5 -> 6)",
		},
		{
			name: "Add note and decrease accent twice",
			commands: []any{
				mappings.NoteAdd,
				mappings.AccentDecrease,
				mappings.AccentDecrease,
			},
			expectedAccent: 7,
			description:    "Should add note and decrease accent by 2 (index 5 -> 7)",
		},
		{
			name: "Add note and decrease accent at boundary",
			commands: []any{
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
		expectedGate int16
		description  string
	}{
		{
			name: "Add note and increase gate",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateIncrease,
			},
			expectedGate: 1,
			description:  "Should add note and increase gate by 1 (index 0 -> 1)",
		},
		{
			name: "Add note and increase gate twice",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateIncrease,
				mappings.GateIncrease,
			},
			expectedGate: 2,
			description:  "Should add note and increase gate by 2 (index 0 -> 2)",
		},
		{
			name: "Add note and increase gate at boundary",
			commands: []any{
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
				mappings.GateDecrease,
			},
			expectedGate: 0,
			description:  "Should add note and decrease gate stays at minimum index (0)",
		},
		{
			name: "Add note, increase then decrease gate",
			commands: []any{
				mappings.NoteAdd,
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
			config.LongGates = config.GetGateLengths(1) // Create a shorter boundary

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedGate, currentNote.GateIndex, tt.description+" - gate value")
		})
	}
}

func TestGateBigIncrease(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		expectedGate  int16
		description   string
		maxGateLength int
	}{
		{
			name: "Add note and big increase gate",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease,
			},
			expectedGate:  8,
			maxGateLength: 32,
			description:   "Should add note and big increase gate by 8",
		},
		{
			name: "Add note and big increase gate from non-zero",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateIncrease,
				mappings.GateBigIncrease,
			},
			expectedGate:  9,
			maxGateLength: 32,
			description:   "Should add note and big increase gate by 9",
		},
		{
			name: "Add note and big decrease gate",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease,
				mappings.GateBigDecrease,
			},
			expectedGate:  0,
			maxGateLength: 32,
			description:   "Should add note, big increase then big decrease gate (index 0 -> 7 -> 0)",
		},
		{
			name: "Add note and big increase gate at boundary",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease,
			},
			expectedGate:  7,
			maxGateLength: 1, // these work in multiples of 8
			description:   "Should add note, big increase then big decrease gate (index 0 -> 7 -> 0)",
		},
		{
			name: "Add note and big decrease gate from max",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease,
				mappings.GateBigDecrease,
			},
			expectedGate: 0,
			description:  "Should add note and big decrease gate by 8 (index 7 -> 0, capped at min)",
		},
		{
			name: "Add note and big decrease gate at boundary",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigDecrease,
			},
			expectedGate: 0,
			description:  "Should add note and big decrease gate stays at minimum index (0)",
		},
		{
			name: "Add note, big increase then regular decrease",
			commands: []any{
				mappings.NoteAdd,
				mappings.GateBigIncrease,
				mappings.GateDecrease,
			},
			expectedGate: 6,
			description:  "Should add note, big increase then decrease gate (index 0 -> 7 -> 6)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()
			config.LongGates = config.GetGateLengths(tt.maxGateLength)

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
				mappings.NoteAdd,
				mappings.WaitIncrease,
			},
			expectedWait: 1,
			description:  "Should add note and increase wait by 1 (index 0 -> 1)",
		},
		{
			name: "Add note and increase wait twice",
			commands: []any{
				mappings.NoteAdd,
				mappings.WaitIncrease,
				mappings.WaitIncrease,
			},
			expectedWait: 2,
			description:  "Should add note and increase wait by 2 (index 0 -> 2)",
		},
		{
			name: "Add note and increase wait at boundary",
			commands: []any{
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
				mappings.WaitDecrease,
			},
			expectedWait: 0,
			description:  "Should add note and decrease wait stays at minimum index (0)",
		},
		{
			name: "Add note, increase then decrease wait",
			commands: []any{
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
				mappings.RotateDown,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 4},
			expectedPos: grid.GridKey{Line: 3, Beat: 4},
			description: "Should rotate pattern down by one line",
		},
		{
			name: "Rotate down at boundary on next overlay",
			commands: []any{
				mappings.NoteAdd,
				mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter,
				mappings.RotateDown,
			},
			initialPos:  grid.GridKey{Line: 0, Beat: 0},
			expectedPos: grid.GridKey{Line: 1, Beat: 0},
			description: "Should rotate pattern down wrapping to top",
		},
		{
			name: "Rotate down at boundary",
			commands: []any{
				mappings.NoteAdd,
				mappings.RotateDown,
			},
			initialPos:  grid.GridKey{Line: 7, Beat: 4},
			expectedPos: grid.GridKey{Line: 0, Beat: 4},
			description: "Should rotate pattern down wrapping to top",
		},
		{
			name: "Rotate up",
			commands: []any{
				mappings.NoteAdd,
				mappings.RotateUp,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 4},
			expectedPos: grid.GridKey{Line: 1, Beat: 4},
			description: "Should rotate pattern up by one line",
		},
		{
			name: "Rotate up at boundary",
			commands: []any{
				mappings.NoteAdd,
				mappings.RotateUp,
			},
			initialPos:  grid.GridKey{Line: 0, Beat: 4},
			expectedPos: grid.GridKey{Line: 7, Beat: 4},
			description: "Should rotate pattern up wrapping to bottom",
		},
		{
			name: "Rotate up at boundary on next overlay",
			commands: []any{
				mappings.NoteAdd,
				mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter,
				mappings.RotateUp,
			},
			initialPos:  grid.GridKey{Line: 7, Beat: 0},
			expectedPos: grid.GridKey{Line: 6, Beat: 0},
			description: "Should rotate pattern up wrapping to bottom",
		},
		{
			name: "Rotate up at visual boundary",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleVisualLineMode,
				mappings.CursorDown,
				mappings.RotateUp,
			},
			initialPos:  grid.GridKey{Line: 0, Beat: 4},
			expectedPos: grid.GridKey{Line: 1, Beat: 4},
			description: "Should rotate pattern up wrapping to bottom",
		},
		{
			name: "Rotate down at visual boundary",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleVisualLineMode,
				mappings.CursorDown,
				mappings.RotateDown,
			},
			initialPos:  grid.GridKey{Line: 0, Beat: 4},
			expectedPos: grid.GridKey{Line: 1, Beat: 4},
			description: "Should rotate pattern up wrapping to bottom",
		},
		{
			name: "Rotate right",
			commands: []any{
				mappings.NoteAdd,
				mappings.RotateRight,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 4},
			expectedPos: grid.GridKey{Line: 2, Beat: 5},
			description: "Should rotate pattern right by one beat",
		},
		{
			name: "Rotate right at boundary",
			commands: []any{
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
				mappings.RotateLeft,
			},
			initialPos:  grid.GridKey{Line: 2, Beat: 0},
			expectedPos: grid.GridKey{Line: 2, Beat: 31},
			description: "Should rotate pattern left wrapping to end",
		},
		{
			name: "Multiple rotations",
			commands: []any{
				mappings.NoteAdd,
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
			m := createTestModel(WithGridCursor(tt.initialPos))

			m, _ = processCommand(mappings.NoteAdd, m)
			initialNote, exists := m.currentOverlay.GetNote(tt.initialPos)
			assert.True(t, exists, tt.description+" - note should exist at initial position")

			rotateCommands := tt.commands[1:]
			m, _ = processCommands(rotateCommands, m)

			rotatedNote, exists := m.currentOverlay.GetNote(tt.expectedPos)
			assert.True(t, exists, tt.description+" - note should exist at expected position")
			assert.Equal(t, initialNote, rotatedNote, tt.description+" - note should be the same")

			if tt.initialPos != tt.expectedPos {
				note, stillExists := m.currentOverlay.GetNote(tt.initialPos)
				assert.True(t, !stillExists || note == zeronote, tt.description+" - note should not exist at initial position")
			}
		})
	}
}

func TestNoteAddRemoveAndOverlayRemoveOnRootOverlay(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		initialPos     grid.GridKey
		shouldHaveNote bool
		description    string
	}{
		{
			name:           "NoteAdd creates note at cursor position",
			commands:       []any{mappings.NoteAdd},
			initialPos:     grid.GridKey{Line: 0, Beat: 4},
			shouldHaveNote: true,
			description:    "Should add a note at the cursor position",
		},
		{
			name:           "NoteRemove removes note at cursor position",
			commands:       []any{mappings.NoteAdd, mappings.NoteRemove},
			initialPos:     grid.GridKey{Line: 1, Beat: 8},
			shouldHaveNote: false,
			description:    "Should remove the note at the cursor position",
		},
		{
			name:           "OverlayNoteRemove removes note from overlay",
			commands:       []any{mappings.NoteAdd, mappings.OverlayNoteRemove},
			initialPos:     grid.GridKey{Line: 2, Beat: 12},
			shouldHaveNote: false,
			description:    "Should remove the note from the overlay",
		},
		{
			name:           "NoteAdd on existing note position",
			commands:       []any{mappings.NoteAdd, mappings.NoteAdd},
			initialPos:     grid.GridKey{Line: 3, Beat: 16},
			shouldHaveNote: true,
			description:    "Should still have a note after adding to existing position",
		},
		{
			name:           "NoteRemove on empty position does nothing",
			commands:       []any{mappings.NoteRemove},
			initialPos:     grid.GridKey{Line: 4, Beat: 20},
			shouldHaveNote: false,
			description:    "Should remain empty after removing from empty position",
		},
		{
			name:           "OverlayNoteRemove on empty position does nothing",
			commands:       []any{mappings.OverlayNoteRemove},
			initialPos:     grid.GridKey{Line: 5, Beat: 24},
			shouldHaveNote: false,
			description:    "Should remain empty after overlay removing from empty position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.initialPos))

			initialNote, initialExists := m.CurrentNote()
			assert.False(t, initialExists || initialNote != zeronote, "Should start with no note at position")

			m, _ = processCommands(tt.commands, m)

			_, exists := m.CurrentNote()

			assert.Equal(t, tt.shouldHaveNote, exists, tt.description)

		})
	}
}

func TestNoteAddRemoveAndOverlayRemoveOnNonRootOverlay(t *testing.T) {
	nonRootOverlayKey := overlaykey.OverlayPeriodicity{
		Shift:      2,
		Interval:   4,
		Width:      0,
		StartCycle: 0,
	}

	tests := []struct {
		name           string
		commands       []any
		initialPos     grid.GridKey
		shouldHaveNote bool
		description    string
	}{
		{
			name:           "NoteAdd creates note at cursor position",
			commands:       []any{mappings.NoteAdd},
			initialPos:     grid.GridKey{Line: 0, Beat: 4},
			shouldHaveNote: true,
			description:    "Should add a note at the cursor position",
		},
		{
			name:           "NoteRemove removes note at cursor position",
			commands:       []any{mappings.NoteAdd, mappings.NoteRemove},
			initialPos:     grid.GridKey{Line: 1, Beat: 8},
			shouldHaveNote: false,
			description:    "Should remove the note at the cursor position",
		},
		{
			name:           "OverlayNoteRemove removes note from overlay",
			commands:       []any{mappings.NoteAdd, mappings.OverlayNoteRemove},
			initialPos:     grid.GridKey{Line: 2, Beat: 12},
			shouldHaveNote: false,
			description:    "Should remove the note from the overlay",
		},
		{
			name:           "NoteAdd on existing note position",
			commands:       []any{mappings.NoteAdd, mappings.NoteAdd},
			initialPos:     grid.GridKey{Line: 3, Beat: 16},
			shouldHaveNote: true,
			description:    "Should still have a note after adding to existing position",
		},
		{
			name:           "NoteRemove on empty position does nothing",
			commands:       []any{mappings.NoteRemove},
			initialPos:     grid.GridKey{Line: 4, Beat: 20},
			shouldHaveNote: false,
			description:    "Should remain empty after removing from empty position",
		},
		{
			name:           "OverlayNoteRemove on empty position does nothing",
			commands:       []any{mappings.OverlayNoteRemove},
			initialPos:     grid.GridKey{Line: 5, Beat: 24},
			shouldHaveNote: false,
			description:    "Should remain empty after overlay removing from empty position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.initialPos), WithNonRootOverlay(nonRootOverlayKey))

			initialNote, initialExists := m.CurrentNote()
			assert.False(t, initialExists || initialNote != zeronote, "Should start with no note at position")

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			hasNote := exists && currentNote != zeronote

			assert.Equal(t, tt.shouldHaveNote, hasNote, tt.description)

			if tt.shouldHaveNote {
				assert.NotEqual(t, zeronote, currentNote, "Note should not be zeronote when expected to have note")
				assert.True(t, exists, "Note should exist when expected")
			} else {
				assert.True(t, !exists || currentNote == zeronote, "Should not have note when not expected")
			}
		})
	}
}

func TestNoteAddRemoveAndOverlayRemoveAcrossOverlays(t *testing.T) {
	tests := []struct {
		name               string
		commands           []any
		initialPos         grid.GridKey
		shouldHaveNote     bool
		shouldHaveZeroNote bool
		description        string
	}{
		{
			name:               "NoteAdd creates note at cursor position when note exists on root overlay",
			commands:           []any{mappings.NoteAdd, mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter, mappings.NoteAdd},
			initialPos:         grid.GridKey{Line: 0, Beat: 4},
			shouldHaveNote:     true,
			shouldHaveZeroNote: false,
			description:        "Should add a note at the cursor position",
		},
		{
			name:               "NoteRemove does not zero note at cursor position when note does not exist on root overlay",
			commands:           []any{mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter, mappings.NoteRemove},
			initialPos:         grid.GridKey{Line: 1, Beat: 8},
			shouldHaveNote:     false,
			shouldHaveZeroNote: false,
			description:        "Should remove the note at the cursor position",
		},
		{
			name:               "NoteRemove zeroes note at cursor position when note exists on root overlay",
			commands:           []any{mappings.NoteAdd, mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter, mappings.NoteRemove},
			initialPos:         grid.GridKey{Line: 1, Beat: 8},
			shouldHaveNote:     true,
			shouldHaveZeroNote: true,
			description:        "Should remove the note at the cursor position",
		},
		{
			name:               "OverlayNoteRemove does nothing when note exists on root overlay but not on above overlay",
			commands:           []any{mappings.NoteAdd, mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter, mappings.OverlayNoteRemove},
			initialPos:         grid.GridKey{Line: 2, Beat: 12},
			shouldHaveNote:     true,
			shouldHaveZeroNote: false,
			description:        "Should remove the note from the overlay",
		},
		{
			name:               "OverlayNoteRemove removes note from overlay when nothing exists below",
			commands:           []any{mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter, mappings.NoteAdd, mappings.OverlayNoteRemove},
			initialPos:         grid.GridKey{Line: 2, Beat: 12},
			shouldHaveNote:     false,
			shouldHaveZeroNote: false,
			description:        "Should remove the note from the overlay",
		},
		{
			name:               "NoteRemove on empty position on root overlay does nothing",
			commands:           []any{mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter, mappings.NoteRemove},
			initialPos:         grid.GridKey{Line: 4, Beat: 20},
			shouldHaveNote:     false,
			shouldHaveZeroNote: false,
			description:        "Should remain empty after removing from empty position",
		},
		{
			name:               "OverlayNoteRemove on empty position does nothing",
			commands:           []any{mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter, mappings.OverlayNoteRemove},
			initialPos:         grid.GridKey{Line: 5, Beat: 24},
			shouldHaveNote:     false,
			shouldHaveZeroNote: false,
			description:        "Should remain empty after overlay removing from empty position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.initialPos))

			initialNote, initialExists := m.CurrentNote()
			assert.False(t, initialExists || initialNote != zeronote, "Should start with no note at position")

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()

			assert.Equal(t, tt.shouldHaveNote, exists, tt.description)

			if tt.shouldHaveNote {
				assert.Equal(t, tt.shouldHaveZeroNote, currentNote == zeronote, "Note should be zeronote when expected to have zeronote")
			}
		})
	}
}

func TestRatchetInputSwitch(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedSelection operation.Selection
		description       string
	}{
		{
			name:              "RatchetInputSwitch when cursor note no note",
			commands:          []any{mappings.RatchetInputSwitch},
			expectedSelection: operation.SelectGrid,
			description:       "First ratchet input switch does nothing if not on a note",
		},
		{
			name:              "RatchetInputSwitch when cursor on note",
			commands:          []any{mappings.NoteAdd, mappings.RatchetInputSwitch},
			expectedSelection: operation.SelectRatchets,
			description:       "First ratchet input switch should select ratchet",
		},
		{
			name:              "Second RatchetInputSwitch when cursor on note",
			commands:          []any{mappings.NoteAdd, mappings.RatchetInputSwitch, mappings.RatchetInputSwitch},
			expectedSelection: operation.SelectRatchetSpan,
			description:       "Second ratchet input switch should select ratchet span",
		},
		{
			name:              "Third RatchetInputSwitch",
			commands:          []any{mappings.NoteAdd, mappings.RatchetInputSwitch, mappings.RatchetInputSwitch, mappings.RatchetInputSwitch},
			expectedSelection: operation.SelectGrid,
			description:       "Second ratchet input switch should select ratchet span",
		},
		{
			name:              "Escape Ratchet Input",
			commands:          []any{mappings.NoteAdd, mappings.RatchetInputSwitch, mappings.Escape},
			expectedSelection: operation.SelectGrid,
			description:       "Second ratchet input switch should select ratchet span",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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
				mappings.NoteAdd,
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

func TestSpecificValueActionAndCursorMovement(t *testing.T) {
	tests := []struct {
		name              string
		setupCommands     []any
		testCommands      []any
		initialPos        grid.GridKey
		moveToPos         grid.GridKey
		expectedSelection operation.Selection
		expectedValue     uint8
		description       string
	}{
		{
			name: "Add specific value action and test cursor movement sets selection indicator",
			setupCommands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase, // Set to ProgramChange
				mappings.Escape,
			},
			testCommands: []any{
				mappings.ActionAddSpecificValue,
			},
			initialPos:        grid.GridKey{Line: 0, Beat: 4},
			expectedSelection: operation.SelectSpecificValue,
			expectedValue:     0,
			description:       "Should add specific value action and set selection indicator when cursor is on note",
		},
		{
			name: "Increase specific value note",
			setupCommands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
			},
			testCommands: []any{
				mappings.ActionAddSpecificValue,
				mappings.Increase,
			},
			initialPos:        grid.GridKey{Line: 0, Beat: 4},
			expectedSelection: operation.SelectSpecificValue,
			expectedValue:     1,
			description:       "Should increase specific value by 1",
		},
		{
			name: "Increase then decrease specific value note",
			setupCommands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
			},
			testCommands: []any{
				mappings.ActionAddSpecificValue,
				mappings.Increase,
				mappings.Increase,
				mappings.Decrease,
			},
			initialPos:        grid.GridKey{Line: 0, Beat: 4},
			expectedSelection: operation.SelectSpecificValue,
			expectedValue:     1,
			description:       "Should increase specific value by 2 then decrease by 1",
		},
		{
			name: "Move cursor away from specific value note resets selection indicator",
			setupCommands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
			},
			testCommands: []any{
				mappings.ActionAddSpecificValue,
				mappings.CursorRight,
			},
			initialPos:        grid.GridKey{Line: 0, Beat: 4},
			moveToPos:         grid.GridKey{Line: 0, Beat: 5},
			expectedSelection: operation.SelectGrid,
			expectedValue:     0,
			description:       "Should reset selection indicator when cursor moves away from specific value note",
		},
		{
			name: "Move cursor back to specific value note sets selection indicator",
			setupCommands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
			},
			testCommands: []any{
				mappings.ActionAddSpecificValue,
				mappings.Increase,
				mappings.CursorRight,
				mappings.CursorLeft,
			},
			initialPos:        grid.GridKey{Line: 0, Beat: 4},
			expectedSelection: operation.SelectSpecificValue,
			expectedValue:     1,
			description:       "Should set selection indicator when cursor moves back to specific value note",
		},
		{
			name: "Multiple specific value notes with cursor movement",
			setupCommands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
			},
			testCommands: []any{
				mappings.ActionAddSpecificValue,
				mappings.Increase,
				mappings.Increase,
				mappings.CursorRight,
				mappings.ActionAddSpecificValue,
				mappings.Increase,
			},
			initialPos:        grid.GridKey{Line: 0, Beat: 4},
			moveToPos:         grid.GridKey{Line: 0, Beat: 5},
			expectedSelection: operation.SelectSpecificValue,
			expectedValue:     1,
			description:       "Should handle multiple specific value notes correctly",
		},
		{
			name: "undo specific value",
			setupCommands: []any{
				mappings.SetupInputSwitch,
				mappings.SetupInputSwitch,
				mappings.Increase,
				mappings.Escape,
			},
			testCommands: []any{
				mappings.ActionAddSpecificValue,
				mappings.Increase,
				mappings.Increase,
				mappings.Enter,
				mappings.Undo,
			},
			initialPos:        grid.GridKey{Line: 0, Beat: 0},
			expectedSelection: operation.SelectSpecificValue,
			expectedValue:     0,
			description:       "Should handle multiple specific value notes correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(WithGridCursor(tt.initialPos))

			// Run setup commands to configure the line
			m, _ = processCommands(tt.setupCommands, m)

			// Verify initial state
			assert.Equal(t, operation.SelectGrid, m.selectionIndicator, "Initial selection should be nothing")

			// Run test commands
			m, _ = processCommands(tt.testCommands, m)

			// Check final state
			assert.Equal(t, tt.expectedSelection, m.selectionIndicator, tt.description+" - selection indicator")

			// If we moved cursor, verify we're at the expected position
			if tt.moveToPos != (grid.GridKey{}) {
				assert.Equal(t, tt.moveToPos, m.gridCursor, tt.description+" - cursor position")
			}

			// Verify the note exists and has correct properties
			currentNote, exists := m.CurrentNote()
			if tt.expectedSelection == operation.SelectSpecificValue {
				assert.True(t, exists, tt.description+" - note should exist")
				assert.Equal(t, grid.ActionSpecificValue, currentNote.Action, tt.description+" - action type")
				assert.Equal(t, tt.expectedValue, currentNote.AccentIndex, tt.description+" - specific value")
			}

			// Test that the first note still exists with correct value if we created multiple notes
			if tt.moveToPos != (grid.GridKey{}) && tt.moveToPos != tt.initialPos {
				firstNote, firstExists := m.currentOverlay.GetNote(tt.initialPos)
				if tt.name == "Multiple specific value notes with cursor movement" {
					assert.True(t, firstExists, tt.description+" - first note should still exist")
					assert.Equal(t, grid.ActionSpecificValue, firstNote.Action, tt.description+" - first note action type")
					assert.Equal(t, uint8(2), firstNote.AccentIndex, tt.description+" - first note specific value should be 2")
				}
			}
		})
	}
}
