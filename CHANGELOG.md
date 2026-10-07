# Changelog

All notable changes follow Keep a Changelog. This project uses semantic
versioning once released.

## Unreleased

## 1.1.3 - 2026-10-08

### Changed

- Adopt OpenTelemetry API and trace v1.47.0 in the root and sibling harness,
  with metric and the newly selected log API at v1.47.0 in the harness.
  Keep its SDK components at v1.45.0; this update does not adopt later SDK
  fixes. The supplier Go 1.26 minimum remains below this module's Go 1.27
  requirement. Typed attributes and trace identity remain compatible;
  Stringer-based attribute diagnostics use the new JSON representation.

- Select Identifier v2.0.1 for default UUIDv4 generation, retaining canonical
  IDs and classified entropy failures. Align indirect pgx v5.11.0 and ULID
  v2.1.2 dependencies in the root and sibling harness while preserving the
  harness's historical Identifier v1 dependency. Applications implementing
  custom `pgx.Rows` must provide the new `TypeMap()` method.

- Update the owned sibling interoperability harness to the published Log
  `/v2` API and align the companion catalog. This is a harness adoption,
  not a production correlation logging migration.

## 1.1.2 - 2026-10-01

### Changed

- Update OpenTelemetry API and trace dependencies in both modules, and
  metric dependencies in the sibling harness, to v1.46.0 while retaining
  the harness SDK components at v1.45.0.

- Keep reusable CI and its checked-out tooling on the same v1.8.5 source
  while retaining the checksum-verified v1.6.1 CLI bootstrap.

### Security

- Upgrade OpenTelemetry dependencies to v1.45.0 in both the root module and
  sibling integration harness, including the SDK fix for exporter
  configuration logging that could disclose endpoint URLs
  (GHSA-8wmf-6v46-5gfg).

## 1.1.1 - 2026-09-30

### Security

- Use the published `go-identifier/v2` UUIDv4 generator for default hop IDs.
  Its classified entropy failures omit caller-reader error text while the
  correlation API, canonical UUIDv4 output, and generator ownership remain
  unchanged. The published sibling integration harness retains its separate
  historical v1.0.0 dependency.

## 1.1.0 - 2026-09-09

### Added

- Add target-oriented `adapters/http`, `adapters/jsonrpc`, `adapters/queue`,
  `adapters/schedule`, `adapters/slog`, and `adapters/otel` packages.
- Add `Factory.Create`, `schedule.Adapter.Create`, and
  `schedule.Adapter.Receive` as the canonical value-creation and scheduled
  metadata-receipt operations.

### Changed

- Publish schema-v2 cohesion metadata for the root module and its transport,
  logging, and observability packages.
- Adopt the checksum-verified `go-library-tools` v1.6.1 CLI, local cohesion
  validation, and immutable proportional hosted-workflow enforcement.
- Reconcile the root and interoperability-harness `go-identifier` v1.0.0
  checksums and all bootstrap-shadowed harness dependency checksums with their
  immutable public proxy and checksum database identities.
- Adopt the versioned shared `golib` repository contract for local and hosted
  verification while retaining package-owned API and mutation evidence.
- Reject typed-nil custom generators during factory construction while keeping
  literal nil as the cryptographic-default selection.
- Keep `Factory.Start`, `schedule.Adapter.Start`, and `schedule.Adapter.Run` as
  exact compatibility aliases for the new canonical names.
- Retain the released `http`, `jsonrpc`, `queue`, `schedule`, `log`, and
  `telemetry` paths as deprecated compatibility implementations with unchanged
  exported type, reflection, sentinel, and runtime behavior.

### Documentation

- Add the stable-v1 install and support contract, compiler-checked quick-start
  navigation, complete lifecycle and concurrency ownership, operational limits,
  actionable troubleshooting, and a repository-bound documentation gate.
- Link consumers to the immutable v1.5.3 Golib ecosystem index and Foundations
  package-family guidance, and publish the internal sibling interoperability
  harness entry point for engineering navigation.
- Document the lifetime, synchronization, and ownership obligations for a
  caller-supplied generator retained by a factory.
- Remove completed implementation plans from the release tree and retain
  package-owned documentation as the maintained reference.

## 1.0.0 - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Fixed

- Amortize system entropy reads through a bounded factory-owned buffer while
  preserving cryptographic UUIDv4 request identities.
- Skip inbound carrier parsing when the default HTTP policy replaces all
  untrusted metadata, and write canonical response headers without reparsing
  their names.
- Trust the canonical default UUID generator while validating custom generator
  output once with a byte-oriented ASCII policy scan.
- Return fresh correlation and request identifiers on explicitly rejected
  malformed HTTP metadata without invoking application handlers.
- Avoid allocating header-value storage when optional inbound correlation
  metadata is absent and reuse canonical header storage when it is safe.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-correlation` identity while preserving its documented API and behavior.
- Replace the obsolete owned-module pseudo-version pin with the monorepo's
  local `v0.0.0` source-proxy coordinate; release tooling continues to emit
  the exact `v1.0.0` dependency version.
- Pin the owned identifier module to an immutable source revision so
  correlation resolves from a clean external consumer without `go.work`.

- Normalized standalone module metadata against the canonical owned dependency
  graph, including complete checksums for clean consumer resolution.

### Added

- Distinct correlation, request, causation, and external identifier types.
- Secure `identifier` generation and explicit deterministic strategies.
- Context, carrier, HTTP, JSON-RPC, queue, schedule, webhook, log, telemetry,
  and request ID middleware adapters.
- Trust, privacy, multi-hop, retry, fuzz, race, mutation, coverage, allocation,
  compatibility, documentation, and CI gates.
