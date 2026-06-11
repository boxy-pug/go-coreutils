# ccwc

A Go clone of the Unix `wc` (word count) command, made for the [Build Your Own wc Tool](https://codingchallenges.fyi/challenges/challenge-wc) challenge.

ccwc counts lines, words, bytes, and characters in files or from stdin. Unicode-aware character counting via runes.

## Usage

```sh
# Default: count lines, words, and bytes
ccwc test.txt

# Count only lines
ccwc -l test.txt

# Count only words
ccwc -w test.txt

# Count only bytes
ccwc -c test.txt

# Count characters (Unicode code points)
ccwc -m test.txt

# Read from stdin
cat test.txt | ccwc -l
```

## Flags

- `-l` — count lines
- `-w` — count words
- `-c` — count bytes
- `-m` — count characters (runes)

If no flags are given, ccwc defaults to `-l -w -c`.

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/ccwc@latest
```

## What I learned

- `io.Reader` is great for reading from many different sources, including `os.Stdin` and `os.File`.
- When opening files in a function, return a cleanup function to close them, and use `defer` in `main` to execute it.
- Ranging over values in Go creates copies; to modify originals, access them by index.
- Runes represent characters in Go; one emoji is one rune but several bytes: `len([]rune("😊")) == 1` but `len("😊") == 4`.
- To provide a custom `--help` message with the `flag` package, redefine `flag.Usage`.
- `os.Stderr` is used for error messages and help text — visible even when stdout is redirected.
- I used to do integration tests against the system wc like `exec.Command("wc", ...)`, but that was brittle and broke bcs of different column widths in different implementations. Hardcoding expectations makes it deterministic.
