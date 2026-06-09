# go-coreutils

> Go implementations of Unix core utilities — built as a learning project from
> [CodingChallenges.fyi](https://codingchallenges.fyi/).

Each tool is a hand-written, dependency-minimal clone of a classic Unix command,
focused on understanding both the original tool's behavior and idiomatic Go patterns.

## Tools

| Tool | Unix | Description | Key Flags | Challenge |
|------|------|-------------|-----------|-----------|
| [`cccat`](cmd/cccat/) | `cat` | Concatenate and print files | `-n`, `-b` | [#1](https://codingchallenges.fyi/challenges/challenge-cat) |
| [`ccgrep`](cmd/ccgrep/) | `grep` | Search files for patterns | `-r`, `-v`, `-i` | [#2](https://codingchallenges.fyi/challenges/challenge-grep) |
| [`cchead`](cmd/cchead/) | `head` | Output the first part of files | `-n`, `-c` | [#3](https://codingchallenges.fyi/challenges/challenge-head) |
| [`ccsed`](cmd/ccsed/) | `sed` | Stream editor for filtering/transforming text | `-n`, `-i` | [#4](https://codingchallenges.fyi/challenges/challenge-sed) |
| [`ccsort`](cmd/ccsort/) | `sort` | Sort lines using multiple algorithms | `--algo`, `--test`, `-u`, `-v` | [#5](https://codingchallenges.fyi/challenges/challenge-sort) |
| [`cctail`](cmd/cctail/) | `tail` | Output the last part of files | `-n` | [#6](https://codingchallenges.fyi/challenges/challenge-tail) |
| [`cctr`](cmd/cctr/) | `tr` | Translate or delete characters | `-d`, `-s` | [#7](https://codingchallenges.fyi/challenges/challenge-tr) |
| [`ccuniq`](cmd/ccuniq/) | `uniq` | Report or filter repeated lines | `-c`, `-d`, `-u` | [#8](https://codingchallenges.fyi/challenges/challenge-uniq) |
| [`ccwc`](cmd/ccwc/) | `wc` | Count lines, words, bytes, and characters | `-c`, `-l`, `-w`, `-m` | [#9](https://codingchallenges.fyi/challenges/challenge-wc) |
| [`ccxxd`](cmd/ccxxd/) | `xxd` | Hex dump and reverse hex dump | `-e`, `-r`, `-g`, `-c`, `-l`, `-s` | [#10](https://codingchallenges.fyi/challenges/challenge-xxd) |

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

- **Hand-written** — every line of Go is mine, written to learn, not generated.
- **Stdlib-first** — only one external dependency (`go-cmp` for ccsort tests).
- **Single-binary tools** — compile with `go build`, no frameworks, no installers.
- **Honest learning** — each tool's README has a "What I Learned" section documenting real debugging experiences, gotchas, and insights.

## About

This monorepo consolidates what were originally 10 separate repositories into one
discoverable project. [CodingChallenges.fyi](https://codingchallenges.fyi/challenges/)
provides the challenge prompts; I provide the implementations — warts, lessons,
and all.

## License

MIT — see [LICENSE](LICENSE).