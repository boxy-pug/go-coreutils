package main

import (
	"testing"
)

/*
import "testing"

func TestToDigit(t *testing.T) {
	t.Run("letter to digit conversion", func(t *testing.T) {
		got := ToDigit('h')
		want := '7'

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
*/

func TestToLetter(t *testing.T) {
	tests := []struct {
		input    rune
		expected rune
	}{
		{'a', 'a'}, // already a letter
		{'A', 'A'}, // already a letter
		{'1', 'x'}, // maps to 'b' (since '1' % 26 = 1)
		{'!', 'h'}, // maps to 'j' (since '!' % 26 = 9)
		{' ', 'g'}, // maps to 'a' (since ' ' % 26 = 0)
	}

	for _, test := range tests {
		got := ToLetter(test.input)
		want := test.expected

		assertEqualRunes(t, got, want)
	}
}

func assertEqualRunes(t testing.TB, got, want rune) {
	t.Helper()

	if got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}
