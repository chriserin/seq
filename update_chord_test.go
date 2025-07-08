package main

import (
	"slices"
	"strconv"
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/theory"
	"github.com/stretchr/testify/assert"
)

func TestChordTriads(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedChord     uint32
		expectedIntervals []uint8
		description       string
	}{
		{
			name:              "MajorTriad creates major triad chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad},
			expectedChord:     theory.MajorTriad,
			expectedIntervals: []uint8{0, 4, 7},
			description:       "Should create major triad chord",
		},
		{
			name:              "MinorTriad creates minor triad chord",
			commands:          []any{mappings.CursorLastLine, mappings.MinorTriad},
			expectedChord:     theory.MinorTriad,
			expectedIntervals: []uint8{0, 3, 7},
			description:       "Should create minor triad chord",
		},
		{
			name:              "AugmentedTriad creates augmented triad chord",
			commands:          []any{mappings.CursorLastLine, mappings.AugmentedTriad},
			expectedChord:     theory.AugmentedTriad,
			expectedIntervals: []uint8{0, 4, 8},
			description:       "Should create augmented triad chord",
		},
		{
			name:              "DiminishedTriad creates diminished triad chord",
			commands:          []any{mappings.CursorLastLine, mappings.DiminishedTriad},
			expectedChord:     theory.DiminishedTriad,
			expectedIntervals: []uint8{0, 3, 6},
			description:       "Should create diminished triad chord",
		},
		{
			name:              "MinorSeventh adds minor seventh to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MinorSeventh},
			expectedChord:     theory.MajorTriad | theory.MinorSeventh,
			expectedIntervals: []uint8{0, 4, 7, 10},
			description:       "Should add minor seventh to existing chord",
		},
		{
			name:              "MajorSeventh adds major seventh to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorSeventh},
			expectedChord:     theory.MajorTriad | theory.MajorSeventh,
			expectedIntervals: []uint8{0, 4, 7, 11},
			description:       "Should add major seventh to existing chord",
		},
		{
			name:              "AugFifth adds augmented fifth to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.AugFifth},
			expectedChord:     theory.Root | theory.Major3 | theory.Aug5,
			expectedIntervals: []uint8{0, 4, 8},
			description:       "Should add augmented fifth to existing chord",
		},
		{
			name:              "DimFifth adds diminished fifth to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.DimFifth},
			expectedChord:     theory.Root | theory.Major3 | theory.Dim5,
			expectedIntervals: []uint8{0, 4, 6},
			description:       "Should add diminished fifth to existing chord",
		},
		{
			name:              "PerfectFifth adds perfect fifth to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.DimFifth, mappings.PerfectFifth},
			expectedChord:     theory.MajorTriad | theory.Perfect5,
			expectedIntervals: []uint8{0, 4, 7},
			description:       "Should add perfect fifth to existing chord",
		},
		{
			name:              "IncreaseInversions increments chord inversion",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.IncreaseInversions},
			expectedChord:     theory.MajorTriad,
			expectedIntervals: []uint8{4, 7, 12},
			description:       "Should increment chord inversion from 0 to 1",
		},
		{
			name:              "DecreaseInversions decrements chord inversion",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.IncreaseInversions, mappings.IncreaseInversions, mappings.DecreaseInversions},
			expectedChord:     theory.MajorTriad,
			expectedIntervals: []uint8{4, 7, 12},
			description:       "Should decrement chord inversion from 2 to 1",
		},
		{
			name:              "Multiple inversion increases",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.IncreaseInversions, mappings.IncreaseInversions},
			expectedChord:     theory.MajorTriad,
			expectedIntervals: []uint8{7, 12, 16},
			description:       "Should increment chord inversion to 3",
		},
		{
			name:              "Multiple inversion increases wraps around",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.IncreaseInversions, mappings.IncreaseInversions, mappings.IncreaseInversions},
			expectedChord:     theory.MajorTriad,
			expectedIntervals: []uint8{0, 4, 7},
			description:       "Should increment chord inversion to 3",
		},
		{
			name:              "OmitRoot removes root from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.OmitRoot},
			expectedIntervals: []uint8{4, 7},
			expectedChord:     theory.MajorTriad & ^theory.Root,
			description:       "Should remove root note from major triad",
		},
		{
			name:              "OmitThird removes third from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.OmitThird},
			expectedIntervals: []uint8{0, 7},
			expectedChord:     theory.MajorTriad & ^theory.Major3,
			description:       "Should remove third note from major triad",
		},
		{
			name:              "OmitFifth removes fifth from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.OmitFifth},
			expectedIntervals: []uint8{0, 4},
			expectedChord:     theory.MajorTriad & ^theory.Perfect5,
			description:       "Should remove fifth note from major triad",
		},
		{
			name:              "OmitSeventh removes seventh from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorSeventh, mappings.OmitSeventh},
			expectedIntervals: []uint8{0, 4, 7},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove seventh note from major seventh chord",
		},
		{
			name:              "Double adds double note to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.NextDouble},
			expectedIntervals: []uint8{0, 4, 7, 12},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove seventh note from major seventh chord",
		},
		{
			name:              "Double twice adds double note to chord twice",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.NextDouble, mappings.NextDouble},
			expectedIntervals: []uint8{0, 4, 7, 12, 16},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove seventh note from major seventh chord",
		},
		{
			name:              "Double thrice adds double note to chord thrice",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.NextDouble, mappings.NextDouble, mappings.NextDouble},
			expectedIntervals: []uint8{0, 4, 7, 12, 16, 19},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove seventh note from major seventh chord",
		},
		{
			name:              "Double four times wraps around",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.NextDouble, mappings.NextDouble, mappings.NextDouble, mappings.NextDouble},
			expectedIntervals: []uint8{0, 4, 7},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove seventh note from major seventh chord",
		},
		{
			name:              "Prev Double wraps around to max doubles",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.PrevDouble},
			expectedIntervals: []uint8{0, 4, 7, 12, 16, 19},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove seventh note from major seventh chord",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithPolyphony(),
			)

			m, _ = processCommands(tt.commands, m)

			for currentInterval, line := uint8(0), uint8(len(m.definition.lines))-1; line != 255; line, currentInterval = line-1, currentInterval+1 {
				key := grid.GridKey{Line: line, Beat: 0}
				overlayChord, exists := m.currentOverlay.FindChord(key)
				if slices.Contains(tt.expectedIntervals, currentInterval) {
					assert.True(t, exists, tt.description+" - chord should exist - line "+strconv.Itoa(int(line)))
					if exists {
						assert.Equal(t, tt.expectedChord, overlayChord.GridChord.Chord.Notes)
					}
				} else {
					assert.False(t, exists, tt.description+" - chord should not exist - line "+strconv.Itoa(int(line)))
				}
			}
		})
	}
}

func TestRemoveChord(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name:        "RemoveChord removes existing chord",
			commands:    []any{mappings.CursorLineEnd, mappings.MajorTriad, mappings.RemoveChord},
			description: "Should remove existing chord",
		},
		{
			name:        "RemoveChord on empty position does nothing",
			commands:    []any{mappings.CursorLineEnd, mappings.RemoveChord},
			description: "Should not crash when removing non-existent chord",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithPolyphony(),
			)

			m, _ = processCommands(tt.commands, m)

			chord := m.CurrentChord()
			assert.False(t, chord.HasValue(), tt.description+" - chord should not exist after removal")
		})
	}
}

func TestConvertToNotes(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name:        "ConvertToNotes converts chord to individual notes",
			commands:    []any{mappings.CursorLineEnd, mappings.MajorTriad, mappings.ConvertToNotes},
			description: "Should convert chord to individual notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithPolyphony(),
			)

			m, _ = processCommands(tt.commands, m)

			chord := m.CurrentChord()
			assert.False(t, chord.HasValue(), tt.description+" - chord should not exist after conversion")

			_, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - at least one note should exist after conversion")
		})
	}
}
