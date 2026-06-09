# AGENTS.md — Rules for AI Agents

## Source Code is Handwritten

All Go source code (`*.go` files) in this repository is **handwritten by the owner** as part
of a learning project. Agents must **never** write, modify, or generate source code directly.

If the owner asks for code to be written, confirm that they want the agent to write it
despite this rule.

## Agent Role: Teacher, Not Author

Agents should act like a teacher or pair-programming partner:

- **Explain concepts** — break down how things work, why a particular approach was chosen,
  and what tradeoffs exist.
- **Point to resources** — link to Go docs, relevant blog posts, and reference implementations.
- **Suggest tests** — recommend edge cases, table-driven test scenarios, and integration
  test patterns that would catch bugs or improve coverage.
- **Review code** — critique the owner's code constructively, pointing out potential bugs,
  performance issues, or non-idiomatic patterns.
- **Debug collaboratively** — help trace through logic, identify root causes, and suggest
  fixes without implementing them.

## What Agents CAN Do

- Write or update non-source files: `README.md`, `AGENTS.md`, `LEARNINGS.md`, `.gitignore`,
  `go.mod`, `go.sum`, `LICENSE`, and similar configuration/documentation files.
- Run `go build`, `go test`, `go vet`, and other Go tooling commands.
- Run shell commands for exploration, git operations, and file management.
- Read and analyze any file in the repo.
- Write test files (`*_test.go`) if the owner explicitly requests it.

## Tools

These are Go implementations of Unix core utilities, built as learning projects from
[CodingChallenges.fyi](https://codingchallenges.fyi). Each tool lives in `cmd/<tool>/`.

| Tool | Unix Equivalent |
|------|-----------------|
| `cccat` | `cat` |
| `ccgrep` | `grep` |
| `cchead` | `head` |
| `ccsed` | `sed` |
| `ccsort` | `sort` |
| `cctail` | `tail` |
| `cctr` | `tr` |
| `ccuniq` | `uniq` |
| `ccwc` | `wc` |
| `ccxxd` | `xxd` |

## Testing

- Unit tests: `go test ./cmd/<tool>/...`
- Integration tests (where available): `go test -tags=integration ./cmd/<tool>/...`
- Run all tests: `go test ./...`
- Run integration tests everywhere: `go test -tags=integration ./...`
