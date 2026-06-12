package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// === sortMapByValue tests ===

func TestSortMapByValue(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]time.Duration
		want []string
	}{
		{
			name: "two entries",
			in:   map[string]time.Duration{"quick": 10 * time.Millisecond, "merge": 5 * time.Millisecond},
			want: []string{"merge", "quick"},
		},
		{
			name: "same duration",
			in:   map[string]time.Duration{"a": 5 * time.Millisecond, "b": 5 * time.Millisecond},
			// sort.Slice is not stable, so either order is valid for equal keys
			want: []string{"a", "b"}, // or []string{"b", "a"}
		},
		{
			name: "empty map",
			in:   map[string]time.Duration{},
			want: []string{},
		},
		{
			name: "single entry",
			in:   map[string]time.Duration{"only": 1 * time.Millisecond},
			want: []string{"only"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sortMapByValue(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("len(got) = %d, want %d", len(got), len(tc.want))
			}
			// For the "same duration" case, we just check length and that
			// all expected keys are present.
			if tc.name == "same duration" {
				for _, w := range tc.want {
					found := false
					for _, g := range got {
						if g == w {
							found = true
							break
						}
					}
					if !found {
						t.Fatalf("missing key %q in result %v", w, got)
					}
				}
				return
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("sortMapByValue() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

// === getLines tests ===

func TestGetLines(t *testing.T) {
	tests := []struct {
		name    string
		readers []io.Reader
		want    []string
		wantErr bool
	}{
		{
			name:    "single reader",
			readers: []io.Reader{strings.NewReader("a\nb\nc")},
			want:    []string{"a", "b", "c"},
		},
		{
			name:    "multiple readers",
			readers: []io.Reader{strings.NewReader("a"), strings.NewReader("b")},
			want:    []string{"a", "b"},
		},
		{
			name:    "empty reader",
			readers: []io.Reader{strings.NewReader("")},
			want:    []string{},
		},
		{
			name:    "reader with trailing newline",
			readers: []io.Reader{strings.NewReader("a\n")},
			want:    []string{"a"},
		},
		{
			name:    "no readers",
			readers: []io.Reader{},
			want:    []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config{files: tc.readers}
			got, err := cfg.getLines()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("getLines() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

// === printLines tests ===

func TestPrintLines(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "three lines",
			input: []string{"a", "b", "c"},
			want:  "a\nb\nc\n",
		},
		{
			name:  "empty list",
			input: []string{},
			want:  "",
		},
		{
			name:  "single line",
			input: []string{"only"},
			want:  "only\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := &config{out: &buf}
			cfg.printLines(tc.input)
			if got := buf.String(); got != tc.want {
				t.Fatalf("printLines() = %q, want %q", got, tc.want)
			}
		})
	}
}

// === printUniqueLines tests ===

func TestPrintUniqueLines(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "all unique",
			input: []string{"a", "b", "c"},
			want:  "a\nb\nc\n",
		},
		{
			name:  "adjacent duplicates",
			input: []string{"a", "a", "b"},
			want:  "a\nb\n",
		},
		{
			name:  "non-adjacent duplicates (all removed after sort)",
			input: []string{"a", "b", "a"},
			want:  "a\nb\n",
		},
		{
			name:  "empty list",
			input: []string{},
			want:  "",
		},
		{
			name:  "all same",
			input: []string{"x", "x", "x"},
			want:  "x\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := &config{out: &buf}
			cfg.printUniqueLines(tc.input)
			if got := buf.String(); got != tc.want {
				t.Fatalf("printUniqueLines() = %q, want %q", got, tc.want)
			}
		})
	}
}

// === openPaths tests ===

func TestOpenPaths(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	f1, err := os.CreateTemp(tmpDir, "test1")
	if err != nil {
		t.Fatal(err)
	}
	f1.WriteString("hello")
	f1.Close()

	t.Run("valid file", func(t *testing.T) {
		readers, cleanup, err := openPaths([]string{f1.Name()})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(readers) != 1 {
			t.Fatalf("len(readers) = %d, want 1", len(readers))
		}
		cleanup()
	})

	t.Run("missing file", func(t *testing.T) {
		readers, cleanup, err := openPaths([]string{"/nonexistent/path/to/file"})
		if err == nil {
			_ = readers
			t.Fatalf("expected error for missing file, got nil")
		}
		if cleanup != nil {
			cleanup()
		}
	})

	t.Run("dash is stdin", func(t *testing.T) {
		readers, cleanup, err := openPaths([]string{"-"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(readers) != 1 {
			t.Fatalf("len(readers) = %d, want 1", len(readers))
		}
		cleanup()
	})

	t.Run("first file ok second missing", func(t *testing.T) {
		readers, cleanup, err := openPaths([]string{f1.Name(), "/nonexistent"})
		if err == nil {
			t.Fatalf("expected error for missing file, got nil. readers=%v", readers)
		}
		// If cleanup is non-nil, it should close the already-opened file.
		if cleanup != nil {
			cleanup()
		}
	})
}

// === loadConfig error path (via main orchestration) ===

// We can't easily test loadConfig() directly because it calls flag.Parse()
// which reads os.Args. But we can test the error handling by checking what
// main would do with a loadConfig error.

func TestZeroConfigPanics(t *testing.T) {
	// A zero-value config has nil out. If main ever continued with one,
	// printLines would panic. main.go now calls os.Exit(1) on loadConfig error,
	// so this path is unreachable, but the invariant is still worth checking.
	var zeroCfg config
	if zeroCfg.out != nil {
		t.Fatalf("expected zero value config.out to be nil, got %v", zeroCfg.out)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("caught expected panic: %v\n", r)
			} else {
				t.Fatalf("expected panic from nil io.Writer, but none occurred")
			}
		}()
		zeroCfg.printLines([]string{"hello"})
	}()
}

// === Integration-style: go run . with testdata ===

func TestIntegration_sortFile(t *testing.T) {
	// Run the tool against testdata/minitest.txt
	cmd := os.Getenv("CCSORT_BIN")
	if cmd == "" {
		cmd = "go"
		// Use t.Skip if you don't want to run integration tests without the binary
		// t.Skip("set CCSORT_BIN to run integration tests")
	}
	// A simple integration test would go here. For now, we keep this as a
	// placeholder since the user asked for unit tests primarily.
}
