package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Command holds config and state for the wc command
type Command struct {
	Output           io.Writer
	Files            []FileInput
	TotalCounter     WordCounter
	BytesFlag        bool
	LinesFlag        bool
	WordsFlag        bool
	CharsFlag        bool
	FileNameProvided bool
}

// FileInput represents single file or input stream
type FileInput struct {
	FileName string
	Counter  WordCounter
	Input    io.Reader
}

// WordCounter tracks the count
type WordCounter struct {
	Lines int
	Words int
	Chars int
	Bytes int
}

func main() {
	cmd, cleanup, err := loadCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading command:", err)
		os.Exit(1)
	}
	defer cleanup()

	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(cmd.Output, "error running wc command:", err)
		os.Exit(1)
	}
}

// loadCommand parses cmdline flags and input files, returns configured command
func loadCommand() (Command, func(), error) {
	cmd := Command{
		Output: os.Stdout,
	}

	flag.BoolVar(&cmd.BytesFlag, "c", false, "count bytes")
	flag.BoolVar(&cmd.LinesFlag, "l", false, "count lines")
	flag.BoolVar(&cmd.WordsFlag, "w", false, "count words")
	flag.BoolVar(&cmd.CharsFlag, "m", false, "count chars")

	// Customize the usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A simple word count tool similar to Unix wc command.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	// If no flags provided, enable standard wc options: lines, words and bytes
	if !cmd.BytesFlag && !cmd.LinesFlag && !cmd.WordsFlag && !cmd.CharsFlag {
		cmd.LinesFlag, cmd.WordsFlag, cmd.BytesFlag = true, true, true
	}

	var cleanup func() = func() {}

	switch {
	// no files provided: use stdin
	case len(args) == 0:
		cmd.FileNameProvided = false
		cmd.Files = append(cmd.Files, FileInput{
			Input: os.Stdin,
		})
	case len(args) > 0:
		var files []*os.File
		for _, a := range args {
			file, err := os.Open(a)
			if err != nil {
				return cmd, cleanup, fmt.Errorf("couldn't open file %v, error: %v", a, err)
			}
			files = append(files, file)
			cmd.Files = append(cmd.Files, FileInput{
				FileName: file.Name(),
				Input:    file,
			})
		}
		cmd.FileNameProvided = true

		cleanup = func() {
			for _, f := range files {
				f.Close()
			}
		}
	}
	return cmd, cleanup, nil
}

// Run processes each input, updates count and prints result
func (cmd *Command) Run() error {
	for i := range cmd.Files {
		input := cmd.Files[i]
		reader := bufio.NewReader(input.Input)

		for {
			line, err := reader.ReadString('\n')

			// Break if we're at EOF and nothing left to process
			if line == "" && err == io.EOF {
				break
			}

			if cmd.LinesFlag {
				input.Counter.Lines++
			}
			if cmd.WordsFlag {
				input.Counter.Words += len(strings.Fields(line))
			}
			if cmd.BytesFlag {
				input.Counter.Bytes += len(line)
			}
			if cmd.CharsFlag {
				input.Counter.Chars += utf8.RuneCountInString(line)
			}

			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}
		printResult(input.Counter, *cmd, input.FileName)

		if len(cmd.Files) > 1 {
			cmd.addCountToTotal(input.Counter)
		}
	}
	if len(cmd.Files) > 1 {
		printResult(cmd.TotalCounter, *cmd, "total")
	}
	return nil
}

// printResult prints the count for each result and total
func printResult(counter WordCounter, cmd Command, fileName string) {
	if cmd.LinesFlag {
		fmt.Fprintf(cmd.Output, "%8d", counter.Lines)
	}
	if cmd.WordsFlag {
		fmt.Fprintf(cmd.Output, "%8d", counter.Words)
	}
	if cmd.BytesFlag {
		fmt.Fprintf(cmd.Output, "%8d", counter.Bytes)
	}
	if cmd.CharsFlag {
		fmt.Fprintf(cmd.Output, "%8d", counter.Chars)
	}
	if cmd.FileNameProvided {
		fmt.Fprintf(cmd.Output, " %s", fileName)
	}
	fmt.Fprintln(cmd.Output)
}

// addCountToTotal accumulates count for the total line when multiple files are provided
func (cmd *Command) addCountToTotal(input WordCounter) {
	if cmd.LinesFlag {
		cmd.TotalCounter.Lines += input.Lines
	}
	if cmd.WordsFlag {
		cmd.TotalCounter.Words += input.Words
	}
	if cmd.BytesFlag {
		cmd.TotalCounter.Bytes += input.Bytes
	}
	if cmd.CharsFlag {
		cmd.TotalCounter.Chars += input.Chars
	}
}
