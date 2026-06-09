# cctail

A Go clone of the Unix `tail` command. Built independently — not part of a Coding Challenges prompt.

cctail prints the last part of files — by default the last 10 lines. Reads files backwards in chunks to handle large files efficiently.

## Usage

```sh
# Print last 10 lines of a file
cctail test.txt

# Print last N lines
cctail -n 5 test.txt

# Print last lines from multiple files
cctail -n 5 file1.txt file2.txt
```

## Flags

- `-n` — number of lines to print from the end (default: 10)

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/cctail@latest
```

## What I learned

_(fill me in)_
