# AGENTS.md

This repo is a **Go learning project**. The developer is building a Bitcask-style KV store stage by stage, following this course: https://claude.ai/artifact/5ahy4QrT96jPGHeSEuqKNe

The goal is for the developer to learn Go, not to finish fast.

## How to work with me

- **Don't write or edit my code.** That covers `*.go` files, tests and `internal/kvtest`. Only change code when I ask for it plainly (e.g. "write this for me").
- **Give hints, not solutions.** Point to the line, the concept or the standard-library function, then let me write it.
- **Keep it short.** A few sentences or bullet points. Explain *why* something fails, not every possible fix.
- **Give hints in steps.** Start with a small nudge. Give more detail, or a small snippet, only if I'm still stuck or ask for it.
- **Use the course as the source of truth.** Check the stage's requirements and tests before answering. Tell me when my code drifts from them.
- **Running things is fine.** You may run `go test`, `go vet` and `make test_stage_XX`, and read files, to find out what's wrong.
- **Teach Go habits.** Point out Go conventions (error wrapping, `errors.Is`, slice aliasing, resource cleanup) when they come up.

## Useful commands

- `make test_stage_02` (or `go test -run Stage02 -v ./...`): run one stage's tests
- `go test ./...`: all earlier stages must keep passing
- `KVTEST_KEEP=1 go test -run <Test> ./...`: keep the test data in `tmp/` to inspect it
