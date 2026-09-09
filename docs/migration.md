# Migration from Cline Correlation

Inventory the old middleware, header names, context accessors, logger fields,
queue metadata, and any business decisions made from correlation values.
Remove business decisions first; this package intentionally has no equivalent
authorization, tenant, or idempotency API.

Map the old workflow identifier to `CorrelationID`, create a new `RequestID`
at each boundary, and set `CausationID` from the immediate prior request. Do
not preserve a legacy request ID across retries. Install explicit adapters in
parallel, compare propagation in redacted logs, then remove the old global or
ambient accessor. Header aliases should be temporary codec configuration with
a documented retirement date.

For code already using this module, migrate new root creation from
`Factory.Start` to `Factory.Create`. Scheduled integrations should migrate
`Adapter.Start` to `Adapter.Create` and `Adapter.Run` to `Adapter.Receive`.
The old methods remain source-compatible aliases and preserve values, trust
handling, generator calls, errors, and nil-receiver behavior, so migration does
not require a coordinated cutover.

Migrate transport imports independently:

- `github.com/faustbrian/go-correlation/http` to
  `github.com/faustbrian/go-correlation/adapters/http`;
- `github.com/faustbrian/go-correlation/jsonrpc` to
  `github.com/faustbrian/go-correlation/adapters/jsonrpc`;
- `github.com/faustbrian/go-correlation/queue` to
  `github.com/faustbrian/go-correlation/adapters/queue`; and
- `github.com/faustbrian/go-correlation/schedule` to
  `github.com/faustbrian/go-correlation/adapters/schedule`;
- `github.com/faustbrian/go-correlation/log` to
  `github.com/faustbrian/go-correlation/adapters/slog`; and
- `github.com/faustbrian/go-correlation/telemetry` to
  `github.com/faustbrian/go-correlation/adapters/otel`.

`adapters/http` keeps the `httpcorrelation` package identifier. Unless an
explicit import alias is retained, update the other default identifiers from
`jsonrpc`, `queue`, `schedule`, `log`, and `telemetry` to
`correlationjsonrpc`, `correlationqueue`, `correlationschedule`,
`correlationslog`, and `correlationotel`. Exported transport type identities
remain the released identities, so those imports can migrate one boundary at a
time without value conversion. The legacy paths remain supported for the
longer of 180 days after public successor availability and two subsequently
published stable minor releases.
