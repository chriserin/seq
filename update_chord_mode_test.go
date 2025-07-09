package main

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/stretchr/testify/assert"
)

func TestToggleChordMode(t *testing.T) {
	tests := []struct {
		name                  string
		commands              []any
		initialSequencerType  grid.SequencerType
		expectedSequencerType grid.SequencerType
		description           string
	}{
		{
			name:                  "Toggle from trigger to polyphony",
			commands:              []any{mappings.ToggleChordMode},
			initialSequencerType:  grid.SeqtypeTrigger,
			expectedSequencerType: grid.SeqtypePolyphony,
			description:           "Should toggle from trigger mode to polyphony (chord) mode",
		},
		{
			name:                  "Toggle from polyphony to trigger",
			commands:              []any{mappings.ToggleChordMode},
			initialSequencerType:  grid.SeqtypePolyphony,
			expectedSequencerType: grid.SeqtypeTrigger,
			description:           "Should toggle from polyphony (chord) mode to trigger mode",
		},
		{
			name:                  "Multiple toggles return to original state",
			commands:              []any{mappings.ToggleChordMode, mappings.ToggleChordMode},
			initialSequencerType:  grid.SeqtypeTrigger,
			expectedSequencerType: grid.SeqtypeTrigger,
			description:           "Should return to original state after two toggles",
		},
		{
			name:                  "Multiple toggles from polyphony return to original state",
			commands:              []any{mappings.ToggleChordMode, mappings.ToggleChordMode},
			initialSequencerType:  grid.SeqtypePolyphony,
			expectedSequencerType: grid.SeqtypePolyphony,
			description:           "Should return to original polyphony state after two toggles",
		},
		{
			name:                  "Three toggles from trigger ends in polyphony",
			commands:              []any{mappings.ToggleChordMode, mappings.ToggleChordMode, mappings.ToggleChordMode},
			initialSequencerType:  grid.SeqtypeTrigger,
			expectedSequencerType: grid.SeqtypePolyphony,
			description:           "Should end in polyphony mode after three toggles from trigger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.templateSequencerType = tt.initialSequencerType
					return *m
				},
			)

			assert.Equal(t, tt.initialSequencerType, m.definition.templateSequencerType, "Initial sequencer type should match")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedSequencerType, m.definition.templateSequencerType, tt.description+" - sequencer type should match expected value")
		})
	}
}

