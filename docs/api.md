# API guide

Use `NewFactory` and `Create` for a workflow, `Next` for an outbound child, and
`Accept` for an inbound boundary. `Factory.Start` remains a compatibility alias
for `Create`. Use `NewCodec` only for custom carriers;
transport packages already configure their standard field names. `NewExternalID`
records non-correlation identifiers without promoting them.

For scheduled work, use `schedule.Adapter.Create` to start an independent run
and `schedule.Adapter.Receive` to accept application-owned metadata under an
explicit trust decision. The released `Start` and `Run` names remain exact
compatibility aliases for those operations.

Adapters are selected through `adapters/http`, `adapters/jsonrpc`,
`adapters/queue`, `adapters/schedule`, `adapters/slog`, and `adapters/otel`.
The former top-level paths remain deprecated compatibility implementations and
identity authorities. The four transport successors alias the released named
types so values cross old and new imports without conversion.

Use `NewDeterministic` only after a privacy review. Use `Disclose` indirectly
through the log and telemetry packages unless implementing another bounded
observability adapter. The generated [API baseline](../api/baseline.txt) is the
authoritative exported-symbol inventory.
