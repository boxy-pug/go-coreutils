# cctr

A Go clone of the Unix `tr` (translate) command, built as part of the [Build Your Own tr Tool](https://codingchallenges.fyi/challenges/challenge-tr) challenge.

cctr translates, deletes, and squeezes characters from standard input.

## Usage

```sh
# Translate characters
echo "hello" | cctr 'a-z' 'A-Z'

# Translate character classes
echo "hello" | cctr '[:lower:]' '[:upper:]'

# Delete characters
echo "hello world" | cctr -d 'aeiou'

# Squeeze repeated characters
echo "helloooo" | cctr -s 'o'
```

## Flags

- `-d` — delete characters in the given set from input
- `-s` — squeeze multiple occurrences of a character into one

## Install

```sh
go install github.com/boxy-pug/go-coreutils/cmd/cctr@latest
```

## What I learned

- `tr` is basically a char by char translator, like a find and replace that works on individual characters. Give it a target (chars to look for) and translation (what to turn them into). Then it walks through file and goes "is this a match, if so lets translate it".
- Three kinds of expressions:
  - Regular "abc": matches those three chars
  - Range "a-d": matches abcd (basically just syntactic sugar for regular)
  - Function class specifier like "[:lower:]", matches all lowercase letters (semantic, doesnt just expand to a list of chars)
- So there are 4 translation modes, translating to and from Regular (including range) and Function:
  - reg -> reg: build a map that maps targets to translations, easy
  - reg -> func: build a map for targets, when char in amp run the func on it
  - func -> reg: trickiest one. when matching the func (like "[:lower:]") then replace by translation chars, one by one. The last one repeats when your out of translation chars.
  - func -> func: if it matches target func (like "[:lower:]"), run it through translation func like "[:upper:]" for example, for translating all lowercase to uppercase
  - While refactoring I kind of understood that there's a simpler pattern i could use for these tools, to allow for better/easier testing. Basically: use main() as a thin wrapper around a run() function. So run() becomes the new main, kindof. In main you have an err := run(os.Args, io.Reader, io.Writer), and that makes for really simple integration testing. Then inside run you can have separate more or less pure functions that you can unit test, so like cfg = parseFlags(args[1:]), buildTranslator(cfg), translate(cfg, tr) for example. So then you dont have to hack around with os.Args and stuff in tests.
  - So the basic insight: main is kindof an impure shell, that binds to the real world with os.args and in and out. when making a run function that can be a neat testable boundary. everything inside there can be pure funcs, so f ex: Instead of parsing os.Args, you're parsing a []string.
