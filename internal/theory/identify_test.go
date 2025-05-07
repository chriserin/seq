package theory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeChordNotes(t *testing.T) {
	var notes []uint8
	var chord Chord

	notes = Chord{notes: MajorTriad, inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(notes)
	assert.Equal(t, MajorTriad, chord.notes)

	notes = Chord{notes: MajorTriad, inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad, chord.notes)
	assert.Equal(t, int8(1), chord.inversion)

	notes = Chord{notes: MajorTriad, inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad, chord.notes)
	assert.Equal(t, int8(2), chord.inversion)

	notes = Chord{notes: MajorTriad | MajorSeventh, inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|MajorSeventh, chord.notes)
	assert.Equal(t, int8(0), chord.inversion)

	notes = Chord{notes: MajorTriad | MajorSeventh, inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|MajorSeventh, chord.notes)
	assert.Equal(t, int8(2), chord.inversion)

	notes = Chord{notes: MinorTriad | MinorSeventh, inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.notes)
	assert.Equal(t, int8(0), chord.inversion)

	notes = Chord{notes: MinorTriad | MinorSeventh, inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.notes)
	assert.Equal(t, int8(1), chord.inversion)

	notes = Chord{notes: MinorTriad | MinorSeventh, inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.notes)
	assert.Equal(t, int8(2), chord.inversion)

	notes = Chord{notes: MinorTriad | Major9, inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.notes)
	assert.Equal(t, int8(0), chord.inversion)

	notes = Chord{notes: MinorTriad | Major9, inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.notes)
	assert.Equal(t, int8(1), chord.inversion)

	notes = Chord{notes: MinorTriad | Major9, inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.notes)
	assert.Equal(t, int8(2), chord.inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend9, inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend9, chord.notes)
	assert.Equal(t, int8(1), chord.inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend9, inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend9, chord.notes)
	assert.Equal(t, int8(2), chord.inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(0), chord.inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(1), chord.inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(2), chord.inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, inversion: 3}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(3), chord.inversion)
}
