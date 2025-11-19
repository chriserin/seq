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

func TestChordName(t *testing.T) {
	tests := []struct {
		name         string
		chord        Chord
		expectedName string
		description  string
	}{
		// Basic triads
		{
			name:         "Major triad",
			chord:        Chord{Notes: MajorTriad},
			expectedName: "I",
			description:  "Major triad should be named 'I'",
		},
		{
			name:         "Minor triad",
			chord:        Chord{Notes: MinorTriad},
			expectedName: "i",
			description:  "Minor triad should be named 'i'",
		},
		{
			name:         "Diminished triad",
			chord:        Chord{Notes: DiminishedTriad},
			expectedName: "i°",
			description:  "Diminished triad should be named 'i°'",
		},
		{
			name:         "Augmented triad",
			chord:        Chord{Notes: AugmentedTriad},
			expectedName: "I+",
			description:  "Augmented triad should be named 'I+'",
		},
		// Seventh chords
		{
			name:         "Major seventh",
			chord:        Chord{Notes: MajorTriad | MajorSeventh},
			expectedName: "Imaj7",
			description:  "Major triad with major seventh should be 'Imaj7'",
		},
		{
			name:         "Minor seventh",
			chord:        Chord{Notes: MinorTriad | MinorSeventh},
			expectedName: "i7",
			description:  "Minor triad with minor seventh should be 'i7'",
		},
		{
			name:         "Dominant seventh",
			chord:        Chord{Notes: MajorTriad | MinorSeventh},
			expectedName: "I7",
			description:  "Major triad with minor seventh should be 'I7'",
		},
		{
			name:         "Minor-major seventh",
			chord:        Chord{Notes: MinorTriad | MajorSeventh},
			expectedName: "imaj7",
			description:  "Minor triad with major seventh should be 'imaj7'",
		},
		{
			name:         "Half-diminished seventh",
			chord:        Chord{Notes: DiminishedTriad | MinorSeventh},
			expectedName: "i°7",
			description:  "Diminished triad with minor seventh should be 'i°7'",
		},
		{
			name:         "Augmented seventh",
			chord:        Chord{Notes: AugmentedTriad | MinorSeventh},
			expectedName: "I+7",
			description:  "Augmented triad with minor seventh should be 'I+7'",
		},
		// Sixth chords
		{
			name:         "Major sixth",
			chord:        Chord{Notes: MajorTriad | Major6},
			expectedName: "I6",
			description:  "Major triad with major sixth should be 'I6'",
		},
		{
			name:         "Minor sixth",
			chord:        Chord{Notes: MinorTriad | Major6},
			expectedName: "i6",
			description:  "Minor triad with major sixth should be 'i6'",
		},
		// Ninth chords
		{
			name:         "Major ninth",
			chord:        Chord{Notes: MajorTriad | MinorSeventh | Major9},
			expectedName: "I79",
			description:  "Major triad with minor seventh and major ninth should be 'I79'",
		},
		{
			name:         "Minor ninth",
			chord:        Chord{Notes: MinorTriad | MinorSeventh | Major9},
			expectedName: "i79",
			description:  "Minor triad with minor seventh and major ninth should be 'i79'",
		},
		{
			name:         "Major add9",
			chord:        Chord{Notes: MajorTriad | Major9},
			expectedName: "Iadd9",
			description:  "Major triad with major ninth (no seventh) should be 'Iadd9'",
		},
		{
			name:         "Minor add9",
			chord:        Chord{Notes: MinorTriad | Major9},
			expectedName: "iadd9",
			description:  "Minor triad with major ninth (no seventh) should be 'iadd9'",
		},
		{
			name:         "Flat nine chord",
			chord:        Chord{Notes: MajorTriad | MinorSeventh | Minor9},
			expectedName: "I7♭9",
			description:  "Major triad with minor seventh and minor ninth should be 'I7♭9'",
		},
		// Eleventh chords
		{
			name:         "Major eleventh",
			chord:        Chord{Notes: MajorTriad | MinorSeventh | Major9 | Perfect11},
			expectedName: "I7911",
			description:  "Major eleventh chord should be 'I7911'",
		},
		{
			name:         "Major add11",
			chord:        Chord{Notes: MajorTriad | Perfect11},
			expectedName: "Iadd11",
			description:  "Major triad with eleventh (no seventh) should be 'Iadd11'",
		},
		{
			name:         "Sharp eleven chord",
			chord:        Chord{Notes: MajorTriad | MinorSeventh | Major9 | Aug11},
			expectedName: "I79#11",
			description:  "Major triad with sharp eleven should be 'I79#11'",
		},
		// Thirteenth chords
		{
			name:         "Major thirteenth",
			chord:        Chord{Notes: MajorTriad | MinorSeventh | Major9 | Major13},
			expectedName: "I7913",
			description:  "Major thirteenth chord should be 'I7913'",
		},
		{
			name:         "Major add13",
			chord:        Chord{Notes: MajorTriad | Major13},
			expectedName: "Iadd13",
			description:  "Major triad with thirteenth (no seventh) should be 'Iadd13'",
		},
		// Suspended chords
		{
			name:         "Sus4",
			chord:        Chord{Notes: Root | Perfect4 | Perfect5},
			expectedName: "5sus4",
			description:  "Suspended fourth chord should be '5sus4'",
		},
		{
			name:         "Sus2",
			chord:        Chord{Notes: Root | Major2 | Perfect5},
			expectedName: "5sus2",
			description:  "Suspended second chord should be '5sus2'",
		},
		// Power chords and incomplete chords
		{
			name:         "Power chord",
			chord:        Chord{Notes: Root | Perfect5},
			expectedName: "5",
			description:  "Power chord (root and fifth only) should be '5'",
		},
		{
			name:         "Major no5",
			chord:        Chord{Notes: Root | Major3},
			expectedName: "I(no5)",
			description:  "Major third without fifth should be 'I(no5)'",
		},
		{
			name:         "Minor no5",
			chord:        Chord{Notes: Root | Minor3},
			expectedName: "i(no5)",
			description:  "Minor third without fifth should be 'i(no5)'",
		},
		// Complex combinations
		{
			name:         "Major seventh with add9",
			chord:        Chord{Notes: MajorTriad | MajorSeventh | Major9},
			expectedName: "Imaj79",
			description:  "Major seventh with major ninth should be 'Imaj79'",
		},
		{
			name:         "Minor seventh with add9",
			chord:        Chord{Notes: MinorTriad | MinorSeventh | Major9},
			expectedName: "i79",
			description:  "Minor seventh with major ninth should be 'i79'",
		},
		// Add2 and add4 variations
		{
			name:         "Major add2",
			chord:        Chord{Notes: MajorTriad | Major2},
			expectedName: "Iadd2",
			description:  "Major triad with major second should be 'Iadd2'",
		},
		{
			name:         "Minor add2",
			chord:        Chord{Notes: MinorTriad | Major2},
			expectedName: "iadd2",
			description:  "Minor triad with major second should be 'iadd2'",
		},
		{
			name:         "Major add♭2",
			chord:        Chord{Notes: MajorTriad | Minor2},
			expectedName: "Iadd♭2",
			description:  "Major triad with minor second should be 'Iadd♭2'",
		},
		{
			name:         "Minor add♭2",
			chord:        Chord{Notes: MinorTriad | Minor2},
			expectedName: "iadd♭2",
			description:  "Minor triad with minor second should be 'iadd♭2'",
		},
		{
			name:         "Major add4",
			chord:        Chord{Notes: MajorTriad | Perfect4},
			expectedName: "Iadd4",
			description:  "Major triad with perfect fourth should be 'Iadd4'",
		},
		{
			name:         "Minor add4",
			chord:        Chord{Notes: MinorTriad | Perfect4},
			expectedName: "iadd4",
			description:  "Minor triad with perfect fourth should be 'iadd4'",
		},
		{
			name:         "Major7 add2",
			chord:        Chord{Notes: MajorTriad | MajorSeventh | Major2},
			expectedName: "Imaj7add2",
			description:  "Major seventh with major second should be 'Imaj7add2'",
		},
		{
			name:         "Dominant7 add4",
			chord:        Chord{Notes: MajorTriad | MinorSeventh | Perfect4},
			expectedName: "I7add4",
			description:  "Dominant seventh with perfect fourth should be 'I7add4'",
		},
		{
			name:         "Major add2 add4",
			chord:        Chord{Notes: MajorTriad | Major2 | Perfect4},
			expectedName: "Iadd2add4",
			description:  "Major triad with second and fourth should be 'Iadd2add4'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := tt.chord.Name()
			assert.Equal(t, tt.expectedName, name, tt.description)
		})
	}
}
