# Sibling interoperability harness

This internal, non-releasable module compiles and exercises the public
contracts between `correlation`, `http-middleware`, `log`, and `telemetry`. It
owns cross-repository evidence only; it does not own a production API,
transport policy, observability backend, or application runtime.

The harness belongs in the engineering inventory. It is not an installable
library, has no independent release, and must not appear in the consumer
catalog. Its supported toolchain is Go 1.26.6.

## What it proves

- A trusted request identifier created by `http-middleware/requestid` can be
  adopted explicitly through `correlation/http/requestidbridge` without
  sharing a private context key.
- Correlation attributes can be passed to `go-log` through standard `slog`
  attributes without adding a logger dependency to the correlation module.
- A correlation link can be attached to a caller-owned OpenTelemetry span and
  observed through the `go-telemetry` test harness.

The tests use public APIs from independently versioned modules. Production
correlation packages do not import this harness or its sibling libraries, so
the evidence does not add a runtime dependency or dependency cycle.

## Run the harness

From this directory, run the deterministic in-memory scenarios against the
repository's current root source through the intentional parent workspace:

```sh
GOTOOLCHAIN=go1.26.6 go test ./...
```

Disable workspace discovery to verify the exact published module versions in
this module's `go.mod` and `go.sum`:

```sh
GOTOOLCHAIN=go1.26.6 GOWORK=off go test ./...
```

From the repository root, `make check` includes the harness in the aggregate
repository contract. No network service, credentials, environment variables,
database, or global telemetry initialization are required.

## Ownership and lifecycle

The tests collectively construct and own middleware, logging, correlation,
and telemetry fixtures. The request context remains scoped to its HTTP
operation. The telemetry scenario shuts down its test harness during cleanup;
this module starts no independently owned background work and exposes no drain
or shutdown API.

Errors remain owned by their source modules. The harness checks returned
errors directly and does not define a shared error category, retry policy, or
fallback. Fixtures contain no secrets or customer data; future diagnostics
must keep identifier values and telemetry output bounded.

## When to use this module

Use this harness when changing request-ID adoption, correlation logging
attributes, correlation telemetry links, or their sibling dependency versions.
Do not import it into applications or treat its in-memory fixtures as
production configuration, operational assurance, or a compatibility-set
publication.

## References

- [Correlation overview](../../README.md)
- [Adapter map](../../docs/adapters.md)
- [Adoption guidance](../../docs/adoption.md)
- [Verification contract](../../docs/verification.md)
- [Security policy](../../SECURITY.md)
- [Support](../../SUPPORT.md)
- [Contributing](../../CONTRIBUTING.md)
- [Repository changelog](../../CHANGELOG.md)
- [License](../../LICENSE)
- [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.5.3/docs/ecosystem/README.md)
- [Foundations family guidance](https://github.com/faustbrian/go-library-tools/blob/v1.5.3/docs/ecosystem/design-language.md#package-families-and-selection)
