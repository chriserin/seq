package overlaykey

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatches(t *testing.T) {
	testCases := []struct {
		desc    string
		op      OverlayPeriodicity
		results []bool
	}{
		{
			desc:    "Root",
			op:      OverlayPeriodicity{1, 1, 1, 0},
			results: []bool{true, true, true, true, true, true, true, true, true, true},
		},
		{
			desc:    "Every Other",
			op:      OverlayPeriodicity{2, 1, 1, 0},
			results: []bool{false, true, false, true, false, true, false, true, false, true},
		},
		{
			desc:    "Every Other And starts at 5",
			op:      OverlayPeriodicity{2, 1, 1, 5},
			results: []bool{false, false, false, false, false, true, false, true, false, true},
		},
		{
			desc:    "First of every two",
			op:      OverlayPeriodicity{1, 2, 1, 0},
			results: []bool{true, false, true, false, true, false, true, false, true, false},
		},
		{
			desc:    "Every fifth",
			op:      OverlayPeriodicity{5, 1, 1, 0},
			results: []bool{false, false, false, false, true, false, false, false, false, true},
		},
		{
			desc:    "Every fifth plus width",
			op:      OverlayPeriodicity{5, 1, 2, 0},
			results: []bool{false, false, false, false, true, true, false, false, false, true},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			results := make([]bool, 10)
			for i := range results {
				results[i] = tC.op.DoesMatch(i + 1)
			}

			assert.Equal(t, tC.results, results)
		})
	}
}

func TestSort(t *testing.T) {
	testCases := []struct {
		desc   string
		input  []OverlayPeriodicity
		output []OverlayPeriodicity
	}{
		{
			desc:   "Sorts based on interval",
			input:  []OverlayPeriodicity{{0, 1, 0, 0}, {0, 2, 0, 0}},
			output: []OverlayPeriodicity{{0, 2, 0, 0}, {0, 1, 0, 0}},
		},
		{
			desc:   "Sorts based on shift",
			input:  []OverlayPeriodicity{{2, 4, 0, 0}, {3, 4, 0, 0}},
			output: []OverlayPeriodicity{{3, 4, 0, 0}, {2, 4, 0, 0}},
		},
		{
			desc:   "Sorts based on width",
			input:  []OverlayPeriodicity{{3, 7, 2, 0}, {3, 7, 1, 0}},
			output: []OverlayPeriodicity{{3, 7, 1, 0}, {3, 7, 2, 0}},
		},
		{
			desc:   "Sorts based on start",
			input:  []OverlayPeriodicity{{1, 2, 0, 0}, {1, 2, 0, 7}},
			output: []OverlayPeriodicity{{1, 2, 0, 7}, {1, 2, 0, 0}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			slices.SortFunc(tC.input, Compare)
			assert.Equal(t, tC.output, tC.input)
		})
	}
}
