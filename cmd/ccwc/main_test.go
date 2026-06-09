package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestWcUnit(t *testing.T) {
	t.Run("normal wordcount cmd", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := Command{
			Files: []FileInput{
				{
					Input: strings.NewReader(`Hello testing
yes yes cool
ok goodbye now
`),
				},
			},
			Output:    &buf,
			WordsFlag: true,
			LinesFlag: true,
			BytesFlag: true,
		}

		cmd.Run()

		got := buf.String()
		want := "       3       8      42\n"

		assertEqual(t, string(got), string(want))
	})

	t.Run("charcount from file with emojis, with filename", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := Command{
			Files: []FileInput{
				{
					Input: strings.NewReader(`Hello testing😊
yes yes cool
ok goodbye 🌟now
`),
					FileName: "faketest.txt",
				},
			},

			Output:           &buf,
			CharsFlag:        true,
			FileNameProvided: true,
		}

		cmd.Run()

		// This counts Unicode code points (runes), not bytes or grapheme clusters.
		// 1 emoji = 1 rune = 1 char
		got := buf.String()
		want := "      44 faketest.txt\n"

		assertEqual(t, string(got), string(want))
	})

	t.Run("charcount from file with emojis, no trailing newline", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := Command{
			Files: []FileInput{
				{
					Input:    strings.NewReader("Hello testing😊\nyes yes cool\nok goodbye 🌟now"), // no final '\n'
					FileName: "faketest.txt",
				},
			},
			Output:           &buf,
			CharsFlag:        true,
			FileNameProvided: true,
		}

		cmd.Run()

		got := buf.String()
		want := "      43 faketest.txt\n" // This is the expected count if the last line is counted

		assertEqual(t, got, want)
	})

	t.Run("wc counts with multiple empty lines and a regular line", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := Command{
			Files: []FileInput{
				{
					Input:    strings.NewReader("\n\n\nHello world\n"),
					FileName: "faketest.txt",
				},
			},
			Output:           &buf,
			LinesFlag:        true,
			WordsFlag:        true,
			BytesFlag:        true,
			FileNameProvided: true,
		}

		cmd.Run()

		got := buf.String()
		// 3 empty lines + "Hello world\n" = 4 lines
		// Only "Hello world" has words (2)
		// 3 empty lines = 3 bytes, "Hello world\n" = 12 bytes, total = 15 bytes
		want := "       4       2      15 faketest.txt\n"

		assertEqual(t, got, want)
	})

	t.Run("empty file", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := Command{
			Files: []FileInput{
				{
					Input: strings.NewReader(""),
				},
			},
			Output:    &buf,
			WordsFlag: true,
			LinesFlag: true,
			BytesFlag: true,
		}

		cmd.Run()

		got := buf.String()
		want := "       0       0       0\n"

		assertEqual(t, string(got), string(want))
	})

	t.Run("multiple files", func(t *testing.T) {
		var buf bytes.Buffer
		cmd := Command{
			Files: []FileInput{
				{
					Input: strings.NewReader(`Hello testing
yes yes cool
ok goodbye now
`),
					FileName: "file1.txt",
				},
				{
					Input: strings.NewReader(`Hello testing
yes yes cool
ok goodbye now
`),
					FileName: "file2.txt",
				},
			},
			Output:           &buf,
			WordsFlag:        true,
			LinesFlag:        true,
			BytesFlag:        true,
			FileNameProvided: true,
		}

		cmd.Run()

		got := buf.String()
		want := "       3       8      42 file1.txt\n" +
			"       3       8      42 file2.txt\n" +
			"       6      16      84 total\n"

		assertEqual(t, string(got), string(want))
	})
}

// Old integration tests
func TestWcIntegration(t *testing.T) {
	testFiles := getTestFiles("./testdata/")

	t.Run("Test Wc without flags", func(t *testing.T) {
		for _, testFile := range testFiles {
			cmd := exec.Command("./ccwc", testFile)
			got, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command %s failed with error: %v", cmd.String(), err)
			}

			unixCmd := exec.Command("wc", testFile)
			want, err := unixCmd.Output()
			if err != nil {
				t.Fatalf("Command %s failed with error: %v", unixCmd.String(), err)
			}

			assertEqual(t, string(got), string(want))
		}
	})

	t.Run("Test wc with lines flag", func(t *testing.T) {
		for _, testFile := range testFiles {
			cmd := exec.Command("./ccwc", "-l", testFile)
			got, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command %s failed with error: %v", cmd.String(), err)
			}

			unixCmd := exec.Command("wc", "-l", testFile)
			want, err := unixCmd.Output()
			if err != nil {
				t.Fatalf("Command %s failed with error: %v", unixCmd.String(), err)
			}

			assertEqual(t, string(got), string(want))
		}
	})

	t.Run("Test with bytes flag", func(t *testing.T) {
		for _, testFile := range testFiles {
			cmd := exec.Command("./ccwc", "-c", testFile)
			got, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command %s failed with error: %v", cmd.String(), err)
			}

			unixCmd := exec.Command("wc", "-c", testFile)
			want, err := unixCmd.Output()
			if err != nil {
				t.Fatalf("Command %s failed with error: %v", unixCmd.String(), err)
			}

			assertEqual(t, string(got), string(want))
		}
	})

	t.Run("Test wc with multiple files", func(t *testing.T) {
		cmd := exec.Command("./ccwc", testFiles...)
		got, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command %s failed with error: %v", cmd.String(), err)
		}

		unixCmd := exec.Command("wc", testFiles...)
		want, err := unixCmd.Output()
		if err != nil {
			t.Fatalf("Command %s failed with error: %v", unixCmd.String(), err)
		}

		assertEqual(t, string(got), string(want))
	})
}

func assertEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got  %q/nwant %q/n", got, want)
	}
}

func getTestFiles(testFolder string) []string {
	var res []string

	files, err := os.ReadDir(testFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		res = append(res, testFolder+file.Name())
	}
	return res
}
