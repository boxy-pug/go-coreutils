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

## Current behavior vs real `tr` spec

### What works now

- **string1 (target):** all classes expand to **literal ASCII sets**
  - `[:lower:]` → `abcdefghijklmnopqrstuvwxyz`
  - `[:upper:]` → `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
  - `[:alpha:]`, `[:digit:]`, `[:print:]`, `[:punct:]`, `[:space:]` → their ASCII chars
  - Range expansion: `a-c` → `abc`, `A-Z` → `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
- **string2 (translation):** also treated as literal sets (same expansion rules)
  - This means `tr '[:lower:]' '[:upper:]'` happens to work because the sets are the same size
  - But it's not doing real case conversion — it's mapping `a→A`, `b→B`, etc. by position
- **Delete (`-d`):** one string, all classes are literal sets, `-1` sentinel deletes

### What real `tr` does (and what we need to match)

**Translation (two strings, no `-d`):**
- string1: all classes are literal sets (same as us)
- string2: **only `[:lower:]` and `[:upper:]` are allowed**
  - These are **case conversion functions**, not literal sets
  - They must appear in the same relative position as their counterpart in string1
  - `tr '[:lower:]' '[:upper:]'` → lowercase chars become uppercase via `unicode.ToUpper`
  - `tr '[:upper:]' '[:lower:]'` → uppercase chars become lowercase via `unicode.ToLower`
  - `tr '[:lower:][:upper:]' '[:upper:][:lower:]'` → swap case
  - Any other class in string2 (`[:digit:]`, `[:alpha:]`, etc.) → **error**
  - Literal target with `[:upper:]`/`[:lower:]` in string2 → **error**
  - e.g., `tr 'abc' '[:upper:]'` → "misaligned" error

**Delete (`-d`):**
- One string only. All classes are literal sets. (same as us)

**Squeeze (`-s`):**
- One string only. Replaces N consecutive identical chars with 1.
- `tr -s 'abc'` → squeeze a, b, or c
- `tr -s 'abc' 'xyz'` → squeeze a→x, b→y, c→z

**Delete + Squeeze (`-d -s`):**
- Two strings: string1 = delete, string2 = squeeze
- Any class can be in string2 (since string2 is for squeezing, not translating)
- `tr -d -s '[:lower:]' '[:upper:]'` → delete lowercase, then squeeze uppercase
- All classes are ASCII-only in the C locale

## Todo

- [ ] `[:upper:]` and `[:lower:]` in string2 should be **case conversion functions**, not literal sets
  - `tr '[:lower:]' '[:upper:]'` → lowercase chars become uppercase via `unicode.ToUpper`
  - `tr '[:upper:]' '[:lower:]'` → uppercase chars become lowercase via `unicode.ToLower`
  - Currently they expand to literal ASCII sets (same as string1)
- [ ] Add misaligned error: other classes in string2 (`[:digit:]`, `[:alpha:]`, etc.) → error
- [ ] Add misaligned error: literal target with `[:upper:]`/`[:lower:]` in string2 → error
- [ ] `[:upper:]`/`[:lower:]` in string2 must match same relative position as counterpart in string1
  - `tr '[:lower:][:upper:]' '[:upper:][:lower:]'` → swap case (valid)
  - `tr '[:upper:]' '[:lower:]'` → valid
  - `tr 'abc' '[:upper:]'` → "misaligned" error
- [ ] Squeeze (`-s`) — one string only, replaces N consecutive identical chars with 1
- [ ] Delete + Squeeze (`-d -s`) — two strings: string1 = delete, string2 = squeeze
- [ ] Track which positions are class specifiers vs literal chars so buildTranslator knows which to treat as functions
