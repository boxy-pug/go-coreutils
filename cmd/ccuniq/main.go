package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var writeToFile bool
var countCol bool
var repeatedOnly bool
var uniqueOnly bool
var inFile *os.File
var outFile *os.File

func main() {
	var err error
	flag.BoolVar(&countCol, "c", false, "enable count column")
	flag.BoolVar(&repeatedOnly, "d", false, "output only repeated")
	flag.BoolVar(&uniqueOnly, "u", false, "output only unique")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 || args[0] == string('-') {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(string(args[0]))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if len(args) > 1 && args[0] == string('-') {
		writeToFile = true
		outFile, err = os.Create(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if uniqueOnly {
		printUniqueOnly(inFile, outFile)
		return
	}

	removeDuplicateAdjecent(inFile, outFile)

}

func removeDuplicateAdjecent(inFile, outFile *os.File) {
	//uniqueLines := make(map[string]struct{})
	lineOccurences := make(map[string]int)
	scanner := bufio.NewScanner(inFile)
	var currentLine string
	var line string

	for scanner.Scan() {
		line = scanner.Text()

		if line == currentLine {
			lineOccurences[line]++
		} else {
			writeUniq(currentLine, outFile, lineOccurences)
			currentLine = line
			lineOccurences[line] = 1
		}
	}
	writeUniq(line, outFile, lineOccurences)
}

func writeUniq(str string, outFile *os.File, lineOccurences map[string]int) {
	alreadyPrinted := make(map[string]bool)
	if lineOccurences[str] == 0 {
		return
	}
	if countCol {
		fmt.Printf("%d ", lineOccurences[str])
	}
	if repeatedOnly && !alreadyPrinted[str] {
		if lineOccurences[str] > 1 {
			fmt.Println(str)
			alreadyPrinted[str] = true
			return
		} else {
			return
		}
	}
	if !writeToFile {
		fmt.Println(str)
	} else {
		_, err := outFile.WriteString(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func printUniqueOnly(inFile, outFile *os.File) {
	uniqueLines := make(map[string]struct{})
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		if _, exists := uniqueLines[line]; exists {
			delete(uniqueLines, line)
		} else {
			uniqueLines[line] = struct{}{}
		}
	}
	for key := range uniqueLines {
		fmt.Println(key)
	}
}
