package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var testFiles = getTestFiles("./testdata/")

func TestHeadUnit(t *testing.T) {
	testText := `at www.gutenberg.org. If you are not located in the United States,
you will have to check the laws of the country where you are located
before using this eBook.
Title: The Art of War
😊`

	var buf bytes.Buffer

	testCases := []struct {
		name    string
		reader  func() io.Reader
		command command
		want    string
	}{
		{
			name: "regular head, line flag 3 lines",
			reader: func() io.Reader {
				return strings.NewReader(testText)
			},
			command: command{
				lines:         3,
				bytes:         0,
				useLines:      true,
				useBytes:      false,
				multipleFiles: false,
				stdIn:         false,
				output:        &buf,
			},
			want: "at www.gutenberg.org. If you are not located in the United States,\nyou will have to check the laws of the country where you are located\nbefore using this eBook.\n",
		},
		{
			name: "byte flag -c 50",
			reader: func() io.Reader {
				return strings.NewReader(testText)
			},
			command: command{
				lines:         0,
				bytes:         50,
				useLines:      false,
				useBytes:      true,
				multipleFiles: false,
				stdIn:         false,
				output:        &buf,
			},
			want: "at www.gutenberg.org. If you are not located in th",
		},
		{
			name: "10 lines head, fewer lines provided",
			reader: func() io.Reader {
				return strings.NewReader(testText)
			},
			command: command{
				lines:         10,
				bytes:         0,
				useLines:      true,
				useBytes:      false,
				multipleFiles: false,
				stdIn:         false,
				output:        &buf,
			},
			want: "at www.gutenberg.org. If you are not located in the United States,\nyou will have to check the laws of the country where you are located\nbefore using this eBook.\nTitle: The Art of War\n😊",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			tc.command.files = []inputFile{
				{
					reader: tc.reader(),
				},
			}
			err := tc.command.run()
			assertNoError(t, err)

			got := buf.String()

			assertEqual(t, got, tc.want)
		})
	}
}

func TestHeadIntegration(t *testing.T) {
	t.Run("head cmd with no flags", func(t *testing.T) {
		for _, testFile := range testFiles {
			cmd := exec.Command("./cchead", testFile)
			got, err := cmd.Output()
			assertNoError(t, err)

			unixCmd := exec.Command("head", testFile)
			want, err := unixCmd.Output()
			assertNoError(t, err)

			assertEqual(t, string(got), string(want))
		}
	})

	t.Run("head cmd with bytes flag", func(t *testing.T) {
		for _, testFile := range testFiles {
			cmd := exec.Command("./cchead", "-c", "30", testFile)
			got, err := cmd.Output()
			assertNoError(t, err)

			unixCmd := exec.Command("head", "-c", "30", testFile)
			want, err := unixCmd.Output()
			assertNoError(t, err)

			assertEqual(t, string(got), string(want))
		}
	})

	t.Run("head cmd five lines", func(t *testing.T) {
		for _, testFile := range testFiles {
			cmd := exec.Command("./cchead", "-n", "5", testFile)
			got, err := cmd.Output()
			assertNoError(t, err)

			unixCmd := exec.Command("head", "-n", "5", testFile)
			want, err := unixCmd.Output()
			assertNoError(t, err)

			assertEqual(t, string(got), string(want))
		}
	})

	t.Run("head multiple files", func(t *testing.T) {
		cmd := exec.Command("./cchead", testFiles...)
		got, err := cmd.Output()
		assertNoError(t, err)

		unixCmd := exec.Command("head", testFiles...)
		want, err := unixCmd.Output()
		assertNoError(t, err)

		assertEqual(t, string(got), string(want))
	})
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("did not expect error: %v", err)
	}
}

func assertEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("GOT: %q\nWANT: %q\n", got, want)
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
