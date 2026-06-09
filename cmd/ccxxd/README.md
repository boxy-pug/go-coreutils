# ccxxd

A Go clone of the Unix `xxd` hex dump utility, built as part of the [Build Your Own Xxd](https://codingchallenges.fyi/challenges/challenge-xxd) challenge.

ccxxd creates hex dumps of files and reverses them back to binary. Supports grouping, custom column widths, little-endian output, and seeking.

## Usage

```sh
# Hex dump a file
ccxxd myfile.bin

# Custom columns and byte groups
ccxxd -c 8 -g 4 myfile.bin

# Little-endian output
ccxxd -e myfile.bin

# Limit output length
ccxxd -l 32 myfile.bin

# Start at an offset
ccxxd -s 512 -l 16 myfile.bin

# Revert hex dump back to binary
ccxxd -r hexdump.txt > restored.bin

# Hex dump from stdin
cat myfile.bin | ccxxd
```

## Flags

- `-e` — little-endian byte order within each group
- `-r` — revert: convert hex dump back to binary
- `-g` — group bytes per group (default: 2)
- `-c` — bytes per line (default: 16)
- `-l` — limit output to N bytes
- `-s` — skip N bytes from the start before dumping

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/ccxxd@latest
```

## Testing

Unit tests: `go test ./cmd/ccxxd/`

Integration tests compare output to system `xxd`:
```sh
go test -tags=integration ./cmd/ccxxd/
```

## What I learned

- All files are just bytes. Text, images, programs — different sequences of bytes on disk. Reverting a hex dump restores the original exactly, even executables work after `chmod +x`.
- `echo -n` suppresses trailing newlines when creating test data. `printf` is more consistent and also doesn't add one.
- In shell, single quotes treat everything as literal. Double quotes interpret escape sequences like `\n`.
- `io.ReadFull` keeps reading until the buffer is full or EOF. `bufio.Reader.Read` can return fewer bytes than requested mid-stream, which caused short lines in the middle of dumps.
- Use `strings.Builder` to assemble output in memory before a single write to stdout. Many small writes are expensive.
- Power of two check: `n > 0 && (n & (n-1)) == 0`. Only binary numbers with one positive bit are powers of two.
