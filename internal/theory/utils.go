package theory

import "slices"

// shiftToZero reduces all notes to within a single octave (0-11)
// and returns a sorted, unique list
func ShiftToZero(notes []uint8) []uint8 {
	slices.Sort(notes)
	firstNote := notes[0]
	normalized := make([]uint8, len(notes))
	for i, n := range notes {
		normalized[i] = n - firstNote
	}

	return normalized
}
