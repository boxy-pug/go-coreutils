# cchead

A Go clone of the Unix `head` command, built as part of the [Build Your Own head](https://codingchallenges.fyi/challenges/challenge-head) challenge.

cchead prints the first part of files — by default the first 10 lines.

## Usage

```sh
# Print first 10 lines of a file
cchead test.txt

# Print first N lines
cchead -n 5 test.txt

# Print first N bytes
cchead -c 100 test.txt

# Print first lines from multiple files
cchead -n 5 test.txt test2.txt
```

## Flags

- `-n` — number of lines to print (default: 10)
- `-c` — number of bytes to print

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/cchead@latest
```

## What I learned

- `bufio.Scanner` vs. `bufio.Reader`: Scanner is simpler but strips newlines, which causes issues with `\r\n` vs `\n`. Reader gives more control but requires careful EOF and error handling.
- For integration testing CLI tools, use the `exec` package to run local commands and capture their output, to compare your tool's output to the real command — very handy when building a clone.
