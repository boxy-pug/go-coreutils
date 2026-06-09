# ccuniq

A Go clone of the Unix `uniq` command, built as part of the [Build Your Own uniq Tool](https://codingchallenges.fyi/challenges/challenge-uniq) challenge.

ccuniq filters adjacent matching lines from input, with options to count, show only repeated, or show only unique lines.

## Usage

```sh
# Remove adjacent duplicate lines
ccuniq test.txt

# Count occurrences of each line
ccuniq -c test.txt

# Show only repeated lines
ccuniq -d test.txt

# Show only unique lines
ccuniq -u test.txt

# Read from stdin, write to file
cat test.txt | ccuniq - out.txt
```

## Flags

- `-c` — prefix lines with the number of occurrences
- `-d` — only print repeated lines
- `-u` — only print unique (non-repeated) lines

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/ccuniq@latest
```

## What I learned

_(fill me in)_
