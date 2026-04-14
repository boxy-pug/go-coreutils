package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type command struct {
	files         []inputFile
	output        io.Writer
	lines         int
	bytes         int
	useLines      bool
	useBytes      bool
	multipleFiles bool
	stdIn         bool
}

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
		// Standard configuration, no flags provided
		useLines: true,
		output:   os.Stdout,
	}

	flag.IntVar(&cmd.lines, "n", 10, "number of lines to print")
	flag.IntVar(&cmd.bytes, "c", 0, "number of bytes to print")

	flag.Parse()
	args := flag.Args()

	cleanup := func() {}

	switch {
	case len(args) == 0:
		cmd.files = append(cmd.files, inputFile{
			reader: os.Stdin,
		})
		cmd.multipleFiles = false
		cmd.stdIn = true
	case len(args) > 0:
		var files []*os.File
		for _, a := range args {
			file, err := os.Open(a)
			if err != nil {
				return cmd, cleanup, fmt.Errorf("could not open %v as file, error: %v", a, err)
			}
			files = append(files, file)
			cmd.files = append(cmd.files, inputFile{
				name:   file.Name(),
				reader: file,
			})
		}
		cleanup = func() {
			for _, f := range files {
				f.Close()
			}
		}
		cmd.multipleFiles = len(args) > 1
	}

	// Disable line flag when byte flag is enabled
	if cmd.bytes > 0 && cmd.lines == 10 {
		cmd.useLines = false
		cmd.useBytes = true
	}
	return cmd, cleanup, nil
}

func (cmd *command) run() error {
	for i, file := range cmd.files {
		if cmd.multipleFiles {
			fmt.Fprintf(cmd.output, "==> %s <==\n", file.name)
		}

		var err error
		if cmd.useBytes {
			err = printHeadBytes(file.reader, cmd.output, cmd.bytes)
		} else {
			err = printHeadLines(file.reader, cmd.output, cmd.lines)
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
			w.Write(line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
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
