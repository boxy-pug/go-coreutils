// cccat concatenates files, prints to standard out.
// It copies bytes unchanged by default. Flags -n or -b
// reads line by line and adds numbering
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// config holds flags, file readers and out writer
type config struct {
	numberLines    bool        // -n flag
	numberNonBlank bool        // -b flag
	files          []io.Reader // files to read
	out            io.Writer   // stdout, makes testing simpler
}

func main() {

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cccat: loading config: %s\n", err)
		os.Exit(1)
	}

	for _, file := range cfg.files {
		err := cfg.catInput(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cccat: %s\n", err)
			os.Exit(1)
		}
		// Close if it's a closer, not stdin
		if closer, ok := file.(io.Closer); ok && file != os.Stdin {
			closer.Close()
		}
	}
}

func loadConfig() (config, error) {
	cfg := config{
		out: os.Stdout,
	}
	flag.BoolVar(&cfg.numberLines, "n", false, "number lines")
	flag.BoolVar(&cfg.numberNonBlank, "b", false, "number non blank lines")

	flag.Parse()
	paths := flag.Args()

	// "-" is posix for "read from stdin", same as empty file arg
	// TODO: support - for stdin mixed with files, like "cat file1 - file2"
	if len(paths) == 0 || paths[0] == "-" {
		cfg.files = []io.Reader{os.Stdin}
	} else {
		for _, path := range paths {
			f, err := os.Open(path)
			if err != nil {
				return config{}, fmt.Errorf("opening %s: %w", path, err)
			}
			cfg.files = append(cfg.files, f)
		}
	}
	return cfg, nil
}

// catInput reads single file or stdin and writes to cfg.out (stdout normally)
// Uses io.Copy for plain output and bufio.Reader for line by line + numbering.
func (cfg *config) catInput(file io.Reader) error {

	// Plain cat, just stream bytes
	// io.Copy preserves exact bytes and handles big files.
	if !cfg.numberLines && !cfg.numberNonBlank {
		_, err := io.Copy(cfg.out, file)
		if err != nil {
			return fmt.Errorf("reading input: %w", err)
		}
		return nil
	}

	// Numbered modes: read line-by-line to get the count
	// bufio.NewReader and ReadString('\n') preserves newlines
	reader := bufio.NewReader(file)
	lineIndex := 1

	for {
		line, err := reader.ReadString('\n')

		// needed this check so it doesnt print number when input is empty
		// ReadString return data and io.EOF on the last read.
		// If input is empty or if the file ends with newline, we must skip printing.
		// so len(line) > 0 guards against printing line number for empty data
		if len(line) > 0 {
			if cfg.numberNonBlank && line == "\n" {
				fmt.Fprintf(cfg.out, "%s", line)
			} else {
				fmt.Fprintf(cfg.out, "%6d\t%s", lineIndex, line)
				lineIndex += 1
			}
		}

		// ReadString returns EOF on last line, so we need to print before checking err, so we don't drop the last chunk.
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("reading input: %w", err)
		}
	}
	return nil
}
