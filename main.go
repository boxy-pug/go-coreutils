package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	numberedLines          = flag.Bool("n", false, "number lines")
	numberedLinesJumpEmpty = flag.Bool("b", false, "number lines jump empty")
)

func main() {
	flag.Parse()

	if len(os.Args) == 1 || os.Args[1] == "-" {
		catInput(os.Stdin)
	} else {
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(strings.TrimSpace(arg), "-") {
				continue
			}
			// fmt.Printf("Parsing file: %s\n", arg)
			file, err := os.Open(arg)
			if err != nil {
				fmt.Printf("error opening %s: %v\n", arg, err)
			}
			defer file.Close()
			catInput(file)
		}
	}
}

func catInput(file *os.File) {
	scanner := bufio.NewScanner(file)
	lineIndex := 0
	// lastLine := ""
	// lastLineHadNewline := false

	for scanner.Scan() {
		line := scanner.Text()
		// lastLine = line
		// lastLineHadNewline = strings.HasSuffix(lastLine, "\n")

		if *numberedLines || *numberedLinesJumpEmpty {
			if line == "" && *numberedLinesJumpEmpty {
				fmt.Printf("%s\n", line)
				continue
			}
			lineIndex += 1
			fmt.Printf("     %d\t\n", lineIndex)
		} else {
			fmt.Println(line)
		}

	}
}
