package theory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeChordNotes(t *testing.T) {
	var notes []uint8
	var chord Chord

	notes = Chord{Notes: MajorTriad, Inversion: 0}.Intervals()
	chord = FromNotesWithAnalysis(notes)
	assert.Equal(t, MajorTriad, chord.Notes)

	notes = Chord{Notes: MajorTriad, Inversion: 1}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad, chord.Notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{Notes: MajorTriad, Inversion: 2}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad, chord.Notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{Notes: MajorTriad | MajorSeventh, Inversion: 0}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|MajorSeventh, chord.Notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{Notes: MajorTriad | MajorSeventh, Inversion: 2}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|MajorSeventh, chord.Notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{Notes: MinorTriad | MinorSeventh, Inversion: 0}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.Notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{Notes: MinorTriad | MinorSeventh, Inversion: 1}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.Notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{Notes: MinorTriad | MinorSeventh, Inversion: 2}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|MinorSeventh, chord.Notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{Notes: MinorTriad | Major9, Inversion: 0}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.Notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{Notes: MinorTriad | Major9, Inversion: 1}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.Notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{Notes: MinorTriad | Major9, Inversion: 2}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MinorTriad|Major9, chord.Notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{Notes: MajorTriad | Minor7 | Extend9, Inversion: 1}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend9, chord.Notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{Notes: MajorTriad | Minor7 | Extend9, Inversion: 2}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend9, chord.Notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{Notes: MajorTriad | Minor7 | Extend911, Inversion: 0}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.Notes)
	assert.Equal(t, int8(0), chord.Inversion)

	notes = Chord{Notes: MajorTriad | Minor7 | Extend911, Inversion: 1}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.Notes)
	assert.Equal(t, int8(1), chord.Inversion)

	notes = Chord{Notes: MajorTriad | Minor7 | Extend911, Inversion: 2}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.Notes)
	assert.Equal(t, int8(2), chord.Inversion)

	notes = Chord{Notes: MajorTriad | Minor7 | Extend911, Inversion: 3}.Intervals()
	chord = FromNotesWithAnalysis(ShiftToZero(notes))
	assert.Equal(t, MajorTriad|Minor7|Extend911, chord.Notes)
	assert.Equal(t, int8(3), chord.Inversion)
}
