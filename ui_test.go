package main

import (
	"slices"
	"testing"

	overlaykey "github.com/chriserin/seq/internal/overlaykey"
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
