package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc      string
		keys      []overlayKey
		keyCycles int
		expected  []overlayKey
	}{
		{
			desc:      "",
			keyCycles: 1,
			keys:      []overlayKey{{2, 1}},
			expected:  []overlayKey{},
		},
		{
			desc:      "",
			keyCycles: 5,
			keys:      []overlayKey{{2, 1}},
			expected:  []overlayKey{},
		},
		{
			desc:      "",
			keyCycles: 3,
			keys:      []overlayKey{{3, 1}, {2, 1}},
			expected:  []overlayKey{{3, 1}},
		},
		{
			desc:      "",
			keyCycles: 3,
			keys:      []overlayKey{{3, 2}},
			expected:  []overlayKey{{3, 2}},
		},
		{
			desc:      "",
			keyCycles: 1,
			keys:      []overlayKey{{3, 2}},
			expected:  []overlayKey{},
		},
		{
			desc:      "",
			keyCycles: 11,
			keys:      []overlayKey{{3, 8}, {3, 4}},
			expected:  []overlayKey{{3, 8}, {3, 4}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			result := GetMatchingOverlays(tC.keyCycles, tC.keys)
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
