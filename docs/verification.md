# Verification

`golib check --all` runs the full configured local gate set. This alone
does not prove release readiness: release metadata checks and clean module
resolution are exercised by the separate workflow dispatch with
`release_dry_run` enabled. Release coverage is measured per production
package and must be exactly 100.0%. Mutation tests invert trust,
fresh request generation, immediate causation, duplicate precedence, overwrite,
proxy trust, transport bounds, custom JSON-RPC field validation, deterministic
versioning, typed-nil generator rejection, root-creation ordering and failure
atomicity, legacy-method delegation, and redaction decisions; every mutant must
be killed.

The race detector covers factory, context, and transport tests. Fuzz smoke
tests exercise typed parsing, carrier extraction, HTTP headers, and raw JSON-RPC
metadata. Allocation benchmarks report hop generation, carrier round trips, and
bounded oversized-value rejection. API and documentation checks are
deterministic checked-in gates. Ordinary hosted CI exercises its selected
module contract; it is not evidence that the separate release rehearsal has
passed.

The sibling integration module compiles and executes the request ID bridge,
`log` attributes, and `telemetry` links against the published versions in its
module manifest when workspace discovery is disabled, as in hosted CI. This
does not inject the candidate root source or prove current fleet composition.
The intentional parent workspace selects the local correlation root during
default local harness runs; the other sibling libraries remain published
dependencies.

NilAway is advisory during full checks; its diagnostics do not fail the gate.
Ordinary hosted local-mode checks do not run it.
