package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInBounds(t *testing.T) {
	testCases := []struct {
		desc   string
		bounds Bounds
		value  gridKey
		result bool
	}{
		{
			desc:   "A",
			bounds: Bounds{1, 1, 1, 1},
			value:  gridKey{1, 1},
			result: true,
		},
		{
			desc:   "B",
			bounds: Bounds{0, 1, 1, 0},
			value:  gridKey{1, 1},
			result: true,
		},
		{
			desc:   "C",
			bounds: Bounds{0, 0, 1, 1},
			value:  gridKey{1, 1},
			result: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.result, tC.bounds.InBounds(tC.value))
		})
	}
}
