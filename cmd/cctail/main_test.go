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

func TestRegularTail(t *testing.T) {
	for _, testFile := range testFiles {
		cmd := exec.Command("go", "run", ".", testFile)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		unixCmd := exec.Command("tail", testFile)
		unixOutput, err := unixCmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		if string(output) != string(unixOutput) {
			t.Errorf("\tEXPECTED: %q\n\tGOT: %q\n", string(unixOutput), string(output))
		}

	}

}

func TestRegularTailMultiple(t *testing.T) {
	cmd := exec.Command("./cctail", testFiles...)
	got, err := cmd.Output()
	if err != nil {
		t.Fatalf("Command %s failed with error: %v, got: %s", cmd.String(), err, string(got))
	}

	unixCmd := exec.Command("tail", testFiles...)
	want, err := unixCmd.Output()
	if err != nil {
		t.Fatalf("Command failed with error: %v", err)
	}

	if string(got) != string(want) {
		t.Errorf("\tEXPECTED: %q\n\tGOT: %q\n", string(want), string(got))
	}

}
