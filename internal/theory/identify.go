package theory

import (
	"slices"
)

// FromNotesWithAnalysis creates a chord from a slice of integers, analyzing the notes
// to determine the most likely chord type based on common chord structures
func FromNotesWithAnalysis(notes []uint8) Chord {
	intNotes := make([]int, len(notes))
	for i, n := range notes {
		intNotes[i] = int(n)
	}
	normalizedNotes := shiftToZero(intNotes)

	bestChord := identifyChord(normalizedNotes)

	return bestChord
}

func Rotate(s []int, x int) []int {
	shiftedNotes := s
	for range x {
		n := shiftedNotes[len(shiftedNotes)-1]
		rotatedNote := n - (12 * ((n + 12) / 12))
		shiftedNotes = append([]int{rotatedNote}, shiftedNotes[:len(shiftedNotes)-1]...)
		shiftedNotes = shiftToZero(shiftedNotes)
	}
	return shiftedNotes
}

func MoveToFront(s []int, x int) []int {
	// fmt.Println("Input", movedNotes)
	movedNotes := make([]int, len(s))
	copy(movedNotes, s)
	index := len(s) - (1 + x)
	newFrontNote := s[index]
	movedNotes = slices.Delete(movedNotes, index, index+1)
	movedNotes = append([]int{newFrontNote - 12}, movedNotes...)
	movedNotes = shiftToZero(movedNotes)
	slices.Sort(movedNotes)
	// fmt.Println("Output", movedNotes)
	return movedNotes
}

func shiftToZero(notes []int) []int {
	slices.Sort(notes)
	firstNote := notes[0]
	normalized := make([]int, len(notes))
	for i, n := range notes {
		normalized[i] = n - firstNote
	}

	return normalized
}

// identifyChord analyzes normalized notes and returns the bit pattern
// representing the most likely chord
func identifyChord(normalizedNotes []int) Chord {
	// One note is just a single note
	if len(normalizedNotes) == 1 {
		return Chord{Notes: uint32(1 << normalizedNotes[0])}
	}

	return GetBestScoreChord(normalizedNotes, 0)
}

func InitPossibleChord(notes []int, inversions int) Chord {
	return Chord{
		Notes:     FromNotes(notes),
		Inversion: int8(inversions),
		score:     scoreChord(notes),
	}
}

// FromNotes creates a chord from a slice of integers representing notes in the chromatic scale
func FromNotes(notes []int) uint32 {
	var chordBits uint32

	for _, note := range notes {
		if note >= 0 && note < 32 {
			chordBits |= uint32(1 << note)
		}
	}

	return chordBits
}

func GetBestScoreChord(notes []int, reversions int) Chord {
	var bestChord = InitPossibleChord(notes, reversions)
	if len(notes) <= reversions {
		return bestChord
	}

	revertedNotes := Rotate(notes, 1)
	newChord := GetBestScoreChord(revertedNotes, reversions+1)

	moveToFrontChord := CycleMoveToFront(notes, reversions)

	if newChord.Score() > bestChord.Score() && newChord.Score() > moveToFrontChord.Score() {
		return newChord
	} else if bestChord.Score() > moveToFrontChord.Score() {
		return bestChord
	} else {
		return moveToFrontChord
	}
}

func CycleMoveToFront(notes []int, reversions int) Chord {
	var bestChord Chord
	for i := 0; i < len(notes)-1; i++ {
		moveToFrontNotes := MoveToFront(notes, i+1)
		chord := GetBestScoreChord(moveToFrontNotes, reversions+1)
		if chord.Score() > bestChord.Score() {
			bestChord = chord
		}
	}

	return bestChord
}

func GetIntervals(notes []int) []int {
	var intervals = make([]int, len(notes)-1)

	for i := 0; i < len(notes)-1; i++ {
		intervals[i] = notes[i+1] - notes[i]
	}
	return intervals
}

func ScoreIntervals(notes []int) int {

	var score = 0
	var intervals = GetIntervals(notes)
	if slices.Contains(intervals, 1) {
		score -= 10
	}

	if slices.Contains(intervals, 2) {
		for _, v := range intervals {
			if v == 2 {
				score -= 2
			}
		}
	}
	return score
}

// scoreChord evaluates a set of intervals assuming 0 is the root
// and returns a chord bit pattern and a score indicating how well it matches common patterns
func scoreChord(notes []int) int {
	var score int

	// Major chord: root, major third (4), perfect fifth (7)
	if containsAllIntervals(notes, []int{0, 4, 7}) {
		score = 100 // Base score for major triad

		// Check for extensions
		if contains(notes, 11) { // Major 7th
			score += 5 // Common extension
		} else if contains(notes, 10) { // Minor 7th
			score += 10 // Very common extension
		}

		if contains(notes, 9) { // 6th
			score += 5
		}

		if contains(notes, 2, 14) { // 9th (actually 2nd in this octave)
			score += 5
		}

		if contains(notes, 5, 17) { // 11th (actually 4th in this octave)
			score += 5 // More common in minor chords
		}

		if contains(notes, 9, 21) { // 13th (actually 6th in this octave)
			score += 3
		}

		if contains(notes, 8) { // Augmented 5th
			score -= 10 // Less common alteration
		} else if contains(notes, 6) { // Diminished 5th
			score -= 10 // Less common alteration
		}

		return score + ScoreIntervals(notes)
	}

	// Minor chord: root, minor third (3), perfect fifth (7)
	if containsAllIntervals(notes, []int{0, 3, 7}) {
		score = 100 // Base score for minor triad

		// Check for extensions
		if contains(notes, 11) { // Major 7th
			score += 5 // Common extension
		} else if contains(notes, 10) { // Minor 7th
			score += 10 // Very common extension
		}

		if contains(notes, 9) { // 6th
			score += 3
		}

		if contains(notes, 2, 14) { // 9th (actually 2nd in this octave)
			score += 5
		}

		if contains(notes, 5, 17) { // 11th (actually 4th in this octave)
			score += 5 // More common in minor chords
		}

		if contains(notes, 9, 21) { // 13th (actually 6th in this octave)
			score += 3
		}

		if contains(notes, 8) { // Augmented 5th
			score -= 10 // Less common alteration
		} else if contains(notes, 6) { // Diminished 5th
			score -= 5 // Less uncommon in minor contexts
		}

		return score + ScoreIntervals(notes)
	}

	// Diminished chord: root, minor third (3), diminished fifth (6)
	if containsAllIntervals(notes, []int{0, 3, 6}) {
		score = 90 // Fairly common

		// Check for diminished seventh
		if contains(notes, 9) { // dim7 (enharmonic with Major6)
			score += 10
		}

		return score + ScoreIntervals(notes)
	}

	// Augmented chord: root, major third (4), augmented fifth (8)
	if containsAllIntervals(notes, []int{0, 4, 8}) {
		score = 80 // Less common
		return score + ScoreIntervals(notes)
	}

	// Sus4 chord: root, perfect fourth (5), perfect fifth (7)
	if containsAllIntervals(notes, []int{0, 5, 7}) {
		score = 85

		// Check for 7sus4
		if contains(notes, 10) {
			score += 10
		}

		return score + ScoreIntervals(notes)
	}

	// Sus2 chord: root, major second (2), perfect fifth (7)
	if containsAllIntervals(notes, []int{0, 2, 7}) {
		score = 85
		return score + ScoreIntervals(notes)
	}

	// Power chord: root and fifth
	if containsAllIntervals(notes, []int{0, 7}) && len(notes) == 2 {
		score = 80
		return score + ScoreIntervals(notes)
	}

	// Generic case - score based on presence of common intervals
	if contains(notes, 0) {
		score += 30
	} // Root
	if contains(notes, 7) {
		score += 25
	} // Fifth
	if contains(notes, 4) {
		score += 20
	} // Major third
	if contains(notes, 3) {
		score += 20
	} // Minor third
	if contains(notes, 10) {
		score += 10
	} // Minor seventh
	if contains(notes, 11) {
		score += 5
	} // Major seventh

	return score
}

// contains checks if a slice contains a specific value
func contains(slice []int, val ...int) bool {
	for _, v := range val {
		if slices.Contains(slice, v) {
			return true
		}
	}
	return false
}

// containsAllIntervals checks if all specified intervals are present
func containsAllIntervals(actual []int, expected []int) bool {
	for _, interval := range expected {
		if !contains(actual, interval) {
			return false
		}
	}
	return true
}

// GetChordComponents returns a slice of constants that can be used to construct the chord
func (c Chord) GetChordComponents() []uint32 {
	components := make([]uint32, 0)

	// Check for triads first
	if (c.Notes & MinorTriad) == MinorTriad {
		components = append(components, MinorTriad)
	} else if (c.Notes & MajorTriad) == MajorTriad {
		components = append(components, MajorTriad)
	} else if (c.Notes & AugmentedTriad) == AugmentedTriad {
		components = append(components, AugmentedTriad)
	} else if (c.Notes & DiminishedTriad) == DiminishedTriad {
		components = append(components, DiminishedTriad)
	} else {
		// Add individual notes if no specific triad is found
		if (c.Notes & Root) != 0 {
			components = append(components, Root)
		}
		if (c.Notes & Minor3) != 0 {
			components = append(components, Minor3)
		} else if (c.Notes & Major3) != 0 {
			components = append(components, Major3)
		}
		if (c.Notes & Dim5) != 0 {
			components = append(components, Dim5)
		} else if (c.Notes & Perfect5) != 0 {
			components = append(components, Perfect5)
		} else if (c.Notes & Aug5) != 0 {
			components = append(components, Aug5)
		}
	}

	// Check for extensions
	if (c.Notes & MinorSeventh) != 0 {
		components = append(components, MinorSeventh)
	} else if (c.Notes & MajorSeventh) != 0 {
		components = append(components, MajorSeventh)
	}

	// Check for 9th, 11th, 13th extensions
	if (c.Notes & Extend91113) == Extend91113 {
		components = append(components, Extend91113)
	} else if (c.Notes & Extend911) == Extend911 {
		components = append(components, Extend911)
	} else if (c.Notes & Extend9) == Extend9 {
		components = append(components, Extend9)
	}

	// Check for alterations
	if (c.Notes & Aug11) != 0 {
		components = append(components, Aug11)
	} else if (c.Notes & Dim11) != 0 {
		components = append(components, Dim11)
	}

	if (c.Notes & Aug13) != 0 {
		components = append(components, Aug13)
	} else if (c.Notes & Dim13) != 0 {
		components = append(components, Dim13)
	}

	return components
}

// GetChordName returns a string representation of the chord
func (c Chord) GetChordName() string {
	components := c.GetChordComponents()
	if len(components) == 0 {
		return "Empty Chord"
	}

	var name string

	// Start with the triad type
	for _, comp := range components {
		switch comp {
		case MinorTriad:
			name = "Minor"
		case MajorTriad:
			name = "Major"
		case AugmentedTriad:
			name = "Augmented"
		case DiminishedTriad:
			name = "Diminished"
		}
		if name != "" {
			break
		}
	}

	// If no standard triad, construct from individual notes
	if name == "" {
		for _, comp := range components {
			if comp == Root {
				name = "Root"
			}
		}

		for _, comp := range components {
			switch comp {
			case Minor3:
				name += " Minor"
				break
			case Major3:
				name += " Major"
				break
			}
		}

		for _, comp := range components {
			switch comp {
			case Dim5:
				name += " Diminished 5th"
				break
			case Perfect5:
				break // Perfect 5th is implied
			case Aug5:
				name += " Augmented 5th"
				break
			}
		}
	}

	// Add seventh chord quality
	for _, comp := range components {
		switch comp {
		case MinorSeventh:
			name += "7"
			break
		case MajorSeventh:
			name += "Maj7"
			break
		}
	}

	// Add extensions
	for _, comp := range components {
		switch comp {
		case Extend9:
			name += "9"
			break
		case Extend911:
			name += "11"
			break
		case Extend91113:
			name += "13"
			break
		}
	}

	// Add altered notes
	for _, comp := range components {
		switch comp {
		case Aug11:
			name += " #11"
			break
		case Dim11:
			name += " b11"
			break
		case Aug13:
			name += " #13"
			break
		case Dim13:
			name += " b13"
			break
		}
	}

	return name
}
