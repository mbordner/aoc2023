package ranges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Subtract(t *testing.T) {
	testCases := []struct {
		name     string
		orig     []int
		sub      []int
		expected []int
	}{
		{
			name:     "test1",
			orig:     []int{0, 10},
			sub:      []int{4, 6},
			expected: []int{0, 3, 7, 10},
		},
		{
			name:     "test2",
			orig:     []int{0, 4, 6, 8, 10, 12},
			sub:      []int{4, 6},
			expected: []int{0, 3, 7, 8, 10, 12},
		},
		{
			name:     "test3",
			orig:     []int{0, 4, 6, 8, 10, 12},
			sub:      []int{4, 6, 11, 12},
			expected: []int{0, 3, 7, 8, 10, 10},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Collection[int]{}
			for i := 0; i < len(tc.orig); i += 2 {
				_, _ = c.Add(tc.orig[i], tc.orig[i+1])
			}
			s := &Collection[int]{}
			for i := 0; i < len(tc.sub); i += 2 {
				_, _ = s.Add(tc.sub[i], tc.sub[i+1])
			}
			results := c.Subtract(s)
			assert.True(t, ArrayEqual[int](results.ValuePairs(), tc.expected))
		})
	}
}

func Test_Remove(t *testing.T) {
	testCases := []struct {
		name           string
		pairs          [][2]int
		removes        [][2]int
		expectedValues [][]int
		expectedLength int
	}{
		{
			name: "no overlap",
			pairs: [][2]int{
				{0, 10},
			},
			removes: [][2]int{
				{12, 14},
			},
			expectedValues: [][]int{
				{0, 10},
			},
			expectedLength: 11,
		},
		{
			name: "overlap 1 start",
			pairs: [][2]int{
				{0, 5},
				{15, 20},
			},
			removes: [][2]int{
				{4, 8},
			},
			expectedValues: [][]int{
				{0, 3, 15, 20},
			},
			expectedLength: 10,
		},
		{
			name: "overlap 1 end",
			pairs: [][2]int{
				{5, 10},
				{15, 20},
			},
			removes: [][2]int{
				{4, 8},
			},
			expectedValues: [][]int{
				{9, 10, 15, 20},
			},
			expectedLength: 8,
		},
		{
			name: "fully overlap 1a",
			pairs: [][2]int{
				{0, 5},
				{15, 20},
			},
			removes: [][2]int{
				{0, 5},
			},
			expectedValues: [][]int{
				{15, 20},
			},
			expectedLength: 6,
		},
		{
			name: "fully overlap 1b",
			pairs: [][2]int{
				{0, 5},
				{15, 20},
			},
			removes: [][2]int{
				{15, 20},
			},
			expectedValues: [][]int{
				{0, 5},
			},
			expectedLength: 6,
		},
		{
			name: "fully overlap 1c",
			pairs: [][2]int{
				{0, 5},
				{15, 20},
				{25, 30},
			},
			removes: [][2]int{
				{15, 20},
			},
			expectedValues: [][]int{
				{0, 5, 25, 30},
			},
			expectedLength: 12,
		},
		{
			name: "fully overlap without removal 1",
			pairs: [][2]int{
				{0, 5},
				{15, 20},
			},
			removes: [][2]int{
				{3, 5},
			},
			expectedValues: [][]int{
				{0, 2, 15, 20},
			},
			expectedLength: 9,
		},
		{
			name: "fully overlap without removal 2",
			pairs: [][2]int{
				{0, 5},
				{15, 20},
			},
			removes: [][2]int{
				{0, 3},
			},
			expectedValues: [][]int{
				{4, 5, 15, 20},
			},
			expectedLength: 8,
		},
		{
			name: "fully overlap without removal 3",
			pairs: [][2]int{
				{0, 5},
				{15, 20},
			},
			removes: [][2]int{
				{4, 4},
			},
			expectedValues: [][]int{
				{0, 3, 5, 5, 15, 20},
			},
			expectedLength: 11,
		},
		{
			name: "overlap 2",
			pairs: [][2]int{
				{0, 5},
				{7, 10},
			},
			removes: [][2]int{
				{4, 8},
			},
			expectedValues: [][]int{
				{0, 3, 9, 10},
			},
			expectedLength: 6,
		},
		{
			name: "overlap 2b",
			pairs: [][2]int{
				{0, 5},
				{7, 10},
				{15, 20},
			},
			removes: [][2]int{
				{4, 17},
			},
			expectedValues: [][]int{
				{0, 3, 18, 20},
			},
			expectedLength: 7,
		},
		{
			name: "overlap all",
			pairs: [][2]int{
				{5, 10},
				{15, 20},
			},
			removes: [][2]int{
				{4, 22},
			},
			expectedValues: [][]int{
				{},
			},
			expectedLength: 0,
		},
		{
			name: "overlap all 2",
			pairs: [][2]int{
				{5, 10},
				{15, 20},
				{25, 30},
			},
			removes: [][2]int{
				{12, 22},
			},
			expectedValues: [][]int{
				{5, 10, 25, 30},
			},
			expectedLength: 12,
		},
		{
			name: "overlap all 3",
			pairs: [][2]int{
				{5, 10},
				{15, 20},
				{25, 30},
			},
			removes: [][2]int{
				{22, 32},
			},
			expectedValues: [][]int{
				{5, 10, 15, 20},
			},
			expectedLength: 12,
		},
		{
			name: "weird cases 1",
			pairs: [][2]int{
				{5, 5},
			},
			removes: [][2]int{
				{4, 6},
			},
			expectedValues: [][]int{
				{},
			},
			expectedLength: 0,
		},
		{
			name: "weird cases 2",
			pairs: [][2]int{
				{5, 5},
			},
			removes: [][2]int{
				{5, 5},
			},
			expectedValues: [][]int{
				{},
			},
			expectedLength: 0,
		},
		{
			name: "weird cases 3",
			pairs: [][2]int{
				{0, 4},
				{5, 5},
				{6, 10},
			},
			removes: [][2]int{
				{5, 5},
			},
			expectedValues: [][]int{
				{0, 4, 6, 10},
			},
			expectedLength: 10,
		},
		{
			name: "weird cases 4",
			pairs: [][2]int{
				{10, 20},
			},
			removes: [][2]int{
				{10, 19},
			},
			expectedValues: [][]int{
				{20, 20},
			},
			expectedLength: 1,
		},
		{
			name: "weird cases 5",
			pairs: [][2]int{
				{10, 20},
			},
			removes: [][2]int{
				{21, 25},
			},
			expectedValues: [][]int{
				{10, 20},
			},
			expectedLength: 11,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rc := &Collection[int]{}

			for _, p := range tc.pairs {
				_, err := rc.Add(p[0], p[1])
				assert.NoError(t, err)
			}

			for r, remove := range tc.removes {
				after, err := rc.Remove(remove[0], remove[1])
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expectedValues[r], after)
			}

			assert.Equal(t, tc.expectedLength, rc.Len())

		})
	}
}

func Test_Add(t *testing.T) {

	testCases := []struct {
		name           string
		pairs          [][2]int
		expectedValues [][]int
		expectedLength int
	}{
		{
			name: "initial add",
			pairs: [][2]int{
				{0, 10},
			},
			expectedValues: [][]int{
				{0, 10},
			},
			expectedLength: 11,
		},
		{
			name: "length 1",
			pairs: [][2]int{
				{1, 1},
			},
			expectedValues: [][]int{
				{1, 1},
			},
			expectedLength: 1,
		},
		{
			name: "test existing overlaps new",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
			},
			expectedLength: 11,
		},
		{
			name: "test extending existing to left and right (overlaps all)",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
				{-5, 15},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
				{-5, 15},
			},
			expectedLength: 21,
		},
		{
			name: "testing extending to right",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
				{5, 15},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
				{0, 15},
			},
			expectedLength: 16,
		},
		{
			name: "test extending to left",
			pairs: [][2]int{
				{0, 10},
				{2, 8},
				{-5, 5},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10},
				{-5, 10},
			},
			expectedLength: 16,
		},
		{
			name: "test adding to right",
			pairs: [][2]int{
				{0, 10},
				{15, 20},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10, 15, 20},
			},
			expectedLength: 17,
		},
		{
			name: "test adding to left",
			pairs: [][2]int{
				{0, 10},
				{-10, -5},
			},
			expectedValues: [][]int{
				{0, 10},
				{-10, -5, 0, 10},
			},
			expectedLength: 17,
		},
		{
			name: "test joining from the middle of two",
			pairs: [][2]int{
				{0, 10},
				{15, 20},
				{5, 17},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10, 15, 20},
				{0, 20},
			},
			expectedLength: 21,
		},
		{
			name: "test length minus one",
			pairs: [][2]int{
				{0, 10},
				{12, 20},
			},
			expectedValues: [][]int{
				{0, 10},
				{0, 10, 12, 20},
			},
			expectedLength: 20,
		},
		{
			name: "test merging",
			pairs: [][2]int{
				{2, 3},
				{4, 5},
			},
			expectedValues: [][]int{
				{2, 3},
				{2, 5},
			},
			expectedLength: 4,
		},
		{
			name: "test merging 2",
			pairs: [][2]int{
				{5, 5},
				{6, 8},
			},
			expectedValues: [][]int{
				{5, 5},
				{5, 8},
			},
			expectedLength: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rc := &Collection[int]{}

			for j, p := range tc.pairs {
				values, err := rc.Add(p[0], p[1])
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedValues[j], values)
			}

			assert.Equal(t, tc.expectedLength, rc.Len())

		})
	}

}
