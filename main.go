package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/boxy-pug/ccsort/sortalgos"
)

type Config struct {
	file   *os.File
	list   []string
	unique bool
	//stdLibSorted []string
	bubbleSort    bool
	verbose       bool
	mergeSort     bool
	insertionSort bool
	quickSort     bool
}

func main() {
	var start time.Time
	var duration time.Duration

	cfg := loadConfig()

	if cfg.verbose {
		start = time.Now()
	}

	if cfg.bubbleSort {
		cfg.list = sortalgos.Bubble(cfg.list)
	} else if cfg.mergeSort {
		cfg.list = sortalgos.MergeSort(cfg.list)
	} else if cfg.insertionSort {
		cfg.list = sortalgos.InsertionSort(cfg.list)
	} else if cfg.quickSort {
		low := 0
		high := len(cfg.list) - 1
		sortalgos.QuickSort(cfg.list, low, high)
	} else {
		cfg.list = sortalgos.StdLib(cfg.list)
	}

	if cfg.verbose {
		duration = time.Since(start)
	}

	printLines(cfg.list)

	if cfg.verbose {
		fmt.Printf("Sorting %s took %v\n", cfg.file.Name(), duration)
	}

}

func loadConfig() Config {
	var err error
	var cfg Config
	lineSet := make(map[string]struct{})

	flag.BoolVar(&cfg.unique, "u", false, "only output unique lines")
	flag.BoolVar(&cfg.verbose, "v", false, "verbose mode")
	flag.BoolVar(&cfg.bubbleSort, "bubble", false, "use bubble sort")
	flag.BoolVar(&cfg.mergeSort, "merge", false, "merge sort")
	flag.BoolVar(&cfg.insertionSort, "insertion", false, "insertion sort")
	flag.BoolVar(&cfg.quickSort, "quick", false, "quick sort")
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

func printLines(list []string) {
	for _, line := range list {
		fmt.Println(line)
	}
}
