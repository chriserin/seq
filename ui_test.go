package main

import (
	"slices"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

var OK = overlaykey.InitOverlayKey

func TestOverlayKeySort(t *testing.T) {
	testCases := []struct {
		desc     string
		keys     []overlayKey
		expected overlayKey
	}{
		{
			desc:     "TEST A",
			keys:     []overlayKey{OK(2, 1), OK(3, 1)},
			expected: OK(3, 1),
		},
		{
			desc:     "TEST B",
			keys:     []overlayKey{OK(3, 1), OK(2, 1)},
			expected: OK(3, 1),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			newSlice := make([]overlayKey, len(tC.keys))
			copy(newSlice, tC.keys)
			slices.SortFunc(newSlice, overlaykey.Compare)
			assert.Equal(t, tC.expected, newSlice[0])
		})
	}
}

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
			focus:          FOCUS_GRID, // Start with grid focus
		}

		// // Create a goroutine to consume from the channel so it doesn't block
		// go func() {
		// 	<-m.programChannel
		// }()

		initialNodeCount := m.arrangement.Root.CountEndNodes()
		updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
		modelPtr := updatedModel.(model)

		assert.Equal(t, FOCUS_ARRANGEMENT_EDITOR, modelPtr.focus, "Model should have arrangement editor focus")
		assert.True(t, modelPtr.arrangement.Focus, "Arrangement model should have focus flag set to true")

		updatedModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyCtrlCloseBracket})
		modelPtr = updatedModel.(model)

		assert.Equal(t, FOCUS_ARRANGEMENT_EDITOR, modelPtr.focus, "Model should have arrangement editor focus")
		assert.True(t, modelPtr.arrangement.Focus, "Arrangement model should have focus flag set to true")

		updatedModelAfterPart, _ := updatedModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
		finalModel := updatedModelAfterPart.(model)

		finalNodeCount := finalModel.arrangement.Root.CountEndNodes()
		assert.Greater(t, finalNodeCount, initialNodeCount, "Arrangement should have more end nodes after part creation")

		assert.Equal(t, FOCUS_ARRANGEMENT_EDITOR, finalModel.focus, "Model should still have arrangement editor focus")
		assert.True(t, finalModel.arrangement.Focus, "Arrangement model should still have focus flag set to true")
		assert.Equal(t, SELECT_ARRANGEMENT_EDITOR, finalModel.selectionIndicator, "Selection indicator should be reset to nothing")
	})
}

func TestSolo(t *testing.T) {
	t.Run("First Solo", func(t *testing.T) {
		playStates := []linestate{
			{groupPlayState: PLAY_STATE_PLAY},
			{groupPlayState: PLAY_STATE_PLAY},
		}
		newPlayStates := Solo(playStates, 0)
		assert.Equal(t, newPlayStates[0].groupPlayState, PLAY_STATE_SOLO)
		assert.Equal(t, newPlayStates[1].groupPlayState, PLAY_STATE_PLAY)
	})

	t.Run("First UnSolo", func(t *testing.T) {
		playStates := []linestate{
			{groupPlayState: PLAY_STATE_SOLO},
			{groupPlayState: PLAY_STATE_PLAY},
		}
		newPlayStates := Solo(playStates, 0)
		assert.Equal(t, newPlayStates[0].groupPlayState, PLAY_STATE_PLAY)
		assert.Equal(t, newPlayStates[1].groupPlayState, PLAY_STATE_PLAY)
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
		m.playing = PLAY_STANDARD
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PLAY_STOPPED {
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
		m.playing = PLAY_STANDARD
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PLAY_STOPPED {
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
		m.playing = PLAY_STANDARD
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PLAY_STOPPED {
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
		m.playing = PLAY_STANDARD
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PLAY_STOPPED {
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
		m.playing = PLAY_STANDARD
		m.arrangement.Cursor.ResetIterations()
		for m.playing != PLAY_STOPPED {
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
