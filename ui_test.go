package main

import (
	"slices"
	"testing"

	overlaykey "github.com/chriserin/seq/overlayKey"
	"github.com/stretchr/testify/assert"
)

var OK = overlaykey.InitOverlayKey

func TestGetMatchingKeys(t *testing.T) {
	testCases := []struct {
		desc       string
		keys       []overlayKey
		keyCycles  int
		expected   []overlayKey
		run        bool
		definition Definition
	}{
		{
			desc:       "A",
			keyCycles:  1,
			keys:       []overlayKey{OK(2, 1)},
			expected:   []overlayKey{},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(1, 1): {false, false}}},
		},
		{
			desc:       "B",
			keyCycles:  5,
			keys:       []overlayKey{OK(2, 1)},
			expected:   []overlayKey{},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(1, 1): {false, false}}},
		},
		{
			desc:       "C",
			keyCycles:  3,
			keys:       []overlayKey{OK(3, 1), OK(2, 1)},
			expected:   []overlayKey{OK(3, 1)},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(1, 1): {false, false}}},
		},
		{
			desc:       "D",
			keyCycles:  3,
			keys:       []overlayKey{OK(3, 2)},
			expected:   []overlayKey{OK(3, 2)},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(1, 1): {false, false}}},
		},
		{
			desc:       "E",
			keyCycles:  1,
			keys:       []overlayKey{OK(3, 2)},
			expected:   []overlayKey{},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(1, 1): {false, false}}},
		},
		{
			desc:       "F",
			keyCycles:  11,
			keys:       []overlayKey{OK(3, 8), OK(3, 4)},
			expected:   []overlayKey{OK(3, 8), OK(3, 4)},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(3, 4): {true, false}}},
		},
		{
			desc:       "G",
			keyCycles:  11,
			keys:       []overlayKey{OK(3, 8), OK(3, 4)},
			expected:   []overlayKey{OK(3, 8)},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(1, 1): {false, false}}},
		},
		{
			desc:      "H",
			keyCycles: 12,
			keys:      []overlayKey{OK(3, 1), OK(2, 1), OK(1, 1)},
			expected:  []overlayKey{OK(3, 1), OK(2, 1), OK(1, 1)},
			definition: Definition{
				metaOverlays: map[overlayKey]metaOverlay{
					OK(1, 1): {true, false},
					OK(2, 1): {true, false},
				},
			},
		},
		{
			desc:       "I",
			keyCycles:  9,
			keys:       []overlayKey{OK(3, 1), OK(2, 1), OK(1, 1)},
			expected:   []overlayKey{OK(3, 1), OK(2, 1)},
			definition: Definition{metaOverlays: map[overlayKey]metaOverlay{OK(3, 1): {false, true}}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := tC.definition.GetMatchingOverlays(tC.keyCycles, tC.keys)
			assert.Equal(t, tC.expected, result)
		})
	}
}

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
			slices.SortFunc(newSlice, overlaykey.Sort)
			assert.Equal(t, tC.expected, newSlice[0])
		})
	}
}

func TestCombinePattern(t *testing.T) {
	testCases := []struct {
		desc       string
		definition Definition
		result     overlay
		keys       []overlayKey
	}{
		{
			desc: "A",
			definition: Definition{
				overlays: overlays{
					OK(1, 1): overlay{
						{1, 1}: note{Action: 1},
					},
					OK(1, 2): overlay{
						{1, 2}: note{Action: 2},
					},
				},
			},
			keys: []overlayKey{OK(1, 2), OK(1, 1)},
			result: overlay{
				{1, 1}: note{Action: 1},
				{1, 2}: note{Action: 2},
			},
		},
		{
			desc: "B",
			definition: Definition{
				overlays: overlays{
					OK(1, 1): overlay{
						{1, 1}: note{Action: 1},
					},
					OK(2, 1): overlay{
						{1, 2}: note{Action: 2},
					},
					OK(3, 1): overlay{
						{1, 3}: note{Action: 3},
					},
				},
			},
			keys: []overlayKey{OK(3, 1), OK(2, 1), OK(1, 1)},
			result: overlay{
				{1, 1}: note{Action: 1},
				{1, 2}: note{Action: 2},
				{1, 3}: note{Action: 3},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.result, tC.definition.CombinedPattern(tC.keys))
		})
	}
}

func TestEditKeys(t *testing.T) {
	testCases := []struct {
		desc   string
		m      model
		result []overlayKey
		focus  bool
	}{
		{
			desc:   "A",
			m:      model{overlayKey: OK(1, 1)},
			result: []overlayKey{OK(1, 1)},
		},
		{
			desc: "B",
			m: model{
				overlayKey: OK(1, 1),
				definition: Definition{
					overlays: overlays{
						OK(1, 1): make(overlay),
						OK(2, 1): make(overlay),
					}}},
			result: []overlayKey{OK(1, 1)},
		},
		{
			desc: "C",
			m: model{
				overlayKey: OK(2, 1),
				definition: Definition{
					overlays: overlays{
						OK(1, 1): make(overlay),
						OK(2, 1): make(overlay),
					},
					metaOverlays: map[overlayKey]metaOverlay{
						OK(1, 1): {PressUp: true},
					},
				},
			},
			result: []overlayKey{OK(2, 1), OK(1, 1)},
		},
		{
			desc: "D",
			m: model{
				overlayKey: OK(3, 1),
				definition: Definition{
					overlays: overlays{
						OK(1, 1): make(overlay),
						OK(2, 1): make(overlay),
						OK(3, 1): make(overlay),
					}}},
			result: []overlayKey{OK(3, 1)},
		},
		{
			desc: "E",
			m: model{
				overlayKey: OK(4, 1),
				definition: Definition{
					overlays: overlays{
						// overlayKey{1, 1}: make(overlay),
						// overlayKey{2, 1}: make(overlay),
						// overlayKey{4, 1}: make(overlay),
					}}},
			result: []overlayKey{OK(4, 1)},
		},
		{
			desc: "F",
			m: model{
				overlayKey: OK(2, 1),
				definition: Definition{
					overlays: overlays{
						OK(1, 1): make(overlay),
					},
					metaOverlays: map[overlayKey]metaOverlay{
						OK(1, 1): {PressUp: true},
					},
				},
			},
			result: []overlayKey{OK(2, 1), OK(1, 1)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if true || tC.focus {
				keys := tC.m.EditKeys()
				assert.Equal(t, tC.result, keys)
			}
		})
	}
}
