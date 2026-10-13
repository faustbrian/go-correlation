# Documentation

[Repository entry](../README.md) · [Versioned ecosystem](https://github.com/faustbrian/go-library-tools/blob/v1.5.3/docs/ecosystem/README.md)

Start with the [README quick start](../README.md#quick-start) and its
[compiler-checked example](../example_test.go). The root module is the sole
releasable module; its public packages are:

- [`correlation`](../README.md) for values, factories, codecs, trust policy,
  disclosure, context storage, and transport-neutral propagation;
- [`adapters/http`](../adapters/http/) and
  [`http/requestidbridge`](../http/requestidbridge/) for HTTP boundaries;
- [`adapters/jsonrpc`](../adapters/jsonrpc/),
  [`adapters/queue`](../adapters/queue/),
  [`adapters/schedule`](../adapters/schedule/), and
  [`webhook`](../webhook/) for protocol and delivery boundaries; and
- [`adapters/slog`](../adapters/slog/) and
  [`adapters/otel`](../adapters/otel/) for bounded observability projection.

The former `http`, `jsonrpc`, `queue`, `schedule`, `log`, and `telemetry` paths
remain deprecated compatibility implementations and identity authorities.

## Learn and adopt

- [Semantics](semantics.md)
- [Generation and generator ownership](generation.md)
- [Trust](trust.md)
- [Propagation matrix](propagation.md)
- [Adapters](adapters.md)
- [Adoption](adoption.md)
- [Migration from Cline Correlation](migration.md)

## Reference and operations

- [API](api.md) and [pkg.go.dev](https://pkg.go.dev/github.com/faustbrian/go-correlation)
- [Architecture](architecture.md)
- [Observability](observability.md) and [privacy](privacy.md)
- [Operations and performance](operations.md)
- [Troubleshooting](troubleshooting.md)
- [FAQ](faq.md)
- [Compatibility](compatibility.md)
- [Verification](verification.md) and [assurance](assurance.md)
- [Source documentation audit](source-documentation-audit.md)
- [Sibling interoperability harness](../integration/siblings/README.md)

## Project and ecosystem

- [Changelog](../CHANGELOG.md), [contributing](../CONTRIBUTING.md), and
  [license](../LICENSE)
- [Support](../SUPPORT.md) and [private security reporting](../SECURITY.md)
- [Versioned Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.5.3/docs/ecosystem/README.md)
- [Foundations package family](https://github.com/faustbrian/go-library-tools/blob/v1.5.3/docs/ecosystem/design-language.md#package-families-and-selection)

## Complete documentation map

- [Adapters](adapters.md) — `docs/adapters.md`
- [Adoption](adoption.md) — `docs/adoption.md`
- [API guide](api.md) — `docs/api.md`
- [Architecture](architecture.md) — `docs/architecture.md`
- [Assurance](assurance.md) — `docs/assurance.md`
- [Compatibility](compatibility.md) — `docs/compatibility.md`
- [FAQ](faq.md) — `docs/faq.md`
- [Generation](generation.md) — `docs/generation.md`
- [Migration from Cline Correlation](migration.md) — `docs/migration.md`
- [Tracing, logging, and metrics](observability.md) — `docs/observability.md`
- [Operations](operations.md) — `docs/operations.md`
- [Package map](package-map.md) — `docs/package-map.md`
- [Privacy](privacy.md) — `docs/privacy.md`
- [Propagation matrix](propagation.md) — `docs/propagation.md`
- [Semantics](semantics.md) — `docs/semantics.md`
- [Source documentation audit](source-documentation-audit.md) — `docs/source-documentation-audit.md`
- [Troubleshooting](troubleshooting.md) — `docs/troubleshooting.md`
- [Trust and precedence](trust.md) — `docs/trust.md`
- [Verification](verification.md) — `docs/verification.md`
