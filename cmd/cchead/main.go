// cchead prints the first lines or bytes of a file.
// Defaults to first ten lines.
// Supports line-count (-n), byte-count (-c) and multiple files
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// command holds the parsed flags and input files.
// Holds io.Writer to simplify testing.
type command struct {
	files         []inputFile
	output        io.Writer
	lineCount     int
	byteCount     int
	wantBytes     bool // Toggle lines vs bytes, mutually exclusive
	multipleFiles bool
}

// inputFile pairs display name with io.Reader
// Needed for multiple files headers
type inputFile struct {
	name   string
	reader io.Reader
}

func main() {
	cmd, cleanup, err := loadCommand()
	if err != nil {
		fmt.Println("error loading command:", err)
		os.Exit(1)
	}
	defer cleanup()

	err = cmd.run()
	if err != nil {
		fmt.Fprintln(cmd.output, "error running command:", err)
	}
}

func loadCommand() (command, func(), error) {
	cmd := command{
		// Standard configuration, no flags provided, 10 lines, no byte limit
		wantBytes: false,
		output:    os.Stdout,
	}

	flag.IntVar(&cmd.lineCount, "n", 10, "number of lines to print")
	flag.IntVar(&cmd.byteCount, "c", 0, "number of bytes to print")

	flag.Parse()
	args := flag.Args()

	cleanup := func() {}

	switch {
	case len(args) == 0:
		cmd.files = append(cmd.files, inputFile{
			reader: os.Stdin,
		})
		cmd.multipleFiles = false
	case len(args) > 0:
		var openFiles []*os.File
		for _, path := range args {
			file, err := os.Open(path)
			if err != nil {
				return cmd, cleanup, fmt.Errorf("could not open %v as file, error: %v", path, err)
			}
			openFiles = append(openFiles, file)
			cmd.files = append(cmd.files, inputFile{
				name:   file.Name(),
				reader: file,
			})
		}
		cleanup = func() {
			for _, f := range openFiles {
				f.Close()
			}
		}
		cmd.multipleFiles = len(args) > 1
	}

	// Hack: if bytes is set and lines is still default, user probably want bytes.
	// Doesn't work for edge case -c 0 (user wants 0 bytes but gets 10 lines)
	// There actually is something called flag.Visit i could use for this
	if cmd.byteCount > 0 && cmd.lineCount == 10 {
		cmd.wantBytes = true
	}
	return cmd, cleanup, nil
}

func (cmd *command) run() error {
	for i, file := range cmd.files {

		// Print headers if multiple files
		if cmd.multipleFiles {
			fmt.Fprintf(cmd.output, "==> %s <==\n", file.name)
		}

		var err error
		if cmd.wantBytes {
			err = printHeadBytes(file.reader, cmd.output, cmd.byteCount)
		} else {
			err = printHeadLines(file.reader, cmd.output, cmd.lineCount)
		}

		if err != nil {
			return fmt.Errorf("error reading %q: %w", file.name, err)
		}

		// Print space between multiple files
		if cmd.multipleFiles && i < len(cmd.files)-1 {
			fmt.Println()
		}
	}
	return nil
}

func printHeadLines(r io.Reader, w io.Writer, n int) error {
	reader := bufio.NewReader(r)

	for range n {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			if _, werr := w.Write(line); werr != nil {
				return werr
			}
		}
		// We must write data to w before checking err,
		// otherwise we would drop the last line on EOF
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func printHeadBytes(r io.Reader, w io.Writer, n int) error {
	_, err := io.CopyN(w, r, int64(n))
	if err == io.EOF {
		return nil
	}
	return err
}
