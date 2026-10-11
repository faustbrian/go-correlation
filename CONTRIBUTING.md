# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's documentation.
2. Run ordinary Go tests for the affected package before editing.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Selected mutation gates must finish with zero surviving viable mutants.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and CI; both modules retain their Go 1.27.0
minimum. An ordinary checkout is sufficient for root package tests and a focused
active fuzz smoke, without sibling repositories or goal ledgers:

```bash
GOWORK=off go test ./...
GOWORK=off go test . -run '^$' -fuzz '^FuzzParseCorrelationID$' -fuzztime=2s
```

Select the native target and package affected by your change; ordinary Go tests
already exercise deterministic fuzz seeds. For shared checks, install the public
`golib` v2.1.1 release declared in `.golib.yaml` and verify its checksums. CI uses
that same verifier, not a sibling source build:

```bash
make check
make ci
```

Both Make routes use bounded local checks, including one count-1, two-minute
root race selection: `TestDefaultFactoryAndContextAreConcurrentAndUnique`
asserts unique generated IDs and exact context values across concurrent callers
sharing the default buffered UUID factory. Neither module declares services or
test tags. Optional mutation analysis can select one affected package:

```bash
golib mutation --module . --package ./http
```

Use `make check-full` for explicit full assurance with configured fuzz budgets,
API compatibility, full risk gates and the published-dependency sibling
integration harness. For candidate-source composition, the intentional workspace
selects the local root; `GOWORK=off` in the harness instead exercises its published
root dependency and does not prove candidate composition. Select affected
transport and direct consumer tests for public-contract changes. Scheduled and
release qualification retain their full selected gates and the explicit
`release_dry_run` workflow dispatch. Report unavailable or failing checks without
claiming release readiness. CI/configuration-only changes do not require a
library release.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
