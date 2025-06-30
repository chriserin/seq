package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
	"github.com/stretchr/testify/assert"
)

func TestReadWrite(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "seq-readwrite-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("Simple model with basic settings", func(t *testing.T) {
		// Create a model with basic settings
		originalModel := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:  "TestPart",
						Beats: 16,
					},
				},
				tempo:           140,
				subdivisions:    4,
				keyline:         2,
				instrument:      "synth",
				template:        "custom",
				templateUIStyle: "light",
			},
		}

		// Write it to a file
		filename := filepath.Join(tempDir, "simple_model.txt")
		err := Write(originalModel, filename)
		assert.NoError(t, err)

		// Read it back
		readDef, err := Read(filename)
		assert.NoError(t, err)
		assert.NotNil(t, readDef)

		// Verify settings are preserved
		assert.Equal(t, 140, readDef.tempo)
		assert.Equal(t, 4, readDef.subdivisions)
		assert.Equal(t, uint8(2), readDef.keyline)
		assert.Equal(t, "synth", readDef.instrument)
		assert.Equal(t, "custom", readDef.template)
		assert.Equal(t, "light", readDef.templateUIStyle)

		// Verify parts
		assert.NotNil(t, readDef.parts)
		assert.Len(t, *readDef.parts, 1)
		assert.Equal(t, "TestPart", (*readDef.parts)[0].Name)
		assert.Equal(t, uint8(16), (*readDef.parts)[0].Beats)
	})

	t.Run("Model with lines and accents", func(t *testing.T) {
		// Create a model with lines and accents
		originalModel := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:  "TestPart",
						Beats: 8,
					},
				},
				lines: []grid.LineDefinition{
					{Channel: 1, Note: 60, MsgType: 0},
					{Channel: 2, Note: 67, MsgType: 1},
				},
				accents: patternAccents{
					Diff:   5,
					Start:  50,
					Target: AccentTargetVelocity,
					Data: []config.Accent{
						{Value: 100},
						{Value: 80},
					},
				},
			},
		}

		// Write it to a file
		filename := filepath.Join(tempDir, "model_with_lines.txt")
		err := Write(originalModel, filename)
		assert.NoError(t, err)

		// Read the file content to verify it's correctly written
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)
		contentStr := string(content)

		// Verify accent data is written properly
		assert.Contains(t, contentStr, "Accent 0:")
		assert.Contains(t, contentStr, "Accent 1:")

		// Read it back
		readDef, err := Read(filename)
		assert.NoError(t, err)
		assert.NotNil(t, readDef)

		// Verify lines
		assert.Len(t, readDef.lines, 2)
		assert.Equal(t, uint8(1), readDef.lines[0].Channel)
		assert.Equal(t, uint8(60), readDef.lines[0].Note)
		assert.Equal(t, grid.MessageType(0), readDef.lines[0].MsgType)
		assert.Equal(t, uint8(2), readDef.lines[1].Channel)
		assert.Equal(t, uint8(67), readDef.lines[1].Note)
		assert.Equal(t, grid.MessageType(1), readDef.lines[1].MsgType)

		// Verify accents
		assert.Equal(t, uint8(5), readDef.accents.Diff)
		assert.Equal(t, uint8(50), readDef.accents.Start)
		assert.Equal(t, AccentTargetVelocity, readDef.accents.Target)
		assert.Len(t, readDef.accents.Data, 2)
	})

	t.Run("Model with multiple parts", func(t *testing.T) {
		// Create a model with lines and accents
		originalModel := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:  "TestPart",
						Beats: 8,
					},
					{
						Name:  "TestPart2",
						Beats: 8,
					},
				},
			},
		}

		// Write it to a file
		filename := filepath.Join(tempDir, "model_with_parts.txt")
		err := Write(originalModel, filename)
		assert.NoError(t, err)

		// Read it back
		readDef, err := Read(filename)
		assert.NoError(t, err)
		assert.NotNil(t, readDef)

		// Verify lines
		assert.Len(t, *readDef.parts, 2)
	})

	t.Run("Model with overlays and notes", func(t *testing.T) {
		// Create overlay with notes
		key := overlaykey.OverlayPeriodicity{
			Shift:      2,
			Interval:   4,
			Width:      1,
			StartCycle: 0,
		}
		overlay := overlays.InitOverlay(key, nil)
		overlay.PressUp = true
		overlay.PressDown = false

		// Add some notes to the overlay
		note1 := grid.InitNote()
		note1.AccentIndex = 1
		note1.GateIndex = 2
		gridKey1 := grid.GridKey{Line: 0, Beat: 0}
		overlay.SetNote(gridKey1, note1)

		note2 := grid.InitNote()
		note2.AccentIndex = 3
		note2.WaitIndex = 1
		note2.Ratchets.Hits = 3
		note2.Ratchets.Length = 2
		note2.Ratchets.Span = 1
		gridKey2 := grid.GridKey{Line: 1, Beat: 2}
		overlay.SetNote(gridKey2, note2)

		// Create a model with overlays
		originalModel := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:     "PartWithOverlay",
						Beats:    8,
						Overlays: overlay,
					},
				},
			},
		}

		// Write it to a file
		filename := filepath.Join(tempDir, "model_with_overlays.txt")
		err := Write(originalModel, filename)
		assert.NoError(t, err)

		// Read it back
		readDef, err := Read(filename)
		assert.NoError(t, err)
		assert.NotNil(t, readDef)

		// Verify overlay structure
		assert.NotNil(t, (*readDef.parts)[0].Overlays)
		readOverlay := (*readDef.parts)[0].Overlays

		// Verify overlay key
		assert.Equal(t, uint8(2), readOverlay.Key.Shift)
		assert.Equal(t, uint8(4), readOverlay.Key.Interval)
		assert.Equal(t, uint8(1), readOverlay.Key.Width)
		assert.Equal(t, uint8(0), readOverlay.Key.StartCycle)
		assert.Equal(t, true, readOverlay.PressUp)
		assert.Equal(t, false, readOverlay.PressDown)

		// Verify notes
		assert.NotEmpty(t, readOverlay.Notes)
		assert.Len(t, readOverlay.Notes, 2)

		// Verify first note
		if note, exists := readOverlay.Notes[grid.GridKey{Line: 0, Beat: 0}]; assert.True(t, exists) {
			assert.Equal(t, uint8(1), note.AccentIndex)
			assert.Equal(t, uint8(2), note.GateIndex)
		}

		// Verify second note
		if note, exists := readOverlay.Notes[grid.GridKey{Line: 1, Beat: 2}]; assert.True(t, exists) {
			assert.Equal(t, uint8(3), note.AccentIndex)
			assert.Equal(t, uint8(1), note.WaitIndex)
			assert.Equal(t, uint8(3), note.Ratchets.Hits)
			assert.Equal(t, uint8(2), note.Ratchets.Length)
			assert.Equal(t, uint8(1), note.Ratchets.Span)
		}
	})

	t.Run("Basic arrangement", func(t *testing.T) {
		// Create a simple model with basic settings
		originalModel := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:  "TestPart",
						Beats: 16,
					},
				},
				tempo:           140,
				subdivisions:    4,
				keyline:         2,
				instrument:      "synth",
				template:        "custom",
				templateUIStyle: "light",
			},
		}

		// Create arrangement data
		section := arrangement.InitSongSection(0)
		section.Cycles = 2
		section.StartBeat = 1

		// Modify original model to add arrangement
		originalModel.definition.arrangement = &arrangement.Arrangement{
			Iterations: 1,
			Section:    section,
		}

		// Write it to a file
		filename := filepath.Join(tempDir, "basic_arrangement.txt")
		err := Write(originalModel, filename)
		assert.NoError(t, err)

		// Read it back
		readDef, err := Read(filename)
		assert.NoError(t, err)
		assert.NotNil(t, readDef)

		// Verify arrangement data
		assert.NotNil(t, readDef.arrangement)
		assert.Equal(t, 1, readDef.arrangement.Iterations)
		assert.Equal(t, 0, readDef.arrangement.Section.Part)
		assert.Equal(t, 2, readDef.arrangement.Section.Cycles)
		assert.Equal(t, 1, readDef.arrangement.Section.StartBeat)
	})
}

func TestReadFileWithChords(t *testing.T) {
	// Test reading the checkchord.seq file which contains chord definitions
	readDef, err := Read("checkchord.seq")
	assert.NoError(t, err)
	assert.NotNil(t, readDef)

	// Verify basic settings from checkchord.seq
	assert.Equal(t, 120, readDef.tempo)
	assert.Equal(t, 2, readDef.subdivisions)
	assert.Equal(t, uint8(0), readDef.keyline)
	assert.Equal(t, "Standard", readDef.instrument)
	assert.Equal(t, "Piano2", readDef.template)
	assert.Equal(t, "blackwhite", readDef.templateUIStyle)

	// Verify parts exist
	assert.NotNil(t, readDef.parts)
	assert.Len(t, *readDef.parts, 1)
	part := (*readDef.parts)[0]
	assert.Equal(t, "Part 1", part.Name)
	assert.Equal(t, uint8(32), part.Beats)

	// Verify overlay exists
	assert.NotNil(t, part.Overlays)
	overlay := part.Overlays
	assert.Equal(t, uint8(1), overlay.Key.Shift)
	assert.Equal(t, uint8(1), overlay.Key.Interval)
	assert.Equal(t, uint8(0), overlay.Key.Width)
	assert.Equal(t, uint8(0), overlay.Key.StartCycle)
	assert.Equal(t, true, overlay.PressUp)
	assert.Equal(t, false, overlay.PressDown)

	// Verify chord exists
	assert.NotEmpty(t, overlay.Chords)
	assert.Len(t, overlay.Chords, 1)
	chord := overlay.Chords[0]

	// Verify chord properties from checkchord.seq
	expectedGridKey := grid.GridKey{Line: 24, Beat: 0}
	assert.Equal(t, expectedGridKey, chord.Root)
	assert.Equal(t, overlays.Arp(2), chord.Arppegio)
	assert.Equal(t, uint8(3), chord.Double)
	assert.Equal(t, uint32(137), chord.Chord.Notes)
	assert.Equal(t, int8(0), chord.Chord.Inversion)

	// Verify chord has beat notes
	assert.Len(t, chord.Notes, 6)
	for i, beatNote := range chord.Notes {
		assert.Equal(t, i, beatNote.Beat)
		assert.Equal(t, uint8(5), beatNote.Note.AccentIndex)
		assert.Equal(t, uint8(1), beatNote.Note.Ratchets.Hits)
		assert.Equal(t, uint8(0), beatNote.Note.Ratchets.Length)
		assert.Equal(t, uint8(0), beatNote.Note.Ratchets.Span)
		assert.Equal(t, grid.Action(0), beatNote.Note.Action)
		assert.Equal(t, uint8(0), beatNote.Note.GateIndex)
		assert.Equal(t, uint8(0), beatNote.Note.WaitIndex)
	}

	// Verify lines exist (25 lines from 0-24)
	assert.Len(t, readDef.lines, 25)

	// Verify accent configuration
	assert.Equal(t, uint8(15), readDef.accents.Diff)
	assert.Equal(t, uint8(120), readDef.accents.Start)
	assert.Equal(t, AccentTargetVelocity, readDef.accents.Target)
	assert.Len(t, readDef.accents.Data, 9)

	// Verify arrangement exists
	assert.NotNil(t, readDef.arrangement)
	assert.Equal(t, 1, readDef.arrangement.Iterations)
	assert.Len(t, readDef.arrangement.Nodes, 1)
	sectionNode := readDef.arrangement.Nodes[0]
	assert.Equal(t, 1, sectionNode.Iterations)
	assert.Equal(t, 0, sectionNode.Section.Part)
	assert.Equal(t, 1, sectionNode.Section.Cycles)
	assert.Equal(t, 0, sectionNode.Section.StartBeat)
	assert.Equal(t, 1, sectionNode.Section.StartCycles)
	assert.Equal(t, false, sectionNode.Section.KeepCycles)
}

func TestReadFileError(t *testing.T) {
	// Test reading from a non-existent file
	_, err := Read("/nonexistent/path/to/file.txt")
	assert.Error(t, err)
}
