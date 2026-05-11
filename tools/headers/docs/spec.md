# Mimir Headers — MVP Spec

**Tool:** `mimir-headers`
**Repo path:** `tools/headers/`
**Status:** Spec — pre-implementation
**Owner:** Mimir suite
**Mythology line:** *Shows defenders which protective HTTP response headers their site is missing or misconfigured.*

---

## 1. Product definition

`mimir-headers` is a CLI tool that fetches a URL, evaluates its HTTP security response headers against a documented checklist, and emits a prioritized human- or machine-readable report.

## 2. Target user

A defensive engineer, security-minded developer, or pentester who wants a fast, scriptable, locally-run check of a single web target's security header posture — without uploading the URL to a third-party service and without operating a full scanning platform.

Concretely:

- **Primary:** application security engineers and developers running a quick pre-deploy or pre-PR check against a staging or production URL.
- **Secondary:** pentesters, bug bounty hunters, and CTF players who need scriptable header analysis they can pipe into other tools.
- **Tertiary:** instructors and students using a transparent, open-source reference implementation.

## 3. Problem it solves

HTTP security headers (CSP, HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy, etc.) are one of the cheapest and highest-leverage defenses a web application has, and they are routinely missing or misconfigured. Existing options have drawbacks:

- **Hosted scanners** (securityheaders.com, Mozilla Observatory) require sending the target URL to a third party and are awkward to script or run against internal/staging hosts.
- **Manual `curl -I` inspection** is fast but requires the operator to remember the full checklist and severity model.
- **Larger scanners** (ZAP, nuclei templates, Burp) include header checks but are heavyweight for a focused, repeatable, single-URL question.

`mimir-headers` fills the gap: a single static binary, no network egress beyond the target, scriptable JSON output, no third-party dependency.

## 4. MVP features

| # | Feature | Notes |
|---|---|---|
| F1 | Single-URL fetch via `net/http` | HTTPS default; HTTP allowed only with explicit opt-in flag |
| F2 | Configurable timeout | Default 10s, settable via `--timeout` |
| F3 | Bounded redirect following | Default max 5, settable via `--max-redirects`. `0` disables |
| F4 | Custom User-Agent | Default `mimir-headers/<version>`, settable via `--user-agent` |
| F5 | Static checklist of security headers | See §10. Catalog is versioned in code, not user-extensible |
| F6 | Per-check verdict | Each check returns `{name, status, observed_value, severity, message, remediation}` |
| F7 | Human-readable text output | Default. Grouped by severity. Color via auto-detected TTY; disabled with `--no-color` |
| F8 | JSON output | `--format json`. Schema documented in README and stable within `0.x` minor versions |
| F9 | Deterministic exit codes | `0` clean, `1` findings present, `2` operational error |
| F10 | `--version` and `--help` | Standard CLI hygiene |

Each MVP feature passes the mythology test: every check exists to give the defender clearer wisdom about a specific protection their site does or does not have.

## 5. Out of scope (explicit)

These are rejected for v0.1.0. They live in this list, not in the backlog. Reopening requires a new spec entry.

- Crawling, sitemap walking, or evaluating any URL beyond the one provided.
- Multi-target / batch / stdin-list / file-list mode.
- Authenticated requests, cookies, request-body sending, custom request headers.
- Proxy support (`HTTP_PROXY` env var explicitly *not* honored in MVP).
- TLS / certificate inspection. Separate tool's job.
- Cookie attribute analysis (`Secure`, `HttpOnly`, `SameSite`). Future tool or v0.2 enhancement.
- Parsing of `<meta http-equiv>` security policies in HTML. MVP looks at response headers only.
- CSP semantic deep analysis (unsafe-inline detection beyond a binary flag, source-list reachability, etc.). MVP only checks for presence and a small set of well-known weakening directives.
- HTML report output. Markdown report is a candidate later enhancement, not MVP.
- Diff / comparison mode against a previous scan.
- Configuration files. Flags only, per Mimir conventions.
- Plugin / custom-rule system. The catalog is static and code-reviewed.
- TUI, web dashboard, or any non-CLI interface.
- Any AI / LLM integration.
- Extraction into a `mimir-core` shared package. Deferred until 3+ tools justify it.
- Final logo and color palette. Tracked as `[UNDECIDED]` at the project level. Headers ships with placeholder styling.

## 6. CLI interface design

### 6.1 Synopsis

```
mimir-headers [flags] <url>
```

### 6.2 Flags

| Flag | Type | Default | Purpose |
|---|---|---|---|
| `--format` | enum: `text`, `json` | `text` | Output format |
| `--timeout` | duration | `10s` | Per-request timeout (Go duration syntax) |
| `--max-redirects` | int | `5` | Cap on redirects followed; `0` disables |
| `--allow-http` | bool | `false` | Permit `http://` URLs. Without this, plain HTTP is rejected before fetch |
| `--user-agent` | string | `mimir-headers/<version>` | Override User-Agent |
| `--no-color` | bool | `false` | Disable ANSI color in text output |
| `--min-severity` | enum: `info`, `low`, `medium`, `high` | `info` | Suppress findings below this severity |
| `--version` | bool | — | Print version and exit |
| `--help`, `-h` | bool | — | Print usage and exit |

### 6.3 Argument rules

- Exactly one positional argument: the target URL.
- URL must include scheme. A bare host like `example.com` is rejected with a clear error suggesting `https://example.com`.
- Anything other than `http`/`https` is rejected.

### 6.4 Exit codes

| Code | Meaning |
|---|---|
| `0` | Fetch succeeded, no findings at or above `--min-severity` |
| `1` | Fetch succeeded, one or more findings at or above `--min-severity` |
| `2` | Operational error (network, DNS, timeout, TLS, invalid URL, redirect loop, non-HTTP scheme, etc.) |

Exit code is determined by the *highest-severity finding observed*, not by count.

## 7. Input / output examples

### 7.1 Text output, healthy site (illustrative)

```
$ mimir-headers https://example.com
mimir-headers v0.1.0  —  https://example.com  (200 OK, 312ms)

HIGH      (0)
MEDIUM    (1)
  - Content-Security-Policy        missing
    A CSP significantly reduces the impact of XSS. Add a policy starting from
    `default-src 'self'` and tighten from there.

LOW       (2)
  - Permissions-Policy             missing
  - Referrer-Policy                missing

INFO      (1)
  - Server                         present (value disclosed: "nginx/1.25.3")

Summary: 4 findings (0 high, 1 medium, 2 low, 1 info)
Exit: 1
```

### 7.2 JSON output

```
$ mimir-headers --format json https://example.com
```

See §8 for the full schema.

### 7.3 Operational error

```
$ mimir-headers https://does-not-resolve.invalid
error: dns lookup failed for "does-not-resolve.invalid": no such host
Exit: 2
```

```
$ mimir-headers http://example.com
error: refusing to fetch plain http URL; pass --allow-http to override
Exit: 2
```

```
$ mimir-headers example.com
error: url must include scheme (e.g. https://example.com)
Exit: 2
```

## 8. JSON output format

The JSON schema is the contract. It is versioned via `schema_version`. Breaking changes require a major bump and a changelog entry.

### 8.1 Top-level object

```json
{
  "schema_version": "1",
  "tool": "mimir-headers",
  "tool_version": "0.1.0",
  "target": {
    "requested_url": "https://example.com",
    "final_url": "https://example.com/",
    "redirect_count": 0,
    "status_code": 200,
    "elapsed_ms": 312,
    "fetched_at": "2026-05-05T14:22:31Z"
  },
  "summary": {
    "total": 4,
    "by_severity": { "high": 0, "medium": 1, "low": 2, "info": 1 }
  },
  "findings": [
    {
      "id": "csp-missing",
      "header": "Content-Security-Policy",
      "status": "missing",
      "observed_value": null,
      "severity": "medium",
      "message": "Content-Security-Policy header is not set.",
      "remediation": "Define a CSP starting from `default-src 'self'` and refine for your application's needs."
    }
  ],
  "exit_code": 1
}
```

### 8.2 Field definitions

| Field | Type | Notes |
|---|---|---|
| `schema_version` | string | Currently `"1"`. Stable across `0.x` minor releases |
| `tool` | string | Always `"mimir-headers"` |
| `tool_version` | string | SemVer of the binary |
| `target.requested_url` | string | The URL as supplied on the command line |
| `target.final_url` | string | URL after redirect chain |
| `target.redirect_count` | int | Number of redirects followed (0–`max_redirects`) |
| `target.status_code` | int | Final HTTP status |
| `target.elapsed_ms` | int | Wall clock from request start to response headers received |
| `target.fetched_at` | string | RFC 3339 UTC timestamp |
| `summary.total` | int | Total findings emitted (after `--min-severity` filter) |
| `summary.by_severity` | object | Counts per severity |
| `findings[]` | array | Ordered: high → medium → low → info, then alphabetical by `id` |
| `findings[].id` | string | Stable, kebab-case identifier (e.g. `csp-missing`, `hsts-short-max-age`) |
| `findings[].header` | string | Canonical header name |
| `findings[].status` | enum | `present`, `missing`, `malformed`, `weak` |
| `findings[].observed_value` | string \| null | Raw header value if present, else `null` |
| `findings[].severity` | enum | `high`, `medium`, `low`, `info` |
| `findings[].message` | string | One-line, human-readable |
| `findings[].remediation` | string | Short, actionable |
| `exit_code` | int | Mirrors process exit code |

### 8.3 Stability rules

- Adding a new finding `id` is a **minor** change.
- Adding a new field is a **minor** change.
- Renaming or removing a field, or changing a finding `id`'s severity, is a **major** change.
- The tool MUST emit valid JSON to stdout when `--format json` is set, even on operational errors. Errors take this shape:

```json
{
  "schema_version": "1",
  "tool": "mimir-headers",
  "tool_version": "0.1.0",
  "error": {
    "kind": "dns",
    "message": "dns lookup failed for \"does-not-resolve.invalid\": no such host"
  },
  "exit_code": 2
}
```

## 9. Edge cases

The engine must handle each of these explicitly. Each row maps to at least one test in §12.

| # | Case | Required behavior |
|---|---|---|
| E1 | Bare hostname (no scheme) | Reject before fetch, exit 2, message suggests `https://` |
| E2 | `http://` without `--allow-http` | Reject before fetch, exit 2 |
| E3 | Non-HTTP scheme (`ftp://`, `file://`) | Reject, exit 2 |
| E4 | DNS failure | Exit 2, `error.kind = "dns"` |
| E5 | Connection refused | Exit 2, `error.kind = "connection"` |
| E6 | TLS handshake failure | Exit 2, `error.kind = "tls"` |
| E7 | Timeout exceeded | Exit 2, `error.kind = "timeout"` |
| E8 | Redirect loop / exceeds `--max-redirects` | Exit 2, `error.kind = "redirect"` |
| E9 | Cross-protocol redirect (HTTPS → HTTP) without `--allow-http` | Treat as redirect failure; exit 2 |
| E10 | Non-2xx final status (e.g., 403, 500) | Still evaluate headers normally; status code reflected in output |
| E11 | Header appears multiple times | Concatenate per RFC 9110 (`,` join), evaluate combined value, note duplication in `observed_value` |
| E12 | Header value with non-ASCII or unprintable bytes | Preserve in JSON via valid UTF-8/escapes; in text output, replace control chars with `·` |
| E13 | Extremely long header value (e.g. 64 KB CSP) | No truncation in JSON; text output truncates to 200 chars with `…` and notes truncation |
| E14 | HSTS present but `max-age` < 15552000 (180 days) | `status: "weak"`, severity low |
| E15 | HSTS with `max-age=0` | `status: "weak"`, severity medium (actively disables HSTS) |
| E16 | CSP present but contains `unsafe-inline` for `script-src` or `default-src` | `status: "weak"`, severity medium, finding id `csp-unsafe-inline` |
| E17 | `X-Frame-Options` present alongside CSP `frame-ancestors` | Informational note; not a finding |
| E18 | Server header discloses detailed version | `severity: info`, finding id `server-version-disclosure` |
| E19 | Empty response body, only headers | Evaluate normally |
| E20 | Server returns HTTP/2 or HTTP/3 | Evaluate normally; protocol noted in `target` (later enhancement — for MVP, just don't break) |

## 10. Security header catalog (MVP)

Locked for v0.1.0. Each entry has at least one finding `id`. Severities are tool defaults and may be discussed in review but not silently changed.

| Header | Checks | Severities |
|---|---|---|
| `Strict-Transport-Security` | missing; `max-age` < 180d; `max-age=0`; missing `includeSubDomains` (info) | high / low / medium / info |
| `Content-Security-Policy` | missing; contains `unsafe-inline` in script/default; contains `unsafe-eval` | medium / medium / low |
| `X-Frame-Options` | missing AND no CSP `frame-ancestors`; value other than `DENY`/`SAMEORIGIN` | medium / low |
| `X-Content-Type-Options` | missing; value not `nosniff` | low / low |
| `Referrer-Policy` | missing; permissive value (`unsafe-url`, `no-referrer-when-downgrade`) | low / low |
| `Permissions-Policy` | missing | low |
| `Cross-Origin-Opener-Policy` | missing | info |
| `Cross-Origin-Resource-Policy` | missing | info |
| `Server` | discloses version | info |
| `X-Powered-By` | present at all | info |

If a header is present and not weak, the tool emits no finding for it. The text output's "clean" line lists which checked headers were OK so the operator can see what was evaluated.

## 11. Acceptance criteria

The build is done when **all** of the following are true.

- [ ] `go build ./...` succeeds with no warnings.
- [ ] `go test ./...` passes; engine package coverage ≥ 80%.
- [ ] Engine package has no import of `os`, `flag`, `fmt.Print*`, or any CLI/output concern.
- [ ] CLI binary `mimir-headers` is produced at `tools/headers/cmd/mimir-headers/`.
- [ ] All MVP features F1–F10 are implemented and demonstrated in tests.
- [ ] All edge cases E1–E20 are covered by at least one test (see §12).
- [ ] All catalog entries in §10 are implemented with their finding IDs.
- [ ] JSON output validates against a JSON Schema committed at `tools/headers/schema/findings.schema.json`.
- [ ] Exit codes match §6.4 across at least three integration tests (clean, findings, error).
- [ ] No third-party dependencies. `go.mod` lists stdlib only.
- [ ] CI (GitHub Actions) runs `go vet`, `go test`, and `go build` on push and PR; green on `main`.
- [ ] README at `tools/headers/README.md` covers all sections in §13.
- [ ] Top-level Mimir README updated with Headers row and current status.
- [ ] Demo command from §14 runs end-to-end against a real public site on a clean checkout.
- [ ] asciinema recording or screenshot linked in `tools/headers/README.md`.
- [ ] Git tag `headers-v0.1.0` pushed.
- [ ] One-paragraph retrospective added to `decisions.md`.

## 12. Test cases

Tests are table-driven where possible. The engine layer is tested with synthetic `http.Response` objects; CLI integration tests use `httptest.Server`.

### 12.1 Engine unit tests (`tools/headers/engine/`)

| ID | Test | Asserts |
|---|---|---|
| T-E-01 | All headers present and strong | Zero findings, exit code would be 0 |
| T-E-02 | All checked headers missing | Findings for each missing header at correct severity |
| T-E-03 | HSTS `max-age=0` | Finding `hsts-disabled`, severity medium |
| T-E-04 | HSTS `max-age=15551999` (just under 180d) | Finding `hsts-short-max-age`, severity low |
| T-E-05 | HSTS `max-age=31536000; includeSubDomains` | No finding |
| T-E-06 | CSP with `script-src 'self' 'unsafe-inline'` | Finding `csp-unsafe-inline`, severity medium |
| T-E-07 | CSP with `'unsafe-eval'` | Finding `csp-unsafe-eval`, severity low |
| T-E-08 | XFO missing but CSP has `frame-ancestors 'none'` | No XFO finding |
| T-E-09 | XFO `ALLOW-FROM https://x` | Finding `xfo-deprecated-value`, severity low |
| T-E-10 | XCTO present but value `nosniff, nosniff` (duplicated) | No finding; canonicalization handled |
| T-E-11 | Referrer-Policy `unsafe-url` | Finding `referrer-policy-permissive`, severity low |
| T-E-12 | Server `nginx/1.25.3` | Finding `server-version-disclosure`, severity info |
| T-E-13 | `X-Powered-By` present | Finding `x-powered-by-disclosure`, severity info |
| T-E-14 | Header appears twice (`X-Frame-Options: DENY` and `X-Frame-Options: SAMEORIGIN`) | Combined per RFC 9110; finding flags duplication |
| T-E-15 | 64 KB CSP value | No crash; full value preserved in JSON |
| T-E-16 | Header value with `\x00` and non-ASCII bytes | No crash; JSON valid UTF-8 |
| T-E-17 | `--min-severity=high` filter | Only high findings emitted |
| T-E-18 | Findings ordered high → info, then alphabetical by id | Order is stable across runs |

### 12.2 CLI / integration tests (`tools/headers/cmd/mimir-headers/`)

| ID | Test | Asserts |
|---|---|---|
| T-C-01 | `httptest.Server` with all-good headers | Exit 0, text output contains `0 high, 0 medium...` |
| T-C-02 | `httptest.Server` with one medium finding | Exit 1 |
| T-C-03 | DNS failure target | Exit 2, error JSON valid when `--format json` |
| T-C-04 | Timeout (server sleeps past `--timeout`) | Exit 2, `error.kind="timeout"` |
| T-C-05 | Redirect chain longer than `--max-redirects` | Exit 2, `error.kind="redirect"` |
| T-C-06 | HTTPS → HTTP redirect without `--allow-http` | Exit 2 |
| T-C-07 | Bare hostname argument | Exit 2 with scheme-suggestion message |
| T-C-08 | `ftp://` URL | Exit 2 |
| T-C-09 | `--format json` produces output that parses and validates against schema | Schema validation passes |
| T-C-10 | `--no-color` produces ANSI-free output | No `\x1b[` bytes in stdout |
| T-C-11 | `--version` exits 0 with version string | — |
| T-C-12 | `--help` exits 0 with usage | — |

### 12.3 Non-tests (deliberately omitted)

- No live network tests in CI. The demo command in §14 is run manually before tagging a release.
- No fuzzing in MVP. Listed as later enhancement.

## 13. README sections

`tools/headers/README.md` must contain, in order:

1. **Title and one-line description** with the mythology line.
2. **Status badge** (build, version).
3. **What it sees** — one paragraph on the threat model and what `mimir-headers` does and does not check.
4. **Install** — `go install` line and prebuilt binary note.
5. **Usage** — synopsis, flags table, exit codes table.
6. **Examples** — text output, JSON output, error case.
7. **Header catalog** — link or inline copy of §10.
8. **JSON schema** — link to `schema/findings.schema.json` and stability rules.
9. **Demo** — asciinema or screenshot.
10. **Limitations / out of scope** — short list with link to spec.
11. **License**.

## 14. Demo scenario

Single command, runnable on a clean machine with Go installed and the binary on `PATH`:

```
mimir-headers https://example.com
```

Then the same target with JSON, piped to `jq`:

```
mimir-headers --format json https://example.com | jq '.summary'
```

Recording requirements:

- Use asciinema (preferred) or a static screenshot.
- Show the text run *and* the JSON-piped-to-`jq` run in the same recording.
- Recording is committed under `tools/headers/demo/` and linked from the README.
- Recording is reproducible: the README documents the exact command and target so a reviewer can re-run.

The demo target for the recording is `https://example.com` because it is stable, public, low-risk, and produces a non-trivial finding set. Real-world demos against the author's own staging hosts are encouraged but not committed to the repo.

---

## Appendix A — Definition of Done checklist (mirrors Project Knowledge Summary)

- [ ] Code compiles, no warnings.
- [ ] Tests exist for core logic and pass.
- [ ] CI runs tests and passes.
- [ ] README has description, install, usage, demo command, screenshot/asciinema.
- [ ] Demo command runs end-to-end on a clean checkout.
- [ ] Tool used at least once on a real input.
- [ ] Top-level Mimir README links to it with status updated.
- [ ] Git tag `headers-v0.1.0` pushed.
- [ ] Retrospective paragraph in `decisions.md`.

## Appendix B — Files this spec authorizes

```
tools/headers/
  README.md
  cmd/mimir-headers/
    main.go
    main_test.go
  engine/
    engine.go
    engine_test.go
    catalog.go
    catalog_test.go
  schema/
    findings.schema.json
  demo/
    headers.cast        # asciinema, or
    headers.png         # screenshot
.github/workflows/
  headers.yml           # if not already covered by repo-wide CI
```

Anything outside this tree requires a spec amendment.