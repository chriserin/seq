package main

import (
	"testing"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
	"github.com/stretchr/testify/assert"
)

func setupTestModel() *model {
	var parts = InitParts()
	var arr = InitArrangement(parts)
	def := Definition{
		arrangement:  arr,
		parts:        &parts,
		keyline:      0,
		tempo:        120,
		subdivisions: 4,
		lines:        make([]grid.LineDefinition, 1),
	}

	m := &model{
		definition:     def,
		programChannel: make(chan midiEventLoopMsg),
		playState:      InitLineStates(1, []linestate{}, 0),
		undoStack:      NilStack,
		redoStack:      NilStack,
		arrangement:    arrangement.InitModel(def.arrangement, def.parts),
	}

	m.currentOverlay = parts[0].Overlays
	return m
}

func TestUndoStack(t *testing.T) {
	t.Run("Push and pop undoables", func(t *testing.T) {
		m := setupTestModel()

		// Create test undoables with the same return signature as required by Undoable interface
		undo := UndoGridNote{
			overlayKey:     overlaykey.InitOverlayKey(1, 1),
			cursorPosition: GK(0, 0),
			gridNote:       GridNote{gridKey: GK(0, 0), note: note{Action: 0}},
		}
		redo := UndoGridNote{
			overlayKey:     overlaykey.InitOverlayKey(1, 1),
			cursorPosition: GK(0, 0),
			gridNote:       GridNote{gridKey: GK(0, 0), note: note{Action: 0}},
		}

		// Push undoables to stack
		m.PushUndoables(undo, redo)

		// Verify stack structure
		assert.NotEqual(t, NilStack, m.undoStack)
		assert.Equal(t, undo, m.undoStack.undo)
		assert.Equal(t, redo, m.undoStack.redo)
		assert.Nil(t, m.undoStack.next)

		// Push another pair
		undo2 := UndoGridNote{
			overlayKey:     overlaykey.InitOverlayKey(1, 1),
			cursorPosition: GK(0, 0),
			gridNote:       GridNote{gridKey: GK(0, 0), note: note{Action: 0}},
		}
		redo2 := UndoGridNote{
			overlayKey:     overlaykey.InitOverlayKey(1, 1),
			cursorPosition: GK(0, 0),
			gridNote:       GridNote{gridKey: GK(0, 0), note: note{Action: 0}},
		}
		m.PushUndoables(undo2, redo2)

		// Verify stack structure
		assert.NotEqual(t, NilStack, m.undoStack)
		assert.Equal(t, undo2, m.undoStack.undo)
		assert.Equal(t, redo2, m.undoStack.redo)
		assert.NotNil(t, m.undoStack.next)
		assert.Equal(t, undo, (*m.undoStack.next).undo)
	})
}

func TestUndoBeats(t *testing.T) {
	t.Run("Apply undo for beats change", func(t *testing.T) {
		m := setupTestModel()
		// Initial beats value
		(*m.definition.parts)[0].Beats = 8

		// Create and apply undoable
		undoable := UndoBeats{beats: 4, ArrCursor: m.arrangement.Cursor}
		undoable.ApplyUndo(m)

		// Verify beats were changed
		assert.Equal(t, uint8(4), (*m.definition.parts)[0].Beats)
	})
}

func TestUndoTempo(t *testing.T) {
	t.Run("Apply undo for tempo change", func(t *testing.T) {
		m := setupTestModel()
		// Initial tempo value
		m.definition.tempo = 130

		// Create and apply undoable
		undoable := UndoTempo{tempo: 120}
		undoable.ApplyUndo(m)

		// Verify tempo was changed
		assert.Equal(t, 120, m.definition.tempo)
	})
}

func TestUndoSubdivisions(t *testing.T) {
	t.Run("Apply undo for subdivisions change", func(t *testing.T) {
		m := setupTestModel()
		// Initial subdivisions value
		m.definition.subdivisions = 8

		// Create and apply undoable
		undoable := UndoSubdivisions{subdivisions: 4}
		undoable.ApplyUndo(m)

		// Verify subdivisions were changed
		assert.Equal(t, 4, m.definition.subdivisions)
	})
}

func TestUndoGridNote(t *testing.T) {
	t.Run("Apply undo for grid note change", func(t *testing.T) {
		m := setupTestModel()
		// Setup overlay and key
		ok := overlaykey.InitOverlayKey(1, 1)
		gk := GK(0, 0)
		testNote := note{Action: grid.ActionNothing}

		// Create undoable with a grid note
		undoable := UndoGridNote{
			overlayKey:     ok,
			cursorPosition: gk,
			gridNote:       GridNote{gridKey: gk, note: testNote},
			ArrCursor:      m.arrangement.Cursor,
		}

		// Apply the undo
		location := undoable.ApplyUndo(m)

		// Verify return values
		assert.Equal(t, ok, location.OverlayKey)
		assert.Equal(t, gk, location.GridKey)

		// Verify the note was set in the overlay
		assert.Equal(t, testNote, m.currentOverlay.Notes[gk])
	})
}

func TestUndoLineGridNotes(t *testing.T) {
	t.Run("Apply undo for line grid notes", func(t *testing.T) {
		m := setupTestModel()
		// Setup overlay and notes
		ok := overlaykey.InitOverlayKey(1, 1)
		line := uint8(2)

		// Create some test notes
		gk1 := GK(line, 0)
		gk2 := GK(line, 1)
		note1 := note{AccentIndex: 3}
		note2 := note{AccentIndex: 4}
		gridNotes := []GridNote{
			{gridKey: gk1, note: note1},
			{gridKey: gk2, note: note2},
		}

		// Add some existing notes on the line to be removed
		m.currentOverlay.SetNote(GK(line, 0), note{AccentIndex: 5})
		m.currentOverlay.SetNote(GK(line, 1), note{AccentIndex: 5})
		m.currentOverlay.SetNote(GK(line, 2), note{AccentIndex: 5})

		// Create undoable
		undoable := UndoLineGridNotes{
			overlayKey:     ok,
			cursorPosition: gk1,
			line:           line,
			gridNotes:      gridNotes,
			ArrCursor:      m.arrangement.Cursor,
		}

		// Set beats for the current part
		(*m.definition.parts)[0].Beats = 3

		// Apply the undo
		location := undoable.ApplyUndo(m)

		// Verify return values
		assert.Equal(t, ok, location.OverlayKey)
		assert.Equal(t, gk1, location.GridKey)

		// Verify the notes were set correctly
		assert.Equal(t, note1, m.currentOverlay.Notes[gk1])
		assert.Equal(t, note2, m.currentOverlay.Notes[gk2])
		_, exists := m.currentOverlay.Notes[GK(line, 2)]
		assert.False(t, exists, "Note at gridKey(2,2) should be removed")
	})
}

func TestUndoBounds(t *testing.T) {
	t.Run("Apply undo for bounds change", func(t *testing.T) {
		m := setupTestModel()
		// Setup overlay and bounds
		ok := overlaykey.InitOverlayKey(1, 1)
		bounds := InitBounds(GK(1, 0), GK(2, 2))

		// Add some existing notes in the bounds to be removed
		m.currentOverlay.SetNote(GK(1, 0), note{Action: grid.ActionNothing})
		m.currentOverlay.SetNote(GK(1, 1), note{Action: grid.ActionNothing})
		m.currentOverlay.SetNote(GK(2, 0), note{Action: grid.ActionNothing})

		// Create test grid notes to restore
		gk1 := GK(1, 0)
		gk2 := GK(2, 1)
		note1 := note{Action: grid.ActionNothing}
		note2 := note{Action: grid.ActionNothing}
		gridNotes := []GridNote{
			{gridKey: gk1, note: note1},
			{gridKey: gk2, note: note2},
		}

		// Create undoable
		undoable := UndoBounds{
			overlayKey:     ok,
			cursorPosition: gk1,
			bounds:         bounds,
			gridNotes:      gridNotes,
			ArrCursor:      m.arrangement.Cursor,
		}

		// Apply the undo
		location := undoable.ApplyUndo(m)

		// Verify return values
		assert.Equal(t, ok, location.OverlayKey)
		assert.Equal(t, gk1, location.GridKey)

		// Verify notes were set correctly
		assert.Equal(t, note1, m.currentOverlay.Notes[gk1])
		assert.Equal(t, note2, m.currentOverlay.Notes[gk2])

		// Check that other notes in bounds were removed
		_, exists := m.currentOverlay.Notes[GK(1, 1)]
		assert.False(t, exists, "Note at gridKey(1,1) should be removed")
		_, exists = m.currentOverlay.Notes[GK(2, 0)]
		assert.False(t, exists, "Note at gridKey(2,0) should be removed")
	})
}

func TestUndoGridNotes(t *testing.T) {
	t.Run("Apply undo for multiple grid notes", func(t *testing.T) {
		m := setupTestModel()
		// Setup overlay
		ok := overlaykey.InitOverlayKey(1, 1)

		// Create test grid notes
		gk1 := GK(1, 1)
		gk2 := GK(2, 2)
		note1 := note{Action: grid.ActionNothing}
		note2 := note{Action: grid.ActionNothing}
		gridNotes := []GridNote{
			{gridKey: gk1, note: note1},
			{gridKey: gk2, note: note2},
		}

		// Create undoable
		undoable := UndoGridNotes{
			overlayKey: ok,
			gridNotes:  gridNotes,
			ArrCursor:  m.arrangement.Cursor,
		}

		// Apply the undo
		location := undoable.ApplyUndo(m)

		// Verify return values
		assert.Equal(t, ok, location.OverlayKey)
		assert.Equal(t, gk1, location.GridKey)

		// Verify the notes were set correctly
		assert.Equal(t, note1, m.currentOverlay.Notes[gk1])
		assert.Equal(t, note2, m.currentOverlay.Notes[gk2])
	})
}

func TestUndoToNothing(t *testing.T) {
	t.Run("Apply undo to remove a note", func(t *testing.T) {
		m := setupTestModel()
		// Setup overlay and add a note to be removed
		ok := overlaykey.InitOverlayKey(1, 1)
		gk := GK(1, 1)
		m.currentOverlay.SetNote(gk, note{Action: grid.ActionNothing})

		// Create undoable
		undoable := UndoToNothing{
			overlayKey: ok,
			location:   gk,
			ArrCursor:  m.arrangement.Cursor,
		}

		// Apply the undo
		location := undoable.ApplyUndo(m)

		// Verify return values
		assert.Equal(t, ok, location.OverlayKey)
		assert.Equal(t, gk, location.GridKey)

		// Verify the note was removed
		_, exists := m.currentOverlay.Notes[gk]
		assert.False(t, exists, "Note should be removed")
	})
}

func TestUndoLineToNothing(t *testing.T) {
	t.Run("Apply undo to remove a line of notes", func(t *testing.T) {
		m := setupTestModel()
		// Setup overlay and add notes on a line to be removed
		ok := overlaykey.InitOverlayKey(1, 1)
		line := uint8(2)

		m.currentOverlay.SetNote(GK(line, 0), note{Action: grid.ActionNothing})
		m.currentOverlay.SetNote(GK(line, 1), note{Action: grid.ActionNothing})
		m.currentOverlay.SetNote(GK(line, 2), note{Action: grid.ActionNothing})

		// Create undoable
		undoable := UndoLineToNothing{
			overlayKey:     ok,
			cursorPosition: GK(line, 0),
			line:           line,
			ArrCursor:      m.arrangement.Cursor,
		}

		// Set beats for the current part
		(*m.definition.parts)[0].Beats = 3

		// Apply the undo
		location := undoable.ApplyUndo(m)

		// Verify return values
		assert.Equal(t, ok, location.OverlayKey)
		assert.Equal(t, GK(line, 0), location.GridKey)

		// Verify all notes on the line were removed
		_, exists := m.currentOverlay.Notes[GK(line, 0)]
		assert.False(t, exists, "Note at gridKey(2,0) should be removed")
		_, exists = m.currentOverlay.Notes[GK(line, 1)]
		assert.False(t, exists, "Note at gridKey(2,1) should be removed")
		_, exists = m.currentOverlay.Notes[GK(line, 2)]
		assert.False(t, exists, "Note at gridKey(2,2) should be removed")
	})
}

func TestUndoNewOverlay(t *testing.T) {
	t.Run("Apply undo for new overlay", func(t *testing.T) {
		m := setupTestModel()
		// Setup overlay and parts
		ok := overlaykey.InitOverlayKey(2, 2)

		// Add overlay to the current part
		previousTopOverlay := m.CurrentPart().Overlays
		newTopOverlay := &overlays.Overlay{
			Key:   ok,
			Notes: make(map[gridKey]note),
			Below: previousTopOverlay,
		}

		(*m.definition.parts)[m.CurrentPartID()].Overlays = newTopOverlay
		// Create undoable
		undoable := UndoNewOverlay{
			overlayKey:     ok,
			cursorPosition: GK(0, 0),
			ArrCursor:      m.arrangement.Cursor,
		}

		m.arrangement.NewPart(-1, false, false)
		m.currentOverlay = m.CurrentPart().Overlays

		// Apply the undo
		location := undoable.ApplyUndo(m)

		// Verify return values
		assert.Equal(t, ok, location.OverlayKey)
		assert.Equal(t, GK(0, 0), location.GridKey)

		// Verify the overlay was removed from current part
		foundOverlay := m.CurrentPart().Overlays.FindOverlay(location.OverlayKey)
		assert.Nil(t, foundOverlay, "Overlay should be removed")

		// Verify the overlay was stored in the part definition
		assert.Equal(t, previousTopOverlay, (*m.definition.parts)[0].Overlays)
	})
}
