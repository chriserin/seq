package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/mappings"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/seqmidi"
	"github.com/stretchr/testify/assert"
)

// TestUpdateCursorMovements tests all cursor movement keybindings

func createTestModel(modelFns ...modelFunc) model {

	m := InitModel("", seqmidi.MidiConnection{}, "", "", MlmStandAlone, "default")

	for _, fn := range modelFns {
		m = fn(&m)
	}

	return m
}

type TestKey struct {
	Keys string
}

func processCommands(commands []any, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	for _, command := range commands {
		switch c := command.(type) {
		case mappings.Command:
			m, cmd = processCommand(c, m)
		case mappings.Mapping:
			m, cmd = processMapping(c, m)
		case TestKey:
			m, cmd = processTestKey(c, m)
		}
		if cmd != nil {
			updateModel, _ := m.Update(cmd())
			switch um := updateModel.(type) {
			case model:
				m = um
			}
		}
	}
	return m, cmd
}

func processTestKey(testKey TestKey, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	var updateModel tea.Model
	for _, key := range testKey.Keys {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}}
		updateModel, cmd = m.Update(keyMsg)
		switch um := updateModel.(type) {
		case model:
			m = um
		}
	}
	return m, cmd
}

func processCommand(command mappings.Command, m model) (model, tea.Cmd) {
	keyMsgs := getKeyMsgs(command)
	var cmd tea.Cmd
	for _, keyMsg := range keyMsgs {
		var updateModel tea.Model
		updateModel, cmd = m.Update(keyMsg)
		switch um := updateModel.(type) {
		case model:
			m = um
		}
	}
	return m, cmd
}

func processMapping(mapping mappings.Mapping, m model) (model, tea.Cmd) {
	var cmd tea.Cmd
	var updateModel tea.Model
	switch mapping.Command {
	case mappings.NumberPattern:
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(mapping.LastValue)}
		updateModel, cmd = m.Update(keyMsg)
		switch um := updateModel.(type) {
		case model:
			m = um
		}
	}
	return m, cmd
}

func getKeyMsgs(command mappings.Command) []tea.KeyMsg {
	keys := mappings.KeysForCommand(command)
	var keyMsgs []tea.KeyMsg
	for _, key := range keys {
		keyMsgs = append(keyMsgs, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	}
	return keyMsgs
}

type modelFunc func(m *model) model

func WithCurosrPos(pos grid.GridKey) modelFunc {
	return func(m *model) model {
		m.cursorPos = pos
		return *m
	}
}

func WithNonRootOverlay(overlayKey overlaykey.OverlayPeriodicity) modelFunc {
	return func(m *model) model {
		(*m.definition.parts)[0].Overlays = m.CurrentPart().Overlays.Add(overlayKey)
		m.currentOverlay = m.CurrentPart().Overlays.FindAboveOverlay(overlayKey)
		return *m
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name        string
		command     mappings.Command
		description string
	}{
		{
			name:        "Save With Filename",
			command:     mappings.Save,
			description: "Should save file when filename is set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary filename
			tempDir := t.TempDir()
			testFilename := filepath.Join(tempDir, "test_save.seq")

			m := createTestModel(
				func(m *model) model {
					m.filename = testFilename
					return *m
				},
			)

			_, err := os.Stat(testFilename)
			assert.True(t, os.IsNotExist(err), "File should not exist initially")

			processCommand(tt.command, m)

			_, err = os.Stat(testFilename)
			assert.NoError(t, err, tt.description+" - file should be created")

			fileInfo, err := os.Stat(testFilename)
			assert.NoError(t, err, "Should be able to get file info")
			assert.Greater(t, fileInfo.Size(), int64(0), "File should not be empty")
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "New Sequence Clears Notes",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.New,
				mappings.Enter,
			},
			description: "Should clear all notes when creating new sequence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			note, exists := m.CurrentNote()
			assert.False(t, exists, tt.description+" - note should not exist after new sequence")
			assert.Equal(t, zeronote, note, tt.description+" - note should be zero note")

			assert.Equal(t, grid.GridKey{Line: 0, Beat: 0}, m.cursorPos, "Cursor should be reset to origin")
		})
	}
}

func TestNextPrevTheme(t *testing.T) {
	tests := []struct {
		name          string
		commands      []any
		initialTheme  string
		expectedTheme string
		description   string
	}{
		{
			name:          "NextTheme from default advances to seafoam",
			commands:      []any{mappings.NextTheme},
			initialTheme:  "default",
			expectedTheme: "seafoam",
			description:   "Should advance from default to seafoam theme",
		},
		{
			name:          "NextTheme from last theme wraps to first",
			commands:      []any{mappings.NextTheme},
			initialTheme:  "miles",
			expectedTheme: "default",
			description:   "Should wrap from last theme (miles) to first theme (default)",
		},
		{
			name:          "PrevTheme from default wraps to last",
			commands:      []any{mappings.PrevTheme},
			initialTheme:  "default",
			expectedTheme: "miles",
			description:   "Should wrap from first theme (default) to last theme (miles)",
		},
		{
			name:          "PrevTheme from seafoam goes back to default",
			commands:      []any{mappings.PrevTheme},
			initialTheme:  "seafoam",
			expectedTheme: "default",
			description:   "Should go back from seafoam to default theme",
		},
		{
			name:          "Multiple NextTheme commands cycle correctly",
			commands:      []any{mappings.NextTheme, mappings.NextTheme, mappings.NextTheme},
			initialTheme:  "default",
			expectedTheme: "springtime",
			description:   "Should advance from default -> seafoam -> dynamite -> springtime",
		},
		{
			name:          "NextTheme then PrevTheme returns to original",
			commands:      []any{mappings.NextTheme, mappings.PrevTheme},
			initialTheme:  "cyberpunk",
			expectedTheme: "cyberpunk",
			description:   "Should return to original theme after next then prev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				func(m *model) model {
					m.theme = tt.initialTheme
					return *m
				},
			)

			assert.Equal(t, tt.initialTheme, m.theme, "Initial theme should match")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedTheme, m.theme, tt.description+" - theme should match expected value")
		})
	}
}

func TestClearLine(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		cursorPos   grid.GridKey
		description string
	}{
		{
			name: "Clear line from beginning cursor position",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLineStart,
				mappings.ClearLine,
			},
			cursorPos:   grid.GridKey{Line: 0, Beat: 0},
			description: "Should clear all notes from cursor position to end of line",
		},
		{
			name: "Clear line from middle cursor position",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLeft,
				mappings.ClearLine,
			},
			cursorPos:   grid.GridKey{Line: 0, Beat: 1},
			description: "Should keep notes before cursor position and clear from cursor to end",
		},
		{
			name: "Clear line from end cursor position",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.ClearLine,
			},
			cursorPos:   grid.GridKey{Line: 0, Beat: 2},
			description: "Should keep notes before cursor position and clear only the cursor position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			for beat := uint8(0); beat < m.CurrentPart().Beats; beat++ {
				m.cursorPos = grid.GridKey{Line: tt.cursorPos.Line, Beat: beat}
				_, exists := m.CurrentNote()

				if beat < tt.cursorPos.Beat {
					assert.True(t, exists, tt.description+" - note should exist before cursor at beat %d", beat)
				} else {
					assert.False(t, exists, tt.description+" - note should not exist at or after cursor at beat %d", beat)
				}
			}
		})
	}
}

func TestClearOverlay(t *testing.T) {
	tests := []struct {
		name        string
		commands    []any
		description string
	}{
		{
			name: "Clear overlay removes current overlay",
			commands: []any{
				mappings.NoteAdd,
				mappings.ClearOverlay,
			},
			description: "Should remove the current overlay from the part",
		},
		{
			name: "Clear overlay with multiple overlays",
			commands: []any{
				mappings.NoteAdd,
				mappings.NextOverlay,
				mappings.NoteAdd,
				mappings.ClearOverlay,
			},
			description: "Should remove the current overlay and switch to next available overlay",
		},
	}

	overlayKey := overlaykey.OverlayPeriodicity{
		Shift:      2,
		Interval:   4,
		Width:      0,
		StartCycle: 0,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithNonRootOverlay(overlayKey),
			)

			initialOverlayKey := m.currentOverlay.Key
			keys := make([]overlaykey.OverlayPeriodicity, 0)
			m.CurrentPart().Overlays.CollectKeys(&keys)
			initialOverlayCount := len(keys)

			m, _ = processCommands(tt.commands, m)

			keys = make([]overlaykey.OverlayPeriodicity, 0)
			m.CurrentPart().Overlays.CollectKeys(&keys)
			finalOverlayCount := len(keys)
			assert.Equal(t, initialOverlayCount-1, finalOverlayCount, tt.description+" - overlay count should decrease by 1")

			assert.NotEqual(t, initialOverlayKey, m.currentOverlay.Key, tt.description+" - should switch to different overlay")
		})
	}
}

func TestActionMappings(t *testing.T) {
	tests := []struct {
		name           string
		commands       []any
		expectedAction grid.Action
		description    string
	}{
		{
			name:           "ActionAddLineReset sets ActionLineReset",
			commands:       []any{mappings.ActionAddLineReset},
			expectedAction: grid.ActionLineReset,
			description:    "Should set note action to ActionLineReset",
		},
		{
			name:           "ActionAddLineReverse sets ActionLineReverse",
			commands:       []any{mappings.ActionAddLineReverse},
			expectedAction: grid.ActionLineReverse,
			description:    "Should set note action to ActionLineReverse",
		},
		{
			name:           "ActionAddSkipBeat sets ActionLineSkipBeat",
			commands:       []any{mappings.ActionAddSkipBeat},
			expectedAction: grid.ActionLineSkipBeat,
			description:    "Should set note action to ActionLineSkipBeat",
		},
		{
			name:           "ActionAddReset sets ActionReset",
			commands:       []any{mappings.ActionAddReset},
			expectedAction: grid.ActionReset,
			description:    "Should set note action to ActionReset",
		},
		{
			name:           "ActionAddLineBounce sets ActionLineBounce",
			commands:       []any{mappings.ActionAddLineBounce},
			expectedAction: grid.ActionLineBounce,
			description:    "Should set note action to ActionLineBounce",
		},
		{
			name:           "ActionAddLineDelay sets ActionLineDelay",
			commands:       []any{mappings.ActionAddLineDelay},
			expectedAction: grid.ActionLineDelay,
			description:    "Should set note action to ActionLineDelay",
		},
		{
			name:           "ActionAddSpecificValue sets ActionSpecificValue",
			commands:       []any{mappings.SetupInputSwitch, mappings.SetupInputSwitch, mappings.Increase, mappings.Increase, mappings.Escape, mappings.ActionAddSpecificValue},
			expectedAction: grid.ActionSpecificValue,
			description:    "Should set note action to ActionSpecificValue when line is ProgramChange",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			note, exists := m.CurrentNote()
			assert.True(t, exists, tt.description+" - note should exist after adding action")
			assert.Equal(t, tt.expectedAction, note.Action, tt.description+" - note action should match expected value")
		})
	}
}

func TestSelectKeyLine(t *testing.T) {
	tests := []struct {
		name            string
		cursorPos       grid.GridKey
		expectedKeyline uint8
		description     string
	}{
		{
			name:            "SelectKeyLine sets keyline to cursor line 0",
			cursorPos:       grid.GridKey{Line: 0, Beat: 0},
			expectedKeyline: 0,
			description:     "Should set keyline to 0 when cursor is on line 0",
		},
		{
			name:            "SelectKeyLine sets keyline to cursor line 2",
			cursorPos:       grid.GridKey{Line: 2, Beat: 5},
			expectedKeyline: 2,
			description:     "Should set keyline to 2 when cursor is on line 2",
		},
		{
			name:            "SelectKeyLine sets keyline to cursor line 7",
			cursorPos:       grid.GridKey{Line: 7, Beat: 3},
			expectedKeyline: 7,
			description:     "Should set keyline to 7 when cursor is on line 7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel(
				WithCurosrPos(tt.cursorPos),
			)

			m, _ = processCommands([]any{mappings.SelectKeyLine}, m)

			assert.Equal(t, tt.expectedKeyline, m.definition.keyline, tt.description+" - keyline should match cursor line")
		})
	}
}

func TestOverlayStackToggle(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		expectedPressUp   bool
		expectedPressDown bool
		description       string
	}{
		{
			name:              "Toggle from neither to PressUp",
			commands:          []any{mappings.OverlayStackToggle},
			expectedPressUp:   false,
			expectedPressDown: true,
			description:       "Should set PressUp to true when neither is set",
		},
		{
			name:              "Toggle from PressUp to PressDown",
			commands:          []any{mappings.OverlayStackToggle, mappings.OverlayStackToggle},
			expectedPressUp:   false,
			expectedPressDown: false,
			description:       "Should set PressDown to true when PressUp is set",
		},
		{
			name:              "Toggle from PressDown to neither",
			commands:          []any{mappings.OverlayStackToggle, mappings.OverlayStackToggle, mappings.OverlayStackToggle},
			expectedPressUp:   true,
			expectedPressDown: false,
			description:       "Should set both to false when PressDown is set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			assert.Equal(t, true, m.currentOverlay.PressUp, tt.description+" - Initial PressUp should be true")
			assert.Equal(t, false, m.currentOverlay.PressDown, tt.description+" - Initial PressDown should be false")

			m, _ = processCommands(tt.commands, m)

			assert.Equal(t, tt.expectedPressUp, m.currentOverlay.PressUp, tt.description+" - PressUp should match expected value")
			assert.Equal(t, tt.expectedPressDown, m.currentOverlay.PressDown, tt.description+" - PressDown should match expected value")
		})
	}
}

func TestYankAndPaste(t *testing.T) {
	tests := []struct {
		name              string
		commands          []any
		description       string
		expectedNoteBeats []uint8
	}{
		{
			name: "Yank single note and paste",
			commands: []any{
				mappings.NoteAdd,
				mappings.Yank,
				mappings.CursorDown,
				mappings.Paste,
			},
			expectedNoteBeats: []uint8{0},
			description:       "Should yank a single note and paste it at cursor position",
		},
		{
			name: "Yank multiple notes in visual mode and paste",
			commands: []any{
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorLineStart,
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.Yank,
				mappings.CursorDown,
				mappings.CursorLineStart,
				mappings.Paste,
			},
			expectedNoteBeats: []uint8{0, 1, 2},
			description:       "Should yank multiple notes in visual mode and paste them",
		},
		{
			name: "Yank and paste with cursor movement",
			commands: []any{
				mappings.NoteAdd,
				mappings.AccentIncrease,
				mappings.AccentIncrease,
				mappings.Yank,
				mappings.CursorDown,
				mappings.CursorRight,
				mappings.Paste,
			},
			expectedNoteBeats: []uint8{1},
			description:       "Should yank note with modifications and paste at different location",
		},
		{
			name: "Yank empty selection should not crash",
			commands: []any{
				mappings.Yank,
				mappings.CursorDown,
				mappings.Paste,
			},
			expectedNoteBeats: []uint8{},
			description:       "Should handle yanking empty selection gracefully",
		},
		{
			name: "Yank does not capture empty space", //NOTE: Should it? This wasn't an intentional behavior, but how it fell out
			commands: []any{
				mappings.ToggleVisualMode,
				mappings.CursorRight,
				mappings.CursorRight,
				mappings.Yank,
				mappings.CursorDown,
				mappings.CursorLineStart,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.NoteAdd,
				mappings.CursorRight,
				mappings.CursorLineStart,
				mappings.Paste,
			},
			expectedNoteBeats: []uint8{0, 1, 2},
			description:       "Should not paste empty space when yanking in visual mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := createTestModel()

			m, _ = processCommands(tt.commands, m)

			for beat := range uint8(32) {
				_, exists := m.currentOverlay.GetNote(grid.GridKey{Line: m.cursorPos.Line, Beat: beat})
				if slices.Contains(tt.expectedNoteBeats, uint8(beat)) {
					fmt.Println("Checking beat:", beat)
					assert.True(t, exists, tt.description+" - note should exist at beat "+string(beat))
				} else {
					fmt.Println("Checking no exist beat:", beat)
					assert.False(t, exists, tt.description+" - note should not exist at beat "+string(beat))
				}
			}
		})
	}
}
