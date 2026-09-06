# Compatibility

The stable root module is installed as
`github.com/faustbrian/go-correlation@v1`; root releases use `vX.Y.Z` tags.
The latest published release is `v1.0.0`. The module requires and is tested
with exactly Go 1.26.6. Public API compatibility is checked against
[`api/baseline.txt`](../api/baseline.txt).

The runtime dependencies are the published `identifier` module for secure
default generation and OpenTelemetry trace contracts for optional links. The
repository has no permanent module replacement and resolves as a standalone
consumer outside `go.work`.

Transport wire keys are stable but remain configurable for controlled
migrations. Additive APIs documented on `main` are not available to consumers
until they appear in a semantic-version release; select the published tag in
`go.mod`, not a branch name.

See the repository-wide [compatibility policy](../COMPATIBILITY.md) for the
stable-v1 contract and [changelog](../CHANGELOG.md) for release chronology.
