# Mimir — Suite-wide briefing

Mimir is a suite of focused, defensive security CLI tools. Each tool gives
defenders specific, actionable wisdom about one well-scoped concern. The name
comes from Norse mythology: Mimir is the keeper of wisdom at the well beneath
Yggdrasil. These tools give defenders wisdom.

## Language and runtime

Go 1.22+. Stdlib only. No third-party dependencies without a written decision
entry in `decisions.md`.

## Repo layout

Monorepo. Each tool is self-contained under `tools/<name>/` with its own
`go.mod`. No shared library yet — `mimir-core` is deferred until at least
three tools ship and a genuine shared abstraction is evident.

```
tools/
  headers/          # mimir-headers (in progress)
docs/
  project-knowledge.md
decisions.md
```

Directory names carry no `mimir-` prefix. Binary names keep the `mimir-`
prefix (e.g. `mimir-headers`).

## Hard rules

**Functional core / imperative shell.** Every tool's engine is a pure
function: no I/O, no global state, no side effects. Engine forbidden imports:
`os`, `flag`, `log`, `fmt.Print*`, `net.Dial`, and anything that performs
I/O. `net/http` types are allowed as data. The CLI layer owns all I/O.

**Stdlib only.** No third-party deps without a written justification in
`decisions.md`.

**Flags only.** No config files, no environment variables driving behavior.
All runtime configuration comes from CLI flags.

**No AI in tools.** The only tool permitted to use AI is `mimir-threat-feed`
(not yet started). All other tools are purely deterministic.

**CLI-first.** Every tool ships a CLI binary named `mimir-<name>`. Library
use is not a design goal.

**Structured output required.** Every tool must support both human-readable
text output and machine-readable JSON output.

## Go conventions

- Package names: lowercase, single word.
- Acronyms: `ID`, `URL`, `HTTP` — not `Id`, `Url`, `Http`.
- Errors: wrapped with `%w`, never swallowed.
- All exported identifiers have doc comments.
- No `init()` functions.
- No `interface{}` / `any` where a concrete type is known.
- Tests live next to the code they test. Use table-driven tests.
- Module path pattern: `github.com/LaylaAbboud02/mimir/tools/<name>`.

## Current state

- `mimir-headers`: engine scaffold in place, CLI not started.
- All other tools: not started.

## Docs

| File | Purpose |
|------|---------|
| `docs/project-knowledge.md` | Suite-wide background, goals, non-goals |
| `tools/<name>/CLAUDE.md` | Tool-specific briefing |
| `tools/<name>/docs/spec.md` | Tool MVP spec |
| `decisions.md` | Architecture and dependency decisions |

## Working discipline

One task per session. For non-trivial work, end the session with a handoff
prompt. Read existing code before adding new code. Do not refactor code
unrelated to the current task.
