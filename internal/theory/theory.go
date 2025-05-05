package theory

import (
	"slices"
)

func InitChord() Chord {
	return Chord{}
}

type Chord struct {
	notes     uint32
	inversion int8
	score     int
}

// Individual note constants in the chromatic scale (relative to root)
const (
	Root      = uint32(1 << 0)  // Root note (0)
	Minor2    = uint32(1 << 1)  // Minor 2nd (1)
	Major2    = uint32(1 << 2)  // Major 2nd (2)
	Minor3    = uint32(1 << 3)  // Minor 3rd (3)
	Major3    = uint32(1 << 4)  // Major 3rd (4)
	Perfect4  = uint32(1 << 5)  // Perfect 4th (5)
	Dim5      = uint32(1 << 6)  // Diminished 5th (6)
	Perfect5  = uint32(1 << 7)  // Perfect 5th (7)
	Aug5      = uint32(1 << 8)  // Augmented 5th (8)
	Major6    = uint32(1 << 9)  // Major 6th (9)
	Minor7    = uint32(1 << 10) // Minor 7th (10)
	Major7    = uint32(1 << 11) // Major 7th (11)
	Octave    = uint32(1 << 12) // Octave (12)
	Minor10   = uint32(1 << 15) // Minor 10th (15)
	Major10   = uint32(1 << 16) // Major 10th (16)
	Perfect12 = uint32(1 << 19) // Perfect 12th (19)

	Minor9    = uint32(1 << 13) // Minor 9th (13)
	Major9    = uint32(1 << 14) // Major 9th (14)
	Dim11     = uint32(1 << 16) // Diminished 11th (16 - same as Major10)
	Perfect11 = uint32(1 << 17) // Perfect 11th (17)
	Aug11     = uint32(1 << 18) // Augmented 11th (18)
	Dim13     = uint32(1 << 20) // Diminished 13th (20)
	Major13   = uint32(1 << 21) // Major 13th (21)
	Aug13     = uint32(1 << 22) // Augmented 13th (22)
)

// Chord type constants composed from individual notes
const (
	MinorTriad      = Root | Minor3 | Perfect5
	MinorTriadI1    = (Minor3 | Perfect5 | Octave) >> 3
	MinorTriadI2    = (Perfect5 | Octave | (Minor3 << 12)) >> 7
	MajorTriad      = Root | Major3 | Perfect5
	MajorTriadI1    = (Major3 | Perfect5 | Octave) >> 4
	MajorTriadI2    = (Perfect5 | Octave | Dim11) >> 7
	AugmentedTriad  = Root | Major3 | Aug5
	DiminishedTriad = Root | Minor3 | Dim5

	MinorSeventh = Minor7
	MajorSeventh = Major7

	Extend9     = Major9
	Extend911   = Major9 | Perfect11
	Extend91113 = Major9 | Perfect11 | Major13
)

// AddNotes adds specified notes to the chord
func (c *Chord) AddNotes(noteConstant uint32) Chord {
	oldNotes := c.notes
	if IsTriad(noteConstant) {
		currentTriad := c.CurrentTriad()
		c.Replace(currentTriad, noteConstant)
	} else {
		c.notes |= noteConstant
	}

	return Chord{notes: oldNotes}
}

var triads = []uint32{MajorTriad, MinorTriad, DiminishedTriad, AugmentedTriad}

func IsTriad(noteConstant uint32) bool {
	for _, t := range triads {
		if ContainsBits(noteConstant, t) {
			return true
		}
	}
	return false
}

func (c Chord) CurrentTriad() uint32 {
	for _, t := range triads {
		if ContainsBits(c.notes, t) {
			return t
		}
	}
	return 0
}

func (c *Chord) Replace(oldNotes uint32, newNotes uint32) {
	c.notes ^= oldNotes
	c.notes |= newNotes
}

func (c Chord) UninvertedNotes() []int {
	notes := make([]int, 0)
	for i := 0; i < 32; i++ {
		if c.notes&(1<<i) != 0 {
			notes = append(notes, i)
		}
	}
	return notes
}

// Notes returns a slice of integers representing the notes in the chord
// If the chord has an inversion value, the appropriate number of notes
// from the bottom of the chord are moved to the top
func (c Chord) Notes() []int {
	// First collect all notes without considering inversion
	notes := c.UninvertedNotes()

	// Apply inversion if needed
	noteCount := len(notes)
	if noteCount > 0 && c.inversion > 0 && int(c.inversion) < noteCount {
		// Move the first 'inversion' notes to the end, raising them by an octave
		invertedNotes := make([]int, 0, noteCount)

		// Add the remaining notes first (notes after the inversion point)
		invertedNotes = append(invertedNotes, notes[c.inversion:]...)

		// Add the inverted notes (notes before the inversion point), raised by an octave (12 semitones)
		for i := 0; i < int(c.inversion); i++ {
			// For the second inversion of a C major triad (0,4,7),
			// we want to move 0 and 4 up an octave, resulting in (7,12,16)
			invertedNotes = append(invertedNotes, notes[i]+12)
		}
		slices.Sort(invertedNotes)
		return invertedNotes
	}

	return notes
}

func ContainsBits(source, pattern uint32) bool {
	return (source & pattern) == pattern
}

func ChordFromNotes(notes []uint8) Chord {
	var chordBits uint32

	for _, note := range notes {
		if note < 32 {
			chordBits |= uint32(1 << note)
		}
	}

	return Chord{notes: chordBits}
}
