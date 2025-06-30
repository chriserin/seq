package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/grid"
	"github.com/stretchr/testify/assert"
)

func TestWriteAdditionalProperties(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "seq-write-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("Write with full definition", func(t *testing.T) {
		// Create line definitions
		lines := []grid.LineDefinition{
			{Channel: 1, Note: 60, MsgType: 0}, // C4 on channel 1
			{Channel: 2, Note: 67, MsgType: 0}, // G4 on channel 2
		}

		// Create accents
		accents := patternAccents{
			Diff:   5,
			Start:  10,
			Target: AccentTargetVelocity,
			Data: []config.Accent{
				{Value: 100},
				{Value: 80},
			},
		}

		// Create a model with all definition fields populated
		m := &model{
			definition: Definition{
				parts: &[]arrangement.Part{
					{
						Name:  "TestPart",
						Beats: 16,
					},
				},
				lines:           lines,
				tempo:           120,
				subdivisions:    4,
				keyline:         2,
				accents:         accents,
				instrument:      "piano",
				template:        "default",
				templateUIStyle: "dark",
			},
		}

		// Create test file path
		filename := filepath.Join(tempDir, "full_definition.txt")

		// Call Write function
		err := Write(m, filename)
		assert.NoError(t, err)

		// Read file content
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)

		// Check content
		contentStr := string(content)

		// Check global settings
		assert.Contains(t, contentStr, "GLOBAL SETTINGS")
		assert.Contains(t, contentStr, "Tempo: 120")
		assert.Contains(t, contentStr, "Subdivisions: 4")
		assert.Contains(t, contentStr, "Keyline: 2")
		assert.Contains(t, contentStr, "Instrument: piano")
		assert.Contains(t, contentStr, "Template: default")
		assert.Contains(t, contentStr, "TemplateUIStyle: dark")

		// Check lines
		assert.Contains(t, contentStr, "LINES")
		assert.Contains(t, contentStr, "Line 0: Channel=1, Note=60")
		assert.Contains(t, contentStr, "Line 1: Channel=2, Note=67")

		// Check accents
		assert.Contains(t, contentStr, "ACCENTS")
		assert.Contains(t, contentStr, "Diff: 5")
		assert.Contains(t, contentStr, "Start: 10")
		assert.Contains(t, contentStr, "Target: VELOCITY")
		assert.Contains(t, contentStr, "ACCENT DATA")
		assert.Contains(t, contentStr, "Accent 0:")
		assert.Contains(t, contentStr, "Value=100")
		assert.Contains(t, contentStr, "Accent 1:")
		assert.Contains(t, contentStr, "Value=80")

		// Check parts
		assert.Contains(t, contentStr, "PARTS")
		assert.Contains(t, contentStr, "PART TestPart")
	})
}
