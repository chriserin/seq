package main

import (
	"slices"
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
			m := createTestModel(t.Context())

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
			m := createTestModel(t.Context())

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
			m := createTestModel(t.Context())

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
			m := createTestModel(t.Context())

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
			m := createTestModel(t.Context(), func(m *model) model {
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
			m := createTestModel(t.Context(), func(m *model) model {
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
			m := createTestModel(t.Context(), func(m *model) model {
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
			m := createTestModel(t.Context(), func(m *model) model {
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

func TestPatternModeMonoSpaceIncrease(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		expectedNotes []grid.GridKey
		description   string
	}{
		{
			name: "Add note, switch to mono mode, fill with 1",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleMonoMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedNotes: []grid.GridKey{GK(0, 1), GK(0, 2)},
			description:   "Should add note, switch to mono mode, and add note with every space pattern",
		},
		{
			name: "Add note, switch to mono mode, move to right no change",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleMonoMode,
				mappings.CursorRight,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(0, 1), GK(0, 2)},
			description:   "Should add note, switch to mono mode, move cursor right and note should remain",
		},
		{
			name: "Add notes in mono mode",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(0, 1), GK(0, 2)},
			description:   "Should add note, switch to mono mode, and handle increase then decrease",
		},
		{
			name: "Add notes on two lines in mono mode",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
				mappings.CursorDown,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedNotes: []grid.GridKey{GK(1, 0), GK(0, 1), GK(1, 2)},
			description:   "Should add note, switch to mono mode, and handle increase then decrease",
		},
		{
			name: "Add notes on two lines in mono mode with a 1 and 3",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
				mappings.CursorDown,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "3"},
			},
			expectedNotes: []grid.GridKey{GK(1, 0), GK(0, 1), GK(0, 2)},
			description:   "Should add note, switch to mono mode, and handle increase then decrease",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(t.Context(), WithGridSize(3, 2))

			m, _ = processCommands(tt.commands, m)

			for b := range 3 {
				for l := range 2 {
					gk := grid.GK(uint8(l), uint8(b))
					shouldExist := slices.Contains(tt.expectedNotes, gk)
					_, noteExists := m.currentOverlay.GetNote(gk)
					assert.Equal(t, shouldExist, noteExists, tt.description+" - note existence at grid location "+gk.String())
				}
			}
		})
	}
}

func TestPatternModeMonoAccentIncreasePattern(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedAccents []uint8
		description     string
	}{
		{
			name: "Switch to mono mode, increase accent every other space pattern",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.ToggleAccentMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedAccents: []uint8{4, 5, 4, 5},
			description:     "Should switch to mono mode and apply every other space pattern",
		},
		{
			name: "Switch to mono mode, decrease accent every 2nd note",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.ToggleAccentMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedAccents: []uint8{6, 5, 6, 5},
			description:     "Should switch to mono mode and decrease accent every 2nd note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(t.Context(), func(m *model) model {
				m.currentOverlay.AddNote(grid.GK(0, 0), note{AccentIndex: 5})
				m.currentOverlay.AddNote(grid.GK(0, 1), note{AccentIndex: 5})
				m.currentOverlay.AddNote(grid.GK(1, 2), note{AccentIndex: 5})
				m.currentOverlay.AddNote(grid.GK(1, 3), note{AccentIndex: 5})
				return *m
			}, WithGridSize(4, 2))

			m, _ = processCommands(tt.commands, m)

			gridLocations := []grid.GridKey{
				grid.GK(0, 0),
				grid.GK(0, 1),
				grid.GK(1, 2),
				grid.GK(1, 3),
			}

			for i, gk := range gridLocations {
				n, exists := m.currentOverlay.GetNote(gk)
				assert.True(t, exists, tt.description+" - note should exist at grid location "+gk.String())
				assert.Equal(t, tt.expectedAccents[i], n.AccentIndex, tt.description+" - note has expected accent "+gk.String())
			}
		})
	}
}

func TestPatternModeMonoGatePattern(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		expectedGates []int16
		description   string
	}{
		{
			name: "Switch to mono mode, increase accent every other space pattern",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedGates: []int16{6, 5, 6, 5},
			description:   "Should switch to mono mode and apply every other space pattern",
		},
		{
			name: "Switch to mono mode, decrease accent every 2nd note",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.ToggleGateMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedGates: []int16{4, 5, 4, 5},
			description:   "Should switch to mono mode and decrease accent every 2nd note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(t.Context(), func(m *model) model {
				m.currentOverlay.AddNote(grid.GK(0, 0), note{AccentIndex: 5, GateIndex: 5})
				m.currentOverlay.AddNote(grid.GK(0, 1), note{AccentIndex: 5, GateIndex: 5})
				m.currentOverlay.AddNote(grid.GK(1, 2), note{AccentIndex: 5, GateIndex: 5})
				m.currentOverlay.AddNote(grid.GK(1, 3), note{AccentIndex: 5, GateIndex: 5})
				return *m
			}, WithGridSize(4, 2))

			m, _ = processCommands(tt.commands, m)

			gridLocations := []grid.GridKey{
				grid.GK(0, 0),
				grid.GK(0, 1),
				grid.GK(1, 2),
				grid.GK(1, 3),
			}

			for i, gk := range gridLocations {
				n, exists := m.currentOverlay.GetNote(gk)
				assert.True(t, exists, tt.description+" - note should exist at grid location "+gk.String())
				assert.Equal(t, tt.expectedGates[i], n.GateIndex, tt.description+" - note has expected accent "+gk.String())
			}
		})
	}
}

func TestPatternNoteModeMonoAccentPattern(t *testing.T) {
	tests := []struct {
		name            string
		commands        []any
		expectedAccents []uint8
		description     string
	}{
		{
			name: "Switch to mono mode, increase accent every other space pattern",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.ToggleAccentNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedAccents: []uint8{4, 5, 4, 5},
			description:     "Should switch to mono mode and apply every other space pattern",
		},
		{
			name: "Switch to mono mode, decrease accent every 2nd note",
			commands: []any{
				mappings.ToggleMonoMode,
				mappings.ToggleAccentNoteMode,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "2"},
			},
			expectedAccents: []uint8{6, 5, 6, 5},
			description:     "Should switch to mono mode and decrease accent every 2nd note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(t.Context(), func(m *model) model {
				m.currentOverlay.AddNote(grid.GK(0, 0), note{AccentIndex: 5})
				m.currentOverlay.AddNote(grid.GK(1, 2), note{AccentIndex: 5})
				m.currentOverlay.AddNote(grid.GK(1, 3), note{AccentIndex: 5})
				m.currentOverlay.AddNote(grid.GK(1, 8), note{AccentIndex: 5})
				return *m
			}, WithGridSize(4, 2))

			m, _ = processCommands(tt.commands, m)

			gridLocations := []grid.GridKey{
				grid.GK(0, 0),
				grid.GK(1, 2),
				grid.GK(1, 3),
				grid.GK(1, 8),
			}

			for i, gk := range gridLocations {
				n, exists := m.currentOverlay.GetNote(gk)
				assert.True(t, exists, tt.description+" - note should exist at grid location "+gk.String())
				assert.Equal(t, tt.expectedAccents[i], n.AccentIndex, tt.description+" - note has expected accent "+gk.String())
			}
		})
	}
}

func TestPatternModeFillSpace(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		expectedNotes []grid.GridKey
		description   string
	}{
		{
			name: "Space pattern fill",
			commands: []any{
				mappings.NoteAdd,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(0, 1), GK(0, 2)},
			description:   "Add note, then fill with space pattern",
		},
		{
			name: "Space pattern fill, with 2",
			commands: []any{
				mappings.NoteAdd,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(0, 1)},
			description:   "Add note, then fill with space pattern 2",
		},
		{
			name: "Space pattern fill doesn't overwrite existing notes",
			commands: []any{
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(0, 1), GK(0, 2)},
			description:   "Add note, then fill with space pattern 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(t.Context(), WithGridSize(3, 2))

			m, _ = processCommands(tt.commands, m)

			for b := range 3 {
				for l := range 2 {
					gk := grid.GK(uint8(l), uint8(b))
					shouldExist := slices.Contains(tt.expectedNotes, gk)
					_, noteExists := m.currentOverlay.GetNote(gk)
					assert.Equal(t, shouldExist, noteExists, tt.description+" - note existence at grid location "+gk.String())
				}
			}
		})
	}
}

func TestPatternModeMonoFillSpace(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		expectedNotes []grid.GridKey
		description   string
	}{
		{
			name: "Space pattern fill",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleMonoMode,
				mappings.CursorDown,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(1, 1), GK(1, 2)},
			description:   "Add note, then fill with space pattern",
		},
		{
			name: "Space pattern fill, with 2",
			commands: []any{
				mappings.NoteAdd,
				mappings.ToggleMonoMode,
				mappings.CursorDown,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "@"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(1, 1)},
			description:   "Add note, then fill with space pattern 2",
		},
		{
			name: "Space pattern fill doesn't overwrite existing notes",
			commands: []any{
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "1"},
				mappings.ToggleMonoMode,
				mappings.CursorDown,
				mappings.Mapping{Command: mappings.NumberPattern, LastValue: "!"},
			},
			expectedNotes: []grid.GridKey{GK(0, 0), GK(0, 1), GK(0, 2)},
			description:   "Add note, then fill with space pattern 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(t.Context(), WithGridSize(3, 2))

			m, _ = processCommands(tt.commands, m)

			for b := range 3 {
				for l := range 2 {
					gk := grid.GK(uint8(l), uint8(b))
					shouldExist := slices.Contains(tt.expectedNotes, gk)
					_, noteExists := m.currentOverlay.GetNote(gk)
					assert.Equal(t, shouldExist, noteExists, tt.description+" - note existence at grid location "+gk.String())
				}
			}
		})
	}
}
