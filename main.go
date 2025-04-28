package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Config struct {
	file         *os.File
	list         []string
	unique       bool
	stdLibSorted []string
}

func main() {

	cfg := loadConfig()

	cfg.stdLibSort()

	cfg.print()

}

func loadConfig() Config {
	var err error
	var cfg Config
	lineSet := make(map[string]struct{})

	flag.BoolVar(&cfg.unique, "u", false, "only output unique lines")
	flag.Parse()
	args := flag.Args()

	switch len(args) {
	case 0:
		cfg.file = os.Stdin
	case 1:
		cfg.file, err = os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	scanner := bufio.NewScanner(cfg.file)
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.unique {
			if _, exists := lineSet[line]; exists {
				continue
			}
			lineSet[line] = struct{}{}

		}
		cfg.list = append(cfg.list, line)

	}
	return cfg
}

func (cfg *Config) stdLibSort() {
	slices.SortFunc(cfg.list, func(a, b string) int {
		return strings.Compare(a, b)
	})
}

func (cfg *Config) print() {
	for _, line := range cfg.list {
		fmt.Println(line)
	}
}
