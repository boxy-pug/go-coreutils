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

func TestUniqCommandFile(t *testing.T) {
	for _, testFile := range testFiles {
		cmd := exec.Command("./ccuniq", testFile)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		unixCmd := exec.Command("uniq", testFile)
		unixOutput, err := unixCmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		if string(output) != string(unixOutput)+"\n" {
			t.Errorf("\tEXPECTED: %q\n\tGOT: %q\n", string(unixOutput), string(output))
		}
	}
}

func TestUniqCommandStdin(t *testing.T) {
	for _, testFile := range testFiles {
		cmd := exec.Command("cat", testFile, "|", "./ccuniq")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		unixCmd := exec.Command("cat", testFile, "|", "uniq")
		unixOutput, err := unixCmd.Output()
		if err != nil {
			t.Fatalf("Command failed with error: %v", err)
		}

		if string(output) != string(unixOutput)+"\n" {
			t.Errorf("\tEXPECTED: %q\n\tGOT: %q\n", string(unixOutput), string(output))
		}
	}
}
