// ccsort sort lines of input text
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	var start time.Time
	var duration time.Duration
	var err error

	cfg, cleanup, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading config: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	lines, err := cfg.getLines()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading file: %v", err)
		os.Exit(1)
	}

	if cfg.test {
		testAllSortingFuncs(lines)
		return
	}

	if cfg.verbose {
		start = time.Now()
	}

	// already checked that it exists in loadConfig
	sortFunc := sortFunctions[cfg.algo]
	sortFunc(lines)

	if cfg.verbose {
		duration = time.Since(start)
	}

	if cfg.unique {
		cfg.printUniqueLines(lines)
	} else {
		cfg.printLines(lines)
	}

	if cfg.verbose {
		// possible bug printing cfg.files doesnt print the names?
		fmt.Fprintf(os.Stderr, "Sorting %v with %v sort algo took %v\n", cfg.files, cfg.algo, duration)
	}
}

func (cfg *config) printLines(list []string) {
	for _, line := range list {
		fmt.Fprintln(cfg.out, line)
	}
}

func (cfg *config) printUniqueLines(list []string) {
	seen := make(map[string]bool)
	for _, line := range list {
		if seen[line] {
			continue
		}
		fmt.Fprintln(cfg.out, line)
		seen[line] = true
	}
}

func testAllSortingFuncs(originalList []string) {
	var wg sync.WaitGroup
	var results = make(map[string]time.Duration)
	listLength := len(originalList)
	mu := sync.Mutex{}

	for name, algo := range sortFunctions {
		wg.Add(1)
		go func(name string, algo func([]string)) {
			defer wg.Done()

			list := make([]string, listLength)
			copy(list, originalList)

			start := time.Now()
			algo(list)
			dur := time.Since(start)

			mu.Lock()
			results[name] = dur
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

// helper for sorting the output of testing the algos by the time they took to execute
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

// getLines loads all lines from the io.Readers into a []string.
func (cfg *config) getLines() ([]string, error) {
	lines := []string{}

	for _, file := range cfg.files {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading input: %w", err)
		}
	}
	return lines, nil
}
