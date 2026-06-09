# cctail – Tail command clone

This program is a Go implementation of the Unix `tail` command, which prints the last few lines of one or more files. It efficiently handles large files by reading them backwards in chunks.

## Features

-  **Print Last N Lines**: By default, prints the last 10 lines of each file. You can specify a different number of lines using the `-n` flag.
-  **Multiple Files**: Supports reading from multiple files and prints a header for each file when more than one file is specified.
-  **Standard Input**: Reads from standard input if no files are provided.

## Usage

```bash
./cctail [options] [file...]
```

### Options

-  `-n <number>`: Specify the number of lines to print from the end of each file (default is 10).

### Examples

-  Print the last 10 lines of `file.txt`:
  ```bash
  ./cctail file.txt
  ```

-  Print the last 5 lines of `file1.txt` and `file2.txt`:
  ```bash
  ./cctail -n 5 file1.txt file2.txt
  ```

## Implementation Details

-  **Backward Reading**: The program seeks to the end of each file and reads backwards in chunks, counting newlines until the desired number of lines is found.
-  **Handles Newlines at EOF**: Properly accounts for files that end with a newline, ensuring accurate line counting.
-  **Efficient Chunk Handling**: Uses a buffer to read chunks of data, minimizing memory usage and improving performance.

## License

This project is open source and available under the MIT License.

---

This README provides a clear overview of your `tail` command clone, explaining its features, usage, and implementation details.

