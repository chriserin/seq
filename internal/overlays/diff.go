package overlays

import (
	"fmt"
	"maps"

	"github.com/chriserin/sq/internal/grid"
	"github.com/chriserin/sq/internal/theory"
)

// OverlayDiff represents the differences between two overlays
type OverlayDiff struct {
	// Notes
	AddedNotes    grid.Pattern
	RemovedNotes  grid.Pattern
	ModifiedNotes map[grid.GridKey]NoteDiff
	// Chords
	AddedChords    []GridChord
	RemovedChords  []GridChord
	ModifiedChords map[*GridChord]ChordDiff
	// Options
	OptionsDiff OptionsDiff
}

func (od OverlayDiff) IsEmpty() bool {
	return (len(od.AddedNotes)+
		len(od.RemovedNotes)+
		len(od.ModifiedNotes)+
		len(od.AddedChords)+
		len(od.RemovedChords)+
		len(od.ModifiedChords)) == 0 &&
		!od.OptionsDiff.PressDownChanged &&
		!od.OptionsDiff.PressUpChanged
}

// NoteDiff represents the difference between two notes
type NoteDiff struct {
	OldNote grid.Note
	NewNote grid.Note
}

// ChordDiff represents the difference between two chords
type ChordDiff struct {
	OldChord GridChord
	Modified bool
}

// OptionsDiff represents differences in overlay options
type OptionsDiff struct {
	PressUpChanged   bool
	PressDownChanged bool
}

func InitDiff() OverlayDiff {
	return OverlayDiff{
		AddedNotes:     make(grid.Pattern),
		RemovedNotes:   make(grid.Pattern),
		ModifiedNotes:  make(map[grid.GridKey]NoteDiff),
		ModifiedChords: make(map[*GridChord]ChordDiff),
	}
}

// DiffOverlays compares two overlays and returns the differences between them
func DiffOverlays(original, modified *Overlay) OverlayDiff {
	diff := InitDiff()

	// Diff notes
	diffNotes(original, modified, &diff)

	// Diff chords
	diffChords(original, modified, &diff)

	// Diff options
	diff.OptionsDiff = OptionsDiff{
		PressUpChanged:   original.PressUp != modified.PressUp,
		PressDownChanged: original.PressDown != modified.PressDown,
	}

	return diff
}

// diffNotes compares the Notes in two overlays and updates the diff
func diffNotes(original, modified *Overlay, diff *OverlayDiff) {
	// Check for removed or modified notes
	for gridKey, originalNote := range original.Notes {
		if modifiedNote, exists := modified.Notes[gridKey]; exists {
			// Note exists in both - check if it's different
			if !noteEqual(originalNote, modifiedNote) {
				diff.ModifiedNotes[gridKey] = NoteDiff{
					OldNote: originalNote,
					NewNote: modifiedNote,
				}
			}
		} else {
			// Note exists in original but not in modified - it was removed
			diff.RemovedNotes[gridKey] = originalNote
		}
	}

	// Check for added notes
	for gridKey, modifiedNote := range modified.Notes {
		if _, exists := original.Notes[gridKey]; !exists {
			// Note exists in modified but not in original - it was added
			diff.AddedNotes[gridKey] = modifiedNote
		}
	}
}

// diffChords compares the Chords in two overlays and updates the diff
func diffChords(original, modified *Overlay, diff *OverlayDiff) {
	// Track which chords we've already processed
	processedOriginalChords := make(map[*GridChord]bool)

	// Check for modified or removed chords
	for _, originalChord := range original.Chords {
		// Try to find a matching chord in the modified overlay
		matched := false

		for _, modifiedChord := range modified.Chords {
			if originalChord.Root == modifiedChord.Root {
				// We found a match - check for differences
				chordDiff := compareChords(originalChord, modifiedChord)
				if chordDiff.hasDifferences() {
					diff.ModifiedChords[originalChord] = chordDiff
				}

				processedOriginalChords[originalChord] = true
				matched = true
				break
			}
		}

		if !matched {
			// Chord exists in original but not in modified - it was removed
			diff.RemovedChords = append(diff.RemovedChords, originalChord.DeepCopy())
		}
	}

	// Check for added chords
	for _, modifiedChord := range modified.Chords {
		found := false
		for _, originalChord := range original.Chords {
			if originalChord.Root == modifiedChord.Root {
				found = true
				break
			}
		}

		if !found {
			// Chord exists in modified but not in original - it was added
			diff.AddedChords = append(diff.AddedChords, modifiedChord.DeepCopy())
		}
	}
}

// noteEqual returns true if two notes are equivalent
func noteEqual(a, b grid.Note) bool {
	return a.AccentIndex == b.AccentIndex &&
		a.Ratchets.Hits == b.Ratchets.Hits &&
		a.Ratchets.Length == b.Ratchets.Length &&
		a.Ratchets.Span == b.Ratchets.Span &&
		a.Action == b.Action &&
		a.GateIndex == b.GateIndex &&
		a.WaitIndex == b.WaitIndex
}

func chordsMatch(a *GridChord, b GridChord) bool {
	// For simplicity, we'll consider chords the same if they have the same root position
	return a.Root == b.Root
}

func compareChords(original, modified *GridChord) ChordDiff {
	diff := ChordDiff{
		OldChord: modified.DeepCopy(),
		Modified: original.Arpeggio != modified.Arpeggio ||
			!alterationsEqual(original.Chord, modified.Chord) ||
			!notesMatch(original.Notes, modified.Notes),
	}

	return diff
}

func notesMatch(a, b []BeatNote) bool {
	if len(a) != len(b) {
		return false
	}

	for i, an := range a {
		if an != b[i] {
			return false
		}
	}

	return true
}

// alterationsEqual checks if two chords have the same alterations
func alterationsEqual(a, b theory.Chord) bool {
	return a.Notes == b.Notes
}

// hasDifferences returns true if the ChordDiff contains actual differences
func (cd ChordDiff) hasDifferences() bool {
	return cd.Modified
}

// String returns a string representation of the OverlayDiff
func (od OverlayDiff) String() string {
	result := "Overlay Diff:\n"

	if len(od.AddedNotes) > 0 {
		result += fmt.Sprintf("  Added Notes: %d\n", len(od.AddedNotes))
	}

	if len(od.RemovedNotes) > 0 {
		result += fmt.Sprintf("  Removed Notes: %d\n", len(od.RemovedNotes))
	}

	if len(od.ModifiedNotes) > 0 {
		result += fmt.Sprintf("  Modified Notes: %d\n", len(od.ModifiedNotes))
	}

	if len(od.AddedChords) > 0 {
		result += fmt.Sprintf("  Added Chords: %d\n", len(od.AddedChords))
	}

	if len(od.RemovedChords) > 0 {
		result += fmt.Sprintf("  Removed Chords: %d\n", len(od.RemovedChords))
	}

	if len(od.ModifiedChords) > 0 {
		result += fmt.Sprintf("  Modified Chords: %d\n", len(od.ModifiedChords))
	}

	if od.OptionsDiff.PressUpChanged || od.OptionsDiff.PressDownChanged {
		result += "  Options Changed\n"
	}

	return result
}

// Apply applies the diff to an overlay, transforming it to match the modified overlay
func (od OverlayDiff) Apply(overlay *Overlay) {
	// Apply note changes
	for gridKey := range od.RemovedNotes {
		overlay.RemoveNote(gridKey)
	}

	for gridKey, note := range od.AddedNotes {
		overlay.MoveNoteTo(gridKey, note)
	}

	for gridKey, noteDiff := range od.ModifiedNotes {
		overlay.MoveNoteTo(gridKey, noteDiff.NewNote)
	}

	// Apply chord changes
	for _, chord := range od.RemovedChords {
		// Find and remove the chord
		for i, c := range overlay.Chords {
			if chordsMatch(c, chord) {
				overlay.Chords = append(overlay.Chords[:i], overlay.Chords[i+1:]...)
				break
			}
		}
	}

	for _, chord := range od.AddedChords {
		// Deep copy the chord
		newChord := chord
		notes := make([]BeatNote, len(chord.Notes))
		copy(notes, chord.Notes)
		newChord.Notes = notes

		overlay.Chords = append(overlay.Chords, &newChord)
	}

	for _, chordDiff := range od.ModifiedChords {
		// Find the chord in the current overlay
		for i, c := range overlay.Chords {
			if c.Root == chordDiff.OldChord.Root {
				remainder := overlay.Chords[i+1:]
				overlay.Chords = append(overlay.Chords[:i], &chordDiff.OldChord)
				overlay.Chords = append(overlay.Chords, remainder...)

				break
			}
		}
	}

	// Apply options changes
	if od.OptionsDiff.PressUpChanged {
		overlay.PressUp = !overlay.PressUp
	}

	if od.OptionsDiff.PressDownChanged {
		overlay.PressDown = !overlay.PressDown
	}
}

// DeepCopy creates a complete deep copy of the Overlay struct
func DeepCopy(ol *Overlay) *Overlay {
	if ol == nil {
		return nil
	}

	// Create a new Overlay struct
	clone := &Overlay{
		Key:       ol.Key,
		Notes:     make(grid.Pattern),
		Chords:    make([]*GridChord, 0, len(ol.Chords)),
		Blockers:  make([]*GridChord, len(ol.Blockers)),
		PressUp:   ol.PressUp,
		PressDown: ol.PressDown,
	}

	// Deep copy the Notes map
	maps.Copy(clone.Notes, ol.Notes)

	// Deep copy the Chords slice
	for _, chord := range ol.Chords {
		chordCopy := &GridChord{
			Chord:    chord.Chord,
			Root:     chord.Root,
			Arpeggio: chord.Arpeggio,
		}

		// Deep copy the Notes slice in GridChord
		chordCopy.Notes = make([]BeatNote, len(chord.Notes))
		copy(chordCopy.Notes, chord.Notes)

		clone.Chords = append(clone.Chords, chordCopy)
	}

	// Deep copy the blockers slice
	for i, blocker := range ol.Blockers {
		if blocker != nil {
			blockerCopy := &GridChord{
				Chord:    blocker.Chord,
				Root:     blocker.Root,
				Arpeggio: blocker.Arpeggio,
			}

			// Deep copy the Notes slice in GridChord
			blockerCopy.Notes = make([]BeatNote, len(blocker.Notes))
			copy(blockerCopy.Notes, blocker.Notes)

			clone.Blockers[i] = blockerCopy
		}
	}

	return clone
}
