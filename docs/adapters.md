# Adapters

New integrations use the target-oriented `adapters/http`, `adapters/jsonrpc`,
`adapters/queue`, `adapters/schedule`, `adapters/slog`, and `adapters/otel`
packages. The former top-level paths remain deprecated compatibility
implementations and identity authorities for their released exported types,
reflection identities, sentinels, methods, and behavior.

HTTP uses `X-Correlation-ID`, `X-Request-ID`, and `X-Causation-ID`. The adapter
finds case-insensitive duplicates, sanitizes request headers to the accepted
values, mirrors them to the response, and stores immutable context values.

JSON-RPC operates on an explicit `Metadata` map of raw JSON values so a strict
envelope decoder can preserve duplicate members. It does not add a `meta`
member or otherwise rewrite protocol envelopes.

Queue and scheduler adapters use application-owned string maps. A scheduler
uses `Create` for an independent invocation and `Receive` for an explicitly
trusted or rejected metadata boundary; the legacy `Start` and `Run` methods
delegate to those names. The webhook adapter names the HTTP send/receive
boundary. The request ID bridge accepts a
bound lookup function; pass a closure around
`requestid.FromContext(ctx, requestid.Request)` after the middleware source is
explicitly trusted.
