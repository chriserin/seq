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
	Arpeggio Arp
	Double   uint8
}

func (gc GridChord) InBounds(gridKey grid.GridKey) bool {
	return gc.HasNote(gridKey)
}

type BeatNote struct {
	Beat int
	Note grid.Note
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
	for i, interval := range gc.Chord.Intervals() {
		beatnote := gc.Notes[i]
		potentialNotePosition := gc.Key(interval, beatnote)
		if potentialNotePosition == fromKey {
			newRelativeBeat := int(toKey.Beat) - int(gc.Root.Beat)
			gc.Notes[i] = BeatNote{newRelativeBeat, beatnote.Note}
		}
	}
}

func (gc *GridChord) SetChordNote(position grid.GridKey, note grid.Note) {
	for i, interval := range gc.ArpeggioIntervals() {
		beatnote := gc.Notes[i]
		potentialNotePosition := gc.Key(interval, beatnote)
		if potentialNotePosition == position {
			if note == zeronote {
				gc.Notes = slices.Delete(gc.Notes, i, i+1)
				gc.Chord.OmitInterval(interval)
			} else {
				gc.Notes[i] = BeatNote{beatnote.Beat, note}
			}
			break
		}
	}
}

func (gc GridChord) HasNote(position grid.GridKey) bool {
	return slices.Contains(gc.Positions(), position)
}

func (gc GridChord) Positions() []grid.GridKey {
	arpInts := gc.ArpeggioIntervals()
	positions := make([]grid.GridKey, len(arpInts))
	for i, interval := range arpInts {
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
	chordNotes := chord.Intervals()
	beatNotes := make([]BeatNote, len(chordNotes))

	for i := range chord.Intervals() {
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
	gc.ApplyArpeggiation()
}

type Arp int

const (
	ArpNothing Arp = iota
	ArpUp
	ArpReverse
)

func (gc *GridChord) NextArp() {
	gc.Arpeggio = (gc.Arpeggio + 1) % 3
	gc.ApplyArpeggiation()
}

func (gc *GridChord) PrevArp() {
	var newArp Arp
	if gc.Arpeggio == 0 {
		newArp = ArpReverse
	} else {
		newArp = gc.Arpeggio - 1
	}
	gc.Arpeggio = newArp
	gc.ApplyArpeggiation()
}

func (gc *GridChord) NextDouble() {
	gc.Double = (gc.Double + 1) % (uint8(len(gc.Chord.Intervals())) + 1)
	gc.ApplyArpeggiation()
}

func (gc *GridChord) PrevDouble() {
	var newDouble uint8
	if gc.Double == 0 {
		newDouble = uint8(len(gc.Chord.Intervals()))
	} else {
		newDouble = gc.Double - 1
	}
	gc.Double = newDouble
	gc.ApplyArpeggiation()
}

func (gc *GridChord) ApplyArpeggiation() {
	intervals := gc.ArpeggioIntervals()
	existingNotes := gc.Notes
	gc.Notes = make([]BeatNote, 0, len(intervals))
	for i := range intervals {
		var step int
		if gc.Arpeggio == ArpNothing {
			step = 0
		} else {
			step = i
		}
		var newNote grid.Note
		if len(existingNotes) > i {
			newNote = existingNotes[i].Note
		} else {
			newNote = existingNotes[i-len(existingNotes)].Note
		}
		gc.Notes = append(gc.Notes, BeatNote{step, newNote})
	}
}

func (gc GridChord) ArpeggioIntervals() []uint8 {
	intervals := gc.Chord.Intervals()
	doubledIntervals := gc.ApplyDoubles(intervals)
	switch gc.Arpeggio {
	case ArpNothing:
		return doubledIntervals
	case ArpUp:
		return doubledIntervals
	case ArpReverse:
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

// DeepCopy creates a deep copy of the GridChord
func (gc GridChord) DeepCopy() GridChord {
	// Create a new GridChord
	copy := GridChord{
		Root:     gc.Root,     // GridKey is a simple struct, so a direct copy is fine
		Chord:    gc.Chord,    // Direct copy of the Chord
		Arpeggio: gc.Arpeggio, // arp is just an int, so direct copy is fine
		Double:   gc.Double,   // uint8 is a simple type, so direct copy is fine
	}

	// Deep copy the Notes slice
	copy.Notes = make([]BeatNote, len(gc.Notes))
	for i, bn := range gc.Notes {
		copy.Notes[i] = BeatNote{
			Beat: bn.Beat,
			Note: grid.Note{
				AccentIndex: bn.Note.AccentIndex,
				Ratchets:    bn.Note.Ratchets,
				Action:      bn.Note.Action,
				GateIndex:   bn.Note.GateIndex,
				WaitIndex:   bn.Note.WaitIndex,
			},
		}
	}

	return copy
}
