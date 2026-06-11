package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var testFiles = getTestFiles("./testdata/")

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

func TestFileInput(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testdata/test3.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	expected := `hello
goodbye

yes man great!`
	if string(output) != expected {
		t.Errorf("Expected %q, got %q", expected, string(output))
	}
}

func TestMultipleFiles(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testdata/test3.txt", "./testdata/test4.txt")
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	unixCmd := exec.Command("cat", "./testdata/test3.txt", "./testdata/test4.txt")
	want, err := unixCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("Expected %q, got %q", string(want), string(got))
	}
}

func TestNumberedLines(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-n", "./testdata/test3.txt")
	got, err := cmd.Output()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	unixCmd := exec.Command("cat", "-n", "./testdata/test3.txt")
	want, err := unixCmd.Output()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("Expected %q, got %q", string(want), string(got))
	}
}

func TestNumberedLinesRegression(t *testing.T) {
	// Regression test: old bug printed line numbers without actual line content.
	// This test would fail with fmt.Printf("%6d\t\n", lineIndex).
	cmd := exec.Command("go", "run", ".", "-n", "./testdata/test3.txt")
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	want := "     1\thello\n     2\tgoodbye\n     3\t\n     4\tyes man great!"
	if string(got) != want {
		t.Errorf("Expected %q, got %q", want, string(got))
	}
}

func TestNumberedNonBlank(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-b", "./testdata/test3.txt")
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	unixCmd := exec.Command("cat", "-b", "./testdata/test3.txt")
	want, err := unixCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("Expected %q, got %q", string(want), string(got))
	}
}

func TestNumberedNonBlankRegression(t *testing.T) {
	// Regression test: old bug numbered empty lines or omitted line content.
	// This test would fail with broken -b logic that numbers empty lines.
	cmd := exec.Command("go", "run", ".", "-b", "./testdata/test3.txt")
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	want := "     1\thello\n     2\tgoodbye\n\n     3\tyes man great!"
	if string(got) != want {
		t.Errorf("Expected %q, got %q", want, string(got))
	}
}

func TestCatCloneVsUnixCat(t *testing.T) {
	for _, testFile := range testFiles {
		cmd := exec.Command("go", "run", ".", testFile)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		unixCmd := exec.Command("cat", testFile)
		unixOutput, err := unixCmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		if string(output) != string(unixOutput) {
			t.Errorf("\tEXPECTED: %q\n\tGOT: %q\n", string(unixOutput), string(output))
		}

	}

}

// Unit tests for catInput
// Uses strings.NewReader as input and bytes.Buffer as output.

func TestCatInputPlain(t *testing.T) {
	input := strings.NewReader("hello\nworld\n")
	var buf bytes.Buffer
	cfg := config{out: &buf}

	if err := cfg.catInput(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello\nworld\n"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCatInputNumbered(t *testing.T) {
	input := strings.NewReader("a\n\nb\n")
	var buf bytes.Buffer
	cfg := config{out: &buf, numberLines: true}

	if err := cfg.catInput(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "     1\ta\n     2\t\n     3\tb\n"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCatInputNonBlank(t *testing.T) {
	input := strings.NewReader("a\n\nb\n")
	var buf bytes.Buffer
	cfg := config{out: &buf, numberNonBlank: true}

	if err := cfg.catInput(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "     1\ta\n\n     2\tb\n"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCatInputEmpty(t *testing.T) {
	input := strings.NewReader("")
	var buf bytes.Buffer
	cfg := config{out: &buf, numberLines: true}

	if err := cfg.catInput(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := ""
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCatInputNoTrailingNewline(t *testing.T) {
	input := strings.NewReader("no newline")
	var buf bytes.Buffer
	cfg := config{out: &buf}

	if err := cfg.catInput(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "no newline"
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
