package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/operation"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

var OK = overlaykey.InitOverlayKey

func TestUpdateArrangementFocus(t *testing.T) {
	t.Run("switch to arrangement view and create part", func(t *testing.T) {
		// Setup a model with a basic arrangement
		var parts = InitParts()
		var arr = InitArrangement(parts)
		def := Definition{
			arrangement: arr,
			parts:       &parts,
			keyline:     0,
		}

		m := model{
			arrangement:    arrangement.InitModel(def.arrangement, def.parts),
			definition:     def,
			programChannel: make(chan midiEventLoopMsg),
			playState: playState{
				lineStates: InitLineStates(1, []linestate{}, 0),
			},
			focus: operation.FocusGrid, // Start with grid focus
		}

		initialNodeCount := m.arrangement.Root.CountEndNodes()
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlA})
		modelPtr := updatedModel.(model)

		assert.Equal(t, operation.FocusArrangementEditor, modelPtr.focus, "Model should have arrangement editor focus")
		assert.True(t, modelPtr.arrangement.Focus, "Arrangement model should have focus flag set to true")

		updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyCtrlCloseBracket})
		modelPtr = updatedModel.(model)

		assert.Equal(t, operation.FocusArrangementEditor, modelPtr.focus, "Model should have arrangement editor focus")
		assert.True(t, modelPtr.arrangement.Focus, "Arrangement model should have focus flag set to true")

		updatedModelAfterPart, _ := updatedModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
		finalModel := updatedModelAfterPart.(model)

		finalNodeCount := finalModel.arrangement.Root.CountEndNodes()
		assert.Greater(t, finalNodeCount, initialNodeCount, "Arrangement should have more end nodes after part creation")

		assert.Equal(t, operation.FocusArrangementEditor, finalModel.focus, "Model should still have arrangement editor focus")
		assert.True(t, finalModel.arrangement.Focus, "Arrangement model should still have focus flag set to true")
		assert.Equal(t, operation.SelectGrid, finalModel.selectionIndicator, "Selection indicator should be reset to nothing")
	})
}

func TestSolo(t *testing.T) {
	t.Run("First Solo", func(t *testing.T) {
		playStates := []linestate{
			{groupPlayState: PlayStatePlay},
			{groupPlayState: PlayStatePlay},
		}
		newPlayStates := Solo(playStates, 0)
		assert.Equal(t, newPlayStates[0].groupPlayState, PlayStateSolo)
		assert.Equal(t, newPlayStates[1].groupPlayState, PlayStatePlay)
	})

	t.Run("First UnSolo", func(t *testing.T) {
		playStates := []linestate{
			{groupPlayState: PlayStateSolo},
			{groupPlayState: PlayStatePlay},
		}
		newPlayStates := Solo(playStates, 0)
		assert.Equal(t, newPlayStates[0].groupPlayState, PlayStatePlay)
		assert.Equal(t, newPlayStates[1].groupPlayState, PlayStatePlay)
	})
}
