package ranges

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Overlaps(t *testing.T) {

	testCases := []struct {
		input    []int
		expected []int
	}{
		{
			input:    []int{-3, 4, 0, 4, -3, 7},
			expected: []int{-3, 4},
		},
		{
			input:    []int{1, 3, 2, 4, 0, 6},
			expected: []int{1, 4},
		},
		{
			input:    []int{0, 4, 8, 9, -3, 5},
			expected: []int{0, 4},
		},
		{
			input:    []int{1, 3, 2, 4},
			expected: []int{2, 3},
		},
		{
			input:    []int{1, 3, 2, 5, 4, 6},
			expected: []int{2, 5},
		},
		{
			input:    []int{10, 20, 15, 40},
			expected: []int{15, 20},
		},
		{
			input:    []int{10, 20, 15, 25, 30, 45, 40, 50},
			expected: []int{15, 20, 40, 45},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test %d: %v", i, tc.input), func(t *testing.T) {
			results := Overlaps[int](tc.input)
			assert.NotNil(t, results)
			assert.True(t, ArrayEqual[int](results, tc.expected))
		})
	}
}
