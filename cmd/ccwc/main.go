// ccwc counts lines, words, bytes, and characters in files or from stdin.
// Unicode-aware character counting via runes.
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

// command holds config and state for the wc command
type command struct {
	out              io.Writer
	files            []fileInput
	total            counter
	bytesFlag        bool // -c for bytes
	linesFlag        bool // -l for lines
	wordsFlag        bool // -w for words
	charsFlag        bool // -m for chars
	fileNameProvided bool
}

// fileInput represents single file or input stream
type fileInput struct {
	name    string
	counter counter
	input   io.Reader
}

// counter tracks the count
type counter struct {
	lines int
	words int
	chars int
	bytes int
}

func main() {
	cmd, cleanup, err := loadCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading command: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	err = cmd.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "running wc command: %v\n", err)
		os.Exit(1)
	}
}

// loadCommand parses cmdline flags and input files, returns configured command
func loadCommand() (command, func(), error) {
	cmd := command{
		out: os.Stdout,
	}

	flag.BoolVar(&cmd.bytesFlag, "c", false, "count bytes")
	flag.BoolVar(&cmd.linesFlag, "l", false, "count lines")
	flag.BoolVar(&cmd.wordsFlag, "w", false, "count words")
	flag.BoolVar(&cmd.charsFlag, "m", false, "count chars")

	// Customize the usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A simple word count tool similar to Unix wc command.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	paths := flag.Args()

	// If no flags provided, enable standard wc options: lines, words and bytes
	if !cmd.bytesFlag && !cmd.linesFlag && !cmd.wordsFlag && !cmd.charsFlag {
		cmd.linesFlag, cmd.wordsFlag, cmd.bytesFlag = true, true, true
	}

	var cleanup func() = func() {}

	switch {
	// no files provided: use stdin
	case len(paths) == 0:
		cmd.fileNameProvided = false
		cmd.files = append(cmd.files, fileInput{
			input: os.Stdin,
		})
	case len(paths) > 0:
		var files []*os.File
		for _, path := range paths {
			// Support "-" for stdin, can be combined with files, with name "-"
			if path == "-" {
				cmd.fileNameProvided = true
				cmd.files = append(cmd.files, fileInput{
					name:  path,
					input: os.Stdin,
				})
				continue
			}
			file, err := os.Open(path)
			if err != nil {
				return cmd, cleanup, fmt.Errorf("open %q: %w", path, err)
			}
			files = append(files, file)
			cmd.files = append(cmd.files, fileInput{
				name:  file.Name(),
				input: file,
			})
		}
		cmd.fileNameProvided = true

		cleanup = func() {
			for _, f := range files {
				f.Close()
			}
		}
	}
	return cmd, cleanup, nil
}

// run processes each input, updates count and prints result
func (cmd *command) run() error {
	for i := range cmd.files {
		input := cmd.files[i]
		reader := bufio.NewReader(input.input)

		for {
			line, err := reader.ReadString('\n')

			// Break if we're at EOF and nothing left to process
			if line == "" && err == io.EOF {
				break
			}

			if cmd.linesFlag {
				input.counter.lines++
			}
			if cmd.wordsFlag {
				input.counter.words += len(strings.Fields(line))
			}
			if cmd.bytesFlag {
				input.counter.bytes += len(line)
			}
			if cmd.charsFlag {
				input.counter.chars += utf8.RuneCountInString(line)
			}

			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("reading input: %w", err)
			}
		}
		printResult(input.counter, *cmd, input.name)

		if len(cmd.files) > 1 {
			cmd.addCountToTotal(input.counter)
		}
	}
	if len(cmd.files) > 1 {
		printResult(cmd.total, *cmd, "total")
	}
	return nil
}

// printResult prints the count for each result and total
func printResult(counter counter, cmd command, fileName string) {
	if cmd.linesFlag {
		fmt.Fprintf(cmd.out, "%8d", counter.lines)
	}
	if cmd.wordsFlag {
		fmt.Fprintf(cmd.out, "%8d", counter.words)
	}
	if cmd.bytesFlag {
		fmt.Fprintf(cmd.out, "%8d", counter.bytes)
	}
	if cmd.charsFlag {
		fmt.Fprintf(cmd.out, "%8d", counter.chars)
	}
	if cmd.fileNameProvided {
		fmt.Fprintf(cmd.out, " %s", fileName)
	}
	fmt.Fprintln(cmd.out)
}

// addCountToTotal accumulates count for the total line when multiple files are provided
func (cmd *command) addCountToTotal(input counter) {
	if cmd.linesFlag {
		cmd.total.lines += input.lines
	}
	if cmd.wordsFlag {
		cmd.total.words += input.words
	}
	if cmd.bytesFlag {
		cmd.total.bytes += input.bytes
	}
	if cmd.charsFlag {
		cmd.total.chars += input.chars
	}
}
