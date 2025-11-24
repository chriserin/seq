package grid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateEuclideanRhythm(t *testing.T) {
	testCases := []struct {
		desc     string
		hits     int
		steps    int
		expected []bool
	}{
		{
			desc:     "3 hits in 8 steps (Cuban tresillo)",
			hits:     3,
			steps:    8,
			expected: []bool{true, false, false, true, false, false, true, false},
		},
		{
			desc:     "5 hits in 8 steps (Cuban cinquillo)",
			hits:     5,
			steps:    8,
			expected: []bool{true, false, true, true, false, true, true, false},
		},
		{
			desc:     "4 hits in 16 steps (basic four on the floor)",
			hits:     4,
			steps:    16,
			expected: []bool{true, false, false, false, true, false, false, false, true, false, false, false, true, false, false, false},
		},
		{
			desc:     "5 hits in 12 steps",
			hits:     5,
			steps:    12,
			expected: []bool{true, false, false, true, false, true, false, false, true, false, true, false},
		},
		{
			desc:     "7 hits in 16 steps",
			hits:     7,
			steps:    16,
			expected: []bool{true, false, false, true, false, true, false, true, false, false, true, false, true, false, true, false},
		},
		{
			desc:     "2 hits in 5 steps",
			hits:     2,
			steps:    5,
			expected: []bool{true, false, true, false, false},
		},
		{
			desc:     "1 hit in 4 steps",
			hits:     1,
			steps:    4,
			expected: []bool{true, false, false, false},
		},
		{
			desc:     "4 hits in 4 steps (all hits)",
			hits:     4,
			steps:    4,
			expected: []bool{true, true, true, true},
		},
		{
			desc:     "0 hits in 8 steps (empty pattern)",
			hits:     0,
			steps:    8,
			expected: []bool{false, false, false, false, false, false, false, false},
		},
		{
			desc:     "hits greater than steps (invalid)",
			hits:     10,
			steps:    8,
			expected: []bool{false, false, false, false, false, false, false, false},
		},
		{
			desc:     "negative hits (invalid)",
			hits:     -1,
			steps:    8,
			expected: []bool{false, false, false, false, false, false, false, false},
		},
		{
			desc:     "zero steps (invalid)",
			hits:     3,
			steps:    0,
			expected: []bool{},
		},
		{
			desc:     "3 hits in 4 steps",
			hits:     3,
			steps:    4,
			expected: []bool{true, true, true, false},
		},
		{
			desc:     "5 hits in 6 steps",
			hits:     5,
			steps:    6,
			expected: []bool{true, true, true, true, true, false},
		},
		{
			desc:     "2 hits in 3 steps",
			hits:     2,
			steps:    3,
			expected: []bool{true, true, false},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := GenerateEuclideanRhythm(tC.hits, tC.steps)
			assert.Equal(t, tC.expected, result, "Pattern should match expected Euclidean rhythm")

			// Verify the length is correct
			assert.Equal(t, tC.steps, len(result), "Result length should match steps")

			// For valid inputs, verify hit count
			if tC.hits > 0 && tC.steps > 0 && tC.hits <= tC.steps {
				hitCount := 0
				for _, hit := range result {
					if hit {
						hitCount++
					}
				}
				assert.Equal(t, tC.hits, hitCount, "Number of hits should match requested hits")
			}
		})
	}
}

func TestGenerateEuclideanRhythm_Properties(t *testing.T) {
	t.Run("Result length equals steps", func(t *testing.T) {
		testCases := []struct {
			hits  int
			steps int
		}{
			{3, 8},
			{5, 12},
			{7, 16},
			{2, 7},
		}

		for _, tC := range testCases {
			result := GenerateEuclideanRhythm(tC.hits, tC.steps)
			assert.Equal(t, tC.steps, len(result), "Result length should always equal steps")
		}
	})

	t.Run("Correct number of hits", func(t *testing.T) {
		testCases := []struct {
			hits  int
			steps int
		}{
			{3, 8},
			{5, 12},
			{7, 16},
			{2, 7},
			{4, 4},
		}

		for _, tC := range testCases {
			result := GenerateEuclideanRhythm(tC.hits, tC.steps)
			hitCount := 0
			for _, hit := range result {
				if hit {
					hitCount++
				}
			}
			assert.Equal(t, tC.hits, hitCount, "Number of hits should match requested hits")
		}
	})

	t.Run("First element is always true for valid inputs", func(t *testing.T) {
		testCases := []struct {
			hits  int
			steps int
		}{
			{1, 8},
			{3, 8},
			{5, 12},
			{7, 16},
		}

		for _, tC := range testCases {
			result := GenerateEuclideanRhythm(tC.hits, tC.steps)
			assert.True(t, result[0], "First element should be true for valid Euclidean rhythms")
		}
	})
}
