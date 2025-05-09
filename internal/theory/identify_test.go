package theory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeChordNotes(t *testing.T) {
	var notes []uint8
	var chord Chord

	notes = Chord{notes: MajorTriad, Inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(notes)
	assert.Equal(t, MajorTriad, chord.notes)

	notes = Chord{notes: MajorTriad, Inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad, chord.notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{notes: MajorTriad, Inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad, chord.notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{notes: MajorTriad | MajorSeventh, Inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|MajorSeventh, chord.notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{notes: MajorTriad | MajorSeventh, Inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|MajorSeventh, chord.notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{notes: MinorTriad | MinorSeventh, Inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{notes: MinorTriad | MinorSeventh, Inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{notes: MinorTriad | MinorSeventh, Inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{notes: MinorTriad | Major9, Inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{notes: MinorTriad | Major9, Inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{notes: MinorTriad | Major9, Inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend9, Inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend9, chord.notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend9, Inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend9, chord.notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, Inversion: 0}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, Inversion: 1}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, Inversion: 2}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{notes: MajorTriad | Minor7 | Extend911, Inversion: 3}.Notes()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.notes)
	assert.Equal(t, int8(3), chord.Inversion)
}
