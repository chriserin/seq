package overlays

import (
	"slices"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/theory"
)

type Chords []*GridChord

type GridChord struct {
	Chord    theory.Chord
	Notes    []BeatNote
	Root     grid.GridKey
	Arppegio arp
	Double   uint8
}

func (gc GridChord) InBounds(gridKey grid.GridKey) bool {
	return gc.HasNote(gridKey)
}

type BeatNote struct {
	beat int
	note grid.Note
}

func (cs Chords) Remove(gridChord *GridChord) Chords {
	index := slices.Index(cs, gridChord)
	return slices.Delete(cs, index, index+1)
}

func (cs Chords) FindChord(position grid.GridKey) (*GridChord, bool) {
	for _, gridChord := range cs {
		if gridChord.InBounds(position) {
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

func (gc *GridChord) Move(fromKey grid.GridKey, toKey grid.GridKey) {
	for i, interval := range gc.Chord.Notes() {
		beatnote := gc.Notes[i]
		potentialNotePosition := gc.Key(interval, beatnote)
		if potentialNotePosition == fromKey {
			newRelativeBeat := int(toKey.Beat) - int(gc.Root.Beat)
			gc.Notes[i] = BeatNote{newRelativeBeat, beatnote.note}
		}
	}
}

func (gc *GridChord) SetChordNote(position grid.GridKey, note grid.Note) {
	for i, interval := range gc.ArppegioIntervals() {
		beatnote := gc.Notes[i]
		potentialNotePosition := gc.Key(interval, beatnote)
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

	for i, interval := range gc.ArppegioIntervals() {
		beatnote := gc.Notes[i]
		positions[i] = gc.Key(interval, beatnote)
	}
	return positions
}

func (ol *Overlay) CreateChord(root grid.GridKey, alteration uint32) {
	ol.Chords = append(ol.Chords, InitChord(root, alteration))
}

func (ol *Overlay) PasteChord(root grid.GridKey, gridChord *GridChord) {
	dst := new(GridChord)
	*dst = *gridChord
	notes := make([]BeatNote, len(gridChord.Notes))
	copy(notes, gridChord.Notes)
	dst.Notes = notes
	dst.Root = root
	ol.Chords = append(ol.Chords, dst)
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
	gc.ApplyArppegiation()
}

type arp int

const (
	ARP_NOTHING arp = iota
	ARP_UP
	ARP_REVERSE
)

func (gc *GridChord) NextArp() {
	gc.Arppegio = (gc.Arppegio + 1) % 3
	gc.ApplyArppegiation()
}

func (gc *GridChord) PrevArp() {
	var newArp arp
	if gc.Arppegio == 0 {
		newArp = 3
	} else {
		newArp = gc.Arppegio - 1
	}
	gc.Arppegio = newArp
	gc.ApplyArppegiation()
}

func (gc *GridChord) NextDouble() {
	gc.Double = (gc.Double + 1) % (uint8(len(gc.Chord.Notes())) + 1)
	gc.ApplyArppegiation()
}

func (gc *GridChord) PrevDouble() {
	var newDouble uint8
	if gc.Double == 0 {
		newDouble = uint8(len(gc.Chord.Notes()))
	} else {
		newDouble = gc.Double - 1
	}
	gc.Double = newDouble
	gc.ApplyArppegiation()
}

func (gc *GridChord) ApplyArppegiation() {
	intervals := gc.ArppegioIntervals()
	existingNotes := gc.Notes
	gc.Notes = make([]BeatNote, 0, len(intervals))
	for i := range intervals {
		var step int
		if gc.Arppegio == ARP_NOTHING {
			step = 0
		} else {
			step = i
		}
		var newNote grid.Note
		if len(existingNotes) > i {
			newNote = existingNotes[i].note
		} else {
			newNote = existingNotes[i-len(existingNotes)].note
		}
		gc.Notes = append(gc.Notes, BeatNote{step, newNote})
	}
}

func (gc GridChord) ArppegioIntervals() []uint8 {
	intervals := gc.Chord.Notes()
	doubledIntervals := gc.ApplyDoubles(intervals)
	switch gc.Arppegio {
	case ARP_NOTHING:
		return doubledIntervals
	case ARP_UP:
		return doubledIntervals
	case ARP_REVERSE:
		slices.Reverse(doubledIntervals)
		return doubledIntervals
	}
	return doubledIntervals
}

func (gc GridChord) ApplyDoubles(intervals []uint8) []uint8 {
	grownIntervals := slices.Grow(intervals, int(gc.Double))
	for i := range gc.Double {
		if int(i) < len(intervals) {
			grownIntervals = append(grownIntervals, intervals[i]+12)
		}
	}
	return grownIntervals
}
