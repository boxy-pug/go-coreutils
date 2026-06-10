# cccat

A Go clone of the Unix `cat` command, built as part of the [Build Your Own cat Tool](https://codingchallenges.fyi/challenges/challenge-cat) challenge.

cccat reads files sequentially and writes their contents to standard output. Supports multiple files, numbering lines, numbering non-empty lines only and reading from stdin.

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

- io.Copy is a zero-copy pipe that reads from reader and writes to writer in chunks, perfect for a cat tool, doesn't hold the whole files in memory.
- I first used a bufio.Scanner for reading line by line. Very handy for breaking input into tokens ("line" is the default token). But it drops the delimiter (strips the newline), so you dont get byte-by-byte fidelity, you need to re-add the newlines, can get messy fast, for trailing newlines, windows style \r\n etc.
- bufio.Reader is the better choice for keeping the newlines around, it preserves the raw bytes. ReadString('\n') gives you the line including the newline.
