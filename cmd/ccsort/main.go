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

	cfg := loadConfig()

	lines, err := getLines(cfg.filePaths, cfg.fromStdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %s", err)
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
		printUniqueLines(lines)
	} else {
		printLines(lines)
	}

	if cfg.verbose {
		fmt.Fprintf(os.Stderr, "Sorting %s with %v sort algo took %v\n", cfg.filePaths, cfg.algo, duration)
	}
}

func printLines(list []string) {
	for _, line := range list {
		fmt.Println(line)
	}
}

func printUniqueLines(list []string) {
	seen := make(map[string]bool)
	for _, line := range list {
		if seen[line] {
			continue
		}
		fmt.Println(line)
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

// getLines loads lines from stdin or files into a []string for sorting later.
func getLines(paths []string, fromStdin bool) ([]string, error) {
	var files []*os.File

	if fromStdin {
		files = []*os.File{os.Stdin}
	} else {
		for _, path := range paths {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
	}

	var lines []string
	for _, file := range files {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading input: %w", err)
		}
		if file != os.Stdin {
			defer file.Close()
		}
	}
	return lines, nil
}
