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

_(fill me in)_
