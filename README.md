# go-coreutils

> Go implementations of Unix core utilities — built as a learning project from
> [CodingChallenges.fyi](https://codingchallenges.fyi/).

Each tool is a hand-written clone of a classic Unix command,
focused on understanding the original tool's behavior and Go coding.

## Tools

| Tool                    | Unix   | Description                                   | Challenge                                                                                    |
| ----------------------- | ------ | --------------------------------------------- | -------------------------------------------------------------------------------------------- |
| [`cccat`](cmd/cccat/)   | `cat`  | Concatenate and print files                   | [Build Your Own cat Tool](https://codingchallenges.fyi/challenges/challenge-cat)             |
| [`ccgrep`](cmd/ccgrep/) | `grep` | Search files for patterns                     | [Build Your Own grep](https://codingchallenges.fyi/challenges/challenge-grep)                |
| [`cchead`](cmd/cchead/) | `head` | Output the first part of files                | [Build Your Own head](https://codingchallenges.fyi/challenges/challenge-head)                |
| [`ccsed`](cmd/ccsed/)   | `sed`  | Stream editor for filtering/transforming text | [Build Your Own Sed](https://codingchallenges.fyi/challenges/challenge-sed)                  |
| [`ccsort`](cmd/ccsort/) | `sort` | Sort lines using multiple algorithms          | [Build Your Own Sort Tool](https://codingchallenges.fyi/challenges/challenge-sort)           |
| [`cctail`](cmd/cctail/) | `tail` | Output the last part of files                 | [Build Your Own tail](https://codingchallenges.fyi/challenges/challenge-tail)                |
| [`cctr`](cmd/cctr/)     | `tr`   | Translate or delete characters                | [Build Your Own tr Tool](https://codingchallenges.fyi/challenges/challenge-tr)               |
| [`ccuniq`](cmd/ccuniq/) | `uniq` | Report or filter repeated lines               | [Build Your Own uniq Tool](https://codingchallenges.fyi/challenges/challenge-uniq)           |
| [`ccwc`](cmd/ccwc/)     | `wc`   | Count lines, words, bytes, and characters     | [Build Your Own wc Tool](https://codingchallenges.fyi/challenges/challenge-wc)               |
| [`ccxxd`](cmd/ccxxd/)   | `xxd`  | Hex dump and reverse hex dump                 | [Build Your Own Xxd](https://codingchallenges.fyi/challenges/challenge-xxd)                  |

## Quick Start

```bash
# Clone the repo
git clone https://github.com/boxy-pug/go-coreutils.git
cd go-coreutils

# Build any tool
go build ./cmd/ccwc/

# Or build them all
go build ./cmd/...

# Run tests
go test ./...                     # all unit tests
go test -tags=integration ./...   # include integration tests
```

## Install via `go install`

```bash
go install github.com/boxy-pug/go-coreutils/cmd/ccwc@latest
go install github.com/boxy-pug/go-coreutils/cmd/ccxxd@latest
go install github.com/boxy-pug/go-coreutils/cmd/cchead@latest
# ... etc
```

## Project Philosophy

- Hand written and not perfect. Minimal ai use, for helping with explaining concepts, pointing to good std lib funcs, helping with tests, not generating code.
- Each tool's README has a "What i learned" section documenting real debugging experiences, gotchas, and insights.

## About

This monorepo consolidates what were originally 10 separate repositories into one
project.

## License

MIT — see [LICENSE](LICENSE).

