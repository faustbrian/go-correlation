# Troubleshooting

## A trusted inbound correlation ID was replaced

Confirm that the immediate peer passed the adapter's explicit trust check and
that the configured field contains exactly one canonical ASCII value within
the policy length. Extraction alone never grants trust. With replacement
policy, malformed or untrusted metadata creates a fresh root; with rejection
policy, malformed HTTP metadata returns a bounded client error before the
application handler runs.

## Injection reports an overwrite or invalid carrier

Create a new outbound carrier or clear application-owned fields before sending.
Injection refuses an already-populated field for each non-empty identifier it
would write, including malformed existing content. A direct `Codec.Inject`
with partial `Values` neither inspects nor clears target fields corresponding
to empty identifiers, so clear any unrelated stale fields explicitly.

## Retries have different request IDs

This is expected. A retry or redelivery is a distinct attempt: correlation is
preserved, the prior request becomes causation, and a new request ID is
generated. If an application needs duplicate suppression, use the idempotency
owner rather than reusing request identity.

## Generation fails

Match `correlation.ErrGeneration` with `errors.Is`, then inspect the preserved
cause. With the default factory, investigate operating-system entropy. With a
custom generator, verify its lifetime, concurrency safety, returned alphabet,
and configured length. Failed two-value creation returns no partial `Values`.

## Identifiers are missing after an asynchronous boundary

Contexts do not cross queue, schedule, JSON-RPC, or webhook boundaries by
themselves. Queue and JSON-RPC `Send` mutate the supplied metadata map; transport
that map and call `Receive` with an explicit trust decision. Scheduled work uses
`Enqueue` and `Receive` with the same map ownership. HTTP and webhook `Inject`
write the outbound request headers, while receiving `Wrap` middleware extracts
and applies its configured trust callback. The returned `Values` describe the
new local hop; they are not the transported carrier. Do not copy a request ID
unchanged across attempts.

## Logs expose raw values or metrics have high cardinality

Use redacted disclosure by default. Keyed hashes require an explicit key and
raw disclosure requires the explicit exposure mode. Metrics should use only
`telemetry.MetricAttributes`, which emits fixed-cardinality presence flags.

For reproducible non-security defects, use
[GitHub Issues](https://github.com/faustbrian/go-correlation/issues/new/choose).
Use [GitHub Discussions](https://github.com/faustbrian/go-correlation/discussions)
for adoption questions and follow the [security policy](../SECURITY.md) for
private vulnerability reports.
