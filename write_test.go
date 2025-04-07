package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlays"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "seq-write-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("Write empty parts", func(t *testing.T) {
		// Create a model with empty parts
		m := &model{
			definition: Definition{
				parts: &[]arrangement.Part{},
			},
		}

		// Create test file path
		filename := filepath.Join(tempDir, "empty_parts.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)

		// Check content
		assert.Empty(t, content, "File should be empty for empty parts")
	})

	t.Run("Write nil parts", func(t *testing.T) {
		// Create a model with nil parts
		m := &model{
			definition: Definition{
				parts: nil,
			},
		}

		// Create test file path
		filename := filepath.Join(tempDir, "nil_parts.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Check if file exists - it should be created but empty
		_, err = os.Stat(filename)
		assert.NoError(t, err, "File should be created even with nil parts")

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)
		assert.Empty(t, content, "File should be empty for nil parts")
	})

	t.Run("Write with one part and no overlays", func(t *testing.T) {
		// Create a model with one part
		m := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:     "TestPart",
						Beats:    16,
						Overlays: nil,
					},
				},
			},
		}

		// Create test file path
		filename := filepath.Join(tempDir, "one_part.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)

		// Check if content contains part name and beats
		contentStr := string(content)
		assert.Contains(t, contentStr, "PARTS")
		assert.Contains(t, contentStr, "PART TestPart")
		assert.Contains(t, contentStr, "Name: TestPart")
		assert.Contains(t, contentStr, "Beats: 16")
		assert.NotContains(t, contentStr, "OVERLAY")
	})

	t.Run("Write with parts and overlays", func(t *testing.T) {
		// Create overlay key
		key := overlaykey.OverlayPeriodicity{
			Shift:      2,
			Interval:   4,
			Width:      1,
			StartCycle: 0,
		}

		// Create overlay with notes
		overlay := overlays.InitOverlay(key, nil)
		note := grid.InitNote()
		note.AccentIndex = 3
		note.GateIndex = 2
		gridKey := grid.GridKey{Line: 1, Beat: 2}
		overlay.SetNote(gridKey, note)

		// Create a model with part containing overlay
		m := &model{
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

		// Create test file path
		filename := filepath.Join(tempDir, "part_with_overlay.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)

		// Check content
		contentStr := string(content)
		assert.Contains(t, contentStr, "PARTS")
		assert.Contains(t, contentStr, "PART PartWithOverlay")
		assert.Contains(t, contentStr, "Name: PartWithOverlay")
		assert.Contains(t, contentStr, "Beats: 8")
		assert.Contains(t, contentStr, "OVERLAY")
		assert.Contains(t, contentStr, "Shift: 2")
		assert.Contains(t, contentStr, "Interval: 4")
		assert.Contains(t, contentStr, "Width: 1")
		assert.Contains(t, contentStr, "NOTES")
		assert.Contains(t, contentStr, "GridKey(1,2)")
		assert.Contains(t, contentStr, "AccentIndex=3")
	})

	t.Run("Write with multiple parts and nested overlays", func(t *testing.T) {
		// Create first key and overlay
		key1 := overlaykey.OverlayPeriodicity{
			Shift:      1,
			Interval:   4,
			Width:      0,
			StartCycle: 0,
		}
		
		// Create second key and overlay
		key2 := overlaykey.OverlayPeriodicity{
			Shift:      2,
			Interval:   8,
			Width:      0,
			StartCycle: 0,
		}
		
		// Create the nested overlay structure
		overlay2 := overlays.InitOverlay(key2, nil)
		note2 := grid.InitNote()
		note2.AccentIndex = 4
		gridKey2 := grid.GridKey{Line: 2, Beat: 3}
		overlay2.SetNote(gridKey2, note2)
		
		overlay1 := overlays.InitOverlay(key1, overlay2)
		note1 := grid.InitNote()
		note1.AccentIndex = 2
		gridKey1 := grid.GridKey{Line: 1, Beat: 1}
		overlay1.SetNote(gridKey1, note1)
		
		// Create a model with multiple parts
		m := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:     "Part1",
						Beats:    8,
						Overlays: overlay1,
					},
					{
						Name:     "Part2",
						Beats:    4,
						Overlays: nil,
					},
				},
			},
		}

		// Create test file path
		filename := filepath.Join(tempDir, "multiple_parts.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)

		// Check content
		contentStr := string(content)
		
		// Check for part 1
		assert.Contains(t, contentStr, "PART Part1")
		assert.Contains(t, contentStr, "Name: Part1")
		assert.Contains(t, contentStr, "Beats: 8")
		
		// Check overlay hierarchy
		count := strings.Count(contentStr, "OVERLAY")
		assert.Equal(t, 2, count, "Should have two overlay sections")
		
		// Check for first overlay
		assert.Contains(t, contentStr, "Shift: 1")
		assert.Contains(t, contentStr, "Interval: 4")
		assert.Contains(t, contentStr, "GridKey(1,1)")
		
		// Check for second overlay
		assert.Contains(t, contentStr, "Shift: 2")
		assert.Contains(t, contentStr, "Interval: 8")
		assert.Contains(t, contentStr, "GridKey(2,3)")
		
		// Check for part 2
		assert.Contains(t, contentStr, "PART Part2")
		assert.Contains(t, contentStr, "Name: Part2")
		assert.Contains(t, contentStr, "Beats: 4")
	})
}

// TestWriteFileError tests error handling when file creation fails
func TestWriteFileError(t *testing.T) {
	// Create a model
	m := &model{
		definition: Definition{
			parts: &[]arrangement.Part{
				{
					Name:  "TestPart",
					Beats: 4,
				},
			},
		},
	}

	// Use an invalid file path
	filename := "/invalid/path/that/does/not/exist/file.txt"

	// Call Write function and expect an error
	err := Write(m, filename)
	assert.Error(t, err, "Write should return an error for invalid file path")
}