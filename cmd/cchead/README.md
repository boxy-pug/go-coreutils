# 📚 cchead

cchead is a simple Go implementation of the Unix `head` command, created as part of a coding challenge on [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-head).

---

## ✅ Features

-  Print a specified number of lines or bytes from the start of files.
-  Supports multiple files, displaying headers for each.

---

## 🛠️ Usage

-  `-n <number>`: Number of lines to print (default: 10).
-  `-c <number>`: Number of bytes to print.

```sh
# Print the first 10 lines of `./testdata/test.txt`:
  ./cchead -n 10 ./testdata/test.txt

# Print the first 100 bytes of `./testdata/test.txt`:
  ./cchead -c 100 ./testdata/test.txt

# Print the first 5 lines of `./testdata/test.txt` and `./testdata/test2.txt`:
  ./cchead -n 5 ./testdata/test.txt ./testdata/test2.txt

```

---

## 🧐 What I learned

-  **`bufio.Scanner` vs. `bufio.Reader`:** `Scanner` is the simpler version that doesn't retain newlines, which might cause confusion with Windows-style carriage returns `/r/n` and Mac-like newlines `/n`. `Reader` is more advanced, allowing you to read data until a specific delimiter, but you also have to treat the EOF and error cases more carefully.

-  For integration testing CLI tools, you can use the `exec package` to compare the actual output of your command to the output of another command – very handy when making a clone.

---

## 🧑‍💻 Installation

```sh
go install github.com/boxy-pug/cchead@latest
```

Or clone and build manually:

```sh
git clone https://github.com/boxy-pug/cchead.git
cd cchead
go build -o cchead
```

---

