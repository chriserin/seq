package main

import (
	"testing"

	"github.com/chriserin/seq/internal/mappings"
	"github.com/stretchr/testify/assert"
)

func TestAccentInputSwitchDiffAndData(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []any
		expectedAccentDiff   uint8
		expectedAccentValues []uint8
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase",
			commands:             []any{mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentDiff:   16,
			expectedAccentValues: []uint8{0, 120, 104, 88, 72, 56, 40, 24, 8},
			description:          "Should select accent input and increase should set accent values",
		},
		{
			name:                 "Accent Input Switch with Decrease",
			commands:             []any{mappings.AccentInputSwitch, mappings.Decrease},
			expectedAccentDiff:   14,
			expectedAccentValues: []uint8{0, 120, 106, 92, 78, 64, 50, 36, 22},
			description:          "Should select accent input and decrease should set accent values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectAccentDiff, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentDiff, m.definition.accents.Diff, tt.description+" - accent diff value")
			for i, value := range tt.expectedAccentValues {
				assert.Equalf(t, value, m.definition.accents.Data[i].Value, "accent values should match expected values - %d == %d", value, m.definition.accents.Data[0].Value)
			}
		})
	}
}

func TestAccentInputSwitchTarget(t *testing.T) {
	tests := []struct {
		name                 string
		commands             []any
		expectedAccentTarget accentTarget
		description          string
	}{
		{
			name:                 "Accent Input Switch with Increase on Target",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			expectedAccentTarget: AccentTargetNote,
			description:          "Should select accent input and increase should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			expectedAccentTarget: AccentTargetNote,
			description:          "Should select accent input and decrease should set accent target to Note",
		},
		{
			name:                 "Accent Input Switch with Decrease on Target Wraparound",
			commands:             []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease, mappings.Decrease},
			expectedAccentTarget: AccentTargetVelocity,
			description:          "Should select accent input and decrease should set accent target to Velocity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectAccentTarget, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentTarget, m.definition.accents.Target, tt.description)
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
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  120,
			expectedAccentStart: 121,
			description:         "Should select accent input and increase should set accent start value to 0",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			initialAccentStart:  120,
			expectedAccentStart: 119,
			description:         "Should select accent input and decrease should set accent start value to 127",
		},
		{
			name:                "Accent Input Switch with Increase on Start Value at Boundary",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Increase},
			initialAccentStart:  127,
			expectedAccentStart: 127,
			description:         "Should select accent input and increase should not go above 127",
		},
		{
			name:                "Accent Input Switch with Decrease on Start Value at Boundary",
			commands:            []any{mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.AccentInputSwitch, mappings.Decrease},
			initialAccentStart:  105,
			expectedAccentStart: 105,
			description:         "Should select accent input and decrease should not go below 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.definition.accents.Start = tt.initialAccentStart
					m.definition.accents.ReCalc()
					return *m
				},
			)

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, SelectAccentStart, m.selectionIndicator, tt.description+" - selection state")
			assert.Equal(t, tt.expectedAccentStart, m.definition.accents.Start, tt.description+" - accent start value")
		})
	}
}
