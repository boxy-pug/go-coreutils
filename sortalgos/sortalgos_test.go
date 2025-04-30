package sortalgos_test

import (
	"testing"

	"github.com/boxy-pug/ccsort/sortalgos"
	"github.com/google/go-cmp/cmp"
)

func TestSortingAlgorithms(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Short list of fruits",
			input:    []string{"banana", "apple", "orange"},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			name:     "Short numeric strings",
			input:    []string{"3", "1", "2"},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "Long list of mixed words",
			input:    []string{"zebra", "apple", "orange", "banana", "grape", "pear", "peach", "kiwi", "melon", "lemon"},
			expected: []string{"apple", "banana", "grape", "kiwi", "lemon", "melon", "orange", "peach", "pear", "zebra"},
		},
		{
			name:     "Long list of numbers",
			input:    []string{"10", "1", "3", "2", "5", "4", "9", "8", "7", "6"},
			expected: []string{"1", "10", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
		// Add more complex or edge case test cases as needed
	}

	for name, sortFunc := range sortalgos.SortFunctions {
		for _, tc := range testCases {
			t.Run(tc.name+"_"+string(name), func(t *testing.T) {
				list := make([]string, len(tc.input))
				copy(list, tc.input)
				sortFunc(list)
				if !cmp.Equal(list, tc.expected) {
					t.Errorf("%s() = %v, want %v", name, list, tc.expected)
				}
			})
		}
	}
}
