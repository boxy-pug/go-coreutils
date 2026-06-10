# ccsort

A Go clone of the Unix `sort` command, built as part of the [Build Your Own Sort Tool](https://codingchallenges.fyi/challenges/challenge-sort) challenge.

ccsort sorts lines from files using multiple sorting algorithms. Includes a test mode and verbose output for comparing algorithm performance.

## Usage

```sh
# Sort a file with default (stdlib) sort
ccsort input.txt

# Sort using a specific algorithm
ccsort --algo quick input.txt

# Sort with unique output and verbose mode
ccsort -u -v --algo merge input.txt

# Test all sorting algorithms
ccsort --test input.txt
```

## Flags

- `--algo` — sorting algorithm to use: `stdlib`, `bubble`, `merge`, `insertion`, `quick`, `selection`
- `--test` — benchmark all sorting algorithms on the input
- `-u` — only output unique lines
- `-v` — verbose mode (show timing and details)

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/ccsort@latest
```

## What I learned

- The sort from coreutil compares strings in a way that accounts for dictonary order, where upper vs lowercase is secondary. I've built a pure byte-by-byte version like LC_COLLATE=C, the c or posix locale: strict ascii/byte by byte ordering.
- log.Fatalf adds timestamp prefixes and is otherwise the same as fmt.Fprintln(os.Stderr, ...) followed by os.Exit(1). For use on stuff running on a server for example.
- Posix convention for reading from stdin is "-" as the filename.
- When using bufio.Scanner.Scan() it returns false for both EOF and read errors. Check scanner.Err() after each loop.
- Diagnostic messages like --verbose -v stuff should go to os.Stderr. Otherwise piping mixes timestamps into the data.
