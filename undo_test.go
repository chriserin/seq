package main

import (
	"testing"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
	"github.com/chriserin/seq/internal/playstate"
	"github.com/chriserin/seq/internal/sequence"
	"github.com/stretchr/testify/assert"
)

func initPlayState() playstate.PlayState {
	return playstate.PlayState{
		LineStates: playstate.InitLineStates(1, []playstate.LineState{}, 0),
	}
}

func setupTestModel() *model {
	var parts = InitParts()
	var arr = InitArrangement(parts)
	def := sequence.Sequence{
		Arrangement:  arr,
		Parts:        &parts,
		Keyline:      0,
		Tempo:        120,
		Subdivisions: 4,
		Lines:        make([]grid.LineDefinition, 1),
	}

	m := &model{
		definition:  def,
		playState:   initPlayState(),
		undoStack:   EmptyStack,
		redoStack:   EmptyStack,
		arrangement: arrangement.InitModel(def.Arrangement, def.Parts),
	}

	m.currentOverlay = parts[0].Overlays
	return m
}

func TestUndoBeats(t *testing.T) {
	t.Run("Apply undo for beats change", func(t *testing.T) {
		m := setupTestModel()
		// Initial beats value
		(*m.definition.Parts)[0].Beats = 8

		// Create and apply undoable
		undoable := UndoBeats{beats: 4, ArrCursor: m.arrangement.Cursor}
		undoable.ApplyUndo(m)

		// Verify beats were changed
		assert.Equal(t, uint8(4), (*m.definition.Parts)[0].Beats)
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

		(*m.definition.Parts)[m.CurrentPartID()].Overlays = newTopOverlay
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
		assert.Equal(t, previousTopOverlay, (*m.definition.Parts)[0].Overlays)
	})
}
