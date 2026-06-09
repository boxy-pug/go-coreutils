# ccsed

A Go clone of the Unix `sed` stream editor, built as part of the [Build Your Own Sed](https://codingchallenges.fyi/challenges/challenge-sed) challenge.

ccsed parses and transforms text using pattern-based substitution commands.

## Usage

```sh
# Substitute pattern in a file
ccsed 's/old/new/g' test.txt

# Print only lines matching a pattern
ccsed -n '/pattern/p' test.txt

# Print a range of lines
ccsed -n '2,4p' test.txt

# Edit file in place
ccsed -i 's/old/new/g' test.txt
```

## Flags

- `-n` — only print explicitly selected lines
- `-i` — edit files in place

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/ccsed@latest
```

## What I learned

_(fill me in)_
