package overlays

import (
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/theory"
)

type Chords []*GridChord

type GridChord struct {
	Chord theory.Chord
	Notes []BeatNote
	Root  grid.GridKey
}

type BeatNote struct {
	beat uint8
	note grid.Note
}

func (cs Chords) FindChord(position grid.GridKey) (*GridChord, bool) {
	for _, gridChord := range cs {
		if gridChord.ChordBounds().InBounds(position) {
			return gridChord, true
		}
	}
	return nil, false
}

func (cs Chords) FindChordWithNote(gridKey grid.GridKey) (*GridChord, bool) {
	for _, gridChord := range cs {
		if gridChord.HasNote(gridKey) {
			return gridChord, true
		}
	}
	return &GridChord{}, false
}

func (gc GridChord) SetChordNote(position grid.GridKey, note grid.Note) {

	for i, interval := range gc.Chord.Notes() {
		beatnote := gc.Notes[i]
		potentialNotePosition := grid.GridKey{Line: gc.Root.Line - interval, Beat: gc.Root.Beat + beatnote.beat}
		if potentialNotePosition == position {
			gc.Notes[i] = BeatNote{beatnote.beat, note}
		}
	}
}

func (gc GridChord) HasNote(position grid.GridKey) bool {
	for _, gk := range gc.Positions() {
		if gk == position {
			return true
		}
	}
	return false
}

func (gc GridChord) Positions() []grid.GridKey {
	positions := make([]grid.GridKey, len(gc.Notes))

	for i, interval := range gc.Chord.Notes() {
		beatnote := gc.Notes[i]
		positions[i] = grid.GridKey{Line: gc.Root.Line - interval, Beat: gc.Root.Beat + beatnote.beat}
	}
	return positions
}

func (gc GridChord) ChordBounds() grid.Bounds {
	notes := theory.ShiftToZero(gc.Chord.Notes())
	return grid.Bounds{
		Top:    min(gc.Root.Line-notes[len(notes)-1], 0),
		Right:  gc.Root.Beat + gc.Notes[len(gc.Notes)-1].beat,
		Bottom: gc.Root.Line + 2,
		Left:   gc.Root.Beat,
	}
}

func (ol *Overlay) CreateChord(root grid.GridKey, alteration uint32) {
	ol.Chords = append(ol.Chords, InitChord(root, alteration))
}

func InitChord(root grid.GridKey, alteration uint32) *GridChord {
	chord := theory.InitChord(alteration)
	chordNotes := chord.Notes()
	beatNotes := make([]BeatNote, len(chordNotes))

	for i := range chord.Notes() {
		note := grid.InitNote()
		beatNotes[i] = BeatNote{0, note}
	}

	return &GridChord{
		Root:  root,
		Chord: chord,
		Notes: beatNotes,
	}
}

func (gc *GridChord) ApplyAlteration(alteration uint32) {
	gc.Chord.AddNotes(alteration)
	notes := gc.Chord.Notes()
	if len(notes) < len(gc.Notes) {
		for range len(gc.Notes) - len(notes) {
			gc.Notes = append(gc.Notes, BeatNote{0, grid.InitNote()})
		}
	}
}
