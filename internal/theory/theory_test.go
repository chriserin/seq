package theory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNotes(t *testing.T) {
	tests := []struct {
		name          string
		initialChord  Chord
		noteToAdd     uint32
		expectedNotes uint32
		description   string
	}{
		{
			name:          "Add minor second to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Minor2,
			expectedNotes: Minor2,
			description:   "Should add minor second to empty chord",
		},
		{
			name:          "Add major second to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Major2,
			expectedNotes: Major2,
			description:   "Should add major second to empty chord",
		},
		{
			name:          "Replace minor second with major second",
			initialChord:  Chord{Notes: Minor2},
			noteToAdd:     Major2,
			expectedNotes: Major2,
			description:   "Should replace minor second with major second",
		},
		{
			name:          "Replace major second with minor second",
			initialChord:  Chord{Notes: Major2},
			noteToAdd:     Minor2,
			expectedNotes: Minor2,
			description:   "Should replace major second with minor second",
		},
		{
			name:          "Add minor third to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Minor3,
			expectedNotes: Minor3,
			description:   "Should add minor third to empty chord",
		},
		{
			name:          "Add major third to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Major3,
			expectedNotes: Major3,
			description:   "Should add major third to empty chord",
		},
		{
			name:          "Replace minor third with major third",
			initialChord:  Chord{Notes: Minor3},
			noteToAdd:     Major3,
			expectedNotes: Major3,
			description:   "Should replace minor third with major third",
		},
		{
			name:          "Replace major third with minor third",
			initialChord:  Chord{Notes: Major3},
			noteToAdd:     Minor3,
			expectedNotes: Minor3,
			description:   "Should replace major third with minor third",
		},
		{
			name:          "Replace major triad with minor triad",
			initialChord:  Chord{Notes: MajorTriad},
			noteToAdd:     MinorTriad,
			expectedNotes: MinorTriad,
			description:   "Should replace major triad with minor triad",
		},
		{
			name:          "Replace minor triad with major triad",
			initialChord:  Chord{Notes: MinorTriad},
			noteToAdd:     MajorTriad,
			expectedNotes: MajorTriad,
			description:   "Should replace minor triad with major triad",
		},
		{
			name:          "Replace major triad with diminished triad",
			initialChord:  Chord{Notes: MajorTriad},
			noteToAdd:     DiminishedTriad,
			expectedNotes: DiminishedTriad,
			description:   "Should replace major triad with diminished triad",
		},
		{
			name:          "Replace major triad with augmented triad",
			initialChord:  Chord{Notes: MajorTriad},
			noteToAdd:     AugmentedTriad,
			expectedNotes: AugmentedTriad,
			description:   "Should replace major triad with augmented triad",
		},
		{
			name:          "Replace perfect fifth with diminished fifth",
			initialChord:  Chord{Notes: Perfect5},
			noteToAdd:     Dim5,
			expectedNotes: Dim5,
			description:   "Should replace perfect fifth with diminished fifth",
		},
		{
			name:          "Replace perfect fifth with augmented fifth",
			initialChord:  Chord{Notes: Perfect5},
			noteToAdd:     Aug5,
			expectedNotes: Aug5,
			description:   "Should replace perfect fifth with augmented fifth",
		},
		{
			name:          "Add major sixth to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Major6,
			expectedNotes: Major6,
			description:   "Should add major sixth to empty chord",
		},
		{
			name:          "Add major sixth to chord that already has major sixth",
			initialChord:  Chord{Notes: Major6},
			noteToAdd:     Major6,
			expectedNotes: Major6,
			description:   "Should not change chord that already has major sixth",
		},
		{
			name:          "Add octave to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Octave,
			expectedNotes: Octave,
			description:   "Should add octave to empty chord",
		},
		{
			name:          "Add octave to chord that already has octave",
			initialChord:  Chord{Notes: Octave},
			noteToAdd:     Octave,
			expectedNotes: Octave,
			description:   "Should not change chord that already has octave",
		},
		{
			name:          "Add minor ninth to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Minor9,
			expectedNotes: Minor9,
			description:   "Should add minor ninth to empty chord",
		},
		{
			name:          "Add major ninth to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Major9,
			expectedNotes: Major9,
			description:   "Should add major ninth to empty chord",
		},
		{
			name:          "Replace minor ninth with major ninth",
			initialChord:  Chord{Notes: Minor9},
			noteToAdd:     Major9,
			expectedNotes: Major9,
			description:   "Should replace minor ninth with major ninth",
		},
		{
			name:          "Replace major ninth with minor ninth",
			initialChord:  Chord{Notes: Major9},
			noteToAdd:     Minor9,
			expectedNotes: Minor9,
			description:   "Should replace major ninth with minor ninth",
		},
		{
			name:          "Replace minor seventh with major seventh",
			initialChord:  Chord{Notes: MinorSeventh},
			noteToAdd:     MajorSeventh,
			expectedNotes: MajorSeventh,
			description:   "Should replace minor seventh with major seventh",
		},
		{
			name:          "Replace major seventh with minor seventh",
			initialChord:  Chord{Notes: MajorSeventh},
			noteToAdd:     MinorSeventh,
			expectedNotes: MinorSeventh,
			description:   "Should replace major seventh with minor seventh",
		},
		{
			name:          "Add perfect fourth to empty chord",
			initialChord:  Chord{Notes: 0},
			noteToAdd:     Perfect4,
			expectedNotes: Perfect4,
			description:   "Should add perfect fourth to empty chord",
		},
		{
			name:          "Add perfect fourth to major triad",
			initialChord:  Chord{Notes: MajorTriad},
			noteToAdd:     Perfect4,
			expectedNotes: MajorTriad | Perfect4,
			description:   "Should add perfect fourth to major triad",
		},
		{
			name:          "Add multiple intervals to major triad",
			initialChord:  Chord{Notes: MajorTriad},
			noteToAdd:     MinorSeventh | Major9,
			expectedNotes: MajorTriad | MinorSeventh | Major9,
			description:   "Should add multiple intervals to major triad",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chord := tt.initialChord
			chord.AddNotes(tt.noteToAdd)
			assert.Equal(t, tt.expectedNotes, chord.Notes, tt.description)
		})
	}
}

func TestAddNotesComplexScenarios(t *testing.T) {
	tests := []struct {
		name          string
		initialChord  Chord
		notesToAdd    []uint32
		expectedNotes uint32
		description   string
	}{
		{
			name:          "Build major seventh chord step by step",
			initialChord:  Chord{Notes: 0},
			notesToAdd:    []uint32{Root, Major3, Perfect5, MajorSeventh},
			expectedNotes: MajorTriad | MajorSeventh,
			description:   "Should build major seventh chord by adding notes sequentially",
		},
		{
			name:          "Build minor ninth chord step by step",
			initialChord:  Chord{Notes: 0},
			notesToAdd:    []uint32{MinorTriad, MinorSeventh, Major9},
			expectedNotes: MinorTriad | MinorSeventh | Major9,
			description:   "Should build minor ninth chord by adding notes sequentially",
		},
		{
			name:          "Change triad quality in existing chord",
			initialChord:  Chord{Notes: MajorTriad | MajorSeventh},
			notesToAdd:    []uint32{MinorTriad},
			expectedNotes: MinorTriad | MajorSeventh,
			description:   "Should change triad quality while preserving seventh",
		},
		{
			name:          "Add and replace seconds in sequence",
			initialChord:  Chord{Notes: MajorTriad},
			notesToAdd:    []uint32{Minor2, Major2},
			expectedNotes: MajorTriad | Major2,
			description:   "Should add minor second then replace with major second",
		},
		{
			name:          "Add and replace thirds in sequence",
			initialChord:  Chord{Notes: Root | Perfect5},
			notesToAdd:    []uint32{Minor3, Major3},
			expectedNotes: Root | Major3 | Perfect5,
			description:   "Should add minor third then replace with major third",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chord := tt.initialChord
			for _, noteToAdd := range tt.notesToAdd {
				chord.AddNotes(noteToAdd)
			}
			assert.Equal(t, tt.expectedNotes, chord.Notes, tt.description)
		})
	}
}

