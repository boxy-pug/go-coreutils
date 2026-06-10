package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

// catInput cats the file or stdin input to stdout
func catInput(file *os.File) {

	// Regular cat is just io.Copy
	if !*numberedLines && !*numberedLinesJumpEmpty {
		_, err := io.Copy(os.Stdout, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %s", err)
			os.Exit(1)
		}
		return
	}

	// Make a reader to read line by line
	reader := bufio.NewReader(file)
	lineIndex := 1

	for {
		line, err := reader.ReadString('\n') // This preserves the newline in line

		if *numberedLinesJumpEmpty && (line == "" || line == "\n") {
			fmt.Printf("%s", line)
		} else {
			fmt.Printf("%6d\t%s", lineIndex, line)
			lineIndex += 1
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "error reading file: %s", err)
			os.Exit(1)
		}
	}
}
