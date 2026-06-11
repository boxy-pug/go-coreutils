// ccuniq reads an input file comparing adjacent lines, and writes each unique
// line to stdout. Supports counting repeated lines with -c,
// repeated lines (-d) or unique lines only with -u.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// config holds flags and input reader and output writer
type config struct {
	writeToFile  bool
	countCol     bool      // -c, --count
	repeatedOnly bool      // -d, --repeated
	uniqueOnly   bool      // -u --unique
	in           io.Reader // uniq only supports one input file
	out          io.Writer
}

func main() {
	var err error

	cfg, cleanup, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading config: %v", err)
		os.Exit(1)
	}
	defer cleanup()

	err = cfg.runUniq()
	if err != nil {
		fmt.Fprintf(os.Stderr, "running uniq: %v\n", err)
		os.Exit(1)
	}
}

// runUniq reads input line by line, tracking identical lines.
// Previous line is printed if it meets the conditions set by flags and linecount.
func (cfg *config) runUniq() error {
	scanner := bufio.NewScanner(cfg.in)

	var prevCount int
	var prevLine string

	for scanner.Scan() {
		line := scanner.Text()

		// line is same as previous, still in same group
		// increment counter and continue
		if line == prevLine {
			prevCount++
			continue
		}

		// New group started bcs line != prevLine, print the previous if it qualifies
		if cfg.shouldPrint(prevCount) {
			if cfg.countCol {
				fmt.Fprintf(cfg.out, "%7d ", prevCount)
			}
			fmt.Fprintln(cfg.out, prevLine)
		}

		// Start a new group
		prevLine = line
		prevCount = 1
	}

	// Handle last group after EOF
	if cfg.shouldPrint(prevCount) {
		if cfg.countCol {
			fmt.Fprintf(cfg.out, "%7d ", prevCount)
		}
		fmt.Fprintln(cfg.out, prevLine)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	return nil
}

// shouldPrint determines if a completed group of lines
// should be printed, based on count and flags.
func (cfg *config) shouldPrint(count int) bool {
	if count == 0 {
		return false // first iteration guard
	}
	if cfg.uniqueOnly {
		return count == 1
	}
	if cfg.repeatedOnly {
		return count >= 2
	}
	return true // default uniq: print any group once
}

// loadConfig parses flags and opens input file and optional output file.
func loadConfig() (config, func(), error) {
	cfg := config{out: os.Stdout}
	cleanup := func() {} // need to initialize this to not crash when no file paths provided
	var err error

	flag.BoolVar(&cfg.countCol, "c", false, "enable count column")
	flag.BoolVar(&cfg.countCol, "count", false, "enable count column")
	flag.BoolVar(&cfg.repeatedOnly, "d", false, "output only repeated")
	flag.BoolVar(&cfg.repeatedOnly, "repeated", false, "output only repeated")
	flag.BoolVar(&cfg.uniqueOnly, "u", false, "output only unique")
	flag.BoolVar(&cfg.uniqueOnly, "unique", false, "output only unique")

	flag.Parse()
	paths := flag.Args()

	switch len(paths) {
	case 0: // no path provided -> read from stdin
		cfg.in = os.Stdin
	case 1: // one path provided
		// if its "-" read from stdin
		if paths[0] == "-" {
			cfg.in = os.Stdin
			break
		}
		// if not, open path as file
		file, err := os.Open(paths[0])
		if err != nil {
			return config{}, cleanup, fmt.Errorf("opening file: %w", err)
		}
		cfg.in = file
		cleanup = func() { file.Close() }
	default: // 2 or more paths provided
		// uniq is a bit unique (lol) bcs it treats first path as in and second as out
		// not following the common coreutil pattern of multiple inputs
		if paths[0] == "-" {
			cfg.in = os.Stdin
		} else {
			file, err := os.Open(paths[0])
			if err != nil {
				return config{}, cleanup, fmt.Errorf("opening file: %w", err)
			}
			cfg.in = file
			cleanup = func() { file.Close() }
		}
		// open second path as file to write to
		cfg.writeToFile = true
		cfg.out, err = os.Create(paths[1])
		if err != nil {
			return config{}, cleanup, fmt.Errorf("creating file: %w", err)
		}
	}
	return cfg, cleanup, nil
}
