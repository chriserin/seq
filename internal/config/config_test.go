package config

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	"github.com/stretchr/testify/assert"
)

func TestProcessConfig(t *testing.T) {
	t.Run("creates templates", func(t *testing.T) {
		ProcessConfig("./testdata/AddTemplate.lua")
		template := GetTemplate("Drums")
		assert.Equal(t, "Drums", template.Name)
		assert.Equal(t, 8, len(template.Lines))
		assert.Equal(t, uint8(41), template.Lines[5].Note)
		assert.Equal(t, grid.MESSAGE_TYPE_NOTE, template.Lines[4].MsgType)
		assert.Equal(t, "plain", template.UIStyle)
		assert.Equal(t, 32, template.MaxGateLength)
	})

	t.Run("adds instruments", func(t *testing.T) {
		ProcessConfig("./testdata/AddInstrument.lua")
		instrument := GetInstrument("Prophet 10")
		assert.Equal(t, "Prophet 10", instrument.Name)
		assert.Equal(t, 59, len(instrument.CCs))
		assert.Equal(t, "GLIDE RATE", instrument.CCs[11].Name)
		assert.Equal(t, uint8(26), instrument.CCs[11].Value)
		assert.Equal(t, uint8(120), instrument.CCs[11].UpperLimit)
	})
}
