package main

import (
	"testing"

	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

func TestUpdateOverlayKey(t *testing.T) {
	tests := []struct {
		name        string
		expectedKey overlaykey.OverlayPeriodicity
		commands    []any
		description string
	}{
		{
			name:        "New Overlay Key with shift",
			commands:    []any{mappings.OverlayInputSwitch, TestKey{Keys: "2"}, mappings.Enter},
			expectedKey: overlaykey.InitOverlayKey(2, 1),
			description: "Should create a new overlay key with shift",
		},
		{
			name:        "New Overlay Key with width",
			commands:    []any{mappings.OverlayInputSwitch, TestKey{Keys: "3:2"}, mappings.Enter},
			expectedKey: overlaykey.OverlayPeriodicity{Shift: 3, Interval: 1, Width: 2, StartCycle: 0},
			description: "Should create a new overlay key with width",
		},
		{
			name:        "New Overlay Key with interval",
			commands:    []any{mappings.OverlayInputSwitch, TestKey{Keys: "3/2"}, mappings.Enter},
			expectedKey: overlaykey.OverlayPeriodicity{Shift: 3, Interval: 2, Width: 0, StartCycle: 0},
			description: "Should create a new overlay key with interval",
		},
		{
			name:        "New Overlay Key with StartCycle",
			commands:    []any{mappings.OverlayInputSwitch, TestKey{Keys: "3S2"}, mappings.Enter},
			expectedKey: overlaykey.OverlayPeriodicity{Shift: 3, Interval: 1, Width: 0, StartCycle: 2},
			description: "Should create a new overlay key with StartCycle",
		},
		{
			name:        "New Overlay Key with StartCycle lower case s",
			commands:    []any{mappings.OverlayInputSwitch, TestKey{Keys: "3S2s"}, mappings.Enter},
			expectedKey: overlaykey.OverlayPeriodicity{Shift: 3, Interval: 1, Width: 0, StartCycle: 0},
			description: "Should create a new overlay key with StartCycle lower case s",
		},
		{
			name:        "New Overlay Key with all attributes",
			commands:    []any{mappings.OverlayInputSwitch, TestKey{Keys: "3:3/3S3"}, mappings.Enter},
			expectedKey: overlaykey.OverlayPeriodicity{Shift: 3, Interval: 3, Width: 3, StartCycle: 3},
			description: "Should create a new overlay key with StartCycle lower case s",
		},
		{
			name:        "Escape from overlay key edit",
			commands:    []any{mappings.OverlayInputSwitch, TestKey{Keys: "3:3/3S3"}, mappings.Escape},
			expectedKey: overlaykey.OverlayPeriodicity{Shift: 1, Interval: 1, Width: 0, StartCycle: 0},
			description: "Should escape from overlay key edit and return to current key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedKey, m.currentOverlay.Key, tt.description)
			assert.Equal(t, tt.expectedKey, m.overlayKeyEdit.GetKey(), tt.description)
		})
	}
}
