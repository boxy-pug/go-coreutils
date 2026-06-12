package main

import (
	"testing"

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
		// Edge cases
		{
			name:     "Empty list",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Single element",
			input:    []string{"only"},
			expected: []string{"only"},
		},
		{
			name:     "Already sorted",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Reverse sorted",
			input:    []string{"c", "b", "a"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "All duplicates",
			input:    []string{"x", "x", "x", "x"},
			expected: []string{"x", "x", "x", "x"},
		},
		{
			name:     "Unicode strings",
			input:    []string{"ø", "å", "æ"},
			expected: []string{"å", "æ", "ø"},
		},
		{
			name:     "Mixed case",
			input:    []string{"Banana", "apple", "Orange"},
			expected: []string{"Banana", "Orange", "apple"},
		},
		{
			name:     "Two elements swapped",
			input:    []string{"b", "a"},
			expected: []string{"a", "b"},
		},
	}

	for name, sortFunc := range sortFunctions {
		for _, tc := range testCases {
			t.Run(tc.name+"_"+name, func(t *testing.T) {
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
