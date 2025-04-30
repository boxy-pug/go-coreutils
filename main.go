package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boxy-pug/ccsort/config"
	"github.com/boxy-pug/ccsort/sortalgos"
)

func main() {
	var start time.Time
	var duration time.Duration

	cfg := config.LoadConfig()

	if cfg.Verbose {
		start = time.Now()
	}

	sortFunc, exists := sortalgos.SortFunctions[cfg.SortingAlgo]
	if exists {
		sortFunc(cfg.List)
	} else {
		log.Fatalf("could not find %v in sort functions", cfg.SortingAlgo)
	}

	if cfg.Verbose {
		duration = time.Since(start)
	}

	printLines(cfg.List)

	if cfg.Verbose {
		fmt.Printf("Sorting %s with %v sort algo took %v\n", cfg.File.Name(), cfg.SortingAlgo, duration)
	}
}

func printLines(list []string) {
	for _, line := range list {
		fmt.Println(line)
	}
}
