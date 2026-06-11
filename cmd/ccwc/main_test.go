package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestWcUnit(t *testing.T) {
	t.Run("normal wordcount cmd", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := command{
			files: []fileInput{
				{
					input: strings.NewReader(`Hello testing
yes yes cool
ok goodbye now
`),
				},
			},
			out:       &buf,
			wordsFlag: true,
			linesFlag: true,
			bytesFlag: true,
		}

		cmd.run()

		got := buf.String()
		want := "       3       8      42\n"

		assertEqual(t, string(got), string(want))
	})

	t.Run("charcount from file with emojis, with filename", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := command{
			files: []fileInput{
				{
					input: strings.NewReader(`Hello testing😊
yes yes cool
ok goodbye 🌟now
`),
					name: "faketest.txt",
				},
			},

			out:              &buf,
			charsFlag:        true,
			fileNameProvided: true,
		}

		cmd.run()

		// This counts Unicode code points (runes), not bytes or grapheme clusters.
		// 1 emoji = 1 rune = 1 char
		got := buf.String()
		want := "      44 faketest.txt\n"

		assertEqual(t, string(got), string(want))
	})

	t.Run("charcount from file with emojis, no trailing newline", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := command{
			files: []fileInput{
				{
					input: strings.NewReader("Hello testing😊\nyes yes cool\nok goodbye 🌟now"), // no final '\n'
					name:  "faketest.txt",
				},
			},
			out:              &buf,
			charsFlag:        true,
			fileNameProvided: true,
		}

		cmd.run()

		got := buf.String()
		want := "      43 faketest.txt\n" // This is the expected count if the last line is counted

		assertEqual(t, got, want)
	})

	t.Run("wc counts with multiple empty lines and a regular line", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := command{
			files: []fileInput{
				{
					input: strings.NewReader("\n\n\nHello world\n"),
					name:  "faketest.txt",
				},
			},
			out:              &buf,
			linesFlag:        true,
			wordsFlag:        true,
			bytesFlag:        true,
			fileNameProvided: true,
		}

		cmd.run()

		got := buf.String()
		// 3 empty lines + "Hello world\n" = 4 lines
		// Only "Hello world" has words (2)
		// 3 empty lines = 3 bytes, "Hello world\n" = 12 bytes, total = 15 bytes
		want := "       4       2      15 faketest.txt\n"

		assertEqual(t, got, want)
	})

	t.Run("empty file", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := command{
			files: []fileInput{
				{
					input: strings.NewReader(""),
				},
			},
			out:       &buf,
			wordsFlag: true,
			linesFlag: true,
			bytesFlag: true,
		}

		cmd.run()

		got := buf.String()
		want := "       0       0       0\n"

		assertEqual(t, string(got), string(want))
	})

	t.Run("multiple files", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := command{
			files: []fileInput{
				{
					input: strings.NewReader(`Hello testing
yes yes cool
ok goodbye now
`),
					name: "file1.txt",
				},
				{
					input: strings.NewReader(`Hello testing
yes yes cool
ok goodbye now
`),
					name: "file2.txt",
				},
			},
			out:              &buf,
			wordsFlag:        true,
			linesFlag:        true,
			bytesFlag:        true,
			fileNameProvided: true,
		}

		cmd.run()

		got := buf.String()
		want := "       3       8      42 file1.txt\n" +
			"       3       8      42 file2.txt\n" +
			"       6      16      84 total\n"

		assertEqual(t, string(got), string(want))
	})
}

// Integration tests using the actual binary output.
// We compare against hardcoded strings instead of the system wc
// because different systems ship different wc implementations (gnu vs bsd)
func TestWcIntegration(t *testing.T) {
	t.Run("no flags - single file", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "./testdata/emoji.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "       7       8      60 ./testdata/emoji.txt\n"
		assertEqual(t, string(got), want)
	})

	t.Run("lines flag", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-l", "./testdata/test3.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "      11 ./testdata/test3.txt\n"
		assertEqual(t, string(got), want)
	})

	t.Run("bytes flag", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-c", "./testdata/test2.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "  163530 ./testdata/test2.txt\n"
		assertEqual(t, string(got), want)
	})

	t.Run("multiple files with total", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".",
			"./testdata/emoji.txt",
			"./testdata/test3.txt",
		)
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "       7       8      60 ./testdata/emoji.txt\n" +
			"      11       8      58 ./testdata/test3.txt\n" +
			"      18      16     118 total\n"
		assertEqual(t, string(got), want)
	})

	t.Run("stdin with no flags", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".")
		cmd.Stdin = strings.NewReader("Hello testing\nyes yes cool\nok goodbye now\n")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "       3       8      42\n"
		assertEqual(t, string(got), want)
	})

	t.Run("chars flag with emoji file", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-m", "./testdata/emoji.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "      54 ./testdata/emoji.txt\n"
		assertEqual(t, string(got), want)
	})

	t.Run("words flag alone", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-w", "./testdata/test3.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "       8 ./testdata/test3.txt\n"
		assertEqual(t, string(got), want)
	})

	t.Run("mixed flags -l and -w", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-l", "-w", "./testdata/test3.txt")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "      11       8 ./testdata/test3.txt\n"
		assertEqual(t, string(got), want)
	})

	t.Run("multiple stdin dashes", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "-", "-")
		cmd.Stdin = strings.NewReader("test\n")
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}
		want := "       1       1       5 -\n" +
			"       0       0       0 -\n" +
			"       1       1       5 total\n"
		assertEqual(t, string(got), want)
	})

	// Note: `go run` appends "exit status 1" to stderr when the program exits
	// non-zero. If running a pre-built binary instead, that line won't appear.
	t.Run("non-existent file exits with error", func(t *testing.T) {
		cmd := exec.Command("go", "run", ".", "./testdata/nonexistent.txt")
		got, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("expected command to fail, but it succeeded")
		}
		want := "loading command: open \"./testdata/nonexistent.txt\": open ./testdata/nonexistent.txt: no such file or directory\nexit status 1\n"
		assertEqual(t, string(got), want)
	})
}

func assertEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got  %q\nwant %q\n", got, want)
	}
}
