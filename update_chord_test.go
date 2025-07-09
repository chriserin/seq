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
			name:              "OmitSecond removes second from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorSecond, mappings.OmitSecond},
			expectedIntervals: []uint8{0, 4, 7},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove second note from chord",
		},
		{
			name:              "OmitFourth removes fourth from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.PerfectFourth, mappings.OmitFourth},
			expectedIntervals: []uint8{0, 4, 7},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove fourth note from chord",
		},
		{
			name:              "OmitSixth removes sixth from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorSixth, mappings.OmitSixth},
			expectedIntervals: []uint8{0, 4, 7},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove sixth note from chord",
		},
		{
			name:              "OmitNinth removes ninth from chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorNinth, mappings.OmitNinth},
			expectedIntervals: []uint8{0, 4, 7},
			expectedChord:     theory.MajorTriad,
			description:       "Should remove ninth note from chord",
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
		{
			name:              "MinorSecond adds minor second interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MinorSecond},
			expectedIntervals: []uint8{0, 1, 4, 7},
			expectedChord:     theory.MajorTriad | theory.Minor2,
			description:       "Should add minor second interval to major triad",
		},
		{
			name:              "MajorSecond adds major second interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorSecond},
			expectedIntervals: []uint8{0, 2, 4, 7},
			expectedChord:     theory.MajorTriad | theory.Major2,
			description:       "Should add major second interval to major triad",
		},
		{
			name:              "MinorThird adds minor third interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MinorThird},
			expectedIntervals: []uint8{0, 3, 7},
			expectedChord:     theory.Root | theory.Minor3 | theory.Perfect5,
			description:       "Should add minor third interval to chord",
		},
		{
			name:              "MajorThird adds major third interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MinorTriad, mappings.MajorThird},
			expectedIntervals: []uint8{0, 4, 7},
			expectedChord:     theory.Root | theory.Major3 | theory.Perfect5,
			description:       "Should add major third interval to chord",
		},
		{
			name:              "PerfectFourth adds perfect fourth interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.PerfectFourth},
			expectedIntervals: []uint8{0, 4, 5, 7},
			expectedChord:     theory.MajorTriad | theory.Perfect4,
			description:       "Should add perfect fourth interval to major triad",
		},
		{
			name:              "MajorSixth adds major sixth interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorSixth},
			expectedIntervals: []uint8{0, 4, 7, 9},
			expectedChord:     theory.MajorTriad | theory.Major6,
			description:       "Should add major sixth interval to major triad",
		},
		{
			name:              "Octave adds octave interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.Octave},
			expectedIntervals: []uint8{0, 4, 7, 12},
			expectedChord:     theory.MajorTriad | theory.Octave,
			description:       "Should add octave interval to major triad",
		},
		{
			name:              "MinorNinth adds minor ninth interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MinorNinth},
			expectedIntervals: []uint8{0, 4, 7, 13},
			expectedChord:     theory.MajorTriad | theory.Minor9,
			description:       "Should add minor ninth interval to major triad",
		},
		{
			name:              "MajorNinth adds major ninth interval to chord",
			commands:          []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.MajorNinth},
			expectedIntervals: []uint8{0, 4, 7, 14},
			expectedChord:     theory.MajorTriad | theory.Major9,
			description:       "Should add major ninth interval to major triad",
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

func TestArpeggioMappings(t *testing.T) {
	tests := []struct {
		name         string
		commands     []any
		expectedKeys []grid.GridKey
		description  string
	}{
		{
			name:         "NextArpeggio changes from ArpNothing to ArpUp",
			commands:     []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.NextArpeggio},
			expectedKeys: []grid.GridKey{{Line: 23, Beat: 0}, {Line: 19, Beat: 1}, {Line: 16, Beat: 2}},
			description:  "Should change arpeggio from ArpNothing to ArpUp",
		},
		{
			name:         "NextArpeggio changes from ArpUp to ArpReverse",
			commands:     []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.NextArpeggio, mappings.NextArpeggio},
			expectedKeys: []grid.GridKey{{Line: 23, Beat: 2}, {Line: 19, Beat: 1}, {Line: 16, Beat: 0}},
			description:  "Should change arpeggio from ArpUp to ArpReverse",
		},
		{
			name:         "NextArpeggio wraps around from ArpReverse to ArpNothing",
			commands:     []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.NextArpeggio, mappings.NextArpeggio, mappings.NextArpeggio},
			expectedKeys: []grid.GridKey{{Line: 23, Beat: 0}, {Line: 19, Beat: 0}, {Line: 16, Beat: 0}},
			description:  "Should wrap around from ArpReverse to ArpNothing",
		},
		{
			name:         "PrevArpeggio changes from ArpNothing to ArpReverse",
			commands:     []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.PrevArpeggio},
			expectedKeys: []grid.GridKey{{Line: 23, Beat: 2}, {Line: 19, Beat: 1}, {Line: 16, Beat: 0}},
			description:  "Should change arpeggio from ArpNothing to ArpReverse",
		},
		{
			name:         "PrevArpeggio changes from ArpReverse to ArpUp",
			commands:     []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.PrevArpeggio, mappings.PrevArpeggio},
			expectedKeys: []grid.GridKey{{Line: 23, Beat: 0}, {Line: 19, Beat: 1}, {Line: 16, Beat: 2}},
			description:  "Should change arpeggio from ArpReverse to ArpUp",
		},
		{
			name:         "PrevArpeggio wraps around from ArpUp to ArpNothing",
			commands:     []any{mappings.CursorLastLine, mappings.MajorTriad, mappings.PrevArpeggio, mappings.PrevArpeggio, mappings.PrevArpeggio},
			expectedKeys: []grid.GridKey{{Line: 23, Beat: 0}, {Line: 19, Beat: 0}, {Line: 16, Beat: 0}},
			description:  "Should wrap around from ArpUp to ArpNothing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithPolyphony(),
			)

			m, _ = processCommands(tt.commands, m)

			chord := m.CurrentChord()
			assert.True(t, chord.HasValue(), tt.description+" - chord should exist")

			gridPattern := make(grid.Pattern)
			m.currentOverlay.CombinePattern(&gridPattern, 1)

			for key := range gridPattern {
				assert.Contains(t, tt.expectedKeys, key, tt.description+" - expected key should exist in pattern")
			}
			for _, key := range tt.expectedKeys {
				assert.Contains(t, gridPattern, key, tt.description+" - pattern should contain expected key "+key.String())
			}
		})
	}
}
