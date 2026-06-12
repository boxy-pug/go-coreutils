# ccsort

ccsort sorts lines from files using multiple sorting algorithms.

A Go clone of the Unix `sort` command, built for the [Build Your Own Sort Tool](https://codingchallenges.fyi/challenges/challenge-sort) challenge. I started out just making a `sort` clone, but found that this was a fun way to try out the different sorting algorithms i learned in the [Boot.dev course](https://www.boot.dev/courses/learn-data-structures-and-algorithms-python) on algorithms, and implementing them in Go instead of Python. I also read the [Wengrow book on algorithms](https://www.amazon.com/Common-Sense-Guide-Structures-Algorithms-Second/) to understand some of this stuff better. I also made a --test flag for running all the algos on (a copy of) the same input in parallel and timing it, to see how long the different algorithms actually take.

## Usage

```sh
# Sort a file with default (stdlib) sort
ccsort input.txt

# Sort using a specific algorithm
ccsort --algo quick input.txt

# Sort with unique output and verbose mode, use mergesort
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

- The sort from coreutil compares strings in a way that accounts for dictonary order, where upper vs lowercase is secondary. It has to do with "locale aware" string comp. I've built a pure byte-by-byte version like LC_COLLATE=C, the c or posix locale: strict ascii/byte by byte ordering. So a line starting with "B" (ascii 66) will come before "a" (ascii 97).
- log.Fatalf adds timestamp prefixes and is otherwise the same as fmt.Fprintln(os.Stderr, ...) followed by os.Exit(1). For use on stuff running on a server for example.
- Posix convention for reading from stdin is "-" as the filename.
- When using bufio.Scanner.Scan() it returns false for both EOF and read errors. Check scanner.Err() after each loop.
- Diagnostic messages like --verbose -v stuff should go to os.Stderr. Otherwise piping mixes timestamps into the data.
- I have learned a lot about the different sorting algorithms making this, check out comments in sort.go for info on the different algos.
