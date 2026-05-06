# Mimir Decisions Log

One paragraph per shipped tool, plus cross-cutting decisions worth remembering.

## Cross-cutting

### 2026-05-07 — Engine-first architecture

All Mimir tools follow a functional core, imperative shell pattern.
Engine package is a pure function: input struct in, result struct out, no I/O.
CLI package owns all I/O, formatting, flag parsing, and exit codes.
Reasoning: testability, future TUI/report modes, eventual Mimir Core extraction
once 3+ tools justify it.

## Tools

(empty until Headers ships)