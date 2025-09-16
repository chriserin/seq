package main

import (
	"testing"

	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/operation"
	"github.com/chriserin/seq/internal/sequence"
	"github.com/stretchr/testify/assert"
)

func TestAccentInputSwitchDiffAndData(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []any
		expectedAccentStart  uint8
		expectedAccentValues []config.Accent
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentStart:  121,
			expectedAccentValues: []config.Accent{0, 121, 106, 91, 76, 60, 45, 30, 15},
			description:          "Should select accent input and increase should set accent values",
		},
		{
			name:                 "Accent Input Switch with Decrease",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			expectedAccentStart:  119,
			expectedAccentValues: []config.Accent{0, 119, 104, 89, 74, 60, 45, 30, 15},
			description:          "Should select accent input and decrease should set accent values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectAccentStart, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentStart, m.definition.Accents.Start, tt.description+" - accent End value")

			for i, value := range tt.expectedAccentValues {
				assert.Equalf(t, value, m.definition.Accents.Data[i], "accent values should match expected values - %d == %d", value, m.definition.Accents.Data[i])
			}
		})
	}
}

func TestAccentInputSwitchTarget(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []any
		expectedAccentTarget sequence.AccentTarget
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase on Target",
			commands:             []any{mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentTarget: sequence.AccentTargetNote,
			description:          "Should select accent input and increase should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target",
			commands:             []any{mappings.AccentInputSwitch, mappings.Decrease},
			expectedAccentTarget: sequence.AccentTargetNote,
			description:          "Should select accent input and decrease should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target Wraparound",
			commands:             []any{mappings.AccentInputSwitch, mappings.Decrease, mappings.Decrease},
			expectedAccentTarget: sequence.AccentTargetVelocity,
			description:          "Should select accent input and decrease should set accent target to Velocity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectAccentTarget, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentTarget, m.definition.Accents.Target, tt.description)
		})
	}
}

func TestAccentInputSwitchStartValue(t *testing.T) {
	tests := []struct {
		name                string
		commands            []any
		initialAccentStart  uint8
		expectedAccentStart uint8
		description         string
	}{
		{
			name:                "Accent Input Switch with Increase on Start Value",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  120,
			expectedAccentStart: 121,
			description:         "Should select accent input and increase should set accent start value to 0",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			initialAccentStart:  120,
			expectedAccentStart: 119,
			description:         "Should select accent input and decrease should set accent start value to 127",
		},
		{
			name:                "Accent Input Switch with Increase on Start Value at Boundary",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  127,
			expectedAccentStart: 127,
			description:         "Should select accent input and increase should not go above 127",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value at Boundary",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			initialAccentStart:  2,
			expectedAccentStart: 2,
			description:         "Should select accent input and decrease should not go below 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.Accents.Start = tt.initialAccentStart
					m.definition.Accents.End = 1
					m.definition.Accents.ReCalc()
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, operation.SelectAccentStart, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentStart, m.definition.Accents.Start, tt.description+" - accent start value")
		})
	}
}
