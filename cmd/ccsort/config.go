package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type config struct {
	files   []io.Reader
	out     io.Writer
	unique  bool
	verbose bool
	test    bool
	algo    string
}

func loadConfig() (config, func(), error) {
	cfg := config{out: os.Stdout}
	cleanup := func() {}
	var err error

	flag.BoolVar(&cfg.unique, "u", false, "only output unique lines")
	flag.BoolVar(&cfg.verbose, "v", false, "verbose mode")
	flag.BoolVar(&cfg.test, "test", false, "test all algos")
	flag.StringVar(&cfg.algo, "algo", "stdlib", "choose sorting algo: stdlib, bubble, merge, insertion, quick, selection")
	// TODO: support flags like --mergesort, --qsort etc instead of the --algo solution?

	flag.Parse()
	paths := flag.Args()

	cfg.algo = strings.TrimSpace(strings.ToLower(cfg.algo))

	// If the algo name doesn't exist, default to stdlib implementation
	if _, exists := sortFunctions[cfg.algo]; !exists {
		cfg.algo = "stdlib"
	}

	if cfg.verbose {
		fmt.Fprintf(os.Stderr, "Using %v sorting algo\n", cfg.algo)
	}

	switch {
	case len(paths) == 0: // zero args, read from stdin
		cfg.files = append(cfg.files, os.Stdin)
	default: // multiple paths, append them all, doesn't support paths and "-" mixed
		cfg.files, cleanup, err = openPaths(paths)
		if err != nil {
			return config{}, cleanup, fmt.Errorf("opening input paths: %w", err)
		}
	}
	return cfg, cleanup, nil
}

// openPaths takes a []string of paths and opens them as io.Reader. Also handles
// "-" as stdin, and returns a cleanup func for closing files, and an error.
func openPaths(paths []string) ([]io.Reader, func(), error) {
	var files []io.Reader
	cleanup := func() {}

	for _, path := range paths {
		if path == "-" {
			files = append(files, os.Stdin)
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			cleanup() // close already opened files if error
			return nil, nil, fmt.Errorf("opening file %s: %w", path, err)
		}
		files = append(files, f)
	}

	// cleanup func that closes all inputs that satisfy closer interface
	cleanup = func() {
		for _, f := range files {
			if closer, ok := f.(io.Closer); ok && f != os.Stdin {
				closer.Close()
			}
		}
	}
	return files, cleanup, nil
}
