package main

import (
	"testing"

	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/operation"
	"github.com/stretchr/testify/assert"
)

func TestToggleChordMode(t *testing.T) {
	tests := []struct {
		name                  string
		commands              []any
		initialSequencerType  operation.SequencerMode
		expectedSequencerType operation.SequencerMode
		description           string
	}{
		{
			name:                  "Toggle from trigger to polyphony",
			commands:              []any{mappings.ToggleChordMode},
			initialSequencerType:  operation.SeqModeLine,
			expectedSequencerType: operation.SeqModeChord,
			description:           "Should toggle from trigger mode to polyphony (chord) mode",
		},
		{
			name:                  "Toggle from polyphony to trigger",
			commands:              []any{mappings.ToggleChordMode},
			initialSequencerType:  operation.SeqModeChord,
			expectedSequencerType: operation.SeqModeLine,
			description:           "Should toggle from polyphony (chord) mode to trigger mode",
		},
		{
			name:                  "Multiple toggles return to original state",
			commands:              []any{mappings.ToggleChordMode, mappings.ToggleChordMode},
			initialSequencerType:  operation.SeqModeLine,
			expectedSequencerType: operation.SeqModeLine,
			description:           "Should return to original state after two toggles",
		},
		{
			name:                  "Multiple toggles from polyphony return to original state",
			commands:              []any{mappings.ToggleChordMode, mappings.ToggleChordMode},
			initialSequencerType:  operation.SeqModeChord,
			expectedSequencerType: operation.SeqModeChord,
			description:           "Should return to original polyphony state after two toggles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.TemplateSequencerType = tt.initialSequencerType
					return *m
				},
			)

			assert.Equal(t, tt.initialSequencerType, m.definition.TemplateSequencerType, "Initial sequencer type should match")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedSequencerType, m.definition.TemplateSequencerType, tt.description+" - sequencer type should match expected value")
		})
	}
}
