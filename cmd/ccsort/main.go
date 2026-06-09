package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/boxy-pug/ccsort/config"
	"github.com/boxy-pug/ccsort/sortalgos"
)

func main() {
	var start time.Time
	var duration time.Duration

	cfg := config.LoadConfig()

	if cfg.Test {
		testAllSortingFuncs(cfg.List)
		return
	}

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

func testAllSortingFuncs(originalList []string) {
	var wg sync.WaitGroup
	var results = make(map[string]time.Duration)
	listLength := len(originalList)
	mu := sync.Mutex{}

	for name, algo := range sortalgos.SortFunctions {
		wg.Add(1)
		go func(name config.SortingAlgo, algo func([]string)) {
			defer wg.Done()

			list := make([]string, listLength)
			copy(list, originalList)

			start := time.Now()
			algo(list)
			dur := time.Since(start)

			mu.Lock()
			results[string(name)] = dur
			mu.Unlock()

		}(name, algo)
	}
	wg.Wait()

	resList := sortMapByValue(results)

	for i, name := range resList {
		dur := results[name]

		fmt.Printf("%d: %v time for %v sort algo to sort list of %d elements\n", i+1, dur, name, listLength)
	}
}

func sortMapByValue(resMap map[string]time.Duration) []string {
	resList := make([]string, 0, len(resMap))
	for key := range resMap {
		resList = append(resList, key)
	}
	sort.Slice(resList, func(i, j int) bool {
		return resMap[resList[i]] < resMap[resList[j]]
	})
	return resList
}
