package main

import (
	"bytes"
	"strings"
	"testing"
)

// --- parseFlags ---

func TestParseFlags(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    config
		wantErr bool
	}{
		{
			name: "target and translation",
			args: []string{"abc", "xyz"},
			want: config{target: "abc", translation: "xyz"},
		},
		{
			name: "delete flag",
			args: []string{"-d", "abc"},
			want: config{deleteFlag: true, target: "abc"},
		},
		{
			name: "squeeze flag",
			args: []string{"-s", "abc"},
			want: config{squeezeFlag: true, target: "abc"},
		},
		{
			name:    "too few args",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name:    "delete and squeeze",
			args:    []string{"-d", "-s", "abc"},
			wantErr: true,
		},
		{
			name:    "too many args",
			args:    []string{"abc", "xyz", "q"},
			wantErr: true,
		},
		{
			name:    "delete with two args",
			args:    []string{"-d", "abc", "xyz"},
			wantErr: true,
		},
		{
			name:    "squeeze with two args",
			args:    []string{"-s", "abc", "xyz"},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseFlags(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseFlags(%v): expected error, got nil", tc.args)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseFlags(%v): unexpected error: %v", tc.args, err)
			}
			if got.target != tc.want.target || got.translation != tc.want.translation || got.deleteFlag != tc.want.deleteFlag || got.squeezeFlag != tc.want.squeezeFlag {
				t.Fatalf("parseFlags(%v) = %+v, want %+v", tc.args, got, tc.want)
			}
		})
	}
}

// --- buildTranslator ---

func TestBuildTranslator(t *testing.T) {
	cases := []struct {
		name    string
		cfg     config
		in      string
		want    string
		wantErr bool
	}{
		{
			name: "literal translation",
			cfg:  config{target: "abc", translation: "xyz"},
			in:   "abc",
			want: "xyz",
		},
		{
			name: "delete",
			cfg:  config{target: "abc", deleteFlag: true},
			in:   "abc123",
			want: "123",
		},
		{
			name: "range expansion",
			cfg:  config{target: "a-c", translation: "XYZ"},
			in:   "abc",
			want: "XYZ",
		},
		{
			name: "class expansion",
			cfg:  config{target: "[:lower:]", translation: "[:upper:]"},
			in:   "abc",
			want: "ABC",
		},
		{
			name: "last char repeats",
			cfg:  config{target: "abc", translation: "xy"},
			in:   "abc",
			want: "xyy",
		},
		{
			name: "unchanged chars",
			cfg:  config{target: "abc", translation: "ABC"},
			in:   "abc123",
			want: "ABC123",
		},
		{
			name: "class expansion digit",
			cfg:  config{target: "[:digit:]", translation: "?"},
			in:   "a1b2c3",
			want: "a?b?c?",
		},
		{
			name:    "range expansion first",
			cfg:     config{target: "A-Z", translation: "[:lower:]"},
			wantErr: true,
		},
		// Misalignment errors: these should fail until we add validation
		{
			name:    "alpha to upper: invalid class in string2",
			cfg:     config{target: "[:alpha:]", translation: "[:upper:]"},
			wantErr: true,
		},
		{
			name:    "digit to lower: invalid class in string2",
			cfg:     config{target: "[:digit:]", translation: "[:lower:]"},
			wantErr: true,
		},
		{
			name:    "lower to digit: invalid class in string2",
			cfg:     config{target: "[:lower:]", translation: "[:digit:]"},
			wantErr: true,
		},
		{
			name:    "literal to upper: misaligned",
			cfg:     config{target: "b", translation: "[:upper:]"},
			wantErr: true,
		},
		{
			name:    "upper to lower class: valid case conversion",
			cfg:     config{target: "[:upper:]", translation: "[:lower:]"},
			in:      "ABC",
			want:    "abc",
		},
		{
			name:    "lower to lower class: valid no-op",
			cfg:     config{target: "[:lower:]", translation: "[:lower:]"},
			in:      "abc",
			want:    "abc",
		},
		{
			name: "empty target delete",
			cfg:  config{target: "", deleteFlag: true},
			in:   "abc",
			want: "abc",
		},
		{
			name: "empty target squeeze",
			cfg:  config{target: "", squeezeFlag: true},
			in:   "abc",
			want: "abc",
		},
		{
			name: "multi char range",
			cfg:  config{target: "ab-c", translation: "XYZ"},
			in:   "abc",
			want: "XYZ",
		},
		{
			name: "range dash at end",
			cfg:  config{target: "abc-", translation: "xyz-"},
			in:   "abc-",
			want: "xyz-",
		},
		{
			name:    "non-empty target with empty translation",
			cfg:     config{target: "abc", translation: ""},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr, err := buildTranslator(tc.cfg)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("buildTranslator(%+v): expected error, got nil", tc.cfg)
				}
				return
			}
			if err != nil {
				t.Fatalf("buildTranslator(%+v): unexpected error: %v", tc.cfg, err)
			}

			var out bytes.Buffer
			cfg := config{in: strings.NewReader(tc.in), out: &out}
			if err := cfg.translate(tr); err != nil {
				t.Fatalf("translate: %v", err)
			}
			if got := out.String(); got != tc.want {
				t.Fatalf("translate(%q): got %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// --- run ---

func TestRun(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		in      string
		want    string
		wantErr bool
	}{
		{
			name: "basic translation",
			args: []string{"cmd", "abc", "xyz"},
			in:   "abc",
			want: "xyz",
		},
		{
			name: "delete",
			args: []string{"cmd", "-d", "abc"},
			in:   "abc123",
			want: "123",
		},
		{
			name: "range expansion",
			args: []string{"cmd", "a-c", "XYZ"},
			in:   "abc",
			want: "XYZ",
		},
		{
			name: "class expansion",
			args: []string{"cmd", "[:lower:]", "[:upper:]"},
			in:   "abc",
			want: "ABC",
		},
		{
			name: "unchanged chars",
			args: []string{"cmd", "abc", "ABC"},
			in:   "abc123",
			want: "ABC123",
		},
		{
			name: "digit class",
			args: []string{"cmd", "[:digit:]", "X"},
			in:   "a1b2c3",
			want: "aXbXcX",
		},
		{
			name: "alpha class",
			args: []string{"cmd", "[:alpha:]", "?"},
			in:   "a1b2c3",
			want: "?1?2?3",
		},
		{
			name: "punct class",
			args: []string{"cmd", "[:punct:]", " "},
			in:   "a,b!c?d",
			want: "a b c d",
		},
		{
			name: "space class",
			args: []string{"cmd", "[:space:]", "_"},
			in:   "a b\tc\nd",
			want: "a_b_c_d",
		},
		{
			name: "upper to lower",
			args: []string{"cmd", "[:upper:]", "[:lower:]"},
			in:   "ABCD",
			want: "abcd",
		},
		{
			name:    "missing translation",
			args:    []string{"cmd", "abc"},
			wantErr: true,
		},
		// Misalignment errors: these should fail until we add validation
		{
			name:    "alpha to upper class: invalid",
			args:    []string{"cmd", "[:alpha:]", "[:upper:]"},
			wantErr: true,
		},
		{
			name:    "digit to lower class: invalid",
			args:    []string{"cmd", "[:digit:]", "[:lower:]"},
			wantErr: true,
		},
		{
			name:    "lower to digit class: invalid",
			args:    []string{"cmd", "[:lower:]", "[:digit:]"},
			wantErr: true,
		},
		{
			name:    "literal to upper class: misaligned",
			args:    []string{"cmd", "b", "[:upper:]"},
			wantErr: true,
		},
		// Squeeze tests: will fail until squeeze is implemented
		{
			name: "squeeze one char",
			args: []string{"cmd", "-s", "o"},
			in:   "helloooo",
			want: "hello",
		},
		{
			name: "squeeze set",
			args: []string{"cmd", "-s", "abc"},
			in:   "aaabbbccc",
			want: "abc",
		},
		{
			name: "squeeze mixed chars",
			args: []string{"cmd", "-s", "lo"},
			in:   "hello",
			want: "helo",
		},
		{
			name: "squeeze leaves non-target chars alone",
			args: []string{"cmd", "-s", "ab"},
			in:   "aaXXbb",
			want: "aXXb",
		},
		// Edge cases
		// BUG: empty translation causes panic in buildTranslator (translation[len(translation)-1])
		{
			name:    "non-empty target with empty translation",
			args:    []string{"cmd", "abc", ""},
			in:      "abc",
			wantErr: true,
		},
		{
			name: "empty input",
			args: []string{"cmd", "abc", "xyz"},
			in:   "",
			want: "",
		},
		{
			name: "delete empty target",
			args: []string{"cmd", "-d", ""},
			in:   "abc",
			want: "abc",
		},
		{
			name: "squeeze empty target",
			args: []string{"cmd", "-s", ""},
			in:   "hello",
			want: "hello",
		},
		{
			name:    "delete with two args",
			args:    []string{"cmd", "-d", "abc", "xyz"},
			in:      "abc123",
			wantErr: true,
		},
		{
			name:    "squeeze with two args",
			args:    []string{"cmd", "-s", "abc", "xyz"},
			in:      "aaabbbccc",
			wantErr: true,
		},
		// Delete with classes
		{
			name: "delete digit class",
			args: []string{"cmd", "-d", "[:digit:]"},
			in:   "a1b2c3",
			want: "abc",
		},
		{
			name: "delete alpha class",
			args: []string{"cmd", "-d", "[:alpha:]"},
			in:   "a1b2c3",
			want: "123",
		},
		{
			name: "delete punct class",
			args: []string{"cmd", "-d", "[:punct:]"},
			in:   "a,b!c?d",
			want: "abcd",
		},
		{
			name: "delete space class",
			args: []string{"cmd", "-d", "[:space:]"},
			in:   "a b\tc\nd",
			want: "abcd",
		},
		{
			name: "delete upper class",
			args: []string{"cmd", "-d", "[:upper:]"},
			in:   "ABCabc",
			want: "abc",
		},
		{
			name: "delete lower class",
			args: []string{"cmd", "-d", "[:lower:]"},
			in:   "ABCabc",
			want: "ABC",
		},
		{
			name: "delete print class",
			args: []string{"cmd", "-d", "[:print:]"},
			in:   "abc",
			want: "",
		},
		// Squeeze with classes
		{
			name: "squeeze digit class",
			args: []string{"cmd", "-s", "[:digit:]"},
			in:   "a111b222c",
			want: "a1b2c",
		},
		{
			name: "squeeze alpha class",
			args: []string{"cmd", "-s", "[:alpha:]"},
			in:   "aaabbbccc",
			want: "abc",
		},
		{
			name: "squeeze lower class",
			args: []string{"cmd", "-s", "[:lower:]"},
			in:   "aaAAbb",
			want: "aAAb",
		},
		{
			name: "squeeze upper class",
			args: []string{"cmd", "-s", "[:upper:]"},
			in:   "AAbbCC",
			want: "AbbC",
		},
		{
			name: "squeeze space class spaces",
			args: []string{"cmd", "-s", "[:space:]"},
			in:   "a  b",
			want: "a b",
		},
		{
			name: "squeeze space class tabs",
			args: []string{"cmd", "-s", "[:space:]"},
			in:   "a\t\tb",
			want: "a\tb",
		},
		{
			name: "squeeze space class newlines",
			args: []string{"cmd", "-s", "[:space:]"},
			in:   "a\n\nb",
			want: "a\nb",
		},
		{
			name: "squeeze space class mixed no squeeze",
			args: []string{"cmd", "-s", "[:space:]"},
			in:   "a \t\nb",
			want: "a \t\nb",
		},
		{
			name: "squeeze punct class",
			args: []string{"cmd", "-s", "[:punct:]"},
			in:   "a!!!b",
			want: "a!b",
		},
		{
			name: "squeeze print class",
			args: []string{"cmd", "-s", "[:print:]"},
			in:   "aa",
			want: "a",
		},
		// Range edge cases
		{
			name: "range upper to lower",
			args: []string{"cmd", "A-Z", "a-z"},
			in:   "ABC",
			want: "abc",
		},
		{
			name: "range lower to upper",
			args: []string{"cmd", "a-z", "A-Z"},
			in:   "abc",
			want: "ABC",
		},
		{
			name: "range upper to lower on numbers",
			args: []string{"cmd", "A-Z", "a-z"},
			in:   "123",
			want: "123",
		},
		{
			name: "range lower to upper on numbers",
			args: []string{"cmd", "a-z", "A-Z"},
			in:   "123",
			want: "123",
		},
		// Class translation edge cases
		{
			name: "lower to upper class",
			args: []string{"cmd", "[:lower:]", "[:upper:]"},
			in:   "abc",
			want: "ABC",
		},
		{
			name: "upper to lower class on mixed",
			args: []string{"cmd", "[:upper:]", "[:lower:]"},
			in:   "AbC",
			want: "abc",
		},
		{
			name: "lower to upper class on mixed",
			args: []string{"cmd", "[:lower:]", "[:upper:]"},
			in:   "AbC",
			want: "ABC",
		},
		{
			name: "upper to lower class on empty",
			args: []string{"cmd", "[:upper:]", "[:lower:]"},
			in:   "",
			want: "",
		},
		{
			name: "lower to upper class on empty",
			args: []string{"cmd", "[:lower:]", "[:upper:]"},
			in:   "",
			want: "",
		},
		{
			name: "upper to lower class on numbers",
			args: []string{"cmd", "[:upper:]", "[:lower:]"},
			in:   "123",
			want: "123",
		},
		{
			name: "lower to upper class on numbers",
			args: []string{"cmd", "[:lower:]", "[:upper:]"},
			in:   "123",
			want: "123",
		},
		// No-match cases
		{
			name: "delete no match",
			args: []string{"cmd", "-d", "xyz"},
			in:   "abc",
			want: "abc",
		},
		{
			name: "squeeze no match",
			args: []string{"cmd", "-s", "xyz"},
			in:   "abc",
			want: "abc",
		},
		{
			name: "translate no match",
			args: []string{"cmd", "xyz", "abc"},
			in:   "123",
			want: "123",
		},
		// Newline and tab handling
		{
			name: "delete newline",
			args: []string{"cmd", "-d", "\n"},
			in:   "a\nb\nc",
			want: "abc",
		},
		{
			name: "squeeze newline",
			args: []string{"cmd", "-s", "\n"},
			in:   "a\n\n\nb",
			want: "a\nb",
		},
		{
			name: "delete tab",
			args: []string{"cmd", "-d", "\t"},
			in:   "a\tb\tc",
			want: "abc",
		},
		{
			name: "squeeze tab",
			args: []string{"cmd", "-s", "\t"},
			in:   "a\t\t\tb",
			want: "a\tb",
		},
		// Single char
		{
			name: "single char translate",
			args: []string{"cmd", "a", "b"},
			in:   "a",
			want: "b",
		},
		{
			name: "single char delete",
			args: []string{"cmd", "-d", "a"},
			in:   "a",
			want: "",
		},
		{
			name: "single char squeeze",
			args: []string{"cmd", "-s", "a"},
			in:   "a",
			want: "a",
		},
		{
			name: "single char squeeze multiple",
			args: []string{"cmd", "-s", "a"},
			in:   "aaa",
			want: "a",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(tc.args, strings.NewReader(tc.in), &out)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("run(%v, %q): expected error, got nil", tc.args, tc.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("run(%v, %q): unexpected error: %v", tc.args, tc.in, err)
			}
			if got := out.String(); got != tc.want {
				t.Fatalf("run(%v, %q): got %q, want %q", tc.args, tc.in, got, tc.want)
			}
		})
	}
}
