# ccgrep

A Go clone of the Unix `grep` command, built as part of the [Build Your Own grep](https://codingchallenges.fyi/challenges/challenge-grep) challenge.

ccgrep searches input files for lines matching a pattern and prints matching lines.

## Usage

```sh
# Search for a pattern in a file
ccgrep "hello" test.txt

# Search with inverted match
ccgrep -v "hello" test.txt

# Case-insensitive search
ccgrep -i "HELLO" test.txt

# Recursive search in a directory
ccgrep -r "pattern" .
```

## Flags

- `-r` — recurse directory tree
- `-v` — invert match (print non-matching lines)
- `-i` — case insensitive matching

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/ccgrep@latest
```

## What I learned

_(fill me in)_
