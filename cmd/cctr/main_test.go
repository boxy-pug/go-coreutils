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
			name: "range expansion first",
			cfg:  config{target: "A-Z", translation: "[:lower:]"},
			in:   "ABC",
			want: "abc",
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
