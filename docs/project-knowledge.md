# Mimir — Project Knowledge Summary

*Source of truth. Finalized decisions only. Items marked `[UNDECIDED]` are explicitly open.*

---

## What Mimir Is

- A **suite of small, focused defensive cybersecurity tools** sharing a name, visual identity, and philosophy.
- **Through-line:** "Tools that give defenders wisdom" (Norse Mimir / Odin's eye framing).
- **Structure:** Single GitHub monorepo named `mimir`. Each tool is a subdirectory under `tools/<name>/` with its own README and demo.
- **Audience:** Recruiters, security-minded engineers, and the student-author's portfolio.
- **Primary language:** Go (1.22+). Other languages allowed only when they clearly fit better.

## What Mimir Is Not

- Not a framework. No "Mimir Core" until 3+ tools have shipped.
- Not a SaaS, startup, SOC platform, or commercial product replacement.
- Not academic research. Engineering polish over novelty.
- Not a dumping ground for unrelated security scripts.
- Not a single monolithic application. The suite is the *collection* of standalone tools.

## Conventions (locked in)

- **CLI naming:** `mimir-<name>` (e.g., `mimir-headers`, `mimir-leaks`).
- **Repo path per tool:** `tools/<name>/`.
- **Each tool ships with:** own README, install instructions, usage example, demo command, screenshot or asciinema recording.
- **Top-level README:** Mythology framing paragraph + table of `tool | what it sees | status | demo link`.
- **Shared visual identity:** Same logo, color palette, README structure across tools. `[UNDECIDED]` exact palette and logo.
- **Dependencies:** stdlib-first. New dependencies require a one-line justification in the spec.
- **Configuration:** Flags only. No config files unless a tool's spec explicitly requires one.

## Tooling Workflow

| Tool | Used for |
|---|---|
| Claude Chat | Brainstorming, specs, READMEs, handoff prompts, diff review, manual updates |
| Claude Projects | One project named "Mimir" — holds manual, conventions, per-tool specs, decisions log |
| Artifacts | Persistent edited documents (specs, READMEs, manuals). Naming: `mimir-<thing>` |
| Claude Code | All implementation: code, tests, CI, refactors. Driven by handoff prompts written in Chat |
| Claude Cowork | Not used unless a slide deck or large tracker is needed later |

## Chat Discipline

**Stay in the same chat when:** same tool, same phase, same document; conversation is still coherent and useful.

**Start a new chat when:** switching tools, switching modes (planning → review → marketing), Claude is forgetting earlier context, about to write a Claude Code handoff prompt, or current chat went off the rails.

When starting a new chat, paste in or reference relevant pinned project docs.

## Handoff Prompt Template (Chat → Claude Code)

Every handoff prompt must contain: Goal, Context, Constraints, Inputs/Outputs, Acceptance Criteria, Out of Scope, Files it may touch.

```
GOAL: <one sentence>

CONTEXT:
- Tool: mimir-<name>
- Repo path: tools/<name>/
- Existing files: <list>
- Relevant spec: <link or paste>

CONSTRAINTS:
- Language: Go 1.22+
- Dependencies: stdlib only unless listed
- Follow Mimir conventions

INPUTS / OUTPUTS:
- CLI: mimir-<name> [flags] <args>
- Flags: ...
- Output format: ...

ACCEPTANCE:
- [ ] go build ./... succeeds
- [ ] go test ./... passes
- [ ] README updated with usage example
- [ ] Demo command in spec produces expected output

OUT OF SCOPE:
- <explicit list>
```

If the template can't be filled in, the spec is not ready and no handoff happens.

## Definition of Done — Planning

- [ ] Written artifact exists (spec, decision, or manual update)
- [ ] Acceptance criteria for implementation are listed
- [ ] Out-of-scope items are listed explicitly
- [ ] A handoff prompt can be written from it without further brainstorming
- [ ] A future reader could understand the intent cold

## Definition of Done — Coding (per tool)

- [ ] Code compiles, no warnings
- [ ] Tests exist for core logic and pass
- [ ] CI runs tests and passes
- [ ] README has: description, install, usage example, demo command, screenshot/asciinema
- [ ] Demo command runs end-to-end on a clean machine or Docker container
- [ ] Tool has been used at least once on a real input
- [ ] Top-level Mimir README links to it with status updated
- [ ] Git tag pushed (e.g., `headers-v0.1.0`)

## Scope Creep Rules

- Out-of-scope list lives in each tool's spec and is updated whenever an idea is rejected.
- Per-tool LOC budget set in the spec. Going over triggers a stop-and-think.
- Two-week soft cap per tool phase. Past it, ship what exists.
- No new dependencies without one-line justification.
- No "while I'm here" refactors. New work needs a new branch and new spec.
- **Mythology test:** If you can't say in one sentence what kind of "wisdom" the tool gives defenders, it doesn't ship under Mimir.
- The Hub and the Core framework are deferred until 3+ tools exist.

## Tool N+1 Gate

Tool N+1 starts only when **all** of these are true for Tool N:

- [ ] Definition of Done — Coding is fully met
- [ ] Tool N has been used at least once on a real input
- [ ] Demo (asciinema/GIF/screenshot/video) is linked in its README
- [ ] Linked from top-level README with current status
- [ ] One-paragraph retrospective written in `decisions.md`

No exceptions. "Almost done" tools do not unlock the next one.

## Operating Loop (per milestone)

1. **Plan** — Open/update spec artifact. Goal, scope, out-of-scope, acceptance. Pass mythology test.
2. **Hand Off** — Write handoff prompt in Chat. Start fresh Claude Code session. Paste prompt.
3. **Build** — Claude Code implements. Review diffs as they come. Stop after two off-rails cycles and rewrite the prompt.
4. **Verify** — Walk Definition of Done. Run demo on clean checkout. Use the tool on something real.
5. **Ship** — Tag release. Update top-level README. Record demo. Write retrospective line.
6. **Stop** — Close the tool's chat. Take at least a day before starting Tool N+1.

## Rollout Order

1. **Mimir Headers** — warm-up; sets visual template, README format, CI pipeline.
2. **Mimir Leaks** — first real Go project; concurrency, API integration.
3. **Mimir Recon** — flagship; demo-video tool.
4. **Mimir Threat Feed** — first AI-assisted tool.
5. **Mimir Memory** OR **Mimir Prompt Shield** — `[UNDECIDED]` showpiece pick.

Tools beyond #5 are not committed to.

## Open Items

- `[UNDECIDED]` Exact color palette, logo, and typography for shared visual identity.
- `[UNDECIDED]` Showpiece choice between Mimir Memory and Mimir Prompt Shield.
- `[UNDECIDED]` Whether to use GitHub Actions or another CI provider (defaulting to GitHub Actions until reason to change).