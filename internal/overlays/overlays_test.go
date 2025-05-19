package overlays

import (
	"testing"

	"github.com/chriserin/seq/internal/grid"
	overlaykey "github.com/chriserin/seq/internal/overlaykey"
	"github.com/stretchr/testify/assert"
)

var secondKey = Key{
	Shift:      2,
	Interval:   1,
	Width:      0,
	StartCycle: 0,
}

var thirdKey = Key{
	Shift:      3,
	Interval:   1,
	Width:      0,
	StartCycle: 0,
}

func TestAddOverlays(t *testing.T) {
	t.Run("None added", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		keys := []Key{}
		overlay.CollectKeys(&keys)
		assert.Equal(t, keys, []Key{overlaykey.ROOT})
	})

	t.Run("One Higher key Added", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		keys := []Key{}
		newOverlay.CollectKeys(&keys)
		assert.Equal(t, keys, []Key{secondKey, overlaykey.ROOT})
	})

	t.Run("One Lower key Added", func(t *testing.T) {
		overlay := InitOverlay(thirdKey, nil)
		newOverlay := overlay.Add(secondKey)
		keys := []Key{}
		newOverlay.CollectKeys(&keys)
		assert.Equal(t, keys, []Key{thirdKey, secondKey})
	})

	t.Run("One key added in between", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(thirdKey)
		newOverlay = newOverlay.Add(secondKey)
		keys := []Key{}
		newOverlay.CollectKeys(&keys)
		assert.Equal(t, keys, []Key{thirdKey, secondKey, overlaykey.ROOT})
	})

	t.Run("Three keys added in order", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		nextOverlay := newOverlay.Add(thirdKey)
		assert.Equal(t, nextOverlay.Below.Key, secondKey)
	})
}

func TestFindAboveOverlay(t *testing.T) {
	t.Run("Chords len should be the same", func(t *testing.T) {
		overlayA := InitOverlay(overlaykey.InitOverlayKey(1, 1), nil)
		overlayB := InitOverlay(overlaykey.InitOverlayKey(2, 1), overlayA)
		overlayB.CreateChord(grid.GridKey{Line: 1, Beat: 1}, 0)
		assert.Equal(t, 1, len(overlayB.Chords))
		resultOverlay := overlayB.FindAboveOverlay(overlaykey.InitOverlayKey(2, 1))
		assert.Equal(t, 1, len(resultOverlay.Chords))
	})
}

func TestHighestMatchingOverlay(t *testing.T) {
	t.Run("Get Highest of 1", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		highest := overlay.HighestMatchingOverlay(1)
		key := (*highest).Key
		assert.Equal(t, overlaykey.ROOT, key)
	})

	t.Run("Get Highest of 2", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		highest := newOverlay.HighestMatchingOverlay(2)
		key := (*highest).Key
		assert.Equal(t, secondKey, key)
	})

	t.Run("Get Highest matching of 2 when highest doesn't match", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		highest := newOverlay.HighestMatchingOverlay(3)
		key := (*highest).Key
		assert.Equal(t, overlaykey.ROOT, key)
	})

	t.Run("Get Highest matching of 3 when highest doesn't match", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		newOverlay = newOverlay.Add(thirdKey)
		highest := newOverlay.HighestMatchingOverlay(2)
		key := highest.Key
		assert.Equal(t, secondKey, key)
	})
}

func TestRemoveOverlay(t *testing.T) {
	t.Run("Remove top overlay", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		nextOverlay := newOverlay.Add(thirdKey)
		reducedOverlay := nextOverlay.Remove(thirdKey)

		keys := []Key{}
		reducedOverlay.CollectKeys(&keys)
		assert.Equal(t, keys, []Key{secondKey, overlaykey.ROOT})
	})

	t.Run("Remove middle overlay", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		nextOverlay := newOverlay.Add(thirdKey)
		reducedOverlay := nextOverlay.Remove(secondKey)

		keys := []Key{}
		reducedOverlay.CollectKeys(&keys)
		assert.Equal(t, keys, []Key{thirdKey, overlaykey.ROOT})
	})

	t.Run("Remove bottom", func(t *testing.T) {
		overlay := InitOverlay(overlaykey.ROOT, nil)
		newOverlay := overlay.Add(secondKey)
		nextOverlay := newOverlay.Add(thirdKey)
		reducedOverlay := nextOverlay.Remove(overlaykey.ROOT)

		keys := []Key{}
		reducedOverlay.CollectKeys(&keys)
		assert.Equal(t, keys, []Key{thirdKey, secondKey})
	})
}
