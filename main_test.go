package main

import (
	"log"
	"os"
	"os/exec"
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
	cmd := exec.Command("go", "run", "main.go", "./testdata/test3.txt", "./testdata/test4.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	expected := `hello
goodbye

yes man great!
checking 

this is good
`
	if string(output) != expected {
		t.Errorf("Expected %q, got %q", expected, string(output))
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
