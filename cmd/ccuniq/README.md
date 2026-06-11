# ccuniq

A Go clone of the Unix `uniq` command, built for the [Build Your Own uniq Tool](https://codingchallenges.fyi/challenges/challenge-uniq) challenge.

ccuniq filters out adjacent matching lines from input, with options to count, show only repeated, show only unique lines or do case insensitive comparison.

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
                                                                                                               t
# Case-insensitive comparison
ccuniq -i test.txt

# Read from stdin, write to file
cat test.txt | ccuniq - out.txt
```

## Flags

- `-c`, `--count` — prefix lines with the number of occurrences
- `-d`, `--repeated` — only print repeated lines
- `-u`, `--unique` — only print unique (non-repeated) lines
- `-i`, `--ignore-case` — ignore case when comparing lines

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/ccuniq@latest
```

## What I learned

- The `uniq` command doesn't sort the input, if you want to find duplicated lines across your input use sort -u instead. `uniq` only checks for adjacent duplicate lines.
- That means that a map is not a good match for implementing uniq, better to do it in a single pass and always checking if line == prevLine and the count of prevLineCount as the main way to decide what to do.
- `uniq` also only takes one input, if you provide two paths it will use the second ones as output, so like `uniq input.txt output.txt` not really matching the classic coreutils pattern.
- For adding short and longform flags in the go flag package, just make duplicate flag entries, pointing to same variable. go/flag also doesn't care about "-" vs "--" hyphens before a flag. and as a result you cannot combine short form flags like "-cd". So that's not optimal, should prbably hand roll the flag parsing, make my own package to import, or just use external lib.
- The cleanup func crashed when no files are opened. The fix is to initialize it as a "no-op" just empty func: func() {}.
- For fixed width column in go formatting: ("%7d", num) -> width will always be 7 digits wide, left padded.
- `strings.EqualFold(a, b string) bool` is great for checking case insensitive equality. Simple unicode case folding. If you do a classic strings.ToLower(s) comparison it creates new strings, EqualFold compares bytes in place. Also handles more unicode eddge cases.
