package sequence

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

	t.Run("Write with full sequence", func(t *testing.T) {
		// Create line definitions
		lines := []grid.LineDefinition{
			{Channel: 1, Note: 60, MsgType: 0}, // C4 on channel 1
			{Channel: 2, Note: 67, MsgType: 0}, // G4 on channel 2
		}

		// Create accents
		accents := PatternAccents{
			Target: AccentTargetVelocity,
			Start:  10,
			End:    5,
			Data: []config.Accent{
				{Value: 100},
				{Value: 80},
			},
		}

		// Create a model with all sequence fields populated
		sequence := Sequence{
			Parts: &[]arrangement.Part{
				{
					Name:  "TestPart",
					Beats: 16,
				},
			},
			Lines:           lines,
			Tempo:           120,
			Subdivisions:    4,
			Keyline:         2,
			Accents:         accents,
			Instrument:      "piano",
			Template:        "default",
			TemplateUIStyle: "dark",
		}

		// Create test file path
		filename := filepath.Join(tempDir, "full_sequence.txt")

		// Call Write function
		err := Write(sequence, filename)
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
		assert.Contains(t, contentStr, "Target: VELOCITY")
		assert.Contains(t, contentStr, "Start: 10")
		assert.Contains(t, contentStr, "End: 5")
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
