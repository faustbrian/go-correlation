# Compatibility

The stable root module is installed as
`github.com/faustbrian/go-correlation@v1`; root releases use `vX.Y.Z` tags.
The module requires Go 1.26.6 and is tested with Go 1.27.0. Public API
compatibility is checked against
[`api/baseline.txt`](../api/baseline.txt).

The runtime dependencies are the published `identifier` module for secure
default generation and OpenTelemetry trace contracts for optional links. The
repository has no permanent module replacement and resolves as a standalone
consumer outside `go.work`.

Transport wire keys are stable but remain configurable for controlled
migrations. Additive APIs documented on `main` are not available to consumers
until they appear in a semantic-version release; select the published tag in
`go.mod`, not a branch name.

The deprecated `http`, `jsonrpc`, `queue`, `schedule`, `log`, and `telemetry`
packages preserve the released API, named and reflection type identity,
sentinels, error traversal, defaults, lifecycle, ownership, concurrency, and
serialization behavior. For the four packages with named types, the
`adapters/...` successors alias the legacy identity authority and delegate to
the same implementation. The slog and OpenTelemetry successors own their
stateless implementations while the old functions delegate to them. The legacy
paths remain supported for the longer of 180 days after public successor
availability and two subsequently published stable minor releases; removal
requires a separately authorized major release.

See the repository-wide [compatibility policy](../COMPATIBILITY.md) for the
stable-v1 contract and [changelog](../CHANGELOG.md) for release chronology.
