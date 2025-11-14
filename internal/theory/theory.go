// Package theory provides music theory functionality for the sequencer.
// It implements chord construction, interval calculations, inversions, and
// musical pattern recognition using bit manipulation for efficient chord
// representation and manipulation. Supports various chord types, voicings,
// and harmonic analysis for polyphonic sequencing.
package theory

func InitChord(alteration uint32) Chord {
	return Chord{Notes: alteration}
}

type Chord struct {
	Notes uint32
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
func (c *Chord) AddNotes(noteConstant uint32) {

	if IsTriad(noteConstant) {
		currentTriad := c.CurrentTriad()
		c.Replace(currentTriad, noteConstant)
	} else if IsSecond(noteConstant) {
		currentSecond := c.Notes & (Minor2 | Major2)
		c.Replace(currentSecond, noteConstant)
	} else if IsThird(noteConstant) {
		currentThird := c.Notes & (Minor3 | Major3)
		c.Replace(currentThird, noteConstant)
	} else if IsFifth(noteConstant) {
		currentFifth := c.CurrentFifth()
		c.Replace(currentFifth, noteConstant)
	} else if IsSixth(noteConstant) {
		c.Replace(Major6, noteConstant)
	} else if IsOctave(noteConstant) {
		c.Replace(Octave, noteConstant)
	} else if IsNinth(noteConstant) {
		c.Replace(c.Notes&(Minor9|Major9), noteConstant)
	} else if IsSeventh(noteConstant) {
		currentSeventh := c.CurrentSeventh()
		c.Replace(currentSeventh, noteConstant)
	} else {
		c.Notes |= noteConstant
	}
}

const (
	OmitRoot       = Root
	OmitSecond     = Minor2 | Major2
	OmitThird      = Minor3 | Major3
	OmitFourth     = Perfect4
	OmitFifth      = Dim5 | Perfect5 | Aug5
	OmitSixth      = Major6
	OmitSeventh    = Minor7 | Major7
	OmitOctave     = Octave
	OmitNinth      = Minor9 | Major9
	OmitEleventh   = Dim11 | Perfect11 | Aug11
	OmitThirteenth = Dim13 | Major13 | Aug13
)

func (c *Chord) OmitNote(omitNoteConstant uint32) {
	c.Notes &= ^omitNoteConstant
}

func (c *Chord) OmitInterval(omitNoteConstant uint8) {
	c.Notes &= ^(uint32(1 << omitNoteConstant))
}

var triads = []uint32{MajorTriad, MinorTriad, DiminishedTriad, AugmentedTriad}

var seconds = []uint32{Minor2, Major2}

func IsSecond(noteConstant uint32) bool {
	for _, t := range seconds {
		if ContainsBits(noteConstant, t) {
			return true
		}
	}
	return false
}

var thirds = []uint32{Minor3, Major3}

func IsThird(noteConstant uint32) bool {
	for _, t := range thirds {
		if ContainsBits(noteConstant, t) {
			return true
		}
	}
	return false
}

func IsFourth(noteConstant uint32) bool {
	return ContainsBits(noteConstant, Perfect4)
}

func IsSixth(noteConstant uint32) bool {
	return ContainsBits(noteConstant, Major6)
}

func IsOctave(noteConstant uint32) bool {
	return ContainsBits(noteConstant, Octave)
}

var ninths = []uint32{Minor9, Major9}

func IsNinth(noteConstant uint32) bool {
	for _, t := range ninths {
		if ContainsBits(noteConstant, t) {
			return true
		}
	}
	return false
}

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
		if ContainsBits(c.Notes, t) {
			return t
		}
	}
	return 0
}

var fifths = []uint32{Perfect5, Dim5, Aug5}

func IsFifth(noteConstant uint32) bool {
	for _, t := range fifths {
		if ContainsBits(noteConstant, t) {
			return true
		}
	}
	return false
}

func (c Chord) CurrentFifth() uint32 {
	for _, t := range fifths {
		if ContainsBits(c.Notes, t) {
			return t
		}
	}
	return 0
}

var sevenths = []uint32{MajorSeventh, MinorSeventh}

func IsSeventh(noteConstant uint32) bool {
	for _, t := range sevenths {
		if ContainsBits(noteConstant, t) {
			return true
		}
	}
	return false
}

func (c Chord) CurrentSeventh() uint32 {
	for _, t := range sevenths {
		if ContainsBits(c.Notes, t) {
			return t
		}
	}
	return 0
}

func (c *Chord) Replace(oldNotes uint32, newNotes uint32) {
	c.Notes ^= oldNotes
	c.Notes |= newNotes
}

func (c Chord) UninvertedNotes() []uint8 {
	notes := make([]uint8, 0)
	for i := range uint8(32) {
		if c.Notes&(1<<i) != 0 {
			notes = append(notes, i)
		}
	}
	return notes
}

func (c Chord) FirstInterval() uint8 {
	for i := range uint8(32) {
		if c.Notes&(1<<i) != 0 {
			return uint8(i)
		}
	}
	return 33
}

func (c Chord) LastInterval() uint8 {
	for i := uint8(31); i >= 0; i-- {
		if c.Notes&(1<<i) != 0 {
			return uint8(i)
		}
	}
	return 33
}

func (c *Chord) FirstIntervalNotDoubled() uint8 {
	for i := range uint8(32) {
		if c.Notes&(1<<i) != 0 {
			if !ContainsBits(c.Notes, 1<<(i+12)) {
				return i
			}
		}
	}
	return 33
}

func (c *Chord) LastIntervalDoubled() uint8 {
	for i := uint8(31); i >= 0; i-- {
		if c.Notes&(1<<i) != 0 {
			if ContainsBits(c.Notes, 1<<(i-12)) {
				return i
			}
		}
	}
	return 33
}

func (c Chord) NamedIntervals() []string {
	notes := make([]string, 0)
	for i := 31; i >= 0; i-- {
		if c.Notes&(1<<i) != 0 {
			notes = append(notes, interval(i))
		}
	}
	return notes
}

func interval(n int) string {
	switch n {
	case 0:
		return "P1"
	case 1:
		return "m2"
	case 2:
		return "M2"
	case 3:
		return "m3"
	case 4:
		return "M3"
	case 5:
		return "P4"
	case 6:
		return "d5"
	case 7:
		return "P5"
	case 8:
		return "m6"
	case 9:
		return "M6"
	case 10:
		return "m7"
	case 11:
		return "M7"
	case 12:
		return "P8"
	case 13:
		return "m9"
	case 14:
		return "M9"
	case 15:
		return "m10"
	case 16:
		return "M10"
	case 17:
		return "P11"
	case 18:
		return "d12"
	case 19:
		return "P12"
	case 20:
		return "m13"
	case 21:
		return "M13"
	case 22:
		return "m14"
	case 23:
		return "M14"
	case 24:
		return "P15"
	}
	return ""
}

// Intervals returns a slice of integers representing the notes in the chord
// If the chord has an inversion value, the appropriate number of notes
// from the bottom of the chord are moved to the top
func (c Chord) Intervals() []uint8 {
	notes := c.UninvertedNotes()

	return notes
}

func ContainsBits(source, pattern uint32) bool {
	return (source & pattern) == pattern
}

func (c *Chord) NextInversion() {
	firstInterval := c.FirstInterval()
	firstNoteBit := uint32(1 << firstInterval)
	upOctaveBit := firstNoteBit << 12
	if upOctaveBit < (1 << 31) {
		c.AddNotes(upOctaveBit)
		c.OmitInterval(firstInterval)
	}
}

func (c *Chord) PreviousInversion() {
	lastInterval := c.LastInterval()
	lastNoteBit := uint32(1 << lastInterval)
	downOctaveBit := lastNoteBit >> 12
	if downOctaveBit > 0 {
		c.AddNotes(downOctaveBit)
		c.OmitInterval(lastInterval)
	}
}

func (c Chord) Inversions() int {
	// Count the number of notes in the chord above 12 semitones
	count := 0
	for i := uint8(12); i < 32; i++ {
		if c.Notes&(1<<i) != 0 {
			count++
		}
	}
	return count
}

func (c *Chord) NextDouble() {
	firstInterval := c.FirstIntervalNotDoubled()
	firstNoteBit := uint32(1 << firstInterval)
	upOctaveBit := firstNoteBit << 12
	if upOctaveBit < (1 << 31) {
		c.AddNotes(upOctaveBit)
	}
}

func (c *Chord) PreviousDouble() {
	lastInterval := c.LastIntervalDoubled()
	if lastInterval >= 12 && lastInterval < 32 {
		c.OmitInterval(lastInterval)
	}
}
