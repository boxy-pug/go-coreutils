# cccat

A Go clone of the Unix `cat` command, built as part of the [Build Your Own cat Tool](https://codingchallenges.fyi/challenges/challenge-cat) challenge.

cccat reads files sequentially and writes their contents to standard output. Supports numbering lines and reading from stdin.

## Usage

```sh
# Print a file
cccat test.txt

# Concatenate multiple files
cccat test.txt test2.txt

# Read from stdin
echo "hello" | cccat

# Read from stdin explicitly
echo "hello" | cccat -
```

## Flags

- `-n` — number all output lines
- `-b` — number non-empty output lines (overrides `-n`)

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/cccat@latest
```

## What I learned

_(fill me in)_
