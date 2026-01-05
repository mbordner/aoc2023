package geom3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLine_IntersectionPoint(t *testing.T) {
	testCases := []struct {
		l1 []int
		l2 []int
		p1 *Point[float64]
	}{
		{
			l1: []int{15, 2, -1, 4, 1, -1},
			l2: []int{13, -5, -4, -5, 2, 3},
			p1: &Point[float64]{X: 3, Y: -1, Z: 2},
		},
		{
			l1: []int{19, 13, 0, -2, 1, 0},
			l2: []int{18, 19, 0, -1, -1, 0},
			p1: &Point[float64]{X: 14.333333333333333, Y: 15.333333333333333, Z: 0},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			l1 := NewLineFromVals(tc.l1)
			l2 := NewLineFromVals(tc.l2)
			p := l1.IntersectionPoint(l2)
			if tc.p1 == nil {
				assert.Nil(t, p)
			} else {
				assert.True(t, tc.p1.Equal(p))
			}
		})
	}
}

func TestLine_Contains(t *testing.T) {
	testCases := []struct {
		l1       []int
		p1       []int
		expected bool
	}{
		{
			l1:       []int{1, 1, 1, 1, 1, 1},
			p1:       []int{1, 1, 1},
			expected: true,
		},
		{
			l1:       []int{1, 1, 1, 1, 1, 1},
			p1:       []int{5, 9, -4},
			expected: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			p1 := Point[int]{X: tc.p1[0], Y: tc.p1[1], Z: tc.p1[2]}
			l1 := NewLineFromVals(tc.l1)
			assert.Equal(t, tc.expected, l1.Contains(&p1))
		})
	}
}

func TestLine_Coincident(t *testing.T) {
	testCases := []struct {
		l1       []int
		l2       []int
		expected bool
	}{
		{
			l1:       []int{1, 1, 1, 1, 1, 1},
			l2:       []int{2, 2, 2, 1, 1, 1},
			expected: true,
		},
		{
			l1:       []int{1, 1, 1, 1, 1, 1},
			l2:       []int{2, 2, 0, 3, 5, 7},
			expected: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			l1 := NewLineFromVals(tc.l1)
			l2 := NewLineFromVals(tc.l2)
			assert.Equal(t, tc.expected, l1.Coincident(l2))
		})
	}
}

func TestLine_ParallelTo(t *testing.T) {

	testCases := []struct {
		l1       []int
		l2       []int
		expected bool
	}{
		{
			l1:       []int{1, 1, 1, 1, 1, 1},
			l2:       []int{2, 2, 0, 2, 2, 2},
			expected: true,
		},
		{
			l1:       []int{1, 1, 1, 1, 1, 1},
			l2:       []int{2, 2, 0, 3, 5, 7},
			expected: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			l1 := NewLineFromVals(tc.l1)
			l2 := NewLineFromVals(tc.l2)
			assert.Equal(t, tc.expected, l1.Parallel(l2))
		})
	}
}
