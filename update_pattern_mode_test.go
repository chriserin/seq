package main

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/stretchr/testify/assert"
)

func TestPatternModeAccentIncrease(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		expectedAccent uint8
		description    string
	}{
		{
			name: "Add note, switch to accent mode, increase accent by 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleAccentMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedAccent: 4,
			description:    "Should add note, switch to accent mode, and increase accent by 1",
		},
		{
			name: "Add note, switch to accent mode, move to right no change",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleAccentMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedAccent: 5,
			description:    "Should add note, switch to accent mode, move cursor right and accent should remain 0",
		},
		{
			name: "Add note, switch to accent mode, increase accent by 1, decrease accent by 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleAccentMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedAccent: 6,
			description:    "Should add note, switch to accent mode, and decrease accent by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			// NOTE: Initial value for accent is 5
			// and "increasing" accent index means moving to 4
			// as accents are defined backwards
			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedAccent, currentNote.AccentIndex, tt.description+" - accent value")
		})
	}
}

func TestPatternModeRatchetIncrease(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedRatchet grid.Ratchet
		description     string
	}{
		{
			name: "Add note, switch to ratchet mode, increase ratchet by 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleRatchetMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   3, // hits are 0b11, so 3
				Length: 1,
			},
			description: "Should add note, switch to ratchet mode, and increase ratchet by 1",
		},
		{
			name: "Add note, switch to ratchet mode, move to right no change",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleRatchetMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1,
				Length: 0,
			},
			description: "Should add note, switch to ratchet mode, move cursor right and ratchet should remain default",
		},
		{
			name: "Add note, switch to ratchet mode, increase then decrease ratchet",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleRatchetMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedRatchet: grid.Ratchet{
				Span:   0,
				Hits:   1,
				Length: 0,
			},
			description: "Should add note, switch to ratchet mode, and set ratchet to 2",
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

func TestPatternModeWaitIncrease(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		expectedWait uint8
		description  string
	}{
		{
			name: "Add note, switch to wait mode, increase wait by 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleWaitMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedWait: 1,
			description:  "Should add note, switch to wait mode, and increase wait by 1",
		},
		{
			name: "Add note, switch to wait mode, move to right no change",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleWaitMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedWait: 0,
			description:  "Should add note, switch to wait mode, move cursor right and wait should remain 0",
		},
		{
			name: "Add note, switch to wait mode, increase wait by 1, decrease wait by 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleWaitMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedWait: 0,
			description:  "Should add note, switch to wait mode, and decrease wait by 1",
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
