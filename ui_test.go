package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMatchingKeys(t *testing.T) {
	testCases := []struct {
		desc      string
		keys      []overlayKey
		keyCycles int
		expected  []overlayKey
		run       bool
		model     model
	}{
		{
			desc:      "A",
			keyCycles: 1,
			keys:      []overlayKey{{2, 1}},
			expected:  []overlayKey{},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{}},
		},
		{
			desc:      "B",
			keyCycles: 5,
			keys:      []overlayKey{{2, 1}},
			expected:  []overlayKey{},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{}},
		},
		{
			desc:      "C",
			keyCycles: 3,
			keys:      []overlayKey{{3, 1}, {2, 1}},
			expected:  []overlayKey{{3, 1}},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{}},
		},
		{
			desc:      "D",
			keyCycles: 3,
			keys:      []overlayKey{{3, 2}},
			expected:  []overlayKey{{3, 2}},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{}},
		},
		{
			desc:      "E",
			keyCycles: 1,
			keys:      []overlayKey{{3, 2}},
			expected:  []overlayKey{},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{}},
		},
		{
			desc:      "F",
			keyCycles: 11,
			keys:      []overlayKey{{3, 8}, {3, 4}},
			expected:  []overlayKey{{3, 8}, {3, 4}},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{{3, 4}}},
		},
		{
			desc:      "G",
			keyCycles: 11,
			keys:      []overlayKey{{3, 8}, {3, 4}},
			expected:  []overlayKey{{3, 8}},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{}},
		},
		{
			desc:      "H",
			keyCycles: 12,
			keys:      []overlayKey{{3, 1}, {2, 1}, {1, 1}},
			expected:  []overlayKey{{3, 1}, {2, 1}, {1, 1}},
			model:     model{pressedDownKeys: []overlayKey{}, stackedupKeys: []overlayKey{{1, 1}, {2, 1}}},
		},
		{
			desc:      "I",
			keyCycles: 9,
			keys:      []overlayKey{{3, 1}, {2, 1}, {1, 1}},
			expected:  []overlayKey{{3, 1}, {2, 1}},
			model:     model{pressedDownKeys: []overlayKey{{3, 1}}, stackedupKeys: []overlayKey{}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := tC.model.GetMatchingOverlays(tC.keyCycles, tC.keys)
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
			keys:     []overlayKey{{2, 1}, {3, 1}},
			expected: overlayKey{3, 1},
		},
		{
			desc:     "TEST B",
			keys:     []overlayKey{{3, 1}, {2, 1}},
			expected: overlayKey{3, 1},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			newSlice := make([]overlayKey, len(tC.keys))
			copy(newSlice, tC.keys)
			slices.SortFunc(newSlice, OverlayKeySort)
			assert.Equal(t, tC.expected, newSlice[0])
		})
	}
}

func TestCombinePattern(t *testing.T) {
	testCases := []struct {
		desc   string
		model  model
		result overlay
		keys   []overlayKey
	}{
		{
			desc: "A",
			model: model{
				overlays: overlays{
					{1, 1}: overlay{
						{1, 1}: note{action: 1},
					},
					{1, 2}: overlay{
						{1, 2}: note{action: 2},
					},
				},
			},
			keys: []overlayKey{{1, 2}, {1, 1}},
			result: overlay{
				{1, 1}: note{action: 1},
				{1, 2}: note{action: 2},
			},
		},
		{
			desc: "B",
			model: model{
				overlays: overlays{
					{1, 1}: overlay{
						{1, 1}: note{action: 1},
					},
					{2, 1}: overlay{
						{1, 2}: note{action: 2},
					},
					{3, 1}: overlay{
						{1, 3}: note{action: 3},
					},
				},
			},
			keys: []overlayKey{{3, 1}, {2, 1}, {1, 1}},
			result: overlay{
				{1, 1}: note{action: 1},
				{1, 2}: note{action: 2},
				{1, 3}: note{action: 3},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.result, tC.model.CombinedPattern(tC.keys))
		})
	}
}
