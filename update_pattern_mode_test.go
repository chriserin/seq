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

func TestPatternModeGateIncrease(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		expectedGate int16
		description  string
	}{
		{
			name: "Add note, switch to gate mode, increase gate by 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedGate: 1,
			description:  "Should add note, switch to gate mode, and increase gate by 1",
		},
		{
			name: "Add note, switch to gate mode, move to right no change",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleGateMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.CursorLeft, // move back to the note
			},
			expectedGate: 0,
			description:  "Should add note, switch to gate mode, move cursor right and gate should remain 0",
		},
		{
			name: "Add note, switch to gate mode, increase gate by 1, decrease gate by 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedGate: 0,
			description:  "Should add note, switch to gate mode, increase then decrease gate by 1",
		},
		{
			name: "Add note, switch to gate mode, increase by 1 every other beat",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedGate: 1,
			description:  "Should add note, switch to gate mode, and increase every other gate by 1",
		},
		{
			name: "Add note, switch to gate mode, decrease gate by 2",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedGate: 0, // gate can't go below 0
			description:  "Should add note, switch to gate mode, and attempt to decrease gate by 2 (clamped to 0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			currentNote, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist")
			assert.Equal(t, tt.expectedGate, currentNote.GateIndex, tt.description+" - gate value")
		})
	}
}

func TestPatternModeNoteAccentIncrease(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedAccents []uint8
		description     string
	}{
		{
			name: "Add note, switch to note accent mode, increase accent by 1",
			commands: []any{
				mappings.ToggleAccentNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedAccents: []uint8{2, 3, 2},
			description:     "Should add note, switch to note accent mode, and increase accent by 1",
		},
		{
			name: "Add note, switch to note accent mode, decrease accent by 1",
			commands: []any{
				mappings.ToggleAccentNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedAccents: []uint8{4, 3, 4},
			description:     "Should add note, switch to note accent mode, and increase accent by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.currentOverlay.AddNote(grid.GK(0, 0), note{AccentIndex: 3})
				m.currentOverlay.AddNote(grid.GK(0, 2), note{AccentIndex: 3})
				m.currentOverlay.AddNote(grid.GK(0, 3), note{AccentIndex: 3})
				return *m
			})

			m, _ = processCommands(tt.commands, m)

			gridLocations := []grid.GridKey{
				grid.GK(0, 0),
				grid.GK(0, 2),
				grid.GK(0, 3),
			}

			for i, gk := range gridLocations {
				note, exists := m.currentOverlay.GetNote(gk)
				assert.True(t, exists, tt.description+" - note should exist at grid location "+gk.String())
				assert.Equal(t, tt.expectedAccents[i], note.AccentIndex, tt.description+" - accent value at grid location "+gk.String())
			}
		})
	}
}

func TestPatternModeNoteWaitIncrease(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		expectedWaits []uint8
		description   string
	}{
		{
			name: "Add note, switch to note accent mode, increase accent by 1",
			commands: []any{
				mappings.ToggleWaitNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedWaits: []uint8{4, 3, 4},
			description:   "Should add note, switch to note accent mode, and increase accent by 1",
		},
		{
			name: "Add note, switch to note accent mode, decrease accent by 1",
			commands: []any{
				mappings.ToggleWaitNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedWaits: []uint8{2, 3, 2},
			description:   "Should add note, switch to note accent mode, and increase accent by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.currentOverlay.AddNote(grid.GK(0, 0), note{WaitIndex: 3})
				m.currentOverlay.AddNote(grid.GK(0, 2), note{WaitIndex: 3})
				m.currentOverlay.AddNote(grid.GK(0, 3), note{WaitIndex: 3})
				return *m
			})

			m, _ = processCommands(tt.commands, m)

			gridLocations := []grid.GridKey{
				grid.GK(0, 0),
				grid.GK(0, 2),
				grid.GK(0, 3),
			}

			for i, gk := range gridLocations {
				note, exists := m.currentOverlay.GetNote(gk)
				assert.True(t, exists, tt.description+" - note should exist at grid location "+gk.String())
				assert.Equal(t, tt.expectedWaits[i], note.WaitIndex, tt.description+" - accent value at grid location "+gk.String())
			}
		})
	}
}

func TestPatternModeNoteGateIncrease(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		expectedGates []int16
		description   string
	}{
		{
			name: "Add note, switch to note accent mode, increase accent by 1",
			commands: []any{
				mappings.ToggleGateNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedGates: []int16{4, 3, 4},
			description:   "Should add note, switch to note accent mode, and increase accent by 1",
		},
		{
			name: "Add note, switch to note accent mode, decrease accent by 1",
			commands: []any{
				mappings.ToggleGateNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedGates: []int16{2, 3, 2},
			description:   "Should add note, switch to note accent mode, and increase accent by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.currentOverlay.AddNote(grid.GK(0, 0), note{GateIndex: 3})
				m.currentOverlay.AddNote(grid.GK(0, 2), note{GateIndex: 3})
				m.currentOverlay.AddNote(grid.GK(0, 3), note{GateIndex: 3})
				return *m
			})

			m, _ = processCommands(tt.commands, m)

			gridLocations := []grid.GridKey{
				grid.GK(0, 0),
				grid.GK(0, 2),
				grid.GK(0, 3),
			}

			for i, gk := range gridLocations {
				note, exists := m.currentOverlay.GetNote(gk)
				assert.True(t, exists, tt.description+" - note should exist at grid location "+gk.String())
				assert.Equal(t, tt.expectedGates[i], note.GateIndex, tt.description+" - accent value at grid location "+gk.String())
			}
		})
	}
}

func TestPatternModeNoteRatchetIncrease(t *testing.T) {
	tests := []struct {
		name                   string
		commands               []any
		expectedRatchetLengths []uint8
		description            string
	}{
		{
			name: "Add note, switch to note accent mode, increase accent by 1",
			commands: []any{
				mappings.ToggleRatchetNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedRatchetLengths: []uint8{1, 0, 1},
			description:            "Should add note, switch to note accent mode, and increase accent by 1",
		},
		{
			name: "Add note, switch to note accent mode, decrease accent by 1",
			commands: []any{
				mappings.ToggleRatchetNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedRatchetLengths: []uint8{1, 0, 1},
			description:            "Should add note, switch to note accent mode, and increase accent by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(func(m *model) model {
				m.currentOverlay.AddNote(grid.GK(0, 0), note{AccentIndex: 1})
				m.currentOverlay.AddNote(grid.GK(0, 2), note{AccentIndex: 1})
				m.currentOverlay.AddNote(grid.GK(0, 3), note{AccentIndex: 1})
				return *m
			})

			m, _ = processCommands(tt.commands, m)

			gridLocations := []grid.GridKey{
				grid.GK(0, 0),
				grid.GK(0, 2),
				grid.GK(0, 3),
			}

			for i, gk := range gridLocations {
				note, exists := m.currentOverlay.GetNote(gk)
				assert.True(t, exists, tt.description+" - note should exist at grid location "+gk.String())
				assert.Equal(t, tt.expectedRatchetLengths[i], note.Ratchets.Length, tt.description+" - accent value at grid location "+gk.String())
			}
		})
	}
}
