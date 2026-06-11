package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestUniqUnit(t *testing.T) {
	assertEqual := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q\nwant %q\n", got, want)
		}
	}

	t.Run("default uniq", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:  strings.NewReader("a\na\nb\na\n"),
			out: &buf,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "a\nb\na\n")
	})

	t.Run("default uniq no trailing newline", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:  strings.NewReader("a\na\nb"),
			out: &buf,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "a\nb\n")
	})

	t.Run("count column", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:       strings.NewReader("a\na\nb\na\n"),
			out:      &buf,
			countCol: true,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "      2 a\n      1 b\n      1 a\n")
	})

	t.Run("repeated only", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:           strings.NewReader("a\na\nb\na\n"),
			out:          &buf,
			repeatedOnly: true,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "a\n")
	})

	t.Run("unique only", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:         strings.NewReader("a\na\nb\na\n"),
			out:        &buf,
			uniqueOnly: true,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "b\na\n")
	})

	t.Run("repeated with count", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:           strings.NewReader("a\na\nb\na\n"),
			out:          &buf,
			countCol:     true,
			repeatedOnly: true,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "      2 a\n")
	})

	t.Run("unique with count", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:         strings.NewReader("a\na\nb\na\n"),
			out:        &buf,
			countCol:   true,
			uniqueOnly: true,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "      1 b\n      1 a\n")
	})

	t.Run("all unique lines", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:         strings.NewReader("a\nb\nc\n"),
			out:        &buf,
			uniqueOnly: true,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "a\nb\nc\n")
	})

	t.Run("all repeated lines", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:           strings.NewReader("a\na\na\n"),
			out:          &buf,
			repeatedOnly: true,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "a\n")
	})

	t.Run("empty input", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:  strings.NewReader(""),
			out: &buf,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "")
	})

	t.Run("single line", func(t *testing.T) {
		var buf bytes.Buffer
		cfg := config{
			in:  strings.NewReader("a"),
			out: &buf,
		}
		if err := cfg.runUniq(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertEqual(t, buf.String(), "a\n")
	})
}

func TestUniqIntegration(t *testing.T) {
	assertEqual := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q\nwant %q\n", got, want)
		}
	}

	t.Run("default with testdata", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "./testdata/test.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "line1\nline2\nline3\nline4\nline5\nline2\n"
		assertEqual(t, string(got), want)
	})

	t.Run("unique only with testdata", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-u", "./testdata/test.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "line1\nline3\nline5\nline2\n"
		assertEqual(t, string(got), want)
	})

	t.Run("repeated only with testdata", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-d", "./testdata/test.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "line2\nline4\n"
		assertEqual(t, string(got), want)
	})

	t.Run("count with testdata", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-c", "./testdata/test.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "      1 line1\n      2 line2\n      1 line3\n      4 line4\n      1 line5\n      1 line2\n"
		assertEqual(t, string(got), want)
	})

	t.Run("repeated only with countries", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-d", "./testdata/countries.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		// Brazil appears twice (lines 31,32), Italy appears twice (107,108), Turkey appears twice (228,229)
		want := "Brazil\nItaly\nTurkey\n"
		assertEqual(t, string(got), want)
	})
}
