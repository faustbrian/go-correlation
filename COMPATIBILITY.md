# Compatibility Policy

The repository has one releasable root module,
`github.com/faustbrian/go-correlation`, and follows semantic versioning. Root
releases use `vX.Y.Z` tags. The module is stable at v1, requires Go 1.27.0,
and is installed with `go get github.com/faustbrian/go-correlation@v1`.

Before `v1`, minor releases MAY contain reviewed breaking changes, but every
break MUST be documented with migration guidance. Patch releases MUST remain
backward compatible. At and after `v1`, incompatible exported API or documented
behavior changes require a new major version.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. Deprecated APIs
follow [`DEPRECATION.md`](DEPRECATION.md).
