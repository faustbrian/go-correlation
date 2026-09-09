# correlation

[![CI](https://github.com/faustbrian/go-correlation/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-correlation/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-correlation/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-correlation.svg)](https://pkg.go.dev/github.com/faustbrian/go-correlation)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-correlation?sort=semver)](https://github.com/faustbrian/go-correlation/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`correlation` is the transport-neutral owner of correlation, request, and
causation identifiers. It carries those semantics through HTTP, JSON-RPC,
queues, scheduled work, webhooks, logs, and OpenTelemetry without creating a
global propagator or redefining trace and idempotency concepts.

The module is a stable v1 library. It requires and is tested with Go 1.26.6.
Stable releases are selected through the `v1` module line.

## Install

```sh
go get github.com/faustbrian/go-correlation@v1
```

## Semantics

- `CorrelationID` groups one logical interaction or workflow.
- `RequestID` identifies exactly one transport hop or delivery attempt.
- `CausationID` identifies that hop's immediate parent request or message.

These types are deliberately distinct. Trace and span IDs remain owned by
OpenTelemetry. Idempotency keys and fingerprints remain owned by
`idempotency`. None of these values is authentication, authorization,
tenancy, uniqueness, replay protection, or idempotency evidence.

## Quick start

```go
factory, err := correlation.NewFactory(correlation.FactoryOptions{})
if err != nil {
    return err
}

root, err := factory.Start()
if err != nil {
    return err
}
child, err := factory.Next(root)
if err != nil {
    return err
}
```

The complete compiler-checked released-v1 version is the package
[`Example`](example_test.go). Further examples cover
[`Factory.Create`](example_test.go) and
[`Factory.Next`](example_test.go).

`Create` is the canonical operation for fresh root values. `Start` remains a
compatibility alias with identical values and errors. Scheduled work follows
the same boundary: `schedule.Adapter.Create` and `schedule.Adapter.Receive`
are canonical, while `Start` and `Run` remain compatible aliases.

The default factory uses an explicitly owned, bounded entropy buffer around the
cryptographic UUIDv4 generator from `identifier`. A caller-supplied generator
remains instance scoped and
must return canonical text accepted by the configured policy.

Inbound metadata is never trusted by extraction alone:

```go
inbound, err := codec.Extract(carrier)
if err != nil {
    return err
}
values, err := factory.Accept(inbound, correlation.InboundPolicy{
    TrustCorrelation: true,
    TrustRequestAsCausation: true,
})
```

The application must establish that trust from an authenticated immediate
transport boundary first. Every accepted hop receives a new request ID.

## Adapters

- [`adapters/http`](adapters/http/) sanitizes inbound headers, applies explicit proxy trust,
  installs immutable context values, and injects outbound hops.
- [`http/requestidbridge`](http/requestidbridge/) explicitly adopts a trusted
  `http-middleware/requestid` value without importing its private key.
- [`adapters/jsonrpc`](adapters/jsonrpc/) reads and writes a separate metadata object without
  altering JSON-RPC envelopes.
- [`adapters/queue`](adapters/queue/) preserves workflow identity while generating a distinct
  request ID for every retry or redelivery.
- [`adapters/schedule`](adapters/schedule/) starts independent runs unless metadata is
  deliberately propagated.
- [`webhook`](webhook/) gives outbound and inbound webhook hops HTTP semantics.
- [`adapters/slog`](adapters/slog/) supplies redacted, keyed-hash, or explicitly raw `slog` attrs.
- [`adapters/otel`](adapters/otel/) attaches attributes to telemetry-owned links and
  exposes only fixed-cardinality presence flags to metrics.

The former `http`, `jsonrpc`, `queue`, `schedule`, `log`, and `telemetry`
import paths remain deprecated compatibility implementations and identity
authorities. New code should use the target-oriented `adapters/...` paths.

W3C Trace Context and Baggage remain optional application-owned propagation.
They may be linked to these values, but correlation IDs never become trace or
span IDs.

## When to use it

Use this module when work crosses process, protocol, queue, retry, or scheduled
boundaries and operators need to relate those attempts without conflating
their identities. Do not use it as an authorization, tenancy, uniqueness,
deduplication, replay-protection, or idempotency mechanism, and do not replace
OpenTelemetry trace or span identity with correlation values.

## Construction, errors, and lifecycle

`NewFactory` and `NewCodec` validate their complete value configuration.
`NewPropagator` and transport constructors reject missing dependencies and
validate the policy options they own. The zero `Policy` accepts canonical
ASCII identifiers up to 128 bytes. Stable sentinel categories such as
`ErrInvalidID`, `ErrInvalidFactory`, and `ErrGeneration` are matched with
`errors.Is`; wrapped causes remain available where disclosure is safe.

Factory, codec, and adapter operations perform no network I/O, retry, or
detached work. `WithValues`, `FromContext`, and HTTP middleware use only the
caller-owned context passed to that call; they do not retain it. The module
starts no goroutines, opens no resources, and owns no `Close` or `Shutdown`
operation. Applications retain ownership of carriers, headers, maps, contexts,
telemetry providers, and supplied callbacks or generators.

The built-in factory and immutable strategies are safe for concurrent use.
A caller-supplied generator and HTTP `Trust` callback are borrowed for the
owning object's lifetime, so callers must keep them valid and provide any
synchronization they require. Extraction and receive operations only read
carriers. Send, inject, and HTTP wrap operations update the carrier or headers
supplied for that call; callers must not concurrently reuse those mutable
values without their own synchronization.

Shared construction, ownership, lifecycle, and composition expectations are in
the versioned [Golib ecosystem
index](https://github.com/faustbrian/go-library-tools/blob/v1.5.3/docs/ecosystem/README.md)
and its [foundations package
guidance](https://github.com/faustbrian/go-library-tools/blob/v1.5.3/docs/ecosystem/design-language.md#package-families-and-selection).

## Deterministic correlation

`NewDeterministic` is an explicit opt-in for stable business workflows. It uses
HMAC-SHA-256, a versioned domain, length-delimited input, and bounded output.
Use a secret key when the input is private or comes from a small input space.
Deterministic correlation is linkable and is never the factory default.

## Security defaults

Identifiers use the canonical `[A-Za-z0-9_-]` alphabet and a default maximum
of 128 bytes. Carriers reject empty, oversized, malformed, control-bearing,
Unicode, and conflicting values. Injection refuses a populated target field
for each non-empty identifier it would write. A direct partial `Codec.Inject`
leaves fields corresponding to empty identifiers unchanged, so callers must
clear any unrelated stale fields first. Observability output is redacted unless
disclosure is explicitly enabled, and metrics never contain identifier values.

## Verification

Run the local release-equivalent gate:

```sh
make inventory
make cohesion
make check
make ci
```

It verifies formatting, module tidiness, vet, unit and integration tests, the
race detector, 100% production statement coverage in every package, fuzz
smoke tests, mutation tests, allocation benchmarks, documentation, API
compatibility, linting, Staticcheck, vulnerability analysis, and NilAway.

See the [documentation index](docs/README.md), [security policy](SECURITY.md),
and [changelog](CHANGELOG.md). Operational limits and performance caveats are
in the [operations guide](docs/operations.md); symptom-driven recovery steps
are in [troubleshooting](docs/troubleshooting.md). Adoption questions belong in
[GitHub Discussions](https://github.com/faustbrian/go-correlation/discussions),
and reproducible defects belong in
[GitHub Issues](https://github.com/faustbrian/go-correlation/issues/new/choose).
