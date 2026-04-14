# 📝 ccwc
A Go implementation of the classic Unix `wc` (word count) tool, created for the [codingchallenges.fyi wc challenge](https://codingchallenges.fyi/challenges/challenge-wc/).

---

## 🚀 Features

-  Counts **lines**, **words**, **bytes**, and **characters** (Unicode-aware)
-  Supports **multiple files** and **stdin**
-  Output style matches GNU `wc`
-  Unit and integration tests

---

## 🛠️ Usage

```sh
# Count lines, words, and bytes (default) in one or more files
ccwc [file1] [file2] ...

# Count only lines
ccwc -l file.txt

# Count only words
ccwc -w file.txt

# Count only bytes
ccwc -c file.txt

# Count only characters (Unicode code points)
ccwc -m file.txt

# Combine flags (order doesn't matter)
ccwc -l -w -c file.txt

# Use stdin
cat file.txt | ccwc
```

---

## 🚩 Flags

| Flag | Description                    |
|------|--------------------------------|
| -l   | Count lines                    |
| -w   | Count words                    |
| -c   | Count bytes                    |
| -m   | Count characters (runes)       |

*If no flags are given, `ccwc` defaults to `-l -w -c` (lines, words, and bytes).*

---

## 💻 Example

```sh
$ ccwc -l -w -c test.txt
      10      42     512 test.txt
```

---

## 🧠 What I learned

-  **`io.Reader`** is great for reading from many different sources, including `os.Stdin` and `os.File`.

-  When opening files in a function, return a cleanup function to close the files, and use `defer` in `main` or parent func to execute it.

-  Ranging over values in Go creates copies; to modify originals, access them by index.

-  Runes represent chars in Go; one emoji or special char is one rune but several bytes: `len([]rune("😊")) == 1` but `len("😊") == 4`. (Unicode, utf8 encoding)

-  To provide a custom `--help` or usage message with the `flag` package, redefine `flag.Usage` with a function to print what you want.

-  `os.Stderr` is used for error messages and help text (visible even if standard output is redirected).

---

## 🧑‍💻 Installation

```sh
go install github.com/boxy-pug/ccwc@latest
```

Or clone and build manually:

```sh
git clone https://github.com/boxy-pug/ccwc.git
cd ccwc
go build -o ccwc
```

---

