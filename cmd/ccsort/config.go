package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type config struct {
	filePaths []string
	unique    bool
	verbose   bool
	test      bool
	fromStdin bool
	algo      string
}

func loadConfig() config {
	var cfg config
	var algo string

	flag.BoolVar(&cfg.unique, "u", false, "only output unique lines")
	flag.BoolVar(&cfg.verbose, "v", false, "verbose mode")
	flag.BoolVar(&cfg.test, "test", false, "test all algos")
	flag.StringVar(&algo, "algo", "stdlib", "choose sorting algo: stdlib, bubble, merge, insertion, quick, selection")

	flag.Parse()
	args := flag.Args()

	cfg.algo = strings.TrimSpace(strings.ToLower(algo))

	// If the algo name doesn't exist, default to stdlib implementation
	if _, exists := sortFunctions[cfg.algo]; !exists {
		cfg.algo = "stdlib"
	}

	if cfg.verbose {
		fmt.Fprintf(os.Stderr, "Using %v sorting algo\n", cfg.algo)
	}

	// Just support one filepath or stdin for now
	switch {
	case len(args) == 0:
		cfg.fromStdin = true
	case len(args) == 1 && args[0] == "-":
		cfg.fromStdin = true
	default:
		for _, path := range args {
			cfg.filePaths = append(cfg.filePaths, strings.TrimSpace(path))
		}
	}
	return cfg
}
