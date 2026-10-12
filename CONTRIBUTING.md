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

Do not fork shared verification workflows, add permanent replacements, machine-specific
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
`golib` v2.4.0 release declared in `.golib.yaml` and verify its checksums. CI uses
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
transport and direct consumer tests for public-contract changes. Explicit full and
release qualification retain their selected gates and the explicit
`release_dry_run` workflow dispatch. Report unavailable or failing checks without
claiming release readiness. CI/configuration-only changes do not require a
library release.

## Selected Monthly Assurance

The owner-level [selected-assurance caller](.github/workflows/assurance.yml) is
configured for main only, on day 24 of each month at 05:07 UTC. It selects one
of the 4 current declared root-module fuzz targets using the shared UTC-week
rotation, with 30 seconds of active exploration. The selected risks are:
correlation identifier identity, carrier extraction/reinjection, inbound
HTTP headers and custom JSON-RPC metadata.
The HTTP target uses in-process request/response fixtures, not a deployed
service. The separate sibling integration harness retains candidate-workspace
and public-dependency composition boundaries; selected fuzzing is not their
qualification.

Base-free ordinary checks execute deterministic seeds, not active exploration;
committed-range CI selects affected active smoke. Manual assurance requires an
exact anchored selector from the root fuzz operation in `.golib.yaml`, for
example `^FuzzCarrierExtraction$`.
Empty, unknown and malformed selections refuse. The shared job bounds setup
and execution to ten minutes, active workers and CPU to two, owned memory to
2 GiB, native output to 1 MiB and private-cache lifetime. Failure, cancellation,
exhausted budgets and cleanup refusal remain attributable failures. This lane
does not run mutation or replace Required, original full fuzz budgets,
explicit `make check-full`, or relevant release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.

## Source Documentation

Install the checksum-pinned public verifier declared in `.golib.yaml`.
`make source-docs-check` validates objective Go comment, marker, link and
generated-producer structure; `make check`, `make check-full` and `make ci`
require it before their existing checks. Shared CI runs the same strict
contract and includes its result in Required. This is structural enforcement,
not a substitute for source-truth review or behavioral verification.
