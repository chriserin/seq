package overlays

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/theory"
	"github.com/stretchr/testify/assert"
)

func TestDiffOverlays(t *testing.T) {
	t.Run("diffing identical overlays should show no differences", func(t *testing.T) {
		// Create a test overlay
		key := overlaykey.InitOverlayKey(2, 1)
		original := InitOverlay(key, nil)

		// Add a note
		note := grid.InitNote()
		gridKey := grid.GridKey{Line: 1, Beat: 2}
		original.MoveNoteTo(gridKey, note)

		// Add a chord
		chordKey := grid.GridKey{Line: 3, Beat: 0}
		original.CreateChord(chordKey, 0)

		// Create a copy (deep copy)
		modified := InitOverlay(key, nil)
		modified.MoveNoteTo(gridKey, note)
		modified.CreateChord(chordKey, 0)

		// Diff the overlays
		diff := DiffOverlays(original, modified)

		// Verify no differences
		assert.Empty(t, diff.AddedNotes)
		assert.Empty(t, diff.RemovedNotes)
		assert.Empty(t, diff.ModifiedNotes)
		assert.Empty(t, diff.AddedChords)
		assert.Empty(t, diff.RemovedChords)
		assert.Empty(t, diff.ModifiedChords)
		assert.False(t, diff.OptionsDiff.PressUpChanged)
		assert.False(t, diff.OptionsDiff.PressDownChanged)
	})

	t.Run("should detect added, removed, and modified notes", func(t *testing.T) {
		// Create original overlay
		key := overlaykey.InitOverlayKey(2, 1)
		original := InitOverlay(key, nil)

		// Add some notes to original
		note1 := grid.InitNote()
		note2 := grid.InitNote()
		note3 := grid.InitNote()
		original.MoveNoteTo(grid.GridKey{Line: 1, Beat: 2}, note1)
		original.MoveNoteTo(grid.GridKey{Line: 3, Beat: 4}, note2)
		original.MoveNoteTo(grid.GridKey{Line: 5, Beat: 6}, note3)

		// Create modified overlay
		modified := InitOverlay(key, nil)

		// Keep note1
		modified.MoveNoteTo(grid.GridKey{Line: 1, Beat: 2}, note1)

		// Modify note2
		modifiedNote2 := grid.InitNote()
		modifiedNote2.AccentIndex = 2 // Change accent index
		modified.MoveNoteTo(grid.GridKey{Line: 3, Beat: 4}, modifiedNote2)

		// Remove note3 (don't add it to modified)

		// Add a new note
		note4 := grid.InitNote()
		modified.MoveNoteTo(grid.GridKey{Line: 7, Beat: 8}, note4)

		// Diff the overlays
		diff := DiffOverlays(original, modified)

		// Verify differences
		assert.Len(t, diff.AddedNotes, 1)
		assert.Len(t, diff.RemovedNotes, 1)
		assert.Len(t, diff.ModifiedNotes, 1)

		// Check specific differences
		assert.Contains(t, diff.AddedNotes, grid.GridKey{Line: 7, Beat: 8})
		assert.Contains(t, diff.RemovedNotes, grid.GridKey{Line: 5, Beat: 6})
		assert.Contains(t, diff.ModifiedNotes, grid.GridKey{Line: 3, Beat: 4})
	})

	t.Run("should detect added, removed, and modified chords", func(t *testing.T) {
		// Create original overlay
		key := overlaykey.InitOverlayKey(2, 1)
		original := InitOverlay(key, nil)

		// Add chords to original
		chord1Key := grid.GridKey{Line: 1, Beat: 0}
		chord2Key := grid.GridKey{Line: 3, Beat: 0}
		original.CreateChord(chord1Key, 0)
		original.CreateChord(chord2Key, 0)

		// Create modified overlay
		modified := InitOverlay(key, nil)

		// Keep chord1
		modified.CreateChord(chord1Key, 0)

		// Modify chord2 - create with different alteration
		modified.CreateChord(chord2Key, 1)

		// Add a new chord
		chord3Key := grid.GridKey{Line: 5, Beat: 0}
		modified.CreateChord(chord3Key, 0)

		// Diff the overlays
		diff := DiffOverlays(original, modified)

		// Verify differences - for chords we expect:
		// - Chord1 is unchanged (not in any difference lists)
		// - Chord2 is modified (different alteration)
		// - Chord3 is added
		assert.Len(t, diff.AddedChords, 1)
		assert.Empty(t, diff.RemovedChords) // None removed, one was modified
		assert.Len(t, diff.ModifiedChords, 1)

		// The added chord should be at position chord3Key
		for _, chord := range diff.AddedChords {
			assert.Equal(t, chord3Key, chord.Root)
		}
	})

	t.Run("should detect moved chords", func(t *testing.T) {
		// Create original overlay
		key := overlaykey.InitOverlayKey(2, 1)
		original := InitOverlay(key, nil)

		// Add chords to original
		chord1Key := grid.GridKey{Line: 1, Beat: 2}
		chord2Key := grid.GridKey{Line: 1, Beat: 1}
		original.CreateChord(chord1Key, 0)

		// Create modified overlay
		modified := InitOverlay(key, nil)

		// Keep chord1
		modified.CreateChord(chord2Key, 0)

		// Diff the overlays
		diff := DiffOverlays(original, modified)

		assert.Len(t, diff.AddedChords, 1)
		assert.Len(t, diff.RemovedChords, 1)

		diff.Apply(original)

		assert.Len(t, modified.Chords, 1)
	})

	t.Run("should detect changes in overlay options", func(t *testing.T) {
		// Create original overlay
		key := overlaykey.InitOverlayKey(2, 1)
		original := InitOverlay(key, nil)
		original.PressUp = true
		original.PressDown = false

		// Create modified overlay with different options
		modified := InitOverlay(key, nil)
		modified.PressUp = false
		modified.PressDown = true

		// Diff the overlays
		diff := DiffOverlays(original, modified)

		// Verify differences
		assert.True(t, diff.OptionsDiff.PressUpChanged)
		assert.True(t, diff.OptionsDiff.PressDownChanged)
	})

	t.Run("Apply function should transform an overlay according to the diff", func(t *testing.T) {
		// Create original overlay
		key := overlaykey.InitOverlayKey(2, 1)
		original := InitOverlay(key, nil)

		// Add a note
		note := grid.InitNote()
		gridKey := grid.GridKey{Line: 1, Beat: 2}
		original.SetNote(gridKey, note)

		another := DeepCopy(original)

		// Add a different note
		newNote := grid.InitNote()
		newNote.AccentIndex = 5
		newGridKey := grid.GridKey{Line: 3, Beat: 4}
		another.SetNote(newGridKey, newNote)

		// Generate diff
		diff := DiffOverlays(original, another)

		// Apply the diff
		diff.Apply(original)

		comparedDiff := DiffOverlays(original, another)

		// Verify target now matches modified
		assert.Equal(t, comparedDiff, InitDiff())
	})

	t.Run("Apply function should transform an overlay with modified chord according to the diff", func(t *testing.T) {
		// Create original overlay
		key := overlaykey.InitOverlayKey(2, 1)
		original := InitOverlay(key, nil)

		// Add a chord
		gridKey := grid.GridKey{Line: 1, Beat: 2}
		original.CreateChord(gridKey, theory.MajorTriad)

		another := DeepCopy(original)

		newNote := grid.InitNote()
		newNote.AccentIndex = 7
		newGridKey := grid.GridKey{Line: 1, Beat: 2}
		another.SetNote(newGridKey, newNote)

		// Generate diff
		diff := DiffOverlays(original, another)

		assert.Len(t, diff.ModifiedChords, 1)

		assert.Equal(t, uint8(5), original.Chords[0].Notes[0].Note.AccentIndex)

		// Apply the diff
		diff.Apply(original)

		assert.Equal(t, uint8(7), original.Chords[0].Notes[0].Note.AccentIndex)

		comparedDiff := DiffOverlays(original, another)

		// Verify target now matches modified
		assert.Equal(t, InitDiff(), comparedDiff)
	})
}

func TestDeepCopy(t *testing.T) {
	// Test with a nil overlay
	t.Run("Nil overlay", func(t *testing.T) {
		var overlay *Overlay
		clone := DeepCopy(overlay)
		assert.Nil(t, clone)
	})

	// Test with an empty overlay
	t.Run("Empty overlay", func(t *testing.T) {
		key := overlaykey.ROOT
		original := InitOverlay(key, nil)

		clone := DeepCopy(original)

		// Check that it's a different pointer
		assert.NotSame(t, original, clone)

		// Check that the values are equal
		assert.Equal(t, original.Key, clone.Key)
		assert.Equal(t, len(original.Notes), len(clone.Notes))
		assert.Equal(t, len(original.Chords), len(clone.Chords))
		assert.Equal(t, original.PressUp, clone.PressUp)
		assert.Equal(t, original.PressDown, clone.PressDown)
		assert.Nil(t, clone.Below)
	})

	// Test with an overlay containing notes
	t.Run("Overlay with notes", func(t *testing.T) {
		key := overlaykey.ROOT
		original := InitOverlay(key, nil)

		// Add some notes
		note1 := grid.InitNote()
		note1.AccentIndex = 1
		gridKey1 := grid.GridKey{Beat: 1, Line: 1}
		original.MoveNoteTo(gridKey1, note1)

		note2 := grid.InitNote()
		note2.AccentIndex = 2
		gridKey2 := grid.GridKey{Beat: 2, Line: 2}
		original.MoveNoteTo(gridKey2, note2)

		clone := DeepCopy(original)

		// Check that the notes are equal but not the same map
		assert.Equal(t, len(original.Notes), len(clone.Notes))
		assert.Equal(t, original.Notes[gridKey1], clone.Notes[gridKey1])
		assert.Equal(t, original.Notes[gridKey2], clone.Notes[gridKey2])

		// Modify the original notes and verify clone is unchanged
		note1.AccentIndex = 3
		original.MoveNoteTo(gridKey1, note1)
		assert.NotEqual(t, original.Notes[gridKey1], clone.Notes[gridKey1])
	})

	// Test with an overlay containing chords
	t.Run("Overlay with chords", func(t *testing.T) {
		key := overlaykey.ROOT
		original := InitOverlay(key, nil)

		// Add a chord
		root := grid.GridKey{Beat: 4, Line: 4}
		original.CreateChord(root, 0) // Chord with default alteration

		// Verify the chord was added
		assert.Greater(t, len(original.Chords), 0, "Chord should have been added")

		// Get a reference to the chord
		chord := original.Chords[0]

		// Modify the chord
		chord.Double = 2
		chord.Arppegio = ARP_UP

		clone := DeepCopy(original)

		// Check that chords are equal but not the same references
		assert.Equal(t, len(original.Chords), len(clone.Chords))
		assert.NotSame(t, original.Chords[0], clone.Chords[0])
		assert.Equal(t, original.Chords[0].Root, clone.Chords[0].Root)
		assert.Equal(t, original.Chords[0].Double, clone.Chords[0].Double)
		assert.Equal(t, original.Chords[0].Arppegio, clone.Chords[0].Arppegio)

		// Modify original chord and verify clone is unchanged
		original.Chords[0].Double = 3
		assert.NotEqual(t, original.Chords[0].Double, clone.Chords[0].Double)
	})

	// Test with blockers
	t.Run("Overlay with blockers", func(t *testing.T) {
		key := overlaykey.ROOT
		original := InitOverlay(key, nil)

		// Add a blocker
		blocker := grid.GridKey{Beat: 3, Line: 3}
		original.blockers = append(original.blockers, blocker)

		clone := DeepCopy(original)

		// Check that blockers are copied
		assert.Equal(t, len(original.blockers), len(clone.blockers))
		assert.Equal(t, original.blockers[0], clone.blockers[0])

		// Modify original and check that clone is unchanged
		original.blockers[0] = grid.GridKey{Beat: 5, Line: 5}
		assert.NotEqual(t, original.blockers[0], clone.blockers[0])
	})
}
