package main

import (
	"os"
	"path/filepath"
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
