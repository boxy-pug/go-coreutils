package main

//
// import (
// 	"testing"
// )
//
// // ToLetter tests are kept here as they test a pure function
// // that is independent of the translator refactor.
//
// func TestToLetter(t *testing.T) {
// 	tests := []struct {
// 		input    rune
// 		expected rune
// 	}{
// 		{'a', 'a'}, // already a letter
// 		{'A', 'A'}, // already a letter
// 		{'1', 'b'}, // maps to 'b' (since '1' % 26 = 1)
// 		{'!', 'j'}, // maps to 'j' (since '!' % 26 = 9)
// 		{' ', 'a'}, // maps to 'a' (since ' ' % 26 = 0)
// 	}
//
// 	for _, test := range tests {
// 		got := ToLetter(test.input)
// 		want := test.expected
//
// 		assertEqualRunes(t, got, want)
// 	}
// }
//
// func assertEqualRunes(t testing.TB, got, want rune) {
// 	t.Helper()
//
// 	if got != want {
// 		t.Fatalf("got %v want %v", got, want)
// 	}
// }
