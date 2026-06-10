package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	numberedLines          = flag.Bool("n", false, "number lines")
	numberedLinesJumpEmpty = flag.Bool("b", false, "number lines jump empty")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 || args[0] == "-" {
		catInput(os.Stdin)
	} else {
		for _, arg := range args {
			file, err := os.Open(arg)
			if err != nil {
				fmt.Printf("error opening %s: %v\n", arg, err)
			}
			catInput(file)
			file.Close()
		}
	}
}

func catInput(file *os.File) {
	scanner := bufio.NewScanner(file)
	lineIndex := 1

	for scanner.Scan() {
		line := scanner.Text()

		if *numberedLinesJumpEmpty && line == "" {
			fmt.Println("")
			continue
		}

		if *numberedLines || *numberedLinesJumpEmpty {
			fmt.Printf("%6d\t%s\n", lineIndex, line)
		} else {
			fmt.Println(line)
		}

		lineIndex += 1

		if err := scanner.Err(); err != nil {
			fmt.Printf("error reading input: %v\n", err)
			os.Exit(1)
		}
	}
}
