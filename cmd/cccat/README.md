# cccat

A Go clone of the Unix `cat` command, built for the [Build Your Own cat Tool](https://codingchallenges.fyi/challenges/challenge-cat) challenge.

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

- io.Copy is a zero-copy pipe that reads from reader and writes to writer in chunks, perfect for a cat tool, doesn't hold the whole files in memory, teh data flows through in small chunks, doenst accumulate into a single buffer.
- I first used a bufio.Scanner for reading line by line. Very handy for breaking input into tokens ("line" is the default token). But it drops the delimiter (strips the newline), so you dont get byte-by-byte fidelity, you need to re-add the newlines, can get messy fast, for trailing newlines, windows style \r\n etc.
- bufio.Reader is the better choice for keeping the newlines around, it preserves the raw bytes. ReadString('\n') gives you the line including the newline.
- Error printing, when to use %w vs %v. %w wraps the original error type, use this to build an error to return up the call stack. User can then do errors.Is(err, os.ErrNotExist) for example, handling specific errors. %v formats the error as a string, use it for printing to the user.
- Don't spell out "error" inside error strings, caller adds that later. Just describe the action that failed, like ("opening %s: %w", path, err) or "reading input", descriptive stuff like that.
- Don't os.Exit() in functions, only main shoudl call that, makes testing simpler, cleaner. Funcs return errors up the call stack, main can handle shutting down the program.
- Using io.Reader/Writer vs os.File and printing to stdout: io.Reader and Writer makes it more testable, easy to pass in a strings.NewReader to read from and bytes.Buffer to read to, instead of creating real files. io.reader accepts anything readable.
- Closing an io.Reader: not every io.Reader can be closed, like stdin or strings.NewReader. So you should use this type assertion pattern to check and only close files, this checks if it implements the Closer interface, which consists of one method: Close() -> error

```go
// Pattern for safely closing only files that were opened by code, skipping stdin and test readers
if closer, ok := file.(io.Closer); ok && file != os.Stdin {
closer.Close()
}
```

- I moved from package level var flags to a config struct. The `var Numberlines = flag.Bool()` thing worked, but it creates global mutable state, every function can reach out and read them, makes testing messy, bcs you can't easily just pass another config in testing without changing global vars.
- So when moving to a config struct, you put the flags and input and output readers and writers inside the struct and pass that around. Makes unittesting straight forward: Create a config with flags, pass it a strings.NewReader and assert against a bytes.Buffer.
- So when moving to a config struct, you put the flags and input and output readers and writers inside the struct and pass that around. Makes unittesting straight forward: Create a config with flags, pass it a strings.NewReader and assert against a bytes.Buffer.
