# Changelog

All notable changes follow Keep a Changelog. This project uses semantic
versioning once released.

## Unreleased

### Changed

- Publish schema-v2 cohesion metadata for the root module and its transport,
  logging, and observability packages.
- Adopt the checksum-verified `go-library-tools` v1.3.0 CLI, local cohesion
  validation, and immutable hosted-workflow enforcement.
- Adopt the versioned shared `golib` repository contract for local and hosted
  verification while retaining package-owned API and mutation evidence.

### Documentation

- Link consumers to the immutable v1.3.0 Golib ecosystem guidance and package
  selection contract.
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
