package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
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

		// Check if content contains all required sections
		contentStr := string(content)
		assert.Contains(t, contentStr, "GLOBAL SETTINGS")
		assert.Contains(t, contentStr, "ACCENTS")
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
		overlay.PressDown = false
		overlay.PressUp = true
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

		// Check content for all sections
		contentStr := string(content)
		assert.Contains(t, contentStr, "PARTS")
		assert.Contains(t, contentStr, "PART PartWithOverlay")
		assert.Contains(t, contentStr, "Name: PartWithOverlay")
		assert.Contains(t, contentStr, "Beats: 8")
		assert.Contains(t, contentStr, "OVERLAY")
		assert.Contains(t, contentStr, "Shift: 2")
		assert.Contains(t, contentStr, "Interval: 4")
		assert.Contains(t, contentStr, "Width: 1")
		assert.Contains(t, contentStr, "PressUp: true")
		assert.Contains(t, contentStr, "PressDown: false")
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

	t.Run("Write with arrangement", func(t *testing.T) {
		// Create an arrangement tree
		// Create a leaf node (section/part node)
		section := arrangement.InitSongSection(0)
		section.Cycles = 2
		section.StartBeat = 1

		leaf := &arrangement.Arrangement{
			Section:    section,
			Iterations: 1,
		}

		// Create a group node with the leaf as a child
		group := &arrangement.Arrangement{
			Nodes:      []*arrangement.Arrangement{leaf},
			Iterations: 3,
		}

		// Create root node with the group as a child
		root := &arrangement.Arrangement{
			Nodes:      []*arrangement.Arrangement{group},
			Iterations: 1,
		}

		// Create a model with parts and arrangement
		m := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:     "TestPart",
						Beats:    16,
						Overlays: nil,
					},
				},
				arrangement: root,
			},
		}

		// Create test file path
		filename := filepath.Join(tempDir, "with_arrangement.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)

		// Check content
		contentStr := string(content)

		// Check both parts and arrangement sections exist
		assert.Contains(t, contentStr, "PARTS")
		assert.Contains(t, contentStr, "ARRANGEMENT")

		// Check arrangement structure
		assert.Contains(t, contentStr, "ROOT NODE")
		assert.Contains(t, contentStr, "GROUP #")
		assert.Contains(t, contentStr, "SECTION #")

		// Check iterations
		assert.Contains(t, contentStr, "Iterations: 3")

		// Check section properties
		assert.Contains(t, contentStr, "Cycles: 2")
		assert.Contains(t, contentStr, "StartBeat: 1")
		assert.Contains(t, contentStr, "Part: 0")
	})

	t.Run("Write with complex arrangement tree", func(t *testing.T) {
		// Create a complex arrangement tree with multiple sections and groups
		// Create leaf nodes (section/part nodes)
		section1 := arrangement.InitSongSection(0)
		section1.Cycles = 2

		section2 := arrangement.InitSongSection(1)
		section2.Cycles = 1
		section2.StartBeat = 4

		section3 := arrangement.InitSongSection(1)
		section3.Cycles = 3
		section3.KeepCycles = true

		leaf1 := &arrangement.Arrangement{
			Section:    section1,
			Iterations: 1,
		}

		leaf2 := &arrangement.Arrangement{
			Section:    section2,
			Iterations: 1,
		}

		leaf3 := &arrangement.Arrangement{
			Section:    section3,
			Iterations: 1,
		}

		// Create a nested group structure
		group1 := &arrangement.Arrangement{
			Nodes:      []*arrangement.Arrangement{leaf1, leaf2},
			Iterations: 2,
		}

		group2 := &arrangement.Arrangement{
			Nodes:      []*arrangement.Arrangement{group1, leaf3},
			Iterations: 4,
		}

		// Create root node
		root := &arrangement.Arrangement{
			Nodes:      []*arrangement.Arrangement{group2},
			Iterations: 1,
		}

		// Create a model with parts and arrangement
		m := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:  "Part1",
						Beats: 16,
					},
					{
						Name:  "Part2",
						Beats: 8,
					},
				},
				arrangement: root,
			},
		}

		// Create test file path
		filename := filepath.Join(tempDir, "complex_arrangement.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)

		// Check content
		contentStr := string(content)

		// Check parts section
		assert.Contains(t, contentStr, "PARTS")
		assert.Contains(t, contentStr, "PART Part1")
		assert.Contains(t, contentStr, "PART Part2")

		// Check arrangement structure hierarchy
		assert.Contains(t, contentStr, "ARRANGEMENT")
		assert.Contains(t, contentStr, "ROOT NODE")

		// Count the number of GROUP and PART nodes
		groupCount := strings.Count(contentStr, "GROUP #")
		partCount := strings.Count(contentStr, "SECTION #")
		assert.Equal(t, 2, groupCount, "Should have 2 group nodes")
		assert.Equal(t, 3, partCount, "Should have 3 part nodes")

		// Check iterations
		assert.Contains(t, contentStr, "Iterations: 2")
		assert.Contains(t, contentStr, "Iterations: 4")

		// Check various section properties
		assert.Contains(t, contentStr, "Cycles: 2")
		assert.Contains(t, contentStr, "Cycles: 1")
		assert.Contains(t, contentStr, "Cycles: 3")
		assert.Contains(t, contentStr, "StartBeat: 4")
		assert.Contains(t, contentStr, "KeepCycles: true")
		assert.Contains(t, contentStr, "KeepCycles: false")

		definition, err := Read(filename)
		if err == nil {
			assert.Equal(t, 1, len(definition.arrangement.Nodes))
			currentNode := definition.arrangement.Nodes[0]
			assert.Equal(t, 2, len(currentNode.Nodes))
			currentNode = currentNode.Nodes[0]
			assert.Equal(t, 2, len(currentNode.Nodes))
		} else {
			t.Fail()
		}
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
