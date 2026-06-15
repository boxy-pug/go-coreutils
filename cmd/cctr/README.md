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
  - Function class specifier like "[:lower:]", matches all lowercase letters
- So there are 4 translation modes, translating to and from Regular (including range) and Function:
  - reg -> reg: build a map that maps targets to translations, easy
  - reg -> func: build a map for targets, when char in amp run the func on it
  - func -> reg: trickiest one. when matching the func (like "[:lower:]") then replace by translation chars, one by one. The last one repeats when your out of translation chars.
  - func -> func: if it matches target func (like "[:lower:]"), run it through translation func like "[:upper:]" for example, for translating all lowercase to uppercase
- I'm leaving this comments in for now, but actually I understood that there is a simpler way of thinking about this, and that's how tr works: The class specifiers just expand to a regular literal set of chars, they are not "functions" or special in any way, only diff is that its restricted what class specifiers you can put in string2/translation string.
- While refactoring I kind of understood that there's a simpler pattern i could use for these tools, to allow for better/easier testing. Basically: use main() as a thin wrapper around a run() function. So run() becomes the new main, kindof. In main you have an err := run(os.Args, io.Reader, io.Writer), and that makes for really simple integration testing. Then inside run you can have separate more or less pure functions that you can unit test, so like cfg = parseFlags(args[1:]), buildTranslator(cfg), translate(cfg, tr) for example. So then you dont have to hack around with os.Args and stuff in tests.
- So the basic insight: main is kindof an impure shell, that binds to the real world with os.args and in and out. when making a run function that can be a neat testable boundary. everything inside there can be pure funcs, so f ex: Instead of parsing os.Args, you're parsing a []string.

## tr spec, how it works

So the basic idea is: you give `tr` a set of characters to look for (string1/target) and a set of characters to replace them with (string2/translation). Then it walks through the input and goes "is this char in my target set? if so, replace it with the corresponding char from the translation set."

If the target is longer than the translation, the last char of the translation just repeats. So `tr 'abcd' 'xy'` gives you a→x, b→y, c→y, d→y.

**Expressions** — `tr` recognizes three kinds of things in its strings:

- **Literal strings** like `abc` — just those three chars
- **Ranges** like `a-z` — shorthand for all lowercase letters
- **Class specifiers** like `[:lower:]` — these expand to predefined sets of chars

I support these classifiers: `[:lower:]`, `[:upper:]`, `[:alpha:]`, `[:digit:]`, `[:print:]`, `[:punct:]`, `[:space:]`. They just expand to a set of chars, so writing `[:digit:]` is exactly the same as writing `0123456789`

`tr` is strict about what classes you can put in the second/translation string. Only `[:upper:]` and `[:lower:]` are allowed there, and they must match up with a corresponding `[:upper:]` or `[:lower:]` in the target string. Everything else gets rejected with a "misaligned" or "invalid class" error.

So `tr '[:lower:]' '[:upper:]'` is fine (lowercase → uppercase), but `tr '[:alpha:]' '[:upper:]'` is an error because `[:alpha:]` in string1 doesn't align with `[:upper:]` in string2. Same with `tr 'abc' '[:upper:]'` — literal strings in string1 can't pair with class specifiers in string2.

**Delete (`-d`)** — just one string. Any chars in that set get deleted.

## Differences from real `tr`

- My implementation is ascii only, i think real `tr` is locale aware and supports unicode for things like lower and upper.
- I support a single class per string, each arg can be one expression. In real `tr` you can combine stuff like `tr '[:lower:][:upper:]' '[:upper:][:lower:]'`.
- I only support squeeze with one string like `tr -s 'abc'` not translate then squeeze, like `tr -s 'abc' 'xyz'`.
- I don't provide separate error msg for misalignment and invalid class. Error messages should be:
  - **Invalid class in string2** (non-upper/lower class like `[:digit:]`, `[:alpha:]`, etc.): `when translating, the only character classes that may appear in string2 are 'upper' and 'lower'`
  - **Misaligned** (`[:upper:]`/`[:lower:]` in string2 but string1 is literal or different class): `misaligned [:upper:] and/or [:lower:] construct`
- I also don't support delete + squeeze. This is how this combo should work i think:
  - Two strings: string1 = delete, string2 = squeeze
  - Any class can be in string2 (since string2 is for squeezing, not translating)
  - `tr -d -s '[:lower:]' '[:upper:]'` → delete lowercase, then squeeze uppercase
