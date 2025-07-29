package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
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
			playState:      InitLineStates(1, []linestate{}, 0),
			focus:          operation.FocusGrid, // Start with grid focus
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

func TestAdvanceKeyCycles(t *testing.T) {
	t.Run("1 cycle song", func(t *testing.T) {
		var counter int
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
			playState:      InitLineStates(1, []linestate{}, 0),
		}
		go func() {
			<-m.programChannel
		}()
		m.playing = PlayStandard
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PlayStopped {
			m.advanceKeyCycle()
			counter++
		}
		assert.Equal(t, 2, counter)
	})

	t.Run("2 iteration song", func(t *testing.T) {
		var counter int
		var parts = InitParts()
		var arr = InitArrangement(parts)
		def := Definition{
			arrangement: arr,
			parts:       &parts,
			keyline:     0,
			lines:       make([]grid.LineDefinition, 1),
		}
		m := model{
			arrangement:    arrangement.InitModel(def.arrangement, def.parts),
			definition:     def,
			programChannel: make(chan midiEventLoopMsg),
			playState:      InitLineStates(1, []linestate{}, 0),
		}
		m.arrangement.Cursor[0].IncreaseIterations()
		go func() {
			<-m.programChannel
		}()
		m.playing = PlayStandard
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PlayStopped {
			m.advanceKeyCycle()
			counter++
		}
		assert.Equal(t, 3, counter)
	})

	t.Run("2 cycle song", func(t *testing.T) {
		var counter int
		var parts = InitParts()
		var arr = InitArrangement(parts)
		def := Definition{
			arrangement: arr,
			parts:       &parts,
			keyline:     0,
			lines:       make([]grid.LineDefinition, 1),
		}
		m := model{
			arrangement:    arrangement.InitModel(def.arrangement, def.parts),
			definition:     def,
			programChannel: make(chan midiEventLoopMsg),
			playState:      InitLineStates(1, []linestate{}, 0),
		}
		arr = m.arrangement.Cursor.GetCurrentNode()
		(*arr).Section.Cycles = 2
		go func() {
			<-m.programChannel
		}()
		m.playing = PlayStandard
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PlayStopped {
			m.advanceKeyCycle()
			counter++
		}
		assert.Equal(t, 3, counter)
	})

	t.Run("nested cycle song 1 part", func(t *testing.T) {
		var counter int
		var parts = InitParts()

		group1 := &arrangement.Arrangement{
			Iterations: 2,
			Nodes:      make([]*arrangement.Arrangement, 0),
		}

		group2 := &arrangement.Arrangement{
			Iterations: 2,
			Nodes:      make([]*arrangement.Arrangement, 0),
		}

		nodeA := &arrangement.Arrangement{
			Section:    arrangement.SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
			Iterations: 1,
		}

		root := &arrangement.Arrangement{
			Iterations: 1,
			Nodes:      make([]*arrangement.Arrangement, 0),
		}

		root.Nodes = append(root.Nodes, group1)
		group1.Nodes = append(group1.Nodes, group2)
		group2.Nodes = append(group2.Nodes, nodeA)

		def := Definition{
			arrangement: root,
			parts:       &parts,
			keyline:     0,
			lines:       make([]grid.LineDefinition, 1),
		}
		logFile, _ := tea.LogToFile("debug.log", "debug")
		m := model{
			logFile:        logFile,
			arrangement:    arrangement.InitModel(def.arrangement, def.parts),
			definition:     def,
			programChannel: make(chan midiEventLoopMsg),
			playState:      InitLineStates(1, []linestate{}, 0),
		}
		arr := m.arrangement.Cursor.GetCurrentNode()
		(*arr).Section.Cycles = 1
		go func() {
			<-m.programChannel
		}()
		m.playing = PlayStandard
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PlayStopped {
			m.advanceKeyCycle()
			counter++
		}
		assert.Equal(t, 5, counter)
	})

	t.Run("grouped part with following part", func(t *testing.T) {
		var counter int
		var parts = InitParts()

		group1 := &arrangement.Arrangement{
			Iterations: 1,
			Nodes:      make([]*arrangement.Arrangement, 0),
		}

		nodeA := &arrangement.Arrangement{
			Section:    arrangement.SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
			Iterations: 1,
		}

		nodeB := &arrangement.Arrangement{
			Section:    arrangement.SongSection{Part: 3, Cycles: 1, StartBeat: 0, StartCycles: 1},
			Iterations: 1,
		}

		root := &arrangement.Arrangement{
			Iterations: 1,
			Nodes:      make([]*arrangement.Arrangement, 0),
		}

		root.Nodes = append(root.Nodes, group1, nodeB)
		group1.Nodes = append(group1.Nodes, nodeA)
		root.ResetAllPlayCycles()

		def := Definition{
			arrangement: root,
			parts:       &parts,
			keyline:     0,
			lines:       make([]grid.LineDefinition, 1),
		}
		logFile, _ := tea.LogToFile("debug.log", "debug")
		m := model{
			logFile:        logFile,
			arrangement:    arrangement.InitModel(def.arrangement, def.parts),
			definition:     def,
			programChannel: make(chan midiEventLoopMsg),
			playState:      InitLineStates(1, []linestate{}, 0),
		}
		go func() {
			<-m.programChannel
		}()
		m.playing = PlayStandard
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PlayStopped {
			m.advanceKeyCycle()
			counter++
			if counter > 3 {
				assert.Fail(t, "Should not go past 3 iterations")
				break
			}
		}
		assert.Equal(t, 2, counter)
	})
}
