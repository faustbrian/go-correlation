# Security policy

Report vulnerabilities through the repository's
[private vulnerability form](https://github.com/faustbrian/go-correlation/security/advisories/new).
Do not open a public issue for identifier spoofing, trust-boundary bypass,
business-data disclosure, generator collisions, context collision, or
unbounded input.

The stable v1 release line receives security fixes. The latest published
release is `v1.0.0`; unreleased commits on `main` are not a substitute for a
supported tag. Correlation values are diagnostic metadata and must never be
relied on for access control, tenancy, replay protection, or idempotency.
