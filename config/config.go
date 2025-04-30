package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type SortingAlgo string

const (
	StdLib    SortingAlgo = "stdlib"
	Bubble    SortingAlgo = "bubble"
	Merge     SortingAlgo = "merge"
	Insertion SortingAlgo = "insertion"
	Quick     SortingAlgo = "quick"
)

type Config struct {
	File        *os.File
	List        []string
	Unique      bool
	Verbose     bool
	SortingAlgo SortingAlgo
}

func LoadConfig() Config {
	var err error
	var cfg Config
	var algo string
	lineSet := make(map[string]struct{})

	flag.BoolVar(&cfg.Unique, "u", false, "only output unique lines")
	flag.BoolVar(&cfg.Verbose, "v", false, "verbose mode")
	flag.StringVar(&algo, "algo", "stdlib", "choose sorting algo: stdlib, bubble, merge, insertion, quick")

	flag.Parse()
	args := flag.Args()

	cfg.SortingAlgo = parseAlgo(algo)

	if cfg.Verbose {
		fmt.Printf("Using %v sorting algo\n", cfg.SortingAlgo)
	}

	switch len(args) {
	case 0:
		cfg.File = os.Stdin
	case 1:
		cfg.File, err = os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	scanner := bufio.NewScanner(cfg.File)
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.Unique {
			if _, exists := lineSet[line]; exists {
				continue
			}
			lineSet[line] = struct{}{}

		}
		cfg.List = append(cfg.List, line)

	}
	return cfg
}

func parseAlgo(str string) SortingAlgo {
	cleanStr := strings.TrimSpace(strings.ToLower(str))
	switch cleanStr {
	case "stdlib":
		return StdLib
	case "bubble":
		return Bubble
	case "merge":
		return Merge
	case "insertion":
		return Insertion
	case "quick":
		return Quick
	default:
		return StdLib
	}
}
