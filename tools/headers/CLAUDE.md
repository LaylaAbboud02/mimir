# mimir-headers — Tool briefing

See `/CLAUDE.md` for suite-wide rules (architecture, language, conventions).
This file covers only what is specific to mimir-headers.

## What it does

Fetches a single URL, evaluates the HTTP response headers against a static
catalog of security header rules, and emits a report. Shows defenders which
protective HTTP response headers their site is missing or misconfigured.
Intended for a developer doing a pre-deploy check.

**Voice:** advisor, not scanner. Each finding includes an explanation of the
risk and a concrete remediation. The output should feel like advice from a
knowledgeable colleague, not a raw checklist.

**Primary audience:** developer running a pre-deploy check in a terminal.
JSON output supports secondary consumers (scripts, pipelines) but does not
drive design decisions. CI integration is explicitly out of scope for v0.1.0.

## Architecture

```
Engine.Evaluate(Input) Result      pure function, no I/O
Input  { URL string; Header http.Header; ... }
Result { Findings []Finding; Summary Summary }
```

The engine never panics and has no error path. Operational errors (network,
DNS, TLS, bad flags) live entirely in the CLI layer.

The CLI layer owns: flag parsing, HTTP fetch, redirect policy, timeouts,
user-agent, error mapping, output formatting, exit codes, and final report
assembly.

## Package layout

```
cmd/mimir-headers/main.go
internal/cli/{cli.go,flags.go,fetch.go,report.go,output_text.go,output_json.go,exitcode.go}
internal/engine/{engine.go,catalog.go,finding.go,summary.go}
schema/findings.schema.json
demo/
docs/spec.md
```

## Catalog model

One catalog entry = one bad state = one Finding ID. Each `Check` function
handles its own preconditions (header absent, header present but wrong value,
etc.). Mutually exclusive rules are resolved by the rules themselves, not by
engine meta-logic.

```go
type CatalogEntry struct {
    ID          string
    Header      string
    Status      string
    Severity    Severity
    Message     string
    Remediation string
    Check       func(http.Header) (matched bool, observed *string)
}
```

## Finding model

`ID` is a stable kebab-case primary key. IDs are never reused. Renaming an
ID is a major version change. `Summary` severity counts always include all
four levels (`critical`, `high`, `medium`, `info`) even when the count is zero.

## Scope

See `docs/spec.md` for the full v0.1.0 spec. Summary:

**In scope:** single URL fetch, HTTPS default with `--allow-http` opt-in,
bounded redirects, configurable timeout, custom user-agent, static catalog,
text and JSON output, `--min-severity` filter, deterministic exit codes,
stdlib only.

**Out of scope:** crawling, multi-target, auth, cookies, proxies, custom
request headers, TLS inspection, cookie attribute analysis, HTML meta
parsing, config files, plugins, TUI, web dashboard, AI integration, markdown
report, mimir-core extraction.

## Current state

Engine scaffold complete: `x-powered-by-disclosure` catalog entry and tests.
CLI layer not started. Next task: HSTS catalog entries.

Severity rubric is not yet finalized. New catalog entries use placeholder
severities pending a dedicated review pass.
