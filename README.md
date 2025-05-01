# ccsort

Ccsort is a command-line application developed as part of a learning project to explore sorting algorithms and CLI tools in Go, part of the [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-sort) series.

## Features

-  **Multiple Sorting Algorithms**: Choose from various sorting algorithms including standard library sort, bubble sort, merge sort, insertion sort, quick sort, and selection sort.
-  **Unique Line Output**: Option to output only unique lines from the input.
-  **Verbose Mode**: Get more detailed information about the sorting process, including timing.
-  **Test Mode**: Test all sorting algorithms to see their performance on the provided input.

## Installation

To use this tool, you need to have Go installed on your system. Clone the repository and build the tool using the following commands:

```bash
git clone https://github.com/boxy-pug/ccsort.git
go build .
```

## Usage

Run the tool from the command line with the following options:

```bash
./ccsort [flags] [input_file]
```

### Flags

-  `-u`: Only output unique lines.
-  `-v`: Enable verbose mode to display detailed information about the sorting process.
-  `--test`: Test all sorting algorithms on the input data and display their performance.
-  `--algo`: Choose the sorting algorithm to use. Options are `stdlib`, `bubble`, `merge`, `insertion`, `quick`, `selection`.

### Examples

1. **Sort a File Using Quick Sort**:
   ```bash
   ./ccsort --algo quick input.txt
   ```

2. **Sort a File with Unique Lines and Verbose Output**:
   ```bash
   ./ccsort -u -v --algo merge input.txt
   ```

3. **Test All Sorting Algorithms**:
   ```bash
   ./ccsort --test input.txt
   ```

## Learning Objectives

This project was created as a learning exercise to understand:

-  How different sorting algorithms work and their performance characteristics.
-  How to build command-line tools in Go.
-  How to use goroutines for concurrent execution in Go.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests with improvements or new features.

## License

This project is open-source and available under the [MIT License](LICENSE).